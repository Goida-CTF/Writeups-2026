package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	"vibecoded/internal/auth"
	"vibecoded/internal/server"
)

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(server.CtxKeyJWTToken).(*jwt.Token)
	if !ok {
		server.HandleClientError(w, r,
			auth.ErrMissingAuthCookie,
			auth.ErrMissingAuthCookie.Error(), http.StatusUnauthorized)
		return
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		server.HandleClientError(w, r,
			fmt.Errorf("%w: %v", err, auth.ErrInvalidToken),
			auth.ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	log.Debugf("User with ID %s logged out", userIDString)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
}
