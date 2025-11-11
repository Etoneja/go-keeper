package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	UserRepo   UserRepositorier
	SecretRepo SecretRepositorier
}

func NewRepositories() *Repositories {
	userRepo := NewUserRepository()
	secretRepo := NewSecretRepository()
	return &Repositories{
		UserRepo:   userRepo,
		SecretRepo: secretRepo,
	}
}

func (r *Repositories) WithTx(ctx context.Context, db *pgxpool.Pool, fn func(Querier) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
