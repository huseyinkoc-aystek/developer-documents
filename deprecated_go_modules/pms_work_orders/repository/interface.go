package repository

import (
	"context"
	"time"

	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListFilter struct {
	Page       int64
	Limit      int64
	PropertyID *primitive.ObjectID
	UnitID     *primitive.ObjectID
	Status     model.WorkOrderStatus
	Priority   model.WorkOrderPriority
	Q          string
}

type Repository interface {
	InsertWorkOrder(ctx context.Context, item model.WorkOrder) (primitive.ObjectID, error)
	FindWorkOrderByID(ctx context.Context, id primitive.ObjectID) (*model.WorkOrder, error)
	ListWorkOrders(ctx context.Context, filter ListFilter) ([]model.WorkOrder, int64, error)
	UpdateWorkOrder(ctx context.Context, id primitive.ObjectID, item model.WorkOrder) error
	UpdateAssignee(ctx context.Context, id primitive.ObjectID, assignee *primitive.ObjectID, actor *primitive.ObjectID) error
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status model.WorkOrderStatus, closedAt *time.Time, actor *primitive.ObjectID) error

	InsertEvent(ctx context.Context, event model.WorkOrderEvent) error
	ListEvents(ctx context.Context, workOrderID primitive.ObjectID, limit int64) ([]model.WorkOrderEvent, error)

	InsertOutboxEvent(ctx context.Context, event model.OutboxEvent) error
}
