package pms_billing_bridge

import (
	newbridge "github.com/aystek/aystek-core/internal/modules/pms/billing/bridge"
	"github.com/aystek/aystek-core/internal/modules/pms_billing_bridge/handler"
	"github.com/gin-gonic/gin"
)

// Deprecated: use internal/modules/pms/billing/bridge.RegisterRoutes.
func RegisterRoutes(g *gin.RouterGroup, admin *handler.AdminHandler) {
	newbridge.RegisterRoutes(g, admin)
}
