package client

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	piston "github.com/milindmadhukar/go-piston"
	"go.uber.org/zap"
)

const bxxLibName = "Ве_крест_крест.h"

type Client struct {
	c                 *piston.Client
	l                 *zap.Logger
	pistonMemoryLimit uint64
	bxxLib            *piston.Code
}

func New(
	pistonBaseURL, pistonAPIKey string,
	pistonAPITimeout time.Duration,
	pistonMemoryLimit uint64,
	taskDataPath string,
	logger *zap.Logger,
) (*Client, error) {
	c := piston.New(pistonAPIKey,
		&http.Client{Timeout: pistonAPITimeout},
		pistonBaseURL)

	var bxxLibPath = filepath.Join(taskDataPath, bxxLibName)

	bxxLibContent, err := os.ReadFile(bxxLibPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open bxx lib: %w", err)
	}

	return &Client{
		c:                 c,
		l:                 logger,
		pistonMemoryLimit: pistonMemoryLimit,
		bxxLib: &piston.Code{
			Name:    bxxLibName,
			Content: string(bxxLibContent),
		},
	}, nil
}
