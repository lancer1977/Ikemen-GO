package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/ini.v1"
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

	// A nil table short-circuits to ini.Empty(), which is never nil and always
	// carries the implicit DEFAULT section -- but that section must hold no
	// keys, which is what "empty" has to mean here.
	f2, err := luaTableToIniFile(nil)
	if err != nil {
		t.Fatalf("luaTableToIniFile(nil) returned error: %v", err)
	}
	if f2 == nil {
		t.Fatal("luaTableToIniFile(nil) returned a nil file")
	}
	if names := f2.SectionStrings(); len(names) != 1 || names[0] != ini.DefaultSection {
		t.Fatalf("luaTableToIniFile(nil) sections = %v, want only %q", names, ini.DefaultSection)
	}
	if keys := f2.Section("").Keys(); len(keys) != 0 {
		t.Fatalf("luaTableToIniFile(nil) default section has %d keys, want 0", len(keys))
	}
}
