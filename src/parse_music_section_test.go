package main

import (
	"testing"

	"github.com/go-ini/ini"
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
	if len(music) != 1 {
		t.Fatalf("parseMusicSection returned %d prefixes, want 1", len(music))
	}
	if !music.HasPrefix("select") {
		t.Fatalf("parseMusicSection should have select prefix: %#v", music)
	}
}
