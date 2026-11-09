package service

import (
	"context"
	"testing"

	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/model"
	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/repository"
	"github.com/aystek/aystek-core/internal/platform/contextkeys"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepo struct {
	items []model.UnitListItem
	det   *model.UnitDetail
}

func (m *mockRepo) ListUnits(_ context.Context, filter repository.ListFilter) ([]model.UnitListItem, int64, error) {
	return m.items, int64(len(m.items)), nil
}

func (m *mockRepo) GetUnitDetail(_ context.Context, id primitive.ObjectID) (*model.UnitDetail, error) {
	return m.det, nil
}

func scopedCtx() context.Context {
	ctx := context.Background()
	ctx = contextkeys.WithTenantID(ctx, "tenant-1")
	ctx = contextkeys.WithSiteID(ctx, "site-1")
	return ctx
}

func TestListUnitsAppliesBounds(t *testing.T) {
	repo := &mockRepo{items: []model.UnitListItem{{UnitID: primitive.NewObjectID()}}}
	svc := New(repo)
	items, total, err := svc.ListUnits(scopedCtx(), repository.ListFilter{Page: 0, Limit: 1000})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("unexpected result total=%d len=%d", total, len(items))
	}
}

func TestListUnitsRequiresContext(t *testing.T) {
	repo := &mockRepo{}
	svc := New(repo)
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic for missing tenant/site context")
		}
	}()
	_, _, _ = svc.ListUnits(context.Background(), repository.ListFilter{})
}

func TestGetUnitDetail(t *testing.T) {
	detail := &model.UnitDetail{UnitListItem: model.UnitListItem{UnitID: primitive.NewObjectID()}}
	repo := &mockRepo{det: detail}
	svc := New(repo)
	got, err := svc.GetUnitDetail(scopedCtx(), primitive.NewObjectID())
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got == nil || got.UnitID == primitive.NilObjectID {
		t.Fatalf("expected detail")
	}
}
