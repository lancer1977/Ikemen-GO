package main

import (
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestFlattenLuaIniSectionFlattensNestedTablesAndRejectsNonStringKeys(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	root := L.NewTable()
	root.RawSetString("name", lua.LString("ryu"))
	nested := L.NewTable()
	nested.RawSetString("score", lua.LNumber(7))
	root.RawSetString("stats", nested)

	out := make(map[string]string)
	if err := flattenLuaIniSection(root, "char", out); err != nil {
		t.Fatalf("flattenLuaIniSection(root) error = %v", err)
	}
	if out["char.name"] != "ryu" || out["char.stats.score"] != "7" {
		t.Fatalf("flattenLuaIniSection(root) = %#v", out)
	}

	bad := L.NewTable()
	bad.RawSet(lua.LNumber(1), lua.LString("value"))
	if err := flattenLuaIniSection(bad, "", make(map[string]string)); err == nil || !strings.Contains(err.Error(), "INI keys must be strings") {
		t.Fatalf("flattenLuaIniSection(non-string key) error = %v, want INI keys must be strings", err)
	}
}
