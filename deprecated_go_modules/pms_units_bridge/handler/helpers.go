package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/aystek/aystek-core/internal/platform/response"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"github.com/gin-gonic/gin"
)

func parsePageLimit(c *gin.Context) (int64, int64) {
	page := int64(1)
	limit := int64(25)
	if raw := strings.TrimSpace(c.Query("page")); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func handleContextError(c *gin.Context, err error) bool {
	if errors.Is(err, tenantfilter.ErrTenantIDMissing) {
		response.BadRequest(c, "errors.tenant.missing", gin.H{"code": "TENANT_ID_MISSING", "message_key": "errors.tenant.missing"})
		return true
	}
	if errors.Is(err, tenantfilter.ErrTenantIDInvalid) {
		response.BadRequest(c, "errors.tenant.invalid", gin.H{"code": "TENANT_ID_INVALID", "message_key": "errors.tenant.invalid"})
		return true
	}
	if errors.Is(err, tenantfilter.ErrSiteIDMissing) {
		response.BadRequest(c, "errors.site.missing", gin.H{"code": "SITE_ID_MISSING", "message_key": "errors.site.missing"})
		return true
	}
	if errors.Is(err, tenantfilter.ErrSiteIDInvalid) {
		response.BadRequest(c, "errors.site.invalid", gin.H{"code": "SITE_ID_INVALID", "message_key": "errors.site.invalid"})
		return true
	}
	return false
}
