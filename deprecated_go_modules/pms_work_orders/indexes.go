package pms_work_orders

import (
	"github.com/aystek/aystek-core/internal/platform/indexes"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	indexes.Register("pms_work_orders", []indexes.IndexSpec{
		{
			Collection: "pms_work_orders",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "property_id", Value: 1}, {Key: "deleted_at", Value: 1}},
			Name:       "pms_work_orders__tenant_site_property_deleted",
		},
		{
			Collection: "pms_work_orders",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "status", Value: 1}, {Key: "priority", Value: 1}, {Key: "created_at", Value: -1}, {Key: "deleted_at", Value: 1}},
			Name:       "pms_work_orders__tenant_site_status_priority_created_deleted",
		},
		{
			Collection: "pms_work_orders",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "unit_id", Value: 1}, {Key: "created_at", Value: -1}, {Key: "deleted_at", Value: 1}},
			Name:       "pms_work_orders__tenant_site_unit_created_deleted",
		},
		{
			Collection: "pms_work_order_events",
			Keys:       bson.D{{Key: "tenant_id", Value: 1}, {Key: "site_id", Value: 1}, {Key: "work_order_id", Value: 1}, {Key: "created_at", Value: -1}},
			Name:       "pms_work_order_events__tenant_site_order_created",
		},
	})
}
