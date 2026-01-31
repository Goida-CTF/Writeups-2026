package usecases

import (
	"vibecoded/internal/auth"
	"vibecoded/internal/config"
)

type UseCases struct {
	repo   Repo
	cfg    *config.Config
	jwt    *auth.JWTProvider
	hasher *auth.PasswordHasher
}

func NewUseCases(repo Repo, config *config.Config,
	jwt *auth.JWTProvider, hasher *auth.PasswordHasher) *UseCases {
	return &UseCases{
		repo:   repo,
		cfg:    config,
		jwt:    jwt,
		hasher: hasher,
	}
}
