## Why

The current `read_spreadsheet_filtered` tool accepts only a single column filter. Real-world queries almost always
involve multiple conditions — e.g., "active accounts in the West region with balance over $10k". Today, an LLM agent
must fetch broadly and reason over excess data, wasting context window and increasing latency.

## What Changes

- Replace the flat `filter_column` / `operator` / `value` params with a `filters` array of `{column, operator, value}`
  objects. All filters are ANDed together.
- Update `filterRows` to accept and evaluate multiple filter predicates per row.
- Update existing unit tests to use the new signature and add new test cases for multi-filter AND behavior.

## Capabilities

### New Capabilities

_None — this extends the existing filtered-read capability._

### Modified Capabilities

- `filtered-read`: The filter input changes from a single column/operator/value to an array of filters, all ANDed
  together. This is a requirements-level change to how filtering is expressed and evaluated.

## Impact

- **Code**: `FilteredReadParams` struct, `filterRows` function signature and loop, `ReadSpreadsheetFiltered` handler
  call site, all tests in `sheets_filter_test.go`.
- **API**: MCP tool schema for `read_spreadsheet_filtered` changes (breaking). No other tools affected.
- **Dependencies**: None — no new packages required.
