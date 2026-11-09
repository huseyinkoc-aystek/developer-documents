package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitListItem struct {
	UnitID              primitive.ObjectID `json:"unit_id"`
	PropertyID          primitive.ObjectID `json:"property_id"`
	PropertyName        string             `json:"property_name"`
	Block               string             `json:"block,omitempty"`
	Number              string             `json:"number,omitempty"`
	Floor               string             `json:"floor,omitempty"`
	Status              string             `json:"status,omitempty"`
	OccupancySnapshot   OccupancySnapshot  `json:"occupancy_snapshot"`
	DuesSummary         DuesSummary        `json:"dues_summary"`
	LastPayment         PaymentSummary     `json:"last_payment"`
	OpenWorkOrdersCount int64              `json:"open_work_orders_count"`
}

type UnitDetail struct {
	UnitListItem
	LastWorkOrders []WorkOrderMini `json:"last_work_orders"`
}

type OccupancySnapshot struct {
	OccupancyID           primitive.ObjectID `json:"occupancy_id,omitempty"`
	PartyRef              string             `json:"party_ref,omitempty"`
	OccupantType          string             `json:"occupant_type,omitempty"`
	Status                string             `json:"status,omitempty"`
	StartAt               *time.Time         `json:"start_at,omitempty"`
	EndAt                 *time.Time         `json:"end_at,omitempty"`
	FinancialResponsible  bool               `json:"financial_responsible"`
	FinancialSharePercent *float64           `json:"financial_share_percent,omitempty"`
}

type DuesSummary struct {
	TotalDue       float64 `json:"total_due"`
	TotalPaid      float64 `json:"total_paid"`
	TotalRemaining float64 `json:"total_remaining"`
	ChargeCount    int64   `json:"charge_count"`
	Currency       string  `json:"currency,omitempty"`
}

type PaymentSummary struct {
	PaymentRef      string     `json:"payment_ref,omitempty"`
	AllocatedAmount float64    `json:"allocated_amount"`
	Currency        string     `json:"currency,omitempty"`
	AllocatedAt     *time.Time `json:"allocated_at,omitempty"`
}

type WorkOrderMini struct {
	ID        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	Status    string             `json:"status"`
	Priority  string             `json:"priority"`
	CreatedAt time.Time          `json:"created_at"`
}
