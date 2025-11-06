package server

import (
	"fmt"
	"os"
)

type Config struct {
	DBURL     string
	JWTSecret string
}

func LoadCfg() (*Config, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("GOKEEPER_DB_USER"),
		os.Getenv("GOKEEPER_DB_PASSWORD"),
		os.Getenv("GOKEEPER_DB_HOST"),
		os.Getenv("GOKEEPER_DB_PORT"),
		os.Getenv("GOKEEPER_DB_NAME"),
	)

	jwtSecret := os.Getenv("GOKEEPER_JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("GOKEEPER_JWT_SECRET environment variable is required")
	}

	return &Config{
		DBURL:     dbURL,
		JWTSecret: jwtSecret,
	}, nil
}
