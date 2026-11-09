package pms_billing_bridge

import (
	"github.com/aystek/aystek-core/internal/platform/indexes"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	indexes.Register("pms_billing_bridge", []indexes.IndexSpec{
		{
			Collection: "pms_charges",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "property_id", Value: 1}, {Key: "unit_id", Value: 1}, {Key: "period", Value: 1}, {Key: "deleted_at", Value: 1}},
			Name:       "pms_charges__tenant_site_property_unit_period_deleted",
			Unique:     true,
			Partial:    bson.M{"deleted_at": nil},
		},
		{
			Collection: "pms_charges",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "property_id", Value: 1}, {Key: "period", Value: 1}, {Key: "status", Value: 1}, {Key: "updated_at", Value: -1}, {Key: "deleted_at", Value: 1}},
			Name:       "pms_charges__tenant_site_property_period_status_updated_deleted",
		},
		{
			Collection: "pms_payment_allocations",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "payment_ref", Value: 1}, {Key: "created_at", Value: -1}},
			Name:       "pms_payment_allocations__tenant_site_paymentref_created",
		},
		{
			Collection: "pms_billing_idempotency",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "operation", Value: 1}, {Key: "idempotency_key", Value: 1}},
			Name:       "pms_billing_idempotency__tenant_site_operation_key",
			Unique:     true,
		},
	})
}
