package main

import "testing"

func TestExtractPrefixFromMusicKey(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in   string
		want string
	}{
		{"select.bgmusic.volume", "select"},
		{"title.music.loop", "title"},
		{"stage.bgm", "stage"},
		{"fallback.key", "fallback"},
		{"plain", ""},
	} {
		if got := extractPrefixFromMusicKey(tc.in); got != tc.want {
			t.Fatalf("extractPrefixFromMusicKey(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
