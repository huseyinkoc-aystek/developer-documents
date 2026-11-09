package service

import (
	"context"

	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/model"
	"github.com/aystek/aystek-core/internal/modules/pms_units_bridge/repository"
	"github.com/aystek/aystek-core/internal/platform/tenantfilter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListUnits(ctx context.Context, filter repository.ListFilter) ([]model.UnitListItem, int64, error) {
	if err := validateContext(ctx); err != nil {
		return nil, 0, err
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.repo.ListUnits(ctx, filter)
}

func (s *service) GetUnitDetail(ctx context.Context, id primitive.ObjectID) (*model.UnitDetail, error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	return s.repo.GetUnitDetail(ctx, id)
}

func validateContext(ctx context.Context) error {
	if _, err := tenantfilter.RequireTenantID(ctx); err != nil {
		return err
	}
	if _, err := tenantfilter.RequireSiteID(ctx); err != nil {
		return err
	}
	return nil
}

var _ Service = (*service)(nil)
