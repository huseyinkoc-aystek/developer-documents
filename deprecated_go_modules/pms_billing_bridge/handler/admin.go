package handler

import (
	"strings"

	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/dto"
	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/service"
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

func (h *AdminHandler) PublishDues(c *gin.Context) {
	var req dto.PublishDuesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.billing.validation.invalid_payload", err.Error())
		return
	}
	propertyID, err := primitive.ObjectIDFromHex(strings.TrimSpace(req.PropertyID))
	if err != nil {
		response.BadRequest(c, "pms.billing.validation.property_required", nil)
		return
	}
	var planID *primitive.ObjectID
	if req.PlanID != nil && strings.TrimSpace(*req.PlanID) != "" {
		oid, parseErr := primitive.ObjectIDFromHex(strings.TrimSpace(*req.PlanID))
		if parseErr != nil {
			response.BadRequest(c, "pms.billing.validation.invalid_payload", nil)
			return
		}
		planID = &oid
	}
	out, err := h.svc.PublishDues(c.Request.Context(), service.PublishInput{
		PropertyID:     propertyID,
		Period:         strings.TrimSpace(req.Period),
		PlanID:         planID,
		DefaultAmount:  req.DefaultAmount,
		Currency:       strings.TrimSpace(req.Currency),
		IdempotencyKey: strings.TrimSpace(c.GetHeader("Idempotency-Key")),
		ActorID:        actorID(c),
	})
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}

func (h *AdminHandler) AllocatePayment(c *gin.Context) {
	var req dto.AllocatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.billing.validation.invalid_payload", err.Error())
		return
	}
	var paymentID *primitive.ObjectID
	if raw := strings.TrimSpace(req.PaymentID); raw != "" {
		oid, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			response.BadRequest(c, "pms.billing.validation.invalid_payload", nil)
			return
		}
		paymentID = &oid
	}
	var propertyID *primitive.ObjectID
	if raw := strings.TrimSpace(req.PropertyID); raw != "" {
		oid, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			response.BadRequest(c, "pms.billing.validation.invalid_payload", nil)
			return
		}
		propertyID = &oid
	}
	out, err := h.svc.AllocatePayment(c.Request.Context(), service.AllocateInput{
		PaymentID:      paymentID,
		Amount:         req.Amount,
		PayerRef:       strings.TrimSpace(req.PayerRef),
		Date:           req.Date,
		Currency:       strings.TrimSpace(req.Currency),
		PropertyID:     propertyID,
		IdempotencyKey: strings.TrimSpace(c.GetHeader("Idempotency-Key")),
		ActorID:        actorID(c),
	})
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}

func (h *AdminHandler) Overview(c *gin.Context) {
	period := strings.TrimSpace(c.Query("period"))
	var propertyID *primitive.ObjectID
	if raw := strings.TrimSpace(c.Query("property_id")); raw != "" {
		oid, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			response.BadRequest(c, "pms.billing.validation.invalid_payload", nil)
			return
		}
		propertyID = &oid
	}
	out, err := h.svc.Overview(c.Request.Context(), propertyID, period)
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, out, "")
}
