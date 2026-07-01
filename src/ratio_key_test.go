package main

import "testing"

func TestNormalizeRatioKey(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		w, h int64
		want string
	}{
		{16, 9, "16:9"},
		{1920, 1080, "16:9"},
		{-4, 6, "2:3"},
		{0, 9, "invalid"},
	} {
		if got := normalizeRatioKey(tc.w, tc.h); got != tc.want {
			t.Fatalf("normalizeRatioKey(%d,%d) = %q, want %q", tc.w, tc.h, got, tc.want)
		}
	}
}
