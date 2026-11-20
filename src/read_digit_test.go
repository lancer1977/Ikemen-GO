package main

import "testing"

func TestReadDigit_RejectsEmptyLeadingZeroAndNonDigits(t *testing.T) {
	// readDigit rejects: empty strings, leading zeros in multi-digit numbers, non-digits.
	// It accepts single "0" because the check is: len(d) >= 2 && d[0] == '0'.
	for _, s := range []string{"", "012", "1a", "-1"} {
		if got, ok := readDigit(s); ok || got != 0 {
			t.Fatalf("readDigit(%q) = (%v, %v), want (0, false)", s, got, ok)
		}
	}
	// Single "0" is accepted
	if got, ok := readDigit("0"); !ok || got != 0 {
		t.Fatalf("readDigit(\"0\") = (%v, %v), want (0, true)", got, ok)
	}
}

func TestReadDigit_ParsesSingleAndMultiDigitNumbers(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int32
	}{
		{"1", 1},
		{"42", 42},
		{"999", 999},
	} {
		got, ok := readDigit(tc.in)
		if !ok || got != tc.want {
			t.Fatalf("readDigit(%q) = (%v, %v), want (%v, true)", tc.in, got, ok, tc.want)
		}
	}
}
