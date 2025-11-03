package ctl

import (
	"fmt"
	"os"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
)

func CheckFileSize(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if info.Size() > constants.MaxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)",
			info.Size(), constants.MaxFileSize)
	}

	return nil
}
