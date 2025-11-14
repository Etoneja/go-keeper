package ctl

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/etoneja/go-keeper/internal/ctl/types"
)

func exportSecret(secret *types.LocalSecret, exportPath string) error {
	if _, err := os.Stat(exportPath); err == nil {
		return fmt.Errorf("file already exists: %s", exportPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	data, err := secret.ParseData()
	if err != nil {
		return err
	}

	var content []byte
	switch data := data.(type) {
	case types.FileData:
		content, err = base64.StdEncoding.DecodeString(data.Content)
		if err != nil {
			return fmt.Errorf("failed to decode base64 content: %w", err)
		}
	default:
		content, err = json.MarshalIndent(secret, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	}

	tmpPath := exportPath + ".tmp"

	if err := os.WriteFile(tmpPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tmpPath, exportPath); err != nil {
		if removeErr := os.Remove(tmpPath); removeErr != nil {
			log.Printf("failed to remove temporary file %s: %v", tmpPath, removeErr)
		}
		return fmt.Errorf("failed to replace file: %w", err)
	}

	return nil
}
