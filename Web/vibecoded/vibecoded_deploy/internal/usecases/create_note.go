package usecases

import (
	"fmt"

	"vibecoded/internal/models"
)

func (u *UseCases) NewNote(newNote *models.NoteInternal) (*models.NoteDBItem, error) {
	note, err := u.repo.NewNote(newNote)
	if err != nil {
		return nil, fmt.Errorf("u.repo.NewNote: %w", err)
	}

	return note, nil
}
