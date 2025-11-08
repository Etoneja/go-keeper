package repository

import (
	"context"
	"errors"

	"github.com/etoneja/go-keeper/internal/server/types"
)

var (
	ErrSecretNotFound = errors.New("secret not found")
)

type SecretRepository struct{}

func NewSecretRepository() *SecretRepository {
	return &SecretRepository{}
}

func (r *SecretRepository) SetSecret(ctx context.Context, q Querier, secret *types.Secret) error {
	query := `
		INSERT INTO secrets (id, user_id, data, hash, last_modified)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			data = $3,
			hash = $4, 
			last_modified = $5
		WHERE secrets.user_id = $2
	`

	_, err := q.Exec(ctx, query,
		secret.ID,
		secret.UserID,
		secret.Data,
		secret.Hash,
		secret.LastModified,
	)

	return err
}

func (r *SecretRepository) GetSecret(ctx context.Context, q Querier, userID, secretID string) (*types.Secret, error) {
	query := `
		SELECT id, user_id, data, hash, last_modified 
		FROM secrets 
		WHERE user_id = $1 AND id = $2
	`

	var secret types.Secret
	err := q.QueryRow(ctx, query, userID, secretID).Scan(
		&secret.ID,
		&secret.UserID,
		&secret.Data,
		&secret.Hash,
		&secret.LastModified,
	)

	if err != nil {
		return nil, ErrSecretNotFound
	}

	return &secret, nil
}

func (r *SecretRepository) DeleteSecret(ctx context.Context, q Querier, userID, secretID string) error {
	query := `
		DELETE FROM secrets
		WHERE user_id = $1 AND id = $2
	`

	result, err := q.Exec(ctx, query, userID, secretID)
	if err != nil {
		return err
	}

	// Check if any row was affected
	if result.RowsAffected() == 0 {
		return ErrSecretNotFound
	}

	return nil
}

func (r *SecretRepository) ListSecrets(ctx context.Context, q Querier, userID string) ([]*types.Secret, error) {
	query := `
		SELECT id, user_id, hash, last_modified
		FROM secrets 
		WHERE user_id = $1
		ORDER BY last_modified DESC
	`

	rows, err := q.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*types.Secret
	for rows.Next() {
		var secret types.Secret
		err := rows.Scan(
			&secret.ID,
			&secret.UserID,
			&secret.Hash,
			&secret.LastModified,
		)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, &secret)
	}
	return secrets, nil
}
