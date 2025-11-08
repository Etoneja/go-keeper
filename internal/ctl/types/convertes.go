package types

import (
	"encoding/json"

	"github.com/etoneja/go-keeper/internal/crypto"
)

func ConvertLocalSecretToRemoteSecret(cryptor crypto.Cryptor, localSecret *LocalSecret) (*RemoteSecret, error) {
	secretData, err := localSecret.ParseData()
	if err != nil {
		return nil, err
	}

	secretDataContainer := &SecretDataContainer{
		Type:       localSecret.Type,
		Name:       localSecret.Name,
		SecretData: secretData,
	}

	remoteData, err := json.Marshal(secretDataContainer)
	if err != nil {
		return nil, err
	}

	encryptedRemoteData, err := cryptor.EncryptSecretData(remoteData)
	if err != nil {
		return nil, err
	}

	remoteSecret := &RemoteSecret{
		UUID:         localSecret.UUID,
		LastModified: localSecret.LastModified,
		Hash:         localSecret.Hash,
		Data:         encryptedRemoteData,
	}

	return remoteSecret, nil
}

func ConvertRemoteSecretToLocalSecret(cryptor crypto.Cryptor, remoteSecret *RemoteSecret) (*LocalSecret, error) {
	remoteDecryptedData, err := cryptor.DecryptSecretData(remoteSecret.Data)
	if err != nil {
		return nil, err
	}

	var secretDataContainer SecretDataContainer
	if err := json.Unmarshal(remoteDecryptedData, &secretDataContainer); err != nil {
		return nil, err
	}

	localSecret := &LocalSecret{
		UUID:         remoteSecret.UUID,
		Type:         secretDataContainer.Type,
		Name:         secretDataContainer.Name,
		LastModified: remoteSecret.LastModified,
		Hash:         remoteSecret.Hash,
	}

	err = localSecret.SetData(cryptor, secretDataContainer.SecretData)
	if err != nil {
		return nil, err
	}

	return localSecret, nil
}
