package cmd

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test-manager/config"
	"test-manager/repos"
	"test-manager/tasks"
	"test-manager/tasks/task_models"
	"test-manager/utils"
	"time"
)

func init() {
	rootCmd.AddCommand(consumeTasksAggregatorCmd)
}

var consumeTasksAggregatorCmd = &cobra.Command{
	Use:   "tasks-aggregator",
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

		//go func() {
		//	http.ListenAndServe("localhost:8883", nil)
		//}()

		go func() {
			e := echo.New()

			p := prometheus.NewPrometheus("automated_test_tasks", nil)
			p.Use(e)
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			// Start server
			go func() {
				if err := e.Start(":10003"); err != nil && err != http.ErrServerClosed {
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

		srv := asynq.NewServer(
			asynq.RedisClientOpt{
				Addr:        redisClient.Options().Addr,
				DialTimeout: redisClient.Options().DialTimeout,
				Username:    redisClient.Options().Username,
				Password:    redisClient.Options().Password,
			}, asynq.Config{
				Logger:           zLogger,
				Queues:           map[string]int{task_models.QueueEndpointStore: 6},
				GroupAggregator:  asynq.GroupAggregatorFunc(tasks.AggregateEndpointStats),
				GroupGracePeriod: 2 * time.Second,
				GroupMaxDelay:    10 * time.Second,
				GroupMaxSize:     10000,
			},
		)

		psqlDb, err := utils.PostgresConnection(config.Database.Pslq.Host, config.Database.Pslq.Port, config.Database.Pslq.User, config.Database.Pslq.Password, config.Database.Pslq.Database, config.Database.Pslq.Ssl, 2, 2)
		if err != nil {
			panic(err)
		}
		defer psqlDb.Close()
		endpointStatsRepo := repos.NewEndpointStatsRepository(psqlDb)

		mux := asynq.NewServeMux()
		mux.Handle(task_models.TypeEndpointStore, tasks.NewEndpointStoreHandler(endpointStatsRepo, zLogger))

		if err := srv.Run(mux); err != nil {
			zLogger.Fatalf("cant start server: %s", err)
		}
	},
}
