package ctl

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/etoneja/go-keeper/internal/ctl/types"
)

func exportSecret(secret *types.Secret, exportPath string) error {
	data, err := secret.ParseData()
	if err != nil {
		return err
	}
	switch data := data.(type) {
	case types.FileData:
		content, err := base64.StdEncoding.DecodeString(data.Content)
		if err != nil {
			return fmt.Errorf("failed to decode base64 content: %w", err)
		}

		if err := os.WriteFile(exportPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("%s File exported to: %s\n", constants.EmojiSuccess, exportPath)

	default:
		jsonData, err := json.MarshalIndent(secret, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		if err := os.WriteFile(exportPath, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("%s Secret exported to: %s\n", constants.EmojiSuccess, exportPath)
	}

	return nil
}
