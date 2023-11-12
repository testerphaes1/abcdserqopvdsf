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

type PingRepository interface {
	UpdatePing(ctx context.Context, Ping models.Ping) error
	GetPings(ctx context.Context, projectId int) (pingUseCase []*usecase_models.Pings, err error)
	GetPing(ctx context.Context, id int) (pingUseCase *usecase_models.Pings, err error)
	GetActivePings(ctx context.Context) (pingUseCase []*usecase_models.Pings, err error)
	SavePing(ctx context.Context, Ping models.Ping) (int, error)
	DeletePing(ctx context.Context, pingId int) error
}

type pingRepository struct {
	db *sql.DB
}

func NewPingRepository(db *sql.DB) PingRepository {
	return &pingRepository{db: db}
}

func (r *pingRepository) SavePing(ctx context.Context, ping models.Ping) (int, error) {
	err := ping.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return ping.ID, nil
}

func (r *pingRepository) UpdatePing(ctx context.Context, ping models.Ping) error {
	_, err := ping.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *pingRepository) GetPings(ctx context.Context, projectId int) (pingUseCase []*usecase_models.Pings, err error) {
	pings, err := models.Pings(models.PingWhere.ProjectID.EQ(projectId)).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.Pings{}, err
	}

	for _, value := range pings {
		var ping usecase_models.Pings
		err := json.Unmarshal(value.Data.JSON, &ping)
		if err != nil {
			log.Println(err.Error())
		}
		ping.Scheduling.PipelineId = value.ID
		pingUseCase = append(pingUseCase, &ping)
	}
	return pingUseCase, nil
}

func (r *pingRepository) GetPing(ctx context.Context, id int) (pingUseCase *usecase_models.Pings, err error) {
	var ping models.Ping
	err = models.Pings(models.PingWhere.ID.EQ(id)).Bind(ctx, r.db, &ping)
	if err != nil {
		return &usecase_models.Pings{}, err
	}

	err = json.Unmarshal(ping.Data.JSON, &pingUseCase)
	if err != nil {
		log.Println(err.Error())
	}
	pingUseCase.Scheduling.PipelineId = ping.ID

	return pingUseCase, nil
}

func (r *pingRepository) GetActivePings(ctx context.Context) (pingUseCase []*usecase_models.Pings, err error) {
	pings, err := models.Pings(qm.Where("data->'scheduling'->>'is_active' = ? and data->'scheduling'->>'end_at' > ?", "true", time.Now().Format("2006-01-02 15:04:05"))).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.Pings{}, err
	}

	for _, value := range pings {
		var ping usecase_models.Pings
		err := json.Unmarshal(value.Data.JSON, &ping)
		if err != nil {
			log.Println(err.Error())
		}
		ping.Scheduling.PipelineId = value.ID
		pingUseCase = append(pingUseCase, &ping)
	}
	return pingUseCase, nil
}

func (r *pingRepository) DeletePing(ctx context.Context, pingId int) error {
	ping := models.Ping{ID: pingId}
	_, err := ping.Delete(ctx, r.db, false)
	if err != nil {
		return err
	}
	return nil
}
