package repository

import (
	context "context"

	"github.com/etoneja/go-keeper/internal/server/stypes"
	"github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(Querier) error) error
}

type UserRepositorier interface {
	CreateUser(ctx context.Context, q Querier, login, passwordHash string) (*stypes.User, error)
	GetUserByLogin(ctx context.Context, q Querier, login string) (*stypes.User, error)
	GetUserByID(ctx context.Context, q Querier, userID string) (*stypes.User, error)
}

type SecretRepositorier interface {
	SetSecret(ctx context.Context, q Querier, secret *stypes.Secret) error
	GetSecret(ctx context.Context, q Querier, userID, secretID string) (*stypes.Secret, error)
	DeleteSecret(ctx context.Context, q Querier, userID, secretID string) error
	ListSecrets(ctx context.Context, q Querier, userID string) ([]*stypes.Secret, error)
}
