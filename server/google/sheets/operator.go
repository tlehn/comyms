package sheets

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

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
		num, ok := parseFloat(value)
		return gtOp{num: num, numOk: ok}, nil
	case "lt":
		num, ok := parseFloat(value)
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

// gtOp matches when both values are numeric and cell > threshold.
// Non-numeric cells or a non-numeric filter value never match.
type gtOp struct {
	num   float64
	numOk bool
}

func (o gtOp) Match(cell string) bool {
	if !o.numOk {
		return false
	}
	cellNum, ok := parseFloat(cell)
	return ok && cellNum > o.num
}

// ltOp matches when both values are numeric and cell < threshold.
// Non-numeric cells or a non-numeric filter value never match.
type ltOp struct {
	num   float64
	numOk bool
}

func (o ltOp) Match(cell string) bool {
	if !o.numOk {
		return false
	}
	cellNum, ok := parseFloat(cell)
	return ok && cellNum < o.num
}

// likeOp matches when the cell value fuzzy-matches the filter value (case-insensitive).
// Useful for matching spelling variations like "Wholefood", "Wholefds", "Wholefoods".
type likeOp struct{ value string }

func (o likeOp) Match(cell string) bool {
	return fuzzy.MatchFold(o.value, cell) || fuzzy.MatchFold(cell, o.value)
}
