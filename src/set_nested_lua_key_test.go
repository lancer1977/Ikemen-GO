package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestSetNestedLuaKeyPromotesScalarsAndPreservesMetadata(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tbl := L.NewTable()
	setNestedLuaKey(L, tbl, "stats.score", lua.LNumber(7), false)
	if tbl.RawGetString("stats").(*lua.LTable).RawGetString("score").String() != "7" {
		t.Fatalf("setNestedLuaKey nested value = %#v", tbl)
	}

	keep := L.NewTable()
	setNestedLuaKey(L, keep, "stats.score", lua.LNumber(7), true)
	setNestedLuaKey(L, keep, "stats", lua.LString("legacy"), true)
	stats := keep.RawGetString("stats").(*lua.LTable)
	if stats.RawGetString("__value").String() != "legacy" || stats.RawGetString("score").String() != "7" {
		t.Fatalf("setNestedLuaKey keepMeta = %#v", stats)
	}
}
