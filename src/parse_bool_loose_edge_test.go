package main

import "testing"

func TestParseBoolLooseAcceptsCanonicalBooleanAliases(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"1", true},
		{"On", true},
		{"0", false},
		{"No", false},
	} {
		got, ok := parseBoolLoose(tc.in)
		if !ok || got != tc.want {
			t.Fatalf("parseBoolLoose(%q) = (%v, %v), want (%v, true)", tc.in, got, ok, tc.want)
		}
	}
}
