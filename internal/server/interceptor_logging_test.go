package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoggingInterceptor(t *testing.T) {
	interceptor := LoggingInterceptor()

	t.Run("successful request", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			return "response", nil
		}

		resp, err := interceptor(context.Background(), "request", &grpc.UnaryServerInfo{
			FullMethod: "/service.Method",
		}, handler)

		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
	})

	t.Run("failed request", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		resp, err := interceptor(context.Background(), "request", &grpc.UnaryServerInfo{
			FullMethod: "/service.Method",
		}, handler)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("with timeout", func(t *testing.T) {
		handler := func(ctx context.Context, req any) (any, error) {
			time.Sleep(10 * time.Millisecond)
			return "delayed", nil
		}

		resp, err := interceptor(context.Background(), "request", &grpc.UnaryServerInfo{
			FullMethod: "/service.Method",
		}, handler)

		assert.NoError(t, err)
		assert.Equal(t, "delayed", resp)
	})
}
