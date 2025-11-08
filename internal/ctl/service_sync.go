package ctl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/etoneja/go-keeper/internal/ctl/types"
)

func (s *VaultService) getDiff(ctx context.Context) (*types.SecretDiff, error) {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return nil, err
	}

	localSecrets, err := storage.ListSecrets(ctx)
	if err != nil {
		return nil, err
	}

	client, err := s.getClient(ctx)
	if err != nil {
		return nil, err
	}

	remoteSecrets, err := client.ListSecrets(ctx)
	if err != nil {
		return nil, err
	}

	diff := diffSecrets(localSecrets, remoteSecrets)
	return diff, err
}

func (s *VaultService) processDiff(ctx context.Context, diff *types.SecretDiff) error {
	fmt.Printf("local_only: %d, remote_only: %d, both: %d\n",
		len(diff.LocalOnly), len(diff.RemoteOnly), len(diff.Both))

	for _, secret := range diff.LocalOnly {
		err := s.syncLocalSecret(ctx, secret)
		if err != nil {
			return err
		}
	}

	for _, secret := range diff.RemoteOnly {
		err := s.syncRemoteSecret(ctx, secret)
		if err != nil {
			return err
		}
	}

	// proceed both
	for _, pair := range diff.Both {
		if pair.Local.GetHash() != pair.Remote.GetHash() {
			fmt.Printf("Secret %s has different hash: local=%s vs remote=%s\n",
				pair.Local.GetUUID(), pair.Local.GetHash(), pair.Remote.GetHash())
		}
		fmt.Printf("%s - %s\n", pair.Local.GetLastModified(), pair.Remote.GetLastModified())
		if pair.Local.GetLastModified().After(pair.Remote.GetLastModified()) {
			fmt.Printf("Secret %s is newer locally\n", pair.Local.GetUUID())
		}
	}

	return nil
}

func (s *VaultService) deleteLocalSecret(ctx context.Context, secretId string) error {
	fmt.Printf("Deleting local secret '%s'\n", secretId)

	storage, err := s.getStorage(ctx)
	if err != nil {
		return err
	}

	err = storage.DeleteSecret(ctx, secretId)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) createRemoteSecret(ctx context.Context, secretId string) error {
	fmt.Printf("Creating remote secret '%s'\n", secretId)

	storage, err := s.getStorage(ctx)
	if err != nil {
		return err
	}

	localSecret, err := storage.GetSecret(ctx, secretId)
	if err != nil {
		return err
	}

	// TODO: Move to converters + add encryption
	secretData, err := localSecret.ParseData()
	if err != nil {
		return err
	}

	secretDataContainer := &types.SecretDataContainer{
		Type:       localSecret.Type,
		Name:       localSecret.Name,
		SecretData: secretData,
	}

	remoteData, err := json.Marshal(secretDataContainer)
	if err != nil {
		return err
	}

	// TODO: Add encryption

	remoteSecret := &types.RemoteSecret{
		UUID:         localSecret.UUID,
		LastModified: localSecret.LastModified,
		Hash:         localSecret.Hash,
		Data:         remoteData,
	}

	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}

	err = client.SetSecret(ctx, remoteSecret)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) createLocalSecret(ctx context.Context, secretId string) error {
	fmt.Printf("Creating local secret '%s'\n", secretId)

	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}
	remoteSecret, err := client.GetSecret(ctx, secretId)
	if err != nil {
		return err
	}

	// TODO: Add decryption

	var secretDataContainer types.SecretDataContainer
	if err := json.Unmarshal(remoteSecret.GetData(), &secretDataContainer); err != nil {
		return err
	}

	localSecret := &types.Secret{
		UUID:         remoteSecret.GetUUID(),
		Type:         secretDataContainer.Type,
		Name:         secretDataContainer.Name,
		LastModified: remoteSecret.GetLastModified(),
		Hash:         remoteSecret.GetHash(),
	}

	err = localSecret.SetData(secretDataContainer.SecretData, s.crypter)
	if err != nil {
		return err
	}

	storage, err := s.getStorage(ctx)
	if err != nil {
		return err
	}

	_, err = storage.CreateSecret(ctx, localSecret)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) deleteRemoteSecret(ctx context.Context, secretId string) error {
	fmt.Printf("Deleting remote secret '%s'\n", secretId)

	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}

	err = client.DeleteSecret(ctx, secretId)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) syncLocalSecret(ctx context.Context, secret types.Secreter) error {
	action, err := PromptForLocalOnlyAction(secret.GetUUID())
	if err != nil {
		return err
	}

	switch action {
	case ActionDeleteLocal:
		err := s.deleteLocalSecret(ctx, secret.GetUUID())
		if err != nil {
			return err
		}
	case ActionCreateRemote:
		err := s.createRemoteSecret(ctx, secret.GetUUID())
		if err != nil {
			return err
		}
	case ActionSkip:
		fmt.Printf("Ignoring secret '%s'\n", secret.GetUUID())
	}

	return nil
}

func (s *VaultService) syncRemoteSecret(ctx context.Context, secret types.Secreter) error {
	action, err := PromptForRemoteOnlyAction(secret.GetUUID())
	if err != nil {
		return err
	}

	switch action {
	case ActionCreateLocal:
		err := s.createLocalSecret(ctx, secret.GetUUID())
		if err != nil {
			return err
		}
	case ActionDeleteRemote:
		err := s.deleteRemoteSecret(ctx, secret.GetUUID())
		if err != nil {
			return err
		}
	case ActionSkip:
		fmt.Printf("Ignoring secret '%s'\n", secret.GetUUID())
	}

	return nil
}
