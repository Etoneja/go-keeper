package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadCfg(t *testing.T) {
	envVars := []string{
		"GOKEEPER_DB_PATH",
		"GOKEEPER_LOGIN",
		"GOKEEPER_PASSWORD",
		"GOKEEPER_SERVER_ADDR",
	}

	originalEnv := make(map[string]string, len(envVars))
	for _, envVar := range envVars {
		originalEnv[envVar] = os.Getenv(envVar)
	}

	defer func() {
		for envVar, value := range originalEnv {
			if value == "" {
				if err := os.Unsetenv(envVar); err != nil {
					t.Logf("failed to unset env var %s: %v", envVar, err)
				}
			} else {
				if err := os.Setenv(envVar, value); err != nil {
					t.Logf("failed to set env var %s: %v", envVar, err)
				}
			}
		}
	}()

	t.Run("success", func(t *testing.T) {
		envValues := map[string]string{
			"GOKEEPER_DB_PATH":     "/test/db",
			"GOKEEPER_LOGIN":       "testuser",
			"GOKEEPER_PASSWORD":    "testpass",
			"GOKEEPER_SERVER_ADDR": "localhost:8080",
		}

		for k, v := range envValues {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}

		cfg, err := LoadCfg()
		require.NoError(t, err)
		assert.Equal(t, "/test/db", cfg.DBPath)
		assert.Equal(t, "testuser", cfg.Login)
		assert.Equal(t, "testpass", cfg.Password)
		assert.Equal(t, "localhost:8080", cfg.ServerAddress)
	})

	t.Run("missing db path", func(t *testing.T) {
		envValues := map[string]string{
			"GOKEEPER_LOGIN":       "testuser",
			"GOKEEPER_PASSWORD":    "testpass",
			"GOKEEPER_SERVER_ADDR": "localhost:8080",
		}

		err := os.Unsetenv("GOKEEPER_DB_PATH")
		require.NoError(t, err)

		for k, v := range envValues {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}

		cfg, err := LoadCfg()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "GOKEEPER_DB_PATH")
	})

	t.Run("missing login", func(t *testing.T) {
		envValues := map[string]string{
			"GOKEEPER_DB_PATH":     "/test/db",
			"GOKEEPER_PASSWORD":    "testpass",
			"GOKEEPER_SERVER_ADDR": "localhost:8080",
		}

		err := os.Unsetenv("GOKEEPER_LOGIN")
		require.NoError(t, err)

		for k, v := range envValues {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}

		cfg, err := LoadCfg()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "GOKEEPER_LOGIN")
	})

	t.Run("missing password", func(t *testing.T) {
		envValues := map[string]string{
			"GOKEEPER_DB_PATH":     "/test/db",
			"GOKEEPER_LOGIN":       "testuser",
			"GOKEEPER_SERVER_ADDR": "localhost:8080",
		}

		err := os.Unsetenv("GOKEEPER_PASSWORD")
		require.NoError(t, err)

		for k, v := range envValues {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}

		cfg, err := LoadCfg()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "GOKEEPER_PASSWORD")
	})

	t.Run("missing server address", func(t *testing.T) {
		envValues := map[string]string{
			"GOKEEPER_DB_PATH":  "/test/db",
			"GOKEEPER_LOGIN":    "testuser",
			"GOKEEPER_PASSWORD": "testpass",
		}

		err := os.Unsetenv("GOKEEPER_SERVER_ADDR")
		require.NoError(t, err)

		for k, v := range envValues {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}

		cfg, err := LoadCfg()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "GOKEEPER_SERVER_ADDR")
	})
}
