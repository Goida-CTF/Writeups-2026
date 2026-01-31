package usecases

import (
	"fmt"

	"vibecoded/internal/models"
)

func (u *UseCases) GetNoteDBItemsByUserID(userID int) ([]*models.NoteDBItem, error) {
	noteDBItems, err := u.repo.GetNoteDBItemsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("u.repo.GetNoteDBItemsByUserID: %w", err)
	}
	return noteDBItems, nil
}
