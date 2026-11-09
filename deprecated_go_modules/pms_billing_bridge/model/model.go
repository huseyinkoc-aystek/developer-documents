package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChargeStatus string

const (
	ChargeStatusOpen    ChargeStatus = "open"
	ChargeStatusPartial ChargeStatus = "partial"
	ChargeStatusPaid    ChargeStatus = "paid"
)

type PMSCharge struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	TenantID        string              `bson:"tenant_id" json:"tenant_id"`
	SiteID          string              `bson:"site_id" json:"site_id"`
	PropertyID      primitive.ObjectID  `bson:"property_id" json:"property_id"`
	UnitID          primitive.ObjectID  `bson:"unit_id" json:"unit_id"`
	PartyRef        string              `bson:"party_ref" json:"party_ref"`
	Period          string              `bson:"period" json:"period"`
	PlanID          *primitive.ObjectID `bson:"plan_id,omitempty" json:"plan_id,omitempty"`
	AmountDue       float64             `bson:"amount_due" json:"amount_due"`
	AmountPaid      float64             `bson:"amount_paid" json:"amount_paid"`
	AmountRemaining float64             `bson:"amount_remaining" json:"amount_remaining"`
	Currency        string              `bson:"currency" json:"currency"`
	DueDate         time.Time           `bson:"due_date" json:"due_date"`
	Status          ChargeStatus        `bson:"status" json:"status"`
	Source          string              `bson:"source" json:"source"`
	CreatedAt       time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time           `bson:"updated_at" json:"updated_at"`
	CreatedBy       *primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy       *primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	DeletedAt       *time.Time          `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	DeletedBy       *primitive.ObjectID `bson:"deleted_by,omitempty" json:"deleted_by,omitempty"`
}

type PaymentAllocation struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	TenantID        string              `bson:"tenant_id" json:"tenant_id"`
	SiteID          string              `bson:"site_id" json:"site_id"`
	ChargeID        primitive.ObjectID  `bson:"charge_id" json:"charge_id"`
	PropertyID      primitive.ObjectID  `bson:"property_id" json:"property_id"`
	UnitID          primitive.ObjectID  `bson:"unit_id" json:"unit_id"`
	PartyRef        string              `bson:"party_ref" json:"party_ref"`
	PaymentRef      string              `bson:"payment_ref" json:"payment_ref"`
	AllocatedAmount float64             `bson:"allocated_amount" json:"allocated_amount"`
	Currency        string              `bson:"currency" json:"currency"`
	AllocatedAt     time.Time           `bson:"allocated_at" json:"allocated_at"`
	CreatedAt       time.Time           `bson:"created_at" json:"created_at"`
	CreatedBy       *primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
}

type BillingIdempotency struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	TenantID       string             `bson:"tenant_id"`
	SiteID         string             `bson:"site_id"`
	Operation      string             `bson:"operation"`
	IdempotencyKey string             `bson:"idempotency_key"`
	RequestHash    string             `bson:"request_hash"`
	Response       interface{}        `bson:"response"`
	CreatedAt      time.Time          `bson:"created_at"`
}

type OutboxEvent struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty"`
	TenantID   string                 `bson:"tenant_id"`
	SiteID     string                 `bson:"site_id"`
	EventName  string                 `bson:"event_name"`
	EntityID   primitive.ObjectID     `bson:"entity_id"`
	ActorID    *primitive.ObjectID    `bson:"actor_id,omitempty"`
	OccurredAt time.Time              `bson:"occurred_at"`
	Version    int                    `bson:"version"`
	Diff       map[string]interface{} `bson:"diff,omitempty"`
	Payload    map[string]interface{} `bson:"payload,omitempty"`
	Processed  bool                   `bson:"processed"`
	CreatedAt  time.Time              `bson:"created_at"`
}

type BillingOverview struct {
	TotalDue       float64             `json:"total_due"`
	TotalPaid      float64             `json:"total_paid"`
	TotalRemaining float64             `json:"total_remaining"`
	Currency       string              `json:"currency"`
	PropertyID     *primitive.ObjectID `json:"property_id,omitempty"`
	Period         string              `json:"period,omitempty"`
	Units          []UnitOverview      `json:"units"`
}

type UnitOverview struct {
	UnitID           primitive.ObjectID `json:"unit_id"`
	PartyRef         string             `json:"party_ref"`
	TotalDue         float64            `json:"total_due"`
	TotalPaid        float64            `json:"total_paid"`
	TotalRemaining   float64            `json:"total_remaining"`
	LastChargeAt     *time.Time         `json:"last_charge_at,omitempty"`
	LastAllocationAt *time.Time         `json:"last_allocation_at,omitempty"`
}

type AllocationResult struct {
	PaymentRef      string  `json:"payment_ref"`
	RequestedAmount float64 `json:"requested_amount"`
	AllocatedAmount float64 `json:"allocated_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Currency        string  `json:"currency"`
	Allocations     int     `json:"allocations"`
}

type PublishResult struct {
	PropertyID    primitive.ObjectID `json:"property_id"`
	Period        string             `json:"period"`
	PlanID        *primitive.ObjectID `json:"plan_id,omitempty"`
	CreatedCount  int                `json:"created_count"`
	SkippedCount  int                `json:"skipped_count"`
	Currency      string             `json:"currency"`
	AmountPerUnit float64            `json:"amount_per_unit"`
}

type DuesPlanMethod string

const (
	DuesPlanMethodFixed     DuesPlanMethod = "fixed"
	DuesPlanMethodAreaBased DuesPlanMethod = "area_based"
)

type DuesPlanReference struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PropertyID    primitive.ObjectID `bson:"property_id" json:"property_id"`
	Method        DuesPlanMethod     `bson:"method" json:"method"`
	Amount        float64            `bson:"amount" json:"amount"`
	RatePerSQM    float64            `bson:"rate_per_sqm" json:"rate_per_sqm"`
	Currency      string             `bson:"currency" json:"currency"`
	EffectiveFrom time.Time          `bson:"effective_from" json:"effective_from"`
}
