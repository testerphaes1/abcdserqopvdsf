package handlers

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/labstack/gommon/log"
	"strconv"
	"test-manager/monitoring"
	"test-manager/repos"
	"test-manager/tasks/push"
	"test-manager/tasks/task_models"
	"test-manager/usecase_models"
	"time"
)

type schedulerProvider struct {
	endpointRepository   repos.EndpointRepository
	netCatRepository     repos.NetCatRepository
	pageSpeedRepository  repos.PageSpeedRepository
	traceRouteRepository repos.TraceRouteRepository
	pingRepository       repos.PingRepository
}

func NewProvider(
	endpointRepository repos.EndpointRepository,
	netCatRepository repos.NetCatRepository,
	pageSpeedRepository repos.PageSpeedRepository,
	traceRouteRepository repos.TraceRouteRepository,
	pingRepository repos.PingRepository,
) asynq.PeriodicTaskConfigProvider {
	return &schedulerProvider{
		endpointRepository:   endpointRepository,
		netCatRepository:     netCatRepository,
		pageSpeedRepository:  pageSpeedRepository,
		traceRouteRepository: traceRouteRepository,
		pingRepository:       pingRepository,
	}
}

func (c *schedulerProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	endpoints, err := c.endpointRepository.GetActiveEndpoints(context.TODO())
	if err != nil {
		panic(err)
	}
	//netcats, err := c.netCatRepository.GetActiveNetCats(context.TODO())
	//if err != nil {
	//	panic(err)
	//}
	//pagespeeds, err := c.pageSpeedRepository.GetActivePageSpeeds(context.TODO())
	//if err != nil {
	//	panic(err)
	//}
	//pings, err := c.pingRepository.GetActivePings(context.TODO())
	//if err != nil {
	//	panic(err)
	//}
	//traceRoutes, err := c.traceRouteRepository.GetActiveTraceRoutes(context.TODO())
	//if err != nil {
	//	panic(err)
	//}

	var taskConfigs []*asynq.PeriodicTaskConfig
	for _, endpoint := range endpoints {
		payloadBytes, err := json.Marshal(endpoint)
		if err != nil {
			panic(err)
		}
		newTask := asynq.NewTask(task_models.TypeEndpoint, payloadBytes)

		if endpoint.Scheduling.IsHeartBeat {
			endpoint.Scheduling.Duration = 1
		}
		// create periodic task config
		taskConfigs = append(taskConfigs, &asynq.PeriodicTaskConfig{
			Cronspec: calcCronSpec(endpoint.Scheduling.Duration),
			Task:     newTask,
			Opts: []asynq.Option{
				asynq.Unique(24 * time.Hour),
				asynq.Queue(task_models.QueueEndpoint),
			},
		})
	}
	//for _, netcat := range netcats {
	//	payloadBytes, err := json.Marshal(netcat)
	//	if err != nil {
	//		panic(err)
	//	}
	//	newTask := asynq.NewTask(task_models.TypeNetCats, payloadBytes)
	//	if netcat.Scheduling.IsHeartBeat {
	//		netcat.Scheduling.Duration = 1
	//	}
	//	// create periodic task config
	//	taskConfigs = append(taskConfigs, &asynq.PeriodicTaskConfig{
	//		Cronspec: calcCronSpec(netcat.Scheduling.Duration),
	//		Task:     newTask,
	//		Opts: []asynq.Option{
	//			asynq.Unique(24 * time.Hour),
	//			asynq.Queue(task_models.QueueNetCats),
	//		},
	//	})
	//}
	//for _, pagespeed := range pagespeeds {
	//	payloadBytes, err := json.Marshal(pagespeed)
	//	if err != nil {
	//		panic(err)
	//	}
	//	newTask := asynq.NewTask(task_models.TypePageSpeeds, payloadBytes)
	//	if pagespeed.Scheduling.IsHeartBeat {
	//		pagespeed.Scheduling.Duration = 1
	//	}
	//	// create periodic task config
	//	taskConfigs = append(taskConfigs, &asynq.PeriodicTaskConfig{
	//		Cronspec: calcCronSpec(pagespeed.Scheduling.Duration),
	//		Task:     newTask,
	//		Opts: []asynq.Option{
	//			asynq.Unique(24 * time.Hour),
	//			asynq.Queue(task_models.QueuePageSpeeds),
	//		},
	//	})
	//}
	//for _, ping := range pings {
	//	payloadBytes, err := json.Marshal(ping)
	//	if err != nil {
	//		panic(err)
	//	}
	//	newTask := asynq.NewTask(task_models.TypePings, payloadBytes)
	//	if ping.Scheduling.IsHeartBeat {
	//		ping.Scheduling.Duration = 1
	//	}
	//	// create periodic task config
	//	taskConfigs = append(taskConfigs, &asynq.PeriodicTaskConfig{
	//		Cronspec: calcCronSpec(ping.Scheduling.Duration),
	//		Task:     newTask,
	//		Opts: []asynq.Option{
	//			asynq.Unique(24 * time.Hour),
	//			asynq.Queue(task_models.QueuePings),
	//		},
	//	})
	//}
	//for _, traceRoute := range traceRoutes {
	//	payloadBytes, err := json.Marshal(traceRoute)
	//	if err != nil {
	//		panic(err)
	//	}
	//	newTask := asynq.NewTask(task_models.TypeTraceRoutes, payloadBytes)
	//	if traceRoute.Scheduling.IsHeartBeat {
	//		traceRoute.Scheduling.Duration = 1
	//	}
	//	// create periodic task config
	//	taskConfigs = append(taskConfigs, &asynq.PeriodicTaskConfig{
	//		Cronspec: calcCronSpec(traceRoute.Scheduling.Duration),
	//		Task:     newTask,
	//		Opts: []asynq.Option{
	//			asynq.Unique(24 * time.Hour),
	//			asynq.Queue(task_models.QueueTraceRoutes),
	//		},
	//	})
	//}

	monitoring.ActiveSchedulerTasksGauge.Set(float64(len(taskConfigs)))
	return taskConfigs, nil
}

func calcCronSpec(duration int) string {
	// duration is in minutes
	if duration > 59 {
		minute := duration % 60
		hour := duration / 60
		return "*/" + strconv.Itoa(minute) + " */" + strconv.Itoa(hour) + " * * *"
	}
	return "*/" + strconv.Itoa(duration) + " * * * *"
}

type HeartBeatScheduler struct {
	endpointRepository repos.EndpointRepository
	taskPusher         push.TaskPusher
	inspector          *asynq.Inspector
}

func NewHeartBeatScheduler(
	endpointRepository repos.EndpointRepository,
	taskPusher push.TaskPusher,
	inspector *asynq.Inspector,
) *HeartBeatScheduler {
	return &HeartBeatScheduler{
		endpointRepository: endpointRepository,
		taskPusher:         taskPusher,
		inspector:          inspector,
	}
}

func (c *HeartBeatScheduler) HeartBeatScheduling(ctx context.Context) {
	go func() {
		var contextMap = make(map[int]context.CancelFunc)
		var oldActiveEndpoints = new([]*usecase_models.Endpoints)
		var newActiveEndpoints *[]*usecase_models.Endpoints
		for {
			e, err := c.endpointRepository.GetActiveEndpoints(ctx)
			if err != nil {
				log.Error(err)
				time.Sleep(5 * time.Second)
				continue
			}
			newActiveEndpoints = &e

			var activeEndpoints = new([]*usecase_models.Endpoints)
			for _, newActive := range *newActiveEndpoints {
				flag := false
				for _, oldActive := range *oldActiveEndpoints {
					if newActive.Scheduling.PipelineId == oldActive.Scheduling.PipelineId {
						flag = true
						break
					}
				}
				if !flag {
					*activeEndpoints = append(*activeEndpoints, newActive)
				}
			}
			var deActiveEndpoints = new([]*usecase_models.Endpoints)
			for _, oldActive := range *oldActiveEndpoints {
				flag := false
				for _, newActive := range *newActiveEndpoints {
					if oldActive.Scheduling.PipelineId == newActive.Scheduling.PipelineId {
						flag = true
						break
					}
				}
				if !flag {
					*deActiveEndpoints = append(*deActiveEndpoints, oldActive)
				}
			}

			for _, endpoint := range *activeEndpoints {
				hCtx, hCancel := context.WithCancel(context.Background())
				contextMap[endpoint.Scheduling.PipelineId] = hCancel
				go heartBeatHandler(hCtx, endpoint, c.taskPusher, c.inspector)
			}
			if len(*activeEndpoints) > 0 {
				log.Infof("successfully registered %d new heart beat scheduler", len(*activeEndpoints))
			}
			for _, endpoint := range *deActiveEndpoints {
				// cancel func
				contextMap[endpoint.Scheduling.PipelineId]()
				delete(contextMap, endpoint.Scheduling.PipelineId)
			}
			if len(*deActiveEndpoints) > 0 {
				log.Infof("successfully deleted %d heart beat scheduler", len(*deActiveEndpoints))
			}
			oldActiveEndpoints = newActiveEndpoints
			monitoring.ActiveSchedulerTasksGauge.Set(float64(len(*newActiveEndpoints)))
			time.Sleep(10 * time.Second)
		}
	}()
}

func heartBeatHandler(ctx context.Context, endpoint *usecase_models.Endpoints,
	taskPusher push.TaskPusher, inspector *asynq.Inspector) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			handleTaskExecuting(ctx, endpoint, taskPusher, inspector)
			time.Sleep(15 * time.Second)
		}
	}
}

func handleTaskExecuting(ctx context.Context, endpoint *usecase_models.Endpoints,
	taskPusher push.TaskPusher, inspector *asynq.Inspector) {
	taskContext, cf := context.WithTimeout(ctx, 60*time.Second)
	defer cf()
	taskId, err := taskPusher.PushEndpoint(ctx, *endpoint)
	if err != nil {
		log.Error(err)
	}
	for {
		select {
		case <-taskContext.Done():
			err = inspector.DeleteTask(task_models.QueueEndpoint, taskId)
			if err != nil {
				log.Error(err)
			}
			return
		default:
			var err error
			taskInfo, err := inspector.GetTaskInfo(task_models.QueueEndpoint, taskId)
			if err != nil {
				log.Error(err)
				return
			}
			if taskInfo.State == asynq.TaskStateCompleted {
				log.Infof("task completed: %s", taskInfo.ID)
				return
			}
		}
		time.Sleep(1 * time.Second)
	}
}
