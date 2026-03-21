## 1. Core Type and Function Changes

- [x] 1.1 Add `Filter` struct to `sheets.go` with `Column`, `Operator`, `Value` fields and JSON tags
- [x] 1.2 Update `FilteredReadParams` to replace flat filter fields with `Filters []Filter`
- [x] 1.3 Update `filterRows` signature to accept `[]Filter` and `limit int`
- [x] 1.4 Implement upfront validation loop: resolve all column indices and validate all operators before row iteration,
      with indexed error messages (`filters[N]: ...`)
- [x] 1.5 Implement AND evaluation: for each row, check all filters and short-circuit on first non-match
- [x] 1.6 Handle empty filters array (return all rows up to limit)

## 2. Handler Wiring

- [x] 2.1 Update `ReadSpreadsheetFiltered` handler to pass `args.Filters` to `filterRows`

## 3. Unit Tests

- [x] 3.1 Update existing single-filter test cases in `sheets_filter_test.go` to use `[]Filter` wrapper
- [x] 3.2 Add test: two eq filters on different columns (AND match)
- [x] 3.3 Add test: one of two filters does not match (row excluded)
- [x] 3.4 Add test: three filters with mixed operators (eq + contains + gt)
- [x] 3.5 Add test: empty filters array returns all rows up to limit
- [x] 3.6 Add test: range query on same column (gt + lt)
- [x] 3.7 Add test: conflicting filters on same column (empty result)
- [x] 3.8 Add test: invalid column in second filter produces indexed error
- [x] 3.9 Add test: invalid operator in any filter produces indexed error

## 4. Verification

- [x] 4.1 Run `make build` and verify it succeeds
- [x] 4.2 Run `make test` and verify all tests pass
- [x] 4.3 Run `make lint` and verify no lint issues
- [x] 4.4 Run `make format` and verify code is properly formatted
