package service

import (
	"context"

	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	PublishDues(ctx context.Context, in PublishInput) (*model.PublishResult, error)
	AllocatePayment(ctx context.Context, in AllocateInput) (*model.AllocationResult, error)
	Overview(ctx context.Context, propertyID *primitive.ObjectID, period string) (*model.BillingOverview, error)
}

type PublishInput struct {
	PropertyID     primitive.ObjectID
	Period         string
	PlanID         *primitive.ObjectID
	DefaultAmount  *float64
	Currency       string
	IdempotencyKey string
	ActorID        *primitive.ObjectID
}

type AllocateInput struct {
	PaymentID      *primitive.ObjectID
	Amount         *float64
	PayerRef       string
	Date           *string
	Currency       string
	PropertyID     *primitive.ObjectID
	IdempotencyKey string
	ActorID        *primitive.ObjectID
}
