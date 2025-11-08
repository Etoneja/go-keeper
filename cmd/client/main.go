package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/etoneja/go-keeper/internal/ctl/client"
	"github.com/etoneja/go-keeper/internal/ctl/config"
	"github.com/etoneja/go-keeper/internal/ctl/types"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.LoadCfg()
	if err != nil {
		panic(err)
	}
	cli := client.NewGRPCClient(cfg)
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect
	if err := cli.Connect(ctx); err != nil {
		log.Fatal("Connect failed:", err)
	}

	// Register & Login
	fmt.Println("=== Auth ===")
	if _, err := cli.Register(ctx); err != nil {
		log.Printf("Register failed: %v", err)
	}

	if err := cli.Login(ctx); err != nil {
		log.Fatal("Login failed:", err)
	}
	fmt.Println("Logged in successfully")

	// Work with secrets (auto-auth)
	fmt.Println("\n=== Secrets ===")

	// Set secret
	secret := &types.RemoteSecret{
		UUID:         "asdf",
		LastModified: time.Now(),
		Hash:         "sdf",
		Data:         []byte("my data"),
	}
	if err := cli.SetSecret(ctx, secret); err != nil {
		log.Fatal("SetSecret failed:", err)
	}
	fmt.Println("Secret set")

	// List secrets
	secrets, err := cli.ListSecrets(ctx)
	if err != nil {
		log.Fatal("ListSecrets failed:", err)
	}

	fmt.Printf("Found %d secrets:\n", len(secrets))
	for i, secret := range secrets {
		fmt.Printf("%d. %s (hash: %s)\n", i+1, secret.GetUUID(), secret.GetHash())
	}
}
