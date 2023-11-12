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

type PageSpeedHandler interface {
	ExecutePageSpeedRule(ctx context.Context, PageSpeedRules usecase_models.PageSpeeds) error
}

type pageSpeedHandler struct {
	alertSystem        alert_system.AlertHandler
	pageSpeedRepo      repos.PageSpeedRepository
	dataCentersRepo    repos.DataCentersRepository
	projectRepo        repos.ProjectsRepository
	pageSpeedStatsRepo repos.PageSpeedStatsRepository
	taskPusher         push.TaskPusher
	agentHandler       AgentHandler
}

func NewPageSpeedHandler(alertSystem alert_system.AlertHandler, pageSpeedRepo repos.PageSpeedRepository, dataCentersRepo repos.DataCentersRepository, projectRepo repos.ProjectsRepository, pageSpeedStatsRepo repos.PageSpeedStatsRepository, taskPusher push.TaskPusher, agentHandler AgentHandler) PageSpeedHandler {
	return &pageSpeedHandler{
		alertSystem:        alertSystem,
		pageSpeedRepo:      pageSpeedRepo,
		dataCentersRepo:    dataCentersRepo,
		projectRepo:        projectRepo,
		pageSpeedStatsRepo: pageSpeedStatsRepo,
		taskPusher:         taskPusher,
		agentHandler:       agentHandler,
	}
}

func (e *pageSpeedHandler) ExecutePageSpeedRule(ctx context.Context, pageSpeedRules usecase_models.PageSpeeds) error {
	if len(pageSpeedRules.Scheduling.DataCentersIds) == 1 && pageSpeedRules.Scheduling.DataCentersIds[0] == 0 {
		datacenters, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err == nil {
			// random number between one and the number of datacenters
			pageSpeedRules.Scheduling.DataCentersIds = []int{}
			pageSpeedRules.Scheduling.DataCentersIds = append(pageSpeedRules.Scheduling.DataCentersIds, datacenters[rand.Intn(len(datacenters))].ID)
		}
	} else if len(pageSpeedRules.Scheduling.DataCentersIds) == 0 {
		datacenters, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err == nil {
			// random number between one and the number of datacenters
			pageSpeedRules.Scheduling.DataCentersIds = []int{}
			for _, value := range datacenters {
				pageSpeedRules.Scheduling.DataCentersIds = append(pageSpeedRules.Scheduling.DataCentersIds, value.ID)
			}
		}
	}

	isHeart := pageSpeedRules.Scheduling.IsHeartBeat
	currentTime := time.Now()
	for {
		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(pageSpeedRules.Scheduling.DataCentersIds))
		for _, dataC := range pageSpeedRules.Scheduling.DataCentersIds {
			go func(dataCenter int) {
				var urlsCalled []string
				for _, rule := range pageSpeedRules.PageSpeed {
					dataCenter, err := e.dataCentersRepo.GetDataCenter(ctx, dataCenter)
					if err != nil {
						log.Info("error on getting data center in executing page speed rule: ", err)
						waitGroup.Done()
						return
					}

					response, err := e.agentHandler.SendPageSpeed(ctx, dataCenter.Baseurl, usecase_models.AgentPageSpeedRequest{
						Url: rule.Url,
					})
					if err != nil {
						log.Info("error on sending page speed in executing rule: ", err)
						waitGroup.Done()
						return
					}

					if response.Status == 0 {
						urlsCalled = append(urlsCalled, rule.Url)
						err = e.pageSpeedStatsRepo.Write(ctx, time.Now(), repos.WritePageSpeedStatsOptions{
							ProjectId:    pageSpeedRules.Scheduling.ProjectId,
							PageSpeedId:  pageSpeedRules.Scheduling.PipelineId,
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
					urlsCalled = append(urlsCalled, rule.Url)
				}
				err := e.pageSpeedStatsRepo.Write(ctx, time.Now(), repos.WritePageSpeedStatsOptions{
					ProjectId:    pageSpeedRules.Scheduling.ProjectId,
					PageSpeedId:  pageSpeedRules.Scheduling.PipelineId,
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

		sessionIds, err := e.pageSpeedStatsRepo.GetLastNSessionsByPageSpeedId(ctx, 2, pageSpeedRules.Scheduling.PipelineId)
		if err != nil {
			log.Info("problem in fetching session ids: ", err.Error())
			break
		}

		if len(sessionIds) != 2 {
			break
		}
		newSession, err := e.pageSpeedStatsRepo.Read(ctx,
			[]string{"time", "success", "datacenter_id", "url"},
			repos.Filters{repos.Filter{Field: "session_id", Op: repos.FilterOpEq, Value: sessionIds[0]}},
			[]string{models.PageSpeedsStatRels.Datacenter},
		)
		oldSession, err := e.pageSpeedStatsRepo.Read(ctx,
			[]string{"time", "success", "datacenter_id", "url"},
			repos.Filters{repos.Filter{Field: "session_id", Op: repos.FilterOpEq, Value: sessionIds[1]}},
			[]string{models.PageSpeedsStatRels.Datacenter},
		)

		project, err := e.projectRepo.GetProjectWithLoads(ctx, pageSpeedRules.Scheduling.ProjectId)
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

		resolved, failed, newSuccess, oldSuccess := calculatePageSpeedStatState(newSession, oldSession)
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
				Type:                "pagespeed",
				State:               "down",
				ProjectId:           pageSpeedRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        pageSpeedRules.Scheduling.PipelineName,
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
				Type:                "pagespeed",
				State:               "up",
				ProjectId:           pageSpeedRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        pageSpeedRules.Scheduling.PipelineName,
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
				Type:                "pagespeed",
				State:               "diff",
				ProjectId:           pageSpeedRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        pageSpeedRules.Scheduling.PipelineName,
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

func calculatePageSpeedStatState(newSession, oldSession models.PageSpeedsStatSlice) ([]int, []int, bool, bool) {
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
