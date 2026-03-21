# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A CLI (`comyms`) that hosts multiple MCP (Model Context Protocol) servers as subcommands. Currently includes a Google
Sheets MCP server (`comyms google sheets`). Uses `github.com/spf13/cobra` for CLI dispatch,
`github.com/modelcontextprotocol/go-sdk` for MCP, and Google's Sheets/Drive API v4. Authentication is via Google
Application Default Credentials (service account JSON set via `GOOGLE_APPLICATION_CREDENTIALS`).

## Build & Development Commands

```sh
make build   # compile binary: comyms
make run     # go run .
make test    # go test ./...
make lint    # go vet ./... && golangci-lint run ./...
make clean   # remove binary
```

## Architecture

```
comyms/
├── main.go                        ← cmd.Execute()
├── cmd/
│   ├── root.go                    ← cobra root command ("comyms")
│   ├── google.go                  ← "comyms google" parent command (no Run)
│   └── google_sheets.go           ← "comyms google sheets" → sheets.Serve(ctx)
├── server/
│   └── google/
│       ├── auth.go                ← NewSheetsService, NewDriveService (shared ADC)
│       └── sheets/
│           ├── server.go          ← Serve(ctx) — creates MCP server, registers tools, runs stdio
│           ├── tools.go           ← tool handlers (ReadSpreadsheet, ListSheets, etc.)
│           ├── filter.go          ← filterRows, resolvedFilter, Filter type, cellToString
│           ├── operator.go        ← Operator interface + implementations
│           ├── filter_test.go
│           └── operator_test.go
├── go.mod
└── Makefile
```

- **main.go** — Thin entry point that calls `cmd.Execute()`.
- **cmd/** — Cobra command definitions. Each subcommand calls the corresponding server package's `Serve(ctx)` function.
  Server packages must not import cobra.
- **server/google/auth.go** — Shared Google ADC authentication. Exports `NewSheetsService(ctx)` and
  `NewDriveService(ctx)`.
- **server/google/sheets/** — Google Sheets MCP server. `server.go` is the entry point (`Serve`), `tools.go` has MCP
  tool handlers, `filter.go` has in-memory row filtering, `operator.go` has filter operator types.

The server communicates over stdio (`mcp.StdioTransport{}`), so it is launched as a subprocess by the MCP host (e.g.,
Claude Desktop) via `comyms google sheets`.

## Testing Conventions

When writing unit tests, use **Gherkin-style comments** (`// Given`, `// When`, `// Then`) to describe each test case.
Place the comments **above** `t.Run` (or above the test body in standalone tests) so the scenario reads as a preamble
before the code.

Example (standalone test):

```go
func TestReadSpreadsheet(t *testing.T) {
	// Given a spreadsheet with data in Sheet1
	srv := setupMockSheetsService(t)

	// When we read the spreadsheet
	result, err := readSpreadsheet(srv, "spreadsheet-id", "Sheet1!A1:B2")

	// Then we get the expected values without error
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
```

Example (subtests):

```go
func TestEqOp(t *testing.T) {
	// Given an eq operator matching "alice" (lowercase)
	// When the cell value is "ALICE" (uppercase)
	// Then it matches regardless of case
	t.Run("case insensitive", func(t *testing.T) {
		op, _ := newOperator("eq", "alice")

		if !op.Match("ALICE") {
			t.Error("expected match")
		}
	})
}
```
