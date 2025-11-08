package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		// log.Printf("gRPC request started: %s", info.FullMethod)

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		statusCode := status.Code(err)

		if err != nil {
			log.Printf("gRPC request failed: %s, duration: %s, status: %s, error: %v",
				info.FullMethod, duration, statusCode, err)
		} else {
			log.Printf("gRPC request completed: %s, duration: %s, status: %s",
				info.FullMethod, duration, statusCode)
		}

		return resp, err
	}
}
