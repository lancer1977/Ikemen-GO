package main

import "testing"

func TestItoa_FormatsIntegersAsDecimalStringsAgain(t *testing.T) {
	for _, tc := range []struct {
		in   int
		want string
	}{
		{-42, "-42"},
		{0, "0"},
		{12345, "12345"},
	} {
		if got := Itoa(tc.in); got != tc.want {
			t.Fatalf("Itoa(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
