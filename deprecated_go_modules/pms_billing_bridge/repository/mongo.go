package repository

import (
	"context"
	"errors"
	"time"

	paymentsmodel "github.com/aystek/aystek-core/internal/modules/payments/model"
	pmsmodel "github.com/aystek/aystek-core/internal/modules/pms/model"
	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	propertiesColl    *mongo.Collection
	occupanciesColl   *mongo.Collection
	duesPlansColl     *mongo.Collection
	unitsColl         *mongo.Collection
	chargesColl       *mongo.Collection
	allocationsColl   *mongo.Collection
	idempotencyColl   *mongo.Collection
	paymentsColl      *mongo.Collection
	outboxColl        *mongo.Collection
	recurringPlanColl *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		propertiesColl:    db.Collection("pms_properties"),
		occupanciesColl:   db.Collection("pms_occupancies"),
		duesPlansColl:     db.Collection("pms_dues_plans"),
		unitsColl:         db.Collection("units"),
		chargesColl:       db.Collection("pms_charges"),
		allocationsColl:   db.Collection("pms_payment_allocations"),
		idempotencyColl:   db.Collection("pms_billing_idempotency"),
		paymentsColl:      db.Collection("payments"),
		outboxColl:        db.Collection("pms_outbox"),
		recurringPlanColl: db.Collection("recurring_plans"),
	}
}

func (r *mongoRepository) FindDuesPlanByID(ctx context.Context, id primitive.ObjectID) (*model.DuesPlanReference, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var item model.DuesPlanReference
	err = r.duesPlansColl.FindOne(ctx, bson.M{
		"_id":       id,
		"tenant_id": tenantID,
		"site_id":   siteID,
		"deleted_at": nil,
	}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mongoRepository) FindActiveDuesPlan(ctx context.Context, propertyID primitive.ObjectID, at time.Time) (*model.DuesPlanReference, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var item model.DuesPlanReference
	err = r.duesPlansColl.FindOne(ctx, bson.M{
		"tenant_id": tenantID,
		"site_id": siteID,
		"property_id": propertyID,
		"status": "active",
		"effective_from": bson.M{"$lte": at},
		"deleted_at": nil,
	}, options.FindOne().SetSort(bson.D{{Key: "effective_from", Value: -1}, {Key: "updated_at", Value: -1}})).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mongoRepository) FindUnitAreaSQM(ctx context.Context, unitID primitive.ObjectID) (float64, error) {
	tenantID, _, err := tenantSiteIDs(ctx)
	if err != nil {
		return 0, err
	}
	var raw bson.M
	if err := r.unitsColl.FindOne(ctx, bson.M{
		"_id":       unitID,
		"tenant_id": tenantID,
		"deleted_at": bson.M{"$exists": false},
	}, options.FindOne().SetProjection(bson.M{"area_sqm": 1, "sqm": 1, "area": 1})).Decode(&raw); err != nil {
		return 0, err
	}
	for _, key := range []string{"area_sqm", "sqm", "area"} {
		if v, ok := raw[key]; ok {
			switch t := v.(type) {
			case float64:
				return t, nil
			case float32:
				return float64(t), nil
			case int32:
				return float64(t), nil
			case int64:
				return float64(t), nil
			case int:
				return float64(t), nil
			}
		}
	}
	return 0, nil
}

func (r *mongoRepository) FindProperty(ctx context.Context, id primitive.ObjectID) (*pmsmodel.Property, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var item pmsmodel.Property
	if err := r.propertiesColl.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil}).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mongoRepository) FindActiveOccupanciesByPropertyPeriod(ctx context.Context, propertyID primitive.ObjectID, from, to time.Time) ([]pmsmodel.Occupancy, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"tenant_id":   tenantID,
		"site_id":     siteID,
		"property_id": propertyID,
		"deleted_at":  nil,
		"start_at":    bson.M{"$lte": to},
		"$or": bson.A{
			bson.M{"end_at": nil},
			bson.M{"end_at": bson.M{"$gte": from}},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "unit_id", Value: 1}, {Key: "start_at", Value: 1}})
	cur, err := r.occupanciesColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	items := make([]pmsmodel.Occupancy, 0)
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *mongoRepository) FindChargeByUnitPeriod(ctx context.Context, propertyID, unitID primitive.ObjectID, period string) (*model.PMSCharge, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var item model.PMSCharge
	err = r.chargesColl.FindOne(ctx, bson.M{
		"tenant_id":   tenantID,
		"site_id":     siteID,
		"property_id": propertyID,
		"unit_id":     unitID,
		"period":      period,
		"deleted_at":  nil,
	}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mongoRepository) InsertCharge(ctx context.Context, charge model.PMSCharge) (primitive.ObjectID, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return primitive.NilObjectID, err
	}
	now := time.Now()
	charge.ID = primitive.NewObjectID()
	charge.TenantID = tenantID
	charge.SiteID = siteID
	charge.CreatedAt = now
	charge.UpdatedAt = now
	_, err = r.chargesColl.InsertOne(ctx, charge)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return charge.ID, nil
}

func (r *mongoRepository) UpdateChargeAmounts(ctx context.Context, id primitive.ObjectID, paid, remaining float64, status model.ChargeStatus, actor *primitive.ObjectID) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	res, err := r.chargesColl.UpdateOne(ctx,
		bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil},
		bson.M{"$set": bson.M{"amount_paid": paid, "amount_remaining": remaining, "status": status, "updated_at": time.Now(), "updated_by": actor}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *mongoRepository) ListOpenChargesByParty(ctx context.Context, partyRef string, propertyID *primitive.ObjectID) ([]model.PMSCharge, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"tenant_id":        tenantID,
		"site_id":          siteID,
		"party_ref":        partyRef,
		"deleted_at":       nil,
		"amount_remaining": bson.M{"$gt": 0},
	}
	if propertyID != nil && *propertyID != primitive.NilObjectID {
		filter["property_id"] = *propertyID
	}
	opts := options.Find().SetSort(bson.D{{Key: "period", Value: 1}, {Key: "due_date", Value: 1}, {Key: "created_at", Value: 1}})
	cur, err := r.chargesColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	items := make([]model.PMSCharge, 0)
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *mongoRepository) ListCharges(ctx context.Context, propertyID *primitive.ObjectID, period string) ([]model.PMSCharge, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"tenant_id": tenantID, "site_id": siteID, "deleted_at": nil}
	if propertyID != nil && *propertyID != primitive.NilObjectID {
		filter["property_id"] = *propertyID
	}
	if period != "" {
		filter["period"] = period
	}
	opts := options.Find().SetSort(bson.D{{Key: "unit_id", Value: 1}, {Key: "period", Value: -1}})
	cur, err := r.chargesColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	items := make([]model.PMSCharge, 0)
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *mongoRepository) InsertAllocation(ctx context.Context, allocation model.PaymentAllocation) (primitive.ObjectID, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return primitive.NilObjectID, err
	}
	now := time.Now()
	allocation.ID = primitive.NewObjectID()
	allocation.TenantID = tenantID
	allocation.SiteID = siteID
	allocation.CreatedAt = now
	_, err = r.allocationsColl.InsertOne(ctx, allocation)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return allocation.ID, nil
}

func (r *mongoRepository) FindPaymentByID(ctx context.Context, id primitive.ObjectID) (*PaymentRef, error) {
	tenantID, _, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var pay paymentsmodel.Payment
	if err := r.paymentsColl.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID, "deleted_at": nil}).Decode(&pay); err != nil {
		return nil, err
	}
	payerRef := pay.EntityID
	if payerRef == "" {
		return nil, errors.New("payer_ref missing in payment")
	}
	return &PaymentRef{ID: pay.ID, Amount: pay.Amount, Currency: pay.Currency, PayerRef: payerRef, ReceivedAt: pay.ReceivedAt}, nil
}

func (r *mongoRepository) FindIdempotency(ctx context.Context, operation, key string) (*model.BillingIdempotency, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var rec model.BillingIdempotency
	err = r.idempotencyColl.FindOne(ctx, bson.M{
		"tenant_id":       tenantID,
		"site_id":         siteID,
		"operation":       operation,
		"idempotency_key": key,
	}).Decode(&rec)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r *mongoRepository) InsertIdempotency(ctx context.Context, rec model.BillingIdempotency) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	rec.ID = primitive.NewObjectID()
	rec.TenantID = tenantID
	rec.SiteID = siteID
	rec.CreatedAt = time.Now()
	_, err = r.idempotencyColl.InsertOne(ctx, rec)
	return err
}

func (r *mongoRepository) InsertOutboxEvent(ctx context.Context, event model.OutboxEvent) error {
	now := time.Now()
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = now
	}
	if event.ID.IsZero() {
		event.ID = primitive.NewObjectID()
	}
	_, err := r.outboxColl.InsertOne(ctx, event)
	return err
}

var _ Repository = (*mongoRepository)(nil)
