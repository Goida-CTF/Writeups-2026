package usecases

import "vibecoded/internal/models"

type Repo interface {
	UsersRepo
	NotesRepo
}

type UsersRepo interface {
	FindPasswordCollision(hashedPassword string) (string, error)
	NewUser(user *models.UserInternal) (*models.UserDBItem, error)
	Login(username, hashedPassword string) (*models.UserDBItem, error)
	GetUserDBItemByUsername(username string) (*models.UserDBItem, error)
	GetUserDBItemByID(id int) (*models.UserDBItem, error)
}

type NotesRepo interface {
	NewNote(note *models.NoteInternal) (*models.NoteDBItem, error)
	UpdateNote(note *models.NoteUpdate) error
	GetNoteDBItemsByUserID(userID int) ([]*models.NoteDBItem, error)
}
