package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBPath        string
	Login         string
	Password      string
	ServerAddress string
}

func LoadCfg() (*Config, error) {
	dbPath := os.Getenv("GOKEEPER_DB_PATH")
	if dbPath == "" {
		return nil, fmt.Errorf("GOKEEPER_DB_PATH environment variable is required")
	}

	login := os.Getenv("GOKEEPER_LOGIN")
	if login == "" {
		return nil, fmt.Errorf("GOKEEPER_LOGIN environment variable is required")
	}

	password := os.Getenv("GOKEEPER_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("GOKEEPER_PASSWORD environment variable is required")
	}

	serverAddres := os.Getenv("GOKEEPER_SERVER_ADDR")
	if serverAddres == "" {
		return nil, fmt.Errorf("GOKEEPER_SERVER_ADDR environment variable is required")
	}

	return &Config{
		DBPath:        dbPath,
		Login:         login,
		Password:      password,
		ServerAddress: serverAddres,
	}, nil
}
