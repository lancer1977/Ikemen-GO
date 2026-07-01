package main

import (
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaIniValueStringJoinsArrayTablesAndRejectsNestedTables(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	arr := L.NewTable()
	arr.Append(lua.LString("alpha"))
	arr.Append(lua.LNumber(12))
	got, err := luaIniValueString(arr)
	if err != nil {
		t.Fatalf("luaIniValueString(array) error = %v", err)
	}
	if got != "alpha, 12" {
		t.Fatalf("luaIniValueString(array) = %q, want alpha, 12", got)
	}

	obj := L.NewTable()
	obj.RawSetString("key", lua.LString("value"))
	if _, err := luaIniValueString(obj); err == nil || !strings.Contains(err.Error(), "nested non-array table") {
		t.Fatalf("luaIniValueString(object) error = %v, want nested non-array table", err)
	}
}
