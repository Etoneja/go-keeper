package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func CreateNewStorage(ctx context.Context, path string) (Storager, error) {
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("vault already exists at %s", path)
	}

	dbDir := filepath.Dir(path)
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	storage, err := NewSQLiteStorage(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	if err := storage.CreateSchema(ctx); err != nil {
		storage.Close()
		os.Remove(path)
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return storage, nil
}

func NewStorage(ctx context.Context, path string) (Storager, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault is not initialized")
	}

	storage, err := NewSQLiteStorage(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	return storage, nil
}
