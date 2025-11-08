package ctl

import (
	"context"
	"fmt"
	"sync"

	"github.com/etoneja/go-keeper/internal/crypto"
	"github.com/etoneja/go-keeper/internal/ctl/client"
	"github.com/etoneja/go-keeper/internal/ctl/config"
	"github.com/etoneja/go-keeper/internal/ctl/storage"
	"github.com/etoneja/go-keeper/internal/ctl/types"
)

type VaultService struct {
	cfg     *config.Config
	crypter crypto.Crypter

	mu sync.Mutex

	storage storage.Storager
	client  client.Clienter
}

func NewVaultService(cfg *config.Config) *VaultService {
	crypter := crypto.NewCrypto(cfg.Password)

	return &VaultService{
		cfg:     cfg,
		crypter: crypter,
	}
}

func (s *VaultService) getStorage(ctx context.Context) (storage.Storager, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage != nil {
		return s.storage, nil
	}

	storage, err := storage.NewStorage(ctx, s.cfg.DBPath)
	if err != nil {
		return nil, err
	}

	s.storage = storage
	return s.storage, nil
}

func (s *VaultService) getClient(ctx context.Context) (client.Clienter, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		return s.client, nil
	}

	client := client.NewGRPCClient(s.cfg)

	err := client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	s.client = client
	return s.client, nil
}

func (s *VaultService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage != nil {
		return s.storage.Close()
	}

	return nil
}

func (s *VaultService) Initialize(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage != nil {
		return fmt.Errorf("service already initialized")
	}

	storage, err := storage.CreateNewStorage(ctx, s.cfg.DBPath)
	if err != nil {
		return err
	}

	s.storage = storage
	return nil
}

func (s *VaultService) CreateSecret(ctx context.Context, secret *types.Secret) (*types.Secret, error) {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return nil, err
	}

	newsecret, err := storage.CreateSecret(ctx, secret)
	if err != nil {
		return nil, err
	}

	return newsecret, nil
}

func (s *VaultService) GetSecret(ctx context.Context, secretId string) (*types.Secret, error) {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return nil, err
	}

	secret, err := storage.GetSecret(ctx, secretId)
	if err != nil {
		return nil, err
	}

	return secret, err
}

func (s *VaultService) ListSecrets(ctx context.Context) ([]*types.Secret, error) {
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

func (s *VaultService) DeleteSecret(ctx context.Context, secretId string) error {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return err
	}

	_, err = storage.GetSecret(ctx, secretId)
	if err != nil {
		return err
	}

	err = storage.DeleteSecret(ctx, secretId)
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
