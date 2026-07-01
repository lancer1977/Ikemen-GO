package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaKeyToStringFormatsStringsNumbersAndFallbacks(t *testing.T) {
	if got := luaKeyToString(lua.LString("alpha")); got != "alpha" {
		t.Fatalf("luaKeyToString(string) = %q, want alpha", got)
	}
	if got := luaKeyToString(lua.LNumber(12)); got != "12" {
		t.Fatalf("luaKeyToString(int-like number) = %q, want 12", got)
	}
	if got := luaKeyToString(lua.LNumber(12.5)); got != "12.5" {
		t.Fatalf("luaKeyToString(float number) = %q, want 12.5", got)
	}
	if got := luaKeyToString(lua.LBool(true)); got != "true" {
		t.Fatalf("luaKeyToString(fallback) = %q, want true", got)
	}
}
