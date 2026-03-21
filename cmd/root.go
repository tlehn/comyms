package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "comyms",
	Short: "A collection of MCP servers",
}

func init() {
	rootCmd.AddCommand(googleCmd)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
