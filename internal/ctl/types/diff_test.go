package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecretCheckPair_IsIdentical(t *testing.T) {
	now := time.Now()

	t.Run("identical secrets", func(t *testing.T) {
		pair := &SecretCheckPair{
			Local: &LocalSecret{
				Hash:         "abc123",
				LastModified: now,
			},
			Remote: &RemoteSecret{
				Hash:         "abc123",
				LastModified: now,
			},
		}
		assert.True(t, pair.IsIdentical())
	})

	t.Run("different hash", func(t *testing.T) {
		pair := &SecretCheckPair{
			Local: &LocalSecret{
				Hash:         "abc123",
				LastModified: now,
			},
			Remote: &RemoteSecret{
				Hash:         "def456",
				LastModified: now,
			},
		}
		assert.False(t, pair.IsIdentical())
	})

	t.Run("different last modified", func(t *testing.T) {
		pair := &SecretCheckPair{
			Local: &LocalSecret{
				Hash:         "abc123",
				LastModified: now,
			},
			Remote: &RemoteSecret{
				Hash:         "abc123",
				LastModified: now.Add(time.Hour),
			},
		}
		assert.False(t, pair.IsIdentical())
	})

	t.Run("both different", func(t *testing.T) {
		pair := &SecretCheckPair{
			Local: &LocalSecret{
				Hash:         "abc123",
				LastModified: now,
			},
			Remote: &RemoteSecret{
				Hash:         "def456",
				LastModified: now.Add(time.Hour),
			},
		}
		assert.False(t, pair.IsIdentical())
	})
}

func TestSecretsDiff_Empty(t *testing.T) {
	diff := &SecretsDiff{}
	assert.Empty(t, diff.LocalOnly)
	assert.Empty(t, diff.RemoteOnly)
	assert.Empty(t, diff.Both)
}
