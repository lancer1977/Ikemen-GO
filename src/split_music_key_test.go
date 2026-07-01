package main

import "testing"

func TestSplitMusicKey(t *testing.T) {
	t.Parallel()

	prefix, prop := splitMusicKey("bgm.loop")
	if prefix != "" || prop != "bgm.loop" {
		t.Fatalf("bare bgm key = %q %q", prefix, prop)
	}

	prefix, prop = splitMusicKey("select.bgmusic.volume")
	if prefix != "select" || prop != "bgmusic.volume" {
		t.Fatalf("anchored key = %q %q", prefix, prop)
	}

	prefix, prop = splitMusicKey("folder.track.bgm.loopstart")
	if prefix != "folder.track" || prop != "bgm.loopstart" {
		t.Fatalf("dotted prefix key = %q %q", prefix, prop)
	}

	prefix, prop = splitMusicKey("fallback.key")
	if prefix != "fallback" || prop != "key" {
		t.Fatalf("fallback key = %q %q", prefix, prop)
	}
}
