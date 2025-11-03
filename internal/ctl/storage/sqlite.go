package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/etoneja/go-keeper/internal/ctl/errs"
	"github.com/etoneja/go-keeper/internal/ctl/types"
	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
}

// TODO: add storage encryption and decryption
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &SQLiteStorage{db: db}
	return storage, nil
}

func (s *SQLiteStorage) CreateSchema(ctx context.Context) error {
	query := `
	CREATE TABLE secrets (
		uuid TEXT PRIMARY KEY,
		last_modified DATETIME NOT NULL,
		hash TEXT NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		data BLOB NOT NULL,
		metadata TEXT
	);
	`

	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SQLiteStorage) CreateSecret(ctx context.Context, secret *types.Secret) (*types.Secret, error) {
	query := `
	INSERT INTO secrets (uuid, last_modified, hash, name, type, data, metadata)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		secret.UUID,
		secret.LastModified,
		secret.Hash,
		secret.Name,
		secret.Type,
		secret.Data,
		secret.Metadata,
	)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (s *SQLiteStorage) GetSecret(ctx context.Context, uuid string) (*types.Secret, error) {
	query := `
	SELECT uuid, last_modified, hash, name, type, data, metadata
	FROM secrets
	WHERE uuid = ?
	`

	row := s.db.QueryRowContext(ctx, query, uuid)

	secret := &types.Secret{}

	err := row.Scan(
		&secret.UUID,
		&secret.LastModified,
		&secret.Hash,
		&secret.Name,
		&secret.Type,
		&secret.Data,
		&secret.Metadata,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewSecretNotFoundError(uuid)
		}
		return nil, err
	}

	return secret, nil
}

func (s *SQLiteStorage) UpdateSecret(ctx context.Context, secret *types.Secret) error {
	query := `
	UPDATE secrets 
	SET last_modified = ?, hash = ?, name = ?, type = ?, data = ?, metadata = ?
	WHERE uuid = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		secret.LastModified,
		secret.Hash,
		secret.Name,
		secret.Type,
		secret.Data,
		secret.Metadata,
		secret.UUID,
	)
	return err
}

func (s *SQLiteStorage) DeleteSecret(ctx context.Context, uuid string) error {
	query := `DELETE FROM secrets WHERE uuid = ?`
	_, err := s.db.ExecContext(ctx, query, uuid)
	return err
}

func (s *SQLiteStorage) ListSecrets(ctx context.Context) ([]*types.Secret, error) {
	query := `
	SELECT uuid, last_modified, hash, name, type, data, metadata
	FROM secrets
	ORDER BY last_modified DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*types.Secret
	for rows.Next() {
		secret := &types.Secret{}

		err := rows.Scan(
			&secret.UUID,
			&secret.LastModified,
			&secret.Hash,
			&secret.Name,
			&secret.Type,
			&secret.Data,
			&secret.Metadata,
		)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
