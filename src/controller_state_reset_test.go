package main

import (
	"testing"

	"github.com/veandco/go-sdl2/sdl"
)

func TestResetControllerState(t *testing.T) {
	t.Parallel()

	oldInput := input
	defer func() { input = oldInput }()

	input = Input{}
	input.controllerstate[1] = &ControllerState{
		Axes:      [6]int8{1, 2, 3, 4, 5, 6},
		Buttons:   map[sdl.GameControllerButton]byte{sdl.CONTROLLER_BUTTON_A: 1, sdl.CONTROLLER_BUTTON_B: 2},
		HasRumble: true,
	}

	resetControllerState(1)
	cs := input.controllerstate[1]
	if cs == nil {
		t.Fatal("expected controller state to remain allocated")
	}
	if cs.Axes != [6]int8{} {
		t.Fatalf("expected axes reset, got %#v", cs.Axes)
	}
	if len(cs.Buttons) != 0 {
		t.Fatalf("expected buttons cleared, got %#v", cs.Buttons)
	}
	if cs.HasRumble {
		t.Fatal("expected rumble flag to clear")
	}

	resetControllerState(-1)
	resetControllerState(len(input.controllerstate))
}
