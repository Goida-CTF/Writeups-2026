package usecases

import (
	"errors"
)

var (
	ErrComplexityNotSatisfied = errors.New("password complexity unsatisfied")
	ErrUsernameAlreadyExists  = errors.New("username already exists")
	ErrPasswordCollision      = errors.New("password collision")
	ErrWrongUsername          = errors.New("wrong username")
	ErrWrongPassword          = errors.New("wrong password")
)
