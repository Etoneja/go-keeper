package ctl

import (
	"context"
	"log"
	"sync"

	"github.com/etoneja/go-keeper/internal/ctl/client"
	"github.com/etoneja/go-keeper/internal/ctl/config"
	"github.com/etoneja/go-keeper/internal/ctl/crypto"
	"github.com/etoneja/go-keeper/internal/ctl/storage"
)

type VaultService struct {
	cfg     *config.Config
	cryptor crypto.Cryptor

	mu sync.Mutex

	storage storage.Storager
	client  client.Clienter
}

func NewVaultService(cfg *config.Config) *VaultService {
	cryptor := crypto.NewCryptor(cfg.Password, cfg.Login)

	return &VaultService{
		cfg:     cfg,
		cryptor: cryptor,
	}
}

func (s *VaultService) getStorage(ctx context.Context) (storage.Storager, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage != nil {
		return s.storage, nil
	}

	cryptor := crypto.NewCryptor(s.cfg.Password, s.cfg.Login)

	storage, err := storage.NewStorage(ctx, cryptor, s.cfg.DBPath)
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

	serverPassword := s.cryptor.GenerateServerPassword()

	client := client.NewGRPCClient(s.cfg.ServerAddress, s.cfg.Login, serverPassword)

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
		if err := s.storage.Close(); err != nil {
			log.Printf("failed to close storage: %v", err)
		}
	}

	if s.client != nil {
		if err := s.client.Close(); err != nil {
			log.Printf("failed to close client: %v", err)
		}
	}

	return nil
}
