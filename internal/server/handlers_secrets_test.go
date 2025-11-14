package server

import (
	"context"
	"testing"
	"time"

	"github.com/etoneja/go-keeper/internal/proto"
	"github.com/etoneja/go-keeper/internal/server/stypes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSecretHandler_SetSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServicer(ctrl)
	handler := NewSecretHandler(mockService)

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		reqSecret := &proto.Secret{}
		reqSecret.SetId("secret1")
		reqSecret.SetData([]byte("secret data"))
		reqSecret.SetHash("hash123")
		reqSecret.SetLastModified(timestamppb.New(time.Now()))

		req := &proto.SetSecretRequest{}
		req.SetSecret(reqSecret)

		mockService.EXPECT().SetSecret(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := handler.SetSecret(ctx, req)
		require.NoError(t, err)
		assert.True(t, resp.GetSuccess())
	})

	t.Run("unauthorized", func(t *testing.T) {
		reqSecret := &proto.Secret{}
		reqSecret.SetId("secret1")

		req := &proto.SetSecretRequest{}
		req.SetSecret(reqSecret)

		resp, err := handler.SetSecret(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		reqSecret := &proto.Secret{}
		reqSecret.SetId("secret1")

		req := &proto.SetSecretRequest{}
		req.SetSecret(reqSecret)

		mockService.EXPECT().SetSecret(gomock.Any(), gomock.Any()).Return(assert.AnError)

		resp, err := handler.SetSecret(ctx, req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}

func TestSecretHandler_GetSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServicer(ctrl)
	handler := NewSecretHandler(mockService)

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		req := &proto.GetSecretRequest{}
		req.SetSecretId("secret1")

		now := time.Now()
		secret := &stypes.Secret{
			ID:           "secret1",
			UserID:       "user123",
			Data:         []byte("secret data"),
			Hash:         "hash123",
			LastModified: now,
		}

		mockService.EXPECT().GetSecret(gomock.Any(), "user123", "secret1").Return(secret, nil)

		resp, err := handler.GetSecret(ctx, req)
		require.NoError(t, err)

		respSecret := resp.GetSecret()
		assert.Equal(t, "secret1", respSecret.GetId())
		assert.Equal(t, "hash123", respSecret.GetHash())
		assert.Equal(t, []byte("secret data"), respSecret.GetData())
		assert.True(t, respSecret.GetLastModified().AsTime().Equal(now))
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := &proto.GetSecretRequest{}
		req.SetSecretId("secret1")

		resp, err := handler.GetSecret(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		req := &proto.GetSecretRequest{}
		req.SetSecretId("secret1")

		mockService.EXPECT().GetSecret(gomock.Any(), "user123", "secret1").Return(nil, assert.AnError)

		resp, err := handler.GetSecret(ctx, req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}

func TestSecretHandler_DeleteSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServicer(ctrl)
	handler := NewSecretHandler(mockService)

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		req := &proto.DeleteSecretRequest{}
		req.SetSecretId("secret1")

		mockService.EXPECT().DeleteSecret(gomock.Any(), "user123", "secret1").Return(nil)

		resp, err := handler.DeleteSecret(ctx, req)
		require.NoError(t, err)
		assert.True(t, resp.GetSuccess())
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := &proto.DeleteSecretRequest{}
		req.SetSecretId("secret1")

		resp, err := handler.DeleteSecret(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		req := &proto.DeleteSecretRequest{}
		req.SetSecretId("secret1")

		mockService.EXPECT().DeleteSecret(gomock.Any(), "user123", "secret1").Return(assert.AnError)

		resp, err := handler.DeleteSecret(ctx, req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}

func TestSecretHandler_ListSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServicer(ctrl)
	handler := NewSecretHandler(mockService)

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		req := &proto.ListSecretsRequest{}

		now := time.Now()
		secrets := []*stypes.Secret{
			{
				ID:           "secret1",
				Hash:         "hash1",
				LastModified: now,
			},
			{
				ID:           "secret2",
				Hash:         "hash2",
				LastModified: now.Add(-time.Hour),
			},
		}

		mockService.EXPECT().ListSecrets(gomock.Any(), "user123").Return(secrets, nil)

		resp, err := handler.ListSecrets(ctx, req)
		require.NoError(t, err)

		respSecrets := resp.GetSecrets()
		require.Len(t, respSecrets, 2)

		assert.Equal(t, "secret1", respSecrets[0].GetId())
		assert.Equal(t, "hash1", respSecrets[0].GetHash())
		assert.True(t, respSecrets[0].GetLastModified().AsTime().Equal(now))

		assert.Equal(t, "secret2", respSecrets[1].GetId())
		assert.Equal(t, "hash2", respSecrets[1].GetHash())
		assert.True(t, respSecrets[1].GetLastModified().AsTime().Equal(now.Add(-time.Hour)))
	})

	t.Run("empty list", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		req := &proto.ListSecretsRequest{}

		mockService.EXPECT().ListSecrets(gomock.Any(), "user123").Return([]*stypes.Secret{}, nil)

		resp, err := handler.ListSecrets(ctx, req)
		require.NoError(t, err)
		assert.Empty(t, resp.GetSecrets())
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := &proto.ListSecretsRequest{}

		resp, err := handler.ListSecrets(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")

		req := &proto.ListSecretsRequest{}

		mockService.EXPECT().ListSecrets(gomock.Any(), "user123").Return(nil, assert.AnError)

		resp, err := handler.ListSecrets(ctx, req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}
