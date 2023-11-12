package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"test-manager/usecase_models"
	models "test-manager/usecase_models/boiler"
	"time"
)

type PageSpeedRepository interface {
	UpdatePageSpeed(ctx context.Context, PageSpeed models.PageSpeed) error
	GetPageSpeeds(ctx context.Context, projectId int) (pageSpeedUseCase []*usecase_models.PageSpeeds, err error)
	GetPageSpeed(ctx context.Context, id int) (pageSpeedUseCase *usecase_models.PageSpeeds, err error)
	GetActivePageSpeeds(ctx context.Context) (pageSpeedUseCase []*usecase_models.PageSpeeds, err error)
	SavePageSpeed(ctx context.Context, PageSpeed models.PageSpeed) (int, error)
	DeletePageSpeed(ctx context.Context, pagespeedId int) error
}

type pageSpeedRepository struct {
	db *sql.DB
}

func NewPageSpeedRepository(db *sql.DB) PageSpeedRepository {
	return &pageSpeedRepository{db: db}
}

func (r *pageSpeedRepository) SavePageSpeed(ctx context.Context, pageSpeed models.PageSpeed) (int, error) {
	err := pageSpeed.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return pageSpeed.ID, nil
}

func (r *pageSpeedRepository) UpdatePageSpeed(ctx context.Context, pageSpeed models.PageSpeed) error {
	_, err := pageSpeed.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *pageSpeedRepository) GetPageSpeeds(ctx context.Context, projectId int) (pageSpeedUseCase []*usecase_models.PageSpeeds, err error) {
	pageSpeeds, err := models.PageSpeeds(models.PageSpeedWhere.ProjectID.EQ(projectId)).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.PageSpeeds{}, err
	}

	for _, value := range pageSpeeds {
		var pageSpeed usecase_models.PageSpeeds
		err := json.Unmarshal(value.Data.JSON, &pageSpeed)
		if err != nil {
			log.Println(err.Error())
		}
		pageSpeed.Scheduling.PipelineId = value.ID
		pageSpeedUseCase = append(pageSpeedUseCase, &pageSpeed)
	}
	return pageSpeedUseCase, nil
}

func (r *pageSpeedRepository) GetPageSpeed(ctx context.Context, id int) (pageSpeedUseCase *usecase_models.PageSpeeds, err error) {
	var pagespeed models.PageSpeed
	err = models.PageSpeeds(models.PageSpeedWhere.ID.EQ(id)).Bind(ctx, r.db, &pagespeed)
	if err != nil {
		return &usecase_models.PageSpeeds{}, err
	}

	err = json.Unmarshal(pagespeed.Data.JSON, &pageSpeedUseCase)
	if err != nil {
		log.Println(err.Error())
	}
	pageSpeedUseCase.Scheduling.PipelineId = pagespeed.ID

	return pageSpeedUseCase, nil
}

func (r *pageSpeedRepository) GetActivePageSpeeds(ctx context.Context) (pageSpeedUseCase []*usecase_models.PageSpeeds, err error) {
	pageSpeeds, err := models.PageSpeeds(qm.Where("data->'scheduling'->>'is_active' = ? and data->'scheduling'->>'end_at' > ?", "true", time.Now().Format("2006-01-02 15:04:05"))).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.PageSpeeds{}, err
	}

	for _, value := range pageSpeeds {
		var pageSpeed usecase_models.PageSpeeds
		err := json.Unmarshal(value.Data.JSON, &pageSpeed)
		if err != nil {
			log.Println(err.Error())
		}
		pageSpeed.Scheduling.PipelineId = value.ID
		pageSpeedUseCase = append(pageSpeedUseCase, &pageSpeed)
	}
	return pageSpeedUseCase, nil
}

func (r *pageSpeedRepository) DeletePageSpeed(ctx context.Context, pagespeedId int) error {
	pageSpeed := models.PageSpeed{ID: pagespeedId}
	_, err := pageSpeed.Delete(ctx, r.db, false)
	if err != nil {
		return err
	}
	return nil
}
