package main

import (
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestExecFuncHandlesMissingFunctionBooleanAndUnexpectedReturns(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if got, err := ExecFunc(L, "missing"); err == nil || got {
		t.Fatalf("ExecFunc(missing) = (%v, %v), want false error", got, err)
	}

	L.SetGlobal("truthy", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LTrue)
		return 1
	}))
	if got, err := ExecFunc(L, "truthy"); err != nil || !got {
		t.Fatalf("ExecFunc(truthy) = (%v, %v), want true nil", got, err)
	}

	L.SetGlobal("badret", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LString("nope"))
		return 1
	}))
	if got, err := ExecFunc(L, "badret"); err == nil || got || !strings.Contains(err.Error(), "unexpected return type") {
		t.Fatalf("ExecFunc(badret) = (%v, %v), want false unexpected return type", got, err)
	}
}
