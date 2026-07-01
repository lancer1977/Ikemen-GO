package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaIniAppendOrderAppendsUniqueKeysOnly(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tbl := L.NewTable()
	luaIniAppendOrder(L, tbl, "name")
	luaIniAppendOrder(L, tbl, "score")
	luaIniAppendOrder(L, tbl, "name")

	order := tbl.RawGetString("__order").(*lua.LTable)
	if order.Len() != 2 || order.RawGetInt(1).String() != "name" || order.RawGetInt(2).String() != "score" {
		t.Fatalf("luaIniAppendOrder order = %#v", order)
	}
}
