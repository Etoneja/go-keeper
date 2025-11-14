package main

import (
	"log"

	"github.com/etoneja/go-keeper/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	app, err := server.NewApp()
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}
	defer app.Stop()

	if err := app.Run(":50051"); err != nil {
		log.Fatal("Failed to run app:", err)
	}
}
