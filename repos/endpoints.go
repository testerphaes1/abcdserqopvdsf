package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"test-manager/usecase_models"
	models "test-manager/usecase_models/boiler"
	"time"
)

const (
	EndpointSessionCachePrefix = "endpoint:session:endpoint_id:"
)

const (
	DisabledEndpointReasonFunctionality = "functionality problem"
)

type EndpointRepository interface {
	UpdateEndpoint(ctx context.Context, endpoint models.Endpoint) error
	GetEndpoints(ctx context.Context, projectId int) (endpointsUseCase []*usecase_models.Endpoints, err error)
	GetEndpoint(ctx context.Context, id int) (endpointsUseCase *usecase_models.Endpoints, err error)
	GetActiveEndpoints(ctx context.Context) (endpointsUseCase []*usecase_models.Endpoints, err error)
	SaveEndpoint(ctx context.Context, endpoint models.Endpoint) (int, error)
	DisableEndpointWithFunctionalityReason(ctx context.Context, endpointId int) error
	DeleteEndpoint(ctx context.Context, endpointId int) error
}

type endpointRepository struct {
	db *sql.DB
}

func NewEndpointRepository(db *sql.DB) EndpointRepository {
	return &endpointRepository{db: db}
}

func (r *endpointRepository) SaveEndpoint(ctx context.Context, endpoint models.Endpoint) (int, error) {
	err := endpoint.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return endpoint.ID, nil
}

func (r *endpointRepository) UpdateEndpoint(ctx context.Context, endpoint models.Endpoint) error {
	_, err := endpoint.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *endpointRepository) GetEndpoints(ctx context.Context, projectId int) (endpointsUseCase []*usecase_models.Endpoints, err error) {
	endpoints, err := models.Endpoints(models.EndpointWhere.ProjectID.EQ(projectId)).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.Endpoints{}, err
	}

	for _, value := range endpoints {
		var endpoint usecase_models.Endpoints
		err := json.Unmarshal(value.Data.JSON, &endpoint)
		if err != nil {
			log.Println(err.Error())
		}
		endpoint.Scheduling.PipelineId = value.ID
		endpointsUseCase = append(endpointsUseCase, &endpoint)
	}
	return endpointsUseCase, nil
}

func (r *endpointRepository) GetEndpoint(ctx context.Context, id int) (endpointsUseCase *usecase_models.Endpoints, err error) {
	var endpoint models.Endpoint
	err = models.Endpoints(models.EndpointWhere.ID.EQ(id)).Bind(ctx, r.db, &endpoint)
	if err != nil {
		return &usecase_models.Endpoints{}, err
	}

	err = json.Unmarshal(endpoint.Data.JSON, &endpointsUseCase)
	if err != nil {
		log.Println(err.Error())
	}
	endpointsUseCase.Scheduling.PipelineId = endpoint.ID

	return endpointsUseCase, nil
}

func (r *endpointRepository) GetActiveEndpoints(ctx context.Context) (endpointsUseCase []*usecase_models.Endpoints, err error) {
	endpoints, err := models.Endpoints(
		qm.Where("data->'scheduling'->>'is_active' = ? and data->'scheduling'->>'end_at' > ?", "true", time.Now().Format("2006-01-02 15:04:05")),
		qm.Where("p.is_active = true"),
		qm.InnerJoin("projects p on p.id = project_id"),
	).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.Endpoints{}, err
	}

	for _, value := range endpoints {
		var endpoint usecase_models.Endpoints
		err := json.Unmarshal(value.Data.JSON, &endpoint)
		if err != nil {
			log.Println(err.Error())
		}
		endpoint.Scheduling.PipelineId = value.ID
		endpointsUseCase = append(endpointsUseCase, &endpoint)
	}
	return endpointsUseCase, nil
}

func (r *endpointRepository) DisableEndpointWithFunctionalityReason(ctx context.Context, endpointId int) error {
	query := queries.Raw("update endpoints set disabled = $1, data = jsonb_set(data, '{scheduling, is_active}','false') where id = $2;", DisabledEndpointReasonFunctionality, endpointId)
	_, err := query.ExecContext(ctx, r.db)
	if err != nil {
		return err
	}
	fmt.Println("clean up successfully finished")
	return nil
}

func (r *endpointRepository) DeleteEndpoint(ctx context.Context, endpointId int) error {
	endpoint := models.Endpoint{ID: endpointId}
	_, err := endpoint.Delete(ctx, r.db, false)
	if err != nil {
		return err
	}
	return nil
}
