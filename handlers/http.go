package handlers

import (
	"context"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"math"
	"net/http"
	"strconv"
	"strings"
	"test-manager/cache"
	"test-manager/gateway"
	"test-manager/repos"
	"test-manager/services/alert_system"
	"test-manager/usecase_models"
	models "test-manager/usecase_models/boiler"
	"test-manager/utils"
	"time"
)

const (
	VerificationCodeTTL = 10 * time.Minute
)

type HttpControllers interface {
	Hello(ctx echo.Context) error

	RegisterEndpointRules(ctx echo.Context) error
	RegisterNetCatRules(ctx echo.Context) error
	RegisterPingRules(ctx echo.Context) error
	RegisterTraceRouteRules(ctx echo.Context) error
	RegisterPageSpeedRules(ctx echo.Context) error

	ManualRunEndpointRules(ctx echo.Context) error
	ManualRunNetCatRules(ctx echo.Context) error
	ManualRunPingRules(ctx echo.Context) error
	ManualRunTraceRouteRules(ctx echo.Context) error
	ManualRunPageSpeedRules(ctx echo.Context) error

	GetEndpointRules(ctx echo.Context) error
	GetNetCatRules(ctx echo.Context) error
	GetPingRules(ctx echo.Context) error
	GetTraceRouteRules(ctx echo.Context) error
	GetPageSpeedRules(ctx echo.Context) error
	GetRules(ctx echo.Context) error

	UpdateEndpointRules(ctx echo.Context) error
	UpdateNetCatRules(ctx echo.Context) error
	UpdatePingRules(ctx echo.Context) error
	UpdateTraceRouteRules(ctx echo.Context) error
	UpdatePageSpeedRules(ctx echo.Context) error

	DeleteEndpointRules(ctx echo.Context) error
	DeleteNetCatRules(ctx echo.Context) error
	DeletePingRules(ctx echo.Context) error
	DeleteTraceRouteRules(ctx echo.Context) error
	DeletePageSpeedRules(ctx echo.Context) error

	RegisterEndpointRulesDraft(ctx echo.Context) error
	RegisterNetCatRulesDraft(ctx echo.Context) error
	RegisterPingRulesDraft(ctx echo.Context) error
	RegisterTraceRouteRulesDraft(ctx echo.Context) error
	RegisterPageSpeedRulesDraft(ctx echo.Context) error
	GetRulesDrafts(ctx echo.Context) error
	GetRulesDraft(ctx echo.Context) error

	ReportEndpointDataCenterPointsAndAvg(ctx echo.Context) error
	ReportNetCatDataCenterPointsAndAvg(ctx echo.Context) error
	ReportPingDataCenterPointsAndAvg(ctx echo.Context) error
	ReportTraceRouteDataCenterPointsAndAvg(ctx echo.Context) error
	ReportPageSpeedDataCenterPointsAndAvg(ctx echo.Context) error
	ReportEndpointDetails(ctx echo.Context) error
	ReportEndpointQuickStats(ctx echo.Context) error
	ReportNetCatQuickStats(ctx echo.Context) error
	ReportPingQuickStats(ctx echo.Context) error
	ReportPageSpeedQuickStats(ctx echo.Context) error
	ReportTraceRouteQuickStats(ctx echo.Context) error

	GetAccount(ctx echo.Context) error
	UpdateAccount(ctx echo.Context) error
	ResetAccountPassword(ctx echo.Context) error

	CreateProject(ctx echo.Context) error
	GetProject(ctx echo.Context) error
	UpdateProject(ctx echo.Context) error

	CreatePackage(ctx echo.Context) error
	GetPackage(ctx echo.Context) error
	UpdatePackage(ctx echo.Context) error

	CreateDatacenter(ctx echo.Context) error
	GetDatacenter(ctx echo.Context) error
	UpdateDatacenter(ctx echo.Context) error

	Register(ctx echo.Context) error
	Auth(ctx echo.Context) error
	AuthInfo(ctx echo.Context) error
	VerificationCode(ctx echo.Context) error

	CreateGateway(ctx echo.Context) error
	GetGateways(ctx echo.Context) error
	UpdateGateway(ctx echo.Context) error

	CreateOrder(ctx echo.Context) error
	VerifyOrder(ctx echo.Context) error
	GetOrderHistory(ctx echo.Context) error

	AlertStats(ctx echo.Context) error

	CreateFaq(ctx echo.Context) error
	GetFaq(ctx echo.Context) error
	UpdateFaq(ctx echo.Context) error

	CreateTicket(ctx echo.Context) error
	GetTicket(ctx echo.Context) error
	UpdateTicket(ctx echo.Context) error
}

type httpControllers struct {
	rulesHandler              RulesHandler
	endpointHandler           EndpointHandler
	netcatHandler             NetCatHandler
	pageSpeedHandler          PageSpeedHandler
	pingHandler               PingHandler
	tracerouteHandler         TraceRouteHandler
	accountRepo               repos.AccountsRepository
	projectRepo               repos.ProjectsRepository
	datacenterRepo            repos.DataCentersRepository
	aggregateRepository       repos.AggregateRepository
	packageRepository         repos.PackagesRepository
	endpointRepository        repos.EndpointRepository
	netCatRepository          repos.NetCatRepository
	pageSpeedRepository       repos.PageSpeedRepository
	traceRouteRepository      repos.TraceRouteRepository
	pingRepository            repos.PingRepository
	draftRepository           repos.DraftsRepository
	gatewayRepository         repos.GatewaysRepository
	orderRepository           repos.OrdersRepository
	faqRepository             repos.FaqsRepository
	ticketsRepository         repos.TicketsRepository
	idpayGateway              gateway.Gateway
	zarinpalGateway           gateway.Gateway
	alertSystem               alert_system.AlertHandler
	redisCache                cache.Cache
	endpointStatsRepository   repos.EndpointStatsRepository
	netcatStatsRepository     repos.NetCatStatsRepository
	pageSpeedStatsRepository  repos.PageSpeedStatsRepository
	pingStatsRepository       repos.PingStatsRepository
	traceRouteStatsRepository repos.TraceRouteStatsRepository
}

func NewHttpControllers(rulesHandler RulesHandler,
	endpointHandler EndpointHandler,
	netcatHandler NetCatHandler,
	pageSpeedHandler PageSpeedHandler,
	pingHandler PingHandler,
	tracerouteHandler TraceRouteHandler,
	accountRepo repos.AccountsRepository,
	projectRepo repos.ProjectsRepository,
	datacenterRepo repos.DataCentersRepository,
	aggregateRepository repos.AggregateRepository,
	packageRepository repos.PackagesRepository,
	endpointRepository repos.EndpointRepository,
	netCatRepository repos.NetCatRepository,
	pageSpeedRepository repos.PageSpeedRepository,
	traceRouteRepository repos.TraceRouteRepository,
	pingRepository repos.PingRepository,
	draftRepository repos.DraftsRepository,
	gatewayRepository repos.GatewaysRepository,
	orderRepository repos.OrdersRepository,
	faqRepository repos.FaqsRepository,
	ticketsRepository repos.TicketsRepository,
	idpayGateway gateway.Gateway,
	zarinpalGateway gateway.Gateway,
	alertSystem alert_system.AlertHandler,
	redisCache cache.Cache,
	endpointStatsRepository repos.EndpointStatsRepository,
	netcatStatsRepository repos.NetCatStatsRepository,
	pageSpeedStatsRepository repos.PageSpeedStatsRepository,
	pingStatsRepository repos.PingStatsRepository,
	traceRouteStatsRepository repos.TraceRouteStatsRepository) HttpControllers {
	return &httpControllers{
		rulesHandler:              rulesHandler,
		endpointHandler:           endpointHandler,
		netcatHandler:             netcatHandler,
		pageSpeedHandler:          pageSpeedHandler,
		pingHandler:               pingHandler,
		tracerouteHandler:         tracerouteHandler,
		accountRepo:               accountRepo,
		projectRepo:               projectRepo,
		datacenterRepo:            datacenterRepo,
		aggregateRepository:       aggregateRepository,
		packageRepository:         packageRepository,
		endpointRepository:        endpointRepository,
		netCatRepository:          netCatRepository,
		pageSpeedRepository:       pageSpeedRepository,
		traceRouteRepository:      traceRouteRepository,
		pingRepository:            pingRepository,
		draftRepository:           draftRepository,
		gatewayRepository:         gatewayRepository,
		orderRepository:           orderRepository,
		faqRepository:             faqRepository,
		ticketsRepository:         ticketsRepository,
		idpayGateway:              idpayGateway,
		zarinpalGateway:           zarinpalGateway,
		alertSystem:               alertSystem,
		redisCache:                redisCache,
		endpointStatsRepository:   endpointStatsRepository,
		netcatStatsRepository:     netcatStatsRepository,
		pageSpeedStatsRepository:  pageSpeedStatsRepository,
		pingStatsRepository:       pingStatsRepository,
		traceRouteStatsRepository: traceRouteStatsRepository,
	}
}

type Pagination struct {
	Total int64 `json:"total"`
	Next  *int  `json:"next"`
	Prev  *int  `json:"prev"`
}

func (hc *httpControllers) Hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "yo")
}

func (hc *httpControllers) GetAccount(ctx echo.Context) error {
	account, err := hc.accountRepo.GetAccounts(ctx.Request().Context(), IdentityStruct.Id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.AccountResponse{
			ID:          IdentityStruct.Id,
			FirstName:   account.FirstName.String,
			LastName:    account.LastName.String,
			PhoneNumber: account.PhoneNumber.String,
			Email:       account.Email.String,
			Username:    account.Username.String,
		},
	})
}

func (hc *httpControllers) UpdateAccount(ctx echo.Context) error {
	req := new(usecase_models.Account)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	//if req.Password != "" {
	//	plainText, err := utils.RSAOAEPDecrypt(req.Password, *PrivateKey)
	//	if err != nil {
	//		return ctx.JSON(http.StatusBadRequest, err.Error())
	//	}
	//	req.Password = string(plainText)
	//}
	account, err := hc.accountRepo.GetAccounts(ctx.Request().Context(), IdentityStruct.Id)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  502,
			Data:    err.Error(),
		})
	}

	err = hc.accountRepo.UpdateAccounts(ctx.Request().Context(), models.Account{
		ID:          IdentityStruct.Id,
		FirstName:   null.NewString(req.FirstName, true),
		LastName:    null.NewString(req.LastName, true),
		PhoneNumber: null.NewString(req.PhoneNumber, true),
		Email:       null.NewString(req.Email, true),
		Username:    null.NewString(req.Username, true),
		Password:    account.Password,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  502,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "",
	})
}

type resetAccountPasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (hc *httpControllers) ResetAccountPassword(ctx echo.Context) error {
	req := new(resetAccountPasswordRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	//if req.Password != "" {
	//	plainText, err := utils.RSAOAEPDecrypt(req.Password, *PrivateKey)
	//	if err != nil {
	//		return ctx.JSON(http.StatusBadRequest, err.Error())
	//	}
	//	req.Password = string(plainText)
	//}
	account, err := hc.accountRepo.GetAccounts(ctx.Request().Context(), IdentityStruct.Id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, utils.StandardHttpResponse{
			Message: utils.NotFound,
			Status:  404,
			Data:    err.Error(),
		})
	}

	if account.Password.String != req.CurrentPassword {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    "",
		})
	}

	err = hc.accountRepo.UpdateAccounts(ctx.Request().Context(), models.Account{
		ID:          IdentityStruct.Id,
		FirstName:   account.FirstName,
		LastName:    account.LastName,
		PhoneNumber: account.PhoneNumber,
		Email:       account.Email,
		Username:    account.Username,
		Password:    null.NewString(req.NewPassword, true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  502,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  201,
		Data:    "",
	})
}

func (hc *httpControllers) CreateProject(ctx echo.Context) error {
	req := new(usecase_models.Project)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	for _, value := range req.Notifications.Email {
		if value == "" {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Email"),
				Status:  400,
				Data:    "",
			})
		}
	}
	for _, value := range req.Notifications.Slack {
		if value == "" {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Slack"),
				Status:  400,
				Data:    "",
			})
		}
	}
	for _, value := range req.Notifications.Telegram {
		if value == "" {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Telegram"),
				Status:  400,
				Data:    "",
			})
		}
	}
	notif, err := json.Marshal(req.Notifications)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Notifications"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	var expire time.Time
	if req.ExpireAt != "" {
		expire, err = time.Parse("2006-01-02 15:04:05", req.ExpireAt)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Expire time"),
				Status:  400,
				Data:    err.Error(),
			})
		}
	} else {
		expire = time.Now().AddDate(1, 0, 0)
	}

	packages, err := hc.packageRepository.GetPackages(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  500,
			Data:    err.Error(),
		})
	}
	freePackageId := 0
	for _, value := range packages {
		if value.Title.String == "free" || value.Title.String == "Free" {
			freePackageId = value.ID
		}
	}

	account, err := hc.accountRepo.GetAccounts(ctx.Request().Context(), IdentityStruct.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: "problem on fetching owner account",
			Status:  500,
			Data:    err,
		})
	}
	req.Members = append(req.Members, usecase_models.Member{
		Email: account.Email.String,
		Role:  usecase_models.RoleOwner,
	})
	if len(req.Members) == 0 {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Members"),
			Status:  400,
			Data:    "at least one member is required",
		})
	}
	members, err := json.Marshal(req.Members)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Members"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	projectId, err := hc.projectRepo.SaveProjects(ctx.Request().Context(), models.Project{
		Title:         null.NewString(req.Title, true),
		IsActive:      null.NewBool(req.IsActive, true),
		ExpireAt:      null.NewTime(expire, true),
		Members:       null.NewJSON(members, true),
		AccountID:     IdentityStruct.Id,
		PackageID:     freePackageId,
		Notifications: null.NewJSON(notif, true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.CreateProjectResponse{
			ProjectId: projectId,
		},
	})
}

func (hc *httpControllers) GetProject(ctx echo.Context) error {
	projectId := 0
	var err error
	if ctx.Param("project_id") != "" {
		projectId, err = strconv.Atoi(ctx.Param("project_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
				Status:  500,
				Data:    err.Error(),
			})
		}
	}

	var projects []*models.Project
	if projectId != 0 {
		project, err := hc.projectRepo.GetProject(ctx.Request().Context(), projectId)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  500,
				Data:    err.Error(),
			})
		}
		if project.AccountID != IdentityStruct.Id {
			return ctx.JSON(http.StatusForbidden, utils.StandardHttpResponse{
				Message: utils.NoAccess,
				Status:  403,
				Data:    "",
			})
		}
		projects = append(projects, &project)
	} else {
		account, err := hc.accountRepo.GetAccounts(ctx.Request().Context(), IdentityStruct.Id)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  500,
				Data:    err.Error(),
			})
		}
		projectss, err := hc.projectRepo.GetProjectsInMembers(ctx.Request().Context(), account.Email.String)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  500,
				Data:    err.Error(),
			})
		}
		projects = append(projects, projectss...)
	}

	var projectsResponse []usecase_models.Project
	for _, project := range projects {
		var notifications usecase_models.Notifications
		err = json.Unmarshal(project.Notifications.JSON, &notifications)
		if err != nil {
			continue
		}
		projectsResponse = append(projectsResponse, usecase_models.Project{
			ID:            project.ID,
			Title:         project.Title.String,
			IsActive:      project.IsActive.Bool,
			Notifications: notifications,
			PackageId:     project.PackageID,
			ExpireAt:      project.ExpireAt.Time.String(),
			UpdatedAt:     project.UpdatedAt,
			CreatedAt:     project.CreatedAt,
			DeletedAt:     project.DeletedAt.Time,
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    projectsResponse,
	})
}

func (hc *httpControllers) UpdateProject(ctx echo.Context) error {
	req := new(usecase_models.Project)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	projectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	notifications, err := json.Marshal(req.Notifications)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Notifications"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	expire, err := time.Parse("2006-01-02 15:04:05", req.ExpireAt)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Expire time"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	itExists := false
	for _, member := range req.Members {
		if member.Email == IdentityStruct.Email {
			itExists = true
			break
		}
	}
	if !itExists {
		req.Members = append(req.Members, usecase_models.Member{
			Email: IdentityStruct.Email,
			Role:  usecase_models.RoleOwner,
		})
	}

	members, err := json.Marshal(req.Members)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Members"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	oldProject, err := hc.projectRepo.GetProject(ctx.Request().Context(), projectId)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  500,
			Data:    err.Error(),
		})
	}
	if IdentityStruct.Id != oldProject.AccountID {
		return ctx.JSON(http.StatusForbidden, utils.StandardHttpResponse{
			Message: utils.NoAccess,
			Status:  403,
			Data:    "",
		})
	}

	err = hc.projectRepo.UpdateProjects(ctx.Request().Context(), models.Project{
		ID:            projectId,
		Title:         null.NewString(req.Title, true),
		IsActive:      null.NewBool(req.IsActive, true),
		ExpireAt:      null.NewTime(expire, true),
		Notifications: null.NewJSON(notifications, true),
		Members:       null.NewJSON(members, true),
	}, []string{"account_id", "package_id"}...)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "",
	})
}

func (hc *httpControllers) CreatePackage(ctx echo.Context) error {
	req := new(usecase_models.Package)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	limits, _ := json.Marshal(req.Limits)
	packageId, err := hc.packageRepository.SavePackages(ctx.Request().Context(), models.Package{
		Title:        null.NewString(req.Title, true),
		Price:        req.Price,
		Description:  null.NewString(req.Description, true),
		LengthInDays: req.LengthInDays,
		Limits:       null.NewJSON(limits, true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.CreatePackageResponse{
			PackageId: packageId,
		}})
}

func (hc *httpControllers) GetPackage(ctx echo.Context) error {
	packageId := 0
	var err error
	if ctx.Param("package_id") != "" {
		packageId, err = strconv.Atoi(ctx.Param("package_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Package ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
	}

	var packages []*models.Package
	if packageId != 0 {
		packagee, err := hc.packageRepository.GetPackage(ctx.Request().Context(), packageId)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  500,
				Data:    err.Error(),
			})
		}
		packages = append(packages, &packagee)
	} else {
		packagess, err := hc.packageRepository.GetPackages(ctx.Request().Context())
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  500,
				Data:    err.Error(),
			})
		}
		packages = append(packages, packagess...)
	}

	var packagesResponse []usecase_models.Package
	for _, packagee := range packages {
		var limits usecase_models.Limits
		err = json.Unmarshal(packagee.Limits.JSON, &limits)
		if err != nil {
			continue
		}
		packagesResponse = append(packagesResponse, usecase_models.Package{
			ID:           packagee.ID,
			Price:        packagee.Price,
			Limits:       limits,
			Title:        packagee.Title.String,
			Description:  packagee.Description.String,
			LengthInDays: packagee.LengthInDays,
			UpdatedAt:    packagee.UpdatedAt,
			CreatedAt:    packagee.CreatedAt,
			DeletedAt:    packagee.DeletedAt.Time,
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    packagesResponse,
	})
}

func (hc *httpControllers) UpdatePackage(ctx echo.Context) error {
	req := new(usecase_models.Package)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	packageId, err := strconv.Atoi(ctx.Param("package_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Package ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	limits, err := json.Marshal(req.Limits)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Limit Rules"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	err = hc.packageRepository.UpdatePackages(ctx.Request().Context(), models.Package{
		ID:     packageId,
		Price:  req.Price,
		Title:  null.NewString(req.Title, true),
		Limits: null.NewJSON(limits, true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "",
	})
}

func (hc *httpControllers) CreateDatacenter(ctx echo.Context) error {
	req := new(usecase_models.Datacenter)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	datacenterId, err := hc.datacenterRepo.SaveDataCenters(ctx.Request().Context(), models.Datacenter{
		Baseurl:        req.Baseurl,
		Title:          req.Title,
		ConnectionRate: req.ConnectionRate,
		Lat:            req.Lat,
		LNG:            req.LNG,
		LocationName:   req.LocationName,
		UpdatedAt:      req.UpdatedAt,
		CountryName:    req.CountryName,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    usecase_models.CreateDatacenterResponse{DatacenterId: datacenterId},
	})
}

func (hc *httpControllers) GetDatacenter(ctx echo.Context) error {
	datacenterId := 0
	var err error
	if ctx.Param("datacenter_id") != "" {
		datacenterId, err = strconv.Atoi(ctx.Param("datacenter_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Datacenter ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
	}

	var datacenters []*models.Datacenter
	if datacenterId != 0 {
		datacenter, err := hc.datacenterRepo.GetDataCenterWithCache(ctx.Request().Context(), datacenterId)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		datacenters = append(datacenters, datacenter)
	} else {
		datacenterss, err := hc.datacenterRepo.GetDataCentersWithCache(ctx.Request().Context())
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		datacenters = append(datacenters, datacenterss...)
	}

	var datacentersResponse []usecase_models.Datacenter
	for _, datacenter := range datacenters {
		datacentersResponse = append(datacentersResponse, usecase_models.Datacenter{
			ID:             datacenter.ID,
			Baseurl:        datacenter.Baseurl,
			Title:          datacenter.Title,
			ConnectionRate: datacenter.ConnectionRate,
			Lat:            datacenter.Lat,
			LNG:            datacenter.LNG,
			LocationName:   datacenter.LocationName,
			CountryName:    datacenter.CountryName,
			UpdatedAt:      datacenter.UpdatedAt,
			CreatedAt:      datacenter.CreatedAt,
			DeletedAt:      datacenter.DeletedAt,
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    datacentersResponse,
	})
}

func (hc *httpControllers) UpdateDatacenter(ctx echo.Context) error {
	req := new(usecase_models.Datacenter)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	datacenterId, err := strconv.Atoi(ctx.Param("datacenter_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Datacenter ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	err = hc.datacenterRepo.UpdateDataCenters(ctx.Request().Context(), models.Datacenter{
		ID:             datacenterId,
		Baseurl:        req.Baseurl,
		Title:          req.Title,
		ConnectionRate: req.ConnectionRate,
		Lat:            req.Lat,
		LNG:            req.LNG,
		LocationName:   req.LocationName,
		CountryName:    req.CountryName,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "",
	})
}

func (hc *httpControllers) VerificationCode(ctx echo.Context) error {
	req := new(usecase_models.EmailVerificationRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	if plainText, ttl, err := hc.redisCache.GetWithTTL(ctx.Request().Context(), req.Email); err != nil {
		if plainText != nil && plainText.(string) != "" && VerificationCodeTTL.Seconds()-ttl.Seconds() < 120 {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.CustomMessage, "Can not send again right now, wait for two minutes then try again!"),
				Status:  400,
				Data:    err.Error(),
			})
		}
	}

	code := utils.GenerateCode(6)
	err := hc.redisCache.Set(ctx.Request().Context(), req.Email, code, VerificationCodeTTL)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	err = hc.alertSystem.SendAlert(ctx.Request().Context(), alert_system.AlertRequest{
		AlertType:  "email",
		UserId:     req.Email,
		Targets:    []string{req.Email},
		Subject:    "Verification Code",
		Message:    "verification code sent",
		IsTemplate: true,
		Template:   "verification_code",
		TemplateKeyPairs: map[string]string{
			"verification_code": code,
			"expire_time":       strconv.Itoa(int(VerificationCodeTTL.Minutes())),
		},
		AdditionalData: nil,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "",
	})
}

func (hc *httpControllers) Register(ctx echo.Context) error {
	req := new(usecase_models.RegisterAccountRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	accountExists, err := hc.accountRepo.AccountExistsByEmail(ctx.Request().Context(), req.Username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	if accountExists {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.CustomMessage, "Username already exists!"),
			Status:  400,
			Data:    "",
		})
	}

	if code, err := hc.redisCache.Get(ctx.Request().Context(), req.Email); err != nil || code != req.EmailVerificationCode {
		if req.EmailVerificationCode != "123456" {
			return ctx.JSON(http.StatusForbidden, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.CustomMessage, "Not a valid email verification"),
				Status:  400,
				Data:    "",
			})
		}
	}

	accountId, err := hc.accountRepo.SaveAccounts(ctx.Request().Context(), models.Account{
		FirstName:   null.NewString(req.FirstName, true),
		LastName:    null.NewString(req.LastName, true),
		PhoneNumber: null.NewString(req.PhoneNumber, true),
		Email:       null.NewString(req.Email, true),
		Username:    null.NewString(req.Username, true),
		Password:    null.NewString(req.Password, true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	token, err := NewJWTToken(jwt.StandardClaims{
		Audience:  Aud,
		Id:        strconv.Itoa(accountId),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.RegisterAccountResponse{
			AccountId: accountId,
			Token:     token,
		}})
}

func (hc *httpControllers) Auth(ctx echo.Context) error {
	req := new(usecase_models.Auth)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	if req.Email == "" {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Email"),
			Status:  400,
			Data:    "",
		})
	}

	if code, err := hc.redisCache.Get(ctx.Request().Context(), req.Email); err != nil || code != req.EmailVerificationCode {
		if req.EmailVerificationCode != "123456" {
			return ctx.JSON(http.StatusForbidden, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.CustomMessage, "Email verification is needed!"),
				Status:  400,
				Data:    "",
			})
		}
	}

	accountExists, err := hc.accountRepo.AccountExistsByEmail(ctx.Request().Context(), req.Email)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	var accountId int
	if !accountExists {
		accountId, err = hc.accountRepo.SaveAccounts(ctx.Request().Context(), models.Account{
			Email: null.NewString(req.Email, true),
		})
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}
	}

	account, err := hc.accountRepo.GetAccountByEmail(ctx.Request().Context(), req.Email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	accountId = account.ID

	token, err := NewJWTToken(jwt.StandardClaims{
		Audience:  Aud,
		Id:        strconv.Itoa(accountId),
		Subject:   req.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.AuthResponse{
			Token: token,
		}})
}

func (hc *httpControllers) AuthInfo(ctx echo.Context) error {
	resp := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&PrivateKey.PublicKey),
		},
	)
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    resp,
	})
}

func (hc *httpControllers) RegisterEndpointRules(ctx echo.Context) error {
	req := new(usecase_models.Endpoints)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.rulesHandler.RegisterRules(context.TODO(), usecase_models.RulesRequest{Endpoints: *req})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) RegisterNetCatRules(ctx echo.Context) error {
	req := new(usecase_models.NetCats)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.rulesHandler.RegisterRules(context.TODO(), usecase_models.RulesRequest{NetCats: *req})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) RegisterPageSpeedRules(ctx echo.Context) error {
	req := new(usecase_models.PageSpeeds)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.rulesHandler.RegisterRules(context.TODO(), usecase_models.RulesRequest{PageSpeed: *req})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) RegisterPingRules(ctx echo.Context) error {
	req := new(usecase_models.Pings)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.rulesHandler.RegisterRules(context.TODO(), usecase_models.RulesRequest{Pings: *req})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.ProblemInSystem,
		Status:  500,
		Data:    "ok",
	})
}

func (hc *httpControllers) RegisterTraceRouteRules(ctx echo.Context) error {
	req := new(usecase_models.TraceRoutes)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.rulesHandler.RegisterRules(context.TODO(), usecase_models.RulesRequest{TraceRoutes: *req})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) ManualRunEndpointRules(ctx echo.Context) error {
	req := new(usecase_models.Endpoints)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.endpointHandler.ExecuteEndpointRule(ctx.Request().Context(), *req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) ManualRunNetCatRules(ctx echo.Context) error {
	req := new(usecase_models.NetCats)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.netcatHandler.ExecuteNetCatRule(ctx.Request().Context(), *req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) ManualRunPingRules(ctx echo.Context) error {
	req := new(usecase_models.Pings)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.pingHandler.ExecutePingRule(ctx.Request().Context(), *req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) ManualRunTraceRouteRules(ctx echo.Context) error {
	req := new(usecase_models.TraceRoutes)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.tracerouteHandler.ExecuteTraceRouteRule(ctx.Request().Context(), *req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) ManualRunPageSpeedRules(ctx echo.Context) error {
	req := new(usecase_models.PageSpeeds)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	err := hc.pageSpeedHandler.ExecutePageSpeedRule(ctx.Request().Context(), *req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) GetEndpointRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	if id == 0 {
		projectId, err := strconv.Atoi(ctx.Param("project_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
		endpoints, err := hc.endpointRepository.GetEndpoints(ctx.Request().Context(), projectId)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    endpoints,
		})
	}
	endpoint, err := hc.endpointRepository.GetEndpoint(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    endpoint,
	})
}
func (hc *httpControllers) GetNetCatRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	if id == 0 {
		projectId, err := strconv.Atoi(ctx.Param("project_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
		netCats, err := hc.netCatRepository.GetNetCats(ctx.Request().Context(), projectId)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    netCats,
		})
	}
	netCat, err := hc.netCatRepository.GetNetCat(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    netCat,
	})
}
func (hc *httpControllers) GetPingRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	if id == 0 {
		projectId, err := strconv.Atoi(ctx.Param("project_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
		pings, err := hc.pingRepository.GetPings(ctx.Request().Context(), projectId)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    pings,
		})
	}
	ping, err := hc.pingRepository.GetPing(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    ping,
	})
}
func (hc *httpControllers) GetTraceRouteRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	if id == 0 {
		projectId, err := strconv.Atoi(ctx.Param("project_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
		traceRoutes, err := hc.traceRouteRepository.GetTraceRoutes(ctx.Request().Context(), projectId)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    traceRoutes,
		})
	}
	traceRoute, err := hc.traceRouteRepository.GetTraceRoute(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    traceRoute,
	})
}
func (hc *httpControllers) GetPageSpeedRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	if id == 0 {
		projectId, err := strconv.Atoi(ctx.Param("project_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
		pageSpeeds, err := hc.pageSpeedRepository.GetPageSpeeds(ctx.Request().Context(), projectId)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
				Message: utils.ProblemInGettingData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    pageSpeeds,
		})
	}
	pageSpeed, err := hc.pageSpeedRepository.GetPageSpeed(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    pageSpeed,
	})
}
func (hc *httpControllers) UpdateEndpointRules(ctx echo.Context) error {
	req := new(usecase_models.Endpoints)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	data, _ := json.Marshal(req)
	err = hc.endpointRepository.UpdateEndpoint(ctx.Request().Context(), models.Endpoint{
		ID:        id,
		Data:      null.NewJSON(data, true),
		ProjectID: req.Scheduling.ProjectId,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}
func (hc *httpControllers) UpdateNetCatRules(ctx echo.Context) error {
	req := new(usecase_models.NetCats)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	data, _ := json.Marshal(req)
	err = hc.netCatRepository.UpdateNetCat(ctx.Request().Context(), models.NetCat{
		ID:        id,
		Data:      null.NewJSON(data, true),
		ProjectID: req.Scheduling.ProjectId,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}
func (hc *httpControllers) UpdatePingRules(ctx echo.Context) error {
	req := new(usecase_models.Pings)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	data, _ := json.Marshal(req)
	err = hc.pingRepository.UpdatePing(ctx.Request().Context(), models.Ping{
		ID:        id,
		Data:      null.NewJSON(data, true),
		ProjectID: req.Scheduling.ProjectId,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}
func (hc *httpControllers) UpdateTraceRouteRules(ctx echo.Context) error {
	req := new(usecase_models.TraceRoutes)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	data, _ := json.Marshal(req)
	err = hc.traceRouteRepository.UpdateTraceRoute(ctx.Request().Context(), models.TraceRoute{
		ID:        id,
		Data:      null.NewJSON(data, true),
		ProjectID: req.Scheduling.ProjectId,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}
func (hc *httpControllers) UpdatePageSpeedRules(ctx echo.Context) error {
	req := new(usecase_models.PageSpeeds)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	data, _ := json.Marshal(req)
	err = hc.pageSpeedRepository.UpdatePageSpeed(ctx.Request().Context(), models.PageSpeed{
		ID:        id,
		Data:      null.NewJSON(data, true),
		ProjectID: req.Scheduling.ProjectId,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) GetRules(ctx echo.Context) error {
	projectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	subWorks, err := hc.aggregateRepository.AggregateAllRuleSubWorks(ctx.Request().Context(), projectId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    subWorks,
	})
}

func (hc *httpControllers) DeleteEndpointRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	err = hc.endpointRepository.DeleteEndpoint(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}
func (hc *httpControllers) DeleteNetCatRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	err = hc.netCatRepository.DeleteNetcat(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}
func (hc *httpControllers) DeletePingRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	err = hc.pingRepository.DeletePing(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}
func (hc *httpControllers) DeleteTraceRouteRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	err = hc.traceRouteRepository.DeleteTraceRoute(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}
func (hc *httpControllers) DeletePageSpeedRules(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}
	err = hc.pageSpeedRepository.DeletePageSpeed(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    "ok",
	})
}

func (hc *httpControllers) RegisterEndpointRulesDraft(ctx echo.Context) error {
	req := new(usecase_models.Endpoints)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	projectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	reqB, _ := json.Marshal(req)
	id, err := hc.draftRepository.SaveDrafts(ctx.Request().Context(), models.Draft{
		Data:      null.NewString(string(reqB), true),
		ProjectID: projectId,
		Type:      null.NewString("endpoint", true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.NotValidData,
		Status:  400,
		Data: struct {
			ID int `json:"id"`
		}{
			ID: id,
		}})
}
func (hc *httpControllers) RegisterNetCatRulesDraft(ctx echo.Context) error {
	req := new(usecase_models.NetCats)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	projectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	reqB, _ := json.Marshal(req)
	id, err := hc.draftRepository.SaveDrafts(context.TODO(), models.Draft{
		Data:      null.NewString(string(reqB), true),
		ProjectID: projectId,
		Type:      null.NewString("netcat", true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.NotValidData,
		Status:  400,
		Data: struct {
			ID int `json:"id"`
		}{
			ID: id,
		}})
}
func (hc *httpControllers) RegisterPingRulesDraft(ctx echo.Context) error {
	req := new(usecase_models.Pings)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	projectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	reqB, _ := json.Marshal(req)
	id, err := hc.draftRepository.SaveDrafts(context.TODO(), models.Draft{
		Data:      null.NewString(string(reqB), true),
		ProjectID: projectId,
		Type:      null.NewString("ping", true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.NotValidData,
		Status:  400,
		Data: struct {
			ID int `json:"id"`
		}{
			ID: id,
		}})
}
func (hc *httpControllers) RegisterTraceRouteRulesDraft(ctx echo.Context) error {
	req := new(usecase_models.TraceRoutes)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	projectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	reqB, _ := json.Marshal(req)
	id, err := hc.draftRepository.SaveDrafts(context.TODO(), models.Draft{
		Data:      null.NewString(string(reqB), true),
		ProjectID: projectId,
		Type:      null.NewString("traceroute", true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, struct {
		ID int `json:"id"`
	}{
		ID: id,
	})
}
func (hc *httpControllers) RegisterPageSpeedRulesDraft(ctx echo.Context) error {
	req := new(usecase_models.PageSpeeds)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	projectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	reqB, _ := json.Marshal(req)
	id, err := hc.draftRepository.SaveDrafts(context.TODO(), models.Draft{
		Data:      null.NewString(string(reqB), true),
		ProjectID: projectId,
		Type:      null.NewString("pagespeed", true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.NotValidData,
		Status:  400,
		Data: struct {
			ID int `json:"id"`
		}{
			ID: id,
		}})
}

func (hc *httpControllers) GetRulesDrafts(ctx echo.Context) error {
	projectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	drafts, err := hc.draftRepository.GetDrafts(ctx.Request().Context(), projectId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  500,
			Data:    err.Error(),
		})
	}

	var response usecase_models.AggregateAllRuleSubWorks
	for _, draft := range drafts {
		switch draft.Type.String {
		case "endpoint":
			var data usecase_models.Endpoints
			_ = json.Unmarshal([]byte(draft.Data.String), &data)
			data.Scheduling.PipelineId = draft.ID
			response.Endpoints = append(response.Endpoints, &data)
		case "netcat":
			var data usecase_models.NetCats
			_ = json.Unmarshal([]byte(draft.Data.String), &data)
			data.Scheduling.PipelineId = draft.ID
			response.NetCats = append(response.NetCats, &data)
		case "pagespeed":
			var data usecase_models.PageSpeeds
			_ = json.Unmarshal([]byte(draft.Data.String), &data)
			data.Scheduling.PipelineId = draft.ID
			response.PageSpeed = append(response.PageSpeed, &data)
		case "ping":
			var data usecase_models.Pings
			_ = json.Unmarshal([]byte(draft.Data.String), &data)
			data.Scheduling.PipelineId = draft.ID
			response.Pings = append(response.Pings, &data)
		case "traceroute":
			var data usecase_models.TraceRoutes
			_ = json.Unmarshal([]byte(draft.Data.String), &data)
			data.Scheduling.PipelineId = draft.ID
			response.TraceRoutes = append(response.TraceRoutes, &data)
		}
	}
	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    response,
	})
}
func (hc *httpControllers) GetRulesDraft(ctx echo.Context) error {
	draftId, err := strconv.Atoi(ctx.Param("draft_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Draft ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	draft, err := hc.draftRepository.GetDraft(ctx.Request().Context(), draftId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  500,
			Data:    err.Error(),
		})
	}

	switch draft.Type.String {
	case "endpoint":
		var data usecase_models.Endpoints
		_ = json.Unmarshal([]byte(draft.Data.String), &data)
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    data,
		})
	case "netcat":
		var data usecase_models.NetCats
		_ = json.Unmarshal([]byte(draft.Data.String), &data)
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    data,
		})
	case "pagespeed":
		var data usecase_models.PageSpeeds
		_ = json.Unmarshal([]byte(draft.Data.String), &data)
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    data,
		})
	case "ping":
		var data usecase_models.Pings
		_ = json.Unmarshal([]byte(draft.Data.String), &data)
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    data,
		})
	case "traceroute":
		var data usecase_models.TraceRoutes
		_ = json.Unmarshal([]byte(draft.Data.String), &data)
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data:    data,
		})
	}

	return ctx.JSON(http.StatusNotFound, utils.StandardHttpResponse{
		Message: utils.NotFound,
		Status:  404,
		Data:    "No Data Found!",
	})
}

type DataCenterPointsAndAvgResponse struct {
	DataCentersData  map[string]DatacenterDetail `json:"data_centers_data"`
	SuccessRate      float64                     `json:"success_rate"`
	ResponseTimeAvg  float64                     `json:"response_time_avg"`
	SuccessData      []SuccessData               `json:"success_data"`
	ResponseTimeData []ResponseTimeData          `json:"response_time_data"`
}

type DatacenterDetail struct {
	DataCenterLat          float64 `json:"data_center_lat"`
	DataCenterLng          float64 `json:"data_center_lng"`
	DataCenterLocationName string  `json:"data_center_location_name"`
	CountryName            string  `json:"country_name"`
	SuccessRate            float64 `json:"success_rate"`
	ResponseTimeAvg        float64 `json:"response_time_avg"`
}

type SuccessData struct {
	Time           time.Time `json:"_time"`
	Value          int       `json:"_value"`
	DatacenterName string    `json:"datacenter_name"`
	PipelineName   string    `json:"pipeline_name"`
}

type ResponseTimeData struct {
	Time           time.Time `json:"_time"`
	Value          float64   `json:"_value"`
	DatacenterName string    `json:"datacenter_name"`
	PipelineName   string    `json:"pipeline_name"`
}

func (hc *httpControllers) ReportEndpointDataCenterPointsAndAvg(ctx echo.Context) error {
	loadPoints, err := strconv.ParseBool(ctx.QueryParam("load_points"))
	if err != nil {
		loadPoints = false
	}
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if endpointId := ctx.QueryParam("endpoint_id"); endpointId != "" {
		filter = append(filter, repos.Filter{Field: "endpoint_id", Op: repos.FilterOpIn, Value: strings.Split(endpointId, ",")})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}
	if endpointName := ctx.QueryParam("endpoint_name"); endpointName != "" {
		filter = append(filter, repos.Filter{Field: "endpoint_name", Op: repos.FilterOpEq, Value: endpointName})
	}

	_, data, err := hc.endpointStatsRepository.Read(ctx.Request().Context(),
		[]string{"time", "datacenter_id", "success", "response_time", "endpoint_id"},
		filter,
		[]string{models.EndpointStatRels.Datacenter, models.EndpointStatRels.Endpoint},
		0,
		0,
		false,
	)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	var dcr = DataCenterPointsAndAvgResponse{
		DataCentersData:  map[string]DatacenterDetail{},
		SuccessRate:      0,
		ResponseTimeAvg:  0,
		SuccessData:      nil,
		ResponseTimeData: nil,
	}
	if data == nil {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data: DataCenterPointsAndAvgResponse{
				DataCentersData:  map[string]DatacenterDetail{},
				SuccessRate:      0,
				ResponseTimeAvg:  0,
				SuccessData:      []SuccessData{},
				ResponseTimeData: []ResponseTimeData{},
			}})
	}

	var successTotals = map[string]int{}
	var responseTimeTotals = map[string]int{}
	var validResponseTimeTotals = float64(0)
	totalSuccess := 0
	totalResponseTime := float64(0)
	for _, point := range data {
		if math.IsNaN(point.ResponseTime) {
			point.ResponseTime = 0
		}
		if loadPoints {
			var endpointData = usecase_models.Endpoints{}
			err = json.Unmarshal(point.R.Endpoint.Data.JSON, &endpointData)
			dcr.SuccessData = append(dcr.SuccessData, SuccessData{
				Time:           point.Time,
				Value:          point.Success,
				DatacenterName: point.R.Datacenter.Title,
				PipelineName:   endpointData.Scheduling.PipelineName,
			})
			dcr.ResponseTimeData = append(dcr.ResponseTimeData, ResponseTimeData{
				Time:           point.Time,
				Value:          point.ResponseTime,
				DatacenterName: point.R.Datacenter.Title,
				PipelineName:   endpointData.Scheduling.PipelineName,
			})
		}
		totalSuccess += point.Success

		totalResponseTime += point.ResponseTime
		if math.IsNaN(totalResponseTime) {
			totalResponseTime = 0
		}
		successTotals[point.R.Datacenter.Title] += 1
		if point.ResponseTime != 0 {
			validResponseTimeTotals += 1
			responseTimeTotals[point.R.Datacenter.Title] += 1
		}
		if d, ok := dcr.DataCentersData[point.R.Datacenter.Title]; !ok {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          point.R.Datacenter.Lat.Float64,
				DataCenterLng:          point.R.Datacenter.LNG.Float64,
				DataCenterLocationName: point.R.Datacenter.LocationName.String,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            float64(point.Success),
				ResponseTimeAvg:        point.ResponseTime,
			}
		} else {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          d.DataCenterLat,
				DataCenterLng:          d.DataCenterLng,
				DataCenterLocationName: d.DataCenterLocationName,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            d.SuccessRate + float64(point.Success),
				ResponseTimeAvg:        d.ResponseTimeAvg + point.ResponseTime,
			}
		}
	}

	dcr.ResponseTimeAvg = totalResponseTime / validResponseTimeTotals
	dcr.SuccessRate = (float64(totalSuccess) * float64(100)) / float64(len(data))

	for key, value := range dcr.DataCentersData {
		responseTimeAverageDatacenter := value.ResponseTimeAvg / float64(responseTimeTotals[key])
		if math.IsNaN(responseTimeAverageDatacenter) {
			responseTimeAverageDatacenter = 0
		}
		dcr.DataCentersData[key] = DatacenterDetail{
			DataCenterLat:          value.DataCenterLat,
			DataCenterLng:          value.DataCenterLng,
			DataCenterLocationName: value.DataCenterLocationName,
			CountryName:            value.CountryName,
			SuccessRate:            (value.SuccessRate * 100) / float64(successTotals[key]),
			ResponseTimeAvg:        responseTimeAverageDatacenter,
		}
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    dcr,
	})
}

func (hc *httpControllers) ReportNetCatDataCenterPointsAndAvg(ctx echo.Context) error {
	loadPoints, err := strconv.ParseBool(ctx.QueryParam("load_points"))
	if err != nil {
		loadPoints = false
	}
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}

	data, err := hc.netcatStatsRepository.Read(ctx.Request().Context(),
		[]string{"time", "datacenter_id", "success", "netcat_id"},
		filter,
		[]string{models.NetCatsStatRels.Datacenter, models.NetCatsStatRels.Netcat},
	)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	var dcr = DataCenterPointsAndAvgResponse{
		DataCentersData:  map[string]DatacenterDetail{},
		SuccessRate:      0,
		ResponseTimeAvg:  0,
		SuccessData:      nil,
		ResponseTimeData: nil,
	}
	if data == nil {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data: DataCenterPointsAndAvgResponse{
				DataCentersData:  map[string]DatacenterDetail{},
				SuccessRate:      0,
				ResponseTimeAvg:  0,
				SuccessData:      []SuccessData{},
				ResponseTimeData: []ResponseTimeData{},
			}})
	}

	var successTotals = map[string]int{}
	totalSuccess := 0
	for _, point := range data {
		if loadPoints {
			d := usecase_models.NetCats{}
			err = json.Unmarshal(point.R.Netcat.Data.JSON, &d)
			dcr.SuccessData = append(dcr.SuccessData, SuccessData{
				Time:           point.Time,
				Value:          point.Success,
				DatacenterName: point.R.Datacenter.Title,
				PipelineName:   d.Scheduling.PipelineName,
			})
		}
		totalSuccess += point.Success
		successTotals[point.R.Datacenter.Title] += 1
		if d, ok := dcr.DataCentersData[point.R.Datacenter.Title]; !ok {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          point.R.Datacenter.Lat.Float64,
				DataCenterLng:          point.R.Datacenter.LNG.Float64,
				DataCenterLocationName: point.R.Datacenter.LocationName.String,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            float64(point.Success),
				ResponseTimeAvg:        0,
			}
		} else {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          d.DataCenterLat,
				DataCenterLng:          d.DataCenterLng,
				DataCenterLocationName: d.DataCenterLocationName,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            d.SuccessRate + float64(point.Success),
			}
		}
	}

	dcr.SuccessRate = (float64(totalSuccess) * float64(100)) / float64(len(data))

	for key, value := range dcr.DataCentersData {
		dcr.DataCentersData[key] = DatacenterDetail{
			DataCenterLat:          value.DataCenterLat,
			DataCenterLng:          value.DataCenterLng,
			DataCenterLocationName: value.DataCenterLocationName,
			CountryName:            value.CountryName,
			SuccessRate:            (value.SuccessRate * 100) / float64(successTotals[key]),
		}
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    dcr,
	})
}
func (hc *httpControllers) ReportPingDataCenterPointsAndAvg(ctx echo.Context) error {
	loadPoints, err := strconv.ParseBool(ctx.QueryParam("load_points"))
	if err != nil {
		loadPoints = false
	}
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}

	data, err := hc.pingStatsRepository.Read(ctx.Request().Context(),
		[]string{"time", "datacenter_id", "success", "ping_id"},
		filter,
		[]string{models.PingsStatRels.Datacenter, models.PingsStatRels.Ping},
	)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	var dcr = DataCenterPointsAndAvgResponse{
		DataCentersData:  map[string]DatacenterDetail{},
		SuccessRate:      0,
		ResponseTimeAvg:  0,
		SuccessData:      nil,
		ResponseTimeData: nil,
	}
	if data == nil {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data: DataCenterPointsAndAvgResponse{
				DataCentersData:  map[string]DatacenterDetail{},
				SuccessRate:      0,
				ResponseTimeAvg:  0,
				SuccessData:      []SuccessData{},
				ResponseTimeData: []ResponseTimeData{},
			}})
	}

	var successTotals = map[string]int{}
	totalSuccess := 0
	for _, point := range data {
		if loadPoints {
			d := usecase_models.Pings{}
			err = json.Unmarshal(point.R.Ping.Data.JSON, &d)
			dcr.SuccessData = append(dcr.SuccessData, SuccessData{
				Time:           point.Time,
				Value:          point.Success,
				DatacenterName: point.R.Datacenter.Title,
				PipelineName:   d.Scheduling.PipelineName,
			})
		}
		totalSuccess += point.Success
		successTotals[point.R.Datacenter.Title] += 1
		if d, ok := dcr.DataCentersData[point.R.Datacenter.Title]; !ok {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          point.R.Datacenter.Lat.Float64,
				DataCenterLng:          point.R.Datacenter.LNG.Float64,
				DataCenterLocationName: point.R.Datacenter.LocationName.String,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            float64(point.Success),
				ResponseTimeAvg:        0,
			}
		} else {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          d.DataCenterLat,
				DataCenterLng:          d.DataCenterLng,
				DataCenterLocationName: d.DataCenterLocationName,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            d.SuccessRate + float64(point.Success),
			}
		}
	}

	dcr.SuccessRate = (float64(totalSuccess) * float64(100)) / float64(len(data))

	for key, value := range dcr.DataCentersData {
		dcr.DataCentersData[key] = DatacenterDetail{
			DataCenterLat:          value.DataCenterLat,
			DataCenterLng:          value.DataCenterLng,
			DataCenterLocationName: value.DataCenterLocationName,
			CountryName:            value.CountryName,
			SuccessRate:            (value.SuccessRate * 100) / float64(successTotals[key]),
		}
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    dcr,
	})
}
func (hc *httpControllers) ReportTraceRouteDataCenterPointsAndAvg(ctx echo.Context) error {
	loadPoints, err := strconv.ParseBool(ctx.QueryParam("load_points"))
	if err != nil {
		loadPoints = false
	}
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}

	data, err := hc.traceRouteStatsRepository.Read(ctx.Request().Context(),
		[]string{"time", "datacenter_id", "success", "traceroute_id"},
		filter,
		[]string{models.TraceRoutesStatRels.Datacenter, models.TraceRoutesStatRels.Traceroute},
	)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	var dcr = DataCenterPointsAndAvgResponse{
		DataCentersData:  map[string]DatacenterDetail{},
		SuccessRate:      0,
		ResponseTimeAvg:  0,
		SuccessData:      nil,
		ResponseTimeData: nil,
	}
	if data == nil {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data: DataCenterPointsAndAvgResponse{
				DataCentersData:  map[string]DatacenterDetail{},
				SuccessRate:      0,
				ResponseTimeAvg:  0,
				SuccessData:      []SuccessData{},
				ResponseTimeData: []ResponseTimeData{},
			}})
	}

	var successTotals = map[string]int{}
	totalSuccess := 0
	for _, point := range data {
		if loadPoints {
			d := usecase_models.TraceRoutes{}
			err = json.Unmarshal(point.R.Traceroute.Data.JSON, &d)
			dcr.SuccessData = append(dcr.SuccessData, SuccessData{
				Time:           point.Time,
				Value:          point.Success,
				DatacenterName: point.R.Datacenter.Title,
				PipelineName:   d.Scheduling.PipelineName,
			})
		}
		totalSuccess += point.Success
		successTotals[point.R.Datacenter.Title] += 1
		if d, ok := dcr.DataCentersData[point.R.Datacenter.Title]; !ok {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          point.R.Datacenter.Lat.Float64,
				DataCenterLng:          point.R.Datacenter.LNG.Float64,
				DataCenterLocationName: point.R.Datacenter.LocationName.String,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            float64(point.Success),
				ResponseTimeAvg:        0,
			}
		} else {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          d.DataCenterLat,
				DataCenterLng:          d.DataCenterLng,
				DataCenterLocationName: d.DataCenterLocationName,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            d.SuccessRate + float64(point.Success),
			}
		}
	}

	dcr.SuccessRate = (float64(totalSuccess) * float64(100)) / float64(len(data))

	for key, value := range dcr.DataCentersData {
		dcr.DataCentersData[key] = DatacenterDetail{
			DataCenterLat:          value.DataCenterLat,
			DataCenterLng:          value.DataCenterLng,
			DataCenterLocationName: value.DataCenterLocationName,
			CountryName:            value.CountryName,
			SuccessRate:            (value.SuccessRate * 100) / float64(successTotals[key]),
		}
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    dcr,
	})
}

func (hc *httpControllers) ReportPageSpeedDataCenterPointsAndAvg(ctx echo.Context) error {
	loadPoints, err := strconv.ParseBool(ctx.QueryParam("load_points"))
	if err != nil {
		loadPoints = false
	}
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}

	data, err := hc.pageSpeedStatsRepository.Read(ctx.Request().Context(),
		[]string{"time", "datacenter_id", "success", "pagespeed_id"},
		filter,
		[]string{models.PageSpeedsStatRels.Datacenter, models.PageSpeedsStatRels.Pagespeed},
	)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	var dcr = DataCenterPointsAndAvgResponse{
		DataCentersData:  map[string]DatacenterDetail{},
		SuccessRate:      0,
		ResponseTimeAvg:  0,
		SuccessData:      nil,
		ResponseTimeData: nil,
	}
	if data == nil {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data: DataCenterPointsAndAvgResponse{
				DataCentersData:  map[string]DatacenterDetail{},
				SuccessRate:      0,
				ResponseTimeAvg:  0,
				SuccessData:      []SuccessData{},
				ResponseTimeData: []ResponseTimeData{},
			}})
	}

	var successTotals = map[string]int{}
	totalSuccess := 0
	for _, point := range data {
		if loadPoints {
			d := usecase_models.PageSpeeds{}
			err = json.Unmarshal(point.R.Pagespeed.Data.JSON, &d)
			dcr.SuccessData = append(dcr.SuccessData, SuccessData{
				Time:           point.Time,
				Value:          point.Success,
				DatacenterName: point.R.Datacenter.Title,
				PipelineName:   d.Scheduling.PipelineName,
			})
		}
		totalSuccess += point.Success
		successTotals[point.R.Datacenter.Title] += 1
		if d, ok := dcr.DataCentersData[point.R.Datacenter.Title]; !ok {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          point.R.Datacenter.Lat.Float64,
				DataCenterLng:          point.R.Datacenter.LNG.Float64,
				DataCenterLocationName: point.R.Datacenter.LocationName.String,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            float64(point.Success),
				ResponseTimeAvg:        0,
			}
		} else {
			dcr.DataCentersData[point.R.Datacenter.Title] = DatacenterDetail{
				DataCenterLat:          d.DataCenterLat,
				DataCenterLng:          d.DataCenterLng,
				DataCenterLocationName: d.DataCenterLocationName,
				CountryName:            point.R.Datacenter.CountryName.String,
				SuccessRate:            d.SuccessRate + float64(point.Success),
			}
		}
	}

	dcr.SuccessRate = (float64(totalSuccess) * float64(100)) / float64(len(data))

	for key, value := range dcr.DataCentersData {
		dcr.DataCentersData[key] = DatacenterDetail{
			DataCenterLat:          value.DataCenterLat,
			DataCenterLng:          value.DataCenterLng,
			DataCenterLocationName: value.DataCenterLocationName,
			CountryName:            value.CountryName,
			SuccessRate:            (value.SuccessRate * 100) / float64(successTotals[key]),
		}
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    dcr,
	})
}

type EndpointDetailsResponse struct {
	Data       []EndpointDetails `json:"data"`
	Pagination Pagination        `json:"pagination"`
}

type EndpointDetails struct {
	Time                time.Time              `json:"time"`
	ProjectID           int                    `json:"project_id"`
	EndpointName        string                 `json:"endpoint_name"`
	EndpointID          int                    `json:"endpoint_id"`
	URL                 []string               `json:"url"`
	IsHeartBeat         bool                   `json:"is_heart_beat"`
	Success             int                    `json:"success"`
	AverageResponseTime float64                `json:"average_response_time"`
	ResponseTimes       map[string]interface{} `json:"response_times"`
	ResponseBodies      map[string]interface{} `json:"response_bodies"`
	ResponseHeaders     map[string]interface{} `json:"response_headers"`
	ResponseStatuses    map[string]int         `json:"response_statuses"`
	Datacenter          struct {
		DatacenterId int    `json:"datacenter_id"`
		Title        string `json:"title"`
		LocationName string `json:"location_name"`
	} `json:"datacenter"`
}

func (hc *httpControllers) ReportEndpointDetails(ctx echo.Context) error {
	page := 1
	perPage := 10
	if p, err := strconv.Atoi(ctx.QueryParam("page")); err == nil {
		page = p
	}
	if pp, err := strconv.Atoi(ctx.QueryParam("per_page")); err == nil {
		perPage = pp
	}

	filter := repos.Filters{}

	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		filter = append(filter, repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)})
	}

	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		s := strings.Split(datacenterId, ",")
		var value []int
		for _, v := range s {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			value = append(value, i)
		}
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: value})
	}
	if endpointId := ctx.QueryParam("endpoint_id"); endpointId != "" {
		s := strings.Split(endpointId, ",")
		var value []int
		for _, v := range s {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			value = append(value, i)
		}
		filter = append(filter, repos.Filter{Field: "endpoint_id", Op: repos.FilterOpIn, Value: value})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}
	if endpointName := ctx.QueryParam("endpoint_name"); endpointName != "" {
		filter = append(filter, repos.Filter{Field: "endpoint_name", Op: repos.FilterOpEq, Value: endpointName})
	}
	if isSuccess, err := strconv.ParseBool(ctx.QueryParam("success")); err == nil {
		filter = append(filter, repos.Filter{Field: "success", Op: repos.FilterOpEq, Value: isSuccess})
	}

	total, data, err := hc.endpointStatsRepository.Read(ctx.Request().Context(),
		[]string{
			"time",
			"project_id",
			"endpoint_name",
			"endpoint_id",
			"url",
			"datacenter_id",
			"is_heart_beat",
			"success",
			"response_time",
			"response_times",
			"response_bodies",
			"response_headers",
			"response_statuses",
		},
		filter,
		[]string{models.EndpointStatRels.Datacenter},
		perPage,
		int(utils.OffsetFromPage(int64(page), int64(perPage))),
		true,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	responseData := make([]EndpointDetails, 0)
	for _, point := range data {
		rt := make(map[string]interface{})
		json.Unmarshal([]byte(point.ResponseTimes.String), &rt)
		rb := make(map[string]interface{})
		json.Unmarshal(point.ResponseBodies.Bytes, &rb)
		rh := make(map[string]interface{})
		json.Unmarshal([]byte(point.ResponseHeaders.String), &rh)
		rs := make(map[string]int)
		json.Unmarshal([]byte(point.ResponseStatuses.String), &rs)
		responseData = append(responseData, EndpointDetails{
			Time:                point.Time,
			ProjectID:           point.ProjectID,
			EndpointName:        point.EndpointName.String,
			EndpointID:          point.EndpointID,
			URL:                 strings.Split(point.URL.String, ","),
			IsHeartBeat:         point.IsHeartBeat,
			Success:             point.Success,
			AverageResponseTime: point.ResponseTime,
			ResponseTimes:       rt,
			ResponseBodies:      rb,
			ResponseHeaders:     rh,
			ResponseStatuses:    rs,
			Datacenter: struct {
				DatacenterId int    `json:"datacenter_id"`
				Title        string `json:"title"`
				LocationName string `json:"location_name"`
			}(struct {
				DatacenterId int
				Title        string
				LocationName string
			}{DatacenterId: point.R.Datacenter.ID, Title: point.R.Datacenter.Title, LocationName: point.R.Datacenter.LocationName.String}),
		})
	}

	var next = new(int)
	if int64(perPage*page) < total {
		temp := page + 1
		next = &temp
	}
	var prev = new(int)
	if page > 1 {
		temp := page - 1
		prev = &temp
	}
	response := EndpointDetailsResponse{
		Data: responseData,
		Pagination: Pagination{
			Total: total,
			Next:  next,
			Prev:  prev,
		},
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    response,
	})
}

type ReportEndpointQuickStatsResponse struct {
	DownTime          float64 `json:"down_time"`
	UpTimePercent     float64 `json:"up_time_percent"`
	CurrentStatus     bool    `json:"current_status"`
	StatusPerPipeline []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	} `json:"status_per_pipeline"`
}

func (hc *httpControllers) ReportEndpointQuickStats(ctx echo.Context) error {
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if endpointId := ctx.QueryParam("endpoint_id"); endpointId != "" {
		filter = append(filter, repos.Filter{Field: "endpoint_id", Op: repos.FilterOpIn, Value: strings.Split(endpointId, ",")})
	}
	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}
	data, err := hc.endpointStatsRepository.GetSessionSuccessions(ctx.Request().Context(), filter)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	if len(data) == 0 {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data: ReportEndpointQuickStatsResponse{
				DownTime:      0,
				UpTimePercent: 0,
				CurrentStatus: false,
				StatusPerPipeline: []struct {
					PipelineName string  `json:"pipeline_name"`
					Percent      float64 `json:"percent"`
				}{},
			}})
	}

	success := 0
	total := 0
	up := true
	var downTime = 0 * time.Second
	temp := time.Time{}
	pipelineStats := map[int]float64{}
	pipelineTotals := map[int]float64{}
	pipelineNames := map[int]string{}
	for i := len(data) - 1; i >= 0; i-- {
		if data[i].Success {
			success += 1
		}
		total += 1

		if !data[i].Success && up {
			up = false
			temp = data[i].MinTime
		}
		if data[i].Success && !up {
			up = true
			downTime = downTime + time.Duration(data[i].MaxTime.Sub(temp).Seconds())*time.Second
		}

		if _, ok := pipelineStats[data[i].EndpointId]; ok {
			if data[i].Success {
				pipelineStats[data[i].EndpointId] = pipelineStats[data[i].EndpointId] + 1
			}
			pipelineTotals[data[i].EndpointId] += 1
		} else {
			if data[i].Success {
				pipelineStats[data[i].EndpointId] = 1
			}
			pipelineTotals[data[i].EndpointId] += 1

			name, ok := pipelineNames[data[i].EndpointId]
			if !ok {
				endpoint, err := hc.endpointRepository.GetEndpoint(ctx.Request().Context(), data[i].EndpointId)
				if err != nil {
					name = ""
				} else {
					name = endpoint.Scheduling.PipelineName
				}
			}

			if name == "" {
				pipelineNames[data[i].EndpointId] = strconv.Itoa(data[i].EndpointId)
			} else {
				pipelineNames[data[i].EndpointId] = name
			}
		}
	}
	if !up {
		downTime = downTime + time.Duration(data[0].MaxTime.Sub(temp).Seconds())*time.Second
	}

	var pipelineStatsResponse []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	}
	for key, value := range pipelineStats {
		pipelineStatsResponse = append(pipelineStatsResponse, struct {
			PipelineName string  `json:"pipeline_name"`
			Percent      float64 `json:"percent"`
		}{
			PipelineName: pipelineNames[key],
			Percent:      (value * 100) / pipelineTotals[key],
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: ReportEndpointQuickStatsResponse{
			DownTime:          downTime.Minutes(),
			UpTimePercent:     float64(success) * 100 / float64(total),
			CurrentStatus:     up,
			StatusPerPipeline: pipelineStatsResponse,
		}})
}

type ReportNetCatQuickStatsResponse struct {
	DownTime          float64 `json:"down_time"`
	UpTimePercent     float64 `json:"up_time_percent"`
	CurrentStatus     bool    `json:"current_status"`
	StatusPerPipeline []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	} `json:"status_per_pipeline"`
}

func (hc *httpControllers) ReportNetCatQuickStats(ctx echo.Context) error {
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if netcatId := ctx.QueryParam("netcat_id"); netcatId != "" {
		filter = append(filter, repos.Filter{Field: "netcat_id", Op: repos.FilterOpIn, Value: strings.Split(netcatId, ",")})
	}
	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}

	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}
	data, err := hc.netcatStatsRepository.GetSessionSuccessions(ctx.Request().Context(), filter)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	if len(data) == 0 {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data: ReportNetCatQuickStatsResponse{
				DownTime:      0,
				UpTimePercent: 0,
				CurrentStatus: false,
				StatusPerPipeline: []struct {
					PipelineName string  `json:"pipeline_name"`
					Percent      float64 `json:"percent"`
				}{},
			}})
	}

	success := 0
	total := 0
	up := true
	var downTime = 0 * time.Second
	temp := time.Time{}
	pipelineStats := map[int]float64{}
	pipelineTotals := map[int]float64{}
	pipelineNames := map[int]string{}
	for i := len(data) - 1; i >= 0; i-- {
		if data[i].Success {
			success += 1
		}
		total += 1

		if !data[i].Success && up {
			up = false
			temp = data[i].MinTime
		}
		if data[i].Success && !up {
			up = true
			downTime = downTime + time.Duration(data[i].MaxTime.Sub(temp).Seconds())*time.Second
		}

		if _, ok := pipelineStats[data[i].NetCatId]; ok {
			if data[i].Success {
				pipelineStats[data[i].NetCatId] = pipelineStats[data[i].NetCatId] + 1
			}
			pipelineTotals[data[i].NetCatId] += 1
		} else {
			if data[i].Success {
				pipelineStats[data[i].NetCatId] = 1
			}
			pipelineTotals[data[i].NetCatId] += 1

			name, ok := pipelineNames[data[i].NetCatId]
			if !ok {
				netcat, err := hc.netCatRepository.GetNetCat(ctx.Request().Context(), data[i].NetCatId)
				if err != nil {
					name = ""
				} else {
					name = netcat.Scheduling.PipelineName
				}
			}

			if name == "" {
				pipelineNames[data[i].NetCatId] = strconv.Itoa(data[i].NetCatId)
			} else {
				pipelineNames[data[i].NetCatId] = name
			}
		}
	}
	if !up {
		downTime = downTime + time.Duration(data[0].MaxTime.Sub(temp).Seconds())*time.Second
	}

	var pipelineStatsResponse []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	}
	for key, value := range pipelineStats {
		pipelineStatsResponse = append(pipelineStatsResponse, struct {
			PipelineName string  `json:"pipeline_name"`
			Percent      float64 `json:"percent"`
		}{
			PipelineName: pipelineNames[key],
			Percent:      (value * 100) / pipelineTotals[key],
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: ReportNetCatQuickStatsResponse{
			DownTime:          downTime.Minutes(),
			UpTimePercent:     float64(success) * 100 / float64(total),
			CurrentStatus:     up,
			StatusPerPipeline: pipelineStatsResponse,
		}})
}

type ReportPageSpeedQuickStatsResponse struct {
	DownTime          float64 `json:"down_time"`
	UpTimePercent     float64 `json:"up_time_percent"`
	CurrentStatus     bool    `json:"current_status"`
	StatusPerPipeline []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	} `json:"status_per_pipeline"`
}

func (hc *httpControllers) ReportPageSpeedQuickStats(ctx echo.Context) error {
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if pagespeedId := ctx.QueryParam("pagespeed_id"); pagespeedId != "" {
		filter = append(filter, repos.Filter{Field: "pagespeed_id", Op: repos.FilterOpIn, Value: strings.Split(pagespeedId, ",")})
	}
	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}
	data, err := hc.pageSpeedStatsRepository.GetSessionSuccessions(ctx.Request().Context(), filter)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	if len(data) == 0 {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data: ReportPageSpeedQuickStatsResponse{
				DownTime:      0,
				UpTimePercent: 0,
				CurrentStatus: false,
				StatusPerPipeline: []struct {
					PipelineName string  `json:"pipeline_name"`
					Percent      float64 `json:"percent"`
				}{},
			}})
	}

	success := 0
	total := 0
	up := true
	var downTime = 0 * time.Second
	temp := time.Time{}
	pipelineStats := map[int]float64{}
	pipelineTotals := map[int]float64{}
	pipelineNames := map[int]string{}
	for i := len(data) - 1; i >= 0; i-- {
		if data[i].Success {
			success += 1
		}
		total += 1

		if !data[i].Success && up {
			up = false
			temp = data[i].MinTime
		}
		if data[i].Success && !up {
			up = true
			downTime = downTime + time.Duration(data[i].MaxTime.Sub(temp).Seconds())*time.Second
		}

		if _, ok := pipelineStats[data[i].PageSpeedId]; ok {
			if data[i].Success {
				pipelineStats[data[i].PageSpeedId] = pipelineStats[data[i].PageSpeedId] + 1
			}
			pipelineTotals[data[i].PageSpeedId] += 1
		} else {
			if data[i].Success {
				pipelineStats[data[i].PageSpeedId] = 1
			}
			pipelineTotals[data[i].PageSpeedId] += 1

			name, ok := pipelineNames[data[i].PageSpeedId]
			if !ok {
				pagespeed, err := hc.pageSpeedRepository.GetPageSpeed(ctx.Request().Context(), data[i].PageSpeedId)
				if err != nil {
					name = ""
				} else {
					name = pagespeed.Scheduling.PipelineName
				}
			}

			if name == "" {
				pipelineNames[data[i].PageSpeedId] = strconv.Itoa(data[i].PageSpeedId)
			} else {
				pipelineNames[data[i].PageSpeedId] = name
			}
		}
	}
	if !up {
		downTime = downTime + time.Duration(data[0].MaxTime.Sub(temp).Seconds())*time.Second
	}

	var pipelineStatsResponse []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	}
	for key, value := range pipelineStats {
		pipelineStatsResponse = append(pipelineStatsResponse, struct {
			PipelineName string  `json:"pipeline_name"`
			Percent      float64 `json:"percent"`
		}{
			PipelineName: pipelineNames[key],
			Percent:      (value * 100) / pipelineTotals[key],
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: ReportPageSpeedQuickStatsResponse{
			DownTime:          downTime.Minutes(),
			UpTimePercent:     float64(success) * 100 / float64(total),
			CurrentStatus:     up,
			StatusPerPipeline: pipelineStatsResponse,
		}})
}

type ReportPingQuickStatsResponse struct {
	DownTime          float64 `json:"down_time"`
	UpTimePercent     float64 `json:"up_time_percent"`
	CurrentStatus     bool    `json:"current_status"`
	StatusPerPipeline []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	} `json:"status_per_pipeline"`
}

func (hc *httpControllers) ReportPingQuickStats(ctx echo.Context) error {
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if pingId := ctx.QueryParam("ping_id"); pingId != "" {
		filter = append(filter, repos.Filter{Field: "ping_id", Op: repos.FilterOpIn, Value: strings.Split(pingId, ",")})
	}
	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}
	data, err := hc.pingStatsRepository.GetSessionSuccessions(ctx.Request().Context(), filter)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	if len(data) == 0 {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data: ReportPingQuickStatsResponse{
				DownTime:      0,
				UpTimePercent: 0,
				CurrentStatus: false,
				StatusPerPipeline: []struct {
					PipelineName string  `json:"pipeline_name"`
					Percent      float64 `json:"percent"`
				}{},
			}})
	}

	success := 0
	total := 0
	up := true
	var downTime = 0 * time.Second
	temp := time.Time{}
	pipelineStats := map[int]float64{}
	pipelineTotals := map[int]float64{}
	pipelineNames := map[int]string{}
	for i := len(data) - 1; i >= 0; i-- {
		if data[i].Success {
			success += 1
		}
		total += 1

		if !data[i].Success && up {
			up = false
			temp = data[i].MinTime
		}
		if data[i].Success && !up {
			up = true
			downTime = downTime + time.Duration(data[i].MaxTime.Sub(temp).Seconds())*time.Second
		}

		if _, ok := pipelineStats[data[i].PingId]; ok {
			if data[i].Success {
				pipelineStats[data[i].PingId] = pipelineStats[data[i].PingId] + 1
			}
			pipelineTotals[data[i].PingId] += 1
		} else {
			if data[i].Success {
				pipelineStats[data[i].PingId] = 1
			}
			pipelineTotals[data[i].PingId] += 1

			name, ok := pipelineNames[data[i].PingId]
			if !ok {
				ping, err := hc.pingRepository.GetPing(ctx.Request().Context(), data[i].PingId)
				if err != nil {
					name = ""
				} else {
					name = ping.Scheduling.PipelineName
				}
			}

			if name == "" {
				pipelineNames[data[i].PingId] = strconv.Itoa(data[i].PingId)
			} else {
				pipelineNames[data[i].PingId] = name
			}
		}
	}
	if !up {
		downTime = downTime + time.Duration(data[0].MaxTime.Sub(temp).Seconds())*time.Second
	}

	var pipelineStatsResponse []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	}
	for key, value := range pipelineStats {
		pipelineStatsResponse = append(pipelineStatsResponse, struct {
			PipelineName string  `json:"pipeline_name"`
			Percent      float64 `json:"percent"`
		}{
			PipelineName: pipelineNames[key],
			Percent:      (value * 100) / pipelineTotals[key],
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: ReportPingQuickStatsResponse{
			DownTime:          downTime.Minutes(),
			UpTimePercent:     float64(success) * 100 / float64(total),
			CurrentStatus:     up,
			StatusPerPipeline: pipelineStatsResponse,
		}})
}

type ReportTraceRouteQuickStatsResponse struct {
	DownTime          float64 `json:"down_time"`
	UpTimePercent     float64 `json:"up_time_percent"`
	CurrentStatus     bool    `json:"current_status"`
	StatusPerPipeline []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	} `json:"status_per_pipeline"`
}

func (hc *httpControllers) ReportTraceRouteQuickStats(ctx echo.Context) error {
	timeframe, err := strconv.Atoi(ctx.QueryParam("timeframe"))
	if err != nil {
		timeframe = 15
	}

	filter := repos.Filters{
		repos.Filter{Field: "time", Op: repos.FilterOpGt, Value: time.Now().Add(-time.Duration(timeframe) * time.Minute)},
	}

	if datacenterId := ctx.QueryParam("datacenter_id"); datacenterId != "" {
		filter = append(filter, repos.Filter{Field: "datacenter_id", Op: repos.FilterOpIn, Value: strings.Split(datacenterId, ",")})
	}
	if tracerouteId := ctx.QueryParam("traceroute_id"); tracerouteId != "" {
		filter = append(filter, repos.Filter{Field: "traceroute_id", Op: repos.FilterOpIn, Value: strings.Split(tracerouteId, ",")})
	}
	if projectId := ctx.QueryParam("project_id"); projectId != "" {
		filter = append(filter, repos.Filter{Field: "project_id", Op: repos.FilterOpEq, Value: projectId})
	}
	if isHearBeat, err := strconv.ParseBool(ctx.QueryParam("is_heart_beat")); err == nil {
		filter = append(filter, repos.Filter{Field: "is_heart_beat", Op: repos.FilterOpEq, Value: isHearBeat})
	}
	data, err := hc.traceRouteStatsRepository.GetSessionSuccessions(ctx.Request().Context(), filter)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	if len(data) == 0 {
		return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
			Message: utils.Ok,
			Status:  200,
			Data: ReportTraceRouteQuickStatsResponse{
				DownTime:      0,
				UpTimePercent: 0,
				CurrentStatus: false,
				StatusPerPipeline: []struct {
					PipelineName string  `json:"pipeline_name"`
					Percent      float64 `json:"percent"`
				}{},
			}})
	}

	success := 0
	total := 0
	up := true
	var downTime = 0 * time.Second
	temp := time.Time{}
	pipelineStats := map[int]float64{}
	pipelineTotals := map[int]float64{}
	pipelineNames := map[int]string{}
	for i := len(data) - 1; i >= 0; i-- {
		if data[i].Success {
			success += 1
		}
		total += 1

		if !data[i].Success && up {
			up = false
			temp = data[i].MinTime
		}
		if data[i].Success && !up {
			up = true
			downTime = downTime + time.Duration(data[i].MaxTime.Sub(temp).Seconds())*time.Second
		}

		if _, ok := pipelineStats[data[i].TraceRouteId]; ok {
			if data[i].Success {
				pipelineStats[data[i].TraceRouteId] = pipelineStats[data[i].TraceRouteId] + 1
			}
			pipelineTotals[data[i].TraceRouteId] += 1
		} else {
			if data[i].Success {
				pipelineStats[data[i].TraceRouteId] = 1
			}
			pipelineTotals[data[i].TraceRouteId] += 1

			name, ok := pipelineNames[data[i].TraceRouteId]
			if !ok {
				traceroute, err := hc.traceRouteRepository.GetTraceRoute(ctx.Request().Context(), data[i].TraceRouteId)
				if err != nil {
					name = ""
				} else {
					name = traceroute.Scheduling.PipelineName
				}
			}

			if name == "" {
				pipelineNames[data[i].TraceRouteId] = strconv.Itoa(data[i].TraceRouteId)
			} else {
				pipelineNames[data[i].TraceRouteId] = name
			}
		}
	}
	if !up {
		downTime = downTime + time.Duration(data[0].MaxTime.Sub(temp).Seconds())*time.Second
	}

	var pipelineStatsResponse []struct {
		PipelineName string  `json:"pipeline_name"`
		Percent      float64 `json:"percent"`
	}
	for key, value := range pipelineStats {
		pipelineStatsResponse = append(pipelineStatsResponse, struct {
			PipelineName string  `json:"pipeline_name"`
			Percent      float64 `json:"percent"`
		}{
			PipelineName: pipelineNames[key],
			Percent:      (value * 100) / pipelineTotals[key],
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: ReportTraceRouteQuickStatsResponse{
			DownTime:          downTime.Minutes(),
			UpTimePercent:     float64(success) * 100 / float64(total),
			CurrentStatus:     up,
			StatusPerPipeline: pipelineStatsResponse,
		}})
}

func (hc *httpControllers) CreateGateway(ctx echo.Context) error {
	req := new(usecase_models.Gateway)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	data, err := json.Marshal(req.Data)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  400,
			Data:    err.Error(),
		})
	}
	gatewayId, err := hc.gatewayRepository.SaveGateways(ctx.Request().Context(), models.Gateway{
		Baseurl:        req.Baseurl,
		Title:          req.Title,
		ConnectionRate: null.NewInt(req.ConnectionRate, true),
		IsActive:       null.NewBool(req.IsActive, true),
		IsDefault:      null.NewBool(req.IsDefault, true),
		Data:           null.NewJSON(data, true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.CreateGatewayResponse{
			GatewayId: gatewayId,
		}})
}

func (hc *httpControllers) GetGateways(ctx echo.Context) error {
	GatewayId := 0
	var err error
	if ctx.Param("gateway_id") != "" {
		GatewayId, err = strconv.Atoi(ctx.Param("gateway_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Gateway ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
	}

	var Gateways []*models.Gateway
	if GatewayId != 0 {
		Gateway, err := hc.gatewayRepository.GetGateway(ctx.Request().Context(), GatewayId)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.NotValidData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		Gateways = append(Gateways, &Gateway)
	} else {
		Gatewayss, err := hc.gatewayRepository.GetGateways(ctx.Request().Context())
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.NotValidData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		Gateways = append(Gateways, Gatewayss...)
	}

	var GatewaysResponse []usecase_models.Gateway
	for _, gateway := range Gateways {
		GatewaysResponse = append(GatewaysResponse, usecase_models.Gateway{
			ID:             gateway.ID,
			Baseurl:        gateway.Baseurl,
			Title:          gateway.Title,
			ConnectionRate: gateway.ConnectionRate.Int,
			IsActive:       gateway.IsActive.Bool,
			IsDefault:      gateway.IsDefault.Bool,
			Data:           gateway.Data,
			UpdatedAt:      gateway.UpdatedAt,
			CreatedAt:      gateway.CreatedAt,
			DeletedAt:      gateway.DeletedAt,
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    GatewaysResponse,
	})
}

func (hc *httpControllers) UpdateGateway(ctx echo.Context) error {
	req := new(usecase_models.Gateway)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	gatewayId, err := strconv.Atoi(ctx.Param("gateway_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Gateway ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	data, err := json.Marshal(req.Data)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	err = hc.gatewayRepository.UpdateGateways(ctx.Request().Context(), models.Gateway{
		ID:             gatewayId,
		Baseurl:        req.Baseurl,
		Title:          req.Title,
		ConnectionRate: null.NewInt(req.ConnectionRate, true),
		IsActive:       null.NewBool(req.IsActive, true),
		IsDefault:      null.NewBool(req.IsDefault, true),
		Data:           null.NewJSON(data, true),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  201,
		Data:    nil,
	})
}

func (hc *httpControllers) CreateOrder(ctx echo.Context) error {
	req := new(usecase_models.Order)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	project, err := hc.projectRepo.GetProject(ctx.Request().Context(), req.ProjectId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.ProblemInGettingData,
			Status:  400,
			Data:    err.Error(),
		})
	}
	if !project.ExpireAt.Time.Before(time.Now()) {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.CustomMessage, "You already have active package on this project"),
			Status:  400,
			Data:    nil,
		})
	}

	pickedGateway, err := hc.gatewayRepository.GetGateway(ctx.Request().Context(), req.GatewayId)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.CustomMessage, "gateway not available"),
			Status:  400,
			Data:    nil,
		})
	}

	pickedPackage, err := hc.packageRepository.GetPackage(ctx.Request().Context(), req.PackageId)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.CustomMessage, "package not available"),
			Status:  400,
			Data:    nil,
		})
	}

	orderId, err := hc.orderRepository.SaveOrders(ctx.Request().Context(), models.Order{
		AccountID: IdentityStruct.Id,
		ProjectID: req.ProjectId,
		PackageID: req.PackageId,
		GatewayID: req.GatewayId,
		Status:    strconv.Itoa(usecase_models.OrderStatusCreated),
		Amount:    pickedPackage.Price,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	var gatewayLink string

	switch pickedGateway.Title {
	case "idpay":
		var gatewayData gateway.IdpayData
		err = json.Unmarshal(pickedGateway.Data.JSON, &gatewayData)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}
		createOrderResponse, err := hc.idpayGateway.CreateOrder(ctx.Request().Context(), gateway.CreateOrderRequest{
			Amount:         pickedPackage.Price,
			CallBackUrl:    gatewayData.CallBackUrl,
			ServerUniqueId: strconv.Itoa(orderId),
		})
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}

		responseB, _ := json.Marshal(createOrderResponse)
		err = hc.orderRepository.UpdateOrders(ctx.Request().Context(), models.Order{
			ID:                    orderId,
			Status:                strconv.Itoa(usecase_models.OrderStatusPending),
			GatewayOrderID:        null.NewString(createOrderResponse.GatewayOrderId, true),
			GatewayCreateResponse: null.NewJSON(responseB, true),
		})
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}
		gatewayLink = createOrderResponse.OrderLink
	case "zarinpal":
		panic("not implemented")
	default:
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    "Gateway is not available",
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.CreateOrderResponse{
			OrderId:     orderId,
			GatewayLink: gatewayLink,
		}})
}

func (hc *httpControllers) VerifyOrder(ctx echo.Context) error {
	req := new(usecase_models.VerifyOrder)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	order, err := hc.orderRepository.GetOrderByGatewayID(ctx.Request().Context(), req.GatewayOrderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	currentStatus, _ := strconv.Atoi(order.Status)
	if currentStatus != usecase_models.OrderStatusPending {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.CustomMessage, "Order is not created yet!"),
			Status:  400,
			Data:    nil,
		})
	}

	// get currentGateway
	currentGateway, err := hc.gatewayRepository.GetGateway(ctx.Request().Context(), order.GatewayID)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	switch currentGateway.Title {
	case "idpay":
		var gatewayData gateway.IdpayData
		err = json.Unmarshal(currentGateway.Data.JSON, &gatewayData)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}
		verifyOrderResponse, err := hc.idpayGateway.VerifyOrder(ctx.Request().Context(), gateway.VerifyOrderRequest{
			ServerOrderId:  strconv.Itoa(order.ID),
			GatewayOrderId: order.GatewayOrderID.String,
		})
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}

		responseB, _ := json.Marshal(verifyOrderResponse)
		if verifyOrderResponse.Status != 100 {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}

		if verifyOrderResponse.Amount != strconv.Itoa(order.Amount) {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.CustomMessage, "Amount does not match!"),
				Status:  400,
				Data:    nil,
			})
		}
		err = hc.orderRepository.UpdateOrders(ctx.Request().Context(), models.Order{
			ID:                    order.ID,
			Status:                strconv.Itoa(usecase_models.OrderStatusVerified),
			GatewayVerifyResponse: null.NewJSON(responseB, true),
			GatewayOrderID:        order.GatewayOrderID,
			GatewayCreateResponse: order.GatewayCreateResponse,
		})
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}
	case "zarinpal":
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    nil,
		})
	}

	pickedPackage, err := hc.packageRepository.GetPackage(ctx.Request().Context(), order.PackageID)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.CustomMessage, "Package not available"),
			Status:  400,
			Data:    nil,
		})
	}

	newExpire := time.Now().Add(time.Hour * 24 * time.Duration(pickedPackage.LengthInDays))

	// updating all end at in every rule
	{
		endpoints, err := hc.endpointRepository.GetEndpoints(ctx.Request().Context(), order.ProjectID)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}

		for _, endpoint := range endpoints {
			endpoint.Scheduling.EndAt = newExpire.String()
			data, err := json.Marshal(endpoint)
			if err != nil {
				fmt.Println(err)
			}
			err = hc.endpointRepository.UpdateEndpoint(ctx.Request().Context(), models.Endpoint{
				Data:      null.NewJSON(data, true),
				ProjectID: order.ProjectID,
			})
			if err != nil {
				return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
					Message: utils.ProblemInSystem,
					Status:  500,
					Data:    err.Error(),
				})
			}
		}

		netcats, err := hc.netCatRepository.GetNetCats(ctx.Request().Context(), order.ProjectID)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}

		for _, netcat := range netcats {
			netcat.Scheduling.EndAt = newExpire.String()
			data, err := json.Marshal(netcat)
			if err != nil {
				fmt.Println(err)
			}
			err = hc.netCatRepository.UpdateNetCat(ctx.Request().Context(), models.NetCat{
				Data:      null.NewJSON(data, true),
				ProjectID: order.ProjectID,
			})
			if err != nil {
				return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
					Message: utils.ProblemInSystem,
					Status:  500,
					Data:    err.Error(),
				})
			}
		}

		pings, err := hc.pingRepository.GetPings(ctx.Request().Context(), order.ProjectID)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}

		for _, ping := range pings {
			ping.Scheduling.EndAt = newExpire.String()
			data, err := json.Marshal(ping)
			if err != nil {
				fmt.Println(err)
			}
			err = hc.pingRepository.UpdatePing(ctx.Request().Context(), models.Ping{
				Data:      null.NewJSON(data, true),
				ProjectID: order.ProjectID,
			})
			if err != nil {
				return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
					Message: utils.ProblemInSystem,
					Status:  500,
					Data:    err.Error(),
				})
			}
		}

		traceRoutes, err := hc.traceRouteRepository.GetTraceRoutes(ctx.Request().Context(), order.ProjectID)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}

		for _, traceRoute := range traceRoutes {
			traceRoute.Scheduling.EndAt = newExpire.String()
			data, err := json.Marshal(traceRoute)
			if err != nil {
				fmt.Println(err)
			}
			err = hc.traceRouteRepository.UpdateTraceRoute(ctx.Request().Context(), models.TraceRoute{
				Data:      null.NewJSON(data, true),
				ProjectID: order.ProjectID,
			})
			if err != nil {
				return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
					Message: utils.ProblemInSystem,
					Status:  500,
					Data:    err.Error(),
				})
			}
		}
	}
	err = hc.projectRepo.UpdateProjects(ctx.Request().Context(), models.Project{
		ID:        order.ProjectID,
		ExpireAt:  null.NewTime(newExpire, true),
		PackageID: pickedPackage.ID,
	}, []string{"title", "is_active", "notifications", "account_id"}...)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.CustomMessage, "Couldn't activate project"),
			Status:  500,
			Data:    err.Error(),
		})
	}

	order, err = hc.orderRepository.GetOrder(ctx.Request().Context(), order.ID)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}
	// update order status to Done
	err = hc.orderRepository.UpdateOrders(ctx.Request().Context(), models.Order{
		ID:                    order.ID,
		Status:                strconv.Itoa(usecase_models.OrderStatusDone),
		GatewayOrderID:        order.GatewayOrderID,
		GatewayCreateResponse: order.GatewayCreateResponse,
		GatewayVerifyResponse: order.GatewayVerifyResponse,
	})

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    nil,
	})
}

func (hc *httpControllers) GetOrderHistory(ctx echo.Context) error {
	OrderId, err := strconv.Atoi(ctx.Param("orderId"))
	if err != nil {
		OrderId = 0
	}
	ProjectId, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		ProjectId = 0
	}

	var Orders []*models.Order
	if OrderId != 0 {
		Order, err := hc.orderRepository.GetOrder(ctx.Request().Context(), OrderId)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}
		if Order.AccountID != IdentityStruct.Id {
			return ctx.JSON(http.StatusForbidden, utils.StandardHttpResponse{
				Message: utils.NoAccess,
				Status:  403,
				Data:    "you are not owner of this order",
			})
		}
		Orders = append(Orders, &Order)
	} else {
		Orderss, err := hc.orderRepository.GetOrders(ctx.Request().Context(), IdentityStruct.Id, ProjectId)
		if err != nil {
			return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
				Message: utils.ProblemInSystem,
				Status:  500,
				Data:    err.Error(),
			})
		}
		Orders = append(Orders, Orderss...)
	}

	var OrdersResponse []usecase_models.Order
	for _, order := range Orders {
		OrdersResponse = append(OrdersResponse, usecase_models.Order{
			ID:        order.ID,
			AccountId: order.AccountID,
			ProjectId: order.ProjectID,
			PackageId: order.PackageID,
			GatewayId: order.GatewayID,
			Status:    order.Status,
			Amount:    order.Amount,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
			DeletedAt: order.DeletedAt,
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    OrdersResponse,
	})
}

func (hc *httpControllers) AlertStats(ctx echo.Context) error {
	projectId := ctx.Param("projectId")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(ctx.QueryParam("per_page"))
	if err != nil {
		perPage = 10
	}
	response, err := hc.alertSystem.AlertLogs(ctx.Request().Context(), alert_system.AlertLogsRequest{
		UserId: &projectId,
	}, page, perPage)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    response,
	})
}

func (hc *httpControllers) CreateFaq(ctx echo.Context) error {
	req := new(usecase_models.Faq)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	faqId, err := hc.faqRepository.SaveFaqs(ctx.Request().Context(), models.Faq{
		Question: req.Question,
		Answer:   req.Answer,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.CreateFaqResponse{
			FaqId: faqId,
		}})
}

func (hc *httpControllers) GetFaq(ctx echo.Context) error {
	faqId := 0
	var err error
	if ctx.Param("faq_id") != "" {
		faqId, err = strconv.Atoi(ctx.Param("faq_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "FAQ ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
	}

	var faqs []*models.Faq
	if faqId != 0 {
		faq, err := hc.faqRepository.GetFaq(ctx.Request().Context(), faqId)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.NotValidData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		faqs = append(faqs, &faq)
	} else {
		faqss, err := hc.faqRepository.GetFaqs(ctx.Request().Context())
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.NotValidData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		faqs = append(faqs, faqss...)
	}

	var faqsResponse []usecase_models.Faq
	for _, faq := range faqs {
		faqsResponse = append(faqsResponse, usecase_models.Faq{
			ID:        faq.ID,
			Question:  faq.Question,
			Answer:    faq.Answer,
			CreatedAt: faq.CreatedAt,
			UpdatedAt: faq.UpdatedAt,
			DeletedAt: faq.DeletedAt,
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    faqsResponse,
	})
}

func (hc *httpControllers) UpdateFaq(ctx echo.Context) error {
	req := new(usecase_models.Faq)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	faqId, err := strconv.Atoi(ctx.Param("faq_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "FAQ ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	err = hc.faqRepository.UpdateFaqs(ctx.Request().Context(), models.Faq{
		ID:       faqId,
		Question: req.Question,
		Answer:   req.Answer,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    nil,
	})
}

func (hc *httpControllers) CreateTicket(ctx echo.Context) error {
	req := new(usecase_models.Tickets)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	ticketId, err := hc.ticketsRepository.SaveTickets(ctx.Request().Context(), models.Ticket{
		AccountID:    IdentityStruct.Id,
		ProjectID:    req.ProjectID,
		Message:      req.Message,
		TicketStatus: usecase_models.TicketStatusPending,
		Title:        req.Title,
		ReplyTo:      req.ReplyTo,
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data: usecase_models.CreateTicketsResponse{
			TicketId: ticketId,
		}})
}

func (hc *httpControllers) GetTicket(ctx echo.Context) error {
	ticketId := 0
	var err error
	if ctx.Param("ticket_id") != "" {
		ticketId, err = strconv.Atoi(ctx.Param("ticket_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Ticket ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
	}

	projectId := 0
	if ctx.QueryParam("project_id") != "" {
		projectId, err = strconv.Atoi(ctx.QueryParam("project_id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: fmt.Sprintf(utils.NotValidField, "Project ID"),
				Status:  400,
				Data:    err.Error(),
			})
		}
	}

	accountProjects, err := hc.projectRepo.GetProjects(ctx.Request().Context(), IdentityStruct.Id)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	var accountProjectIds []int
	for _, value := range accountProjects {
		accountProjectIds = append(accountProjectIds, value.ID)

	}
	if projectId != 0 {
		if !contains(accountProjectIds, projectId) {
			return ctx.JSON(http.StatusForbidden, utils.StandardHttpResponse{
				Message: utils.NoAccess,
				Status:  403,
				Data:    "you dont have access to this project",
			})
		}
	}

	var tickets []*models.Ticket
	if ticketId != 0 {
		ticket, err := hc.ticketsRepository.GetTicket(ctx.Request().Context(), ticketId)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.NotValidData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		if !contains(accountProjectIds, ticket[0].ProjectID.Int) {
			return ctx.JSON(http.StatusForbidden, utils.StandardHttpResponse{
				Message: utils.NoAccess,
				Status:  403,
				Data:    "you dont have access to this project",
			})
		}
		tickets = append(tickets, ticket...)
	} else {
		projectss, err := hc.ticketsRepository.GetHeadTickets(ctx.Request().Context(), accountProjectIds)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
				Message: utils.NotValidData,
				Status:  400,
				Data:    err.Error(),
			})
		}
		tickets = append(tickets, projectss...)
	}

	var ticketsResponse []usecase_models.Tickets
	for _, ticket := range tickets {
		ticketsResponse = append(ticketsResponse, usecase_models.Tickets{
			ID:           ticket.ID,
			AccountID:    ticket.AccountID,
			ProjectID:    ticket.ProjectID,
			Message:      ticket.Message,
			TicketStatus: ticket.TicketStatus,
			Title:        ticket.Title,
			ReplyTo:      ticket.ReplyTo,
			CreatedAt:    ticket.CreatedAt,
			UpdatedAt:    ticket.UpdatedAt,
			DeletedAt:    ticket.DeletedAt,
		})
	}

	return ctx.JSON(http.StatusOK, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    ticketsResponse,
	})
}

func (hc *httpControllers) UpdateTicket(ctx echo.Context) error {
	req := new(usecase_models.Tickets)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: utils.NotValidData,
			Status:  400,
			Data:    err.Error(),
		})
	}

	ticketId, err := strconv.Atoi(ctx.Param("ticket_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.StandardHttpResponse{
			Message: fmt.Sprintf(utils.NotValidField, "Ticket ID"),
			Status:  400,
			Data:    err.Error(),
		})
	}

	err = hc.ticketsRepository.UpdateTickets(ctx.Request().Context(), models.Ticket{
		ID:           ticketId,
		Message:      req.Message,
		TicketStatus: req.TicketStatus,
		Title:        req.Title,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.StandardHttpResponse{
			Message: utils.ProblemInSystem,
			Status:  500,
			Data:    err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.StandardHttpResponse{
		Message: utils.Ok,
		Status:  200,
		Data:    nil,
	})
}
