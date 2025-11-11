package server

import (
	"context"
	"log"

	"github.com/etoneja/go-keeper/internal/proto"
	"github.com/etoneja/go-keeper/internal/server/stypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SecretHandler struct {
	proto.UnimplementedSecretServiceServer
	service Servicer
}

func NewSecretHandler(service Servicer) *SecretHandler {
	return &SecretHandler{service: service}
}

func (h *SecretHandler) SetSecret(ctx context.Context, req *proto.SetSecretRequest) (*proto.SetSecretResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	reqSecret := req.GetSecret()

	secret := &stypes.Secret{
		ID:           reqSecret.GetId(),
		UserID:       userID,
		Data:         reqSecret.GetData(),
		Hash:         reqSecret.GetHash(),
		LastModified: reqSecret.GetLastModified().AsTime(),
	}

	err = h.service.SetSecret(ctx, secret)
	if err != nil {
		log.Printf("SetSecret failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to set secret")
	}

	resp := &proto.SetSecretResponse{}
	resp.SetSuccess(true)

	return resp, nil
}

func (h *SecretHandler) GetSecret(ctx context.Context, req *proto.GetSecretRequest) (*proto.GetSecretResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	secret, err := h.service.GetSecret(ctx, userID, req.GetSecretId())
	if err != nil {
		// TODO: return NotFound if not found
		log.Printf("GetSecret failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to get secret")
	}

	respSecret := &proto.Secret{}

	respSecret.SetId(secret.ID)
	respSecret.SetHash(secret.Hash)
	respSecret.SetLastModified(timestamppb.New(secret.LastModified))
	respSecret.SetData(secret.Data)

	resp := &proto.GetSecretResponse{}
	resp.SetSecret(respSecret)

	return resp, nil
}

func (h *SecretHandler) DeleteSecret(ctx context.Context, req *proto.DeleteSecretRequest) (*proto.DeleteSecretResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	err = h.service.DeleteSecret(ctx, userID, req.GetSecretId())
	if err != nil {
		// TODO: return NotFound if not found
		log.Printf("DeleteSecret failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete secret")
	}

	resp := &proto.DeleteSecretResponse{}
	resp.SetSuccess(true)

	return resp, nil
}

func (h *SecretHandler) ListSecrets(ctx context.Context, req *proto.ListSecretsRequest) (*proto.ListSecretsResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	secrets, err := h.service.ListSecrets(ctx, userID)
	if err != nil {
		log.Printf("ListSecrets failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to retrieve secrets")
	}

	respSecrets := make([]*proto.Secret, len(secrets))
	for i, secret := range secrets {
		respSecret := &proto.Secret{}

		respSecret.SetId(secret.ID)
		respSecret.SetHash(secret.Hash)
		respSecret.SetLastModified(timestamppb.New(secret.LastModified))

		respSecrets[i] = respSecret
	}

	resp := &proto.ListSecretsResponse{}
	resp.SetSecrets(respSecrets)

	return resp, nil
}
