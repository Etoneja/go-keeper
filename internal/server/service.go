package server

import (
	"context"
	"errors"

	"github.com/etoneja/go-keeper/internal/server/repository"
	"github.com/etoneja/go-keeper/internal/server/token"
	"github.com/etoneja/go-keeper/internal/server/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrSecretTooLarge     = errors.New("secret data too large")
)

type Service struct {
	db           *pgxpool.Pool
	tokenManager token.TokenManager
	repos        *repository.Repositories
}

func NewService(db *pgxpool.Pool, tokenManager token.TokenManager, repos *repository.Repositories) *Service {
	return &Service{
		db:           db,
		tokenManager: tokenManager,
		repos:        repos,
	}
}

func (s *Service) Register(ctx context.Context, login, password string) (*types.User, error) {
	var user *types.User

	err := s.repos.WithTx(ctx, s.db, func(q repository.Querier) error {
		existing, _ := s.repos.UserRepo.GetUserByLogin(ctx, q, login)
		if existing != nil {
			return ErrUserAlreadyExists
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		createdUser, err := s.repos.UserRepo.CreateUser(ctx, q, login, string(passwordHash))
		if err != nil {
			return err
		}

		user = createdUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (string, *types.User, error) {
	user, err := s.repos.UserRepo.GetUserByLogin(ctx, s.db, login)
	if err != nil {
		return "", nil, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.tokenManager.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *Service) SetSecret(ctx context.Context, secret *types.Secret) error {
	if len(secret.Data) > maxSecretSize {
		return ErrSecretTooLarge
	}

	return s.repos.SecretRepo.SetSecret(ctx, s.db, secret)
}

func (s *Service) GetSecret(ctx context.Context, userID, secretID string) (*types.Secret, error) {
	// TODO: check ownership in service
	return s.repos.SecretRepo.GetSecret(ctx, s.db, userID, secretID)
}

func (s *Service) DeleteSecret(ctx context.Context, userID, secretID string) error {
	// TODO: check ownership in service
	return s.repos.SecretRepo.DeleteSecret(ctx, s.db, userID, secretID)
}

func (s *Service) ListSecrets(ctx context.Context, userID string) ([]*types.Secret, error) {
	return s.repos.SecretRepo.ListSecrets(ctx, s.db, userID)
}
