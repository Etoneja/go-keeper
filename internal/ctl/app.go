package ctl

import (
	"context"

	"github.com/etoneja/go-keeper/internal/ctl/config"
	"github.com/spf13/cobra"
)

type contextKey string

const (
	appContextKey contextKey = "app"
)

type App struct {
	cfg     *config.Config
	service *VaultService
}

func NewApp() (*App, error) {
	cfg, err := config.LoadCfg()
	if err != nil {
		return nil, err
	}
	return &App{
		cfg: cfg,
	}, nil
}

func (a *App) GetCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "vault",
		Short:             "Zero-Knowledge secret manager",
		PersistentPreRunE: a.initializeService,
	}

	ctx := context.WithValue(context.Background(), appContextKey, a)
	rootCmd.SetContext(ctx)

	addCommands(rootCmd)
	return rootCmd
}

func (a *App) initializeService(cmd *cobra.Command, args []string) error {
	service := NewVaultService(a.cfg)
	a.service = service
	return nil
}

func (a *App) Close() {
	if a.service != nil {
		a.service.Close()
	}
}

func addCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(syncCmd)
}
