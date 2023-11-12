package handlers

import (
	"context"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"math/rand"
	"strings"
	"sync"
	"test-manager/repos"
	"test-manager/services/alert_system"
	"test-manager/tasks/push"
	"test-manager/tasks/task_models"
	"test-manager/usecase_models"
	models "test-manager/usecase_models/boiler"
	"time"
)

type TraceRouteHandler interface {
	ExecuteTraceRouteRule(ctx context.Context, TraceRouteRules usecase_models.TraceRoutes) error
}

type traceRouteHandler struct {
	alertSystem         alert_system.AlertHandler
	traceRouteRepo      repos.TraceRouteRepository
	dataCentersRepo     repos.DataCentersRepository
	projectRepo         repos.ProjectsRepository
	traceRouteStatsRepo repos.TraceRouteStatsRepository
	taskPusher          push.TaskPusher
	agentHandler        AgentHandler
}

func NewTraceRouteHandler(alertSystem alert_system.AlertHandler, traceRouteRepo repos.TraceRouteRepository, dataCentersRepo repos.DataCentersRepository, projectRepo repos.ProjectsRepository, traceRouteStatsRepo repos.TraceRouteStatsRepository, taskPusher push.TaskPusher, agentHandler AgentHandler) TraceRouteHandler {
	return &traceRouteHandler{
		alertSystem:         alertSystem,
		traceRouteRepo:      traceRouteRepo,
		dataCentersRepo:     dataCentersRepo,
		projectRepo:         projectRepo,
		traceRouteStatsRepo: traceRouteStatsRepo,
		taskPusher:          taskPusher,
		agentHandler:        agentHandler,
	}
}

func (e *traceRouteHandler) ExecuteTraceRouteRule(ctx context.Context, traceRouteRules usecase_models.TraceRoutes) error {
	if len(traceRouteRules.Scheduling.DataCentersIds) == 1 && traceRouteRules.Scheduling.DataCentersIds[0] == 0 {
		datacenters, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err == nil {
			// random number between one and the number of datacenters
			traceRouteRules.Scheduling.DataCentersIds = []int{}
			traceRouteRules.Scheduling.DataCentersIds = append(traceRouteRules.Scheduling.DataCentersIds, datacenters[rand.Intn(len(datacenters))].ID)
		}
	} else if len(traceRouteRules.Scheduling.DataCentersIds) == 0 {
		datacenters, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err == nil {
			// random number between one and the number of datacenters
			traceRouteRules.Scheduling.DataCentersIds = []int{}
			for _, value := range datacenters {
				traceRouteRules.Scheduling.DataCentersIds = append(traceRouteRules.Scheduling.DataCentersIds, value.ID)
			}
		}
	}

	isHeart := traceRouteRules.Scheduling.IsHeartBeat
	currentTime := time.Now()
	for {
		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(traceRouteRules.Scheduling.DataCentersIds))
		for _, dataC := range traceRouteRules.Scheduling.DataCentersIds {
			go func(dataCenter int) {
				var urlsCalled []string
				for _, rule := range traceRouteRules.TraceRouts {
					dataCenter, err := e.dataCentersRepo.GetDataCenter(ctx, dataCenter)
					if err != nil {
						log.Info("error on getting data center in executing trace rote rule: ", err)
						waitGroup.Done()
						return
					}

					response, err := e.agentHandler.SendTraceRoute(ctx, dataCenter.Baseurl, usecase_models.AgentTraceRouteRequest{
						Address: rule.Address,
						Retry:   rule.Retry,
						Hop:     rule.Hop,
					})
					if err != nil {
						log.Info("error on sending trace route in executing rule: ", err)
						waitGroup.Done()
						return
					}

					if response.Status == 0 {
						urlsCalled = append(urlsCalled, rule.Address)
						err = e.traceRouteStatsRepo.Write(ctx, time.Now(), repos.WriteTraceRouteStatsOptions{
							ProjectId:    traceRouteRules.Scheduling.ProjectId,
							TraceRouteId: traceRouteRules.Scheduling.PipelineId,
							IsHeartBeat:  isHeart,
							Url:          strings.Join(urlsCalled, ","),
							DatacenterId: dataCenter.ID,
							Success:      0,
						})
						if err != nil {
							log.Info("error on writing curl report in executing rule: ", err)
						}

						waitGroup.Done()
						return
					}
					urlsCalled = append(urlsCalled, rule.Address)
				}
				err := e.traceRouteStatsRepo.Write(ctx, time.Now(), repos.WriteTraceRouteStatsOptions{
					ProjectId:    traceRouteRules.Scheduling.ProjectId,
					TraceRouteId: traceRouteRules.Scheduling.PipelineId,
					IsHeartBeat:  isHeart,
					Url:          strings.Join(urlsCalled, ","),
					DatacenterId: dataCenter,
					Success:      1,
				})
				if err != nil {
					log.Info("error on writing curl report in executing rule: ", err)
				}
				waitGroup.Done()
			}(dataC)
		}

		waitGroup.Wait()

		sessionIds, err := e.traceRouteStatsRepo.GetLastNSessionsByTraceRouteId(ctx, 2, traceRouteRules.Scheduling.PipelineId)
		if err != nil {
			log.Info("problem in fetching session ids: ", err.Error())
			break
		}

		if len(sessionIds) != 2 {
			break
		}
		newSession, err := e.traceRouteStatsRepo.Read(ctx,
			[]string{"time", "success", "datacenter_id", "url"},
			repos.Filters{repos.Filter{Field: "session_id", Op: repos.FilterOpEq, Value: sessionIds[0]}},
			[]string{models.PageSpeedsStatRels.Datacenter},
		)
		oldSession, err := e.traceRouteStatsRepo.Read(ctx,
			[]string{"time", "success", "datacenter_id", "url"},
			repos.Filters{repos.Filter{Field: "session_id", Op: repos.FilterOpEq, Value: sessionIds[1]}},
			[]string{models.PageSpeedsStatRels.Datacenter},
		)

		project, err := e.projectRepo.GetProjectWithLoads(ctx, traceRouteRules.Scheduling.ProjectId)
		if err != nil {
			log.Info("error on getting project in executing rule: ", err)
			break
		}
		var notifications usecase_models.Notifications
		err = json.Unmarshal(project.Notifications.JSON, &notifications)
		if err != nil {
			log.Info("can not unmarshal notification")
			break
		}

		resolved, failed, newSuccess, oldSuccess := calculateTraceRouteStatState(newSession, oldSession)
		if oldSuccess && !newSuccess {
			// send fail
			address := ""
			rootCause := "working on it"
			for _, value := range oldSession {
				if value.Success == 0 {
					t := strings.Split(value.URL.String, ",")
					address = t[len(t)-1]
					break
				}
			}
			dcs, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
			if err != nil {
				log.Info("problem in getting datacenters in sending alerts: ", err.Error())
			}
			var dcTitles []string
			for _, value := range failed {
				for _, dc := range dcs {
					if dc.ID == value {
						dcTitles = append(dcTitles, dc.Title)
					}
				}
			}
			_, err = e.taskPusher.PushNotifications(ctx, task_models.NotificationsPayload{
				Type:                "traceroute",
				State:               "down",
				ProjectId:           traceRouteRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        traceRouteRules.Scheduling.PipelineName,
				Time:                time.Now().String(),
				Address:             address,
				Datacenters:         strings.Join(dcTitles, ","),
				RootCause:           rootCause,
				ResolvedDatacenters: "",
				FailedDatacenters:   "",
			})
			if err != nil {
				log.Info(err.Error())
			}
		}
		if !oldSuccess && newSuccess {
			// send success
			address := ""
			rootCause := "working on it"
			for _, value := range oldSession {
				if value.Success == 1 {
					t := strings.Split(value.URL.String, ",")
					address = t[len(t)-1]
					break
				}
			}
			dcs, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
			if err != nil {
				log.Info("problem in getting datacenters in sending alerts: ", err.Error())
			}
			var dcTitles []string
			for _, value := range resolved {
				for _, dc := range dcs {
					if dc.ID == value {
						dcTitles = append(dcTitles, dc.Title)
					}
				}
			}
			_, err = e.taskPusher.PushNotifications(ctx, task_models.NotificationsPayload{
				Type:                "traceroute",
				State:               "up",
				ProjectId:           traceRouteRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        traceRouteRules.Scheduling.PipelineName,
				Time:                time.Now().String(),
				Address:             address,
				Datacenters:         strings.Join(dcTitles, ","),
				RootCause:           rootCause,
				ResolvedDatacenters: "",
				FailedDatacenters:   "",
			})
			if err != nil {
				log.Info(err.Error())
			}
		}
		if !oldSuccess && !newSuccess && (len(resolved) > 0 || len(failed) > 0) {
			// send diff
			address := ""
			rootCause := "working on it"
			for _, value := range oldSession {
				if value.Success == 1 {
					t := strings.Split(value.URL.String, ",")
					address = t[len(t)-1]
					break
				}
			}
			dcs, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
			if err != nil {
				log.Info("problem in getting datacenters in sending alerts: ", err.Error())
			}
			var rTitles []string
			var fTitles []string
			for _, value := range resolved {
				for _, dc := range dcs {
					if dc.ID == value {
						rTitles = append(rTitles, dc.Title)
					}
				}
			}
			for _, value := range failed {
				for _, dc := range dcs {
					if dc.ID == value {
						fTitles = append(fTitles, dc.Title)
					}
				}
			}
			_, err = e.taskPusher.PushNotifications(ctx, task_models.NotificationsPayload{
				Type:                "traceroute",
				State:               "diff",
				ProjectId:           traceRouteRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        traceRouteRules.Scheduling.PipelineName,
				Time:                time.Now().String(),
				Address:             address,
				Datacenters:         "",
				RootCause:           rootCause,
				ResolvedDatacenters: strings.Join(rTitles, ","),
				FailedDatacenters:   strings.Join(fTitles, ","),
			})
			if err != nil {
				log.Info(err.Error())
			}
		}

		if isHeart {
			if time.Now().Sub(currentTime) >= 59*time.Second {
				break
			}
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return nil
}

func calculateTraceRouteStatState(newSession, oldSession models.TraceRoutesStatSlice) ([]int, []int, bool, bool) {
	var mapNDcSuccess map[int]bool
	var mapODcSuccess map[int]bool
	var newSuccess = true
	var oldSuccess = true
	for _, value := range newSession {
		mapNDcSuccess[value.DatacenterID] = value.Success != 0
	}
	for _, value := range oldSession {
		if _, ok := mapNDcSuccess[value.DatacenterID]; !ok {
			mapNDcSuccess[value.DatacenterID] = true
		}
		mapODcSuccess[value.DatacenterID] = value.Success != 0
	}

	var resolved []int
	var failed []int
	for key, newValue := range mapNDcSuccess {
		oldValue, ok := mapODcSuccess[key]
		if !ok {
			oldValue = true
		}
		newSuccess = newSuccess && newValue
		oldSuccess = oldSuccess && oldValue

		if newValue && !oldValue {
			// resolved
			resolved = append(resolved, key)
		} else if !newValue && oldValue {
			// error detected
			failed = append(failed, key)
		} else {
			// same shit
		}
	}
	return resolved, failed, newSuccess, oldSuccess
}
