package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	activitylogs "github.com/aystek/aystek-core/internal/modules/activity_logs/service"
	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/model"
	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/repository"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrStateTransition = errors.New("state_transition_invalid")

const (
	errNotFound         = "pms.work_orders.error.not_found"
	errInvalidPayload   = "pms.work_orders.validation.invalid_payload"
	errPropertyRequired = "pms.work_orders.validation.property_required"
	errTitleRequired    = "pms.work_orders.validation.title_required"
	errCategoryInvalid  = "pms.work_orders.validation.category_invalid"
	errPriorityInvalid  = "pms.work_orders.validation.priority_invalid"
	errStatusInvalid    = "pms.work_orders.validation.status_invalid"
	errAssigneeRequired = "pms.work_orders.validation.assignee_required"
	errCommentRequired  = "pms.work_orders.validation.comment_required"
)

type service struct {
	repo   repository.Repository
	logSvc activitylogs.Service
}

func New(repo repository.Repository, logSvc activitylogs.Service) Service {
	return &service{repo: repo, logSvc: logSvc}
}

func (s *service) Create(ctx context.Context, in CreateInput) (*model.WorkOrderDetail, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	item, err := s.buildCreate(in)
	if err != nil {
		return nil, err
	}
	id, err := s.repo.InsertWorkOrder(ctx, *item)
	if err != nil {
		return nil, err
	}
	_ = s.repo.InsertEvent(ctx, model.WorkOrderEvent{WorkOrderID: id, Type: model.EventCreated, Payload: map[string]interface{}{"status": item.Status}, CreatedBy: in.ActorID})
	s.emitEvent(ctx, "pms.work_order.created", id, in.ActorID, map[string]interface{}{"status": item.Status, "property_id": item.PropertyID.Hex(), "unit_id": oidHex(item.UnitID)})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, in.ActorID, "admin", "pms", "pms.work_orders.create", "pms_work_order", &id, "PMS work order created", map[string]interface{}{"status": item.Status}, "", "")
	}
	return s.Get(ctx, id)
}

func (s *service) List(ctx context.Context, filter repository.ListFilter) ([]model.WorkOrder, int64, error) {
	if err := validateContext(ctx); err != nil {
		return nil, 0, err
	}
	if filter.Status != "" && !filter.Status.IsValid() {
		return nil, 0, errors.New(errStatusInvalid)
	}
	if filter.Priority != "" && !filter.Priority.IsValid() {
		return nil, 0, errors.New(errPriorityInvalid)
	}
	return s.repo.ListWorkOrders(ctx, filter)
}

func (s *service) Get(ctx context.Context, id primitive.ObjectID) (*model.WorkOrderDetail, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	item, err := s.repo.FindWorkOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errNotFound)
		}
		return nil, err
	}
	events, err := s.repo.ListEvents(ctx, id, 10)
	if err != nil {
		return nil, err
	}
	return &model.WorkOrderDetail{WorkOrder: *item, Events: events}, nil
}

func (s *service) Update(ctx context.Context, id primitive.ObjectID, in UpdateInput) (*model.WorkOrderDetail, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	item, err := s.repo.FindWorkOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errNotFound)
		}
		return nil, err
	}
	if strings.TrimSpace(in.Title) == "" {
		return nil, errors.New(errTitleRequired)
	}
	if !in.Category.IsValid() {
		return nil, errors.New(errCategoryInvalid)
	}
	if !in.Priority.IsValid() {
		return nil, errors.New(errPriorityInvalid)
	}
	sla, err := parseOptionalTime(in.SLADueAt)
	if err != nil {
		return nil, errors.New(errInvalidPayload)
	}
	item.Title = strings.TrimSpace(in.Title)
	item.Description = strings.TrimSpace(in.Description)
	item.Category = in.Category
	item.Priority = in.Priority
	item.SLADueAt = sla
	item.Tags = normalizeTags(in.Tags)
	item.Attachments = in.Attachments
	item.CostItems = normalizeCostItems(in.CostItems)
	item.UpdatedBy = in.ActorID
	if err := s.repo.UpdateWorkOrder(ctx, id, *item); err != nil {
		return nil, err
	}
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, in.ActorID, "admin", "pms", "pms.work_orders.update", "pms_work_order", &id, "PMS work order updated", nil, "", "")
	}
	return s.Get(ctx, id)
}

func (s *service) Assign(ctx context.Context, id primitive.ObjectID, assignee *primitive.ObjectID, actor *primitive.ObjectID) (*model.WorkOrderDetail, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	if assignee == nil || *assignee == primitive.NilObjectID {
		return nil, errors.New(errAssigneeRequired)
	}
	item, err := s.repo.FindWorkOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errNotFound)
		}
		return nil, err
	}
	if err := s.repo.UpdateAssignee(ctx, id, assignee, actor); err != nil {
		return nil, err
	}
	_ = s.repo.InsertEvent(ctx, model.WorkOrderEvent{WorkOrderID: id, Type: model.EventAssigned, Payload: map[string]interface{}{"assignee_user_id": assignee.Hex()}, CreatedBy: actor})
	s.emitEvent(ctx, "pms.work_order.assigned", id, actor, map[string]interface{}{"assignee_user_id": assignee.Hex(), "property_id": item.PropertyID.Hex(), "unit_id": oidHex(item.UnitID)})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, actor, "admin", "pms", "pms.work_orders.assign", "pms_work_order", &id, "PMS work order assigned", map[string]interface{}{"assignee_user_id": assignee.Hex()}, "", "")
	}
	return s.Get(ctx, id)
}

func (s *service) ChangeStatus(ctx context.Context, id primitive.ObjectID, status model.WorkOrderStatus, actor *primitive.ObjectID) (*model.WorkOrderDetail, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	if !status.IsValid() {
		return nil, errors.New(errStatusInvalid)
	}
	item, err := s.repo.FindWorkOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errNotFound)
		}
		return nil, err
	}
	if !canTransition(item.Status, status) {
		return nil, fmt.Errorf("%w:%s->%s", ErrStateTransition, item.Status, status)
	}
	var closedAt *time.Time
	if status == model.StatusDone {
		now := time.Now()
		closedAt = &now
	}
	if err := s.repo.UpdateStatus(ctx, id, status, closedAt, actor); err != nil {
		return nil, err
	}
	_ = s.repo.InsertEvent(ctx, model.WorkOrderEvent{WorkOrderID: id, Type: model.EventStatusChanged, Payload: map[string]interface{}{"from": item.Status, "to": status}, CreatedBy: actor})
	s.emitEvent(ctx, "pms.work_order.status_changed", id, actor, map[string]interface{}{"from": item.Status, "to": status, "property_id": item.PropertyID.Hex(), "unit_id": oidHex(item.UnitID)})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, actor, "admin", "pms", "pms.work_orders.status", "pms_work_order", &id, "PMS work order status changed", map[string]interface{}{"from": item.Status, "to": status}, "", "")
	}
	return s.Get(ctx, id)
}

func (s *service) Comment(ctx context.Context, id primitive.ObjectID, comment string, actor *primitive.ObjectID) (*model.WorkOrderDetail, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(comment) == "" {
		return nil, errors.New(errCommentRequired)
	}
	item, err := s.repo.FindWorkOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(errNotFound)
		}
		return nil, err
	}
	_ = s.repo.InsertEvent(ctx, model.WorkOrderEvent{WorkOrderID: id, Type: model.EventCommented, Payload: map[string]interface{}{"comment": strings.TrimSpace(comment)}, CreatedBy: actor})
	s.emitEvent(ctx, "pms.work_order.commented", id, actor, map[string]interface{}{"property_id": item.PropertyID.Hex(), "unit_id": oidHex(item.UnitID)})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, actor, "admin", "pms", "pms.work_orders.comment", "pms_work_order", &id, "PMS work order commented", nil, "", "")
	}
	return s.Get(ctx, id)
}

func (s *service) Close(ctx context.Context, id primitive.ObjectID, actor *primitive.ObjectID) (*model.WorkOrderDetail, error) {
	detail, err := s.ChangeStatus(ctx, id, model.StatusDone, actor)
	if err != nil {
		return nil, err
	}
	_ = s.repo.InsertEvent(ctx, model.WorkOrderEvent{WorkOrderID: id, Type: model.EventClosed, Payload: map[string]interface{}{"status": model.StatusDone}, CreatedBy: actor})
	s.emitEvent(ctx, "pms.work_order.closed", id, actor, map[string]interface{}{"status": model.StatusDone, "property_id": detail.WorkOrder.PropertyID.Hex(), "unit_id": oidHex(detail.WorkOrder.UnitID)})
	if s.logSvc != nil {
		_ = s.logSvc.Log(ctx, actor, "admin", "pms", "pms.work_orders.close", "pms_work_order", &id, "PMS work order closed", nil, "", "")
	}
	return s.Get(ctx, id)
}

func (s *service) buildCreate(in CreateInput) (*model.WorkOrder, error) {
	if in.PropertyID == primitive.NilObjectID {
		return nil, errors.New(errPropertyRequired)
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, errors.New(errTitleRequired)
	}
	if !in.Category.IsValid() {
		return nil, errors.New(errCategoryInvalid)
	}
	if !in.Priority.IsValid() {
		return nil, errors.New(errPriorityInvalid)
	}
	sla, err := parseOptionalTime(in.SLADueAt)
	if err != nil {
		return nil, errors.New(errInvalidPayload)
	}
	return &model.WorkOrder{
		PropertyID:  in.PropertyID,
		UnitID:      in.UnitID,
		Title:       title,
		Description: strings.TrimSpace(in.Description),
		Category:    in.Category,
		Priority:    in.Priority,
		Status:      model.StatusOpen,
		SLADueAt:    sla,
		ReporterRef: strings.TrimSpace(in.ReporterRef),
		Tags:        normalizeTags(in.Tags),
		Attachments: in.Attachments,
		CostItems:   normalizeCostItems(in.CostItems),
		CreatedBy:   in.ActorID,
		UpdatedBy:   in.ActorID,
	}, nil
}

func parseOptionalTime(raw *string) (*time.Time, error) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(*raw))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		v := strings.TrimSpace(tag)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func normalizeCostItems(items []model.CostItem) []model.CostItem {
	out := make([]model.CostItem, 0, len(items))
	for _, it := range items {
		if it.Amount < 0 {
			continue
		}
		it.Currency = strings.ToUpper(strings.TrimSpace(it.Currency))
		if it.Currency == "" {
			it.Currency = "TRY"
		}
		it.Note = strings.TrimSpace(it.Note)
		it.VendorRef = strings.TrimSpace(it.VendorRef)
		out = append(out, it)
	}
	return out
}

func canTransition(from, to model.WorkOrderStatus) bool {
	if from == to {
		return true
	}
	allowed := map[model.WorkOrderStatus]map[model.WorkOrderStatus]bool{
		model.StatusOpen: {
			model.StatusInProgress: true,
			model.StatusBlocked:    true,
			model.StatusCancelled:  true,
		},
		model.StatusInProgress: {
			model.StatusBlocked:   true,
			model.StatusDone:      true,
			model.StatusCancelled: true,
		},
		model.StatusBlocked: {
			model.StatusInProgress: true,
			model.StatusCancelled:  true,
		},
	}
	next, ok := allowed[from]
	if !ok {
		return false
	}
	return next[to]
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

func oidHex(v *primitive.ObjectID) string {
	if v == nil || *v == primitive.NilObjectID {
		return ""
	}
	return v.Hex()
}

var _ Service = (*service)(nil)
