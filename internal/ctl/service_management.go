package ctl

import (
	"context"
	"fmt"

	"github.com/etoneja/go-keeper/internal/crypto"
	"github.com/etoneja/go-keeper/internal/ctl/storage"
)

func (s *VaultService) Initialize(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage != nil {
		return fmt.Errorf("service already initialized")
	}

	cryptor := crypto.NewCryptor(s.cfg.Password, s.cfg.Login)

	err := storage.InitializeStorage(ctx, cryptor, s.cfg.DBPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) Register(ctx context.Context) error {
	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.Register(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *VaultService) SyncSecrets(ctx context.Context) error {
	diff, err := s.getDiff(ctx)
	if err != nil {
		return err
	}

	err = s.processDiff(ctx, diff)
	if err != nil {
		return err
	}

	return nil
}
