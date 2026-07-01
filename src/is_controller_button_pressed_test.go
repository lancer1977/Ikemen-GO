package main

import "testing"

func TestIsControllerButtonPressed(t *testing.T) {
	prevSys := sys
	defer func() { sys = prevSys }()

	StringToButtonLUT = map[string]int{
		"LT": 15,
		"A":  4,
	}

	cl := NewCommandList(nil)

	if cl.IsControllerButtonPressed("   ", 0) {
		t.Fatal("expected blank token to be false")
	}
	if cl.IsControllerButtonPressed("missing", 0) {
		t.Fatal("expected unknown token to be false")
	}
	if cl.IsControllerButtonPressed("A", 0) {
		t.Fatal("expected no joystick present to be false")
	}
	if cl.IsControllerButtonPressed("LT", -1) {
		t.Fatal("expected negative controller index to be false without input")
	}
}
