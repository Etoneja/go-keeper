package main

import (
	"log"

	"github.com/etoneja/go-keeper/internal/ctl"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	app, err := ctl.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.GetCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
