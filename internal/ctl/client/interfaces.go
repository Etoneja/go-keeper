package client

import (
	"context"

	"github.com/etoneja/go-keeper/internal/ctl/types"
)

type Clienter interface {
	Connect(ctx context.Context) error
	Close() error

	Login(ctx context.Context) error
	Register(ctx context.Context) (string, error)

	SetSecret(ctx context.Context, secret types.Secreter) error
	GetSecret(ctx context.Context, secretID string) (types.Secreter, error)
	DeleteSecret(ctx context.Context, secretID string) error
	ListSecrets(ctx context.Context) ([]types.Secreter, error)
}
