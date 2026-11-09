package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TxnDirection string

type TxnStatus string

const (
	TxnDirectionCredit TxnDirection = "credit"
	TxnDirectionDebit  TxnDirection = "debit"
)

const (
	TxnStatusNew       TxnStatus = "new"
	TxnStatusMatched   TxnStatus = "matched"
	TxnStatusConfirmed TxnStatus = "confirmed"
	TxnStatusRejected  TxnStatus = "rejected"
)

func (v TxnDirection) IsValid() bool {
	switch v {
	case TxnDirectionCredit, TxnDirectionDebit:
		return true
	default:
		return false
	}
}

func (v TxnStatus) IsValid() bool {
	switch v {
	case TxnStatusNew, TxnStatusMatched, TxnStatusConfirmed, TxnStatusRejected:
		return true
	default:
		return false
	}
}

type MatchCandidate struct {
	ChargeID primitive.ObjectID `bson:"charge_id" json:"charge_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	PartyRef string             `bson:"party_ref" json:"party_ref"`
	Amount   float64            `bson:"amount" json:"amount"`
	Score    float64            `bson:"score" json:"score"`
}

type BankTransaction struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	TenantID          string              `bson:"tenant_id" json:"tenant_id"`
	SiteID            string              `bson:"site_id" json:"site_id"`
	ExternalTxnID     string              `bson:"external_txn_id" json:"external_txn_id"`
	BankAccountID     string              `bson:"bank_account_id" json:"bank_account_id"`
	TxnDate           time.Time           `bson:"txn_date" json:"txn_date"`
	Amount            float64             `bson:"amount" json:"amount"`
	Currency          string              `bson:"currency" json:"currency"`
	Direction         TxnDirection        `bson:"direction" json:"direction"`
	Description       string              `bson:"description" json:"description"`
	Counterparty      string              `bson:"counterparty" json:"counterparty"`
	Status            TxnStatus           `bson:"status" json:"status"`
	MatchCandidates   []MatchCandidate    `bson:"match_candidates,omitempty" json:"match_candidates,omitempty"`
	ConfirmedMatchRef string              `bson:"confirmed_match_ref,omitempty" json:"confirmed_match_ref,omitempty"`
	PaymentID         *primitive.ObjectID `bson:"payment_id,omitempty" json:"payment_id,omitempty"`
	CreatedAt         time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time           `bson:"updated_at" json:"updated_at"`
	CreatedBy         *primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy         *primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	DeletedAt         *time.Time          `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	DeletedBy         *primitive.ObjectID `bson:"deleted_by,omitempty" json:"deleted_by,omitempty"`
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
