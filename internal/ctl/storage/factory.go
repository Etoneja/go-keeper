package storage

import (
	"context"
	"github.com/etoneja/go-keeper/internal/crypto"
)

func InitializeStorage(ctx context.Context, cryptor crypto.Cryptor, path string) error {
	err := initializeSQLiteStorage(ctx, cryptor, path)
	if err != nil {
		return err
	}

	return nil
}

func NewStorage(ctx context.Context, cryptor crypto.Cryptor, path string) (Storager, error) {
	storage, err := openSQLiteStorage(ctx, cryptor, path)
	if err != nil {
		return nil, err
	}

	return storage, nil
}
