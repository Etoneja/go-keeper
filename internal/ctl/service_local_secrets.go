package ctl

import (
	"context"

	"github.com/etoneja/go-keeper/internal/ctl/types"
)

func (s *VaultService) CreateLocalSecret(ctx context.Context, secret *types.LocalSecret) (*types.LocalSecret, error) {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return nil, err
	}

	newSecret, err := storage.CreateSecret(ctx, secret)
	if err != nil {
		return nil, err
	}

	return newSecret, nil
}

func (s *VaultService) GetLocalSecret(ctx context.Context, secretID string) (*types.LocalSecret, error) {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return nil, err
	}

	secret, err := storage.GetSecret(ctx, secretID)
	if err != nil {
		return nil, err
	}

	return secret, err
}

func (s *VaultService) ListLocalSecrets(ctx context.Context) ([]*types.LocalSecret, error) {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return nil, err
	}

	secrets, err := storage.ListSecrets(ctx)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (s *VaultService) DeleteLocalSecret(ctx context.Context, secretID string) error {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return err
	}

	_, err = storage.GetSecret(ctx, secretID)
	if err != nil {
		return err
	}

	err = storage.DeleteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}
