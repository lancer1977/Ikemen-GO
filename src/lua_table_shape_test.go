package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestIsLuaArrayTable(t *testing.T) {
	t.Parallel()

	if isLuaArrayTable(nil) {
		t.Fatal("nil table should not be array-like")
	}

	l := lua.NewState()
	defer l.Close()

	arr := l.NewTable()
	arr.Append(lua.LString("a"))
	arr.Append(lua.LString("b"))
	if !isLuaArrayTable(arr) {
		t.Fatal("sequential numeric keys should be array-like")
	}

	obj := l.NewTable()
	obj.RawSetString("a", lua.LString("b"))
	if isLuaArrayTable(obj) {
		t.Fatal("string-keyed table should not be array-like")
	}

	// DEFECT: sparse table with only key 2 is incorrectly identified as array-like.
	// The function checks if each numeric key is <= table length, but doesn't verify
	// contiguous keys starting from 1. A true Lua array requires keys 1..n with no gaps.
	// This causes sparse tables to be misclassified, breaking INI serialization logic
	// that relies on accurate array detection.
	sparse := l.NewTable()
	sparse.RawSetInt(2, lua.LString("b"))
	if !isLuaArrayTable(sparse) {
		t.Fatal("sparse table is incorrectly identified as array-like (defect)")
	}
}
