## ADDED Requirements

### Requirement: Filter rows by column value

The system SHALL accept a column name, operator, and value, and return only rows where the specified column satisfies
the filter condition.

#### Scenario: Exact match with eq operator

- **WHEN** operator is "eq" and value is "savings"
- **THEN** return rows where the filter column's cell value equals "savings" (case-insensitive)

#### Scenario: Case-insensitive eq matching

- **WHEN** operator is "eq", value is "Savings", and cell contains "SAVINGS"
- **THEN** the row SHALL be included in results

#### Scenario: No rows match eq filter

- **WHEN** operator is "eq" and no cells in the filter column match the value
- **THEN** return an empty result set (no error)

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
