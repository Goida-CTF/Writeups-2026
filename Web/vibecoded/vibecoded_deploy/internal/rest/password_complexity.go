package rest

import (
	"fmt"
	"net/http"

	"vibecoded/internal/models"
	"vibecoded/internal/server"
)

func (s *Service) CheckPasswordComplexity(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordComplexityRequest
	if !server.DecodeJSON(w, r, &req) {
		return
	}

	if !server.CheckAllFieldsExist(w, r, &req) {
		return
	}

	data, err := s.uc.CheckPasswordComplexity(req.Password)
	if err != nil {
		server.HandleInternalServerError(w, r,
			fmt.Errorf("s.uc.CheckPasswordComplexity: %w", err))
		return
	}

	res := &models.JsonResponse{
		Status: models.ResponseStatusOK,
		Data:   data,
	}
	server.RenderJSON(w, r, res)
}
