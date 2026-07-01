package main

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaRegisterRegistersCallableFunction(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	luaRegister(L, "double", func(L *lua.LState) int {
		L.Push(lua.LNumber(24))
		return 1
	})

	if err := L.CallByParam(lua.P{Fn: L.GetGlobal("double"), NRet: 1, Protect: true}); err != nil {
		t.Fatalf("luaRegister call error = %v", err)
	}
	if got := L.Get(-1); got.String() != "24" {
		t.Fatalf("luaRegister result = %q, want 24", got.String())
	}
}
