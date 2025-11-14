package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour)

	t.Run("success", func(t *testing.T) {
		token, err := manager.GenerateToken("user123")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("different users different tokens", func(t *testing.T) {
		token1, _ := manager.GenerateToken("user1")
		token2, _ := manager.GenerateToken("user2")
		assert.NotEqual(t, token1, token2)
	})
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Hour)

	t.Run("valid token", func(t *testing.T) {
		token, _ := manager.GenerateToken("user123")
		userID, err := manager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, "user123", userID)
	})

	t.Run("invalid signature", func(t *testing.T) {
		token, _ := manager.GenerateToken("user123")
		wrongManager := NewJWTManager("wrong-secret", time.Hour)
		userID, err := wrongManager.ValidateToken(token)
		require.Error(t, err)
		assert.Equal(t, "", userID)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("malformed token", func(t *testing.T) {
		userID, err := manager.ValidateToken("malformed.token.here")
		require.Error(t, err)
		assert.Equal(t, "", userID)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("empty token", func(t *testing.T) {
		userID, err := manager.ValidateToken("")
		require.Error(t, err)
		assert.Equal(t, "", userID)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("expired token", func(t *testing.T) {
		shortManager := NewJWTManager("test-secret", time.Millisecond)
		token, _ := shortManager.GenerateToken("user123")
		time.Sleep(10 * time.Millisecond)
		userID, err := manager.ValidateToken(token)
		require.Error(t, err)
		assert.Equal(t, "", userID)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})
}

func TestJWTManager_Integration(t *testing.T) {
	manager := NewJWTManager("integration-secret", time.Hour)

	users := []string{"user1", "user2", "user3"}
	tokens := make([]string, len(users))

	for i, userID := range users {
		token, err := manager.GenerateToken(userID)
		require.NoError(t, err)
		tokens[i] = token
	}

	for i, token := range tokens {
		userID, err := manager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, users[i], userID)
	}
}
