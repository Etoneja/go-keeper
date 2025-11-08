package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/etoneja/go-keeper/internal/crypto"
	"github.com/etoneja/go-keeper/internal/ctl/errs"
	"github.com/etoneja/go-keeper/internal/ctl/types"
)

type SQLiteStorage struct {
	db      *sql.DB
	cryptor crypto.Cryptor
	path    string
	isDirty bool
}

func (s *SQLiteStorage) markClean() {
	s.isDirty = false
}

func (s *SQLiteStorage) markDirty() {
	s.isDirty = true
}

func (s *SQLiteStorage) dump() error {
	if !s.isDirty {
		return nil
	}

	decryptedBytes, err := serializeInMemoryDBToBytes(context.Background(), s.db)
	if err != nil {
		return err
	}

	encryptedData, err := s.cryptor.EncryptStorageData(decryptedBytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt db: %w", err)
	}

	if err := os.WriteFile(s.path, encryptedData, 0600); err != nil {
		return fmt.Errorf("failed to write db file: %w", err)
	}

	s.markClean()

	return nil
}

func (s *SQLiteStorage) Close() error {
	err := s.dump()
	if err != nil {
		return err
	}

	return s.db.Close()
}

func (s *SQLiteStorage) createSchema(ctx context.Context) error {
	query := `
	CREATE TABLE secrets (
		uuid TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		last_modified DATETIME NOT NULL,
		hash TEXT NOT NULL,
		metadata TEXT,
		data BLOB NOT NULL
	);
	`
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	s.markDirty()

	return nil
}

func (s *SQLiteStorage) CreateSecret(ctx context.Context, secret *types.LocalSecret) (*types.LocalSecret, error) {
	query := `
		INSERT INTO secrets (uuid, type, name, last_modified, hash, metadata, data)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		secret.UUID,
		secret.Type,
		secret.Name,
		secret.LastModified,
		secret.Hash,
		secret.Metadata,
		secret.Data,
	)
	if err != nil {
		return nil, err
	}

	s.markDirty()

	return secret, nil
}

func (s *SQLiteStorage) GetSecret(ctx context.Context, uuid string) (*types.LocalSecret, error) {
	query := `
		SELECT uuid, type, name, last_modified, hash, metadata, data
		FROM secrets
		WHERE uuid = ?
	`

	row := s.db.QueryRowContext(ctx, query, uuid)

	secret := &types.LocalSecret{}

	err := row.Scan(
		&secret.UUID,
		&secret.Type,
		&secret.Name,
		&secret.LastModified,
		&secret.Hash,
		&secret.Metadata,
		&secret.Data,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewSecretNotFoundError(uuid)
		}
		return nil, err
	}

	return secret, nil
}

func (s *SQLiteStorage) UpdateSecret(ctx context.Context, secret *types.LocalSecret) error {
	query := `
		UPDATE secrets 
		SET type = ?, name = ?, last_modified = ?, hash = ?, metadata = ?, data = ?
		WHERE uuid = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		secret.Type,
		secret.Name,
		secret.LastModified,
		secret.Hash,
		secret.Metadata,
		secret.Data,
		secret.UUID,
	)
	if err != nil {
		return err
	}

	s.markDirty()

	return nil
}

func (s *SQLiteStorage) DeleteSecret(ctx context.Context, uuid string) error {
	query := `DELETE FROM secrets WHERE uuid = ?`
	_, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	s.markDirty()

	return nil
}

func (s *SQLiteStorage) ListSecrets(ctx context.Context) ([]*types.LocalSecret, error) {
	query := `
		SELECT uuid, type, name, last_modified, hash, metadata
		FROM secrets
		ORDER BY last_modified DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*types.LocalSecret
	for rows.Next() {
		secret := &types.LocalSecret{}

		err := rows.Scan(
			&secret.UUID,
			&secret.Type,
			&secret.Name,
			&secret.LastModified,
			&secret.Hash,
			// NOTE: Do not fetch data for listing
			// &secret.Data,
			&secret.Metadata,
		)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}
