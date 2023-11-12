package cmd

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/mohammad-safakhou/asynqmon"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"test-manager/config"
	"test-manager/utils"
	"time"
)

func init() {
	rootCmd.AddCommand(uiTasksCmd)
}

var uiTasksCmd = &cobra.Command{
	Use:   "ui-tasks",
	Short: "ui consume async tasks events",
	Long:  `ui consume async tasks events`,
	Run: func(cmd *cobra.Command, args []string) {
		_, vConfig, err := config.ViperConfig()
		if err != nil {
			panic(err)
		}
		redisClient, err := utils.CreateRedisConnection(context.TODO(), vConfig.Database.Redis.Host, vConfig.Database.Redis.Port, vConfig.Database.Redis.Database, time.Duration(vConfig.Database.Redis.Timeout)*time.Second)
		if err != nil {
			panic(err)
		}
		defer redisClient.Close()

		h := asynqmon.New(asynqmon.Options{
			RootPath: "/monitoring", // RootPath specifies the root for asynqmon app
			RedisConnOpt: asynq.RedisClientOpt{
				Addr:        redisClient.Options().Addr,
				DialTimeout: redisClient.Options().DialTimeout,
				Username:    redisClient.Options().Username,
				Password:    redisClient.Options().Password,
			},
		})

		// Note: We need the tailing slash when using net/http.ServeMux.
		http.Handle(h.RootPath()+"/", h)

		// Go to http://localhost:8080/monitoring to see asynqmon homepage.
		log.Fatal(http.ListenAndServe(":9999", nil))
	},
}
