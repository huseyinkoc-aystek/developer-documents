package service

import (
	"context"

	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/model"
	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(ctx context.Context, in CreateInput) (*model.WorkOrderDetail, error)
	List(ctx context.Context, filter repository.ListFilter) ([]model.WorkOrder, int64, error)
	Get(ctx context.Context, id primitive.ObjectID) (*model.WorkOrderDetail, error)
	Update(ctx context.Context, id primitive.ObjectID, in UpdateInput) (*model.WorkOrderDetail, error)
	Assign(ctx context.Context, id primitive.ObjectID, assignee *primitive.ObjectID, actor *primitive.ObjectID) (*model.WorkOrderDetail, error)
	ChangeStatus(ctx context.Context, id primitive.ObjectID, status model.WorkOrderStatus, actor *primitive.ObjectID) (*model.WorkOrderDetail, error)
	Comment(ctx context.Context, id primitive.ObjectID, comment string, actor *primitive.ObjectID) (*model.WorkOrderDetail, error)
	Close(ctx context.Context, id primitive.ObjectID, actor *primitive.ObjectID) (*model.WorkOrderDetail, error)
}

type CreateInput struct {
	PropertyID  primitive.ObjectID
	UnitID      *primitive.ObjectID
	Title       string
	Description string
	Category    model.WorkOrderCategory
	Priority    model.WorkOrderPriority
	SLADueAt    *string
	ReporterRef string
	Tags        []string
	Attachments []model.Attachment
	CostItems   []model.CostItem
	ActorID     *primitive.ObjectID
}

type UpdateInput struct {
	Title       string
	Description string
	Category    model.WorkOrderCategory
	Priority    model.WorkOrderPriority
	SLADueAt    *string
	Tags        []string
	Attachments []model.Attachment
	CostItems   []model.CostItem
	ActorID     *primitive.ObjectID
}
