package server

import (
	"context"

	"github.com/etoneja/go-keeper/internal/server/stypes"
)

type Servicer interface {
	Register(ctx context.Context, login, password string) (*stypes.User, error)
	Login(ctx context.Context, login, password string) (string, *stypes.User, error)
	SetSecret(ctx context.Context, secret *stypes.Secret) error
	GetSecret(ctx context.Context, userID, secretID string) (*stypes.Secret, error)
	DeleteSecret(ctx context.Context, userID, secretID string) error
	ListSecrets(ctx context.Context, userID string) ([]*stypes.Secret, error)
}
