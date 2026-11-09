package pms_bank_reconciliation

import (
	"github.com/aystek/aystek-core/internal/platform/indexes"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	indexes.Register("pms_bank_reconciliation", []indexes.IndexSpec{
		{
			Collection: "pms_bank_transactions",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "external_txn_id", Value: 1}, {Key: "deleted_at", Value: 1}},
			Name:       "pms_bank_transactions__tenant_site_external_deleted",
			Unique:     true,
			Partial:    bson.M{"deleted_at": nil},
		},
		{
			Collection: "pms_bank_transactions",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "status", Value: 1}, {Key: "txn_date", Value: -1}, {Key: "deleted_at", Value: 1}},
			Name:       "pms_bank_transactions__tenant_site_status_txndate_deleted",
		},
	})
}
