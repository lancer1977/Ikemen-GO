package main

import (
	"encoding/json"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestToLValueConvertsNumbersPointersAndFallbackScalars(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if got := toLValue(L, nil); got != lua.LNil {
		t.Fatalf("toLValue(nil) = %#v, want nil", got)
	}
	if got := toLValue(L, json.Number("12")); got != lua.LNumber(12) {
		t.Fatalf("toLValue(json int) = %#v, want 12", got)
	}
	if got := toLValue(L, json.Number("12.5")); got != lua.LNumber(12.5) {
		t.Fatalf("toLValue(json float) = %#v, want 12.5", got)
	}

	cmd := NewCommandList(nil)
	if got := toLValue(L, cmd); got.Type() != lua.LTUserData {
		t.Fatalf("toLValue(pointer userdata) = %s, want userdata", got.Type())
	}

	bg := &bgMusic{bgmusic: "intro.ogg", bgmloop: true, bgmvolume: 80, bgmloopstart: 12, bgmloopend: 34, bgmstartposition: 5}
	tbl, ok := toLValue(L, bg).(*lua.LTable)
	if !ok || tbl.RawGetString("bgm").String() != "intro.ogg" || tbl.RawGetString("volume").String() != "80" {
		t.Fatalf("toLValue(bgMusic) = %#v", tbl)
	}

	if got := toLValue(L, "plain"); got.String() != "plain" {
		t.Fatalf("toLValue(string) = %q, want plain", got.String())
	}
}
