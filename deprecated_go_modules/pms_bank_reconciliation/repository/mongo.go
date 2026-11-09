package repository

import (
	"context"
	"strings"
	"time"

	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/model"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	txnColl    *mongo.Collection
	chargeColl *mongo.Collection
	outboxColl *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		txnColl:    db.Collection("pms_bank_transactions"),
		chargeColl: db.Collection("pms_charges"),
		outboxColl: db.Collection("pms_outbox"),
	}
}

func (r *mongoRepository) InsertTransaction(ctx context.Context, txn model.BankTransaction) (primitive.ObjectID, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return primitive.NilObjectID, err
	}
	now := time.Now()
	txn.ID = primitive.NewObjectID()
	txn.TenantID = tenantID
	txn.SiteID = siteID
	txn.CreatedAt = now
	txn.UpdatedAt = now
	if txn.Status == "" {
		txn.Status = model.TxnStatusNew
	}
	_, err = r.txnColl.InsertOne(ctx, txn)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return txn.ID, nil
}

func (r *mongoRepository) FindTransactionByID(ctx context.Context, id primitive.ObjectID) (*model.BankTransaction, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var item model.BankTransaction
	if err := r.txnColl.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil}).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mongoRepository) FindTransactionByExternalID(ctx context.Context, externalTxnID string) (*model.BankTransaction, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var item model.BankTransaction
	err = r.txnColl.FindOne(ctx, bson.M{
		"tenant_id":       tenantID,
		"site_id":         siteID,
		"external_txn_id": strings.TrimSpace(externalTxnID),
		"deleted_at":      nil,
	}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mongoRepository) ListTransactions(ctx context.Context, filter ListFilter) ([]model.BankTransaction, int64, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, 0, err
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	q := bson.M{"tenant_id": tenantID, "site_id": siteID, "deleted_at": nil}
	if filter.Status != "" {
		q["status"] = filter.Status
	}
	if strings.TrimSpace(filter.Q) != "" {
		pattern := primitive.Regex{Pattern: strings.TrimSpace(filter.Q), Options: "i"}
		q["$or"] = bson.A{
			bson.M{"external_txn_id": pattern},
			bson.M{"description": pattern},
			bson.M{"counterparty": pattern},
		}
	}
	opts := options.Find().SetSort(bson.D{{Key: "txn_date", Value: -1}, {Key: "created_at", Value: -1}}).SetSkip((filter.Page - 1) * filter.Limit).SetLimit(filter.Limit)
	cur, err := r.txnColl.Find(ctx, q, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)
	items := make([]model.BankTransaction, 0)
	if err := cur.All(ctx, &items); err != nil {
		return nil, 0, err
	}
	total, err := r.txnColl.CountDocuments(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *mongoRepository) UpdateMatch(ctx context.Context, id primitive.ObjectID, status model.TxnStatus, candidates []model.MatchCandidate, actor *primitive.ObjectID) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	res, err := r.txnColl.UpdateOne(ctx,
		bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil},
		bson.M{"$set": bson.M{"status": status, "match_candidates": candidates, "updated_at": time.Now(), "updated_by": actor}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *mongoRepository) ConfirmTransaction(ctx context.Context, id primitive.ObjectID, confirmedRef string, paymentID *primitive.ObjectID, actor *primitive.ObjectID) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	res, err := r.txnColl.UpdateOne(ctx,
		bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil},
		bson.M{"$set": bson.M{"status": model.TxnStatusConfirmed, "confirmed_match_ref": strings.TrimSpace(confirmedRef), "payment_id": paymentID, "updated_at": time.Now(), "updated_by": actor}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *mongoRepository) RejectTransaction(ctx context.Context, id primitive.ObjectID, actor *primitive.ObjectID) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	res, err := r.txnColl.UpdateOne(ctx,
		bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil},
		bson.M{"$set": bson.M{"status": model.TxnStatusRejected, "updated_at": time.Now(), "updated_by": actor}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
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

func (r *mongoRepository) ListOpenChargesByAmount(ctx context.Context, amount float64, currency string) ([]ChargeCandidate, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	q := bson.M{
		"tenant_id":        tenantID,
		"site_id":          siteID,
		"deleted_at":       nil,
		"amount_remaining": bson.M{"$gt": 0},
	}
	if amount > 0 {
		q["amount_remaining"] = bson.M{"$gte": amount * 0.90, "$lte": amount * 1.10}
	}
	if strings.TrimSpace(currency) != "" {
		q["currency"] = strings.ToUpper(strings.TrimSpace(currency))
	}
	opts := options.Find().SetSort(bson.D{{Key: "period", Value: 1}, {Key: "due_date", Value: 1}}).SetLimit(25)
	cur, err := r.chargeColl.Find(ctx, q, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	type chargeDoc struct {
		ID              primitive.ObjectID `bson:"_id"`
		UnitID          primitive.ObjectID `bson:"unit_id"`
		PartyRef        string             `bson:"party_ref"`
		AmountRemaining float64            `bson:"amount_remaining"`
		Currency        string             `bson:"currency"`
		Period          string             `bson:"period"`
	}
	items := make([]ChargeCandidate, 0)
	for cur.Next(ctx) {
		var doc chargeDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		items = append(items, ChargeCandidate{ChargeID: doc.ID, UnitID: doc.UnitID, PartyRef: doc.PartyRef, AmountRemaining: doc.AmountRemaining, Currency: doc.Currency, DescriptionHint: doc.Period})
	}
	return items, nil
}

func tenantSiteIDs(ctx context.Context) (string, string, error) {
	tenantID, err := tenantfilter.RequireTenantID(ctx)
	if err != nil {
		return "", "", err
	}
	siteID, err := tenantfilter.RequireSiteID(ctx)
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(tenantID), strings.TrimSpace(siteID), nil
}

var _ Repository = (*mongoRepository)(nil)
