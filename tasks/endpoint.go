package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"test-manager/handlers"
	"test-manager/monitoring"
	"test-manager/repos"
	"test-manager/usecase_models"
	"time"
)

type EndpointTaskHandler struct {
	EndpointHandler handlers.EndpointHandler

	Logger *zap.SugaredLogger
}

func NewEndpointTaskHandler(
	EndpointHandler handlers.EndpointHandler,
	Logger *zap.SugaredLogger,
) *EndpointTaskHandler {
	return &EndpointTaskHandler{
		EndpointHandler: EndpointHandler,
		Logger:          Logger,
	}
}

func (c *EndpointTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload usecase_models.Endpoints
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed on endpoint task: %v: %w", err, asynq.SkipRetry)
	}

	start := time.Now()
	err := c.EndpointHandler.ExecuteEndpointRule(ctx, payload)
	if err != nil {
		log.Errorf("executing rule on endpoint task: %v", err)
	}

	monitoring.EndpointTaskDuration.Observe(-time.Until(start).Seconds())

	c.Logger.Info("success on processing endpoint task")
	return nil
}

type EndpointStoreHandler struct {
	endpointStatsRepo repos.EndpointStatsRepository
	Logger            *zap.SugaredLogger
}

func NewEndpointStoreHandler(
	endpointStatsRepo repos.EndpointStatsRepository,
	Logger *zap.SugaredLogger,
) *EndpointStoreHandler {
	return &EndpointStoreHandler{
		endpointStatsRepo: endpointStatsRepo,
		Logger:            Logger,
	}
}

func (c *EndpointStoreHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload []repos.WriteEndpointStatsOptions
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed on endpoint store: %v: %w", err, asynq.SkipRetry)
	}

	start := time.Now()

	err := c.endpointStatsRepo.WriteBulk(ctx, payload)
	if err != nil {
		return fmt.Errorf("problem on writing bulk endpoint stat: %s", err.Error())
	}

	monitoring.EndpointStoreDuration.Observe(-time.Until(start).Seconds())

	c.Logger.Info("success on processing endpoint store task")
	return nil
}
