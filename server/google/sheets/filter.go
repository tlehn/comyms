package sheets

import (
	"fmt"
	"strings"
)

// Filter defines a single column filter predicate.
type Filter struct {
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// cellToString coerces a cell value from the Sheets API into a string.
// The Sheets API returns cells as interface{} — strings pass through directly,
// all other types (float64 for numbers, bool, nil) are formatted via fmt.Sprintf.
func cellToString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

// resolvedFilter holds a validated filter ready for row evaluation.
type resolvedFilter struct {
	colIdx int
	op     Operator
}

// filterRows applies in-memory filters to spreadsheet data rows. All filters are ANDed
// together — a row must satisfy every filter to be included.
//
// Column resolution uses an exact, case-sensitive match against the headers slice.
// All filters are validated upfront before any row iteration. Errors reference the
// filter index (e.g., "filters[0]: column \"X\" not found").
//
// See the Operator interface and its implementations (eqOp, containsOp, gtOp, ltOp, likeOp)
// for operator semantics.
//
// An empty filters slice returns all rows up to limit (vacuous AND).
// A limit of 0 means no limit — all matching rows are returned. A positive limit stops
// collection after that many matches (first-match-wins order, preserving row order).
func filterRows(headers []string, rows [][]any, filters []Filter, limit int) ([][]any, error) {
	resolved := make([]resolvedFilter, len(filters))
	for i, f := range filters {
		colIdx := -1
		for j, h := range headers {
			if h == f.Column {
				colIdx = j
				break
			}
		}
		if colIdx == -1 {
			return nil, fmt.Errorf("filters[%d]: column %q not found; available columns: %s", i, f.Column, strings.Join(headers, ", "))
		}

		op, err := newOperator(f.Operator, f.Value)
		if err != nil {
			return nil, fmt.Errorf("filters[%d]: %w", i, err)
		}

		resolved[i] = resolvedFilter{colIdx: colIdx, op: op}
	}

	var matched [][]any
	for _, row := range rows {
		allMatch := true
		for _, rf := range resolved {
			var cell string
			if rf.colIdx < len(row) {
				cell = cellToString(row[rf.colIdx])
			}
			if !rf.op.Match(cell) {
				allMatch = false
				break
			}
		}

		if allMatch {
			matched = append(matched, row)
			if limit > 0 && len(matched) >= limit {
				break
			}
		}
	}

	return matched, nil
}
