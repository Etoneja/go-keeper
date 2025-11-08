package client

import (
	"context"
	"errors"

	"github.com/etoneja/go-keeper/internal/ctl/types"
	"github.com/etoneja/go-keeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotConnected = errors.New("not connected to server")
)

type Client struct {
	conn         *grpc.ClientConn
	authClient   proto.AuthServiceClient
	secretClient proto.SecretServiceClient

	serverAddress string

	login    string
	password string

	token string
}

func NewGRPCClient(serverAddress string, login string, password string) *Client {

	return &Client{
		serverAddress: serverAddress,
		login:         login,
		password:      password,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	conn, err := grpc.NewClient(c.serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.conn = conn
	c.authClient = proto.NewAuthServiceClient(conn)
	c.secretClient = proto.NewSecretServiceClient(conn)

	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Login(ctx context.Context) error {
	req := &proto.LoginRequest{}
	req.SetLogin(c.login)
	req.SetPassword(c.password)

	resp, err := c.authClient.Login(ctx, req)
	if err != nil {
		return err
	}

	c.token = resp.GetToken()
	return nil
}

func (c *Client) Register(ctx context.Context) (string, error) {
	req := &proto.RegisterRequest{}
	req.SetLogin(c.login)
	req.SetPassword(c.password)

	resp, err := c.authClient.Register(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.GetUserId(), nil
}

func (c *Client) ensureAuth(ctx context.Context) error {
	if c.token != "" {
		return nil
	}

	return c.Login(ctx)
}

func (c *Client) withAuthRetry(ctx context.Context, fn func(context.Context) error) error {
	if err := c.ensureAuth(ctx); err != nil {
		return err
	}

	authCtx := c.createAuthContext(ctx)
	err := fn(authCtx)

	if isUnauthorizedError(err) {
		c.token = ""
		if err := c.ensureAuth(ctx); err != nil {
			return err
		}

		authCtx = c.createAuthContext(ctx)
		return fn(authCtx)
	}

	return err
}

func (c *Client) createAuthContext(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx,
		metadata.Pairs("authorization", c.token),
	)
}

// Secret methods with auto-auth
func (c *Client) SetSecret(ctx context.Context, secret *types.RemoteSecret) error {
	reqSecret := &proto.Secret{}
	reqSecret.SetId(secret.UUID)
	reqSecret.SetLastModified(timestamppb.New(secret.LastModified))
	reqSecret.SetHash(secret.Hash)
	reqSecret.SetData(secret.Data)

	req := &proto.SetSecretRequest{}
	req.SetSecret(reqSecret)

	return c.withAuthRetry(ctx, func(authCtx context.Context) error {
		_, err := c.secretClient.SetSecret(authCtx, req)
		return err
	})
}

func (c *Client) GetSecret(ctx context.Context, secretID string) (*types.RemoteSecret, error) {
	var resp *proto.GetSecretResponse

	req := &proto.GetSecretRequest{}
	req.SetSecretId(secretID)

	err := c.withAuthRetry(ctx, func(authCtx context.Context) error {
		var err error
		resp, err = c.secretClient.GetSecret(authCtx, req)
		return err
	})
	if err != nil {
		return nil, err
	}

	secretResp := resp.GetSecret()
	secret := &types.RemoteSecret{
		UUID:         secretResp.GetId(),
		LastModified: secretResp.GetLastModified().AsTime(),
		Hash:         secretResp.GetHash(),
		Data:         secretResp.GetData(),
	}

	return secret, nil
}

func (c *Client) DeleteSecret(ctx context.Context, secretID string) error {
	req := &proto.DeleteSecretRequest{}
	req.SetSecretId(secretID)

	return c.withAuthRetry(ctx, func(authCtx context.Context) error {
		_, err := c.secretClient.DeleteSecret(authCtx, req)
		return err
	})
}

func (c *Client) ListSecrets(ctx context.Context) ([]*types.RemoteSecret, error) {
	var resp *proto.ListSecretsResponse
	err := c.withAuthRetry(ctx, func(authCtx context.Context) error {
		var err error
		resp, err = c.secretClient.ListSecrets(authCtx, &proto.ListSecretsRequest{})
		return err
	})
	if err != nil {
		return nil, err
	}

	secretsResp := resp.GetSecrets()
	secrets := make([]*types.RemoteSecret, len(secretsResp))
	for i, secretResp := range secretsResp {
		secrets[i] = &types.RemoteSecret{
			UUID:         secretResp.GetId(),
			LastModified: secretResp.GetLastModified().AsTime(),
			Hash:         secretResp.GetHash(),
		}
	}

	return secrets, nil
}

func isUnauthorizedError(err error) bool {
	return err != nil && err.Error() == "rpc error: code = Unauthenticated desc = invalid token"
}
