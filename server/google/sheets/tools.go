package sheets

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/api/drive/v3"
	sheetsapi "google.golang.org/api/sheets/v4"
)

type ReadSpreadsheetParams struct {
	SpreadsheetID string `json:"spreadsheet_id"`
	Range         string `json:"range"`
}

// ReadSpreadsheet returns an MCP tool handler that reads raw cell values from a spreadsheet range.
// The range must be in A1 notation (e.g., "Sheet1!A1:C10" or "A1:B5").
//
// The response is tab-separated text, one row per line. All cell types are coerced to strings.
// If the range contains no data, returns a result with the text "No data found." (not an error).
//
// Errors:
//   - Returns the Sheets API error directly on network failure, invalid spreadsheet ID,
//     malformed range, or insufficient permissions (403).
//   - A range referencing a non-existent sheet name results in a 400 from the Sheets API.
func ReadSpreadsheet(srv *sheetsapi.Service) func(context.Context, *mcp.CallToolRequest, ReadSpreadsheetParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ReadSpreadsheetParams) (*mcp.CallToolResult, any, error) {
		resp, err := srv.Spreadsheets.Values.Get(args.SpreadsheetID, args.Range).Context(ctx).Do()
		if err != nil {
			return nil, nil, err
		}

		if len(resp.Values) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "No data found."},
				},
			}, nil, nil
		}

		var b strings.Builder
		for _, row := range resp.Values {
			for i, cell := range row {
				if i > 0 {
					b.WriteString("\t")
				}
				b.WriteString(cellToString(cell))
			}
			b.WriteString("\n")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: b.String()},
			},
		}, nil, nil
	}
}

type ListSheetsParams struct {
	SpreadsheetID string `json:"spreadsheet_id"`
}

// ListSheets returns an MCP tool handler that lists all sheets (tabs) within a spreadsheet.
// Each line in the result contains the sheet title and its numeric ID, e.g. "Sheet1 (ID: 0)".
//
// If the spreadsheet exists but contains no sheets (unlikely in practice), returns
// "No sheets found." as a successful result.
//
// Errors:
//   - Returns the Sheets API error on network failure, invalid spreadsheet ID, or
//     insufficient permissions (403). A non-existent spreadsheet ID returns 404.
func ListSheets(srv *sheetsapi.Service) func(context.Context, *mcp.CallToolRequest, ListSheetsParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListSheetsParams) (*mcp.CallToolResult, any, error) {
		spreadsheet, err := srv.Spreadsheets.Get(args.SpreadsheetID).Context(ctx).Do()
		if err != nil {
			return nil, nil, err
		}

		if len(spreadsheet.Sheets) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "No sheets found."},
				},
			}, nil, nil
		}

		var lines []string
		for _, s := range spreadsheet.Sheets {
			lines = append(lines, fmt.Sprintf("%s (ID: %d)", s.Properties.Title, s.Properties.SheetId))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: strings.Join(lines, "\n")},
			},
		}, nil, nil
	}
}

type ListSpreadsheetParams struct{}

// ListSpreadsheets returns an MCP tool handler that lists spreadsheets accessible to the
// authenticated service account. Results are ordered by most recently modified, capped at 100.
// Each line contains the spreadsheet name and its ID, e.g. "Budget 2026 (ID: 1BxiMVs...)".
//
// Only spreadsheets explicitly shared with the service account are visible — this does not
// enumerate all spreadsheets in the organization.
//
// If no spreadsheets are accessible, returns "No spreadsheets found." as a successful result.
//
// Errors:
//   - Returns the Drive API error on network failure or insufficient permissions (403).
//   - If the Drive API is not enabled, returns a "googleapi: Error 403" with a
//     service-not-enabled message.
func ListSpreadsheets(driveSvc *drive.Service) func(context.Context, *mcp.CallToolRequest, ListSpreadsheetParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListSpreadsheetParams) (*mcp.CallToolResult, any, error) {
		fileList, err := driveSvc.Files.List().
			Q("mimeType='application/vnd.google-apps.spreadsheet'").
			Fields("files(id, name)").
			OrderBy("modifiedTime desc").
			PageSize(100).
			Context(ctx).
			Do()
		if err != nil {
			return nil, nil, err
		}

		if len(fileList.Files) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "No spreadsheets found."},
				},
			}, nil, nil
		}

		var lines []string
		for _, f := range fileList.Files {
			lines = append(lines, fmt.Sprintf("%s (ID: %s)", f.Name, f.Id))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: strings.Join(lines, "\n")},
			},
		}, nil, nil
	}
}

// FilteredReadParams defines the input for ReadSpreadsheetFiltered.
// The sheet's first row is always treated as headers and used for column name resolution.
// Limit defaults to 100 when omitted (zero value) to prevent unbounded responses.
type FilteredReadParams struct {
	SpreadsheetID string   `json:"spreadsheet_id"`
	Sheet         string   `json:"sheet"`
	Filters       []Filter `json:"filters"`
	Limit         int      `json:"limit,omitempty"`
}

// ReadSpreadsheetFiltered returns an MCP tool handler that reads rows from a single sheet,
// filtered by a column value. It fetches the entire sheet (A:ZZ) and applies filtering
// in-memory via filterRows.
//
// The first row of the sheet is treated as column headers. The filter_column must match
// a header exactly (case-sensitive). See filterRows for operator semantics.
//
// The response is tab-separated text with the header row first, followed by matched data rows.
// If the sheet is empty, returns "No data found." as a successful result. If the filter matches
// zero rows, returns only the header line.
//
// Errors:
//   - Returns the Sheets API error on network failure, invalid spreadsheet ID, non-existent
//     sheet name (400), or insufficient permissions (403).
//   - Returns a filterRows error if filter_column does not match any header or the operator
//     is not one of "eq", "contains", "gt", "lt", "like".
func ReadSpreadsheetFiltered(srv *sheetsapi.Service) func(context.Context, *mcp.CallToolRequest, FilteredReadParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args FilteredReadParams) (*mcp.CallToolResult, any, error) {
		rangeStr := args.Sheet + "!A:ZZ"
		resp, err := srv.Spreadsheets.Values.Get(args.SpreadsheetID, rangeStr).Context(ctx).Do()
		if err != nil {
			return nil, nil, err
		}

		if len(resp.Values) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "No data found."},
				},
			}, nil, nil
		}

		// First row is headers
		headerRow := resp.Values[0]
		headers := make([]string, len(headerRow))
		for i, h := range headerRow {
			headers[i] = cellToString(h)
		}

		dataRows := resp.Values[1:]

		// Default limit to 100 when omitted
		limit := args.Limit
		if limit == 0 {
			limit = 100
		}

		matched, err := filterRows(headers, dataRows, args.Filters, limit)
		if err != nil {
			return nil, nil, err
		}

		// Format as tab-separated text with header
		var result strings.Builder
		for i, h := range headers {
			if i > 0 {
				result.WriteString("\t")
			}
			result.WriteString(h)
		}
		result.WriteString("\n")

		for _, row := range matched {
			for i := range headers {
				if i > 0 {
					result.WriteString("\t")
				}
				if i < len(row) {
					result.WriteString(cellToString(row[i]))
				}
			}
			result.WriteString("\n")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: result.String()},
			},
		}, nil, nil
	}
}
