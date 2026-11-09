package service

import (
	"context"

	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/model"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Ingest(ctx context.Context, in IngestInput) (*model.BankTransaction, bool, error)
	ListQueue(ctx context.Context, filter repository.ListFilter) ([]model.BankTransaction, int64, error)
	ProposeMatch(ctx context.Context, id primitive.ObjectID) (*model.BankTransaction, error)
	Confirm(ctx context.Context, id primitive.ObjectID, in ConfirmInput) (*model.BankTransaction, error)
	Reject(ctx context.Context, id primitive.ObjectID, reason string, actorID *primitive.ObjectID) (*model.BankTransaction, error)
}

type IngestInput struct {
	ExternalTxnID string
	BankAccountID string
	TxnDate       string
	Amount        float64
	Currency      string
	Direction     model.TxnDirection
	Description   string
	Counterparty  string
	ActorID       *primitive.ObjectID
}

type ConfirmInput struct {
	CandidateChargeID *primitive.ObjectID
	CreatePayment     bool
	ActorID           *primitive.ObjectID
}
