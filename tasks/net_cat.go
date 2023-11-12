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

type NetCatTaskHandler struct {
	NetCatHandler handlers.NetCatHandler

	Logger *zap.SugaredLogger
}

func (c *NetCatTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload usecase_models.NetCats
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed on NetCat task: %v: %w", err, asynq.SkipRetry)
	}

	err := c.NetCatHandler.ExecuteNetCatRule(ctx, payload)
	if err != nil {
		c.Logger.Info(err)
		return fmt.Errorf("executing rule on NetCat task: %v", err)
	}

	c.Logger.Info("success on processing net cat task")
	return nil
}

func NewNetCatTaskHandler(
	NetCatHandler handlers.NetCatHandler,
	Logger *zap.SugaredLogger,
) *NetCatTaskHandler {
	return &NetCatTaskHandler{
		NetCatHandler: NetCatHandler,
		Logger:        Logger,
	}
}
