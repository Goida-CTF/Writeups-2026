package server

import (
	"net/http"

	"znanie-drevnih/internal/ws"
)

func New(s *ws.Service, listenAddr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.HandleWS)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
		// WS server, so timeouts are not set there!
	}
	return srv
}
