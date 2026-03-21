## Context

The `read_spreadsheet_filtered` MCP tool currently accepts a single filter (one column, one operator, one value). The
`filterRows` pure function handles all filtering logic and is well-tested with 19 cases. The proposal calls for
replacing the flat filter fields with an array of filter objects, all ANDed together.

Current params struct:

```go
type FilteredReadParams struct {
    SpreadsheetID string `json:"spreadsheet_id"`
    Sheet         string `json:"sheet"`
    FilterColumn  string `json:"filter_column"`
    Operator      string `json:"operator"`
    Value         string `json:"value"`
    Limit         int    `json:"limit,omitempty"`
}
```

Current `filterRows` signature:

```go
func filterRows(headers []string, rows [][]any, column string, op string, value string, limit int) ([][]any, error)
```

## Goals / Non-Goals

**Goals:**

- Support N filters ANDed together in a single `read_spreadsheet_filtered` call.
- Maintain the existing operator semantics (eq, contains, gt, lt) — no behavior changes per-operator.
- Upfront validation of all filters before row iteration.
- Clear, indexed error messages (e.g., `filters[1]: column "Foo" not found`).

**Non-Goals:**

- OR logic or nested boolean expressions.
- New operators (regex, not-eq, etc.).
- Multiple filters on the same column with implicit OR (each filter is independent; two filters on the same column means
  the row must satisfy both).

## Decisions

### Decision 1: Introduce a `Filter` struct and change params to use `[]Filter`

```go
type Filter struct {
    Column   string `json:"column"`
    Operator string `json:"operator"`
    Value    string `json:"value"`
}

type FilteredReadParams struct {
    SpreadsheetID string   `json:"spreadsheet_id"`
    Sheet         string   `json:"sheet"`
    Filters       []Filter `json:"filters"`
    Limit         int      `json:"limit,omitempty"`
}
```

**Rationale**: A typed slice of structs is self-describing in the JSON schema, validates naturally, and avoids the
fragility of parallel arrays. The go-sdk generates the MCP tool schema from struct tags, so this produces a clean schema
automatically.

**Alternative considered**: Parallel arrays (`filter_columns []string`, `operators []string`, `values []string`).
Rejected — arrays can desync and produce confusing errors.

### Decision 2: Validate all filters upfront before iterating rows

The updated `filterRows` will:

1. Resolve all column indices and validate all operators in a first pass.
2. Pre-parse numeric filter values for gt/lt operators.
3. Only then iterate rows.

**Rationale**: Fail-fast with a clear error like `filters[0]: column "Foo" not found` is better than doing partial work
and failing mid-iteration. The upfront pass also lets us build a pre-computed slice of resolved filter state (column
index, parsed numeric value) to avoid repeated work per row.

**Alternative considered**: Lazy validation during row iteration. Rejected — produces partial results before erroring,
which is confusing.

### Decision 3: Change `filterRows` signature to accept `[]Filter`

New signature:

```go
func filterRows(headers []string, rows [][]any, filters []Filter, limit int) ([][]any, error)
```

**Rationale**: Keeps `filterRows` as a pure function with a clean contract. The caller (handler) just passes the filters
slice through. All existing operator logic stays in the same function, just wrapped in an inner loop.

### Decision 4: Update existing tests mechanically, add new multi-filter cases

Existing single-filter test cases change from `filterRows(headers, rows, col, op, val, limit)` to
`filterRows(headers, rows, []Filter{{col, op, val}}, limit)`. This is a mechanical transformation. New test cases cover
multi-filter AND behavior, mixed operators, and edge cases like conflicting filters.

## Risks / Trade-offs

- **[Risk] Empty filters array**: A call with `filters: []` would match all rows (vacuous truth of AND over zero
  predicates). → Mitigation: This is actually reasonable behavior and consistent with "no filter = no restriction".
  Document it.
- **[Risk] Performance with many filters on large sheets**: Each row checks N predicates. → Mitigation: N is expected to
  be small (2-5 filters). Short-circuit on first non-matching filter via early break.
- **[Risk] Error message clarity with multiple filters**: Need to identify which filter failed. → Mitigation: Use
  indexed error messages like `filters[0]: column "X" not found`.
