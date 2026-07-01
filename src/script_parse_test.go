package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestSplitIniListOutsideQuotes(t *testing.T) {
	t.Parallel()

	got := splitIniListOutsideQuotes(`alpha, "bravo, charlie", 'delta, echo', foxtrot`)
	want := []string{`alpha`, `"bravo, charlie"`, `'delta, echo'`, `foxtrot`}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseIniLuaValue(t *testing.T) {
	t.Parallel()

	l := lua.NewState()
	defer l.Close()

	if got := parseIniLuaValue(l, ` "hello" `); got.String() != "hello" {
		t.Fatalf("quoted string = %s", got.String())
	}
	if got := parseIniLuaValue(l, "true"); got != lua.LTrue {
		t.Fatalf("true = %s", got.String())
	}
	if got := parseIniLuaValue(l, "12"); got.String() != "12" {
		t.Fatalf("int = %s", got.String())
	}

	list := parseIniLuaValue(l, `1, "two, too", false`)
	tbl, ok := list.(*lua.LTable)
	if !ok {
		t.Fatalf("expected table, got %T", list)
	}
	if tbl.Len() != 3 {
		t.Fatalf("table len = %d, want 3", tbl.Len())
	}
	if tbl.RawGetInt(1).String() != "1" || tbl.RawGetInt(2).String() != "two, too" || tbl.RawGetInt(3) != lua.LFalse {
		t.Fatalf("unexpected table values: %s, %s, %s", tbl.RawGetInt(1).String(), tbl.RawGetInt(2).String(), tbl.RawGetInt(3).String())
	}
}
