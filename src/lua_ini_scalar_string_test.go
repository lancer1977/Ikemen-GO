package main

import (
	"math"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaIniScalarStringFormatsCommonScalarsAndRejectsUnsupportedTypes(t *testing.T) {
	if got, err := luaIniScalarString(lua.LNil); err != nil || got != "" {
		t.Fatalf("luaIniScalarString(nil) = %q, %v, want empty nil", got, err)
	}
	if got, err := luaIniScalarString(lua.LTrue); err != nil || got != "true" {
		t.Fatalf("luaIniScalarString(true) = %q, %v, want true nil", got, err)
	}
	if got, err := luaIniScalarString(lua.LNumber(12.0)); err != nil || got != "12" {
		t.Fatalf("luaIniScalarString(int-like) = %q, %v, want 12 nil", got, err)
	}
	if got, err := luaIniScalarString(lua.LNumber(12.5)); err != nil || got != "12.5" {
		t.Fatalf("luaIniScalarString(float) = %q, %v, want 12.5 nil", got, err)
	}
	if got, err := luaIniScalarString(lua.LString("plain")); err != nil || got != "plain" {
		t.Fatalf("luaIniScalarString(string) = %q, %v, want plain nil", got, err)
	}
	if _, err := luaIniScalarString(lua.LNumber(math.NaN())); err == nil {
		t.Fatal("luaIniScalarString should reject NaN")
	}
	if _, err := luaIniScalarString(lua.LFunction(nil)); err == nil {
		t.Fatal("luaIniScalarString should reject unsupported types")
	}
}
