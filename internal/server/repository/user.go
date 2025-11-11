package repository

import (
	"context"
	"errors"

	"github.com/etoneja/go-keeper/internal/server/stypes"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(ctx context.Context, q Querier, login, passwordHash string) (*stypes.User, error) {
	query := `
		INSERT INTO users (login, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, login, password_hash, created_at
	`

	var user stypes.User
	err := q.QueryRow(ctx, query, login, passwordHash).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, q Querier, login string) (*stypes.User, error) {
	query := `
		SELECT id, login, password_hash, created_at
		FROM users
		WHERE login = $1
	`

	var user stypes.User
	err := q.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, q Querier, userID string) (*stypes.User, error) {
	query := `
		SELECT id, login, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	var user stypes.User
	err := q.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}
