## Context

The `ReadSpreadsheetFiltered` handler is registered in `main.go` and has a complete spec comment in `sheets.go` (lines
159-195), but the body is a stub. The `FilteredReadParams` struct already exists. All other tool handlers
(`ReadSpreadsheet`, `ListSheets`, `ListSpreadsheets`) follow the same pattern: closure over `*sheets.Service`, fetch
data, format as tab-separated text.

The Google Sheets API v4 does not support server-side filtering — it only accepts A1-notation ranges. All filtering must
happen in-memory after fetching.

## Goals / Non-Goals

**Goals:**

- Implement `ReadSpreadsheetFiltered` with full operator support (eq, contains, gt, lt)
- Achieve high test coverage of filtering logic without mocking Google APIs
- Keep the architecture consistent with existing tool handlers

**Non-Goals:**

- Introducing interfaces or dependency injection for the Sheets service (not needed for testing the pure filtering
  logic)
- Testing Google API client behavior
- Supporting additional operators beyond the four specified
- Pagination of Sheets API responses (fetch all with `Sheet!A:ZZ`)

## Decisions

### Extract `filterRows` as a pure function

**Decision**: The filtering logic lives in a standalone `filterRows(headers, rows, column, op, value, limit)` function.
The handler calls it after fetching data from the API.

**Rationale**: This cleanly separates fetch (untestable without mocks) from filter (pure logic, trivially testable). No
interfaces or DI needed — just call the function directly in tests with fabricated data. This matches the user's
preference: test the filtering, not the Google libraries.

**Alternative considered**: Introduce a `SheetReader` interface and mock the fetch. Rejected — adds abstraction that
isn't needed when the interesting logic is all post-fetch.

### Column matching is strict case

**Decision**: Header lookup for `FilterColumn` uses exact case matching.

**Rationale**: Column names in spreadsheets are user-defined and case matters for disambiguation (e.g., "Status" vs
"STATUS"). Strict matching is predictable. The error message lists available columns, so discovery is easy.

### Ragged rows treated as empty string

**Decision**: If a row has fewer cells than the header row, missing cells are treated as `""` for filtering purposes.

**Rationale**: The Sheets API omits trailing empty cells from rows. Treating them as empty strings is consistent with
how spreadsheets work — an empty cell is an empty value, not an error.

### Invalid operator returns error

**Decision**: An unrecognized operator returns an error rather than matching nothing.

**Rationale**: Silent failure on a typo ("equ" instead of "eq") would be confusing. Explicit errors are better.

## Risks / Trade-offs

- **Large sheets**: Fetching all rows with `Sheet!A:ZZ` could be slow/large for very big sheets. → Acceptable for now;
  the Limit param caps output size even if fetch is large. Pagination is a future concern.
- **Numeric parsing edge cases**: `gt`/`lt` use `strconv.ParseFloat64`. Cells with currency symbols ("$100") or commas
  ("1,000") won't parse. → Skip silently per spec. Users can use `contains` for string matching on those values.
