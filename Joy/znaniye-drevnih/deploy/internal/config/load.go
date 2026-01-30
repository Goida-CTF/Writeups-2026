package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

func Load() (*Config, error) {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("envconfig.Process: %w", err)
	}
	cfg.PistonBaseURL = strings.Trim(cfg.PistonBaseURL, `'"/`) + "/"
	cfg.PistonAPIKey = strings.Trim(cfg.PistonAPIKey, `'"`)
	cfg.TaskDataPath = strings.Trim(cfg.TaskDataPath, `'"`)
	cfg.TaskFlag = strings.Trim(cfg.TaskFlag, `'"`)

	return cfg, nil
}
