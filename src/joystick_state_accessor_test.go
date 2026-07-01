package main

import "testing"

func TestJoystickStateAccessorsReturnEmptyWhenUnset(t *testing.T) {
	var in Input

	if got := in.GetJoystickAxes(0); got != ([6]float32{0, 0, 0, 0, 0, 0}) {
		t.Fatalf("GetJoystickAxes = %#v, want zeros", got)
	}
	if got := in.GetJoystickButtons(0); len(got) != 0 {
		t.Fatalf("GetJoystickButtons len = %d, want 0", len(got))
	}
}
