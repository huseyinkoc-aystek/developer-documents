package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	activitylogs "github.com/aystek/aystek-core/internal/modules/activity_logs/service"
	pmsmodel "github.com/aystek/aystek-core/internal/modules/pms/model"
	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/model"
	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/repository"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	errInvalidPayload      = "pms.billing.validation.invalid_payload"
	errPropertyRequired    = "pms.billing.validation.property_required"
	errPeriodRequired      = "pms.billing.validation.period_required"
	errPeriodInvalid       = "pms.billing.validation.period_invalid"
	errAmountInvalid       = "pms.billing.validation.amount_invalid"
	errPayerRefRequired    = "pms.billing.validation.payer_ref_required"
	errPaymentOrAmount     = "pms.billing.validation.payment_or_amount_required"
	errNoOpenCharges       = "pms.billing.error.no_open_charges"
	errPropertyNotFound    = "pms.properties.error.not_found"
	errIdempotencyConflict = "pms.billing.validation.idempotency_conflict"
)

type service struct {
	repo   repository.Repository
	logSvc activitylogs.Service
}

func New(repo repository.Repository, logSvc activitylogs.Service) Service {
	return &service{repo: repo, logSvc: logSvc}
}

func (s *service) PublishDues(ctx context.Context, in PublishInput) (*model.PublishResult, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	if in.PropertyID == primitive.NilObjectID {
		return nil, errors.New(errPropertyRequired)
	}
	period := strings.TrimSpace(in.Period)
	if period == "" {
		return nil, errors.New(errPeriodRequired)
	}
	periodStart, periodEnd, err := parsePeriod(period)
	if err != nil {
		return nil, errors.New(errPeriodInvalid)
	}
	amountPerUnit := 0.0
	if in.DefaultAmount != nil {
		if *in.DefaultAmount < 0 {
			return nil, errors.New(errAmountInvalid)
		}
		amountPerUnit = *in.DefaultAmount
	}
	currency := strings.ToUpper(strings.TrimSpace(in.Currency))
	if currency == "" {
		currency = "TRY"
	}
	var selectedPlan *model.DuesPlanReference
	if in.PlanID != nil && *in.PlanID != primitive.NilObjectID {
		selectedPlan, err = s.repo.FindDuesPlanByID(ctx, *in.PlanID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, errors.New("pms.billing.validation.plan_not_found")
			}
			return nil, err
		}
	} else {
		if plan, planErr := s.repo.FindActiveDuesPlan(ctx, in.PropertyID, periodStart); planErr == nil {
			selectedPlan = plan
		} else if planErr != nil && !errors.Is(planErr, mongo.ErrNoDocuments) {
			return nil, planErr
		}
	}
	if selectedPlan != nil {
		in.PlanID = &selectedPlan.ID
		if currency == "" || strings.EqualFold(currency, "TRY") {
			currency = strings.ToUpper(strings.TrimSpace(selectedPlan.Currency))
			if currency == "" {
				currency = "TRY"
			}
		}
		if selectedPlan.Method == model.DuesPlanMethodFixed && selectedPlan.Amount >= 0 {
			amountPerUnit = selectedPlan.Amount
		}
	}

	property, err := s.repo.FindProperty(ctx, in.PropertyID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errPropertyNotFound)
		}
		return nil, err
	}

	requestHash := hashPayload(map[string]interface{}{
		"property_id": in.PropertyID.Hex(),
		"period":      period,
		"plan_id":     optionalOID(in.PlanID),
		"amount":      amountPerUnit,
		"currency":    currency,
	})
	if in.IdempotencyKey != "" {
		if replay, err := s.idempotencyReplayPublish(ctx, "publish_dues", in.IdempotencyKey, requestHash); err != nil || replay != nil {
			return replay, err
		}
	}

	occupancies, err := s.repo.FindActiveOccupanciesByPropertyPeriod(ctx, in.PropertyID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}
	byUnit := make(map[primitive.ObjectID][]pmsmodel.Occupancy)
	for _, occ := range occupancies {
		byUnit[occ.UnitID] = append(byUnit[occ.UnitID], occ)
	}

	res := &model.PublishResult{
		PropertyID:    in.PropertyID,
		Period:        period,
		PlanID:        in.PlanID,
		Currency:      currency,
		AmountPerUnit: amountPerUnit,
	}
	for unitID, unitOccs := range byUnit {
		if _, err := s.repo.FindChargeByUnitPeriod(ctx, in.PropertyID, unitID, period); err == nil {
			res.SkippedCount++
			continue
		} else if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}

		responsible := pickResponsible(unitOccs, periodEnd)
		partyRef := responsible.PartyRef
		if strings.TrimSpace(partyRef) == "" {
			partyRef = "unassigned"
		}
		unitAmount := amountPerUnit
		if selectedPlan != nil && selectedPlan.Method == model.DuesPlanMethodAreaBased {
			unitArea, areaErr := s.repo.FindUnitAreaSQM(ctx, unitID)
			if areaErr != nil && !errors.Is(areaErr, mongo.ErrNoDocuments) {
				return nil, areaErr
			}
			if unitArea > 0 {
				unitAmount = selectedPlan.RatePerSQM * unitArea
			} else if in.DefaultAmount == nil {
				res.SkippedCount++
				continue
			}
		}

		due := periodEnd
		status := model.ChargeStatusOpen
		if unitAmount <= 0 {
			status = model.ChargeStatusPaid
		}
		charge := model.PMSCharge{
			PropertyID:      in.PropertyID,
			UnitID:          unitID,
			PartyRef:        partyRef,
			Period:          period,
			PlanID:          in.PlanID,
			AmountDue:       unitAmount,
			AmountPaid:      0,
			AmountRemaining: unitAmount,
			Currency:        currency,
			DueDate:         due,
			Status:          status,
			Source:          "pms_billing_bridge",
			CreatedBy:       in.ActorID,
			UpdatedBy:       in.ActorID,
		}
		id, err := s.repo.InsertCharge(ctx, charge)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				res.SkippedCount++
				continue
			}
			return nil, err
		}
		res.CreatedCount++
		s.emitEvent(ctx, "pms.charge.created", id, in.ActorID, map[string]interface{}{"property_id": in.PropertyID.Hex(), "unit_id": unitID.Hex(), "period": period})
	}

	s.emitEvent(ctx, "pms.dues.published", in.PropertyID, in.ActorID, map[string]interface{}{"period": period, "created": res.CreatedCount, "skipped": res.SkippedCount})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, in.ActorID, "admin", "pms", "pms.billing.publish", "pms_property", &property.ID,
			"PMS dues published", map[string]interface{}{"period": period, "created_count": res.CreatedCount, "skipped_count": res.SkippedCount}, "", "")
	}
	if in.IdempotencyKey != "" {
		_ = s.repo.InsertIdempotency(ctx, model.BillingIdempotency{
			Operation:      "publish_dues",
			IdempotencyKey: in.IdempotencyKey,
			RequestHash:    requestHash,
			Response:       res,
		})
	}
	return res, nil
}

func (s *service) AllocatePayment(ctx context.Context, in AllocateInput) (*model.AllocationResult, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	amount := 0.0
	currency := strings.ToUpper(strings.TrimSpace(in.Currency))
	payerRef := strings.TrimSpace(in.PayerRef)
	paymentRef := ""

	if in.PaymentID != nil && *in.PaymentID != primitive.NilObjectID {
		pay, err := s.repo.FindPaymentByID(ctx, *in.PaymentID)
		if err != nil {
			return nil, err
		}
		amount = pay.Amount
		if currency == "" {
			currency = strings.ToUpper(strings.TrimSpace(pay.Currency))
		}
		if payerRef == "" {
			payerRef = strings.TrimSpace(pay.PayerRef)
		}
		paymentRef = pay.ID.Hex()
	} else {
		if in.Amount == nil {
			return nil, errors.New(errPaymentOrAmount)
		}
		amount = *in.Amount
		paymentRef = fmt.Sprintf("manual:%d", time.Now().UnixNano())
	}
	if amount <= 0 {
		return nil, errors.New(errAmountInvalid)
	}
	if payerRef == "" {
		return nil, errors.New(errPayerRefRequired)
	}
	if currency == "" {
		currency = "TRY"
	}

	requestHash := hashPayload(map[string]interface{}{
		"payment_id": inPaymentID(in.PaymentID),
		"amount":     amount,
		"payer_ref":  payerRef,
		"currency":   currency,
	})
	if in.IdempotencyKey != "" {
		if replay, err := s.idempotencyReplayAllocation(ctx, "allocate_payment", in.IdempotencyKey, requestHash); err != nil || replay != nil {
			return replay, err
		}
	}

	charges, err := s.repo.ListOpenChargesByParty(ctx, payerRef, in.PropertyID)
	if err != nil {
		return nil, err
	}
	if len(charges) == 0 {
		return nil, errors.New(errNoOpenCharges)
	}
	sort.SliceStable(charges, func(i, j int) bool {
		if charges[i].Period == charges[j].Period {
			if charges[i].DueDate.Equal(charges[j].DueDate) {
				return charges[i].CreatedAt.Before(charges[j].CreatedAt)
			}
			return charges[i].DueDate.Before(charges[j].DueDate)
		}
		return charges[i].Period < charges[j].Period
	})

	remaining := amount
	allocCount := 0
	for _, ch := range charges {
		if remaining <= 0 {
			break
		}
		toAllocate := minFloat(remaining, ch.AmountRemaining)
		if toAllocate <= 0 {
			continue
		}
		newPaid := ch.AmountPaid + toAllocate
		newRemaining := ch.AmountRemaining - toAllocate
		newStatus := model.ChargeStatusOpen
		if newRemaining <= 0.000001 {
			newRemaining = 0
			newStatus = model.ChargeStatusPaid
		} else if newPaid > 0 {
			newStatus = model.ChargeStatusPartial
		}
		if err := s.repo.UpdateChargeAmounts(ctx, ch.ID, newPaid, newRemaining, newStatus, in.ActorID); err != nil {
			return nil, err
		}
		_, err := s.repo.InsertAllocation(ctx, model.PaymentAllocation{
			ChargeID:        ch.ID,
			PropertyID:      ch.PropertyID,
			UnitID:          ch.UnitID,
			PartyRef:        payerRef,
			PaymentRef:      paymentRef,
			AllocatedAmount: toAllocate,
			Currency:        currency,
			AllocatedAt:     resolveAllocationTime(in.Date),
			CreatedBy:       in.ActorID,
		})
		if err != nil {
			return nil, err
		}
		s.emitEvent(ctx, "pms.payment.allocated", ch.ID, in.ActorID, map[string]interface{}{"payment_ref": paymentRef, "allocated": toAllocate, "period": ch.Period})
		allocCount++
		remaining -= toAllocate
	}
	allocated := amount - remaining
	res := &model.AllocationResult{
		PaymentRef:      paymentRef,
		RequestedAmount: amount,
		AllocatedAmount: allocated,
		RemainingAmount: remaining,
		Currency:        currency,
		Allocations:     allocCount,
	}
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, in.ActorID, "admin", "pms", "pms.billing.allocate", "payment", nil,
			"PMS payment allocated", map[string]interface{}{"payment_ref": paymentRef, "allocated": allocated, "remaining": remaining, "allocations": allocCount}, "", "")
	}
	if in.IdempotencyKey != "" {
		_ = s.repo.InsertIdempotency(ctx, model.BillingIdempotency{
			Operation:      "allocate_payment",
			IdempotencyKey: in.IdempotencyKey,
			RequestHash:    requestHash,
			Response:       res,
		})
	}
	return res, nil
}

func (s *service) Overview(ctx context.Context, propertyID *primitive.ObjectID, period string) (*model.BillingOverview, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	items, err := s.repo.ListCharges(ctx, propertyID, strings.TrimSpace(period))
	if err != nil {
		return nil, err
	}
	res := &model.BillingOverview{Currency: "TRY", Period: strings.TrimSpace(period), Units: make([]model.UnitOverview, 0)}
	if propertyID != nil && *propertyID != primitive.NilObjectID {
		res.PropertyID = propertyID
	}
	byUnit := make(map[primitive.ObjectID]*model.UnitOverview)
	for _, ch := range items {
		res.TotalDue += ch.AmountDue
		res.TotalPaid += ch.AmountPaid
		res.TotalRemaining += ch.AmountRemaining
		if strings.TrimSpace(ch.Currency) != "" {
			res.Currency = ch.Currency
		}
		u, ok := byUnit[ch.UnitID]
		if !ok {
			u = &model.UnitOverview{UnitID: ch.UnitID, PartyRef: ch.PartyRef}
			byUnit[ch.UnitID] = u
		}
		u.TotalDue += ch.AmountDue
		u.TotalPaid += ch.AmountPaid
		u.TotalRemaining += ch.AmountRemaining
		ts := ch.UpdatedAt
		if u.LastChargeAt == nil || ts.After(*u.LastChargeAt) {
			u.LastChargeAt = &ts
		}
	}
	res.Units = make([]model.UnitOverview, 0, len(byUnit))
	for _, v := range byUnit {
		res.Units = append(res.Units, *v)
	}
	sort.SliceStable(res.Units, func(i, j int) bool {
		return res.Units[i].UnitID.Hex() < res.Units[j].UnitID.Hex()
	})
	return res, nil
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

func parsePeriod(period string) (time.Time, time.Time, error) {
	t, err := time.Parse("2006-01", period)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return start, end, nil
}

func pickResponsible(items []pmsmodel.Occupancy, at time.Time) pmsmodel.Occupancy {
	containsAt := make([]pmsmodel.Occupancy, 0, len(items))
	for _, item := range items {
		if !item.StartAt.After(at) && (item.EndAt == nil || !item.EndAt.Before(at)) {
			containsAt = append(containsAt, item)
		}
	}
	if len(containsAt) == 0 {
		containsAt = items
	}
	for _, item := range containsAt {
		if item.Responsibility.FinancialResponsible {
			return item
		}
	}
	for _, item := range containsAt {
		if item.OccupantType == pmsmodel.OccupantTypeTenant {
			return item
		}
	}
	for _, item := range containsAt {
		if item.OccupantType == pmsmodel.OccupantTypeOwner {
			return item
		}
	}
	return containsAt[0]
}

func resolveAllocationTime(raw *string) time.Time {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return time.Now()
	}
	if t, err := time.Parse(time.RFC3339, strings.TrimSpace(*raw)); err == nil {
		return t
	}
	return time.Now()
}

func (s *service) idempotencyReplayPublish(ctx context.Context, operation, key, requestHash string) (*model.PublishResult, error) {
	rec, err := s.repo.FindIdempotency(ctx, operation, key)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, nil
	}
	if rec.RequestHash != requestHash {
		return nil, errors.New(errIdempotencyConflict)
	}
	var out model.PublishResult
	bytes, _ := json.Marshal(rec.Response)
	if len(bytes) > 0 {
		_ = json.Unmarshal(bytes, &out)
	}
	return &out, nil
}

func (s *service) idempotencyReplayAllocation(ctx context.Context, operation, key, requestHash string) (*model.AllocationResult, error) {
	rec, err := s.repo.FindIdempotency(ctx, operation, key)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, nil
	}
	if rec.RequestHash != requestHash {
		return nil, errors.New(errIdempotencyConflict)
	}
	var out model.AllocationResult
	bytes, _ := json.Marshal(rec.Response)
	if len(bytes) > 0 {
		_ = json.Unmarshal(bytes, &out)
	}
	return &out, nil
}

func (s *service) emitEvent(ctx context.Context, name string, entityID primitive.ObjectID, actor *primitive.ObjectID, diff map[string]interface{}) {
	tenantID, siteID, err := tenantSiteIDs(ctx)
	if err != nil {
		return
	}
	now := time.Now()
	payload := map[string]interface{}{
		"event_id":    primitive.NewObjectID().Hex(),
		"tenant_id":   tenantID,
		"site_id":     siteID,
		"entity_id":   entityID.Hex(),
		"actor_id":    actorHex(actor),
		"occurred_at": now.UTC().Format(time.RFC3339),
		"version":     1,
	}
	if len(diff) > 0 {
		payload["diff"] = diff
	}
	_ = s.repo.InsertOutboxEvent(ctx, model.OutboxEvent{
		TenantID:   tenantID,
		SiteID:     siteID,
		EventName:  name,
		EntityID:   entityID,
		ActorID:    actor,
		OccurredAt: now,
		Version:    1,
		Diff:       diff,
		Payload:    payload,
		Processed:  false,
	})
}

func hashPayload(v interface{}) string {
	b, _ := json.Marshal(v)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func actorHex(actor *primitive.ObjectID) string {
	if actor == nil || *actor == primitive.NilObjectID {
		return ""
	}
	return actor.Hex()
}

func optionalOID(v *primitive.ObjectID) string {
	if v == nil || *v == primitive.NilObjectID {
		return ""
	}
	return v.Hex()
}

func inPaymentID(v *primitive.ObjectID) string {
	if v == nil || *v == primitive.NilObjectID {
		return ""
	}
	return v.Hex()
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

var _ Service = (*service)(nil)
