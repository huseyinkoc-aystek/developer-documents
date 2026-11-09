package pms_units_bridge

import (
	newbridge "github.com/aystek/aystek-core/internal/modules/pms/units/bridge"
	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/handler"
	"github.com/gin-gonic/gin"
)

// Deprecated: use internal/modules/pms/units/bridge.RegisterRoutes.
func RegisterRoutes(g *gin.RouterGroup, admin *handler.AdminHandler) {
	newbridge.RegisterRoutes(g, admin)
}
