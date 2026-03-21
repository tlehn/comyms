## Why

The project currently builds a single-purpose binary (`google-sheets-mcp`) that serves one MCP server. As more MCP
servers are needed (other Google APIs, non-Google services), each would require its own repo, binary, and duplicated
boilerplate. Restructuring into a CLI with subcommands (`comyms google sheets`) lets all MCP servers live in one
project, share common code (auth, utilities), and ship as a single binary.

## What Changes

- Add cobra as a CLI framework with nested subcommand dispatch (`comyms google sheets`)
- Move Google Sheets MCP server logic into `server/google/sheets/` package, exposing a `Serve(ctx) error` entry point
- Extract shared Google auth setup into `server/google/` package
- Replace `main.go` with cobra root command wiring
- Binary name changes from `google-sheets-mcp` to `comyms`; MCP host configs must update to `comyms google sheets`
- Rename Go module from `github.com/todlehn/google-sheets-mcp` to match new project identity

## Capabilities

### New Capabilities

- `cli-dispatch`: Cobra-based CLI with nested subcommands that route to individual MCP server `Serve()` functions
- `shared-google-auth`: Reusable Google ADC authentication extracted into `server/google/` for use by all Google API
  servers

### Modified Capabilities

- `filtered-read`: No requirement changes — implementation moves to `server/google/sheets/` but behavior is identical

## Impact

- **Code**: All existing `.go` files move into new package structure; `main.go` becomes a thin cobra entrypoint
- **Dependencies**: Adds `github.com/spf13/cobra`
- **Build**: Makefile updates for new binary name and package paths
- **MCP host configs**: Must change from `google-sheets-mcp` to `comyms google sheets`
- **Module path**: `go.mod` module name changes
