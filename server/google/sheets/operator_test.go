package sheets

import (
	"testing"
)

// Given an unsupported operator name
// When creating the operator
// Then an error is returned
func TestNewOperator_InvalidName(t *testing.T) {
	_, err := newOperator("regex", "foo")

	if err == nil {
		t.Fatal("expected error for invalid operator name")
	}
}

func TestEqOp(t *testing.T) {
	// Given an eq operator matching "Alice"
	// When the cell value is exactly "Alice"
	// Then it matches
	t.Run("exact match", func(t *testing.T) {
		op, _ := newOperator("eq", "Alice")

		if !op.Match("Alice") {
			t.Error("expected match")
		}
	})

	// Given an eq operator matching "alice" (lowercase)
	// When the cell value is "ALICE" (uppercase)
	// Then it matches regardless of case
	t.Run("case insensitive", func(t *testing.T) {
		op, _ := newOperator("eq", "alice")
		got := op.Match("ALICE")

		if !got {
			t.Error("expected match")
		}
	})

	// Given an eq operator matching "Alice"
	// When the cell value is "Bob"
	// Then it does not match
	t.Run("no match", func(t *testing.T) {
		op, _ := newOperator("eq", "Alice")
		got := op.Match("Bob")

		if got {
			t.Error("expected no match")
		}
	})

	// Given an eq operator matching ""
	// When the cell value is also ""
	// Then it matches
	t.Run("empty strings", func(t *testing.T) {
		op, _ := newOperator("eq", "")
		got := op.Match("")

		if !got {
			t.Error("expected match")
		}
	})

	// Given an eq operator matching ""
	// When the cell value is "x"
	// Then it does not match
	t.Run("value empty cell not", func(t *testing.T) {
		op, _ := newOperator("eq", "")
		got := op.Match("x")
		if got {
			t.Error("expected no match")
		}
	})
}

func TestContainsOp(t *testing.T) {
	// Given a contains operator matching "fund"
	// When the cell value is "The Fund"
	// Then it matches the substring
	t.Run("substring present", func(t *testing.T) {
		op, _ := newOperator("contains", "fund")
		got := op.Match("The Fund")

		if !got {
			t.Error("expected match")
		}
	})

	// Given a contains operator matching "FUND" (uppercase)
	// When the cell value has mixed case "The Fund"
	// Then it matches regardless of case
	t.Run("case insensitive", func(t *testing.T) {
		op, _ := newOperator("contains", "FUND")
		got := op.Match("The Fund")

		if !got {
			t.Error("expected match")
		}
	})

	// Given a contains operator matching "xyz"
	// When the cell value is "hello world"
	// Then it does not match
	t.Run("no match", func(t *testing.T) {
		op, _ := newOperator("contains", "xyz")
		got := op.Match("hello world")

		if got {
			t.Error("expected no match")
		}
	})

	// Given a contains operator matching "" (empty string)
	// When the cell value is any non-empty string
	// Then it matches because every string contains ""
	t.Run("empty substring matches all", func(t *testing.T) {
		op, _ := newOperator("contains", "")
		got := op.Match("anything")

		if !got {
			t.Error("expected match")
		}
	})

	// Given a contains operator matching "hello"
	// When the cell value is exactly "hello"
	// Then it matches
	t.Run("full string match", func(t *testing.T) {
		op, _ := newOperator("contains", "hello")
		got := op.Match("hello")

		if !got {
			t.Error("expected match")
		}
	})
}

func TestGtOp(t *testing.T) {
	// Given a gt operator with threshold 100
	// When the cell value is 250
	// Then it matches because 250 > 100
	t.Run("cell greater", func(t *testing.T) {
		op, _ := newOperator("gt", "100")

		if !op.Match("250") {
			t.Error("expected match")
		}
	})

	// Given a gt operator with threshold 100
	// When the cell value is exactly 100
	// Then it does not match because gt is strictly greater
	t.Run("cell equal", func(t *testing.T) {
		op, _ := newOperator("gt", "100")

		if op.Match("100") {
			t.Error("expected no match")
		}
	})

	// Given a gt operator with threshold 100
	// When the cell value is 50
	// Then it does not match
	t.Run("cell less", func(t *testing.T) {
		op, _ := newOperator("gt", "100")

		if op.Match("50") {
			t.Error("expected no match")
		}
	})

	// Given a gt operator with a non-numeric threshold "abc"
	// When the cell value is numeric
	// Then it does not match because the threshold can't be parsed
	t.Run("non-numeric filter", func(t *testing.T) {
		op, _ := newOperator("gt", "abc")

		if op.Match("100") {
			t.Error("expected no match")
		}
	})

	// Given a gt operator with threshold 100
	// When the cell value is "n/a"
	// Then it does not match because the cell can't be parsed
	t.Run("non-numeric cell", func(t *testing.T) {
		op, _ := newOperator("gt", "100")

		if op.Match("n/a") {
			t.Error("expected no match")
		}
	})

	// Given a gt operator with threshold -10
	// When the cell value is 5
	// Then it matches because 5 > -10
	t.Run("negative numbers", func(t *testing.T) {
		op, _ := newOperator("gt", "-10")

		if !op.Match("5") {
			t.Error("expected match")
		}
	})

	// Given a gt operator with threshold 1.5
	// When the cell value is 1.6
	// Then it matches because 1.6 > 1.5
	t.Run("decimal values", func(t *testing.T) {
		op, _ := newOperator("gt", "1.5")

		if !op.Match("1.6") {
			t.Error("expected match")
		}
	})

	// Given a gt operator with a date threshold "3/20/2026"
	// When the cell value is "3/21/2026"
	// Then it matches because 3/21 is after 3/20
	t.Run("date after", func(t *testing.T) {
		op, _ := newOperator("gt", "3/20/2026")

		if !op.Match("3/21/2026") {
			t.Error("expected match")
		}
	})

	// Given a gt operator with a date threshold "3/20/2026"
	// When the cell value is "3/20/2026" (same day)
	// Then it does not match because gt is strictly after
	t.Run("date equal", func(t *testing.T) {
		op, _ := newOperator("gt", "3/20/2026")

		if op.Match("3/20/2026") {
			t.Error("expected no match")
		}
	})

	// Given a gt operator with a date threshold "3/20/2026"
	// When the cell value is "3/19/2026"
	// Then it does not match
	t.Run("date before", func(t *testing.T) {
		op, _ := newOperator("gt", "3/20/2026")

		if op.Match("3/19/2026") {
			t.Error("expected no match")
		}
	})

	// Given a gt operator with an ISO date threshold "2026-03-20"
	// When the cell value uses M/D/YYYY format "3/21/2026"
	// Then it matches across formats
	t.Run("date mixed formats", func(t *testing.T) {
		op, _ := newOperator("gt", "2026-03-20")

		if !op.Match("3/21/2026") {
			t.Error("expected match")
		}
	})
}

func TestLikeOp(t *testing.T) {
	// Given a like operator matching "Wholefoods"
	// When the cell value is exactly "Wholefoods"
	// Then it matches
	t.Run("exact match", func(t *testing.T) {
		op, _ := newOperator("like", "Wholefoods")

		if !op.Match("Wholefoods") {
			t.Error("expected match")
		}
	})

	// Given a like operator matching "Wholefoods"
	// When the cell value is the misspelling "Wholefds"
	// Then it matches because the characters appear in order
	t.Run("fuzzy misspelling", func(t *testing.T) {
		op, _ := newOperator("like", "Wholefoods")

		if !op.Match("Wholefds") {
			t.Error("expected match")
		}
	})

	// Given a like operator matching "Wholefoods"
	// When the cell value is the shorter variant "Wholefood"
	// Then it matches
	t.Run("missing trailing letter", func(t *testing.T) {
		op, _ := newOperator("like", "Wholefoods")

		if !op.Match("Wholefood") {
			t.Error("expected match")
		}
	})

	// Given a like operator matching "wholefoods" (lowercase)
	// When the cell value is "WHOLEFOODS" (uppercase)
	// Then it matches regardless of case
	t.Run("case insensitive", func(t *testing.T) {
		op, _ := newOperator("like", "wholefoods")

		if !op.Match("WHOLEFOODS") {
			t.Error("expected match")
		}
	})

	// Given a like operator matching "Wholefoods"
	// When the cell value is "Target"
	// Then it does not match
	t.Run("no match", func(t *testing.T) {
		op, _ := newOperator("like", "Wholefoods")

		if op.Match("Target") {
			t.Error("expected no match")
		}
	})
}

func TestLtOp(t *testing.T) {
	// Given a lt operator with threshold 100
	// When the cell value is 50
	// Then it matches because 50 < 100
	t.Run("cell less", func(t *testing.T) {
		op, _ := newOperator("lt", "100")

		if !op.Match("50") {
			t.Error("expected match")
		}
	})

	// Given a lt operator with threshold 100
	// When the cell value is exactly 100
	// Then it does not match because lt is strictly less
	t.Run("cell equal", func(t *testing.T) {
		op, _ := newOperator("lt", "100")

		if op.Match("100") {
			t.Error("expected no match")
		}
	})

	// Given a lt operator with threshold 100
	// When the cell value is 250
	// Then it does not match
	t.Run("cell greater", func(t *testing.T) {
		op, _ := newOperator("lt", "100")

		if op.Match("250") {
			t.Error("expected no match")
		}
	})

	// Given a lt operator with a non-numeric threshold "abc"
	// When the cell value is numeric
	// Then it does not match because the threshold can't be parsed
	t.Run("non-numeric filter", func(t *testing.T) {
		op, _ := newOperator("lt", "abc")

		if op.Match("100") {
			t.Error("expected no match")
		}
	})

	// Given a lt operator with threshold 100
	// When the cell value is "n/a"
	// Then it does not match because the cell can't be parsed
	t.Run("non-numeric cell", func(t *testing.T) {
		op, _ := newOperator("lt", "100")

		if op.Match("n/a") {
			t.Error("expected no match")
		}
	})

	// Given a lt operator with threshold 5
	// When the cell value is -10
	// Then it matches because -10 < 5
	t.Run("negative numbers", func(t *testing.T) {
		op, _ := newOperator("lt", "5")

		if !op.Match("-10") {
			t.Error("expected match")
		}
	})

	// Given a lt operator with threshold 1.6
	// When the cell value is 1.5
	// Then it matches because 1.5 < 1.6
	t.Run("decimal values", func(t *testing.T) {
		op, _ := newOperator("lt", "1.6")

		if !op.Match("1.5") {
			t.Error("expected match")
		}
	})

	// Given a lt operator with a date threshold "3/20/2026"
	// When the cell value is "3/19/2026"
	// Then it matches because 3/19 is before 3/20
	t.Run("date before", func(t *testing.T) {
		op, _ := newOperator("lt", "3/20/2026")

		if !op.Match("3/19/2026") {
			t.Error("expected match")
		}
	})

	// Given a lt operator with a date threshold "3/20/2026"
	// When the cell value is "3/20/2026" (same day)
	// Then it does not match because lt is strictly before
	t.Run("date equal", func(t *testing.T) {
		op, _ := newOperator("lt", "3/20/2026")

		if op.Match("3/20/2026") {
			t.Error("expected no match")
		}
	})

	// Given a lt operator with a date threshold "3/20/2026"
	// When the cell value is "3/21/2026"
	// Then it does not match
	t.Run("date after", func(t *testing.T) {
		op, _ := newOperator("lt", "3/20/2026")

		if op.Match("3/21/2026") {
			t.Error("expected no match")
		}
	})
}
