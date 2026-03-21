package sheets

import (
	"strings"
	"testing"
)

func TestFilterRows(t *testing.T) {
	headers := []string{"Name", "Account", "Amount"}

	rows := [][]any{
		{"Alice", "Savings", "100"},
		{"Bob", "Checking", "250"},
		{"Carol", "savings", "50"},
		{"Dave", "The Lehn Fund", "1000"},
		{"Eve", "SAVINGS", "75"},
	}

	tests := []struct {
		name      string
		headers   []string
		rows      [][]any
		filters   []Filter
		limit     int
		wantCount int
		wantErr   bool
		errSubstr string
	}{
		// Given rows with mixed-case Account values
		// When filtering by Account eq "Savings"
		// Then all case variants match (3 rows)
		{
			name:      "column found exact case",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Account", Operator: "eq", Value: "Savings"}},
			limit:     0,
			wantCount: 3,
		},
		// Given headers that do not include "Status"
		// When filtering by a non-existent column
		// Then an error is returned
		{
			name:    "column not found",
			headers: headers,
			rows:    rows,
			filters: []Filter{{Column: "Status", Operator: "eq", Value: "active"}},
			limit:   0,
			wantErr: true,
		},
		// Given rows with distinct Name values
		// When filtering Name eq "Alice"
		// Then exactly one row matches
		{
			name:      "eq match",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Name", Operator: "eq", Value: "Alice"}},
			limit:     0,
			wantCount: 1,
		},
		// Given rows with distinct Name values
		// When filtering Name eq "Zara"
		// Then no rows match
		{
			name:      "eq no match",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Name", Operator: "eq", Value: "Zara"}},
			limit:     0,
			wantCount: 0,
		},
		// Given rows with Account values "Savings", "savings", and "SAVINGS"
		// When filtering Account eq "SAVINGS"
		// Then all three case variants match
		{
			name:      "eq case insensitive",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Account", Operator: "eq", Value: "SAVINGS"}},
			limit:     0,
			wantCount: 3,
		},
		// Given a row with Account "The Lehn Fund"
		// When filtering Account contains "fund"
		// Then that row matches (case-insensitive)
		{
			name:      "contains substring present",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Account", Operator: "contains", Value: "fund"}},
			limit:     0,
			wantCount: 1,
		},
		// Given rows with known Account values
		// When filtering Account contains "xyz"
		// Then no rows match
		{
			name:      "contains substring absent",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Account", Operator: "contains", Value: "xyz"}},
			limit:     0,
			wantCount: 0,
		},
		// Given a row with Account "The Lehn Fund"
		// When filtering Account contains "LEHN"
		// Then that row matches (case-insensitive)
		{
			name:      "contains case insensitive",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Account", Operator: "contains", Value: "LEHN"}},
			limit:     0,
			wantCount: 1,
		},
		// Given rows with Amount values 100, 250, 50, 1000, 75
		// When filtering Amount gt 100
		// Then two rows match (250, 1000)
		{
			name:      "gt numeric comparison",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Amount", Operator: "gt", Value: "100"}},
			limit:     0,
			wantCount: 2,
		},
		// Given rows with Amount values 100, 250, 50, 1000, 75
		// When filtering Amount lt 100
		// Then two rows match (50, 75)
		{
			name:      "lt numeric comparison",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Amount", Operator: "lt", Value: "100"}},
			limit:     0,
			wantCount: 2,
		},
		// Given rows where one Score value is "n/a"
		// When filtering Score gt 70
		// Then the non-numeric cell is skipped and two rows match
		{
			name:      "gt non-numeric cell skipped",
			headers:   []string{"Name", "Score"},
			rows:      [][]any{{"Alice", "80"}, {"Bob", "n/a"}, {"Carol", "90"}},
			filters:   []Filter{{Column: "Score", Operator: "gt", Value: "70"}},
			limit:     0,
			wantCount: 2,
		},
		// Given numeric Amount values
		// When filtering Amount gt "abc" (non-numeric filter value)
		// Then no rows match
		{
			name:      "gt non-numeric filter value skips all",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Amount", Operator: "gt", Value: "abc"}},
			limit:     0,
			wantCount: 0,
		},
		// Given numeric Amount values
		// When filtering Amount lt "abc" (non-numeric filter value)
		// Then no rows match
		{
			name:      "lt non-numeric filter value skips all",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Amount", Operator: "lt", Value: "abc"}},
			limit:     0,
			wantCount: 0,
		},
		// Given rows where three Account values contain "s"
		// When filtering with limit 0
		// Then all three rows are returned (unlimited)
		{
			name:      "limit zero unlimited",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Account", Operator: "contains", Value: "s"}},
			limit:     0,
			wantCount: 3,
		},
		// Given rows where three Account values contain "s"
		// When filtering with limit 2
		// Then only two rows are returned
		{
			name:      "limit caps results",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Account", Operator: "contains", Value: "s"}},
			limit:     2,
			wantCount: 2,
		},
		// Given one row matching Name eq "Alice"
		// When filtering with limit 100
		// Then only the one matching row is returned
		{
			name:      "fewer matches than limit",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Name", Operator: "eq", Value: "Alice"}},
			limit:     100,
			wantCount: 1,
		},
		// Given a ragged row missing the third column
		// When filtering column C eq ""
		// Then the ragged row matches with an implicit empty string
		{
			name:      "ragged row",
			headers:   []string{"A", "B", "C"},
			rows:      [][]any{{"x", "y"}, {"a", "b", "c"}},
			filters:   []Filter{{Column: "C", Operator: "eq", Value: ""}},
			limit:     0,
			wantCount: 1,
		},
		// Given an empty data set
		// When filtering by any criteria
		// Then no rows are returned
		{
			name:      "empty data set",
			headers:   headers,
			rows:      [][]any{},
			filters:   []Filter{{Column: "Name", Operator: "eq", Value: "Alice"}},
			limit:     0,
			wantCount: 0,
		},
		// Given rows with known Name values
		// When filtering Name eq "Nobody"
		// Then no rows match
		{
			name:      "all rows filtered out",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Name", Operator: "eq", Value: "Nobody"}},
			limit:     0,
			wantCount: 0,
		},
		// Given a valid set of rows
		// When filtering with an unsupported operator "regex"
		// Then an error is returned
		{
			name:    "invalid operator",
			headers: headers,
			rows:    rows,
			filters: []Filter{{Column: "Name", Operator: "regex", Value: ".*"}},
			limit:   0,
			wantErr: true,
		},
		// Given rows with Status and Region columns
		// When filtering Status eq "Active" AND Region eq "West"
		// Then only rows satisfying both conditions are returned
		{
			name:    "two eq filters AND match",
			headers: []string{"Name", "Status", "Region"},
			rows: [][]any{
				{"Alice", "Active", "West"},
				{"Bob", "Active", "East"},
				{"Carol", "Inactive", "West"},
				{"Dave", "Active", "West"},
			},
			filters:   []Filter{{Column: "Status", Operator: "eq", Value: "Active"}, {Column: "Region", Operator: "eq", Value: "West"}},
			limit:     0,
			wantCount: 2,
		},
		// Given a row with Status "Active" but Region "East"
		// When filtering Status eq "Active" AND Region eq "West"
		// Then the row is excluded because one filter does not match
		{
			name:    "one of two filters no match excludes row",
			headers: []string{"Name", "Status", "Region"},
			rows: [][]any{
				{"Bob", "Active", "East"},
			},
			filters:   []Filter{{Column: "Status", Operator: "eq", Value: "Active"}, {Column: "Region", Operator: "eq", Value: "West"}},
			limit:     0,
			wantCount: 0,
		},
		// Given rows with Status, Region, and Balance columns
		// When filtering with eq + contains + gt
		// Then only rows satisfying all three conditions are returned
		{
			name:    "three filters mixed operators",
			headers: []string{"Name", "Status", "Region", "Balance"},
			rows: [][]any{
				{"Alice", "Active", "Northwest", "5000"},
				{"Bob", "Active", "Southwest", "500"},
				{"Carol", "Inactive", "Northwest", "8000"},
				{"Dave", "Active", "East", "3000"},
			},
			filters: []Filter{
				{Column: "Status", Operator: "eq", Value: "Active"},
				{Column: "Region", Operator: "contains", Value: "west"},
				{Column: "Balance", Operator: "gt", Value: "1000"},
			},
			limit:     0,
			wantCount: 1,
		},
		// Given rows with data
		// When filtering with an empty filters array
		// Then all rows are returned up to limit
		{
			name:      "empty filters returns all rows",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{},
			limit:     0,
			wantCount: 5,
		},
		// Given rows with numeric Balance values
		// When filtering Balance gt 100 AND Balance lt 500
		// Then only rows strictly between 100 and 500 match
		{
			name:    "range query same column gt and lt",
			headers: []string{"Name", "Balance"},
			rows: [][]any{
				{"Alice", "50"},
				{"Bob", "250"},
				{"Carol", "100"},
				{"Dave", "500"},
				{"Eve", "300"},
			},
			filters:   []Filter{{Column: "Balance", Operator: "gt", Value: "100"}, {Column: "Balance", Operator: "lt", Value: "500"}},
			limit:     0,
			wantCount: 2,
		},
		// Given rows with Status values
		// When filtering Status eq "Active" AND Status eq "Inactive"
		// Then no rows can satisfy both and the result is empty
		{
			name:    "conflicting filters same column empty result",
			headers: []string{"Name", "Status"},
			rows: [][]any{
				{"Alice", "Active"},
				{"Bob", "Inactive"},
			},
			filters:   []Filter{{Column: "Status", Operator: "eq", Value: "Active"}, {Column: "Status", Operator: "eq", Value: "Inactive"}},
			limit:     0,
			wantCount: 0,
		},
		// Given a valid first filter and an invalid column in the second filter
		// When validating filters
		// Then an indexed error referencing filters[1] is returned
		{
			name:      "invalid column in second filter produces indexed error",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Name", Operator: "eq", Value: "Alice"}, {Column: "BadCol", Operator: "eq", Value: "X"}},
			limit:     0,
			wantErr:   true,
			errSubstr: "filters[1]",
		},
		// Given a valid first filter and an invalid operator in the second filter
		// When validating filters
		// Then an indexed error referencing filters[1] is returned
		{
			name:      "invalid operator in any filter produces indexed error",
			headers:   headers,
			rows:      rows,
			filters:   []Filter{{Column: "Name", Operator: "eq", Value: "Alice"}, {Column: "Account", Operator: "nope", Value: "X"}},
			limit:     0,
			wantErr:   true,
			errSubstr: "filters[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When filtering with the given filters
			got, err := filterRows(tt.headers, tt.rows, tt.filters, tt.limit)

			// Then the result matches expectations
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("got %d rows, want %d", len(got), tt.wantCount)
			}
		})
	}
}
