package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	activitylogs "github.com/aystek/aystek-core/internal/modules/activity_logs/service"
	paymentsmodel "github.com/aystek/aystek-core/internal/modules/payments/model"
	paymentsservice "github.com/aystek/aystek-core/internal/modules/payments/service"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/model"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/repository"
	pmsbillingmodel "github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/model"
	pmsbillingservice "github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/service"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrStateTransition = errors.New("state_transition_invalid")

const (
	errInvalidPayload   = "pms.reconciliation.validation.invalid_payload"
	errExternalRequired = "pms.reconciliation.validation.external_required"
	errAmountInvalid    = "pms.reconciliation.validation.amount_invalid"
	errDirectionInvalid = "pms.reconciliation.validation.direction_invalid"
	errDateInvalid      = "pms.reconciliation.validation.date_invalid"
	errTxnNotFound      = "pms.reconciliation.error.not_found"
	errNoCandidates     = "pms.reconciliation.error.no_candidates"
	errCandidateInvalid = "pms.reconciliation.validation.candidate_invalid"
)

type paymentCreator interface {
	Create(ctx context.Context, in paymentsservice.PaymentInput) (*paymentsmodel.Payment, error)
}

type allocationBridge interface {
	AllocatePayment(ctx context.Context, in pmsbillingservice.AllocateInput) (*pmsbillingmodel.AllocationResult, error)
}

type service struct {
	repo       repository.Repository
	logSvc     activitylogs.Service
	payments   paymentCreator
	billingSvc allocationBridge
}

func New(repo repository.Repository, logSvc activitylogs.Service, paySvc paymentCreator, billingSvc allocationBridge) Service {
	return &service{repo: repo, logSvc: logSvc, payments: paySvc, billingSvc: billingSvc}
}

func (s *service) Ingest(ctx context.Context, in IngestInput) (*model.BankTransaction, bool, error) {
	if err := validateContext(ctx); err != nil {
		return nil, false, err
	}
	external := strings.TrimSpace(in.ExternalTxnID)
	if external == "" {
		return nil, false, errors.New(errExternalRequired)
	}
	if in.Amount <= 0 {
		return nil, false, errors.New(errAmountInvalid)
	}
	if !in.Direction.IsValid() {
		return nil, false, errors.New(errDirectionInvalid)
	}
	if strings.TrimSpace(in.TxnDate) == "" {
		return nil, false, errors.New(errDateInvalid)
	}
	txnDate, err := time.Parse(time.RFC3339, strings.TrimSpace(in.TxnDate))
	if err != nil {
		return nil, false, errors.New(errDateInvalid)
	}
	if existing, err := s.repo.FindTransactionByExternalID(ctx, external); err == nil {
		return existing, false, nil
	} else if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, false, err
	}

	txn := model.BankTransaction{
		ExternalTxnID: strings.TrimSpace(in.ExternalTxnID),
		BankAccountID: strings.TrimSpace(in.BankAccountID),
		TxnDate:       txnDate,
		Amount:        in.Amount,
		Currency:      normalizeCurrency(in.Currency),
		Direction:     in.Direction,
		Description:   strings.TrimSpace(in.Description),
		Counterparty:  strings.TrimSpace(in.Counterparty),
		Status:        model.TxnStatusNew,
		CreatedBy:     in.ActorID,
		UpdatedBy:     in.ActorID,
	}
	id, err := s.repo.InsertTransaction(ctx, txn)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			current, findErr := s.repo.FindTransactionByExternalID(ctx, external)
			if findErr == nil {
				return current, false, nil
			}
		}
		return nil, false, err
	}
	created, err := s.repo.FindTransactionByID(ctx, id)
	if err != nil {
		return nil, false, err
	}
	s.emitEvent(ctx, "pms.bank_txn.ingested", id, in.ActorID, map[string]interface{}{"external_txn_id": external, "amount": in.Amount})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, in.ActorID, "admin", "pms", "pms.bank.ingest", "pms_bank_transaction", &id, "Bank transaction ingested", map[string]interface{}{"external_txn_id": external}, "", "")
	}
	return created, true, nil
}

func (s *service) ListQueue(ctx context.Context, filter repository.ListFilter) ([]model.BankTransaction, int64, error) {
	if err := validateContext(ctx); err != nil {
		return nil, 0, err
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	return s.repo.ListTransactions(ctx, filter)
}

func (s *service) ProposeMatch(ctx context.Context, id primitive.ObjectID) (*model.BankTransaction, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	txn, err := s.repo.FindTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errTxnNotFound)
		}
		return nil, err
	}
	if txn.Status != model.TxnStatusNew && txn.Status != model.TxnStatusRejected {
		return nil, fmt.Errorf("%w:%s->%s", ErrStateTransition, txn.Status, model.TxnStatusMatched)
	}
	candidatesRaw, err := s.repo.ListOpenChargesByAmount(ctx, txn.Amount, txn.Currency)
	if err != nil {
		return nil, err
	}
	candidates := make([]model.MatchCandidate, 0, len(candidatesRaw))
	descLower := strings.ToLower(txn.Description + " " + txn.Counterparty)
	for _, c := range candidatesRaw {
		score := 0.5
		diff := absFloat(txn.Amount - c.AmountRemaining)
		if diff < 0.01 {
			score += 0.35
		} else if diff < txn.Amount*0.05 {
			score += 0.2
		}
		if c.PartyRef != "" && strings.Contains(descLower, strings.ToLower(c.PartyRef)) {
			score += 0.15
		}
		if score > 0.99 {
			score = 0.99
		}
		candidates = append(candidates, model.MatchCandidate{
			ChargeID: c.ChargeID,
			UnitID:   c.UnitID,
			PartyRef: c.PartyRef,
			Amount:   c.AmountRemaining,
			Score:    score,
		})
	}
	sort.SliceStable(candidates, func(i, j int) bool { return candidates[i].Score > candidates[j].Score })
	if err := s.repo.UpdateMatch(ctx, txn.ID, model.TxnStatusMatched, candidates, nil); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindTransactionByID(ctx, txn.ID)
	if err != nil {
		return nil, err
	}
	s.emitEvent(ctx, "pms.bank_match.proposed", txn.ID, nil, map[string]interface{}{"candidates": len(candidates)})
	return updated, nil
}

func (s *service) Confirm(ctx context.Context, id primitive.ObjectID, in ConfirmInput) (*model.BankTransaction, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	txn, err := s.repo.FindTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errTxnNotFound)
		}
		return nil, err
	}
	if txn.Status != model.TxnStatusMatched {
		return nil, fmt.Errorf("%w:%s->%s", ErrStateTransition, txn.Status, model.TxnStatusConfirmed)
	}
	if len(txn.MatchCandidates) == 0 {
		return nil, errors.New(errNoCandidates)
	}
	chosen := txn.MatchCandidates[0]
	if in.CandidateChargeID != nil && *in.CandidateChargeID != primitive.NilObjectID {
		found := false
		for _, c := range txn.MatchCandidates {
			if c.ChargeID == *in.CandidateChargeID {
				chosen = c
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New(errCandidateInvalid)
		}
	}
	var paymentID *primitive.ObjectID
	if in.CreatePayment && s.payments != nil {
		receivedAt := txn.TxnDate.UTC().Format(time.RFC3339)
		created, err := s.payments.Create(ctx, paymentsservice.PaymentInput{
			Amount:          txn.Amount,
			Currency:        txn.Currency,
			PaymentMethod:   paymentsmodel.MethodBankTransfer,
			PaymentReason:   "pms_bank_reconciliation",
			Status:          paymentsmodel.PaymentReceived,
			Description:     strings.TrimSpace(txn.Description),
			EntityType:      "pms_charge",
			EntityID:        chosen.ChargeID.Hex(),
			ReceivedAt:      &receivedAt,
			ReferenceNumber: txn.ExternalTxnID,
			CreatedBy:       actorHex(in.ActorID),
			UpdatedBy:       actorHex(in.ActorID),
		})
		if err != nil {
			return nil, err
		}
		paymentID = &created.ID
		if s.billingSvc != nil {
			_, err = s.billingSvc.AllocatePayment(ctx, pmsbillingservice.AllocateInput{
				PaymentID: paymentID,
				PayerRef:  chosen.PartyRef,
				Currency:  txn.Currency,
				ActorID:   in.ActorID,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	if err := s.repo.ConfirmTransaction(ctx, id, chosen.ChargeID.Hex(), paymentID, in.ActorID); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.emitEvent(ctx, "pms.bank_match.confirmed", id, in.ActorID, map[string]interface{}{"charge_id": chosen.ChargeID.Hex(), "payment_id": oidHex(paymentID)})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, in.ActorID, "admin", "pms", "pms.bank.confirm", "pms_bank_transaction", &id, "Bank transaction confirmed", map[string]interface{}{"charge_id": chosen.ChargeID.Hex()}, "", "")
	}
	return updated, nil
}

func (s *service) Reject(ctx context.Context, id primitive.ObjectID, reason string, actorID *primitive.ObjectID) (*model.BankTransaction, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	txn, err := s.repo.FindTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errTxnNotFound)
		}
		return nil, err
	}
	if txn.Status != model.TxnStatusMatched {
		return nil, fmt.Errorf("%w:%s->%s", ErrStateTransition, txn.Status, model.TxnStatusRejected)
	}
	if err := s.repo.RejectTransaction(ctx, id, actorID); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.emitEvent(ctx, "pms.bank_match.rejected", id, actorID, map[string]interface{}{"reason": strings.TrimSpace(reason)})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, actorID, "admin", "pms", "pms.bank.reject", "pms_bank_transaction", &id, "Bank transaction rejected", map[string]interface{}{"reason": strings.TrimSpace(reason)}, "", "")
	}
	return updated, nil
}

func validateContext(ctx context.Context) error {
	if _, err := tenantfilter.RequireTenantID(ctx); err != nil {
		return err
	}
	if _, err := tenantfilter.RequireSiteID(ctx); err != nil {
		return err
	}
	return nil
}

func (s *service) emitEvent(ctx context.Context, name string, id primitive.ObjectID, actor *primitive.ObjectID, diff map[string]interface{}) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return
	}
	now := time.Now()
	payload := map[string]interface{}{
		"event_id":    primitive.NewObjectID().Hex(),
		"tenant_id":   tenantID,
		"site_id":     siteID,
		"entity_id":   id.Hex(),
		"actor_id":    actorHex(actor),
		"occurred_at": now.UTC().Format(time.RFC3339),
		"version":     1,
	}
	if len(diff) > 0 {
		payload["diff"] = diff
	}
	_ = s.repo.InsertOutboxEvent(ctx, model.OutboxEvent{TenantID: tenantID, SiteID: siteID, EventName: name, EntityID: id, ActorID: actor, OccurredAt: now, Version: 1, Diff: diff, Payload: payload, Processed: false})
}

func normalizeCurrency(v string) string {
	c := strings.ToUpper(strings.TrimSpace(v))
	if c == "" {
		return "TRY"
	}
	return c
}

func tenantSiteIDs(ctx context.Context) (string, string, error) {
	tenantID, err := tenantfilter.RequireTenantID(ctx)
	if err != nil {
		return "", "", err
	}
	siteID, err := tenantfilter.RequireSiteID(ctx)
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(tenantID), strings.TrimSpace(siteID), nil
}

func actorHex(actor *primitive.ObjectID) string {
	if actor == nil || *actor == primitive.NilObjectID {
		return ""
	}
	return actor.Hex()
}

func oidHex(id *primitive.ObjectID) string {
	if id == nil || *id == primitive.NilObjectID {
		return ""
	}
	return id.Hex()
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

var _ Service = (*service)(nil)
