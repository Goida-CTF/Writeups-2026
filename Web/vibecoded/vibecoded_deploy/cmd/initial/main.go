package main

import (
	log "github.com/sirupsen/logrus"

	"vibecoded/cmd/initial/initial"
	"vibecoded/internal/auth"
	"vibecoded/internal/config"
	"vibecoded/internal/database"
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

	uc := usecases.NewUseCases(
		database.NewRepo(db),
		cfg,
		jwt,
		auth.NewPasswordHasher(cfg.HashPepper),
	)

	initial.Initialize(uc)
}
