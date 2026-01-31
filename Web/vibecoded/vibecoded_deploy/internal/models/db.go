package models

import "time"

type dbItem struct {
	ID        int
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type User struct {
	Username string
	IsAdmin  bool
}

type UserInternal struct {
	User
	HashedPassword string
}

type UserDBItem struct {
	dbItem
	UserInternal
}

type NewUser struct {
	User
	Password string
}

type NoteInternal struct {
	Note
	UserID int
}

type NoteDBItem struct {
	dbItem
	NoteInternal
}
