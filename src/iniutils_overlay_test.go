package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestIniutilsIniFirstValue_PrefersFirstShadowAndHandlesNil(t *testing.T) {
	if val, dup := iniFirstValue(nil); val != "" || dup != 0 {
		t.Fatalf("iniFirstValue(nil) = %q %d", val, dup)
	}

	f := ini.Empty(ini.LoadOptions{AllowShadows: true})
	sec, _ := f.NewSection("user")
	k, _ := sec.NewKey("name", "alpha")
	_, _ = sec.NewKey("name", "beta")

	if val, dup := iniFirstValue(k); val != "alpha" || dup != 1 {
		t.Fatalf("iniFirstValue() = %q %d, want alpha 1", val, dup)
	}
}

func TestIniutilsOverlayUserFirstWins_CopiesFirstShadowValues(t *testing.T) {
	dst := ini.Empty()
	_, _ = dst.NewSection("user")

	user := ini.Empty(ini.LoadOptions{AllowShadows: true})
	sec, _ := user.NewSection("user")
	_, _ = sec.NewKey("name", "alpha")
	_, _ = sec.NewKey("name", "beta")
	_, _ = sec.NewKey("rank", "3")

	overlayUserFirstWins(dst, user)

	gotName := dst.Section("user").Key("name").String()
	gotRank := dst.Section("user").Key("rank").String()
	if gotName != "alpha" || gotRank != "3" {
		t.Fatalf("overlayUserFirstWins() = name %q rank %q", gotName, gotRank)
	}
}
