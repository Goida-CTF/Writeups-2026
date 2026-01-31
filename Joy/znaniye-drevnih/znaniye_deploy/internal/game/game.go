package game

import (
	"errors"
	"fmt"
	"time"

	"znanie-drevnih/internal/game/session"
	"znanie-drevnih/internal/models/taskmodels"
)

type Game struct {
	gameData *taskmodels.GameData

	taskPartsRequired uint64
	taskPartTimeout   time.Duration
	taskFlag          string
}

func New(
	taskDataPath string,
	taskPartsRequired uint64,
	taskPartTimeout time.Duration,
	taskFlag string,
) (*Game, error) {
	gameData, err := loadGameData(taskDataPath)
	if err != nil {
		return nil, fmt.Errorf("loadGameData: %w", err)
	}

	return &Game{
		gameData: gameData,

		taskPartsRequired: taskPartsRequired,
		taskPartTimeout:   taskPartTimeout,
		taskFlag:          taskFlag,
	}, nil
}

var errGameNotInit = errors.New("game not initialized")

func (g *Game) StartSession() (*session.Session, error) {
	if g == nil || g.gameData == nil {
		return nil, errGameNotInit
	}
	return session.New(g.gameData, g.taskPartsRequired), nil
}

func (g *Game) TaskPartTimeout() time.Duration {
	return g.taskPartTimeout
}

func (g *Game) TaskFlag() string {
	return g.taskFlag
}
