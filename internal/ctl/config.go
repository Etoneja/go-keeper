package ctl

import (
	"fmt"
	"os"
)

type Config struct {
	DBPath   string
	Login    string
	Password string
}

func LoadCfg() (*Config, error) {
	dbPath := os.Getenv("GOKEEPERCTL_DB_PATH")
	if dbPath == "" {
		return nil, fmt.Errorf("GOKEEPERCTL_DB_PATH environment variable is required")
	}

	login := os.Getenv("GOKEEPERCTL_LOGIN")
	if login == "" {
		return nil, fmt.Errorf("GOKEEPERCTL_LOGIN environment variable is required")
	}

	password := os.Getenv("GOKEEPERCTL_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("GOKEEPERCTL_PASSWORD environment variable is required")
	}

	return &Config{
		DBPath:   dbPath,
		Login:    login,
		Password: password,
	}, nil
}
