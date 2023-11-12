package repos

import (
	"context"
	"database/sql"
	"test-manager/usecase_models"
)

type AggregateRepository interface {
	AggregateAllRuleSubWorks(ctx context.Context, projectId int) (usecase_models.AggregateAllRuleSubWorks, error)
}

type aggregateRepository struct {
	db             *sql.DB
	endpointRepo   EndpointRepository
	netCatRepo     NetCatRepository
	pageSpeedRepo  PageSpeedRepository
	pingRepo       PingRepository
	traceRouteRepo TraceRouteRepository
}

func NewAggregateRepository(db *sql.DB,
	endpointRepo EndpointRepository,
	netCatRepo NetCatRepository,
	pageSpeedRepo PageSpeedRepository,
	pingRepo PingRepository,
	traceRouteRepo TraceRouteRepository) AggregateRepository {
	return &aggregateRepository{
		db:             db,
		endpointRepo:   endpointRepo,
		netCatRepo:     netCatRepo,
		pageSpeedRepo:  pageSpeedRepo,
		pingRepo:       pingRepo,
		traceRouteRepo: traceRouteRepo,
	}
}

func (a *aggregateRepository) AggregateAllRuleSubWorks(ctx context.Context, projectId int) (usecase_models.AggregateAllRuleSubWorks, error) {
	endpoints, err := a.endpointRepo.GetEndpoints(ctx, projectId)
	if err != nil {
		return usecase_models.AggregateAllRuleSubWorks{}, err
	}
	netcats, err := a.netCatRepo.GetNetCats(ctx, projectId)
	if err != nil {
		return usecase_models.AggregateAllRuleSubWorks{}, err
	}
	pageSpeeds, err := a.pageSpeedRepo.GetPageSpeeds(ctx, projectId)
	if err != nil {
		return usecase_models.AggregateAllRuleSubWorks{}, err
	}
	pings, err := a.pingRepo.GetPings(ctx, projectId)
	if err != nil {
		return usecase_models.AggregateAllRuleSubWorks{}, err
	}
	traceRoutes, err := a.traceRouteRepo.GetTraceRoutes(ctx, projectId)
	if err != nil {
		return usecase_models.AggregateAllRuleSubWorks{}, err
	}

	return usecase_models.AggregateAllRuleSubWorks{
		Endpoints:   endpoints,
		TraceRoutes: traceRoutes,
		NetCats:     netcats,
		Pings:       pings,
		PageSpeed:   pageSpeeds,
	}, nil
}
