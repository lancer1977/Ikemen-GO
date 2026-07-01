package main

import "testing"

func TestReadDigit_RejectsEmptyLeadingZeroAndNonDigits(t *testing.T) {
	for _, s := range []string{"", "0", "012", "1a", "-1"} {
		if got, ok := readDigit(s); ok || got != 0 {
			t.Fatalf("readDigit(%q) = (%v, %v), want (0, false)", s, got, ok)
		}
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
