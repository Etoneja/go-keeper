package ctl

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize vault storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := getAppFromCommand(cmd)

		if err := app.service.Initialize(context.Background()); err != nil {
			return err
		}

		fmt.Printf("✅ Vault initialized successfully\n")
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Use constant, flags or git commit/tag
		fmt.Println("Vault CLI v0.1.0")
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register new user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Registration is not implemented in local mode")
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to vault",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Login is not implemented in local mode")
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with remote storage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sync is not implemented yet")
	},
}
