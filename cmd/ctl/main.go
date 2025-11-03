package main

import (
	"log"

	"github.com/etoneja/go-keeper/internal/ctl"
)

func main() {
	app, err := ctl.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.GetCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
