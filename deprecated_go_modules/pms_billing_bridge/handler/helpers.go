package handler

import (
	"errors"
	"strings"

	"github.com/aystek/aystek-core/internal/platform/response"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func actorID(c *gin.Context) *primitive.ObjectID {
	if c == nil {
		return nil
	}
	v := strings.TrimSpace(c.GetString("user_id"))
	if v == "" {
		v = strings.TrimSpace(c.GetString("userID"))
	}
	if v == "" {
		return nil
	}
	oid, err := primitive.ObjectIDFromHex(v)
	if err != nil {
		return nil
	}
	return &oid
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
