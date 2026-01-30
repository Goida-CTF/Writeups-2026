package config

import (
	"fmt"
	"time"
)

type Config struct {
	PistonBaseURL     string        `envconfig:"PISTON_SERVER_URL" default:"http://piston:2000/api/v2"`
	PistonAPIKey      string        `envconfig:"PISTON_API_KEY" default:""`
	PistonAPITimeout  time.Duration `envconfig:"PISTON_API_TIMEOUT" default:"5m"`
	PistonMemoryLimit uint64        `envconfig:"PISTON_MEMORY_LIMIT" default:"4194304"`
	Host              string        `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	Port              uint16        `envconfig:"HTTP_PORT" default:"8080"`
	TaskDataPath      string        `envconfig:"TASK_DATA_PATH" default:"./data/"`
	// Zero value uses all parts in config
	TaskPartsRequired uint64        `envconfig:"TASK_PARTS_REQUIRED" default:"0"`
	TaskPartTimeout   time.Duration `envconfig:"TASK_PART_TIMEOUT" default:"30s"`
	TaskFlag          string        `envconfig:"FLAG" required:"true"`
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
