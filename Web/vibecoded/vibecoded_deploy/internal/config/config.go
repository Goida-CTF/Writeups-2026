package config

import "fmt"

type Config struct {
	LogLevel       string `envconfig:"LOG_LEVEL" default:"ERROR"`
	Host           string `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	Port           uint16 `envconfig:"HTTP_PORT" default:"8080"`
	DBPath         string `envconfig:"DB_PATH" default:"data.db"`
	HashPepper     []byte `envconfig:"HASH_PEPPER"`
	JWTSecret      []byte `envconfig:"JWT_SECRET"`
	HCaptchaSecret string `envconfig:"HCAPTCHA_SECRET"`
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
