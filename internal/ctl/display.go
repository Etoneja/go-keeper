package ctl

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/etoneja/go-keeper/internal/buildinfo"
	"github.com/etoneja/go-keeper/internal/ctl/types"
)

const timeFormat = "2006-01-02 15:04:05"

func displaySecrets(responses []*types.LocalSecret) {
	if len(responses) == 0 {
		fmt.Println("No secrets found")
		return
	}

	fmt.Printf("%-36s %-12s %-12s %s\n", "UUID", "Type", "Name", "Last Modified")
	fmt.Println(strings.Repeat("-", 82))
	for _, resp := range responses {
		fmt.Printf("%-36s %-12s %-12s %s\n",
			resp.UUID,
			resp.Type,
			resp.Name,
			resp.LastModified.Local().Format(timeFormat))
	}
}

func displaySecret(secret *types.LocalSecret, full bool) error {
	fmt.Printf("UUID: %s\n", secret.UUID)
	fmt.Printf("Type: %s\n", secret.Type)
	fmt.Printf("Name: %s\n", secret.Name)
	fmt.Printf("Last Modified: %s\n", secret.LastModified.Local().Format(timeFormat))
	if secret.Metadata != "" {
		fmt.Printf("Metadata: %s\n", secret.Metadata)
	}
	fmt.Println()

	data, err := secret.ParseData()
	if err != nil {
		return err
	}

	switch data := data.(type) {
	case types.LoginData:
		fmt.Printf("Username: %s\n", data.Username)
		if full {
			fmt.Printf("Password: %s\n", data.Password)
		} else {
			fmt.Printf("Password: ********\n")
		}
		if data.URL != "" {
			fmt.Printf("URL: %s\n", data.URL)
		}

	case types.TextData:
		fmt.Printf("Content: %s\n", data.Content)

	case types.FileData:
		fmt.Printf("File Name: %s\n", data.FileName)
		fmt.Printf("File Size: %d bytes\n", data.FileSize)
		fmt.Printf("Content Size: %d bytes (base64)\n", len(data.Content))

	case types.CardData:
		fmt.Printf("Card Number: %s\n", data.Number)
		fmt.Printf("Card Holder: %s\n", data.Holder)
		fmt.Printf("Expiry: %s\n", data.Expiry)
		if full {
			fmt.Printf("CVV: %s\n", data.CVV)
		} else {
			fmt.Printf("CVV: ***\n")
		}

	default:
		return fmt.Errorf("unknown data type: %T", data)
	}

	return nil
}

func displayVersion() {
	commitShort := buildinfo.Commit
	if len(commitShort) > 12 {
		commitShort = commitShort[:12]
	}

	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	fmt.Println("GoKeeper - Secret Manager")
	fmt.Println("┌──────────────────┬─────────────────────────────────┐")
	fmt.Printf("│ Version          │ %-31s │\n", buildinfo.Version)
	fmt.Printf("│ Build Date       │ %-31s │\n", buildinfo.BuildTime)
	fmt.Printf("│ Git Commit       │ %-31s │\n", commitShort)
	fmt.Printf("│ Go Version       │ %-31s │\n", runtime.Version())
	fmt.Printf("│ Platform         │ %-31s │\n", platform)
	fmt.Println("└──────────────────┴─────────────────────────────────┘")
}
