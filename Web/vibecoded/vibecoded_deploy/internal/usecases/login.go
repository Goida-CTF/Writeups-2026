package usecases

import (
	"fmt"

	"vibecoded/internal/models"
)

func (u *UseCases) Login(form *models.LoginForm) (*models.UserDBItem, string, error) {
	user, err := u.repo.GetUserDBItemByUsername(form.Username)
	if err != nil {
		return nil, "", fmt.Errorf("u.repo.GetUserDBItemByUsername %w", err)
	}
	if user == nil {
		return nil, "", ErrWrongUsername
	}

	hashedPassword := u.hasher.HashPassword(form.Password)
	user, err = u.repo.Login(form.Username, hashedPassword)
	if err != nil {
		return nil, "", fmt.Errorf("u.repo.Login: %w", err)
	}
	if user == nil {
		return nil, "", ErrWrongPassword
	}

	tokenString, err := u.jwt.NewJWTToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("u.jwt.NewJWTToken: %w", err)
	}

	return user, tokenString, nil
}
