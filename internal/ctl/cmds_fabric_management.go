package ctl

import (
	"context"
	"github.com/spf13/cobra"
)

func createVersionHandler() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		displayVersion()
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
