package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestParseMusicSection(t *testing.T) {
	if got := parseMusicSection(nil); len(got) != 0 {
		t.Fatalf("parseMusicSection(nil) = %#v, want empty", got)
	}

	section := ini.Empty().Section("Music")
	section.Key("select.bgmusic").SetValue("theme.ogg")
	section.Key("select.bgm.loop").SetValue("12")
	section.Key("folder.track.bgm.volume").SetValue("80")

	music := parseMusicSection(section)
	if len(music) != 2 {
		t.Fatalf("parseMusicSection returned %d prefixes, want 2", len(music))
	}
	if !music.HasPrefix("select") {
		t.Fatalf("parseMusicSection should have select prefix: %#v", music)
	}
	if !music.HasPrefix("folder_track") {
		t.Fatalf("parseMusicSection should have folder_track prefix: %#v", music)
	}
}
