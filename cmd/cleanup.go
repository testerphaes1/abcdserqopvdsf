package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"test-manager/config"
	"test-manager/repos"
	"test-manager/utils"
	"time"
)

func init() {
	rootCmd.AddCommand(cleaUpTasksCmd)
	cleaUpTasksCmd.Flags().Bool("vacuum", false, "Vacuum table")
}

var cleaUpTasksCmd = &cobra.Command{
	Use:   "clean-up",
	Short: "clean up",
	Long:  `database clean up`,
	Run: func(cmd *cobra.Command, args []string) {
		l := utils.ZapLogger()
		zLogger := l.Sugar()

		vacuum, err := cmd.Flags().GetBool("vacuum")

		_, config, err := config.ViperConfig()
		if err != nil {
			panic(err)
		}

		psqlDb, err := utils.PostgresConnection(config.Database.Pslq.Host, config.Database.Pslq.Port, config.Database.Pslq.User, config.Database.Pslq.Password, config.Database.Pslq.Database, config.Database.Pslq.Ssl, 5, 5)
		if err != nil {
			panic(err)
		}

		endpointStatsRepo := repos.NewEndpointStatsRepository(psqlDb)

		if vacuum {
			zLogger.Info("starting to vacuum endpoint stats table")
			err = endpointStatsRepo.VacuumEndpointStats(context.TODO())
			if err != nil {
				zLogger.Fatal(err)
			}
		} else {
			zLogger.Info("starting to clean up endpoint stats table")
			err = endpointStatsRepo.CleanUpResponseBodies(context.TODO(), time.Now().Add(-1*time.Hour))
			if err != nil {
				zLogger.Fatal(err)
			}
		}
		zLogger.Info("finished cleaning up")
	},
}
