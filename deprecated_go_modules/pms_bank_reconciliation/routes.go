package pms_bank_reconciliation

import (
	newreconciliation "github.com/aystek/aystek-core/internal/modules/pms/reconciliation"
	"github.com/aystek/aystek-core/internal/modules/pms_bank_reconciliation/handler"
	"github.com/gin-gonic/gin"
)

// Deprecated: use internal/modules/pms/reconciliation.RegisterRoutes.
func RegisterRoutes(g *gin.RouterGroup, admin *handler.AdminHandler) {
	newreconciliation.RegisterRoutes(g, admin)
}
