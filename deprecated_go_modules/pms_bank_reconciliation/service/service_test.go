package service

import (
	"context"
	"errors"
	"testing"
	"time"

	paymentsmodel "github.com/aystek/aystek-core/internal/modules/payments/model"
	paymentsservice "github.com/aystek/aystek-core/internal/modules/payments/service"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/model"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/repository"
	pmsbillingmodel "github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/model"
	pmsbillingservice "github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/service"
	"github.com/aystek/aystek-core/internal/platform/contextkeys"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockRepo struct {
	txns   map[primitive.ObjectID]model.BankTransaction
	extMap map[string]primitive.ObjectID
	open   []repository.ChargeCandidate
	outbox []model.OutboxEvent
}

func newMockRepo() *mockRepo {
	return &mockRepo{txns: map[primitive.ObjectID]model.BankTransaction{}, extMap: map[string]primitive.ObjectID{}}
}

func (m *mockRepo) InsertTransaction(_ context.Context, txn model.BankTransaction) (primitive.ObjectID, error) {
	if _, ok := m.extMap[txn.ExternalTxnID]; ok {
		return primitive.NilObjectID, mongo.WriteException{}
	}
	if txn.ID == primitive.NilObjectID {
		txn.ID = primitive.NewObjectID()
	}
	m.txns[txn.ID] = txn
	m.extMap[txn.ExternalTxnID] = txn.ID
	return txn.ID, nil
}

func (m *mockRepo) FindTransactionByID(_ context.Context, id primitive.ObjectID) (*model.BankTransaction, error) {
	txn, ok := m.txns[id]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	copy := txn
	return &copy, nil
}

func (m *mockRepo) FindTransactionByExternalID(_ context.Context, externalTxnID string) (*model.BankTransaction, error) {
	id, ok := m.extMap[externalTxnID]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	txn := m.txns[id]
	copy := txn
	return &copy, nil
}

func (m *mockRepo) ListTransactions(_ context.Context, filter repository.ListFilter) ([]model.BankTransaction, int64, error) {
	items := make([]model.BankTransaction, 0)
	for _, v := range m.txns {
		if filter.Status != "" && v.Status != filter.Status {
			continue
		}
		items = append(items, v)
	}
	return items, int64(len(items)), nil
}

func (m *mockRepo) UpdateMatch(_ context.Context, id primitive.ObjectID, status model.TxnStatus, candidates []model.MatchCandidate, actor *primitive.ObjectID) error {
	txn, ok := m.txns[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	txn.Status = status
	txn.MatchCandidates = candidates
	txn.UpdatedBy = actor
	m.txns[id] = txn
	return nil
}

func (m *mockRepo) ConfirmTransaction(_ context.Context, id primitive.ObjectID, confirmedRef string, paymentID *primitive.ObjectID, actor *primitive.ObjectID) error {
	txn, ok := m.txns[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	txn.Status = model.TxnStatusConfirmed
	txn.ConfirmedMatchRef = confirmedRef
	txn.PaymentID = paymentID
	txn.UpdatedBy = actor
	m.txns[id] = txn
	return nil
}

func (m *mockRepo) RejectTransaction(_ context.Context, id primitive.ObjectID, actor *primitive.ObjectID) error {
	txn, ok := m.txns[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	txn.Status = model.TxnStatusRejected
	txn.UpdatedBy = actor
	m.txns[id] = txn
	return nil
}

func (m *mockRepo) InsertOutboxEvent(_ context.Context, event model.OutboxEvent) error {
	m.outbox = append(m.outbox, event)
	return nil
}

func (m *mockRepo) ListOpenChargesByAmount(_ context.Context, amount float64, currency string) ([]repository.ChargeCandidate, error) {
	return append([]repository.ChargeCandidate(nil), m.open...), nil
}

type mockPay struct {
	created int
}

func (m *mockPay) Create(ctx context.Context, in paymentsservice.PaymentInput) (*paymentsmodel.Payment, error) {
	m.created++
	id := primitive.NewObjectID()
	return &paymentsmodel.Payment{ID: id, Amount: in.Amount, Currency: in.Currency, EntityID: in.EntityID}, nil
}

type mockAlloc struct {
	called int
}

func (m *mockAlloc) AllocatePayment(ctx context.Context, in pmsbillingservice.AllocateInput) (*pmsbillingmodel.AllocationResult, error) {
	if in.PaymentID == nil {
		return nil, errors.New("payment required")
	}
	m.called++
	return &pmsbillingmodel.AllocationResult{AllocatedAmount: 1}, nil
}

func testCtx() context.Context {
	ctx := context.Background()
	ctx = contextkeys.WithTenantID(ctx, "tenant-1")
	ctx = contextkeys.WithSiteID(ctx, "site-1")
	return ctx
}

func TestIngestDuplicateExternalID(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo, nil, &mockPay{}, &mockAlloc{})
	in := IngestInput{ExternalTxnID: "EXT-1", TxnDate: time.Now().UTC().Format(time.RFC3339), Amount: 100, Currency: "TRY", Direction: model.TxnDirectionCredit}
	first, created, err := svc.Ingest(testCtx(), in)
	if err != nil || !created || first == nil {
		t.Fatalf("first ingest expected created, err=%v", err)
	}
	second, created2, err := svc.Ingest(testCtx(), in)
	if err != nil || created2 {
		t.Fatalf("second ingest expected idempotent replay, created=%v err=%v", created2, err)
	}
	if second.ID != first.ID {
		t.Fatalf("expected same txn id for duplicate ingest")
	}
}

func TestInvalidStateTransition(t *testing.T) {
	repo := newMockRepo()
	txnID := primitive.NewObjectID()
	repo.txns[txnID] = model.BankTransaction{ID: txnID, ExternalTxnID: "E2", Status: model.TxnStatusConfirmed, Amount: 10, Currency: "TRY"}
	repo.extMap["E2"] = txnID
	svc := New(repo, nil, &mockPay{}, &mockAlloc{})
	_, err := svc.ProposeMatch(testCtx(), txnID)
	if err == nil || !errors.Is(err, ErrStateTransition) {
		t.Fatalf("expected state transition error, got %v", err)
	}
}

func TestMatchProposalGenerated(t *testing.T) {
	repo := newMockRepo()
	txnID := primitive.NewObjectID()
	repo.txns[txnID] = model.BankTransaction{ID: txnID, ExternalTxnID: "E3", Status: model.TxnStatusNew, Amount: 100, Currency: "TRY", Description: "payer-a"}
	repo.extMap["E3"] = txnID
	repo.open = []repository.ChargeCandidate{{ChargeID: primitive.NewObjectID(), UnitID: primitive.NewObjectID(), PartyRef: "payer-a", AmountRemaining: 100, Currency: "TRY"}}
	svc := New(repo, nil, &mockPay{}, &mockAlloc{})
	out, err := svc.ProposeMatch(testCtx(), txnID)
	if err != nil {
		t.Fatalf("propose error: %v", err)
	}
	if out.Status != model.TxnStatusMatched || len(out.MatchCandidates) == 0 {
		t.Fatalf("expected matched with candidates, got %+v", out)
	}
}

func TestConfirmCreatesPaymentAndEvents(t *testing.T) {
	repo := newMockRepo()
	chargeID := primitive.NewObjectID()
	txnID := primitive.NewObjectID()
	repo.txns[txnID] = model.BankTransaction{
		ID:              txnID,
		ExternalTxnID:   "E4",
		Status:          model.TxnStatusMatched,
		Amount:          150,
		Currency:        "TRY",
		TxnDate:         time.Now().UTC(),
		MatchCandidates: []model.MatchCandidate{{ChargeID: chargeID, UnitID: primitive.NewObjectID(), PartyRef: "payer-z", Amount: 150, Score: 0.9}},
	}
	repo.extMap["E4"] = txnID
	pay := &mockPay{}
	alloc := &mockAlloc{}
	svc := New(repo, nil, pay, alloc)
	out, err := svc.Confirm(testCtx(), txnID, ConfirmInput{CreatePayment: true})
	if err != nil {
		t.Fatalf("confirm error: %v", err)
	}
	if out.Status != model.TxnStatusConfirmed {
		t.Fatalf("expected confirmed status")
	}
	if pay.created != 1 || alloc.called != 1 {
		t.Fatalf("expected payment+allocation called once, got pay=%d alloc=%d", pay.created, alloc.called)
	}
	if len(repo.outbox) == 0 {
		t.Fatalf("expected outbox events")
	}
}

func TestRejectTransition(t *testing.T) {
	repo := newMockRepo()
	txnID := primitive.NewObjectID()
	repo.txns[txnID] = model.BankTransaction{ID: txnID, ExternalTxnID: "E5", Status: model.TxnStatusMatched, Amount: 100, Currency: "TRY"}
	repo.extMap["E5"] = txnID
	svc := New(repo, nil, &mockPay{}, &mockAlloc{})
	out, err := svc.Reject(testCtx(), txnID, "manual", nil)
	if err != nil {
		t.Fatalf("reject error: %v", err)
	}
	if out.Status != model.TxnStatusRejected {
		t.Fatalf("expected rejected")
	}
}
