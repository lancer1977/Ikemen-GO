package main

import "testing"

func TestNormalizeSectionName(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in   string
		want string
	}{
		{"  Select  Screen ", "select_screen"},
		{"TITLE\tSCREEN", "title_screen"},
		{"", ""},
	} {
		if got := normalizeSectionName(tc.in); got != tc.want {
			t.Fatalf("normalizeSectionName(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
