package service

import (
	"context"

	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/model"
	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListUnits(ctx context.Context, filter repository.ListFilter) ([]model.UnitListItem, int64, error)
	GetUnitDetail(ctx context.Context, id primitive.ObjectID) (*model.UnitDetail, error)
}
