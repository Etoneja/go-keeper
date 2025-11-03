package ctl

import (
	"fmt"

	"github.com/etoneja/go-keeper/internal/ctl/types"
)

func displaySecrets(responses []*types.Secret) {
	if len(responses) == 0 {
		fmt.Println("No secrets found")
		return
	}

	fmt.Printf("%-36s %-20s %-12s %s\n", "UUID", "Last Modified", "Type", "Name")
	fmt.Println("-------------------------------------------------------------------------------")
	for _, resp := range responses {
		fmt.Printf("%-36s %-20s %-12s %s\n",
			resp.UUID,
			resp.LastModified.Format("2006-01-02 15:04:05"),
			resp.Type,
			resp.Name)
	}
}

func displaySecret(secret *types.Secret, full bool) error {
	fmt.Printf("UUID: %s\n", secret.UUID)
	fmt.Printf("Name: %s\n", secret.Name)
	fmt.Printf("Type: %s\n", secret.Type)
	fmt.Printf("Last Modified: %s\n", secret.LastModified.Format("2006-01-02 15:04:05"))
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
