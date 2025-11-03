package ctl

import (
	"context"
	"fmt"
	"sync"

	"github.com/etoneja/go-keeper/internal/crypto"
	"github.com/etoneja/go-keeper/internal/ctl/storage"
	"github.com/etoneja/go-keeper/internal/ctl/types"
)

type VaultService struct {
	cfg     *Config
	crypter crypto.Crypter

	storage   storage.Storager
	storageMu sync.Mutex
}

func NewVaultService(cfg *Config) *VaultService {
	crypter := crypto.NewCrypto(cfg.Password)

	return &VaultService{
		cfg:     cfg,
		crypter: crypter,
	}
}

func (s *VaultService) getStorage(ctx context.Context) (storage.Storager, error) {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

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

// Close cleans up resources
func (s *VaultService) Close() error {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	if s.storage != nil {
		return s.storage.Close()
	}
	return nil
}

// Initialize sets up the storage
func (s *VaultService) Initialize(ctx context.Context) error {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

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

// Create methods with models
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

// Get methods
func (s *VaultService) GetSecret(ctx context.Context, uuid string) (*types.Secret, error) {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return nil, err
	}

	return storage.GetSecret(ctx, uuid)
}

func (s *VaultService) ListSecrets(ctx context.Context) ([]*types.Secret, error) {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return nil, err
	}

	return storage.ListSecrets(ctx)
}

// Delete method
func (s *VaultService) DeleteSecret(ctx context.Context, uuid string) error {
	storage, err := s.getStorage(ctx)
	if err != nil {
		return err
	}
	_, err = storage.GetSecret(ctx, uuid)
	if err != nil {
		return err
	}
	return storage.DeleteSecret(ctx, uuid)
}
