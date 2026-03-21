package cmd

import (
	"github.com/spf13/cobra"
	"github.com/todlehn/comyms/server/google/sheets"
)

var googleCmd = &cobra.Command{
	Use:   "google",
	Short: "Google API MCP servers",
}

var googleSheetsCmd = &cobra.Command{
	Use:   "sheets",
	Short: "Google Sheets MCP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sheets.Serve(cmd.Context())
	},
}

func init() {
	googleCmd.AddCommand(googleSheetsCmd)
}
