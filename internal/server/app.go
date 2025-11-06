package server

import (
	"context"
	"log"
	"net"

	"github.com/etoneja/go-keeper/internal/proto"
	"github.com/etoneja/go-keeper/internal/server/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type App struct {
	grpcServer *grpc.Server
	db         *pgxpool.Pool
}

func NewApp(cfg *Config) (*App, error) {
	// Database connection
	db, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		return nil, err
	}

	// Repositories
	repos := repository.NewRepositories()

	// Service
	svc := NewService(db, cfg, repos)

	// gRPC server with auth interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor(svc)),
	)

	// Register services
	proto.RegisterAuthServiceServer(grpcServer, NewAuthHandler(svc))
	proto.RegisterSecretServiceServer(grpcServer, NewSecretHandler(svc))

	return &App{
		grpcServer: grpcServer,
		db:         db,
	}, nil
}

func (a *App) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("Server starting on %s", addr)
	return a.grpcServer.Serve(lis)
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
	a.db.Close()
}
