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

type TraceRouteRepository interface {
	UpdateTraceRoute(ctx context.Context, TraceRoute models.TraceRoute) error
	GetTraceRoutes(ctx context.Context, projectId int) (traceRouteUseCase []*usecase_models.TraceRoutes, err error)
	GetTraceRoute(ctx context.Context, id int) (traceRouteUseCase *usecase_models.TraceRoutes, err error)
	GetActiveTraceRoutes(ctx context.Context) (traceRouteUseCase []*usecase_models.TraceRoutes, err error)
	SaveTraceRoute(ctx context.Context, TraceRoute models.TraceRoute) (int, error)
	DeleteTraceRoute(ctx context.Context, traceRouteId int) error
}

type traceRouteRepository struct {
	db *sql.DB
}

func NewTraceRouteRepository(db *sql.DB) TraceRouteRepository {
	return &traceRouteRepository{db: db}
}

func (r *traceRouteRepository) SaveTraceRoute(ctx context.Context, traceRoute models.TraceRoute) (int, error) {
	err := traceRoute.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return traceRoute.ID, nil
}

func (r traceRouteRepository) UpdateTraceRoute(ctx context.Context, traceRoute models.TraceRoute) error {
	_, err := traceRoute.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *traceRouteRepository) GetTraceRoutes(ctx context.Context, projectId int) (traceRouteUseCase []*usecase_models.TraceRoutes, err error) {
	traceRoutes, err := models.TraceRoutes(models.TraceRouteWhere.ProjectID.EQ(projectId)).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.TraceRoutes{}, err
	}

	for _, value := range traceRoutes {
		var traceRoute usecase_models.TraceRoutes
		err := json.Unmarshal(value.Data.JSON, &traceRoute)
		if err != nil {
			log.Println(err.Error())
		}
		traceRoute.Scheduling.PipelineId = value.ID
		traceRouteUseCase = append(traceRouteUseCase, &traceRoute)
	}
	return traceRouteUseCase, nil
}

func (r *traceRouteRepository) GetTraceRoute(ctx context.Context, id int) (traceRouteUseCase *usecase_models.TraceRoutes, err error) {
	var traceroute models.Ping
	err = models.TraceRoutes(models.TraceRouteWhere.ID.EQ(id)).Bind(ctx, r.db, &traceroute)
	if err != nil {
		return &usecase_models.TraceRoutes{}, err
	}

	err = json.Unmarshal(traceroute.Data.JSON, &traceRouteUseCase)
	if err != nil {
		log.Println(err.Error())
	}
	traceRouteUseCase.Scheduling.PipelineId = traceroute.ID

	return traceRouteUseCase, nil
}

func (r *traceRouteRepository) GetActiveTraceRoutes(ctx context.Context) (traceRouteUseCase []*usecase_models.TraceRoutes, err error) {
	traceRoutes, err := models.TraceRoutes(qm.Where("data->'scheduling'->>'is_active' = ? and data->'scheduling'->>'end_at' > ?", "true", time.Now().Format("2006-01-02 15:04:05"))).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.TraceRoutes{}, err
	}

	for _, value := range traceRoutes {
		var traceRoute usecase_models.TraceRoutes
		err := json.Unmarshal(value.Data.JSON, &traceRoute)
		if err != nil {
			log.Println(err.Error())
		}
		traceRoute.Scheduling.PipelineId = value.ID
		traceRouteUseCase = append(traceRouteUseCase, &traceRoute)
	}
	return traceRouteUseCase, nil
}

func (r *traceRouteRepository) DeleteTraceRoute(ctx context.Context, traceRouteId int) error {
	traceRoute := models.TraceRoute{ID: traceRouteId}
	_, err := traceRoute.Delete(ctx, r.db, false)
	if err != nil {
		return err
	}
	return nil
}
