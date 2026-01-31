package rest

import (
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"vibecoded/internal/auth"
	"vibecoded/internal/models"
	"vibecoded/internal/server"
	"vibecoded/internal/usecases"
)

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginForm
	if !server.DecodeJSON(w, r, &req) {
		return
	}
	req.Username = strings.TrimSpace(req.Username)

	if !server.CheckAllFieldsExist(w, r, &req) {
		return
	}

	user, token, err := s.uc.Login(&req)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrWrongUsername) ||
			errors.Is(err, usecases.ErrWrongPassword):
			server.HandleClientError(w, r,
				err,
				err.Error(), http.StatusOK)
		default:
			server.HandleInternalServerError(w, r, err)
		}
		return
	}

	s.jwt.SetJWTTokenCookie(w, auth.BearerTokenPrefix+token)
	log.Debugf("User \"%s\" logged in", user.Username)

	res := &models.JsonResponse{
		Status: models.ResponseStatusOK,
	}
	server.RenderJSON(w, r, res)
}
