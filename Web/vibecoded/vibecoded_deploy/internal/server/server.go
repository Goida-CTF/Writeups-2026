package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"vibecoded/internal/config"
)

type Server struct {
	cfg *config.Config
	s   *http.Server
}

func NewServer(cfg *config.Config, router *mux.Router) *Server {
	return &Server{
		cfg: cfg,
		s: &http.Server{
			Addr:         cfg.GetAddr(),
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
	}
}

func (s *Server) ListenAndServe(addr string) error {
	s.s.Addr = addr
	return s.s.ListenAndServe()
}
