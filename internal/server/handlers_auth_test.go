package server

import (
	"context"
	"errors"
	"testing"

	"github.com/etoneja/go-keeper/internal/proto"
	"github.com/etoneja/go-keeper/internal/server/stypes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServicer(ctrl)
	handler := NewAuthHandler(mockService)

	t.Run("success", func(t *testing.T) {
		req := &proto.RegisterRequest{}
		req.SetLogin("testuser")
		req.SetPassword("password123")

		mockService.EXPECT().Register(gomock.Any(), "testuser", "password123").
			Return(&stypes.User{ID: "user123"}, nil)

		resp, err := handler.Register(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "user123", resp.GetUserId())
	})

	t.Run("service error", func(t *testing.T) {
		req := &proto.RegisterRequest{}
		req.SetLogin("testuser")
		req.SetPassword("password123")

		mockService.EXPECT().Register(gomock.Any(), "testuser", "password123").
			Return(nil, errors.New("registration failed"))

		resp, err := handler.Register(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServicer(ctrl)
	handler := NewAuthHandler(mockService)

	t.Run("success", func(t *testing.T) {
		req := &proto.LoginRequest{}
		req.SetLogin("testuser")
		req.SetPassword("password123")

		mockService.EXPECT().Login(gomock.Any(), "testuser", "password123").
			Return("token123", &stypes.User{ID: "user123"}, nil)

		resp, err := handler.Login(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "token123", resp.GetToken())
		assert.Equal(t, "user123", resp.GetUserId())
	})

	t.Run("user not found", func(t *testing.T) {
		req := &proto.LoginRequest{}
		req.SetLogin("unknown")
		req.SetPassword("password123")

		mockService.EXPECT().Login(gomock.Any(), "unknown", "password123").
			Return("", nil, ErrUserNotFound)

		resp, err := handler.Login(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("invalid credentials", func(t *testing.T) {
		req := &proto.LoginRequest{}
		req.SetLogin("testuser")
		req.SetPassword("wrongpassword")

		mockService.EXPECT().Login(gomock.Any(), "testuser", "wrongpassword").
			Return("", nil, ErrInvalidCredentials)

		resp, err := handler.Login(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
		assert.Contains(t, err.Error(), "invalid password")
	})

	t.Run("internal error", func(t *testing.T) {
		req := &proto.LoginRequest{}
		req.SetLogin("testuser")
		req.SetPassword("password123")

		mockService.EXPECT().Login(gomock.Any(), "testuser", "password123").
			Return("", nil, errors.New("database error"))

		resp, err := handler.Login(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Contains(t, err.Error(), "internal error")
	})
}
