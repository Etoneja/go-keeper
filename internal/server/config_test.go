package server

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadCfg(t *testing.T) {
	originalEnv := map[string]string{}
	envVars := []string{
		"GOKEEPER_DB_USER",
		"GOKEEPER_DB_PASSWORD",
		"GOKEEPER_DB_HOST",
		"GOKEEPER_DB_PORT",
		"GOKEEPER_DB_NAME",
		"GOKEEPER_JWT_SECRET",
	}

	for _, envVar := range envVars {
		originalEnv[envVar] = os.Getenv(envVar)
	}

	defer func() {
		for envVar, value := range originalEnv {
			if value == "" {
				if err := os.Unsetenv(envVar); err != nil {
					log.Printf("Error unsetting env var %s: %v", envVar, err)
				}
			} else {
				if err := os.Setenv(envVar, value); err != nil {
					log.Printf("Error setting env var %s: %v", envVar, err)
				}
			}
		}
	}()

	t.Run("success", func(t *testing.T) {
		envVars := map[string]string{
			"GOKEEPER_DB_USER":     "testuser",
			"GOKEEPER_DB_PASSWORD": "testpass",
			"GOKEEPER_DB_HOST":     "localhost",
			"GOKEEPER_DB_PORT":     "5432",
			"GOKEEPER_DB_NAME":     "testdb",
			"GOKEEPER_JWT_SECRET":  "jwtsecret123",
		}

		for envVar, value := range envVars {
			if err := os.Setenv(envVar, value); err != nil {
				log.Printf("Error setting env var %s: %v", envVar, err)
			}
		}

		cfg, err := LoadCfg()
		require.NoError(t, err)
		assert.Equal(t, "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable", cfg.DBURL)
		assert.Equal(t, "jwtsecret123", cfg.JWTSecret)
	})

	t.Run("missing jwt secret", func(t *testing.T) {
		envVars := map[string]string{
			"GOKEEPER_DB_USER":     "testuser",
			"GOKEEPER_DB_PASSWORD": "testpass",
			"GOKEEPER_DB_HOST":     "localhost",
			"GOKEEPER_DB_PORT":     "5432",
			"GOKEEPER_DB_NAME":     "testdb",
		}

		for envVar, value := range envVars {
			if err := os.Setenv(envVar, value); err != nil {
				log.Printf("Error setting env var %s: %v", envVar, err)
			}
		}

		if err := os.Unsetenv("GOKEEPER_JWT_SECRET"); err != nil {
			log.Printf("Error unsetting env var GOKEEPER_JWT_SECRET: %v", err)
		}

		cfg, err := LoadCfg()
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "GOKEEPER_JWT_SECRET environment variable is required")
	})

	t.Run("empty db values", func(t *testing.T) {
		envVars := map[string]string{
			"GOKEEPER_DB_USER":     "",
			"GOKEEPER_DB_PASSWORD": "",
			"GOKEEPER_DB_HOST":     "",
			"GOKEEPER_DB_PORT":     "",
			"GOKEEPER_DB_NAME":     "",
			"GOKEEPER_JWT_SECRET":  "jwtsecret123",
		}

		for envVar, value := range envVars {
			if err := os.Setenv(envVar, value); err != nil {
				log.Printf("Error setting env var %s: %v", envVar, err)
			}
		}
		cfg, err := LoadCfg()
		require.NoError(t, err)
		assert.Equal(t, "postgres://:@:/?sslmode=disable", cfg.DBURL)
		assert.Equal(t, "jwtsecret123", cfg.JWTSecret)
	})
}
