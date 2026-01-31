package config

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var ErrEmptyHCaptchaSecret = errors.New("empty HCAPTCHA_SECRET")

func generateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("rand.Read: %w", err)
	}
	return randomBytes, nil
}

func setRandomValueIfNotSet(ev *[]byte, keyName string, size int) error {
	if len(*ev) != 0 {
		return nil
	}

	log.Warnf("%s env variable is not set, generating value...", keyName)

	value, err := generateRandomBytes(size)
	if err != nil {
		return fmt.Errorf("generateRandomBytes: %w", err)
	}
	*ev = value

	return nil
}

func InitConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalln("log.ParseLevel: ", err)
	}
	log.SetLevel(level)

	if err := setRandomValueIfNotSet(
		&cfg.HashPepper, "HASH_PEPPER", 16); err != nil {
		return nil, fmt.Errorf("setRandomValueIfNotSet: %w", err)
	}
	if err := setRandomValueIfNotSet(
		&cfg.JWTSecret, "JWT_SECRET", 32); err != nil {
		return nil, fmt.Errorf("setRandomValueIfNotSet: %w", err)
	}

	if cfg.HCaptchaSecret == "" {
		return nil, ErrEmptyHCaptchaSecret
	}

	return &cfg, nil
}
