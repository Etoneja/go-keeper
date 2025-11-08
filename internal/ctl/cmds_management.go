package ctl

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run:   createVersionHandler(),
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize local storage",
	Run:   withErrorHandling(createInitializeHandler()),
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register new user",
	Run:   withErrorHandling(createRegisterHandler()),
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with remote storage",
	Run:   withErrorHandling(createSyncHandler()),
}
