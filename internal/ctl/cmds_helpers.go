package ctl

import (
	"fmt"
	"log"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/etoneja/go-keeper/internal/ctl/errs"
	"github.com/spf13/cobra"
)

func init() {
	addPasswordCmd.Flags().String("name", "", "Secret name (required)")
	addPasswordCmd.Flags().String("username", "", "Username (required)")
	addPasswordCmd.Flags().String("password", "", "Password (required)")
	addPasswordCmd.Flags().String("url", "", "URL (optional)")
	addPasswordCmd.Flags().String("metadata", "", "Metadata (optional)")
	addPasswordCmd.MarkFlagRequired("name")
	addPasswordCmd.MarkFlagRequired("username")
	addPasswordCmd.MarkFlagRequired("password")

	addTextCmd.Flags().String("name", "", "Secret name (required)")
	addTextCmd.Flags().String("content", "", "Text content (required)")
	addTextCmd.Flags().String("metadata", "", "Metadata (optional)")
	addTextCmd.MarkFlagRequired("name")
	addTextCmd.MarkFlagRequired("content")

	addBinaryCmd.Flags().String("name", "", "Secret name (required)")
	addBinaryCmd.Flags().String("file", "", "File path (required)")
	addBinaryCmd.Flags().String("metadata", "", "Metadata (optional)")
	addBinaryCmd.MarkFlagRequired("name")
	addBinaryCmd.MarkFlagRequired("file")

	addCardCmd.Flags().String("name", "", "Secret name (required)")
	addCardCmd.Flags().String("number", "", "Card number (required)")
	addCardCmd.Flags().String("holder", "", "Card holder name (required)")
	addCardCmd.Flags().String("expiry", "", "Expiry date (required)")
	addCardCmd.Flags().String("cvv", "", "CVV code (required)")
	addCardCmd.Flags().String("metadata", "", "Metadata (optional)")
	addCardCmd.MarkFlagRequired("name")
	addCardCmd.MarkFlagRequired("number")
	addCardCmd.MarkFlagRequired("holder")
	addCardCmd.MarkFlagRequired("expiry")
	addCardCmd.MarkFlagRequired("cvv")

	getCmd.Flags().Bool("full", false, "Show all data including passwords/CVV")
	getCmd.Flags().String("export", "", "Export to file path")

	addCmd.AddCommand(addPasswordCmd)
	addCmd.AddCommand(addTextCmd)
	addCmd.AddCommand(addBinaryCmd)
	addCmd.AddCommand(addCardCmd)
}

func addCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(syncCmd)
}

func getAppFromCommand(cmd *cobra.Command) *App {
	rootCmd := cmd.Root()
	app, ok := rootCmd.Context().Value(appContextKey).(*App)
	if !ok {
		log.Fatal("app not found in command context")
	}
	return app
}

func getStringFlag(cmd *cobra.Command, name string) string {
	value, _ := cmd.Flags().GetString(name)
	return value
}

func withErrorHandling(fn func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := fn(cmd, args)
		short := cmd.Short
		if err != nil {
			emoji := constants.EmojiError
			if errs.IsNotFound(err) {
				emoji = constants.EmojiWarning
			}
			fmt.Printf("%s Failed to %s: %v\n", emoji, short, err)
			return
		}
	}
}
