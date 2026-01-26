package main

import (
	"fmt"
	"log"
	"net/http"

	"oreshnik/internal/app/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", s.Router()))
}
