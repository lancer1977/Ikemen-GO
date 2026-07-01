package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaArgHelpers(t *testing.T) {
	t.Parallel()

	l := lua.NewState()
	defer l.Close()

	l.Push(lua.LNil)
	if !nilArg(l, 1) {
		t.Fatal("nilArg(nil) = false, want true")
	}
	l.Pop(1)

	l.Push(lua.LFalse)
	if nilArg(l, 1) {
		t.Fatal("nilArg(false) = true, want false")
	}
	if !boolArg(l, 1) {
		t.Fatal("boolArg(false) = true, want false")
	}
	l.Pop(1)

	l.Push(lua.LString("hello"))
	if got := strArg(l, 1); got != "hello" {
		t.Fatalf("strArg = %q, want %q", got, "hello")
	}
	l.Pop(1)

	l.Push(lua.LNumber(12.5))
	if got := numArg(l, 1); got != 12.5 {
		t.Fatalf("numArg = %v, want 12.5", got)
	}
	l.Pop(1)

	tbl := l.NewTable()
	tbl.RawSetString("present", lua.LTrue)
	if tableArg(l, 1) != nil {
		t.Fatal("tableArg should read the current top value only")
	}
	l.Push(tbl)
	if got := tableArg(l, 1); got != tbl {
		t.Fatalf("tableArg = %#v, want %#v", got, tbl)
	}
	if !tableHasKey(tbl, "present") {
		t.Fatal("tableHasKey(present) = false, want true")
	}
	if tableHasKey(tbl, "missing") {
		t.Fatal("tableHasKey(missing) = true, want false")
	}
	l.Pop(1)

	ud := newUserData(l, NewCommandList(nil))
	l.Push(ud)
	if got, ok := commandListArg(l, 1); !ok || got == nil {
		t.Fatalf("commandListArg = (%#v, %v), want non-nil command list", got, ok)
	}
	if got := toUserData(l, 1); got == nil {
		t.Fatal("toUserData(userdata) = nil, want value")
	}
	l.Pop(1)

	l.Push(lua.LString("not userdata"))
	if got, ok := commandListArg(l, 1); ok || got != nil {
		t.Fatalf("commandListArg(string) = (%#v, %v), want (nil, false)", got, ok)
	}
}
