## ADDED Requirements

### Requirement: Root command exists

The binary SHALL provide a `comyms` root command that displays usage and available subcommands when invoked without
arguments.

#### Scenario: No arguments

- **WHEN** the user runs `comyms` with no arguments
- **THEN** the CLI prints usage information listing available subcommands

### Requirement: Google parent command groups Google API servers

The CLI SHALL provide a `google` parent command under root that groups all Google API MCP servers. It has no `Run`
function of its own.

#### Scenario: Google with no subcommand

- **WHEN** the user runs `comyms google` with no subcommand
- **THEN** the CLI prints usage information listing available Google subcommands (e.g., `sheets`)

### Requirement: Sheets subcommand starts the Google Sheets MCP server

The CLI SHALL provide a `sheets` subcommand under `google` that starts the Google Sheets MCP server over stdio by
calling `sheets.Serve(ctx)`.

#### Scenario: Start sheets server

- **WHEN** the user runs `comyms google sheets`
- **THEN** the Google Sheets MCP server starts and communicates over stdio

#### Scenario: Server error propagates

- **WHEN** `sheets.Serve(ctx)` returns an error
- **THEN** the CLI exits with a non-zero exit code and prints the error to stderr

### Requirement: Subcommand structure mirrors package tree

Each server package SHALL expose a `Serve(ctx context.Context) error` function. Cobra commands in `cmd/` SHALL call
these functions — server packages MUST NOT import cobra.

#### Scenario: Server package independence

- **WHEN** a server package (e.g., `server/google/sheets`) is compiled
- **THEN** it has no dependency on `github.com/spf13/cobra`
