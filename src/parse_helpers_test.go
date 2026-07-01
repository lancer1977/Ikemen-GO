package main

import (
	"math"
	"testing"
)

func TestIsFiniteIsNumericAndIsInt_HandleCommonInputs(t *testing.T) {
	if !IsFinite(123.5) {
		t.Fatalf("expected finite value to be accepted")
	}
	if IsFinite(float32(math.Inf(1))) {
		t.Fatalf("expected +Inf to be rejected")
	}
	if IsFinite(float32(math.NaN())) {
		t.Fatalf("expected NaN to be rejected")
	}

	if !IsNumeric("  12.5 ") {
		t.Fatalf("expected numeric string with whitespace to be accepted")
	}
	if IsNumeric("12.5x") {
		t.Fatalf("expected invalid numeric string to be rejected")
	}

	if !IsInt("  -42 ") {
		t.Fatalf("expected signed integer string with whitespace to be accepted")
	}
	if IsInt("+") || IsInt("12.5") || IsInt("abc") {
		t.Fatalf("expected invalid integer strings to be rejected")
	}
}
