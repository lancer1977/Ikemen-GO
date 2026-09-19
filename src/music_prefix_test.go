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

	// Map keys are always stored in underscore form: music.go normalizes the
	// prefix when it builds the map, and HasPrefix normalizes the query before
	// looking it up. So a single "round_1" entry must be reachable by both
	// spellings -- that round trip is the whole point of the helper.
	m := Music{
		"round_1": []*bgMusic{{}},
	}
	if !m.HasPrefix("round.1") {
		t.Fatal("HasPrefix should find an underscore key from a dotted prefix")
	}
	if !m.HasPrefix("round_1") {
		t.Fatal("HasPrefix should find an underscore key from its own spelling")
	}
	if m.HasPrefix("round2") {
		t.Fatal("HasPrefix should not report a prefix that is absent")
	}
	// An entry present but empty must not count as a match.
	if (Music{"round_1": []*bgMusic{}}).HasPrefix("round.1") {
		t.Fatal("HasPrefix should reject a present but empty entry")
	}
	if m.HasPrefix("victory") {
		t.Fatal("HasPrefix should report false for missing prefix")
	}
}
