package types

import (
	"encoding/json"
	"testing"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretDataContainer_MarshalJSON(t *testing.T) {
	t.Run("marshal and unmarshal text data", func(t *testing.T) {
		container := &SecretDataContainer{
			Type: constants.SecretTypeText,
			Name: "test secret",
			SecretData: TextData{
				Content: "secret content",
			},
		}

		data, err := json.Marshal(container)
		require.NoError(t, err)

		var unmarshaled SecretDataContainer
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, container.Type, unmarshaled.Type)
		assert.Equal(t, container.Name, unmarshaled.Name)

		textData, ok := unmarshaled.SecretData.(TextData)
		require.True(t, ok)
		assert.Equal(t, "secret content", textData.Content)
	})

	t.Run("marshal and unmarshal file data", func(t *testing.T) {
		container := &SecretDataContainer{
			Type: constants.SecretTypeBinary,
			Name: "test file",
			SecretData: FileData{
				FileName: "document.pdf",
				FileSize: 1024,
				Content:  "base64content",
			},
		}

		data, err := json.Marshal(container)
		require.NoError(t, err)

		var unmarshaled SecretDataContainer
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, container.Type, unmarshaled.Type)

		fileData, ok := unmarshaled.SecretData.(FileData)
		require.True(t, ok)
		assert.Equal(t, "document.pdf", fileData.FileName)
		assert.Equal(t, int64(1024), fileData.FileSize)
	})

	t.Run("unmarshal invalid type", func(t *testing.T) {
		data := []byte(`{"type": "invalid", "name": "test", "secret_data": {}}`)

		var container SecretDataContainer
		err := json.Unmarshal(data, &container)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown secret type")
	})

	t.Run("unmarshal invalid secret_data format", func(t *testing.T) {
		data := []byte(`{"type": "text", "name": "test", "secret_data": "invalid"}`)

		var container SecretDataContainer
		err := json.Unmarshal(data, &container)
		require.Error(t, err)
	})
}
