package cmd

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"net/http"
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"test-manager/cache"
	"test-manager/config"
	"test-manager/handlers"
	"test-manager/repos"
	"test-manager/services/alert_system"
	"test-manager/tasks"
	"test-manager/tasks/push"
	"test-manager/tasks/task_models"
	"test-manager/utils"
	"time"
)

const (
	numWorkersAsynq = 100
)

var (
	// list of queues associated with priority, large numbers indicate higher priority
	queues = map[string]int{
		task_models.QueueEndpoint: 6,
		//task_models.QueueNetCats:       0,
		//task_models.QueuePageSpeeds:    0,
		//task_models.QueuePings:         0,
		//task_models.QueueTraceRoutes:   0,
		task_models.QueueNotification:  6,
		task_models.QueueEndpointStore: 6,
	}
)

func init() {
	rootCmd.AddCommand(consumeTasksCmd)
}

var consumeTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "consume async tasks events",
	Long:  `consume async tasks events`,
	Run: func(cmd *cobra.Command, args []string) {
		l := utils.ZapLogger()
		zLogger := l.Sugar()

		_, config, err := config.ViperConfig()
		if err != nil {
			panic(err)
		}

		redisClient, err := utils.CreateRedisConnection(context.TODO(), config.Database.Redis.Host, config.Database.Redis.Port, config.Database.Redis.Database, time.Duration(config.Database.Redis.Timeout)*time.Second)
		if err != nil {
			panic(err)
		}
		defer redisClient.Close()

		redisCacheClient, err := utils.CreateRedisConnection(context.TODO(), config.Database.RedisCache.Host, config.Database.RedisCache.Port, config.Database.RedisCache.Database, time.Duration(config.Database.RedisCache.Timeout)*time.Second)
		if err != nil {
			panic(err)
		}
		defer redisCacheClient.Close()

		cacheRepo := cache.NewRedisCache(redisCacheClient)

		//go func() {
		//	http.ListenAndServe("localhost:8881", nil)
		//}()

		go func() {
			e := echo.New()

			p := prometheus.NewPrometheus("automated_test_tasks", nil)
			p.Use(e)
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			// Start server
			go func() {
				if err := e.Start(":10002"); err != nil && err != http.ErrServerClosed {
					log.Fatal("shutting down server")
				}
			}()

			// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				e.Logger.Fatal(err)
			}
			//srv = &http.Server{Addr: ":10000"}
			//srv.Handler = http.DefaultServeMux
			//
			//http.Handle("/metrics", promhttp.Handler())
			//
			//if err := srv.ListenAndServe(); err != nil {
			//	if v, ok := err.(*net.OpError); ok {
			//		zLogger.Errorf("error on serving metrics server: %s", v)
			//	}
			//}
		}()
		tasks.ReportMetrics(10*time.Second, redisClient.Options())
		boil.DebugMode = false

		psqlDb, err := utils.PostgresConnection(config.Database.Pslq.Host, config.Database.Pslq.Port, config.Database.Pslq.User, config.Database.Pslq.Password, config.Database.Pslq.Database, config.Database.Pslq.Ssl, 2, 2)
		if err != nil {
			panic(err)
		}
		defer psqlDb.Close()

		srv := asynq.NewServer(
			asynq.RedisClientOpt{
				Addr:        redisClient.Options().Addr,
				DialTimeout: redisClient.Options().DialTimeout,
				Username:    redisClient.Options().Username,
				Password:    redisClient.Options().Password,
			}, asynq.Config{
				Concurrency: numWorkersAsynq,
				Logger:      zLogger,
				Queues:      queues,
			},
		)

		asynqClient := asynq.NewClient(asynq.RedisClientOpt{
			Addr:        redisClient.Options().Addr,
			DialTimeout: redisClient.Options().DialTimeout,
			Username:    redisClient.Options().Username,
			Password:    redisClient.Options().Password,
		})
		taskPusher := push.NewTaskPush(asynqClient)

		// alert system
		alertSystem := alert_system.NewAlertHandler(config.Services.Alert.BaseUrl)

		projectRepo := repos.NewProjectsRepository(psqlDb)
		endpointRepo := repos.NewEndpointRepository(psqlDb)
		//netCatRepo := repos.NewNetCatRepository(psqlDb)
		//pageSpeedRepo := repos.NewPageSpeedRepository(psqlDb)
		//pingRepo := repos.NewPingRepository(psqlDb)
		//traceRouteRepo := repos.NewTraceRouteRepository(psqlDb)
		dataCenterRepo := repos.NewDataCentersRepositoryRepository(cacheRepo, psqlDb)
		endpointStatsRepo := repos.NewEndpointStatsRepository(psqlDb)
		//netCatStatsRepo := repos.NewNetCatStatsRepository(psqlDb)
		//pageSPeedStatsRepo := repos.NewPageSpeedStatsRepository(psqlDb)
		//pingStatsRepo := repos.NewPingStatsRepository(psqlDb)
		//traceRouteStatsRepo := repos.NewTraceRouteStatsRepository(psqlDb)

		agentHandler := handlers.NewAgentHandler()
		endpointHandler := handlers.NewEndpointHandler(alertSystem, endpointRepo, dataCenterRepo,
			projectRepo, endpointStatsRepo, cacheRepo, taskPusher, agentHandler)
		//netCatHandler := handlers.NewNetCatHandler(alertSystem, netCatRepo, dataCenterRepo, projectRepo, netCatStatsRepo, taskPusher, agentHandler)
		//pageSpeedHandler := handlers.NewPageSpeedHandler(alertSystem, pageSpeedRepo, dataCenterRepo, projectRepo, pageSPeedStatsRepo, taskPusher, agentHandler)
		//pingHandler := handlers.NewPingHandler(alertSystem, pingRepo, dataCenterRepo, projectRepo, pingStatsRepo, taskPusher, agentHandler)
		//traceRouteHandler := handlers.NewTraceRouteHandler(alertSystem, traceRouteRepo, dataCenterRepo, projectRepo, traceRouteStatsRepo, taskPusher, agentHandler)

		mux := asynq.NewServeMux()
		// handlers
		mux.Handle(task_models.TypeEndpoint, tasks.NewEndpointTaskHandler(endpointHandler, zLogger))
		//mux.Handle(task_models.TypeNetCats, tasks.NewNetCatTaskHandler(netCatHandler, zLogger))
		//mux.Handle(task_models.TypePageSpeeds, tasks.NewPageSpeedTaskHandler(pageSpeedHandler, zLogger))
		//mux.Handle(task_models.TypePings, tasks.NewPingTaskHandler(pingHandler, zLogger))
		//mux.Handle(task_models.TypeTraceRoutes, tasks.NewTraceRouteTaskHandler(traceRouteHandler, zLogger))
		mux.Handle(task_models.TypeNotification, tasks.NewNotificationTaskHandler(alertSystem, projectRepo))

		if err := srv.Run(mux); err != nil {
			zLogger.Fatalf("cant start server: %s", err)
		}
	},
}
