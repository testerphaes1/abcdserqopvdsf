package cmd

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test-manager/cache"
	"test-manager/config"
	"test-manager/gateway"
	"test-manager/handlers"
	"test-manager/repos"
	"test-manager/services/alert_system"
	"test-manager/tasks/push"
	"test-manager/utils"
	"time"
)

var Debug = true

func init() {
	rootCmd.AddCommand(httpCmd)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		//go func() {
		//	http.ListenAndServe("localhost:8884", nil)
		//}()
		e := echo.New()

		p := prometheus.NewPrometheus("automated_test_http", nil)
		p.Use(e)

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		}))
		e.HTTPErrorHandler = customZapHttpErrorHandler
		e.Validator = &CustomValidator{validator: validator.New()}

		// viper config
		_, config, err := config.ViperConfig()
		if err != nil {
			panic(err)
		}

		redisClient, err := utils.CreateRedisConnection(context.TODO(), config.Database.Redis.Host, config.Database.Redis.Port, config.Database.Redis.Database, time.Duration(config.Database.Redis.Timeout)*time.Second)
		if err != nil {
			panic(err)
		}

		redisCacheClient, err := utils.CreateRedisConnection(context.TODO(), config.Database.RedisCache.Host, config.Database.RedisCache.Port, config.Database.RedisCache.Database, time.Duration(config.Database.RedisCache.Timeout)*time.Second)
		if err != nil {
			panic(err)
		}

		cacheRepo := cache.NewRedisCache(redisCacheClient)

		psqlDb, err := utils.PostgresConnection(config.Database.Pslq.Host, config.Database.Pslq.Port, config.Database.Pslq.User, config.Database.Pslq.Password, config.Database.Pslq.Database, config.Database.Pslq.Ssl, 2, 2)
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

		boil.DebugMode = false
		accountRepo := repos.NewAccountsRepositoryRepository(psqlDb)
		projectRepo := repos.NewProjectsRepository(psqlDb)

		endpointRepo := repos.NewEndpointRepository(psqlDb)
		netCatRepo := repos.NewNetCatRepository(psqlDb)
		pageSpeedRepo := repos.NewPageSpeedRepository(psqlDb)
		pingRepo := repos.NewPingRepository(psqlDb)
		traceRouteRepo := repos.NewTraceRouteRepository(psqlDb)

		aggregateRepo := repos.NewAggregateRepository(psqlDb, endpointRepo, netCatRepo, pageSpeedRepo, pingRepo, traceRouteRepo)

		packageRepo := repos.NewPackagesRepository(psqlDb)
		dataCenterRepo := repos.NewDataCentersRepositoryRepository(cacheRepo, psqlDb)

		redisCache := cache.NewRedisCache(redisCacheClient)

		endpointStatsRepo := repos.NewEndpointStatsRepository(psqlDb)
		netCatStatsRepo := repos.NewNetCatStatsRepository(psqlDb)
		pingStatsRepo := repos.NewPingStatsRepository(psqlDb)
		traceRouteStatsRepo := repos.NewTraceRouteStatsRepository(psqlDb)
		pageSpeedStatsRepo := repos.NewPageSpeedStatsRepository(psqlDb)

		draftRepo := repos.NewDraftsRepository(psqlDb)

		gatewayRepo := repos.NewGatewaysRepository(psqlDb)
		orderRepository := repos.NewOrdersRepository(psqlDb)
		idpayGateway := gateway.NewIdpayGateway(config.Gateways.Idpay.BaseUrl, config.Gateways.Idpay.ApiToken)
		zarinpalGateway := gateway.NewZarinpalGateway(config.Gateways.Zarinpal.BaseUrl, config.Gateways.Zarinpal.ApiToken)

		alertSystem := alert_system.NewAlertHandler(config.Services.Alert.BaseUrl)

		faqRepo := repos.NewFaqsRepository(psqlDb)
		ticketRepo := repos.NewTicketsRepository(psqlDb)

		agentHandler := handlers.NewAgentHandler()
		endpointHandler := handlers.NewEndpointHandler(alertSystem, endpointRepo, dataCenterRepo,
			projectRepo, endpointStatsRepo, cacheRepo, taskPusher, agentHandler)
		netCatHandler := handlers.NewNetCatHandler(alertSystem, netCatRepo, dataCenterRepo, projectRepo, netCatStatsRepo, taskPusher, agentHandler)
		pageSpeedHandler := handlers.NewPageSpeedHandler(alertSystem, pageSpeedRepo, dataCenterRepo, projectRepo, pageSpeedStatsRepo, taskPusher, agentHandler)
		pingHandler := handlers.NewPingHandler(alertSystem, pingRepo, dataCenterRepo, projectRepo, pingStatsRepo, taskPusher, agentHandler)
		traceRouteHandler := handlers.NewTraceRouteHandler(alertSystem, traceRouteRepo, dataCenterRepo, projectRepo, traceRouteStatsRepo, taskPusher, agentHandler)

		//endpointHandler := handlers.NewEndpointHandler(endpointRepo, dataCenterRepo, taskPusher, agentHandler)
		ruleHandler := handlers.NewRulesHandler(projectRepo, endpointRepo, netCatRepo, pageSpeedRepo, pingRepo, traceRouteRepo, dataCenterRepo, taskPusher, agentHandler)
		controllers := handlers.NewHttpControllers(
			ruleHandler,
			endpointHandler,
			netCatHandler,
			pageSpeedHandler,
			pingHandler,
			traceRouteHandler,
			accountRepo,
			projectRepo,
			dataCenterRepo,
			aggregateRepo,
			packageRepo,
			endpointRepo,
			netCatRepo,
			pageSpeedRepo,
			traceRouteRepo,
			pingRepo,
			draftRepo,
			gatewayRepo,
			orderRepository,
			faqRepo,
			ticketRepo,
			idpayGateway,
			zarinpalGateway,
			alertSystem,
			redisCache,
			endpointStatsRepo,
			netCatStatsRepo,
			pageSpeedStatsRepo,
			pingStatsRepo,
			traceRouteStatsRepo)

		e.GET("/", controllers.Hello)
		e.POST("/rules/endpoint/register", controllers.RegisterEndpointRules, handlers.WithAuth())
		e.POST("/rules/netcat/register", controllers.RegisterNetCatRules, handlers.WithAuth())
		e.POST("/rules/ping/register", controllers.RegisterPingRules, handlers.WithAuth())
		e.POST("/rules/traceroute/register", controllers.RegisterTraceRouteRules, handlers.WithAuth())
		e.POST("/rules/pagespeed/register", controllers.RegisterPageSpeedRules, handlers.WithAuth())

		e.POST("/rules/endpoint/manual", controllers.ManualRunEndpointRules, handlers.WithAuth())
		e.POST("/rules/netcat/manual", controllers.ManualRunNetCatRules, handlers.WithAuth())
		e.POST("/rules/ping/manual", controllers.ManualRunPingRules, handlers.WithAuth())
		e.POST("/rules/traceroute/manual", controllers.ManualRunTraceRouteRules, handlers.WithAuth())
		e.POST("/rules/pagespeed/manual", controllers.ManualRunPageSpeedRules, handlers.WithAuth())

		e.GET("/rules/endpoint/:project_id/:id", controllers.GetEndpointRules, handlers.WithAuth())
		e.GET("/rules/netcat/:project_id/:id", controllers.GetNetCatRules, handlers.WithAuth())
		e.GET("/rules/ping/:project_id/:id", controllers.GetPingRules, handlers.WithAuth())
		e.GET("/rules/traceroute/:project_id/:id", controllers.GetTraceRouteRules, handlers.WithAuth())
		e.GET("/rules/pagespeed/:project_id/:id", controllers.GetPageSpeedRules, handlers.WithAuth())

		e.PUT("/rules/endpoint/:id", controllers.UpdateEndpointRules, handlers.WithAuth())
		e.PUT("/rules/netcat/:id", controllers.UpdateNetCatRules, handlers.WithAuth())
		e.PUT("/rules/ping/:id", controllers.UpdatePingRules, handlers.WithAuth())
		e.PUT("/rules/traceroute/:id", controllers.UpdateTraceRouteRules, handlers.WithAuth())
		e.PUT("/rules/pagespeed/:id", controllers.UpdatePageSpeedRules, handlers.WithAuth())

		e.DELETE("/rules/endpoint/:id", controllers.DeleteEndpointRules, handlers.WithAuth())
		e.DELETE("/rules/netcat/:id", controllers.DeleteNetCatRules, handlers.WithAuth())
		e.DELETE("/rules/ping/:id", controllers.DeletePingRules, handlers.WithAuth())
		e.DELETE("/rules/traceroute/:id", controllers.DeleteTraceRouteRules, handlers.WithAuth())
		e.DELETE("/rules/pagespeed/:id", controllers.DeletePageSpeedRules, handlers.WithAuth())

		e.GET("/rules/:project_id", controllers.GetRules, handlers.WithAuth())

		e.POST("/draft/rules/endpoint/register/:project_id", controllers.RegisterEndpointRulesDraft, handlers.WithAuth())
		e.POST("/draft/rules/netcat/register/:project_id", controllers.RegisterNetCatRulesDraft, handlers.WithAuth())
		e.POST("/draft/rules/ping/register/:project_id", controllers.RegisterPingRulesDraft, handlers.WithAuth())
		e.POST("/draft/rules/traceroute/register/:project_id", controllers.RegisterTraceRouteRulesDraft, handlers.WithAuth())
		e.POST("/draft/rules/pagespeed/register/:project_id", controllers.RegisterPageSpeedRulesDraft, handlers.WithAuth())
		e.GET("/draft/rules/:project_id", controllers.GetRulesDrafts, handlers.WithAuth())
		e.GET("/draft/rules/draft/:draft_id", controllers.GetRulesDraft, handlers.WithAuth())

		e.GET("/report/endpoint/datacenter/points_avg", controllers.ReportEndpointDataCenterPointsAndAvg, handlers.WithAuth())
		e.GET("/report/netcat/datacenter/points_avg", controllers.ReportNetCatDataCenterPointsAndAvg, handlers.WithAuth())
		e.GET("/report/ping/datacenter/points_avg", controllers.ReportPingDataCenterPointsAndAvg, handlers.WithAuth())
		e.GET("/report/traceroute/datacenter/points_avg", controllers.ReportTraceRouteDataCenterPointsAndAvg, handlers.WithAuth())
		e.GET("/report/pagespeed/datacenter/points_avg", controllers.ReportPageSpeedDataCenterPointsAndAvg, handlers.WithAuth())
		e.GET("/report/endpoint/details", controllers.ReportEndpointDetails, handlers.WithAuth())
		e.GET("/report/endpoint/quick", controllers.ReportEndpointQuickStats, handlers.WithAuth())
		e.GET("/report/netcat/quick", controllers.ReportNetCatQuickStats, handlers.WithAuth())
		e.GET("/report/ping/quick", controllers.ReportPingQuickStats, handlers.WithAuth())
		e.GET("/report/traceroute/quick", controllers.ReportTraceRouteQuickStats, handlers.WithAuth())
		e.GET("/report/pagespeed/quick", controllers.ReportPageSpeedQuickStats, handlers.WithAuth())

		e.POST("/email/verification", controllers.VerificationCode)
		//e.POST("/register", controllers.Register)
		e.POST("/auth", controllers.Auth)
		e.GET("/auth/info", controllers.AuthInfo)
		e.GET("/accounts", controllers.GetAccount, handlers.WithAuth())
		e.PUT("/accounts", controllers.UpdateAccount, handlers.WithAuth())
		e.POST("/accounts/password/reset/:account_id", controllers.ResetAccountPassword, handlers.WithAuth())

		e.POST("/projects", controllers.CreateProject, handlers.WithAuth())
		e.GET("/projects/:project_id", controllers.GetProject, handlers.WithAuth())
		e.PUT("/projects/:project_id", controllers.UpdateProject, handlers.WithAuth())

		e.POST("/packages", controllers.CreatePackage, handlers.WithAuth())
		e.GET("/packages/:package_id", controllers.GetPackage, handlers.WithAuth())
		e.PUT("/packages/:package_id", controllers.UpdatePackage, handlers.WithAuth())

		e.POST("/datacenters", controllers.CreateDatacenter, handlers.WithAuth())
		e.GET("/datacenters/:datacenter_id", controllers.GetDatacenter, handlers.WithAuth())
		e.PUT("/datacenters/:datacenter_id", controllers.UpdateDatacenter, handlers.WithAuth())

		e.POST("/gateway", controllers.CreateGateway, handlers.WithAuth())
		e.GET("/gateway/:gateway_id", controllers.GetGateways, handlers.WithAuth())
		e.PUT("/gateway/:gateway_id", controllers.UpdateGateway, handlers.WithAuth())

		e.POST("/order/create", controllers.CreateOrder, handlers.WithAuth())
		e.POST("/order/verify", controllers.VerifyOrder, handlers.WithAuth())
		e.GET("/order/:project_id/:order_id", controllers.GetOrderHistory, handlers.WithAuth())

		// TODO: delete alert (with delete all)
		e.GET("/alert/stats/:projectId", controllers.AlertStats, handlers.WithAuth())

		e.POST("/faq", controllers.CreateFaq, handlers.WithAuth())
		e.GET("/faq/:faq_id", controllers.GetFaq, handlers.WithAuth())
		e.PUT("/faq/:faq_id", controllers.UpdateFaq, handlers.WithAuth())

		e.POST("/ticket", controllers.CreateTicket, handlers.WithAuth())
		e.GET("/ticket/:ticket_id", controllers.GetTicket, handlers.WithAuth())
		e.PUT("/ticket/:ticket_id", controllers.UpdateTicket, handlers.WithAuth())

		// Start server
		go func() {
			if err := e.Start(":10000"); err != nil && err != http.ErrServerClosed {
				log.Fatal("shutting down server")
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	},
}

func customZapHttpErrorHandler(err error, c echo.Context) {
	var code = http.StatusInternalServerError
	var message interface{}

	if s, ok := err.(*utils.StandardHttpResponse); ok {
		code = s.Status
		message = s
		c.JSON(code, message)
		return
	} else {
		if c.Response().Committed {
			return
		}

		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}

		// Issue #1426
		code := he.Code
		message := he.Message
		if m, ok := he.Message.(string); ok {
			message = utils.StandardHttpResponse{
				Message: m,
				Status:  0,
				Data:    nil,
			}
		}

		// Send response
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			log.Println(err)
		}
		return
	}
}
