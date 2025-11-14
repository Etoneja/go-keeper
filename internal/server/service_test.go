package server

import (
	"context"
	"testing"

	"github.com/etoneja/go-keeper/internal/server/repository"
	"github.com/etoneja/go-keeper/internal/server/stypes"
	"github.com/etoneja/go-keeper/internal/server/token"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type MockTxManager struct {
	querier repository.Querier
}

func (m *MockTxManager) WithTx(ctx context.Context, fn func(repository.Querier) error) error {
	return fn(m.querier)
}

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenManager := token.NewMockTokenManager(ctrl)
	mockUserRepo := repository.NewMockUserRepositorier(ctrl)
	mockSecretRepo := repository.NewMockSecretRepositorier(ctrl)
	mockQuerier := repository.NewMockQuerier(ctrl)

	repos := &repository.Repositories{
		UserRepo:   mockUserRepo,
		SecretRepo: mockSecretRepo,
	}
	txManager := &MockTxManager{querier: mockQuerier}

	service := NewService(nil, mockTokenManager, txManager, repos)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), mockQuerier, "testuser").
			Return(nil, repository.ErrUserNotFound)
		mockUserRepo.EXPECT().CreateUser(gomock.Any(), mockQuerier, "testuser", gomock.Any()).
			Return(&stypes.User{ID: "123"}, nil)

		user, err := service.Register(context.Background(), "testuser", "pass")
		require.NoError(t, err)
		assert.Equal(t, "123", user.ID)
	})

	t.Run("user exists", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), mockQuerier, "exists").
			Return(&stypes.User{ID: "456"}, nil)

		user, err := service.Register(context.Background(), "exists", "pass")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUserAlreadyExists)
		assert.Nil(t, user)
	})
}

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenManager := token.NewMockTokenManager(ctrl)
	mockUserRepo := repository.NewMockUserRepositorier(ctrl)
	mockSecretRepo := repository.NewMockSecretRepositorier(ctrl)

	repos := &repository.Repositories{
		UserRepo:   mockUserRepo,
		SecretRepo: mockSecretRepo,
	}
	service := NewService(nil, mockTokenManager, nil, repos)

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	testUser := &stypes.User{ID: "123", PasswordHash: string(passwordHash)}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any(), "testuser").
			Return(testUser, nil)
		mockTokenManager.EXPECT().GenerateToken("123").Return("token123", nil)

		token, user, err := service.Login(context.Background(), "testuser", "pass")
		require.NoError(t, err)
		assert.Equal(t, "token123", token)
		assert.Equal(t, "123", user.ID)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any(), "unknown").
			Return(nil, repository.ErrUserNotFound)

		_, _, err := service.Login(context.Background(), "unknown", "pass")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any(), "testuser").
			Return(testUser, nil)

		_, _, err := service.Login(context.Background(), "testuser", "wrong")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
}

func TestService_Secrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenManager := token.NewMockTokenManager(ctrl)
	mockUserRepo := repository.NewMockUserRepositorier(ctrl)
	mockSecretRepo := repository.NewMockSecretRepositorier(ctrl)

	repos := &repository.Repositories{
		UserRepo:   mockUserRepo,
		SecretRepo: mockSecretRepo,
	}
	service := NewService(nil, mockTokenManager, nil, repos)

	secret := &stypes.Secret{ID: "s1", UserID: "u1", Data: []byte("data")}
	secrets := []*stypes.Secret{secret}

	t.Run("SetSecret success", func(t *testing.T) {
		mockSecretRepo.EXPECT().SetSecret(gomock.Any(), gomock.Any(), secret).Return(nil)
		err := service.SetSecret(context.Background(), secret)
		require.NoError(t, err)
	})

	t.Run("SetSecret too large", func(t *testing.T) {
		largeSecret := &stypes.Secret{ID: "s1", UserID: "u1", Data: make([]byte, maxSecretSize+1)}
		err := service.SetSecret(context.Background(), largeSecret)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrSecretTooLarge)
	})

	t.Run("GetSecret", func(t *testing.T) {
		mockSecretRepo.EXPECT().GetSecret(gomock.Any(), gomock.Any(), "u1", "s1").Return(secret, nil)
		result, err := service.GetSecret(context.Background(), "u1", "s1")
		require.NoError(t, err)
		assert.Equal(t, secret, result)
	})

	t.Run("DeleteSecret", func(t *testing.T) {
		mockSecretRepo.EXPECT().DeleteSecret(gomock.Any(), gomock.Any(), "u1", "s1").Return(nil)
		err := service.DeleteSecret(context.Background(), "u1", "s1")
		require.NoError(t, err)
	})

	t.Run("ListSecrets", func(t *testing.T) {
		mockSecretRepo.EXPECT().ListSecrets(gomock.Any(), gomock.Any(), "u1").Return(secrets, nil)
		result, err := service.ListSecrets(context.Background(), "u1")
		require.NoError(t, err)
		assert.Equal(t, secrets, result)
	})
}
