package server

import (
	"context"
	"testing"

	"github.com/etoneja/go-keeper/internal/server/token"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenManager := token.NewMockTokenManager(ctrl)
	mockService := &Service{
		tokenManager: mockTokenManager,
	}

	interceptor := AuthInterceptor(mockService)

	t.Run("skip auth for login", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			return "response", nil
		}

		resp, err := interceptor(context.Background(), "request", &grpc.UnaryServerInfo{
			FullMethod: "/gokeeper.AuthService/Login",
		}, handler)

		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
	})

	t.Run("skip auth for register", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			return "response", nil
		}

		resp, err := interceptor(context.Background(), "request", &grpc.UnaryServerInfo{
			FullMethod: "/gokeeper.AuthService/Register",
		}, handler)

		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
	})

	t.Run("missing metadata", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			return "response", nil
		}

		resp, err := interceptor(context.Background(), "request", &grpc.UnaryServerInfo{
			FullMethod: "/gokeeper.OtherService/Method",
		}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("missing authorization header", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			return "response", nil
		}

		md := metadata.MD{}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := interceptor(ctx, "request", &grpc.UnaryServerInfo{
			FullMethod: "/gokeeper.OtherService/Method",
		}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("valid token", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			userID, err := getUserIDFromContext(ctx)
			assert.NoError(t, err)
			assert.Equal(t, "user123", userID)
			return "response", nil
		}

		md := metadata.MD{"authorization": []string{"valid-token"}}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockTokenManager.EXPECT().ValidateToken("valid-token").Return("user123", nil)

		resp, err := interceptor(ctx, "request", &grpc.UnaryServerInfo{
			FullMethod: "/gokeeper.OtherService/Method",
		}, handler)

		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
	})

	t.Run("invalid token", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			return "response", nil
		}

		md := metadata.MD{"authorization": []string{"invalid-token"}}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockTokenManager.EXPECT().ValidateToken("invalid-token").Return("", assert.AnError)

		resp, err := interceptor(ctx, "request", &grpc.UnaryServerInfo{
			FullMethod: "/gokeeper.OtherService/Method",
		}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})
}

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("userID exists", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user123")
		userID, err := getUserIDFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "user123", userID)
	})

	t.Run("userID missing", func(t *testing.T) {
		userID, err := getUserIDFromContext(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "", userID)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, 123)
		userID, err := getUserIDFromContext(ctx)
		assert.Error(t, err)
		assert.Equal(t, "", userID)
	})
}
