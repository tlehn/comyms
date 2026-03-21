## ADDED Requirements

### Requirement: Filter rows by multiple columns with AND logic

The system SHALL accept an array of filters, each specifying a column, operator, and value. A row SHALL be included in
results only if it satisfies ALL filters (logical AND).

#### Scenario: Two filters both match

- **WHEN** filters are
  `[{column: "Status", operator: "eq", value: "Active"}, {column: "Region", operator: "eq", value: "West"}]`
- **THEN** return only rows where Status equals "Active" AND Region equals "West"

#### Scenario: One of two filters does not match

- **WHEN** filters are
  `[{column: "Status", operator: "eq", value: "Active"}, {column: "Region", operator: "eq", value: "West"}]` and a row
  has Status "Active" but Region "East"
- **THEN** that row SHALL NOT be included in results

#### Scenario: Three filters with mixed operators

- **WHEN** filters are
  `[{column: "Status", operator: "eq", value: "Active"}, {column: "Region", operator: "contains", value: "west"}, {column: "Balance", operator: "gt", value: "1000"}]`
- **THEN** return only rows satisfying all three conditions simultaneously

#### Scenario: Empty filters array

- **WHEN** filters array is empty (`[]`)
- **THEN** return all rows (up to limit), since the vacuous AND is true

#### Scenario: Single filter in array

- **WHEN** filters array contains exactly one filter
- **THEN** behavior SHALL be identical to the previous single-filter behavior

### Requirement: Upfront validation of all filters

The system SHALL validate all filters (column resolution and operator validation) before iterating any data rows.

#### Scenario: Invalid column in second filter

- **WHEN** filters are
  `[{column: "Status", operator: "eq", value: "Active"}, {column: "BadCol", operator: "eq", value: "X"}]`
- **THEN** return an error referencing the filter index: `filters[1]: column "BadCol" not found; available columns: ...`
- **AND** no partial results SHALL be returned

#### Scenario: Invalid operator in any filter

- **WHEN** any filter has an operator other than "eq", "contains", "gt", "lt"
- **THEN** return an error referencing the filter index:
  `filters[N]: invalid operator "X"; valid operators: eq, contains, gt, lt`

### Requirement: Multiple filters on the same column

The system SHALL allow multiple filters targeting the same column. Each filter is evaluated independently and all must
match (AND).

#### Scenario: Range query on same column

- **WHEN** filters are
  `[{column: "Balance", operator: "gt", value: "100"}, {column: "Balance", operator: "lt", value: "500"}]`
- **THEN** return only rows where Balance is strictly between 100 and 500

#### Scenario: Conflicting filters on same column

- **WHEN** filters are
  `[{column: "Status", operator: "eq", value: "Active"}, {column: "Status", operator: "eq", value: "Inactive"}]`
- **THEN** return an empty result set (no row can satisfy both, no error)

## MODIFIED Requirements

### Requirement: Filter rows by column value

The system SHALL accept an array of filter objects, each specifying a column name, operator, and value. A row SHALL be
included in results only if it satisfies ALL filter conditions. Each individual filter's operator semantics (eq,
contains, gt, lt) remain unchanged.

#### Scenario: Exact match with eq operator

- **WHEN** filters contain `{column: "Type", operator: "eq", value: "savings"}`
- **THEN** return rows where the Type column's cell value equals "savings" (case-insensitive)

#### Scenario: Case-insensitive eq matching

- **WHEN** filters contain `{column: "Type", operator: "eq", value: "Savings"}` and cell contains "SAVINGS"
- **THEN** the row SHALL be included in results (assuming all other filters also match)

#### Scenario: No rows match eq filter

- **WHEN** filters contain `{column: "Type", operator: "eq", value: "nonexistent"}` and no cells match
- **THEN** return an empty result set (no error)
