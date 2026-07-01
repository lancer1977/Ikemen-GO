package main

import "testing"

func TestNeedsIniQuotes(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"", true},
		{" spaced", true},
		{"a,b", true},
		{"true", true},
		{"12", true},
		{"12.5", true},
		{"plain", false},
		{"under_score", false},
	} {
		if got := needsIniQuotes(tc.in); got != tc.want {
			t.Fatalf("needsIniQuotes(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
