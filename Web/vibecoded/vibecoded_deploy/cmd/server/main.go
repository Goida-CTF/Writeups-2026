package main

import (
	log "github.com/sirupsen/logrus"

	"vibecoded/internal/auth"
	"vibecoded/internal/config"
	"vibecoded/internal/database"
	"vibecoded/internal/rest"
	"vibecoded/internal/server"
	"vibecoded/internal/usecases"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalln("config.InitConfig:", err)
	}

	db, err := database.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalln("database.Connect:", err)
	}
	if database.InitDB(db) != nil {
		log.Fatalln("database.InitDB:", err)
	}

	jwt := auth.NewJWTProvider(cfg.JWTSecret)
	captcha := auth.NewCaptchaProvider(cfg.HCaptchaSecret)

	uc := usecases.NewUseCases(
		database.NewRepo(db),
		cfg,
		jwt,
		auth.NewPasswordHasher(cfg.HashPepper),
	)

	router := server.NewRouter()

	ctrl := rest.NewRESTCtrl(rest.NewService(uc, jwt), router)
	ctrl.RegisterMiddleware(
		server.LoggingMiddleware,
		server.PanicMiddleware,
	)
	ctrl.RegisterRoutes(
		server.NewJWTMiddleware(jwt),
		server.NewCaptchaMiddleware(captcha))

	s := server.NewServer(cfg, router)

	log.Infof("server listening on http://%s", cfg.GetAddr())
	if err := s.ListenAndServe(cfg.GetAddr()); err != nil {
		log.Fatalln("s.ListenAndServe:", err)
	}
}
