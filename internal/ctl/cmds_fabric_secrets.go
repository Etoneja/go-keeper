package ctl

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/etoneja/go-keeper/internal/ctl/types"
	"github.com/spf13/cobra"
)

func createSecretAddCommand(secretType string) *cobra.Command {
	return &cobra.Command{
		Use:   secretType,
		Short: fmt.Sprintf("Add %s secret", secretType),
		Run:   withErrorHandling(createSecretHandler(secretType)),
	}
}

func createSecretHandler(secretType string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		base := types.BaseSecret{
			Type:     secretType,
			Name:     getStringFlag(cmd, "name"),
			Metadata: getStringFlag(cmd, "metadata"),
		}

		var data types.SecretData
		switch secretType {
		case constants.SecretTypePassword:
			data = types.LoginData{
				Username: getStringFlag(cmd, "username"),
				Password: getStringFlag(cmd, "password"),
				URL:      getStringFlag(cmd, "url"),
			}
		case constants.SecretTypeText:
			data = types.TextData{
				Content: getStringFlag(cmd, "content"),
			}
		case constants.SecretTypeBinary:
			filePath := getStringFlag(cmd, "file")
			if err := CheckFileSize(filePath); err != nil {
				return err
			}
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			data = types.FileData{
				FileName: filepath.Base(filePath),
				FileSize: int64(len(content)),
				Content:  base64.StdEncoding.EncodeToString(content),
			}
		case constants.SecretTypeCard:
			data = types.CardData{
				Number: getStringFlag(cmd, "number"),
				Holder: getStringFlag(cmd, "holder"),
				Expiry: getStringFlag(cmd, "expiry"),
				CVV:    getStringFlag(cmd, "cvv"),
			}
		default:
			return fmt.Errorf("unsupported secret type: %s", secretType)
		}

		app := getAppFromCommand(cmd)

		secret, err := types.NewSecretModel(base, data, app.service.cryptor)
		if err != nil {
			return err
		}

		createdSecret, err := app.service.CreateLocalSecret(context.Background(), secret)
		if err != nil {
			return err
		}

		err = displaySecret(createdSecret, false)
		if err != nil {
			return err
		}

		return nil
	}
}

func createSecretGetCommand() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		uuid := args[0]
		full, _ := cmd.Flags().GetBool("full")
		exportPath, _ := cmd.Flags().GetString("export")

		app := getAppFromCommand(cmd)

		secret, err := app.service.GetLocalSecret(context.Background(), uuid)
		if err != nil {
			return err
		}

		if exportPath != "" {
			err := exportSecret(secret, exportPath)
			if err != nil {
				return err
			}
			return nil
		}

		err = displaySecret(secret, full)
		if err != nil {
			return err
		}

		return nil
	}
}

func createSecretDeleteCommand() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		uuid := args[0]

		app := getAppFromCommand(cmd)

		err := app.service.DeleteLocalSecret(context.Background(), uuid)
		if err != nil {
			return err
		}

		return nil
	}
}

func createSecretsListCommand() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		app := getAppFromCommand(cmd)

		secrets, err := app.service.ListLocalSecrets(context.Background())
		if err != nil {
			return err
		}

		displaySecrets(secrets)

		return nil
	}
}
