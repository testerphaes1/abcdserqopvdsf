package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"test-manager/handlers"
	"test-manager/usecase_models"
)

type TraceRouteTaskHandler struct {
	TraceRouteHandler handlers.TraceRouteHandler

	Logger *zap.SugaredLogger
}

func (c *TraceRouteTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload usecase_models.TraceRoutes
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed on TraceRoute task: %v: %w", err, asynq.SkipRetry)
	}

	err := c.TraceRouteHandler.ExecuteTraceRouteRule(ctx, payload)
	if err != nil {
		c.Logger.Info(err)
		return fmt.Errorf("executing rule on TraceRoute task: %v", err)
	}

	c.Logger.Info("success on processing trace route task")
	return nil
}

func NewTraceRouteTaskHandler(
	TraceRouteHandler handlers.TraceRouteHandler,
	Logger *zap.SugaredLogger,
) *TraceRouteTaskHandler {
	return &TraceRouteTaskHandler{
		TraceRouteHandler: TraceRouteHandler,
		Logger:            Logger,
	}
}
