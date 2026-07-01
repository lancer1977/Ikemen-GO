package main

import (
	"reflect"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaToJsonValueHandlesArraysObjectsAndCircularReferences(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	arr := L.NewTable()
	arr.Append(lua.LString("alpha"))
	arr.Append(lua.LNumber(12))
	got, err := luaToJsonValue(arr, nil)
	if err != nil {
		t.Fatalf("luaToJsonValue(array) error = %v", err)
	}
	if !reflect.DeepEqual(got, []any{"alpha", int64(12)}) {
		t.Fatalf("luaToJsonValue(array) = %#v, want [alpha 12]", got)
	}

	obj := L.NewTable()
	obj.RawSetString("name", lua.LString("ryu"))
	obj.RawSetString("score", lua.LNumber(7))
	got, err = luaToJsonValue(obj, nil)
	if err != nil {
		t.Fatalf("luaToJsonValue(object) error = %v", err)
	}
	m, ok := got.(map[string]any)
	if !ok || m["name"] != "ryu" || m["score"] != int64(7) {
		t.Fatalf("luaToJsonValue(object) = %#v", got)
	}

	circ := L.NewTable()
	circ.RawSetString("self", circ)
	if _, err := luaToJsonValue(circ, nil); err == nil || !strings.Contains(err.Error(), "circular reference") {
		t.Fatalf("luaToJsonValue(circular) error = %v, want circular reference", err)
	}
}
