package storage

import (
	"context"
	"github.com/etoneja/go-keeper/internal/ctl/types"
)

type Storager interface {
	CreateSecret(ctx context.Context, secret *types.LocalSecret) (*types.LocalSecret, error)
	GetSecret(ctx context.Context, secretID string) (*types.LocalSecret, error)
	UpdateSecret(ctx context.Context, secret *types.LocalSecret) error
	DeleteSecret(ctx context.Context, secretID string) error
	ListSecrets(ctx context.Context) ([]*types.LocalSecret, error)

	Close() error
}
