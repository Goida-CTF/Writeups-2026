package rest

import "vibecoded/internal/models"

type UseCases interface {
	UsersUseCases
	NotesUseCases
}

type UsersUseCases interface {
	CheckRegistrationPrerequirements(newUser *models.NewUser) error
	RegisterUser(newUser *models.NewUser) (*models.UserDBItem, error)
	Login(form *models.LoginForm) (*models.UserDBItem, string, error)
	CheckPasswordComplexity(password string) (*models.PasswordComplexityResult, error)
}

type NotesUseCases interface {
	NewNote(newNote *models.NoteInternal) (*models.NoteDBItem, error)
	ReadNotes(userID int) ([]*models.NoteDBItem, error)
	GetNoteDBItemsByUserID(userID int) ([]*models.NoteDBItem, error)
}
