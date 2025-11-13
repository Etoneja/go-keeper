package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCryptorImpl(t *testing.T) {
	cryptor := NewCryptor("masterpass", "testuser")

	t.Run("EncryptStorageData and DecryptStorageData", func(t *testing.T) {
		plainData := []byte("test data for storage")

		encrypted, err := cryptor.EncryptStorageData(plainData)
		require.NoError(t, err)
		require.NotEmpty(t, encrypted)
		require.Greater(t, len(encrypted), storageSaltSize)

		decrypted, err := cryptor.DecryptStorageData(encrypted)
		require.NoError(t, err)
		assert.Equal(t, plainData, decrypted)
	})

	t.Run("DecryptStorageData invalid data", func(t *testing.T) {
		_, err := cryptor.DecryptStorageData([]byte("short"))
		require.Error(t, err)
	})

	t.Run("EncryptSecretData and DecryptSecretData", func(t *testing.T) {
		plainData := []byte("secret data")

		encrypted, err := cryptor.EncryptSecretData(plainData)
		require.NoError(t, err)
		require.NotEmpty(t, encrypted)

		decrypted, err := cryptor.DecryptSecretData(encrypted)
		require.NoError(t, err)
		assert.Equal(t, plainData, decrypted)
	})

	t.Run("CalculateDataHash", func(t *testing.T) {
		data := []byte("test data")
		hash1 := cryptor.CalculateDataHash(data)
		hash2 := cryptor.CalculateDataHash(data)

		assert.NotEmpty(t, hash1)
		assert.Equal(t, hash1, hash2)

		hash3 := cryptor.CalculateDataHash([]byte("different data"))
		assert.NotEqual(t, hash1, hash3)
	})

	t.Run("GenerateServerPassword", func(t *testing.T) {
		password1 := cryptor.GenerateServerPassword()
		password2 := cryptor.GenerateServerPassword()

		assert.NotEmpty(t, password1)
		assert.Equal(t, password1, password2)
	})

	t.Run("storage data different salts", func(t *testing.T) {
		plainData := []byte("same data")

		encrypted1, err := cryptor.EncryptStorageData(plainData)
		require.NoError(t, err)

		encrypted2, err := cryptor.EncryptStorageData(plainData)
		require.NoError(t, err)

		assert.NotEqual(t, encrypted1, encrypted2)

		decrypted1, err := cryptor.DecryptStorageData(encrypted1)
		require.NoError(t, err)
		assert.Equal(t, plainData, decrypted1)

		decrypted2, err := cryptor.DecryptStorageData(encrypted2)
		require.NoError(t, err)
		assert.Equal(t, plainData, decrypted2)
	})

	t.Run("secret data different nonces", func(t *testing.T) {
		plainData := []byte("same secret data")

		encrypted1, err := cryptor.EncryptSecretData(plainData)
		require.NoError(t, err)

		encrypted2, err := cryptor.EncryptSecretData(plainData)
		require.NoError(t, err)

		assert.NotEqual(t, encrypted1, encrypted2)

		decrypted1, err := cryptor.DecryptSecretData(encrypted1)
		require.NoError(t, err)
		assert.Equal(t, plainData, decrypted1)

		decrypted2, err := cryptor.DecryptSecretData(encrypted2)
		require.NoError(t, err)
		assert.Equal(t, plainData, decrypted2)
	})

	t.Run("decrypt with wrong key fails", func(t *testing.T) {
		plainData := []byte("test data")
		encrypted, err := cryptor.EncryptStorageData(plainData)
		require.NoError(t, err)

		wrongCryptor := NewCryptor("wrongpassword", "testuser")
		_, err = wrongCryptor.DecryptStorageData(encrypted)
		require.Error(t, err)
	})
}
