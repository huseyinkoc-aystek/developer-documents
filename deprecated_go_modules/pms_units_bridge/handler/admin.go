package handler

import (
	"strings"

	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/repository"
	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/service"
	"github.com/aystek/aystek-core/internal/platform/pagination"
	"github.com/aystek/aystek-core/internal/platform/response"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminHandler struct {
	svc service.Service
}

func NewAdminHandler(svc service.Service) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) List(c *gin.Context) {
	page, limit := parsePageLimit(c)
	filter := repository.ListFilter{
		Page:  page,
		Limit: limit,
		Q:     strings.TrimSpace(c.Query("q")),
	}
	if raw := strings.TrimSpace(c.Query("property_id")); raw != "" {
		propertyID, err := primitive.ObjectIDFromHex(raw)
		if err != nil {
			response.BadRequest(c, "pms.units.validation.property_invalid", nil)
			return
		}
		filter.PropertyID = &propertyID
	}
	items, total, err := h.svc.ListUnits(c.Request.Context(), filter)
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.WithPagination(c, items, pagination.New(page, limit, total), "")
}

func (h *AdminHandler) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(c.Param("id")))
	if err != nil {
		response.BadRequest(c, "pms.units.validation.invalid_id", nil)
		return
	}
	item, err := h.svc.GetUnitDetail(c.Request.Context(), id)
	if err != nil {
		if handleContextError(c, err) {
			return
		}
		if err == mongo.ErrNoDocuments {
			response.NotFound(c, "pms.units.error.not_found", nil)
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, item, "")
}
