package storage

import (
	"context"
	"github.com/etoneja/go-keeper/internal/ctl/types"
)

// Storager defines the interface for secret storage
type Storager interface {
	// CreateSchema creates the storage schema
	CreateSchema(ctx context.Context) error

	// Secrets management
	CreateSecret(ctx context.Context, secret *types.Secret) (*types.Secret, error)
	GetSecret(ctx context.Context, uuid string) (*types.Secret, error)
	UpdateSecret(ctx context.Context, secret *types.Secret) error
	DeleteSecret(ctx context.Context, uuid string) error
	ListSecrets(ctx context.Context) ([]*types.Secret, error)

	// Close cleans up resources
	Close() error
}
