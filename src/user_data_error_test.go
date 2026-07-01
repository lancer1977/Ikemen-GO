package main

import (
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestUserDataErrorRaisesExpectedMessage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("userDataError should panic via Lua error")
		}
		if !strings.Contains(r.(error).Error(), "Argument 1 is not a userdata of type") {
			t.Fatalf("userDataError panic = %v", r)
		}
	}()

	userDataError(L, 1, NewCommandList(nil))
}
