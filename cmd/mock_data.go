package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"test-manager/config"
	"test-manager/repos"
	"test-manager/utils"
	"time"
)

func init() {
	rootCmd.AddCommand(mockCmd)
}

var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		// viper config
		_, config, err := config.ViperConfig()
		if err != nil {
			panic(err)
		}

		//redisClient, err := utils.CreateRedisConnection(context.TODO(), config.Database.Redis.Host, config.Database.Redis.Port, time.Duration(config.Database.Redis.Timeout)*time.Second)
		//if err != nil {
		//	panic(err)
		//}

		psqlDb, err := utils.PostgresConnection(config.Database.Pslq.Host, config.Database.Pslq.Port, config.Database.Pslq.User, config.Database.Pslq.Password, config.Database.Pslq.Database, config.Database.Pslq.Ssl, 2, 2)
		if err != nil {
			panic(err)
		}
		//asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		//	Addr:        redisClient.Options().Addr,
		//	DialTimeout: redisClient.Options().DialTimeout,
		//	Username:    redisClient.Options().Username,
		//	Password:    redisClient.Options().Password,
		//})
		//taskPusher := push.NewTaskPush(asynqClient)
		//influxClient, writeAPI, queryAPI, err := utils.CreateInfluxDBConnection(context.TODO(), config.Database.Influx.Token, config.Database.Influx.Host+":"+config.Database.Influx.Port, config.Database.Influx.Org, config.Database.Influx.Bucket)
		//if err != nil {
		//	panic(err)
		//}
		//defer influxClient.Close()

		//accountRepo := repos.NewAccountsRepositoryRepository(psqlDb)
		//projectRepo := repos.NewProjectsRepository(psqlDb)
		//
		//endpointRepo := repos.NewEndpointRepository(psqlDb)
		//netCatRepo := repos.NewNetCatRepository(psqlDb)
		//pageSpeedRepo := repos.NewPageSpeedRepository(psqlDb)
		//pingRepo := repos.NewPingRepository(psqlDb)
		//traceRouteRepo := repos.NewTraceRouteRepository(psqlDb)
		//
		//aggregateRepo := repos.NewAggregateRepository(psqlDb, endpointRepo, netCatRepo, pageSpeedRepo, pingRepo, traceRouteRepo)
		//
		//packageRepo := repos.NewPackagesRepository(psqlDb)
		//dataCenterRepo := repos.NewDataCentersRepositoryRepository(psqlDb)
		//
		endpointStatsRepo := repos.NewEndpointStatsRepository(psqlDb)
		//netCatReportRepo := influx.NewNetCatsReportRepository(config.Database.Influx.Bucket, writeAPI, queryAPI, psqlDb)
		//pageSPeedReportRepo := influx.NewPageSpeedReportRepository(config.Database.Influx.Bucket, writeAPI, queryAPI, psqlDb)
		//pingReportRepo := influx.NewPingReportRepository(config.Database.Influx.Bucket, writeAPI, queryAPI, psqlDb)
		//traceRouteReportRepo := influx.NewTraceRouteReportRepository(config.Database.Influx.Bucket, writeAPI, queryAPI, psqlDb)
		//
		//draftRepo := repos.NewDraftsRepository(psqlDb)
		//
		//gatewayRepo := repos.NewGatewaysRepository(psqlDb)
		//orderRepository := repos.NewOrdersRepository(psqlDb)
		//idpayGateway := gateway.NewIdpayGateway(config.Gateways.Idpay.BaseUrl, config.Gateways.Idpay.ApiToken)
		//zarinpalGateway := gateway.NewZarinpalGateway(config.Gateways.Zarinpal.BaseUrl, config.Gateways.Zarinpal.ApiToken)
		//
		//alertSystem := alert_system.NewAlertHandler(config.Services.Alert.BaseUrl)
		//
		//faqRepo := repos.NewFaqsRepository(psqlDb)
		//ticketRepo := repos.NewTicketsRepository(psqlDb)
		//
		//agentHandler := handlers.NewAgentHandler()
		//endpointHandler := handlers.NewEndpointHandler(endpointRepo, dataCenterRepo, taskPusher, agentHandler)
		//ruleHandler := handlers.NewRulesHandler(projectRepo, endpointRepo, netCatRepo, pageSpeedRepo, pingRepo, traceRouteRepo, dataCenterRepo, taskPusher, agentHandler)

		go func() {
			for {
				status := 1
				if rand.Intn(10) < 3 {
					status = 0
				}
				err = endpointStatsRepo.Write(context.TODO(), repos.WriteEndpointStatsOptions{
					ProjectId:        4,
					EndpointName:     "just for test",
					EndpointId:       1,
					IsHeartBeat:      true,
					Url:              "www.digikala.com",
					DatacenterId:     1,
					Success:          status,
					ResponseTime:     float64(rand.Intn(100)),
					ResponseTimes:    "",
					ResponseBodies:   "",
					ResponseHeaders:  "",
					ResponseStatuses: "",
				})
				if err != nil {
					log.Println(err.Error())
				}
				time.Sleep(10000 * time.Millisecond)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		//go func() {
		//	for {
		//		status := 1
		//		if rand.Intn(10) < 3 {
		//			status = 0
		//		}
		//		err := pingReportRepo.WritePingReport(context.TODO(), 1, 1, status, float64(rand.Intn(100)))
		//		if err != nil {
		//			log.Println(err.Error())
		//		}
		//		time.Sleep(100 * time.Millisecond)
		//	}
		//}()
		//go func() {
		//	for {
		//		status := 1
		//		if rand.Intn(10) < 3 {
		//			status = 0
		//		}
		//		err := endpointReportRepo.WriteEndpointReport(context.TODO(), 1, 1, status, float64(rand.Intn(100)))
		//		if err != nil {
		//			log.Println(err.Error())
		//		}
		//		time.Sleep(100 * time.Millisecond)
		//	}
		//}()
		//go func() {
		//	for {
		//		status := 1
		//		if rand.Intn(10) < 3 {
		//			status = 0
		//		}
		//		err := endpointReportRepo.WriteEndpointReport(context.TODO(), 1, 1, status, float64(rand.Intn(100)))
		//		if err != nil {
		//			log.Println(err.Error())
		//		}
		//		time.Sleep(100 * time.Millisecond)
		//	}
		//}()
	},
}
