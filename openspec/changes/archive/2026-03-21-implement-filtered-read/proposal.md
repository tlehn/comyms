## Why

The `ReadSpreadsheetFiltered` tool is registered in the MCP server but returns a stub "not yet implemented" response.
Users need server-side filtering to avoid pulling entire sheets into the LLM context when they only need rows matching a
condition (e.g., "show me transactions for The Lehn Fund").

## What Changes

- Extract a pure `filterRows` function that performs in-memory filtering on spreadsheet data (headers + rows) without
  any Google API dependency.
- Implement the `ReadSpreadsheetFiltered` handler to fetch data via the Sheets API, delegate to `filterRows`, and format
  results as tab-separated text.
- Add comprehensive table-driven unit tests for `filterRows` covering all operators, edge cases, and error conditions.

## Capabilities

### New Capabilities

- `filtered-read`: Server-side row filtering by column value with support for eq, contains, gt, lt operators.

### Modified Capabilities

(none)

## Impact

- **Code**: `sheets.go` — new `filterRows` function, full implementation of `ReadSpreadsheetFiltered` handler. New
  `sheets_filter_test.go` test file.
- **APIs**: No new external APIs. The existing `read_spreadsheet_filtered` MCP tool goes from stub to functional.
- **Dependencies**: No new dependencies. Tests use only the standard library.
