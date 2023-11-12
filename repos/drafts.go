package repos

import (
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	models "test-manager/usecase_models/boiler"
)

type DraftsRepository interface {
	DeleteDraft(ctx context.Context, Draft models.Draft) error
	GetDrafts(ctx context.Context, projectId int) ([]*models.Draft, error)
	GetDraft(ctx context.Context, id int) (models.Draft, error)
	SaveDrafts(ctx context.Context, Draft models.Draft) (int, error)
}

type draftsRepository struct {
	db *sql.DB
}

func NewDraftsRepository(db *sql.DB) DraftsRepository {
	return &draftsRepository{db: db}
}

func (r *draftsRepository) SaveDrafts(ctx context.Context, Draft models.Draft) (int, error) {
	err := Draft.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return Draft.ID, nil
}

func (r *draftsRepository) DeleteDraft(ctx context.Context, Draft models.Draft) error {
	_, err := Draft.Delete(ctx, r.db, false)
	if err != nil {
		return err
	}
	return nil
}

func (r *draftsRepository) GetDrafts(ctx context.Context, projectId int) ([]*models.Draft, error) {
	Draft, err := models.Drafts(models.DraftWhere.ProjectID.EQ(projectId)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return Draft, nil
}

func (r *draftsRepository) GetDraft(ctx context.Context, id int) (models.Draft, error) {
	Draft, err := models.Drafts(models.DraftWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return models.Draft{}, err
	}
	return *Draft, nil
}
