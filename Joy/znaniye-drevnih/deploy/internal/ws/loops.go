package ws

import (
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func readLoop(conn *websocket.Conn, incoming chan<- *clientMessage, errs chan<- error, done <-chan struct{}, logger *zap.Logger) {
	for {
		select {
		case <-done:
			return
		default:
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			errs <- err
			return
		}

		parsed, err := parseClientMessage(msg)
		if err != nil {
			if logger != nil {
				logger.Warn("ws parse failed", zap.Error(err))
			}
			continue
		}

		incoming <- parsed
	}
}

func writeLoop(conn *websocket.Conn, outgoing <-chan any, errs chan<- error, done <-chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case msg, ok := <-outgoing:
			if !ok {
				return
			}
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(msg); err != nil {
				errs <- err
				return
			}
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				errs <- err
				return
			}
		}
	}
}
