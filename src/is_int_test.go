package main

import "testing"

func TestIsInt_RecognizesOptionalSignsAndRejectsNonIntegers(t *testing.T) {
	for _, s := range []string{"0", "42", "+7", "-9"} {
		if !IsInt(s) {
			t.Fatalf("IsInt(%q) = false, want true", s)
		}
	}

	for _, s := range []string{"", " ", "+", "-", "3.14", "12a", "1e3"} {
		if IsInt(s) {
			t.Fatalf("IsInt(%q) = true, want false", s)
		}
	}
}
