package repository

import (
	"context"

	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListFilter struct {
	Page       int64
	Limit      int64
	PropertyID *primitive.ObjectID
	Q          string
}

type Repository interface {
	ListUnits(ctx context.Context, filter ListFilter) ([]model.UnitListItem, int64, error)
	GetUnitDetail(ctx context.Context, id primitive.ObjectID) (*model.UnitDetail, error)
}
