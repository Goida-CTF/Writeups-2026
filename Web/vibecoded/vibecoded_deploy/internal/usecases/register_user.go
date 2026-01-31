package usecases

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"

	"vibecoded/internal/models"
)

func (u *UseCases) RegisterUser(newUser *models.NewUser) (*models.UserDBItem, error) {
	newUser.Username = strings.TrimSpace(newUser.Username)

	if err := u.CheckRegistrationPrerequirements(newUser); err != nil {
		return nil, fmt.Errorf("u.CheckRegistrationPrerequirements: %w", err)
	}

	user, err := u.repo.NewUser(&models.UserInternal{
		User: models.User{
			Username: newUser.Username,
			IsAdmin:  newUser.IsAdmin,
		},
		HashedPassword: u.hasher.HashPassword(newUser.Password),
	})
	if err != nil {
		e := errors.Unwrap(err)
		if e == nil {
			e = err
		}

		if sqliteErr, ok := e.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			switch {
			case strings.Contains(sqliteErr.Error(), "users.username"):
				return nil, ErrUsernameAlreadyExists
			case strings.Contains(sqliteErr.Error(), "users.password"):
				return nil, ErrPasswordCollision
			}
		}
		return nil, fmt.Errorf("u.repo.NewUser: %w", err)
	}

	return user, nil
}
