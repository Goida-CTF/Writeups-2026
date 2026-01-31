package ws

import (
	"context"
	"time"

	"znanie-drevnih/internal/game/session"
	"znanie-drevnih/internal/models"
	"znanie-drevnih/internal/models/taskmodels"
)

type UseCases interface {
	StartSession() (*session.Session, error)
	RunTask(ctx context.Context, task *taskmodels.Task, bxxCode string) (*models.TaskRunResult, error)
	TaskPartTimeout() time.Duration
	TaskFlag() string
}
