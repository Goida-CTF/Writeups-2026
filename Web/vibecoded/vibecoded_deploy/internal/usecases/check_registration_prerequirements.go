package usecases

import (
	"fmt"

	"vibecoded/internal/models"
)

func (u *UseCases) CheckRegistrationPrerequirements(newUser *models.NewUser) error {
	user, err := u.repo.GetUserDBItemByUsername(newUser.Username)
	if err != nil {
		return fmt.Errorf("u.repo.GetUserDBItemByUsername: %w", err)
	}
	if user != nil {
		return ErrUsernameAlreadyExists
	}

	passwordComplexity, err := u.CheckPasswordComplexity(newUser.Password)
	if err != nil {
		return fmt.Errorf("u.CheckPasswordComplexity: %w", err)
	}
	if !passwordComplexity.Ok {
		return fmt.Errorf("%w: %s",
			ErrComplexityNotSatisfied, passwordComplexity.Text,
		)
	}

	return nil
}
