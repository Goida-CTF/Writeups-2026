package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"vibecoded/internal/auth"
	"vibecoded/internal/server"
)

type Service struct {
	uc  UseCases
	jwt *auth.JWTProvider
}

func NewService(uc UseCases, jwt *auth.JWTProvider) *Service {
	return &Service{
		uc:  uc,
		jwt: jwt,
	}
}

type RESTCtrl struct {
	s      *Service
	router *mux.Router
}

func NewRESTCtrl(s *Service, r *mux.Router) *RESTCtrl {
	return &RESTCtrl{
		s:      s,
		router: r,
	}
}

func (s *RESTCtrl) RegisterMiddleware(m ...mux.MiddlewareFunc) {
	s.router.Use(m...)
}

func (s *RESTCtrl) RegisterRoutes(
	authMiddleware *server.JWTMiddleware,
	captchaMiddleware *server.CaptchaMiddleware) {
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/check-password-complexity", s.s.CheckPasswordComplexity).Methods("POST")
	apiRouter.HandleFunc("/login", captchaMiddleware.Handle(
		http.HandlerFunc(s.s.Login))).Methods("POST")
	apiRouter.HandleFunc("/register", captchaMiddleware.Handle(
		http.HandlerFunc(s.s.Register))).Methods("POST")
	apiRouter.HandleFunc("/posts", authMiddleware.Handle(
		http.HandlerFunc(s.s.GetPosts))).Methods("GET")
	apiRouter.HandleFunc("/logout", authMiddleware.Handle(
		http.HandlerFunc(s.s.Logout))).Methods("POST")
}
