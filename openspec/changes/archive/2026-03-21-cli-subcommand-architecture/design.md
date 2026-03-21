## Context

The project is a single-binary MCP server for Google Sheets. All code lives in `package main` with three files:
`main.go` (entrypoint + tool registration), `sheets.go` (API clients, tool handlers, filtering), and `operator.go`
(filter operators). Auth uses Google ADC via `google.FindDefaultCredentials`.

The goal is to restructure into a multi-server CLI (`comyms`) where `comyms google sheets` starts the Sheets MCP server,
and future subcommands add more servers — all sharing common code within one binary.

## Goals / Non-Goals

**Goals:**

- Cobra-based CLI with nested subcommands mirroring the package tree
- Server packages expose `Serve(ctx context.Context) error` — no cobra dependency in server code
- Shared Google ADC auth in `server/google/` reusable across Google API servers
- Preserve all existing MCP tool behavior unchanged

**Non-Goals:**

- Adding new MCP servers in this change (only restructuring google-sheets)
- Global flags beyond what cobra provides by default
- Plugin/dynamic server loading — all servers are compiled in

## Decisions

### Package layout: `server/google/sheets/` hierarchy

```
comyms/
├── main.go                    ← cmd.Execute()
├── cmd/
│   ├── root.go                ← cobra root command ("comyms")
│   ├── google.go              ← "comyms google" parent command (no Run)
│   └── google_sheets.go       ← "comyms google sheets" → sheets.Serve(ctx)
├── server/
│   └── google/
│       ├── auth.go            ← NewSheetsService, NewDriveService (shared ADC)
│       └── sheets/
│           ├── server.go      ← Serve(ctx) — creates MCP server, registers tools, runs stdio
│           ├── tools.go       ← tool handlers (ReadSpreadsheet, ListSheets, etc.)
│           ├── filter.go      ← filterRows, resolvedFilter, Filter type
│           ├── operator.go    ← Operator interface + implementations
│           └── *_test.go
├── go.mod
└── Makefile
```

**Why this over flat `cmd/google-sheets/`**: The hierarchy gives a natural home for shared Google auth
(`server/google/auth.go`) and mirrors the cobra command tree. Adding `server/google/calendar/` later is obvious.

**Alternative considered**: Putting auth in a separate `shared/` or `pkg/` package. Rejected because the shared code is
Google-specific — `server/google/` is more semantically precise and avoids a grab-bag package.

### Server entry point: `Serve(ctx context.Context) error`

Each server package exposes a single function. The cobra command's `RunE` just calls it. This keeps server packages
testable without cobra and makes the contract between `cmd/` and `server/` explicit.

**Alternative considered**: Server packages returning `*cobra.Command`. Rejected because it couples server logic to the
CLI framework and makes integration testing harder.

### Auth extraction: `server/google/auth.go`

Extract `newSheetsService` and `newDriveService` into exported functions in `server/google/`. These return typed API
clients (`*sheets.Service`, `*drive.Service`) from ADC. Each server's `Serve()` calls the auth functions it needs.

This is intentionally thin — just the credential lookup + service creation. No abstraction layer over different Google
APIs.

### File split within `server/google/sheets/`

Current `sheets.go` mixes auth, tool handlers, filtering, and types. Split into:

- `server.go` — `Serve()` function: MCP server creation, tool registration, stdio transport
- `tools.go` — tool handler functions and param structs
- `filter.go` — `filterRows`, `resolvedFilter`, `Filter` type, `cellToString`
- `operator.go` — `Operator` interface and implementations (moved as-is)

**Why**: Each file has a single responsibility. `server.go` is the entry point. `tools.go` is where you add new tools.
`filter.go` is the in-memory query engine. Clean boundaries for testing.

### Module rename

`go.mod` module changes from `github.com/todlehn/google-sheets-mcp` to `github.com/todlehn/comyms`. All internal imports
update accordingly.

## Risks / Trade-offs

- **[Breaking MCP host configs]** → Users must update configs from `google-sheets-mcp` to `comyms google sheets`.
  Mitigated by documenting in README.
- **[Cobra dependency weight]** → Adds cobra + its transitive deps. Acceptable for a CLI binary; the dependency is
  stable and well-maintained.
- **[Package visibility]** → Moving from `package main` to exported packages means types/functions become part of the
  public API. This is intentional — other packages in the binary need to import them. Not publishing as a library, so
  semver concerns don't apply.
