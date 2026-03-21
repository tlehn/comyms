package sheets

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	googleauth "github.com/todlehn/comyms/server/google"
)

// Serve creates the Google Sheets MCP server, registers all tools, and runs
// the stdio transport. It blocks until the transport is closed or an error occurs.
func Serve(ctx context.Context) error {
	sheetsSvc, err := googleauth.NewSheetsService(ctx)
	if err != nil {
		return err
	}

	driveSvc, err := googleauth.NewDriveService(ctx)
	if err != nil {
		return err
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "google-sheets-mcp",
		Version: "0.1.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_spreadsheet",
		Description: "Read data from a Google Sheets spreadsheet",
	}, ReadSpreadsheet(sheetsSvc))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_sheets",
		Description: "List all sheets (tabs) in a Google Sheets spreadsheet",
	}, ListSheets(sheetsSvc))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_spreadsheets",
		Description: "List Google Sheets spreadsheets available to the service account",
	}, ListSpreadsheets(driveSvc))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_spreadsheet_filtered",
		Description: "Read rows from a sheet, filtered by a column value",
	}, ReadSpreadsheetFiltered(sheetsSvc))

	return server.Run(ctx, &mcp.StdioTransport{})
}
