package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"

	"vibecoded/internal/auth"
	"vibecoded/internal/models"
	"vibecoded/internal/server"
)

func (s *Service) GetPosts(w http.ResponseWriter, r *http.Request) {
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

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		server.HandleClientError(w, r,
			fmt.Errorf("%w: %v", err, auth.ErrInvalidToken),
			auth.ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	noteDBItems, err := s.uc.GetNoteDBItemsByUserID(userID)
	if err != nil {
		server.HandleInternalServerError(w, r, err)
		return
	}

	var notes []*models.Note
	for _, note := range noteDBItems[:] {
		notes = append(notes, &models.Note{
			Title:   note.Title,
			Content: note.Content,
		})
	}

	res := &models.JsonResponse{
		Status: models.ResponseStatusOK,
		Data:   notes,
	}
	server.RenderJSON(w, r, res)
}
