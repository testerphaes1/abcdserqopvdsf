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

type PingTaskHandler struct {
	PingHandler handlers.PingHandler

	Logger *zap.SugaredLogger
}

func (c *PingTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload usecase_models.Pings
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed on Ping task: %v: %w", err, asynq.SkipRetry)
	}

	err := c.PingHandler.ExecutePingRule(ctx, payload)
	if err != nil {
		c.Logger.Info(err)
		return fmt.Errorf("executing rule on Ping task: %v", err)
	}

	c.Logger.Info("success on processing ping task")
	return nil
}

func NewPingTaskHandler(
	PingHandler handlers.PingHandler,
	Logger *zap.SugaredLogger,
) *PingTaskHandler {
	return &PingTaskHandler{
		PingHandler: PingHandler,
		Logger:      Logger,
	}
}
