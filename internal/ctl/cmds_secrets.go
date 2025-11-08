package ctl

import (
	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new secret",
}

var addPasswordCmd = createSecretAddCommand(constants.SecretTypePassword)
var addTextCmd = createSecretAddCommand(constants.SecretTypeText)
var addBinaryCmd = createSecretAddCommand(constants.SecretTypeBinary)
var addCardCmd = createSecretAddCommand(constants.SecretTypeCard)

var getCmd = &cobra.Command{
	Use:   "get [uuid]",
	Short: "Get secret by UUID",
	Args:  cobra.ExactArgs(1),
	Run:   withErrorHandling(createSecretGetCommand()),
}

var deleteCmd = &cobra.Command{
	Use:   "delete [uuid]",
	Short: "Delete secret by UUID",
	Args:  cobra.ExactArgs(1),
	Run:   withErrorHandling(createSecretDeleteCommand()),
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets",
	Run:   withErrorHandling(createSecretsListCommand()),
}
