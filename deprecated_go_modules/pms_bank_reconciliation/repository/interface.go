package repository

import (
	"context"

	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListFilter struct {
	Page   int64
	Limit  int64
	Status model.TxnStatus
	Q      string
}

type Repository interface {
	InsertTransaction(ctx context.Context, txn model.BankTransaction) (primitive.ObjectID, error)
	FindTransactionByID(ctx context.Context, id primitive.ObjectID) (*model.BankTransaction, error)
	FindTransactionByExternalID(ctx context.Context, externalTxnID string) (*model.BankTransaction, error)
	ListTransactions(ctx context.Context, filter ListFilter) ([]model.BankTransaction, int64, error)
	UpdateMatch(ctx context.Context, id primitive.ObjectID, status model.TxnStatus, candidates []model.MatchCandidate, actor *primitive.ObjectID) error
	ConfirmTransaction(ctx context.Context, id primitive.ObjectID, confirmedRef string, paymentID *primitive.ObjectID, actor *primitive.ObjectID) error
	RejectTransaction(ctx context.Context, id primitive.ObjectID, actor *primitive.ObjectID) error
	InsertOutboxEvent(ctx context.Context, event model.OutboxEvent) error

	ListOpenChargesByAmount(ctx context.Context, amount float64, currency string) ([]ChargeCandidate, error)
}

type ChargeCandidate struct {
	ChargeID        primitive.ObjectID
	UnitID          primitive.ObjectID
	PartyRef        string
	AmountRemaining float64
	Currency        string
	DescriptionHint string
}
