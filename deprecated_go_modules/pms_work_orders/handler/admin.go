package handler

import (
	"errors"
	"strings"

	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/dto"
	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/model"
	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/repository"
	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/service"
	"github.com/aystek/aystek-core/internal/platform/pagination"
	"github.com/aystek/aystek-core/internal/platform/response"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminHandler struct {
	svc service.Service
}

func NewAdminHandler(svc service.Service) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) List(c *gin.Context) {
	page, limit := parsePageLimit(c)
	filter := repository.ListFilter{Page: page, Limit: limit, Q: strings.TrimSpace(c.Query("q")), Status: model.WorkOrderStatus(strings.TrimSpace(c.Query("status"))), Priority: model.WorkOrderPriority(strings.TrimSpace(c.Query("priority")))}
	if raw := strings.TrimSpace(c.Query("property_id")); raw != "" {
		oid, err := primitive.ObjectIDFromHex(raw)
		if err == nil {
			filter.PropertyID = &oid
		}
	}
	if raw := strings.TrimSpace(c.Query("unit_id")); raw != "" {
		oid, err := primitive.ObjectIDFromHex(raw)
		if err == nil {
			filter.UnitID = &oid
		}
	}
	items, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.WithPagination(c, items, pagination.New(page, limit, total), "")
}

func (h *AdminHandler) Create(c *gin.Context) {
	var req dto.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_payload", err.Error())
		return
	}
	propertyID, err := primitive.ObjectIDFromHex(strings.TrimSpace(req.PropertyID))
	if err != nil {
		response.BadRequest(c, "pms.work_orders.validation.property_required", nil)
		return
	}
	var unitID *primitive.ObjectID
	if req.UnitID != nil && strings.TrimSpace(*req.UnitID) != "" {
		oid, parseErr := primitive.ObjectIDFromHex(strings.TrimSpace(*req.UnitID))
		if parseErr != nil {
			response.BadRequest(c, "pms.work_orders.validation.invalid_payload", nil)
			return
		}
		unitID = &oid
	}
	attachments := make([]model.Attachment, 0, len(req.Attachments))
	for _, a := range req.Attachments {
		attachments = append(attachments, model.Attachment{MediaID: strings.TrimSpace(a.MediaID), FileName: strings.TrimSpace(a.FileName), MimeType: strings.TrimSpace(a.MimeType), SizeBytes: a.SizeBytes})
	}
	costItems := make([]model.CostItem, 0, len(req.CostItems))
	for _, ci := range req.CostItems {
		costItems = append(costItems, model.CostItem{Amount: ci.Amount, Currency: strings.TrimSpace(ci.Currency), Note: strings.TrimSpace(ci.Note), VendorRef: strings.TrimSpace(ci.VendorRef)})
	}
	out, err := h.svc.Create(c.Request.Context(), service.CreateInput{PropertyID: propertyID, UnitID: unitID, Title: req.Title, Description: req.Description, Category: model.WorkOrderCategory(req.Category), Priority: model.WorkOrderPriority(req.Priority), SLADueAt: req.SLADueAt, ReporterRef: req.ReporterRef, Tags: req.Tags, Attachments: attachments, CostItems: costItems, ActorID: actorID(c)})
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, out, "")
}

func (h *AdminHandler) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_id", nil)
		return
	}
	out, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}

func (h *AdminHandler) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_id", nil)
		return
	}
	var req dto.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_payload", err.Error())
		return
	}
	attachments := make([]model.Attachment, 0, len(req.Attachments))
	for _, a := range req.Attachments {
		attachments = append(attachments, model.Attachment{MediaID: strings.TrimSpace(a.MediaID), FileName: strings.TrimSpace(a.FileName), MimeType: strings.TrimSpace(a.MimeType), SizeBytes: a.SizeBytes})
	}
	costItems := make([]model.CostItem, 0, len(req.CostItems))
	for _, ci := range req.CostItems {
		costItems = append(costItems, model.CostItem{Amount: ci.Amount, Currency: strings.TrimSpace(ci.Currency), Note: strings.TrimSpace(ci.Note), VendorRef: strings.TrimSpace(ci.VendorRef)})
	}
	out, err := h.svc.Update(c.Request.Context(), id, service.UpdateInput{Title: req.Title, Description: req.Description, Category: model.WorkOrderCategory(req.Category), Priority: model.WorkOrderPriority(req.Priority), SLADueAt: req.SLADueAt, Tags: req.Tags, Attachments: attachments, CostItems: costItems, ActorID: actorID(c)})
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}

func (h *AdminHandler) Assign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_id", nil)
		return
	}
	var req dto.AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_payload", err.Error())
		return
	}
	assignee, err := primitive.ObjectIDFromHex(strings.TrimSpace(req.AssigneeUserID))
	if err != nil {
		response.BadRequest(c, "pms.work_orders.validation.assignee_required", nil)
		return
	}
	out, err := h.svc.Assign(c.Request.Context(), id, &assignee, actorID(c))
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}

func (h *AdminHandler) Status(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_id", nil)
		return
	}
	var req dto.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_payload", err.Error())
		return
	}
	out, err := h.svc.ChangeStatus(c.Request.Context(), id, model.WorkOrderStatus(strings.TrimSpace(req.Status)), actorID(c))
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		if errors.Is(err, service.ErrStateTransition) {
			response.Error(c, 409, err.Error(), nil)
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}

func (h *AdminHandler) Comment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_id", nil)
		return
	}
	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_payload", err.Error())
		return
	}
	out, err := h.svc.Comment(c.Request.Context(), id, strings.TrimSpace(req.Comment), actorID(c))
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}

func (h *AdminHandler) Close(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.work_orders.validation.invalid_id", nil)
		return
	}
	out, err := h.svc.Close(c.Request.Context(), id, actorID(c))
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		if errors.Is(err, service.ErrStateTransition) {
			response.Error(c, 409, err.Error(), nil)
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}
