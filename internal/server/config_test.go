package server

import (
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
				os.Unsetenv(envVar)
			} else {
				os.Setenv(envVar, value)
			}
		}
	}()

	t.Run("success", func(t *testing.T) {
		os.Setenv("GOKEEPER_DB_USER", "testuser")
		os.Setenv("GOKEEPER_DB_PASSWORD", "testpass")
		os.Setenv("GOKEEPER_DB_HOST", "localhost")
		os.Setenv("GOKEEPER_DB_PORT", "5432")
		os.Setenv("GOKEEPER_DB_NAME", "testdb")
		os.Setenv("GOKEEPER_JWT_SECRET", "jwtsecret123")

		cfg, err := LoadCfg()
		require.NoError(t, err)
		assert.Equal(t, "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable", cfg.DBURL)
		assert.Equal(t, "jwtsecret123", cfg.JWTSecret)
	})

	t.Run("missing jwt secret", func(t *testing.T) {
		os.Setenv("GOKEEPER_DB_USER", "testuser")
		os.Setenv("GOKEEPER_DB_PASSWORD", "testpass")
		os.Setenv("GOKEEPER_DB_HOST", "localhost")
		os.Setenv("GOKEEPER_DB_PORT", "5432")
		os.Setenv("GOKEEPER_DB_NAME", "testdb")
		os.Unsetenv("GOKEEPER_JWT_SECRET")

		cfg, err := LoadCfg()
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "GOKEEPER_JWT_SECRET environment variable is required")
	})

	t.Run("empty db values", func(t *testing.T) {
		os.Setenv("GOKEEPER_DB_USER", "")
		os.Setenv("GOKEEPER_DB_PASSWORD", "")
		os.Setenv("GOKEEPER_DB_HOST", "")
		os.Setenv("GOKEEPER_DB_PORT", "")
		os.Setenv("GOKEEPER_DB_NAME", "")
		os.Setenv("GOKEEPER_JWT_SECRET", "jwtsecret123")

		cfg, err := LoadCfg()
		require.NoError(t, err)
		assert.Equal(t, "postgres://:@:/?sslmode=disable", cfg.DBURL)
		assert.Equal(t, "jwtsecret123", cfg.JWTSecret)
	})
}
