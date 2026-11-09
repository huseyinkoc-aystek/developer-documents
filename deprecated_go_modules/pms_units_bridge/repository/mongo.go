package repository

import (
	"context"
	"strings"
	"time"

	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/model"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	unitsColl       *mongo.Collection
	occupanciesColl *mongo.Collection
	unitSiteIndexColl *mongo.Collection
	propertiesColl  *mongo.Collection
	chargesColl     *mongo.Collection
	allocationsColl *mongo.Collection
	workOrdersColl  *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		unitsColl:       db.Collection("units"),
		occupanciesColl: db.Collection("pms_occupancies"),
		unitSiteIndexColl: db.Collection("pms_unit_site_index"),
		propertiesColl:  db.Collection("pms_properties"),
		chargesColl:     db.Collection("pms_charges"),
		allocationsColl: db.Collection("pms_payment_allocations"),
		workOrdersColl:  db.Collection("pms_work_orders"),
	}
}

type unitSiteIndexDoc struct {
	UnitID          primitive.ObjectID  `bson:"unit_id"`
	PropertyID      primitive.ObjectID  `bson:"property_id"`
	LastOccupancyID *primitive.ObjectID `bson:"last_occupancy_id,omitempty"`
}

func (r *mongoRepository) ListUnits(ctx context.Context, filter ListFilter) ([]model.UnitListItem, int64, error) {
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

	indexFilter := bson.M{
		"tenant_id":  tenantID,
		"site_id":    siteID,
	}
	if filter.PropertyID != nil && *filter.PropertyID != primitive.NilObjectID {
		indexFilter["property_id"] = *filter.PropertyID
	}
	idxCur, err := r.unitSiteIndexColl.Find(ctx, indexFilter)
	if err != nil {
		return nil, 0, err
	}
	defer idxCur.Close(ctx)
	indexRows := make([]unitSiteIndexDoc, 0)
	if err := idxCur.All(ctx, &indexRows); err != nil {
		return nil, 0, err
	}
	if len(indexRows) == 0 {
		return []model.UnitListItem{}, 0, nil
	}
	unitIDs := make([]primitive.ObjectID, 0, len(indexRows))
	indexByUnit := make(map[primitive.ObjectID]unitSiteIndexDoc, len(indexRows))
	for _, row := range indexRows {
		if row.UnitID == primitive.NilObjectID {
			continue
		}
		unitIDs = append(unitIDs, row.UnitID)
		indexByUnit[row.UnitID] = row
	}
	if len(unitIDs) == 0 {
		return []model.UnitListItem{}, 0, nil
	}

	unitFilter := bson.M{
		"tenant_id": tenantID,
		"_id":       bson.M{"$in": unitIDs},
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if strings.TrimSpace(filter.Q) != "" {
		pattern := primitive.Regex{Pattern: strings.TrimSpace(filter.Q), Options: "i"}
		unitFilter["$or"] = bson.A{
			bson.M{"block": pattern},
			bson.M{"number": pattern},
			bson.M{"floor": pattern},
		}
	}
	total, err := r.unitsColl.CountDocuments(ctx, unitFilter)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []model.UnitListItem{}, 0, nil
	}

	opts := options.Find().
		SetSkip((filter.Page - 1) * filter.Limit).
		SetLimit(filter.Limit).
		SetSort(bson.D{{Key: "block", Value: 1}, {Key: "number", Value: 1}})
	cur, err := r.unitsColl.Find(ctx, unitFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	type unitDoc struct {
		ID     primitive.ObjectID `bson:"_id"`
		Block  string             `bson:"block"`
		Number string             `bson:"number"`
		Floor  string             `bson:"floor"`
		Status string             `bson:"status"`
	}
	var docs []unitDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, 0, err
	}
	items := make([]model.UnitListItem, 0, len(docs))
	propertyNames := map[primitive.ObjectID]string{}
	for _, unit := range docs {
		item, err := r.buildUnitCard(ctx, tenantID, siteID, indexByUnit[unit.ID], unit.ID, unit.Block, unit.Number, unit.Floor, unit.Status, propertyNames)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, nil
}

func (r *mongoRepository) GetUnitDetail(ctx context.Context, id primitive.ObjectID) (*model.UnitDetail, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	type unitDoc struct {
		ID     primitive.ObjectID `bson:"_id"`
		Block  string             `bson:"block"`
		Number string             `bson:"number"`
		Floor  string             `bson:"floor"`
		Status string             `bson:"status"`
	}
	var unit unitDoc
	if err := r.unitsColl.FindOne(ctx, bson.M{
		"_id":       id,
		"tenant_id": tenantID,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}).Decode(&unit); err != nil {
		return nil, err
	}

	var idx unitSiteIndexDoc
	if err := r.unitSiteIndexColl.FindOne(ctx, bson.M{"tenant_id": tenantID, "site_id": siteID, "unit_id": id}).Decode(&idx); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}

	card, err := r.buildUnitCard(ctx, tenantID, siteID, idx, unit.ID, unit.Block, unit.Number, unit.Floor, unit.Status, map[primitive.ObjectID]string{})
	if err != nil {
		return nil, err
	}
	lastWorkOrders, err := r.lastWorkOrders(ctx, tenantID, siteID, id, 5)
	if err != nil {
		return nil, err
	}
	return &model.UnitDetail{
		UnitListItem:   *card,
		LastWorkOrders: lastWorkOrders,
	}, nil
}

func (r *mongoRepository) buildUnitCard(
	ctx context.Context,
	tenantID, siteID string,
	indexRow unitSiteIndexDoc,
	unitID primitive.ObjectID,
	block, number, floor, status string,
	propertyNames map[primitive.ObjectID]string,
) (*model.UnitListItem, error) {
	snapshot, propertyID, err := r.latestOccupancySnapshot(ctx, tenantID, siteID, indexRow)
	if err != nil {
		return nil, err
	}
	propertyName := ""
	if propertyID != primitive.NilObjectID {
		if cached, ok := propertyNames[propertyID]; ok {
			propertyName = cached
		} else {
			name, loadErr := r.propertyName(ctx, tenantID, siteID, propertyID)
			if loadErr != nil {
				return nil, loadErr
			}
			propertyName = name
			propertyNames[propertyID] = name
		}
	}
	dues, err := r.duesSummary(ctx, tenantID, siteID, unitID)
	if err != nil {
		return nil, err
	}
	lastPayment, err := r.lastPayment(ctx, tenantID, siteID, unitID)
	if err != nil {
		return nil, err
	}
	openCount, err := r.openWorkOrderCount(ctx, tenantID, siteID, unitID)
	if err != nil {
		return nil, err
	}
	return &model.UnitListItem{
		UnitID:              unitID,
		PropertyID:          propertyID,
		PropertyName:        propertyName,
		Block:               block,
		Number:              number,
		Floor:               floor,
		Status:              status,
		OccupancySnapshot:   snapshot,
		DuesSummary:         dues,
		LastPayment:         lastPayment,
		OpenWorkOrdersCount: openCount,
	}, nil
}

func (r *mongoRepository) latestOccupancySnapshot(ctx context.Context, tenantID, siteID string, indexRow unitSiteIndexDoc) (model.OccupancySnapshot, primitive.ObjectID, error) {
	if indexRow.LastOccupancyID == nil || *indexRow.LastOccupancyID == primitive.NilObjectID {
		return model.OccupancySnapshot{}, indexRow.PropertyID, nil
	}
	type occupancyDoc struct {
		ID             primitive.ObjectID  `bson:"_id"`
		PropertyID     primitive.ObjectID  `bson:"property_id"`
		PartyRef       string              `bson:"party_ref"`
		OccupantType   string              `bson:"occupant_type"`
		Status         string              `bson:"status"`
		StartAt        *primitive.DateTime `bson:"start_at"`
		EndAt          *primitive.DateTime `bson:"end_at"`
		Responsibility struct {
			FinancialResponsible bool     `bson:"financial_responsible"`
			SharePercent         *float64 `bson:"share_percent"`
		} `bson:"responsibility"`
	}
	var occ occupancyDoc
	err := r.occupanciesColl.FindOne(ctx, bson.M{
		"_id":         *indexRow.LastOccupancyID,
		"tenant_id":   tenantID,
		"site_id":     siteID,
		"deleted_at":  nil,
	}).Decode(&occ)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.OccupancySnapshot{}, indexRow.PropertyID, nil
		}
		return model.OccupancySnapshot{}, indexRow.PropertyID, err
	}
	var startAt *time.Time
	var endAt *time.Time
	if occ.StartAt != nil {
		v := occ.StartAt.Time()
		startAt = &v
	}
	if occ.EndAt != nil {
		v := occ.EndAt.Time()
		endAt = &v
	}
	return model.OccupancySnapshot{
		OccupancyID:           occ.ID,
		PartyRef:              occ.PartyRef,
		OccupantType:          occ.OccupantType,
		Status:                occ.Status,
		StartAt:               startAt,
		EndAt:                 endAt,
		FinancialResponsible:  occ.Responsibility.FinancialResponsible,
		FinancialSharePercent: occ.Responsibility.SharePercent,
	}, occ.PropertyID, nil
}

func (r *mongoRepository) propertyName(ctx context.Context, tenantID, siteID string, propertyID primitive.ObjectID) (string, error) {
	type propertyDoc struct {
		Name string `bson:"name"`
	}
	var prop propertyDoc
	err := r.propertiesColl.FindOne(ctx, bson.M{
		"_id":        propertyID,
		"tenant_id":  tenantID,
		"site_id":    siteID,
		"deleted_at": nil,
	}).Decode(&prop)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}
	return prop.Name, nil
}

func (r *mongoRepository) duesSummary(ctx context.Context, tenantID, siteID string, unitID primitive.ObjectID) (model.DuesSummary, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"tenant_id":  tenantID,
			"site_id":    siteID,
			"unit_id":    unitID,
			"deleted_at": nil,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":             nil,
			"total_due":       bson.M{"$sum": "$amount_due"},
			"total_paid":      bson.M{"$sum": "$amount_paid"},
			"total_remaining": bson.M{"$sum": "$amount_remaining"},
			"charge_count":    bson.M{"$sum": 1},
			"currency":        bson.M{"$first": "$currency"},
		}}},
	}
	cur, err := r.chargesColl.Aggregate(ctx, pipeline)
	if err != nil {
		return model.DuesSummary{}, err
	}
	defer cur.Close(ctx)
	type agg struct {
		TotalDue       float64 `bson:"total_due"`
		TotalPaid      float64 `bson:"total_paid"`
		TotalRemaining float64 `bson:"total_remaining"`
		ChargeCount    int64   `bson:"charge_count"`
		Currency       string  `bson:"currency"`
	}
	var out []agg
	if err := cur.All(ctx, &out); err != nil {
		return model.DuesSummary{}, err
	}
	if len(out) == 0 {
		return model.DuesSummary{}, nil
	}
	return model.DuesSummary{
		TotalDue:       out[0].TotalDue,
		TotalPaid:      out[0].TotalPaid,
		TotalRemaining: out[0].TotalRemaining,
		ChargeCount:    out[0].ChargeCount,
		Currency:       out[0].Currency,
	}, nil
}

func (r *mongoRepository) lastPayment(ctx context.Context, tenantID, siteID string, unitID primitive.ObjectID) (model.PaymentSummary, error) {
	type paymentDoc struct {
		PaymentRef      string     `bson:"payment_ref"`
		AllocatedAmount float64    `bson:"allocated_amount"`
		Currency        string     `bson:"currency"`
		AllocatedAt     *time.Time `bson:"allocated_at"`
	}
	var row paymentDoc
	err := r.allocationsColl.FindOne(ctx, bson.M{
		"tenant_id": tenantID,
		"site_id":   siteID,
		"unit_id":   unitID,
	}, options.FindOne().SetSort(bson.D{{Key: "allocated_at", Value: -1}, {Key: "created_at", Value: -1}})).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.PaymentSummary{}, nil
		}
		return model.PaymentSummary{}, err
	}
	return model.PaymentSummary{
		PaymentRef:      row.PaymentRef,
		AllocatedAmount: row.AllocatedAmount,
		Currency:        row.Currency,
		AllocatedAt:     row.AllocatedAt,
	}, nil
}

func (r *mongoRepository) openWorkOrderCount(ctx context.Context, tenantID, siteID string, unitID primitive.ObjectID) (int64, error) {
	return r.workOrdersColl.CountDocuments(ctx, bson.M{
		"tenant_id": tenantID,
		"site_id":   siteID,
		"unit_id":   unitID,
		"status": bson.M{
			"$in": bson.A{"open", "in_progress", "blocked"},
		},
		"deleted_at": nil,
	})
}

func (r *mongoRepository) lastWorkOrders(ctx context.Context, tenantID, siteID string, unitID primitive.ObjectID, limit int64) ([]model.WorkOrderMini, error) {
	if limit <= 0 {
		limit = 5
	}
	cur, err := r.workOrdersColl.Find(ctx, bson.M{
		"tenant_id":  tenantID,
		"site_id":    siteID,
		"unit_id":    unitID,
		"deleted_at": nil,
	}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	type workOrderDoc struct {
		ID        primitive.ObjectID `bson:"_id"`
		Title     string             `bson:"title"`
		Status    string             `bson:"status"`
		Priority  string             `bson:"priority"`
		CreatedAt primitive.DateTime `bson:"created_at"`
	}
	docs := make([]workOrderDoc, 0)
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	items := make([]model.WorkOrderMini, 0, len(docs))
	for _, doc := range docs {
		items = append(items, model.WorkOrderMini{
			ID:        doc.ID,
			Title:     doc.Title,
			Status:    doc.Status,
			Priority:  doc.Priority,
			CreatedAt: doc.CreatedAt.Time(),
		})
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
