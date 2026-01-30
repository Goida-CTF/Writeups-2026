package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4 << 20 // 4 MiB
)

type Service struct {
	uc       UseCases
	l        *zap.Logger
	upgrader *websocket.Upgrader
}

func New(uc UseCases, l *zap.Logger) *Service {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return &Service{
		uc:       uc,
		l:        l,
		upgrader: upgrader,
	}
}
