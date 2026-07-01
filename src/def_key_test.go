package main

import "testing"

func TestNormalizeDefKey(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in   string
		want string
	}{
		{`chars\kfm\kfm.def`, "kfm"},
		{"kfm.def", "kfm"},
		{"kfm", "kfm"},
		{"KFM.DEF", "kfm"},
	} {
		if got := normalizeDefKey(tc.in); got != tc.want {
			t.Fatalf("normalizeDefKey(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
