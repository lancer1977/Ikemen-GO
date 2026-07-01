package main

import "testing"

func TestMusicPrefixHelpers_NormalizeAndDetectEntries(t *testing.T) {
	if got := normalizeMusicPrefix("round.1"); got != "round_1" {
		t.Fatalf("normalizeMusicPrefix() = %q, want round_1", got)
	}
	if got := musicKeyPrefix("victory.bgmusic"); got != "victory" {
		t.Fatalf("musicKeyPrefix(bgmusic) = %q, want victory", got)
	}
	if got := musicKeyPrefix("stage.music"); got != "stage" {
		t.Fatalf("musicKeyPrefix(music) = %q, want stage", got)
	}
	if got := musicKeyPrefix("title"); got != "title" {
		t.Fatalf("musicKeyPrefix(no suffix) = %q, want title", got)
	}

	m := Music{
		"round1": []*bgMusic{{}},
	}
	if !m.HasPrefix("round1") || !m.HasPrefix("round.1") {
		t.Fatal("HasPrefix should normalize dotted prefixes")
	}
	if m.HasPrefix("victory") {
		t.Fatal("HasPrefix should report false for missing prefix")
	}
}
