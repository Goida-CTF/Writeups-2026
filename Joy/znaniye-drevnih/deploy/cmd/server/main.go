package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"znanie-drevnih/internal/client"
	"znanie-drevnih/internal/config"
	"znanie-drevnih/internal/logger"
	"znanie-drevnih/internal/server"
	"znanie-drevnih/internal/usecases"
	"znanie-drevnih/internal/ws"
)

func main() {
	var logLevel = zap.InfoLevel
	logger, err := logger.New(logLevel)
	if err != nil {
		log.Fatalf("logger.New: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("config failed to load", zap.Error(err))
	}

	client, err := client.New(cfg.PistonBaseURL,
		cfg.PistonAPIKey,
		cfg.PistonAPITimeout,
		cfg.PistonMemoryLimit,
		cfg.TaskDataPath,
		logger)
	if err != nil {
		logger.Fatal("piston client failed to create", zap.Error(err))
	}

	uc, err := usecases.New(logger, client,
		cfg.TaskDataPath,
		cfg.TaskPartsRequired,
		cfg.TaskPartTimeout,
		cfg.TaskFlag)
	if err != nil {
		logger.Fatal("usecases failed to create", zap.Error(err))
	}
	svc := ws.New(uc, logger)
	srv := server.New(svc, cfg.ListenAddr())

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	logger.Info("server started", zap.String("addr", srv.Addr))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	const shutdownTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", zap.Error(err))
	} else {
		logger.Info("server shutdown complete")
	}
}
