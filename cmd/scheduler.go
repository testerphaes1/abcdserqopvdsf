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
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"test-manager/config"
	"test-manager/handlers"
	"test-manager/repos"
	"test-manager/tasks/push"
	"test-manager/utils"
	"time"
)

func init() {
	rootCmd.AddCommand(schedulerTasksCmd)
}

var schedulerTasksCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "async tasks scheduler",
	Long:  `async tasks scheduler`,
	Run: func(cmd *cobra.Command, args []string) {
		//go func() {
		//	http.ListenAndServe("localhost:8882", nil)
		//}()
		go func() {
			e := echo.New()

			p := prometheus.NewPrometheus("automated_test_scheduler", nil)
			p.Use(e)
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			// Start server
			go func() {
				if err := e.Start(":10001"); err != nil && err != http.ErrServerClosed {
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
		}()

		_, cfg, err := config.ViperConfig()
		if err != nil {
			panic(err)
		}

		redisClient, err := utils.CreateRedisConnection(context.TODO(), cfg.Database.Redis.Host, cfg.Database.Redis.Port, cfg.Database.Redis.Database, time.Duration(cfg.Database.Redis.Timeout)*time.Second)
		if err != nil {
			panic(err)
		}

		asynqClient := asynq.NewClient(asynq.RedisClientOpt{
			Addr:        redisClient.Options().Addr,
			DialTimeout: redisClient.Options().DialTimeout,
			Username:    redisClient.Options().Username,
			Password:    redisClient.Options().Password,
		})
		taskPusher := push.NewTaskPush(asynqClient)
		inspector := asynq.NewInspector(asynq.RedisClientOpt{
			Addr:        redisClient.Options().Addr,
			DialTimeout: redisClient.Options().DialTimeout,
			Username:    redisClient.Options().Username,
			Password:    redisClient.Options().Password,
		})

		psqlDb, err := utils.PostgresConnection(cfg.Database.Pslq.Host, cfg.Database.Pslq.Port, cfg.Database.Pslq.User, cfg.Database.Pslq.Password, cfg.Database.Pslq.Database, cfg.Database.Pslq.Ssl, 10, 10)
		if err != nil {
			panic(err)
		}

		endpointRepo := repos.NewEndpointRepository(psqlDb)
		heartBeatScheduler := handlers.NewHeartBeatScheduler(endpointRepo, taskPusher, inspector)
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		heartBeatScheduler.HeartBeatScheduling(ctx)

		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx.Done()
	},
}
