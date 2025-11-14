package types

import (
	"testing"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/etoneja/go-keeper/internal/ctl/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSecretModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCryptor := crypto.NewMockCryptor(ctrl)

	t.Run("successful creation", func(t *testing.T) {
		base := BaseSecret{
			Type:     constants.SecretTypeText,
			Name:     "test secret",
			Metadata: "metadata",
		}
		data := TextData{Content: "test content"}

		mockCryptor.EXPECT().
			CalculateDataHash(gomock.Any()).
			Return("hash123")

		secret, err := NewSecretModel(base, data, mockCryptor)
		require.NoError(t, err)

		assert.NotEmpty(t, secret.UUID)
		assert.Equal(t, base.Type, secret.Type)
		assert.Equal(t, base.Name, secret.Name)
		assert.Equal(t, base.Metadata, secret.Metadata)
		assert.Equal(t, "hash123", secret.Hash)
		assert.NotZero(t, secret.LastModified)
		assert.NotEmpty(t, secret.Data)
	})

	t.Run("empty type", func(t *testing.T) {
		base := BaseSecret{Name: "test", Type: ""}
		data := TextData{Content: "content"}
		_, err := NewSecretModel(base, data, mockCryptor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type is required")
	})

	t.Run("empty name", func(t *testing.T) {
		base := BaseSecret{Type: constants.SecretTypeText, Name: ""}
		data := TextData{Content: "content"}
		_, err := NewSecretModel(base, data, mockCryptor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("nil data", func(t *testing.T) {
		base := BaseSecret{Type: constants.SecretTypeText, Name: "test"}
		_, err := NewSecretModel(base, nil, mockCryptor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data is required")
	})

	t.Run("unsupported type", func(t *testing.T) {
		base := BaseSecret{Type: "invalid", Name: "test"}
		data := TextData{Content: "content"}
		_, err := NewSecretModel(base, data, mockCryptor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported secret type")
	})

	t.Run("type mismatch", func(t *testing.T) {
		base := BaseSecret{Type: constants.SecretTypeText, Name: "test"}
		data := FileData{FileName: "file.txt", Content: "content"}
		_, err := NewSecretModel(base, data, mockCryptor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type mismatch")
	})

	t.Run("invalid data", func(t *testing.T) {
		base := BaseSecret{Type: constants.SecretTypeText, Name: "test"}
		data := TextData{Content: ""}
		_, err := NewSecretModel(base, data, mockCryptor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data validation failed")
	})
}

func TestLocalSecret_SetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCryptor := crypto.NewMockCryptor(ctrl)

	t.Run("set data successfully", func(t *testing.T) {
		secret := &LocalSecret{
			Type:     constants.SecretTypeText,
			Metadata: "meta",
		}
		data := TextData{Content: "test content"}

		mockCryptor.EXPECT().
			CalculateDataHash(gomock.Any()).
			Return("hash123")

		err := secret.SetData(mockCryptor, data)
		require.NoError(t, err)

		assert.Equal(t, "hash123", secret.Hash)
		assert.NotEmpty(t, secret.Data)

		// ParseData должен работать после установки Type
		parsedData, err := secret.ParseData()
		require.NoError(t, err)
		assert.IsType(t, TextData{}, parsedData)
		assert.Equal(t, "test content", parsedData.(TextData).Content)
	})
}
