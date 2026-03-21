## ADDED Requirements

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

### Requirement: Filter rows with contains operator

The system SHALL support substring matching via the "contains" operator.

#### Scenario: Substring match

- **WHEN** operator is "contains" and value is "fund"
- **THEN** return rows where the filter column's cell contains "fund" as a substring (case-insensitive)

#### Scenario: Contains with no match

- **WHEN** operator is "contains" and no cells contain the substring
- **THEN** return an empty result set

### Requirement: Filter rows with numeric comparison operators

The system SHALL support "gt" (greater than) and "lt" (less than) operators for numeric comparison.

#### Scenario: Greater than comparison

- **WHEN** operator is "gt" and value is "100"
- **THEN** return rows where the filter column's numeric value is strictly greater than 100

#### Scenario: Less than comparison

- **WHEN** operator is "lt" and value is "50"
- **THEN** return rows where the filter column's numeric value is strictly less than 50

#### Scenario: Non-numeric cell with numeric operator

- **WHEN** operator is "gt" or "lt" and a cell in the filter column is not a valid number
- **THEN** that row SHALL be silently skipped (not included, no error)

#### Scenario: Non-numeric filter value with numeric operator

- **WHEN** operator is "gt" or "lt" and the filter value is not a valid number
- **THEN** all rows SHALL be skipped (empty result, no error)

### Requirement: Column resolution uses strict case matching

The system SHALL match the filter column name against header row values using exact case comparison.

#### Scenario: Column found with exact case

- **WHEN** filter column is "Account" and headers contain "Account"
- **THEN** filtering proceeds on that column

#### Scenario: Column not found

- **WHEN** filter column does not match any header (exact case)
- **THEN** return an error listing available column names

### Requirement: Limit caps result rows

The system SHALL respect a limit parameter that caps the number of returned rows.

#### Scenario: Limit applied

- **WHEN** limit is 5 and 20 rows match the filter
- **THEN** return only the first 5 matching rows

#### Scenario: Limit of zero means unlimited

- **WHEN** limit is 0
- **THEN** return all matching rows

#### Scenario: Fewer matches than limit

- **WHEN** limit is 100 and only 3 rows match
- **THEN** return all 3 matching rows

### Requirement: Ragged rows treated as empty string

The system SHALL treat missing cells in short rows as empty strings for filtering purposes.

#### Scenario: Row shorter than header

- **WHEN** a row has fewer cells than the header row and the filter column index exceeds the row length
- **THEN** treat that cell as empty string for comparison

### Requirement: Invalid operator returns error

The system SHALL return an error for unrecognized operator values.

#### Scenario: Unknown operator

- **WHEN** operator is "regex" or any value other than "eq", "contains", "gt", "lt"
- **THEN** return an error indicating valid operators

### Requirement: Empty data set

The system SHALL handle sheets with no data rows gracefully.

#### Scenario: No data rows after header

- **WHEN** the sheet contains only a header row and no data rows
- **THEN** return an empty result set (no error)
