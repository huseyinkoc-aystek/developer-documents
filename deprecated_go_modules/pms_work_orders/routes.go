package pms_work_orders

import (
	newworkorders "github.com/aystek/aystek-core/internal/modules/pms/work_orders"
	"github.com/aystek/aystek-core/internal/modules/pms_work_orders/handler"
	"github.com/gin-gonic/gin"
)

// Deprecated: use internal/modules/pms/work_orders.RegisterRoutes.
func RegisterRoutes(g *gin.RouterGroup, admin *handler.AdminHandler) {
	newworkorders.RegisterRoutes(g, admin)
}
