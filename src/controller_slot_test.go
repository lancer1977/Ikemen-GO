package main

import (
	"testing"

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

	input.controllers[0] = &sdl.GameController{}
	if got := findFreeControllerSlot(); got != 1 {
		t.Fatalf("findFreeControllerSlot() = %d, want 1 after filling slot 0", got)
	}
}
