package main

import "testing"

func TestMusicAppendAndOverride(t *testing.T) {
	m := Music{
		"stage": []*bgMusic{{bgmusic: "a"}, {bgmusic: "b"}},
	}
	m.Append(Music{
		"stage": []*bgMusic{{bgmusic: "c"}},
		"title": []*bgMusic{{bgmusic: "d"}},
	})

	if len(m["stage"]) != 3 || m["stage"][2].bgmusic != "c" {
		t.Fatalf("Append stage = %#v", m["stage"])
	}
	if len(m["title"]) != 1 || m["title"][0].bgmusic != "d" {
		t.Fatalf("Append title = %#v", m["title"])
	}

	m.Override(Music{
		"stage": []*bgMusic{{bgmusic: "x"}, {bgmusic: "y"}, {bgmusic: "z"}},
		"bonus": []*bgMusic{{bgmusic: "w"}},
	})

	if got := m["stage"]; len(got) != 3 || got[0].bgmusic != "x" || got[1].bgmusic != "y" || got[2].bgmusic != "z" {
		t.Fatalf("Override stage = %#v", got)
	}
	if got := m["bonus"]; len(got) != 1 || got[0].bgmusic != "w" {
		t.Fatalf("Override bonus = %#v", got)
	}
}
