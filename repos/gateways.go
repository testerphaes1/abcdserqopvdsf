package repos

import (
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	models "test-manager/usecase_models/boiler"
)

type GatewaysRepository interface {
	UpdateGateways(ctx context.Context, Gateway models.Gateway) error
	GetGateways(ctx context.Context) ([]*models.Gateway, error)
	GetGateway(ctx context.Context, id int) (models.Gateway, error)
	SaveGateways(ctx context.Context, Gateway models.Gateway) (int, error)
}

type gatewaysRepository struct {
	db *sql.DB
}

func NewGatewaysRepository(db *sql.DB) GatewaysRepository {
	return &gatewaysRepository{db: db}
}

func (r *gatewaysRepository) SaveGateways(ctx context.Context, Gateway models.Gateway) (int, error) {
	err := Gateway.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return Gateway.ID, nil
}

func (r *gatewaysRepository) UpdateGateways(ctx context.Context, Gateway models.Gateway) error {
	_, err := Gateway.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *gatewaysRepository) GetGateways(ctx context.Context) ([]*models.Gateway, error) {
	Gateway, err := models.Gateways().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return Gateway, nil
}

func (r *gatewaysRepository) GetGateway(ctx context.Context, id int) (models.Gateway, error) {
	Gateway, err := models.Gateways(models.GatewayWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return models.Gateway{}, err
	}
	return *Gateway, nil
}
