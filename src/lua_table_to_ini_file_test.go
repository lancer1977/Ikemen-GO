package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaTableToIniFileFlattensSectionsAndMapsDefaultSection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	root := L.NewTable()
	defSec := L.NewTable()
	defSec.RawSetString("name", lua.LString("ryu"))
	stats := L.NewTable()
	stats.RawSetString("score", lua.LNumber(7))
	defSec.RawSetString("stats", stats)
	root.RawSetString("default", defSec)

	f, err := luaTableToIniFile(root)
	if err != nil {
		t.Fatalf("luaTableToIniFile error = %v", err)
	}
	if f.Section("").Key("name").String() != "ryu" {
		t.Fatalf("default section name = %q, want ryu", f.Section("").Key("name").String())
	}
	if f.Section("").Key("stats.score").String() != "7" {
		t.Fatalf("default section stats.score = %q, want 7", f.Section("").Key("stats.score").String())
	}

	if f2, err := luaTableToIniFile(nil); err != nil || len(f2.Sections()) != 0 {
		t.Fatalf("luaTableToIniFile(nil) = (%#v, %v), want empty ini nil", f2, err)
	}
}
