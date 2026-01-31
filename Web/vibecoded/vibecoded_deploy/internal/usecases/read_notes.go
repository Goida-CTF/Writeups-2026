package usecases

import (
	"fmt"

	"vibecoded/internal/models"
)

func (u *UseCases) ReadNotes(userID int) ([]*models.NoteDBItem, error) {
	notes, err := u.repo.GetNoteDBItemsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("u.repo.GetNoteDBItemsByUserID: %w", err)
	}

	return notes, nil
}
