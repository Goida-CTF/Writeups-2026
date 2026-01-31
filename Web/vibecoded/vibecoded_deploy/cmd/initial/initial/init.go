package initial

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"vibecoded/internal/models"
	"vibecoded/internal/usecases"
)

type Init struct {
	uc *usecases.UseCases
}

func newInit(uc *usecases.UseCases) *Init {
	return &Init{
		uc: uc,
	}
}

func (in *Init) createInitialUser(newUser *InitialUser) error {
	user, err := in.uc.RegisterUser(&models.NewUser{
		User: models.User{
			Username: newUser.Username,
			IsAdmin:  newUser.IsAdmin,
		},
		Password: newUser.Password,
	})
	if err != nil {
		if errors.Is(err, usecases.ErrUsernameAlreadyExists) {
			return fmt.Errorf("failed to register initial user \"%s\": user already exists", newUser.Username)
		}
		return fmt.Errorf("in.uc.RegisterUser: %w", err)
	}

	log.Infof("Registered initial user \"%s\" with id %d", user.Username, user.ID)

	for _, note := range newUser.Notes[:] {
		_, err := in.uc.NewNote(&models.NoteInternal{
			Note: models.Note{
				Title:   note.Title,
				Content: note.Content,
			},
			UserID: user.ID,
		})
		if err != nil {
			return fmt.Errorf("in.uc.NewNote: %w", err)
		}

		log.Debugf("Created initial note \"%s\" for user with id %d", note.Title, user.ID)
	}
	return nil
}
