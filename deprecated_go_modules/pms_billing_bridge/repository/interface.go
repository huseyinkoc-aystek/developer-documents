package repository

import (
	"context"
	"time"

	pmsmodel "github.com/aystek/aystek-core/internal/modules/pms/model"
	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	FindProperty(ctx context.Context, id primitive.ObjectID) (*pmsmodel.Property, error)
	FindActiveOccupanciesByPropertyPeriod(ctx context.Context, propertyID primitive.ObjectID, from, to time.Time) ([]pmsmodel.Occupancy, error)
	FindDuesPlanByID(ctx context.Context, id primitive.ObjectID) (*model.DuesPlanReference, error)
	FindActiveDuesPlan(ctx context.Context, propertyID primitive.ObjectID, at time.Time) (*model.DuesPlanReference, error)
	FindUnitAreaSQM(ctx context.Context, unitID primitive.ObjectID) (float64, error)

	FindChargeByUnitPeriod(ctx context.Context, propertyID, unitID primitive.ObjectID, period string) (*model.PMSCharge, error)
	InsertCharge(ctx context.Context, charge model.PMSCharge) (primitive.ObjectID, error)
	UpdateChargeAmounts(ctx context.Context, id primitive.ObjectID, paid, remaining float64, status model.ChargeStatus, actor *primitive.ObjectID) error
	ListOpenChargesByParty(ctx context.Context, partyRef string, propertyID *primitive.ObjectID) ([]model.PMSCharge, error)
	ListCharges(ctx context.Context, propertyID *primitive.ObjectID, period string) ([]model.PMSCharge, error)

	InsertAllocation(ctx context.Context, allocation model.PaymentAllocation) (primitive.ObjectID, error)
	FindPaymentByID(ctx context.Context, id primitive.ObjectID) (*PaymentRef, error)

	FindIdempotency(ctx context.Context, operation, key string) (*model.BillingIdempotency, error)
	InsertIdempotency(ctx context.Context, rec model.BillingIdempotency) error

	InsertOutboxEvent(ctx context.Context, event model.OutboxEvent) error
}

type PaymentRef struct {
	ID         primitive.ObjectID
	Amount     float64
	Currency   string
	PayerRef   string
	ReceivedAt time.Time
}
