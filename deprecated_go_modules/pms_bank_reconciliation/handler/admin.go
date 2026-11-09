package handler

import (
	"errors"
	"strings"

	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/dto"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/model"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/repository"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/service"
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

func (h *AdminHandler) Ingest(c *gin.Context) {
	var req dto.IngestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.reconciliation.validation.invalid_payload", err.Error())
		return
	}
	out, created, err := h.svc.Ingest(c.Request.Context(), service.IngestInput{
		ExternalTxnID: strings.TrimSpace(req.ExternalTxnID),
		BankAccountID: strings.TrimSpace(req.BankAccountID),
		TxnDate:       strings.TrimSpace(req.TxnDate),
		Amount:        req.Amount,
		Currency:      strings.TrimSpace(req.Currency),
		Direction:     model.TxnDirection(strings.TrimSpace(req.Direction)),
		Description:   strings.TrimSpace(req.Description),
		Counterparty:  strings.TrimSpace(req.Counterparty),
		ActorID:       actorID(c),
	})
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if created {
		response.Created(c, out, "")
		return
	}
	response.OK(c, out, "")
}

func (h *AdminHandler) Queue(c *gin.Context) {
	page, limit := parsePageLimit(c)
	items, total, err := h.svc.ListQueue(c.Request.Context(), repository.ListFilter{
		Page:   page,
		Limit:  limit,
		Status: model.TxnStatus(strings.TrimSpace(c.Query("status"))),
		Q:      strings.TrimSpace(c.Query("q")),
	})
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.WithPagination(c, items, pagination.New(page, limit, total), "")
}

func (h *AdminHandler) Propose(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.reconciliation.validation.invalid_id", nil)
		return
	}
	out, err := h.svc.ProposeMatch(c.Request.Context(), id)
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

func (h *AdminHandler) Confirm(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.reconciliation.validation.invalid_id", nil)
		return
	}
	var req dto.ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "pms.reconciliation.validation.invalid_payload", err.Error())
		return
	}
	var candidateID *primitive.ObjectID
	if raw := strings.TrimSpace(req.CandidateChargeID); raw != "" {
		oid, parseErr := primitive.ObjectIDFromHex(raw)
		if parseErr != nil {
			response.BadRequest(c, "pms.reconciliation.validation.invalid_payload", nil)
			return
		}
		candidateID = &oid
	}
	out, err := h.svc.Confirm(c.Request.Context(), id, service.ConfirmInput{
		CandidateChargeID: candidateID,
		CreatePayment:     req.CreatePayment,
		ActorID:           actorID(c),
	})
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

func (h *AdminHandler) Reject(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.reconciliation.validation.invalid_id", nil)
		return
	}
	var req dto.RejectRequest
	_ = c.ShouldBindJSON(&req)
	out, err := h.svc.Reject(c.Request.Context(), id, strings.TrimSpace(req.Reason), actorID(c))
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
