package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/etoneja/go-keeper/internal/crypto"
)

func initializeSQLiteStorage(ctx context.Context, cryptor crypto.Cryptor, dbPath string) error {
	if _, err := os.Stat(dbPath); err == nil {
		return fmt.Errorf("vault already exists at %s", dbPath)
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := openInMemoryDB()
	if err != nil {
		return err
	}

	storage := &SQLiteStorage{
		db:      db,
		cryptor: cryptor,
		path:    dbPath,
		isDirty: false,
	}

	err = storage.createSchema(ctx)
	if err != nil {
		return err
	}

	err = storage.Close()
	if err != nil {
		return err
	}

	return nil
}

func openSQLiteStorage(ctx context.Context, cryptor crypto.Cryptor, dbPath string) (*SQLiteStorage, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("storage is not initialized: %s not found", dbPath)
	}

	encryptedData, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read db file: %w", err)
	}

	decryptedData, err := cryptor.DecryptStorageData(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt db: %w", err)
	}

	db, err := deserializeInMemoryDBFromBytes(ctx, decryptedData)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	storage := &SQLiteStorage{
		db:      db,
		cryptor: cryptor,
		path:    dbPath,
		isDirty: false,
	}

	return storage, nil
}
