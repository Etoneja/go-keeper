package ctl

import (
	"github.com/etoneja/go-keeper/internal/ctl/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiffSecrets(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		diff := diffSecrets(nil, nil)
		assert.Empty(t, diff.LocalOnly)
		assert.Empty(t, diff.RemoteOnly)
		assert.Empty(t, diff.Both)
	})

	t.Run("local only", func(t *testing.T) {
		local := []*types.LocalSecret{
			{UUID: "1", Name: "local1"},
			{UUID: "2", Name: "local2"},
		}

		diff := diffSecrets(local, nil)
		assert.Len(t, diff.LocalOnly, 2)
		assert.Empty(t, diff.RemoteOnly)
		assert.Empty(t, diff.Both)
	})

	t.Run("remote only", func(t *testing.T) {
		remote := []*types.RemoteSecret{
			{UUID: "1"},
			{UUID: "2"},
		}

		diff := diffSecrets(nil, remote)
		assert.Empty(t, diff.LocalOnly)
		assert.Len(t, diff.RemoteOnly, 2)
		assert.Empty(t, diff.Both)
	})

	t.Run("both sides", func(t *testing.T) {
		local := []*types.LocalSecret{
			{UUID: "1", Name: "secret1"},
			{UUID: "2", Name: "secret2"},
		}
		remote := []*types.RemoteSecret{
			{UUID: "1"},
			{UUID: "2"},
		}

		diff := diffSecrets(local, remote)
		assert.Empty(t, diff.LocalOnly)
		assert.Empty(t, diff.RemoteOnly)
		assert.Len(t, diff.Both, 2)
	})

	t.Run("mixed", func(t *testing.T) {
		local := []*types.LocalSecret{
			{UUID: "1", Name: "local-only"},
			{UUID: "2", Name: "both-secret"},
		}
		remote := []*types.RemoteSecret{
			{UUID: "2"},
			{UUID: "3"},
		}

		diff := diffSecrets(local, remote)
		assert.Len(t, diff.LocalOnly, 1)
		assert.Equal(t, "1", diff.LocalOnly[0].UUID)
		assert.Len(t, diff.RemoteOnly, 1)
		assert.Equal(t, "3", diff.RemoteOnly[0].UUID)
		assert.Len(t, diff.Both, 1)
		assert.Equal(t, "2", diff.Both[0].Local.UUID)
		assert.Equal(t, "2", diff.Both[0].Remote.UUID)
	})
}
