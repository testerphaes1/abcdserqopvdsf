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

type PageSpeedTaskHandler struct {
	PageSpeedHandler handlers.PageSpeedHandler

	Logger *zap.SugaredLogger
}

func (c *PageSpeedTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload usecase_models.PageSpeeds
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed on PageSpeeds task: %v: %w", err, asynq.SkipRetry)
	}

	err := c.PageSpeedHandler.ExecutePageSpeedRule(ctx, payload)
	if err != nil {
		c.Logger.Info(err)
		return fmt.Errorf("executing rule on PageSpeeds task: %v", err)
	}

	c.Logger.Info("success on processing page speed task")
	return nil
}

func NewPageSpeedTaskHandler(
	PageSpeedHandler handlers.PageSpeedHandler,
	Logger *zap.SugaredLogger,
) *PageSpeedTaskHandler {
	return &PageSpeedTaskHandler{
		PageSpeedHandler: PageSpeedHandler,
		Logger:           Logger,
	}
}
