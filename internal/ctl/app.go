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
	cmd     *cobra.Command
	service *VaultService
}

func NewApp() (*App, error) {
	cfg, err := config.LoadCfg()
	if err != nil {
		return nil, err
	}

	app := &App{cfg: cfg}
	app.setupCommands()
	return app, nil
}

func (a *App) Run() error {
	return a.cmd.Execute()
}

func (a *App) setupCommands() {
	rootCmd := &cobra.Command{
		Use:              "vault",
		Short:            "Zero-Knowledge secret manager",
		PersistentPreRun: a.initializeService,
	}

	ctx := context.WithValue(context.Background(), appContextKey, a)
	rootCmd.SetContext(ctx)

	addCommands(rootCmd)
	a.cmd = rootCmd
}

func (a *App) initializeService(cmd *cobra.Command, args []string) {
	service := NewVaultService(a.cfg)
	a.service = service
}

func (a *App) Close() {
	if a.service != nil {
		a.service.Close()
	}
}
