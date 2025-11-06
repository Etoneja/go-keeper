package main

import (
	"log"

	"github.com/etoneja/go-keeper/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := server.LoadCfg()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	app, err := server.NewApp(cfg)
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}
	defer app.Stop()

	if err := app.Run(":50051"); err != nil {
		log.Fatal("Failed to run app:", err)
	}
}
