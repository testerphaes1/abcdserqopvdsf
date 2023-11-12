package repos

import (
	"context"
	"database/sql"
	models "test-manager/usecase_models/boiler"
)

type AuthInfoRepository interface {
	GetAuthKeys(ctx context.Context) (string, error)
}

type authInfoRepository struct {
	db *sql.DB
}

func NewAuthInfoRepositoryRepository(db *sql.DB) AuthInfoRepository {
	return &authInfoRepository{db: db}
}

func (r *authInfoRepository) GetAuthKeys(ctx context.Context) (string, error) {
	info, err := models.AuthInfos().One(ctx, r.db)
	if err != nil {
		return "", err
	}
	return info.PrivateKey.String, nil
}
