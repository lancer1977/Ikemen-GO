package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestParseIniLuaValueFallsBackWhenUnquoteFails(t *testing.T) {
	t.Parallel()

	l := lua.NewState()
	defer l.Close()

	got := parseIniLuaValue(l, `"bad\qescape"`)
	if got.String() != `bad\qescape` {
		t.Fatalf("parseIniLuaValue fallback = %q, want %q", got.String(), `bad\qescape`)
	}
}
