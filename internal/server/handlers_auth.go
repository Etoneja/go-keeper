package server

import (
	"context"
	"errors"

	"github.com/etoneja/go-keeper/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	proto.UnimplementedAuthServiceServer
	service *Service
}

func NewAuthHandler(service *Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user, err := h.service.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &proto.RegisterResponse{}
	resp.SetUserId(user.ID)

	return resp, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, user, err := h.service.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		case errors.Is(err, ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid password")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	resp := &proto.LoginResponse{}
	resp.SetToken(token)
	resp.SetUserId(user.ID)

	return resp, nil
}
