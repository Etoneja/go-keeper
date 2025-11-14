package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/etoneja/go-keeper/internal/proto"
	"github.com/etoneja/go-keeper/internal/server/repository"
	"github.com/etoneja/go-keeper/internal/server/token"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type App struct {
	grpcServer *grpc.Server
	db         *pgxpool.Pool
}

func NewApp() (*App, error) {
	cfg, err := LoadCfg()
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories()
	jwtManager := token.NewJWTManager(cfg.JWTSecret, time.Hour)
	svc := NewService(db, jwtManager, nil, repos)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor(),
			AuthInterceptor(svc),
		),
	)

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
