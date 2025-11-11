package ctl

import (
	"context"
	"fmt"

	"github.com/etoneja/go-keeper/internal/ctl/types"
)

func (s *VaultService) getDiff(ctx context.Context) (*types.SecretsDiff, error) {
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

func (s *VaultService) processDiff(ctx context.Context, diff *types.SecretsDiff) error {
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

	for _, pair := range diff.Both {
		err := s.syncSecretCheckPair(ctx, pair)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *VaultService) deleteLocalSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Deleting local secret '%s'\n", secretID)

	storage, err := s.getStorage(ctx)
	if err != nil {
		return err
	}

	err = storage.DeleteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) createRemoteSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Creating remote secret '%s'\n", secretID)

	storage, err := s.getStorage(ctx)
	if err != nil {
		return err
	}

	localSecret, err := storage.GetSecret(ctx, secretID, true)
	if err != nil {
		return err
	}

	remoteSecret, err := types.ConvertLocalSecretToRemoteSecret(s.cryptor, localSecret)
	if err != nil {
		return err
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

func (s *VaultService) createLocalSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Creating local secret '%s'\n", secretID)

	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}
	remoteSecret, err := client.GetSecret(ctx, secretID)
	if err != nil {
		return err
	}

	localSecret, err := types.ConvertRemoteSecretToLocalSecret(s.cryptor, remoteSecret)
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

func (s *VaultService) deleteRemoteSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Deleting remote secret '%s'\n", secretID)

	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}

	err = client.DeleteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) replaceLocalSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Replacing remote secret '%s'\n", secretID)

	err := s.deleteLocalSecret(ctx, secretID)
	if err != nil {
		return err
	}

	err = s.createLocalSecret(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) replaceRemoteSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Replacing remote secret '%s'\n", secretID)

	err := s.deleteRemoteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	err = s.createRemoteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) syncLocalSecret(ctx context.Context, localSecret *types.LocalSecret) error {
	action, err := PromptForLocalOnlyAction(localSecret)
	if err != nil {
		return err
	}

	switch action {
	case ActionDeleteLocal:
		err := s.deleteLocalSecret(ctx, localSecret.UUID)
		if err != nil {
			return err
		}
	case ActionCreateRemote:
		err := s.createRemoteSecret(ctx, localSecret.UUID)
		if err != nil {
			return err
		}
	case ActionSkip:
		fmt.Printf("Ignoring secret '%s'\n", localSecret.UUID)
	}

	return nil
}

func (s *VaultService) syncRemoteSecret(ctx context.Context, remoteSecret *types.RemoteSecret) error {
	action, err := PromptForRemoteOnlyAction(remoteSecret)
	if err != nil {
		return err
	}

	switch action {
	case ActionCreateLocal:
		err := s.createLocalSecret(ctx, remoteSecret.UUID)
		if err != nil {
			return err
		}
	case ActionDeleteRemote:
		err := s.deleteRemoteSecret(ctx, remoteSecret.UUID)
		if err != nil {
			return err
		}
	case ActionSkip:
		fmt.Printf("Ignoring secret '%s'\n", remoteSecret.UUID)
	}

	return nil
}

func (s *VaultService) syncSecretCheckPair(ctx context.Context, checkPair *types.SecretCheckPair) error {
	if checkPair.IsIdentical() {
		return nil
	}

	// TODO: Add option to show diff

	action, err := PromptForConflictCheckPairAction(checkPair)
	if err != nil {
		return err
	}

	switch action {
	case ActionReplaceLocal:
		err := s.replaceLocalSecret(ctx, checkPair.Local.UUID)
		if err != nil {
			return err
		}
	case ActionReplaceRemote:
		err := s.replaceRemoteSecret(ctx, checkPair.Remote.UUID)
		if err != nil {
			return err
		}
	case ActionSkip:
		fmt.Printf("Ignoring secret '%s'\n", checkPair.Local.UUID)
	}

	return nil
}
