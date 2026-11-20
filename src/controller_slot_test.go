package main

import (
	"testing"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

func TestFindFreeControllerSlot(t *testing.T) {
	t.Parallel()

	oldInput := input
	defer func() { input = oldInput }()

	input = Input{}
	if got := findFreeControllerSlot(); got != 0 {
		t.Fatalf("findFreeControllerSlot() = %d, want 0 on empty input", got)
	}

	// sdl.GameController is an incomplete cgo type and cannot be allocated from
	// Go. findFreeControllerSlot only compares slots against nil, so an opaque
	// non-nil pointer is enough and is never dereferenced.
	input.controllers[0] = (*sdl.GameController)(unsafe.Pointer(new(byte)))
	if got := findFreeControllerSlot(); got != 1 {
		t.Fatalf("findFreeControllerSlot() = %d, want 1 after filling slot 0", got)
	}
}
