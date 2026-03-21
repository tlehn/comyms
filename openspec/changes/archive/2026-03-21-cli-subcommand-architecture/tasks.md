## 1. Project Setup

- [x] 1.1 Add `github.com/spf13/cobra` dependency
- [x] 1.2 Rename module in `go.mod` from `github.com/todlehn/google-sheets-mcp` to `github.com/todlehn/comyms`

## 2. Package Structure

- [x] 2.1 Create `server/google/auth.go` — export `NewSheetsService(ctx)` and `NewDriveService(ctx)` from existing
      `newSheetsService`/`newDriveService`
- [x] 2.2 Create `server/google/sheets/server.go` — implement `Serve(ctx context.Context) error` that creates MCP
      server, registers tools, and runs stdio transport
- [x] 2.3 Move tool handlers and param structs to `server/google/sheets/tools.go`
- [x] 2.4 Move `filterRows`, `resolvedFilter`, `Filter`, and `cellToString` to `server/google/sheets/filter.go`
- [x] 2.5 Move `Operator` interface and implementations to `server/google/sheets/operator.go`
- [x] 2.6 Move existing tests to `server/google/sheets/` with updated package names

## 3. Cobra CLI

- [x] 3.1 Create `cmd/root.go` — root cobra command for `comyms`
- [x] 3.2 Create `cmd/google.go` — `google` parent command (no Run, just groups subcommands)
- [x] 3.3 Create `cmd/google_sheets.go` — `sheets` subcommand under `google` that calls `sheets.Serve(ctx)`
- [x] 3.4 Replace `main.go` with cobra entrypoint (`cmd.Execute()`)

## 4. Build & Config

- [x] 4.1 Update `Makefile` — binary name to `comyms`, update package paths
- [x] 4.2 Update `CLAUDE.md` to reflect new project structure and commands
- [x] 4.3 Update `README.md` to reflect new binary name, subcommand usage, and MCP host config

## 5. Verification

- [x] 5.1 Run `make build` and verify it succeeds
- [x] 5.2 Run `make test` and verify all tests pass
- [x] 5.3 Run `make lint` and verify no lint issues
- [x] 5.4 Remove old top-level `main.go`, `sheets.go`, `operator.go` and verify clean build
