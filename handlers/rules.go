package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/volatiletech/null/v8"
	"test-manager/repos"
	"test-manager/tasks/push"
	"test-manager/usecase_models"
	models "test-manager/usecase_models/boiler"
)

type RulesHandler interface {
	RegisterRules(ctx context.Context, rules usecase_models.RulesRequest) error
}

type rulesHandler struct {
	projectRepo     repos.ProjectsRepository
	endpointRepo    repos.EndpointRepository
	netCatRepo      repos.NetCatRepository
	pageSpeedRepo   repos.PageSpeedRepository
	pingRepo        repos.PingRepository
	traceRouteRepo  repos.TraceRouteRepository
	dataCentersRepo repos.DataCentersRepository
	taskPusher      push.TaskPusher
	agentHandler    AgentHandler
}

func NewRulesHandler(
	projectRepo repos.ProjectsRepository,
	endpointRepo repos.EndpointRepository,
	netCatRepo repos.NetCatRepository,
	pageSpeedRepo repos.PageSpeedRepository,
	pingRepo repos.PingRepository,
	traceRouteRepo repos.TraceRouteRepository,
	dataCentersRepo repos.DataCentersRepository,
	taskPusher push.TaskPusher,
	agentHandler AgentHandler,
) RulesHandler {
	return &rulesHandler{
		projectRepo:     projectRepo,
		endpointRepo:    endpointRepo,
		netCatRepo:      netCatRepo,
		pageSpeedRepo:   pageSpeedRepo,
		pingRepo:        pingRepo,
		traceRouteRepo:  traceRouteRepo,
		dataCentersRepo: dataCentersRepo,
		taskPusher:      taskPusher,
		agentHandler:    agentHandler,
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (r *rulesHandler) RegisterRules(ctx context.Context, rules usecase_models.RulesRequest) error {
	projects, err := r.projectRepo.GetProjects(ctx, IdentityStruct.Id)
	if err != nil {
		return err
	}
	var projectIds []int
	for _, value := range projects {
		projectIds = append(projectIds, value.ID)
	}
	if len(rules.Endpoints.Endpoints) != 0 && !contains(projectIds, rules.Endpoints.Scheduling.ProjectId) {
		return errors.New(fmt.Sprintf("project id : %d is not your project", rules.Endpoints.Scheduling.ProjectId))
	}
	if len(rules.NetCats.NetCats) != 0 && !contains(projectIds, rules.NetCats.Scheduling.ProjectId) {
		return errors.New(fmt.Sprintf("project id : %d is not your project", rules.NetCats.Scheduling.ProjectId))
	}
	if len(rules.Pings.Pings) != 0 && !contains(projectIds, rules.Pings.Scheduling.ProjectId) {
		return errors.New(fmt.Sprintf("project id : %d is not your project", rules.Pings.Scheduling.ProjectId))
	}
	if len(rules.TraceRoutes.TraceRouts) != 0 && !contains(projectIds, rules.TraceRoutes.Scheduling.ProjectId) {
		return errors.New(fmt.Sprintf("project id : %d is not your project", rules.TraceRoutes.Scheduling.ProjectId))
	}
	if len(rules.PageSpeed.PageSpeed) != 0 && !contains(projectIds, rules.PageSpeed.Scheduling.ProjectId) {
		return errors.New(fmt.Sprintf("project id : %d is not your project", rules.PageSpeed.Scheduling.ProjectId))
	}

	if rules.Endpoints.Scheduling.IsHeartBeat && len(rules.Endpoints.Endpoints) > 1 {
		return errors.New("more that one endpoint can not be registered if heartbeat is active")
	}
	if rules.NetCats.Scheduling.IsHeartBeat && len(rules.NetCats.NetCats) > 1 {
		return errors.New("more that one net cats can not be registered if heartbeat is active")
	}
	if rules.Pings.Scheduling.IsHeartBeat && len(rules.Pings.Pings) > 1 {
		return errors.New("more that one ping can not be registered if heartbeat is active")
	}
	if rules.TraceRoutes.Scheduling.IsHeartBeat && len(rules.TraceRoutes.TraceRouts) > 1 {
		return errors.New("more that one trace route can not be registered if heartbeat is active")
	}
	if rules.PageSpeed.Scheduling.IsHeartBeat && len(rules.PageSpeed.PageSpeed) > 1 {
		return errors.New("more that one page speed can not be registered if heartbeat is active")
	}

	if len(rules.Endpoints.Endpoints) != 0 {
		var project models.Project
		for _, value := range projects {
			if value.ID == rules.Endpoints.Scheduling.ProjectId {
				project = *value
			}
		}
		rules.Endpoints.Scheduling.EndAt = project.ExpireAt.Time.String()
		j, _ := json.Marshal(rules.Endpoints)
		endpointId, err := r.endpointRepo.SaveEndpoint(ctx, models.Endpoint{
			Data:      null.NewJSON(j, true),
			ProjectID: rules.Endpoints.Scheduling.ProjectId,
		})
		if err != nil {
			return err
		}
		rules.Endpoints.Scheduling.PipelineId = endpointId
	}
	if len(rules.NetCats.NetCats) != 0 {
		var project models.Project
		for _, value := range projects {
			if value.ID == rules.NetCats.Scheduling.ProjectId {
				project = *value
			}
		}
		rules.NetCats.Scheduling.EndAt = project.ExpireAt.Time.String()
		j, _ := json.Marshal(rules.NetCats)
		netcatId, err := r.netCatRepo.SaveNetCat(ctx, models.NetCat{
			Data:      null.NewJSON(j, true),
			ProjectID: rules.NetCats.Scheduling.ProjectId,
		})
		if err != nil {
			return err
		}
		rules.NetCats.Scheduling.PipelineId = netcatId
	}
	if len(rules.PageSpeed.PageSpeed) != 0 {
		var project models.Project
		for _, value := range projects {
			if value.ID == rules.PageSpeed.Scheduling.ProjectId {
				project = *value
			}
		}
		rules.PageSpeed.Scheduling.EndAt = project.ExpireAt.Time.String()
		j, _ := json.Marshal(rules.PageSpeed)
		pagespeedId, err := r.pageSpeedRepo.SavePageSpeed(ctx, models.PageSpeed{
			Data:      null.NewJSON(j, true),
			ProjectID: rules.PageSpeed.Scheduling.ProjectId,
		})
		if err != nil {
			return err
		}
		rules.PageSpeed.Scheduling.PipelineId = pagespeedId
	}
	if len(rules.Pings.Pings) != 0 {
		var project models.Project
		for _, value := range projects {
			if value.ID == rules.Pings.Scheduling.ProjectId {
				project = *value
			}
		}
		rules.Pings.Scheduling.EndAt = project.ExpireAt.Time.String()
		j, _ := json.Marshal(rules.Pings)
		pingId, err := r.pingRepo.SavePing(ctx, models.Ping{
			Data:      null.NewJSON(j, true),
			ProjectID: rules.Pings.Scheduling.ProjectId,
		})
		if err != nil {
			return err
		}
		rules.Pings.Scheduling.PipelineId = pingId
	}
	if len(rules.TraceRoutes.TraceRouts) != 0 {
		var project models.Project
		for _, value := range projects {
			if value.ID == rules.TraceRoutes.Scheduling.ProjectId {
				project = *value
			}
		}
		rules.TraceRoutes.Scheduling.EndAt = project.ExpireAt.Time.String()
		j, _ := json.Marshal(rules.TraceRoutes)
		tracerouteId, err := r.traceRouteRepo.SaveTraceRoute(ctx, models.TraceRoute{
			Data:      null.NewJSON(j, true),
			ProjectID: rules.TraceRoutes.Scheduling.ProjectId,
		})
		if err != nil {
			return err
		}
		rules.TraceRoutes.Scheduling.PipelineId = tracerouteId
	}

	//_, err = r.taskPusher.PushRules(ctx, rules)
	//if err != nil {
	//	return err
	//}

	return nil
}
