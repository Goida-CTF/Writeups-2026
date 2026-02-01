package rest

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"vibecoded/internal/models"
	"vibecoded/internal/server"
)

const textTooFast = "too many passwords sent in a meanwhile, timeout of 10 passwords per second reached"

type passwordComplexityLimiter struct {
	mu          sync.Mutex
	windowStart time.Time
	count       int
}

func (l *passwordComplexityLimiter) allow(now time.Time, limit int, window time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.windowStart.IsZero() || now.Sub(l.windowStart) >= window {
		l.windowStart = now
		l.count = 0
	}

	if l.count >= limit {
		return false
	}

	l.count++
	return true
}

var passwordComplexityRateLimiter passwordComplexityLimiter

func (s *Service) CheckPasswordComplexity(w http.ResponseWriter, r *http.Request) {
	var req models.PasswordComplexityRequest
	if !server.DecodeJSON(w, r, &req) {
		return
	}

	if !server.CheckAllFieldsExist(w, r, &req) {
		return
	}

	if !passwordComplexityRateLimiter.allow(time.Now(), 10, time.Second) {
		res := &models.JsonResponse{
			Status: models.ResponseStatusOK,
			Data: &models.PasswordComplexityResult{
				Ok:    false,
				Level: 1,
				Text:  textTooFast,
			},
		}
		server.RenderJSON(w, r, res)
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
