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

type NetCatHandler interface {
	ExecuteNetCatRule(ctx context.Context, NetCatRules usecase_models.NetCats) error
}

type netCatHandler struct {
	alertSystem     alert_system.AlertHandler
	netCatRepo      repos.NetCatRepository
	dataCentersRepo repos.DataCentersRepository
	projectRepo     repos.ProjectsRepository
	netCatStatsRepo repos.NetCatStatsRepository
	taskPusher      push.TaskPusher
	agentHandler    AgentHandler
}

func NewNetCatHandler(alertSystem alert_system.AlertHandler, netCatRepo repos.NetCatRepository, dataCentersRepo repos.DataCentersRepository, projectRepo repos.ProjectsRepository, netcatStatsRepo repos.NetCatStatsRepository, taskPusher push.TaskPusher, agentHandler AgentHandler) NetCatHandler {
	return &netCatHandler{alertSystem: alertSystem, netCatRepo: netCatRepo, dataCentersRepo: dataCentersRepo, projectRepo: projectRepo, netCatStatsRepo: netcatStatsRepo, taskPusher: taskPusher, agentHandler: agentHandler}
}

func (e *netCatHandler) ExecuteNetCatRule(ctx context.Context, netCatRules usecase_models.NetCats) error {
	if len(netCatRules.Scheduling.DataCentersIds) == 1 && netCatRules.Scheduling.DataCentersIds[0] == 0 {
		datacenters, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err == nil {
			// random number between one and the number of datacenters
			netCatRules.Scheduling.DataCentersIds = []int{}
			netCatRules.Scheduling.DataCentersIds = append(netCatRules.Scheduling.DataCentersIds, datacenters[rand.Intn(len(datacenters))].ID)
		}
	} else if len(netCatRules.Scheduling.DataCentersIds) == 0 {
		datacenters, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err == nil {
			// random number between one and the number of datacenters
			netCatRules.Scheduling.DataCentersIds = []int{}
			for _, value := range datacenters {
				netCatRules.Scheduling.DataCentersIds = append(netCatRules.Scheduling.DataCentersIds, value.ID)
			}
		}
	}
	isHeart := netCatRules.Scheduling.IsHeartBeat
	currentTime := time.Now()
	for {
		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(netCatRules.Scheduling.DataCentersIds))
		for _, dataC := range netCatRules.Scheduling.DataCentersIds {
			go func(dataCenter int) {
				var addressesCalled []string
				for _, rule := range netCatRules.NetCats {
					dataCenter, err := e.dataCentersRepo.GetDataCenter(ctx, dataCenter)
					if err != nil {
						log.Info("error on getting data center in executing net cat rule: ", err)
						waitGroup.Done()
						return
					}

					response, err := e.agentHandler.SendNetCat(ctx, dataCenter.Baseurl, usecase_models.AgentNetCatRequest{
						Address: rule.Address,
						Port:    rule.Port,
						Type:    rule.Type,
						TimeOut: rule.TimeOut,
					})
					if err != nil {
						log.Info("error on sending net cat in executing rule: ", err)
						waitGroup.Done()
						return
					}

					if response.Status == 0 {
						addressesCalled = append(addressesCalled, rule.Address)
						err = e.netCatStatsRepo.Write(ctx, time.Now(), repos.WriteNetCatStatsOptions{
							ProjectId:    netCatRules.Scheduling.ProjectId,
							NetCatId:     netCatRules.Scheduling.PipelineId,
							IsHeartBeat:  isHeart,
							Url:          strings.Join(addressesCalled, ","),
							DatacenterId: dataCenter.ID,
							Success:      0,
						})
						if err != nil {
							log.Info("error on writing curl report in executing rule: ", err)
						}

						waitGroup.Done()
						return
					}
					addressesCalled = append(addressesCalled, rule.Address)
				}
				err := e.netCatStatsRepo.Write(ctx, time.Now(), repos.WriteNetCatStatsOptions{
					ProjectId:    netCatRules.Scheduling.ProjectId,
					NetCatId:     netCatRules.Scheduling.PipelineId,
					IsHeartBeat:  isHeart,
					Url:          strings.Join(addressesCalled, ","),
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

		sessionIds, err := e.netCatStatsRepo.GetLastNSessionsByNetcatId(ctx, 2, netCatRules.Scheduling.PipelineId)
		if err != nil {
			log.Info("problem in fetching session ids: ", err.Error())
			break
		}

		if len(sessionIds) != 2 {
			break
		}
		newSession, err := e.netCatStatsRepo.Read(ctx,
			[]string{"time", "success", "datacenter_id", "url"},
			repos.Filters{repos.Filter{Field: "session_id", Op: repos.FilterOpEq, Value: sessionIds[0]}},
			[]string{models.NetCatsStatRels.Datacenter},
		)
		oldSession, err := e.netCatStatsRepo.Read(ctx,
			[]string{"time", "success", "datacenter_id", "url"},
			repos.Filters{repos.Filter{Field: "session_id", Op: repos.FilterOpEq, Value: sessionIds[1]}},
			[]string{models.NetCatsStatRels.Datacenter},
		)

		project, err := e.projectRepo.GetProjectWithLoads(ctx, netCatRules.Scheduling.ProjectId)
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

		resolved, failed, newSuccess, oldSuccess := calculateNetCatStatState(newSession, oldSession)
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
				Type:                "netcat",
				State:               "down",
				ProjectId:           netCatRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        netCatRules.Scheduling.PipelineName,
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
				Type:                "netcat",
				State:               "up",
				ProjectId:           netCatRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        netCatRules.Scheduling.PipelineName,
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
				Type:                "netcat",
				State:               "diff",
				ProjectId:           netCatRules.Scheduling.ProjectId,
				Username:            project.R.Account.Username.String,
				PipelineName:        netCatRules.Scheduling.PipelineName,
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

func calculateNetCatStatState(newSession, oldSession models.NetCatsStatSlice) ([]int, []int, bool, bool) {
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
