package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DefaultTxManager struct {
	Repos *Repositories
	Db    *pgxpool.Pool
}

func (d *DefaultTxManager) WithTx(ctx context.Context, fn func(Querier) error) error {
	return d.Repos.WithTx(ctx, d.Db, fn)
}
