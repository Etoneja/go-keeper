package ctl

import (
	"context"
	"fmt"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/etoneja/go-keeper/internal/ctl/errs"
	"github.com/spf13/cobra"
)

// Add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new secret",
}

var addPasswordCmd = createAddCommand(constants.TypePassword)
var addTextCmd = createAddCommand(constants.TypeText)
var addBinaryCmd = createAddCommand(constants.TypeBinary)
var addCardCmd = createAddCommand(constants.TypeCard)

// Get command
var getCmd = &cobra.Command{
	Use:   "get [uuid]",
	Short: "Get secret by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		uuid := args[0]
		full, _ := cmd.Flags().GetBool("full")
		exportPath, _ := cmd.Flags().GetString("export")

		app := getAppFromCommand(cmd)

		secret, err := app.service.GetSecret(context.Background(), uuid)
		if err != nil {
			if errs.IsNotFound(err) {
				fmt.Printf("%s Secret %s not found\n", constants.EmojiWarning, uuid)
				return nil
			}
			fmt.Printf("%s Failed to get secret %s: %v\n", constants.EmojiError, uuid, err)
			return nil
		}

		if exportPath != "" {
			return exportSecret(secret, exportPath)
		}

		return displaySecret(secret, full)
	},
}

// Delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [uuid]",
	Short: "Delete secret by UUID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		uuid := args[0]

		app := getAppFromCommand(cmd)

		err := app.service.DeleteSecret(context.Background(), uuid)
		if err != nil {
			if errs.IsNotFound(err) {
				fmt.Printf("%s Secret %s not found\n", constants.EmojiWarning, uuid)
				return nil
			}
			fmt.Printf("%s Failed to delete secret %s: %v\n", constants.EmojiError, uuid, err)
			return nil
		}
		fmt.Printf("%s Secret %s deleted successfully\n", constants.EmojiSuccess, uuid)
		return nil
	},
}

// List command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := getAppFromCommand(cmd)

		secrets, err := app.service.ListSecrets(context.Background())
		if err != nil {
			fmt.Printf("%s Failed to list secrets: %v\n", constants.EmojiError, err)
			return nil
		}

		displaySecrets(secrets)
		return nil
	},
}
