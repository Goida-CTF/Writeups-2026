package usecases

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"znanie-drevnih/internal/client"
	"znanie-drevnih/internal/game"
)

type UseCases struct {
	logger *zap.Logger
	client *client.Client
	game   *game.Game
}

func New(
	l *zap.Logger,
	c *client.Client,
	taskDataPath string,
	taskPartsRequired uint64,
	taskPartTimeout time.Duration,
	taskFlag string,
) (*UseCases, error) {
	game, err := game.New(taskDataPath, taskPartsRequired, taskPartTimeout, taskFlag)
	if err != nil {
		return nil, fmt.Errorf("game.New: %w", err)
	}

	return &UseCases{
		logger: l,
		client: c,
		game:   game,
	}, nil
}
