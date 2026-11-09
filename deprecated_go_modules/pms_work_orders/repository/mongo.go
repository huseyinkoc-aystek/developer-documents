package repository

import (
	"context"
	"strings"
	"time"

	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/model"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	workOrdersColl *mongo.Collection
	eventsColl     *mongo.Collection
	outboxColl     *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		workOrdersColl: db.Collection("pms_work_orders"),
		eventsColl:     db.Collection("pms_work_order_events"),
		outboxColl:     db.Collection("pms_outbox"),
	}
}

func (r *mongoRepository) InsertWorkOrder(ctx context.Context, item model.WorkOrder) (primitive.ObjectID, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return primitive.NilObjectID, err
	}
	now := time.Now()
	item.ID = primitive.NewObjectID()
	item.TenantID = tenantID
	item.SiteID = siteID
	item.CreatedAt = now
	item.UpdatedAt = now
	_, err = r.workOrdersColl.InsertOne(ctx, item)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return item.ID, nil
}

func (r *mongoRepository) FindWorkOrderByID(ctx context.Context, id primitive.ObjectID) (*model.WorkOrder, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	var item model.WorkOrder
	if err := r.workOrdersColl.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil}).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *mongoRepository) ListWorkOrders(ctx context.Context, filter ListFilter) ([]model.WorkOrder, int64, error) {
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
	if filter.PropertyID != nil && *filter.PropertyID != primitive.NilObjectID {
		q["property_id"] = *filter.PropertyID
	}
	if filter.UnitID != nil && *filter.UnitID != primitive.NilObjectID {
		q["unit_id"] = *filter.UnitID
	}
	if filter.Status != "" {
		q["status"] = filter.Status
	}
	if filter.Priority != "" {
		q["priority"] = filter.Priority
	}
	if strings.TrimSpace(filter.Q) != "" {
		pattern := primitive.Regex{Pattern: strings.TrimSpace(filter.Q), Options: "i"}
		q["$or"] = bson.A{
			bson.M{"title": pattern},
			bson.M{"description": pattern},
			bson.M{"reporter_ref": pattern},
		}
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip((filter.Page - 1) * filter.Limit).SetLimit(filter.Limit)
	cur, err := r.workOrdersColl.Find(ctx, q, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)
	items := make([]model.WorkOrder, 0)
	if err := cur.All(ctx, &items); err != nil {
		return nil, 0, err
	}
	total, err := r.workOrdersColl.CountDocuments(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *mongoRepository) UpdateWorkOrder(ctx context.Context, id primitive.ObjectID, item model.WorkOrder) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	res, err := r.workOrdersColl.UpdateOne(ctx,
		bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil},
		bson.M{"$set": bson.M{
			"title":       item.Title,
			"description": item.Description,
			"category":    item.Category,
			"priority":    item.Priority,
			"tags":        item.Tags,
			"sla_due_at":  item.SLADueAt,
			"attachments": item.Attachments,
			"cost_items":  item.CostItems,
			"updated_at":  time.Now(),
			"updated_by":  item.UpdatedBy,
		}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *mongoRepository) UpdateAssignee(ctx context.Context, id primitive.ObjectID, assignee *primitive.ObjectID, actor *primitive.ObjectID) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	res, err := r.workOrdersColl.UpdateOne(ctx,
		bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil},
		bson.M{"$set": bson.M{"assignee_user_id": assignee, "updated_at": time.Now(), "updated_by": actor}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *mongoRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status model.WorkOrderStatus, closedAt *time.Time, actor *primitive.ObjectID) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	set := bson.M{"status": status, "updated_at": time.Now(), "updated_by": actor}
	if closedAt != nil {
		set["closed_at"] = *closedAt
	}
	res, err := r.workOrdersColl.UpdateOne(ctx,
		bson.M{"_id": id, "tenant_id": tenantID, "site_id": siteID, "deleted_at": nil},
		bson.M{"$set": set},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *mongoRepository) InsertEvent(ctx context.Context, event model.WorkOrderEvent) error {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return err
	}
	event.ID = primitive.NewObjectID()
	event.TenantID = tenantID
	event.SiteID = siteID
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	_, err = r.eventsColl.InsertOne(ctx, event)
	return err
}

func (r *mongoRepository) ListEvents(ctx context.Context, workOrderID primitive.ObjectID, limit int64) ([]model.WorkOrderEvent, error) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	cur, err := r.eventsColl.Find(ctx,
		bson.M{"tenant_id": tenantID, "site_id": siteID, "work_order_id": workOrderID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	items := make([]model.WorkOrderEvent, 0)
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *mongoRepository) InsertOutboxEvent(ctx context.Context, event model.OutboxEvent) error {
	now := time.Now()
	if event.ID.IsZero() {
		event.ID = primitive.NewObjectID()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = now
	}
	_, err := r.outboxColl.InsertOne(ctx, event)
	return err
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
