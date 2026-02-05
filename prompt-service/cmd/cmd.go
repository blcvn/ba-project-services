package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "prompt-service",
	Short: "Prompt Service",
	Long:  `Prompt Service for managing AI prompt templates and experiments`,
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
