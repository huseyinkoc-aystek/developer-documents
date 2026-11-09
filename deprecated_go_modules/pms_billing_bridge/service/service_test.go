package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	pmsmodel "github.com/aystek/aystek-core/internal/modules/pms/model"
	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/model"
	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/repository"
	"github.com/aystek/aystek-core/internal/platform/contextkeys"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockRepo struct {
	property    *pmsmodel.Property
	occupancies []pmsmodel.Occupancy
	charges     []model.PMSCharge
	idem        map[string]model.BillingIdempotency
	activePlan  *model.DuesPlanReference
	unitAreas   map[primitive.ObjectID]float64
}

func newMockRepo() *mockRepo {
	return &mockRepo{idem: map[string]model.BillingIdempotency{}, unitAreas: map[primitive.ObjectID]float64{}}
}

func idemKey(op, key string) string { return op + ":" + key }

func (m *mockRepo) FindProperty(_ context.Context, id primitive.ObjectID) (*pmsmodel.Property, error) {
	if m.property == nil || m.property.ID != id {
		return nil, mongo.ErrNoDocuments
	}
	return m.property, nil
}

func (m *mockRepo) FindActiveOccupanciesByPropertyPeriod(_ context.Context, propertyID primitive.ObjectID, from, to time.Time) ([]pmsmodel.Occupancy, error) {
	items := make([]pmsmodel.Occupancy, 0)
	for _, occ := range m.occupancies {
		if occ.PropertyID != propertyID {
			continue
		}
		if occ.StartAt.After(to) {
			continue
		}
		if occ.EndAt != nil && occ.EndAt.Before(from) {
			continue
		}
		items = append(items, occ)
	}
	return items, nil
}

func (m *mockRepo) FindDuesPlanByID(_ context.Context, id primitive.ObjectID) (*model.DuesPlanReference, error) {
	if m.activePlan != nil && m.activePlan.ID == id {
		return m.activePlan, nil
	}
	return nil, mongo.ErrNoDocuments
}

func (m *mockRepo) FindActiveDuesPlan(_ context.Context, propertyID primitive.ObjectID, at time.Time) (*model.DuesPlanReference, error) {
	if m.activePlan != nil && m.activePlan.PropertyID == propertyID {
		return m.activePlan, nil
	}
	return nil, mongo.ErrNoDocuments
}

func (m *mockRepo) FindUnitAreaSQM(_ context.Context, unitID primitive.ObjectID) (float64, error) {
	if v, ok := m.unitAreas[unitID]; ok {
		return v, nil
	}
	return 0, mongo.ErrNoDocuments
}

func (m *mockRepo) FindChargeByUnitPeriod(_ context.Context, propertyID, unitID primitive.ObjectID, period string) (*model.PMSCharge, error) {
	for i := range m.charges {
		if m.charges[i].PropertyID == propertyID && m.charges[i].UnitID == unitID && m.charges[i].Period == period {
			return &m.charges[i], nil
		}
	}
	return nil, mongo.ErrNoDocuments
}

func (m *mockRepo) InsertCharge(_ context.Context, charge model.PMSCharge) (primitive.ObjectID, error) {
	for _, ch := range m.charges {
		if ch.PropertyID == charge.PropertyID && ch.UnitID == charge.UnitID && ch.Period == charge.Period {
			return primitive.NilObjectID, fmt.Errorf("duplicate")
		}
	}
	charge.ID = primitive.NewObjectID()
	if charge.AmountRemaining == 0 && charge.AmountDue > 0 {
		charge.AmountRemaining = charge.AmountDue
	}
	m.charges = append(m.charges, charge)
	return charge.ID, nil
}

func (m *mockRepo) UpdateChargeAmounts(_ context.Context, id primitive.ObjectID, paid, remaining float64, status model.ChargeStatus, actor *primitive.ObjectID) error {
	for i := range m.charges {
		if m.charges[i].ID == id {
			m.charges[i].AmountPaid = paid
			m.charges[i].AmountRemaining = remaining
			m.charges[i].Status = status
			m.charges[i].UpdatedBy = actor
			return nil
		}
	}
	return mongo.ErrNoDocuments
}

func (m *mockRepo) ListOpenChargesByParty(_ context.Context, partyRef string, propertyID *primitive.ObjectID) ([]model.PMSCharge, error) {
	items := make([]model.PMSCharge, 0)
	for _, ch := range m.charges {
		if ch.PartyRef != partyRef || ch.AmountRemaining <= 0 {
			continue
		}
		if propertyID != nil && *propertyID != primitive.NilObjectID && ch.PropertyID != *propertyID {
			continue
		}
		items = append(items, ch)
	}
	return items, nil
}

func (m *mockRepo) ListCharges(_ context.Context, propertyID *primitive.ObjectID, period string) ([]model.PMSCharge, error) {
	items := make([]model.PMSCharge, 0)
	for _, ch := range m.charges {
		if propertyID != nil && *propertyID != primitive.NilObjectID && ch.PropertyID != *propertyID {
			continue
		}
		if period != "" && ch.Period != period {
			continue
		}
		items = append(items, ch)
	}
	return items, nil
}

func (m *mockRepo) InsertAllocation(_ context.Context, allocation model.PaymentAllocation) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (m *mockRepo) FindPaymentByID(_ context.Context, id primitive.ObjectID) (*repository.PaymentRef, error) {
	return nil, errors.New("not implemented")
}

func (m *mockRepo) FindIdempotency(_ context.Context, operation, key string) (*model.BillingIdempotency, error) {
	if rec, ok := m.idem[idemKey(operation, key)]; ok {
		return &rec, nil
	}
	return nil, nil
}

func (m *mockRepo) InsertIdempotency(_ context.Context, rec model.BillingIdempotency) error {
	m.idem[idemKey(rec.Operation, rec.IdempotencyKey)] = rec
	return nil
}

func (m *mockRepo) InsertOutboxEvent(_ context.Context, event model.OutboxEvent) error { return nil }

func testContext() context.Context {
	ctx := context.Background()
	ctx = contextkeys.WithTenantID(ctx, "tenant-a")
	ctx = contextkeys.WithSiteID(ctx, "site-a")
	return ctx
}

func TestPublishIdempotentSameKey(t *testing.T) {
	repo := newMockRepo()
	propertyID := primitive.NewObjectID()
	unitID := primitive.NewObjectID()
	repo.property = &pmsmodel.Property{ID: propertyID}
	repo.occupancies = []pmsmodel.Occupancy{{
		ID:             primitive.NewObjectID(),
		PropertyID:     propertyID,
		UnitID:         unitID,
		PartyRef:       "party-1",
		StartAt:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		OccupantType:   pmsmodel.OccupantTypeTenant,
		Responsibility: pmsmodel.OccupancyResponsibility{FinancialResponsible: true},
	}}
	svc := New(repo, nil)

	amount := 100.0
	in := PublishInput{PropertyID: propertyID, Period: "2026-01", DefaultAmount: &amount, IdempotencyKey: "same-key"}
	first, err := svc.PublishDues(testContext(), in)
	if err != nil {
		t.Fatalf("first publish error: %v", err)
	}
	second, err := svc.PublishDues(testContext(), in)
	if err != nil {
		t.Fatalf("second publish error: %v", err)
	}
	if first.CreatedCount != second.CreatedCount || first.SkippedCount != second.SkippedCount {
		t.Fatalf("expected idempotent same result, got first=%+v second=%+v", first, second)
	}
	if len(repo.charges) != 1 {
		t.Fatalf("expected single charge after idempotent replay, got %d", len(repo.charges))
	}
}

func TestPublishDuplicationGuardWithoutIdempotency(t *testing.T) {
	repo := newMockRepo()
	propertyID := primitive.NewObjectID()
	unitID := primitive.NewObjectID()
	repo.property = &pmsmodel.Property{ID: propertyID}
	repo.occupancies = []pmsmodel.Occupancy{{
		ID:           primitive.NewObjectID(),
		PropertyID:   propertyID,
		UnitID:       unitID,
		PartyRef:     "party-2",
		StartAt:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		OccupantType: pmsmodel.OccupantTypeOwner,
	}}
	svc := New(repo, nil)

	amount := 120.0
	_, err := svc.PublishDues(testContext(), PublishInput{PropertyID: propertyID, Period: "2026-02", DefaultAmount: &amount})
	if err != nil {
		t.Fatalf("first publish error: %v", err)
	}
	second, err := svc.PublishDues(testContext(), PublishInput{PropertyID: propertyID, Period: "2026-02", DefaultAmount: &amount})
	if err != nil {
		t.Fatalf("second publish error: %v", err)
	}
	if second.CreatedCount != 0 || second.SkippedCount != 1 {
		t.Fatalf("expected duplicate guard to skip existing charge, got %+v", second)
	}
}

func TestAllocateFIFOWithPartial(t *testing.T) {
	repo := newMockRepo()
	propertyID := primitive.NewObjectID()
	unitID := primitive.NewObjectID()
	repo.charges = []model.PMSCharge{
		{ID: primitive.NewObjectID(), PropertyID: propertyID, UnitID: unitID, PartyRef: "payer-1", Period: "2026-01", AmountDue: 100, AmountRemaining: 100, Currency: "TRY", DueDate: time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC), Status: model.ChargeStatusOpen},
		{ID: primitive.NewObjectID(), PropertyID: propertyID, UnitID: unitID, PartyRef: "payer-1", Period: "2026-02", AmountDue: 100, AmountRemaining: 100, Currency: "TRY", DueDate: time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC), Status: model.ChargeStatusOpen},
	}
	svc := New(repo, nil)
	pay := 150.0
	out, err := svc.AllocatePayment(testContext(), AllocateInput{Amount: &pay, PayerRef: "payer-1", Currency: "TRY"})
	if err != nil {
		t.Fatalf("allocate error: %v", err)
	}
	if out.AllocatedAmount != 150 || out.RemainingAmount != 0 || out.Allocations != 2 {
		t.Fatalf("unexpected allocation result: %+v", out)
	}
	if repo.charges[0].AmountRemaining != 0 || repo.charges[0].Status != model.ChargeStatusPaid {
		t.Fatalf("expected first charge paid, got %+v", repo.charges[0])
	}
	if repo.charges[1].AmountRemaining != 50 || repo.charges[1].Status != model.ChargeStatusPartial {
		t.Fatalf("expected second charge partial, got %+v", repo.charges[1])
	}
}

func TestPublishResponsibilityMappingOnBoundary(t *testing.T) {
	repo := newMockRepo()
	propertyID := primitive.NewObjectID()
	unitID := primitive.NewObjectID()
	repo.property = &pmsmodel.Property{ID: propertyID}
	endJan := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	repo.occupancies = []pmsmodel.Occupancy{
		{
			ID:           primitive.NewObjectID(),
			PropertyID:   propertyID,
			UnitID:       unitID,
			PartyRef:     "tenant-old",
			OccupantType: pmsmodel.OccupantTypeTenant,
			StartAt:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			EndAt:        &endJan,
		},
		{
			ID:             primitive.NewObjectID(),
			PropertyID:     propertyID,
			UnitID:         unitID,
			PartyRef:       "owner-new",
			OccupantType:   pmsmodel.OccupantTypeOwner,
			StartAt:        time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			Responsibility: pmsmodel.OccupancyResponsibility{FinancialResponsible: true},
		},
	}
	svc := New(repo, nil)
	amount := 80.0
	_, err := svc.PublishDues(testContext(), PublishInput{PropertyID: propertyID, Period: "2026-02", DefaultAmount: &amount})
	if err != nil {
		t.Fatalf("publish error: %v", err)
	}
	if len(repo.charges) != 1 {
		t.Fatalf("expected one charge, got %d", len(repo.charges))
	}
	if repo.charges[0].PartyRef != "owner-new" {
		t.Fatalf("expected responsible owner-new on period boundary, got %s", repo.charges[0].PartyRef)
	}
}

func TestPublishUsesActivePlanWhenPlanMissing(t *testing.T) {
	repo := newMockRepo()
	propertyID := primitive.NewObjectID()
	unitID := primitive.NewObjectID()
	planID := primitive.NewObjectID()
	repo.property = &pmsmodel.Property{ID: propertyID}
	repo.activePlan = &model.DuesPlanReference{
		ID:         planID,
		PropertyID: propertyID,
		Method:     model.DuesPlanMethodFixed,
		Amount:     175,
		Currency:   "TRY",
	}
	repo.occupancies = []pmsmodel.Occupancy{{
		ID:             primitive.NewObjectID(),
		PropertyID:     propertyID,
		UnitID:         unitID,
		PartyRef:       "party-7",
		StartAt:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		OccupantType:   pmsmodel.OccupantTypeTenant,
		Responsibility: pmsmodel.OccupancyResponsibility{FinancialResponsible: true},
	}}
	svc := New(repo, nil)

	out, err := svc.PublishDues(testContext(), PublishInput{PropertyID: propertyID, Period: "2026-03"})
	if err != nil {
		t.Fatalf("publish error: %v", err)
	}
	if out.PlanID == nil || *out.PlanID != planID {
		t.Fatalf("expected active plan id in result")
	}
	if len(repo.charges) != 1 {
		t.Fatalf("expected one charge, got %d", len(repo.charges))
	}
	if repo.charges[0].AmountDue != 175 {
		t.Fatalf("expected fixed amount from active plan, got %f", repo.charges[0].AmountDue)
	}
}
