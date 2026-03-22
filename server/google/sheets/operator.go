package sheets

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// dateFormats lists the date layouts we attempt when parsing cell/filter values.
var dateFormats = []string{
	"2006-01-02",
	"01/02/2006",
	"1/2/2006",
	"1/02/2006",
	"01/2/2006",
	"January 2, 2006",
	"Jan 2, 2006",
	"02-Jan-2006",
	"2-Jan-2006",
}

// Operator evaluates whether a cell value satisfies a filter condition.
type Operator interface {
	Match(cellValue string) bool
}

// newOperator creates an Operator for the given name and comparison value.
func newOperator(name, value string) (Operator, error) {
	switch name {
	case "eq":
		return eqOp{value: value}, nil
	case "contains":
		return containsOp{value: strings.ToLower(value)}, nil
	case "gt":
		num, ok := parseNumeric(value)
		return gtOp{num: num, numOk: ok}, nil
	case "lt":
		num, ok := parseNumeric(value)
		return ltOp{num: num, numOk: ok}, nil
	case "like":
		return likeOp{value: value}, nil
	default:
		return nil, fmt.Errorf("invalid operator %q; valid operators: eq, contains, gt, lt, like", name)
	}
}

func parseFloat(s string) (float64, bool) {
	n, err := strconv.ParseFloat(s, 64)
	return n, err == nil
}

// parseNumeric tries to interpret s as a number. If that fails, it tries
// common date formats and returns the Unix timestamp as a float64.
func parseNumeric(s string) (float64, bool) {
	if n, ok := parseFloat(s); ok {
		return n, true
	}
	s = strings.TrimSpace(s)
	for _, layout := range dateFormats {
		if t, err := time.Parse(layout, s); err == nil {
			return float64(t.Unix()), true
		}
	}
	return 0, false
}

// eqOp matches when the cell equals the value (case-insensitive).
type eqOp struct{ value string }

func (o eqOp) Match(cell string) bool {
	return strings.EqualFold(cell, o.value)
}

// containsOp matches when the cell contains the substring (case-insensitive).
type containsOp struct{ value string }

func (o containsOp) Match(cell string) bool {
	return strings.Contains(strings.ToLower(cell), o.value)
}

// gtOp matches when both values parse as numeric (including dates) and cell > threshold.
type gtOp struct {
	numOk bool
	num   float64
}

func (o gtOp) Match(cell string) bool {
	if !o.numOk {
		return false
	}
	cellNum, ok := parseNumeric(cell)
	return ok && cellNum > o.num
}

// ltOp matches when both values parse as numeric (including dates) and cell < threshold.
type ltOp struct {
	numOk bool
	num   float64
}

func (o ltOp) Match(cell string) bool {
	if !o.numOk {
		return false
	}
	cellNum, ok := parseNumeric(cell)
	return ok && cellNum < o.num
}

// likeOp matches when the cell value fuzzy-matches the filter value (case-insensitive).
// Useful for matching spelling variations like "Wholefood", "Wholefds", "Wholefoods".
type likeOp struct{ value string }

func (o likeOp) Match(cell string) bool {
	return fuzzy.MatchFold(o.value, cell) || fuzzy.MatchFold(cell, o.value)
}
