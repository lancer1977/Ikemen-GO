package main

import "testing"

func TestIsNumeric_RecognizesCommonNumberForms(t *testing.T) {
	for _, s := range []string{"0", "  3.14 ", "-2.5e3", "+7", ".5"} {
		if !IsNumeric(s) {
			t.Fatalf("IsNumeric(%q) = false, want true", s)
		}
	}

	for _, s := range []string{"", " ", "abc", "1,000", "--2"} {
		if IsNumeric(s) {
			t.Fatalf("IsNumeric(%q) = true, want false", s)
		}
	}
}
