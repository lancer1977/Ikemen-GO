package main

import "testing"

func TestJoystickStateAccessorsReturnEmptyWhenIndexOutOfRange(t *testing.T) {
	var in Input

	for _, joy := range []int{-1, MaxPlayerNo, MaxPlayerNo + 5} {
		if got := in.GetJoystickAxes(joy); got != ([6]float32{}) {
			t.Fatalf("GetJoystickAxes(%d) = %#v, want zeros", joy, got)
		}
		if got := in.GetJoystickButtons(joy); len(got) != 0 {
			t.Fatalf("GetJoystickButtons(%d) len = %d, want 0", joy, len(got))
		}
	}
}
