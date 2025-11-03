package ctl

import "github.com/spf13/cobra"

func init() {
	// Add command flags
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

	// Get command flags
	getCmd.Flags().Bool("full", false, "Show all data including passwords/CVV")
	getCmd.Flags().String("export", "", "Export to file path")

	// Add subcommands to add command
	addCmd.AddCommand(addPasswordCmd)
	addCmd.AddCommand(addTextCmd)
	addCmd.AddCommand(addBinaryCmd)
	addCmd.AddCommand(addCardCmd)
}

func getAppFromCommand(cmd *cobra.Command) *App {
	rootCmd := cmd.Root()
	app, ok := rootCmd.Context().Value(appContextKey).(*App)
	if !ok {
		panic("app not found in command context")
	}
	return app
}

// Helper function to safely get string flags
func getStringFlag(cmd *cobra.Command, name string) string {
	value, _ := cmd.Flags().GetString(name)
	return value
}
