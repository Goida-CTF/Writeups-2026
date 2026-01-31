package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"

	"vibecoded/internal/auth"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugln(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				HandleInternalServerError(w, r,
					fmt.Errorf("panic recovered: %v; %s",
						err, string(debug.Stack())))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type JWTMiddleware struct {
	jwtProvider auth.JWTProvider
}

func NewJWTMiddleware(jwtProvider *auth.JWTProvider) *JWTMiddleware {
	return &JWTMiddleware{
		jwtProvider: *jwtProvider,
	}
}

func (m *JWTMiddleware) Handle(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.jwtProvider.ParseJWTToken(r)
		if err != nil {
			handleJWTError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyJWTToken, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type CaptchaMiddleware struct {
	captchaProvider auth.CaptchaProvider
}

func NewCaptchaMiddleware(captchaProvider *auth.CaptchaProvider) *CaptchaMiddleware {
	return &CaptchaMiddleware{
		captchaProvider: *captchaProvider,
	}
}

func (m *CaptchaMiddleware) Handle(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verified, err := m.captchaProvider.HandleVerification(r)
		if err != nil {
			handleCaptchaError(w, r, err)
			return
		}

		if !verified {
			handleCaptchaError(w, r, auth.ErrCaptchaTokenFailed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
