package usecases

import (
	"fmt"

	"vibecoded/internal/auth"
	"vibecoded/internal/models"
)

func (u *UseCases) CheckPasswordComplexity(password string,
) (*models.PasswordComplexityResult, error) {
	username, err := u.repo.FindPasswordCollision(
		u.hasher.HashPassword(password))
	if err != nil {
		return nil, fmt.Errorf("u.repo.FindPasswordCollision: %w", err)
	}

	return auth.CheckPasswordComplexity(password, username), nil
}
