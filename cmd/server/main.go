package main

import (
	"fmt"

	"github.com/etoneja/go-keeper/internal/server"
)

func main() {
	fmt.Println("Starting my-go-app...")
	server.Start()
}
