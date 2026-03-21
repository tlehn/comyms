## 1. Core filtering function

- [x] 1.1 Implement
      `filterRows(headers []string, rows [][]interface{}, column string, op string, value string, limit int) ([][]interface{}, error)`
      in `sheets.go`
- [x] 1.2 Handle column resolution with strict case match; return error listing available columns on miss
- [x] 1.3 Implement "eq" operator (case-insensitive string match)
- [x] 1.4 Implement "contains" operator (case-insensitive substring match)
- [x] 1.5 Implement "gt" and "lt" operators (float64 comparison, skip non-numeric cells)
- [x] 1.6 Return error for invalid operator
- [x] 1.7 Handle ragged rows (short rows → empty string for missing cells)
- [x] 1.8 Apply limit (0 = unlimited, >0 caps result count)

## 2. Wire up the handler

- [x] 2.1 Implement `ReadSpreadsheetFiltered` handler: fetch via `srv.Spreadsheets.Values.Get(id, "Sheet!A:ZZ")`, call
      `filterRows`, format header + matched rows as tab-separated text
- [x] 2.2 Default limit to 100 when omitted (Limit == 0 in handler before calling filterRows)

## 3. Unit tests

- [x] 3.1 Create `sheets_filter_test.go` with table-driven tests for `filterRows`
- [x] 3.2 Test column resolution: exact match found, column not found → error with available columns
- [x] 3.3 Test eq operator: match, no match, case-insensitive
- [x] 3.4 Test contains operator: substring present, absent, case-insensitive
- [x] 3.5 Test gt/lt operators: numeric comparison works, non-numeric cell skipped, non-numeric filter value skips all
- [x] 3.6 Test limit: 0 unlimited, N caps at N, fewer matches than limit
- [x] 3.7 Test edge cases: ragged row, empty data set (no rows), all rows filtered out
- [x] 3.8 Test invalid operator → error

## 4. Cleanup

- [x] 4.1 Simplify the `ReadSpreadsheetFiltered` comment block in `sheets.go` — replace the implementation-specific spec
      (step-by-step algorithm, param docs, registration example) with a concise godoc-style comment describing what the
      tool does

## 5. Verification

- [x] 5.1 Run `make build` and verify it succeeds
- [x] 5.2 Run `make test` and verify all tests pass
- [x] 5.3 Run `make lint` and verify no lint issues
