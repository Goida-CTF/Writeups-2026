package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"vibecoded/internal/models"
	"vibecoded/internal/server"
	"vibecoded/internal/usecases"
)

const RegistrationIsClosed = `Registration of new users is closed for now`

func (s *Service) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegistrationForm
	if !server.DecodeJSON(w, r, &req) {
		return
	}
	req.Username = strings.TrimSpace(req.Username)

	if !server.CheckAllFieldsExist(w, r, &req) {
		return
	}

	if err := s.uc.CheckRegistrationPrerequirements(&models.NewUser{
		User: models.User{
			Username: req.Username,
			IsAdmin:  false,
		},
		Password: req.Password,
	}); err != nil {
		switch {
		case errors.Is(err, usecases.ErrUsernameAlreadyExists) ||
			errors.Is(err, usecases.ErrComplexityNotSatisfied):
			server.HandleClientError(w, r,
				err,
				err.Error(), http.StatusBadRequest)
			return
		}
		server.HandleInternalServerError(w, r, err)
		return
	}

	server.HandleClientError(w, r,
		fmt.Errorf("user tried to register by username \"%s\"", req.Username),
		RegistrationIsClosed, http.StatusForbidden)
}
