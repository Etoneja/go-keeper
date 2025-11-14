package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("GOKEEPER_DB_USER"),
		os.Getenv("GOKEEPER_DB_PASSWORD"),
		os.Getenv("GOKEEPER_DB_HOST"),
		os.Getenv("GOKEEPER_DB_PORT"),
		os.Getenv("GOKEEPER_DB_NAME"),
	)

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration up failed:", err)
		}
		log.Println("Migration up completed successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration down failed:", err)
		}
		log.Println("Migration down completed successfully")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				log.Println("No migrations applied yet")
				return
			}
			log.Fatal("Failed to get version:", err)
		}
		log.Printf("Current version: %d, dirty: %t", version, dirty)
	default:
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down]")
	}
}
