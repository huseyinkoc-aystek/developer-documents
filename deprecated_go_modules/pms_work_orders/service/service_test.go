package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/model"
	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/repository"
	"github.com/aystek/aystek-core/internal/platform/contextkeys"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockRepo struct {
	workOrders map[primitive.ObjectID]model.WorkOrder
	events     []model.WorkOrderEvent
	outbox     []model.OutboxEvent
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		workOrders: map[primitive.ObjectID]model.WorkOrder{},
		events:     make([]model.WorkOrderEvent, 0),
		outbox:     make([]model.OutboxEvent, 0),
	}
}

func (m *mockRepo) InsertWorkOrder(_ context.Context, item model.WorkOrder) (primitive.ObjectID, error) {
	item.ID = primitive.NewObjectID()
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	m.workOrders[item.ID] = item
	return item.ID, nil
}

func (m *mockRepo) FindWorkOrderByID(_ context.Context, id primitive.ObjectID) (*model.WorkOrder, error) {
	item, ok := m.workOrders[id]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	copy := item
	return &copy, nil
}

func (m *mockRepo) ListWorkOrders(_ context.Context, filter repository.ListFilter) ([]model.WorkOrder, int64, error) {
	items := make([]model.WorkOrder, 0)
	for _, item := range m.workOrders {
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		if filter.Priority != "" && item.Priority != filter.Priority {
			continue
		}
		if filter.PropertyID != nil && *filter.PropertyID != primitive.NilObjectID && item.PropertyID != *filter.PropertyID {
			continue
		}
		if filter.UnitID != nil && *filter.UnitID != primitive.NilObjectID {
			if item.UnitID == nil || *item.UnitID != *filter.UnitID {
				continue
			}
		}
		items = append(items, item)
	}
	return items, int64(len(items)), nil
}

func (m *mockRepo) UpdateWorkOrder(_ context.Context, id primitive.ObjectID, item model.WorkOrder) error {
	current, ok := m.workOrders[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	current.Title = item.Title
	current.Description = item.Description
	current.Category = item.Category
	current.Priority = item.Priority
	current.SLADueAt = item.SLADueAt
	current.Tags = item.Tags
	current.Attachments = item.Attachments
	current.CostItems = item.CostItems
	current.UpdatedBy = item.UpdatedBy
	current.UpdatedAt = time.Now()
	m.workOrders[id] = current
	return nil
}

func (m *mockRepo) UpdateAssignee(_ context.Context, id primitive.ObjectID, assignee *primitive.ObjectID, actor *primitive.ObjectID) error {
	current, ok := m.workOrders[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	current.AssigneeUserID = assignee
	current.UpdatedBy = actor
	current.UpdatedAt = time.Now()
	m.workOrders[id] = current
	return nil
}

func (m *mockRepo) UpdateStatus(_ context.Context, id primitive.ObjectID, status model.WorkOrderStatus, closedAt *time.Time, actor *primitive.ObjectID) error {
	current, ok := m.workOrders[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	current.Status = status
	current.ClosedAt = closedAt
	current.UpdatedBy = actor
	current.UpdatedAt = time.Now()
	m.workOrders[id] = current
	return nil
}

func (m *mockRepo) InsertEvent(_ context.Context, event model.WorkOrderEvent) error {
	if event.ID == primitive.NilObjectID {
		event.ID = primitive.NewObjectID()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	m.events = append(m.events, event)
	return nil
}

func (m *mockRepo) ListEvents(_ context.Context, workOrderID primitive.ObjectID, limit int64) ([]model.WorkOrderEvent, error) {
	items := make([]model.WorkOrderEvent, 0)
	for _, event := range m.events {
		if event.WorkOrderID == workOrderID {
			items = append(items, event)
		}
	}
	if limit > 0 && int64(len(items)) > limit {
		return items[len(items)-int(limit):], nil
	}
	return items, nil
}

func (m *mockRepo) InsertOutboxEvent(_ context.Context, event model.OutboxEvent) error {
	m.outbox = append(m.outbox, event)
	return nil
}

func svcCtx() context.Context {
	ctx := context.Background()
	ctx = contextkeys.WithTenantID(ctx, "tenant-1")
	ctx = contextkeys.WithSiteID(ctx, "site-1")
	return ctx
}

func TestCreateWorkOrder(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo, nil)
	propertyID := primitive.NewObjectID()
	out, err := svc.Create(svcCtx(), CreateInput{
		PropertyID: propertyID,
		Title:      "Leak on 3rd floor",
		Category:   model.CategoryPlumbing,
		Priority:   model.PriorityHigh,
		Tags:       []string{"urgent"},
	})
	if err != nil {
		t.Fatalf("create returned error: %v", err)
	}
	if out.WorkOrder.Status != model.StatusOpen {
		t.Fatalf("expected status open, got %s", out.WorkOrder.Status)
	}
	if len(repo.outbox) == 0 || repo.outbox[0].EventName != "pms.work_order.created" {
		t.Fatalf("expected created outbox event")
	}
}

func TestInvalidStatusTransitionReturnsErrStateTransition(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo, nil)
	propertyID := primitive.NewObjectID()
	created, err := svc.Create(svcCtx(), CreateInput{
		PropertyID: propertyID,
		Title:      "Security gate issue",
		Category:   model.CategorySecurity,
		Priority:   model.PriorityNormal,
	})
	if err != nil {
		t.Fatalf("create returned error: %v", err)
	}
	_, err = svc.ChangeStatus(svcCtx(), created.WorkOrder.ID, model.StatusDone, nil)
	if err == nil {
		t.Fatalf("expected transition error")
	}
	if !errors.Is(err, ErrStateTransition) {
		t.Fatalf("expected ErrStateTransition, got %v", err)
	}
}

func TestAssignCreatesEventAndOutbox(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo, nil)
	propertyID := primitive.NewObjectID()
	created, err := svc.Create(svcCtx(), CreateInput{
		PropertyID: propertyID,
		Title:      "Elevator maintenance",
		Category:   model.CategoryOther,
		Priority:   model.PriorityNormal,
	})
	if err != nil {
		t.Fatalf("create returned error: %v", err)
	}
	assignee := primitive.NewObjectID()
	out, err := svc.Assign(svcCtx(), created.WorkOrder.ID, &assignee, nil)
	if err != nil {
		t.Fatalf("assign returned error: %v", err)
	}
	if out.WorkOrder.AssigneeUserID == nil || *out.WorkOrder.AssigneeUserID != assignee {
		t.Fatalf("assignee was not persisted")
	}
	hasAssignedOutbox := false
	for _, ev := range repo.outbox {
		if ev.EventName == "pms.work_order.assigned" {
			hasAssignedOutbox = true
			break
		}
	}
	if !hasAssignedOutbox {
		t.Fatalf("expected assigned outbox event")
	}
}

func TestCloseSetsClosedAtAndEmitsEvent(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo, nil)
	propertyID := primitive.NewObjectID()
	created, err := svc.Create(svcCtx(), CreateInput{
		PropertyID: propertyID,
		Title:      "Generator check",
		Category:   model.CategoryElectric,
		Priority:   model.PriorityHigh,
	})
	if err != nil {
		t.Fatalf("create returned error: %v", err)
	}
	_, err = svc.ChangeStatus(svcCtx(), created.WorkOrder.ID, model.StatusInProgress, nil)
	if err != nil {
		t.Fatalf("status in_progress returned error: %v", err)
	}
	closed, err := svc.Close(svcCtx(), created.WorkOrder.ID, nil)
	if err != nil {
		t.Fatalf("close returned error: %v", err)
	}
	if closed.WorkOrder.Status != model.StatusDone || closed.WorkOrder.ClosedAt == nil {
		t.Fatalf("expected done with closed_at set")
	}
	hasClosedOutbox := false
	for _, ev := range repo.outbox {
		if ev.EventName == "pms.work_order.closed" {
			hasClosedOutbox = true
			break
		}
	}
	if !hasClosedOutbox {
		t.Fatalf("expected closed outbox event")
	}
}

func TestListFiltersByStatus(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo, nil)
	propertyID := primitive.NewObjectID()
	first, _ := svc.Create(svcCtx(), CreateInput{
		PropertyID: propertyID,
		Title:      "Open task",
		Category:   model.CategoryCleaning,
		Priority:   model.PriorityLow,
	})
	second, _ := svc.Create(svcCtx(), CreateInput{
		PropertyID: propertyID,
		Title:      "Task to cancel",
		Category:   model.CategoryOther,
		Priority:   model.PriorityLow,
	})
	_, _ = svc.ChangeStatus(svcCtx(), second.WorkOrder.ID, model.StatusCancelled, nil)
	list, total, err := svc.List(svcCtx(), repository.ListFilter{Status: model.StatusOpen, Page: 1, Limit: 25})
	if err != nil {
		t.Fatalf("list returned error: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].ID != first.WorkOrder.ID {
		t.Fatalf("expected only open item, got total=%d len=%d", total, len(list))
	}
}
