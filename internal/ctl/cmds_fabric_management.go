package ctl

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func createVersionHandler() func(cmd *cobra.Command, args []string) {
	// TODO: use constant or build flags
	return func(cmd *cobra.Command, args []string) {
		fmt.Println("GoKeeper CLI v0.1.0")
	}
}

func createInitializeHandler() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		app := getAppFromCommand(cmd)
		err := app.service.Initialize(context.Background())
		if err != nil {
			return err
		}
		return err
	}
}

func createRegisterHandler() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		app := getAppFromCommand(cmd)
		err := app.service.Register(context.Background())
		if err != nil {
			return err
		}
		return nil
	}
}

func createSyncHandler() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		app := getAppFromCommand(cmd)
		err := app.service.SyncSecrets(context.Background())
		if err != nil {
			return err
		}
		return nil
	}
}
