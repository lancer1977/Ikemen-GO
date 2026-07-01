package main

import (
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestJSONToLuaValueConvertsScalarsArraysAndObjects(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	got, err := jsonToLuaValue(L, strings.NewReader(`null`))
	if err != nil || got != lua.LNil {
		t.Fatalf("jsonToLuaValue(null) = (%#v, %v), want nil nil", got, err)
	}

	got, err = jsonToLuaValue(L, strings.NewReader(`[1, "two", false]`))
	if err != nil {
		t.Fatalf("jsonToLuaValue(array) error = %v", err)
	}
	tbl, ok := got.(*lua.LTable)
	if !ok || tbl.Len() != 3 || tbl.RawGetInt(1).String() != "1" || tbl.RawGetInt(2).String() != "two" || tbl.RawGetInt(3) != lua.LFalse {
		t.Fatalf("jsonToLuaValue(array) = %#v", got)
	}

	got, err = jsonToLuaValue(L, strings.NewReader(`{"name":"ryu","score":7}`))
	if err != nil {
		t.Fatalf("jsonToLuaValue(object) error = %v", err)
	}
	tbl, ok = got.(*lua.LTable)
	if !ok || tbl.RawGetString("name").String() != "ryu" || tbl.RawGetString("score").String() != "7" {
		t.Fatalf("jsonToLuaValue(object) = %#v", got)
	}

	if _, err := jsonToLuaValue(L, strings.NewReader(`{`)); err == nil {
		t.Fatal("jsonToLuaValue should reject invalid JSON")
	}
}
