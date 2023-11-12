package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"test-manager/cache"
	"test-manager/monitoring"
	"test-manager/repos"
	"test-manager/services/alert_system"
	"test-manager/tasks/push"
	"test-manager/tasks/task_models"
	"test-manager/usecase_models"
	"time"
	"unicode/utf8"
)

type EndpointHandler interface {
	ExecuteEndpointRule(ctx context.Context, endpointRules usecase_models.Endpoints) error
}

type endpointHandler struct {
	alertSystem     alert_system.AlertHandler
	endpointRepo    repos.EndpointRepository
	dataCentersRepo repos.DataCentersRepository
	projectRepo     repos.ProjectsRepository
	endpointStats   repos.EndpointStatsRepository
	cacheRepo       cache.Cache
	taskPusher      push.TaskPusher
	agentHandler    AgentHandler
}

func NewEndpointHandler(alertSystem alert_system.AlertHandler,
	endpointRepo repos.EndpointRepository,
	dataCentersRepo repos.DataCentersRepository,
	projectRepo repos.ProjectsRepository,
	endpointStats repos.EndpointStatsRepository,
	cacheRepo cache.Cache,
	taskPusher push.TaskPusher,
	agentHandler AgentHandler) EndpointHandler {
	return &endpointHandler{
		alertSystem:     alertSystem,
		endpointRepo:    endpointRepo,
		dataCentersRepo: dataCentersRepo,
		endpointStats:   endpointStats,
		projectRepo:     projectRepo,
		cacheRepo:       cacheRepo,
		taskPusher:      taskPusher,
		agentHandler:    agentHandler,
	}
}

func (e *endpointHandler) ExecuteEndpointRule(ctx context.Context, endpointRules usecase_models.Endpoints) error {
	if len(endpointRules.Scheduling.DataCentersIds) == 1 && endpointRules.Scheduling.DataCentersIds[0] == 0 {
		datacenters, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err == nil {
			// random number between one and the number of datacenters
			endpointRules.Scheduling.DataCentersIds = []int{}
			endpointRules.Scheduling.DataCentersIds = append(endpointRules.Scheduling.DataCentersIds, datacenters[rand.Intn(len(datacenters))].ID)
		}
	} else if len(endpointRules.Scheduling.DataCentersIds) == 0 {
		datacenters, err := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err == nil {
			// all datacenters
			endpointRules.Scheduling.DataCentersIds = []int{}
			for _, value := range datacenters {
				endpointRules.Scheduling.DataCentersIds = append(endpointRules.Scheduling.DataCentersIds, value.ID)
			}
		}
	}
	session, _ := uuid.NewUUID()
	sessionIsValid := make(chan bool, 1)
	sessionIsValid <- true
	sessionResults := make(chan repos.WriteEndpointStatsOptions, len(endpointRules.Scheduling.DataCentersIds))
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(endpointRules.Scheduling.DataCentersIds))
	for _, dataC := range endpointRules.Scheduling.DataCentersIds {
		go func(dataCenterId int, sessionResults chan repos.WriteEndpointStatsOptions, waitGroup *sync.WaitGroup) {
			sessionDurationTime := prometheus.NewTimer(monitoring.EndpointTaskDatacenterDuration.WithLabelValues(strconv.Itoa(dataCenterId)))
			sessionR := repos.WriteEndpointStatsOptions{}
			defer func() {
				sessionDurationTime.ObserveDuration()
				waitGroup.Done()
				sessionR.ResponseBodies = ""
				sessionResults <- sessionR
			}()
			var responses = usecase_models.EndpointResponses{
				HeaderResponses: map[string]map[string][]string{},
				BodyResponses:   map[string]string{},
				TimeResponses:   map[string]float64{},
				StatusResponses: map[string]int{},
			}
			avgResTime := float64(0)
			var urlsCalled []string
			var endpointNamesCalled []string
			for i, rule := range endpointRules.Endpoints {
				urlsCalled = append(urlsCalled, rule.Url)
				endpointNamesCalled = append(endpointNamesCalled, rule.EndpointName)
				dataCenter, err := e.dataCentersRepo.GetDataCenterWithCache(ctx, dataCenterId)
				if err != nil {
					log.Error("error on getting data center in executing endpoint rule: ", err)
					monitoring.ErrorEndpointTaskFetchingDatacenterGauge.Inc()
					sessionIsValid <- false
					return
				}

				var newHeader = map[string][]string{}
				if i != 0 { // iterate over body and replace with values from responses
					tempKeys := strings.Split(rule.Body, "{{")
					for _, tempKey := range tempKeys {
						identifier := strings.Split(tempKey, "}}")[0]

						actions := strings.Split(identifier, ".")
						switch actions[0] {
						case "header":
							if headerResponse, ok := responses.HeaderResponses[actions[1]]; ok {
								key := strings.Join(actions[2:], ".")
								rule.Body = strings.ReplaceAll(rule.Body, "{{"+identifier+"}}", strings.Join(headerResponse[key], ","))
							} else {
								err = e.endpointRepo.DisableEndpointWithFunctionalityReason(ctx, endpointRules.Scheduling.PipelineId)
								if err != nil {
									log.Error("error on disabling endpoint in executing endpoint rule: ", err)
								}
								log.Info("error on key mapping in executing endpoint header rule: ", identifier)
								monitoring.DisabledEndpointRuleReasonFunctionality.Inc()
								sessionIsValid <- false
								return
							}
						case "body":
							if bodyResponse, ok := responses.BodyResponses[actions[1]]; ok {
								key := strings.Join(actions[2:], ".")
								value := gjson.Get(bodyResponse, key).String()
								rule.Body = strings.ReplaceAll(rule.Body, "{{"+identifier+"}}", value)
							} else {
								err = e.endpointRepo.DisableEndpointWithFunctionalityReason(ctx, endpointRules.Scheduling.PipelineId)
								if err != nil {
									log.Error("error on disabling endpoint in executing endpoint rule: ", err)
								}
								log.Info("error on key mapping in executing endpoint body rule: ", identifier)
								monitoring.DisabledEndpointRuleReasonFunctionality.Inc()
								sessionIsValid <- false
								return
							}
						}
					}

					// iterate over header and replace with values from responses
					for key, value := range rule.Header {
						if strings.Contains(key, "{{") && strings.Contains(key, "}}") {
							tempKeys := strings.Split(key, "{{")
							for _, tempKey := range tempKeys {
								identifier := strings.Split(tempKey, "}}")[0]
								actions := strings.Split(identifier, ".")
								switch actions[0] {
								case "header":
									if headerResponse, ok := responses.HeaderResponses[actions[1]]; ok {
										key := strings.Join(actions[2:], ".")
										rule.Header[key] = strings.Join(headerResponse[key], ",")
									} else {
										err = e.endpointRepo.DisableEndpointWithFunctionalityReason(ctx, endpointRules.Scheduling.PipelineId)
										if err != nil {
											log.Error("error on disabling endpoint in executing endpoint rule: ", err)
										}
										log.Info("error on key mapping in executing endpoint header rule: ", identifier)
										monitoring.DisabledEndpointRuleReasonFunctionality.Inc()
										sessionIsValid <- false
										return
									}
								case "body":
									if bodyResponse, ok := responses.BodyResponses[actions[1]]; ok {
										key := strings.Join(actions[2:], ".")
										value := gjson.Get(bodyResponse, key).String()
										rule.Header[key] = value
									} else {
										err = e.endpointRepo.DisableEndpointWithFunctionalityReason(ctx, endpointRules.Scheduling.PipelineId)
										if err != nil {
											log.Error("error on disabling endpoint in executing endpoint rule: ", err)
										}
										log.Info("error on key mapping in executing endpoint body rule: ", identifier)
										monitoring.DisabledEndpointRuleReasonFunctionality.Inc()
										sessionIsValid <- false
										return
									}
								}
							}
						}
						if strings.Contains(value, "{{") && strings.Contains(value, "}}") {
							tempKeys := strings.Split(value, "{{")
							for _, tempKey := range tempKeys {
								identifier := strings.Split(tempKey, "}}")[0]
								actions := strings.Split(identifier, ".")
								switch actions[0] {
								case "header":
									if headerResponse, ok := responses.HeaderResponses[actions[1]]; ok {
										key := strings.Join(actions[2:], ".")
										rule.Header[key] = strings.Join(headerResponse[key], ",")
									} else {
										err = e.endpointRepo.DisableEndpointWithFunctionalityReason(ctx, endpointRules.Scheduling.PipelineId)
										if err != nil {
											log.Error("error on disabling endpoint in executing endpoint rule: ", err)
										}
										log.Info("error on key mapping in executing endpoint header rule: ", identifier)
										monitoring.DisabledEndpointRuleReasonFunctionality.Inc()
										sessionIsValid <- false
										return
									}
								case "body":
									if bodyResponse, ok := responses.BodyResponses[actions[1]]; ok {
										key := strings.Join(actions[2:], ".")
										value := gjson.Get(bodyResponse, key).String()
										rule.Header[key] = value
									} else {
										err = e.endpointRepo.DisableEndpointWithFunctionalityReason(ctx, endpointRules.Scheduling.PipelineId)
										if err != nil {
											log.Error("error on disabling endpoint in executing endpoint rule: ", err)
										}
										log.Info("error on key mapping in executing endpoint body rule: ", identifier)
										monitoring.DisabledEndpointRuleReasonFunctionality.Inc()
										sessionIsValid <- false
										return
									}
								}
							}
						}
					}

					// iterate over url and replace with values from responses
					tempKeys = strings.Split(rule.Url, "{{")
					for _, tempKey := range tempKeys {
						identifier := strings.Split(tempKey, "}}")[0]
						actions := strings.Split(identifier, ".")
						switch actions[0] {
						case "header":
							if headerResponse, ok := responses.HeaderResponses[actions[1]]; ok {
								key := strings.Join(actions[2:], ".")
								rule.Url = strings.ReplaceAll(rule.Url, "{{"+identifier+"}}", strings.Join(headerResponse[key], ","))
							} else {
								err = e.endpointRepo.DisableEndpointWithFunctionalityReason(ctx, endpointRules.Scheduling.PipelineId)
								if err != nil {
									log.Error("error on disabling endpoint in executing endpoint rule: ", err)
								}
								log.Info("error on key mapping in executing endpoint header rule: ", identifier)
								monitoring.DisabledEndpointRuleReasonFunctionality.Inc()
								sessionIsValid <- false
								return
							}
						case "body":
							if bodyResponse, ok := responses.BodyResponses[actions[1]]; ok {
								key := strings.Join(actions[2:], ".")
								value := gjson.Get(bodyResponse, key).String()
								rule.Url = strings.ReplaceAll(rule.Url, "{{"+identifier+"}}", value)
							} else {
								err = e.endpointRepo.DisableEndpointWithFunctionalityReason(ctx, endpointRules.Scheduling.PipelineId)
								if err != nil {
									log.Error("error on disabling endpoint in executing endpoint rule: ", err)
								}
								log.Info("error on key mapping in executing endpoint body rule: ", identifier)
								monitoring.DisabledEndpointRuleReasonFunctionality.Inc()
								sessionIsValid <- false
								return
							}
						}
					}

					for key, val := range rule.Header {
						newHeader[key] = strings.Split(val, ",")
					}
				}

				respBody, respHeader, respStatus, respTime, err := e.agentHandler.SendCurl(ctx, dataCenter.Baseurl, usecase_models.AgentCurlRequest{
					Url:    rule.Url,
					Method: rule.Method,
					Header: newHeader,
					Body:   rule.Body,
				})
				if err != nil {
					log.Info("error on sending curl in executing rule: ", err)
					responses.BodyResponses[rule.EndpointName] = fmt.Sprintf("error on sending curl in executing rule: %s", err.Error())
					responses.HeaderResponses[rule.EndpointName] = map[string][]string{}
					responses.TimeResponses[rule.EndpointName] = 0
					responses.StatusResponses[rule.EndpointName] = 0

					rt, _ := json.Marshal(responses.TimeResponses)
					rb, _ := json.Marshal(responses.BodyResponses)
					rh, _ := json.Marshal(responses.HeaderResponses)
					rs, _ := json.Marshal(responses.StatusResponses)
					sessionR = repos.WriteEndpointStatsOptions{
						Time:             time.Now(),
						SessionId:        session.String(),
						ProjectId:        endpointRules.Scheduling.ProjectId,
						EndpointName:     strings.Join(endpointNamesCalled, ","),
						EndpointId:       endpointRules.Scheduling.PipelineId,
						IsHeartBeat:      endpointRules.Scheduling.IsHeartBeat,
						Url:              strings.Join(urlsCalled, ","),
						DatacenterId:     dataCenter.ID,
						Success:          0,
						ResponseTime:     avgResTime,
						ResponseTimes:    string(rt),
						ResponseBodies:   string(rb),
						ResponseHeaders:  string(rh),
						ResponseStatuses: string(rs),
					}
					_, err = e.taskPusher.PushEndpointStore(ctx, sessionR)
					return
				}
				t, err := base64.StdEncoding.DecodeString(respBody)
				if err == nil {
					respBody = string(t)
				}

				limit := 100 // 1 megabyte

				truncatedRespBody := forceStringSize(respBody, limit)
				responses.BodyResponses[rule.EndpointName] = truncatedRespBody
				responses.HeaderResponses[rule.EndpointName] = respHeader
				responses.TimeResponses[rule.EndpointName] = respTime
				responses.StatusResponses[rule.EndpointName] = respStatus
				if !curlAcceptanceCriteria(strconv.Itoa(respStatus), []byte(respBody), rule.AcceptanceModel) {
					avgResTime = float64(0)
					c := float64(0)
					for _, value := range responses.TimeResponses {
						if value != 0 {
							avgResTime += value
							c += 1
						}
					}
					if c != 0 {
						avgResTime = avgResTime / c
					}
					rt, _ := json.Marshal(responses.TimeResponses)
					rb, _ := json.Marshal(responses.BodyResponses)
					rh, _ := json.Marshal(responses.HeaderResponses)
					rs, _ := json.Marshal(responses.StatusResponses)
					sessionR = repos.WriteEndpointStatsOptions{
						Time:             time.Now(),
						SessionId:        session.String(),
						ProjectId:        endpointRules.Scheduling.ProjectId,
						EndpointName:     strings.Join(endpointNamesCalled, ","),
						EndpointId:       endpointRules.Scheduling.PipelineId,
						IsHeartBeat:      endpointRules.Scheduling.IsHeartBeat,
						Url:              strings.Join(urlsCalled, ","),
						DatacenterId:     dataCenter.ID,
						Success:          0,
						ResponseTime:     avgResTime,
						ResponseTimes:    string(rt),
						ResponseBodies:   string(rb),
						ResponseHeaders:  string(rh),
						ResponseStatuses: string(rs),
					}
					_, err = e.taskPusher.PushEndpointStore(ctx, sessionR)
					if err != nil {
						log.Info("error on writing curl report in executing rule: ", err)
					}

					return
				}
			}

			{
				avgResTime = float64(0)
				c := float64(0)
				for _, value := range responses.TimeResponses {
					if value != 0 {
						avgResTime += value
						c += 1
					}
				}
				if c != 0 {
					avgResTime = avgResTime / c
				}
				rt, _ := json.Marshal(responses.TimeResponses)
				rb, _ := json.Marshal(responses.BodyResponses)
				rh, _ := json.Marshal(responses.HeaderResponses)
				rs, _ := json.Marshal(responses.StatusResponses)
				sessionR = repos.WriteEndpointStatsOptions{
					Time:             time.Now(),
					SessionId:        session.String(),
					ProjectId:        endpointRules.Scheduling.ProjectId,
					EndpointName:     strings.Join(endpointNamesCalled, ","),
					EndpointId:       endpointRules.Scheduling.PipelineId,
					IsHeartBeat:      endpointRules.Scheduling.IsHeartBeat,
					Url:              strings.Join(urlsCalled, ","),
					DatacenterId:     dataCenterId,
					Success:          1,
					ResponseTime:     avgResTime,
					ResponseTimes:    string(rt),
					ResponseBodies:   string(rb),
					ResponseHeaders:  string(rh),
					ResponseStatuses: string(rs),
				}
				_, err := e.taskPusher.PushEndpointStore(ctx, sessionR)
				if err != nil {
					log.Info("error on pushing curl report in executing rule: ", err)
				}
				return
			}
		}(dataC, sessionResults, &waitGroup)
	}
	waitGroup.Wait()

	if !<-sessionIsValid {
		log.Warn("this session was invalid: ", session.String())
		return nil
	}

	var newSession []repos.WriteEndpointStatsOptions
	for i := 0; i < len(endpointRules.Scheduling.DataCentersIds); i++ {
		newSession = append(newSession, <-sessionResults)
	}

	// rosin means nothing
	rosin, err := e.cacheRepo.Get(ctx, repos.EndpointSessionCachePrefix+strconv.Itoa(endpointRules.Scheduling.PipelineId))
	if err != nil {
		log.Warn("could not find old session on cache: ", err)
	}
	temp, _ := json.Marshal(newSession)
	err = e.cacheRepo.Set(ctx, repos.EndpointSessionCachePrefix+strconv.Itoa(endpointRules.Scheduling.PipelineId),
		temp, time.Duration(2*endpointRules.Scheduling.Duration)*time.Minute)
	if err != nil {
		log.Warn("problem on setting new session on cache: ", err)
	}

	if rosin == "" {
		return errors.New(fmt.Sprintf("first cycle of session started on endpoint id: %d", endpointRules.Scheduling.PipelineId))
	}

	var oldSession []repos.WriteEndpointStatsOptions
	rosinStr, ok := rosin.(string)
	if !ok {
		log.Error("problem on casting old session: ")
		return errors.New(fmt.Sprintf("problem on casting old session: %v", rosin))
	}
	err = json.Unmarshal([]byte(rosinStr), &oldSession)
	if err != nil {
		log.Error("problem on unmarshalling rosinByte: ", err)
		return errors.New(fmt.Sprintf("problem on casting old session: %v", rosinStr))
	}

	resolved, failed, newSuccess, oldSuccess := calculateEndpointStatState(newSession, oldSession)
	if oldSuccess && !newSuccess {
		project, err2 := e.projectRepo.GetProjectWithLoads(ctx, endpointRules.Scheduling.ProjectId)
		if err2 != nil {
			log.Info("error on getting project in executing rule: ", err2)
			return err2
		}
		var notifications usecase_models.Notifications
		err2 = json.Unmarshal(project.Notifications.JSON, &notifications)
		if err2 != nil {
			log.Info("can not unmarshal notification")
			return err2
		}

		// send fail
		address := ""
		rootCause := "working on it"
		for _, value := range oldSession {
			if value.Success == 1 {
				t := strings.Split(value.Url, ",")
				address = t[len(t)-1]
				break
			}
		}
		for _, value := range newSession {
			if value.Success == 0 {
				rootCause = value.ResponseStatuses
			}
		}

		dcs, err2 := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err2 != nil {
			log.Info("problem in getting datacenters in sending alerts: ", err2.Error())
		}
		var dcTitles []string
		for _, value := range failed {
			for _, dc := range dcs {
				if dc.ID == value {
					dcTitles = append(dcTitles, dc.Title)
				}
			}
		}
		_, err2 = e.taskPusher.PushNotifications(ctx, task_models.NotificationsPayload{
			Type:                "endpoint",
			State:               "down",
			ProjectId:           endpointRules.Scheduling.ProjectId,
			Username:            project.R.Account.Username.String,
			PipelineName:        endpointRules.Scheduling.PipelineName,
			Time:                time.Now().String(),
			Address:             address,
			Datacenters:         strings.Join(dcTitles, ","),
			RootCause:           rootCause,
			ResolvedDatacenters: "",
			FailedDatacenters:   "",
		})
		if err2 != nil {
			log.Info(err2.Error())
		}
	}
	if !oldSuccess && newSuccess {
		project, err2 := e.projectRepo.GetProjectWithLoads(ctx, endpointRules.Scheduling.ProjectId)
		if err2 != nil {
			log.Info("error on getting project in executing rule: ", err2)
			return err2
		}
		var notifications usecase_models.Notifications
		err2 = json.Unmarshal(project.Notifications.JSON, &notifications)
		if err2 != nil {
			log.Info("can not unmarshal notification")
			return err2
		}

		// send success
		address := ""
		rootCause := "working on it"
		for _, value := range oldSession {
			if value.Success == 0 {
				rootCause = value.ResponseStatuses
				t := strings.Split(value.Url, ",")
				address = t[len(t)-1]
				break
			}
		}
		dcs, err2 := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err2 != nil {
			log.Info("problem in getting datacenters in sending alerts: ", err2.Error())
		}
		var dcTitles []string
		for _, value := range resolved {
			for _, dc := range dcs {
				if dc.ID == value {
					dcTitles = append(dcTitles, dc.Title)
				}
			}
		}
		_, err2 = e.taskPusher.PushNotifications(ctx, task_models.NotificationsPayload{
			Type:                "endpoint",
			State:               "up",
			ProjectId:           endpointRules.Scheduling.ProjectId,
			Username:            project.R.Account.Username.String,
			PipelineName:        endpointRules.Scheduling.PipelineName,
			Time:                time.Now().String(),
			Address:             address,
			Datacenters:         strings.Join(dcTitles, ","),
			RootCause:           rootCause,
			ResolvedDatacenters: "",
			FailedDatacenters:   "",
		})
		if err2 != nil {
			log.Info(err2.Error())
		}
	}
	if !oldSuccess && !newSuccess && (len(resolved) > 0 || len(failed) > 0) {
		project, err2 := e.projectRepo.GetProjectWithLoads(ctx, endpointRules.Scheduling.ProjectId)
		if err2 != nil {
			log.Info("error on getting project in executing rule: ", err2)
			return err2
		}
		var notifications usecase_models.Notifications
		err2 = json.Unmarshal(project.Notifications.JSON, &notifications)
		if err2 != nil {
			log.Info("can not unmarshal notification")
			return err2
		}

		// send diff
		address := ""
		rootCause := "working on it"
		for _, value := range oldSession {
			rootCause = value.ResponseStatuses
			t := strings.Split(value.Url, ",")
			address = t[len(t)-1]
			break
		}
		dcs, err2 := e.dataCentersRepo.GetDataCentersWithCache(ctx)
		if err2 != nil {
			log.Info("problem in getting datacenters in sending alerts: ", err2.Error())
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
		_, err2 = e.taskPusher.PushNotifications(ctx, task_models.NotificationsPayload{
			Type:                "endpoint",
			State:               "diff",
			ProjectId:           endpointRules.Scheduling.ProjectId,
			Username:            project.R.Account.Username.String,
			PipelineName:        endpointRules.Scheduling.PipelineName,
			Time:                time.Now().String(),
			Address:             address,
			Datacenters:         "",
			RootCause:           rootCause,
			ResolvedDatacenters: strings.Join(rTitles, ","),
			FailedDatacenters:   strings.Join(fTitles, ","),
		})
		if err2 != nil {
			log.Info(err2.Error())
		}
	}
	return nil
}

func curlAcceptanceCriteria(status string, body []byte, acceptRules usecase_models.AcceptanceModel) bool {
	statusCheck := false
	for _, val := range acceptRules.Statuses {
		if val == status {
			statusCheck = true
			break
		}
	}
	if !statusCheck {
		return false
	}

	var respbody map[string]interface{}
	json.Unmarshal(body, &respbody)

	if len(acceptRules.ResponseBodies) == 0 {
		return true
	}
	bodyCheck := true
	for _, val := range acceptRules.ResponseBodies {
		_, ok := respbody[val.Key]
		if !ok {
			bodyCheck = false
			break
		}
		//if reflect.TypeOf(value).String() != val.Value {
		//	bodyCheck = false
		//	break
		//}
	}

	if !bodyCheck {
		return false
	}

	return true
}

func calculateEndpointStatState(newSession, oldSession []repos.WriteEndpointStatsOptions) ([]int, []int, bool, bool) {
	var mapNDcSuccess = make(map[int]bool)
	var mapODcSuccess = make(map[int]bool)
	var newSuccess = true
	var oldSuccess = true
	for _, value := range newSession {
		mapNDcSuccess[value.DatacenterId] = value.Success != 0
	}
	for _, value := range oldSession {
		if _, ok := mapNDcSuccess[value.DatacenterId]; !ok {
			mapNDcSuccess[value.DatacenterId] = true
		}
		mapODcSuccess[value.DatacenterId] = value.Success != 0
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

func GetStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}

func forceStringSize(s string, maxSize int) string {
	if len(s) <= maxSize {
		return s
	}

	// Truncate the string to fit the maximum size
	runes := []rune(s)
	truncatedRunes := runes[:maxSize]
	truncatedString := string(truncatedRunes)

	// Check if the truncated string is valid UTF-8
	for len(truncatedString) > 0 && !utf8.ValidString(truncatedString) {
		// If it's not valid UTF-8, remove the last rune and try again
		truncatedRunes = truncatedRunes[:len(truncatedRunes)-1]
		truncatedString = string(truncatedRunes)
	}

	return truncatedString
}
