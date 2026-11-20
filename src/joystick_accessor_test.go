package main

import "testing"

// Input.controllers and Input.controllerstate are fixed [MaxPlayerNo] arrays,
// not slices, so every accessor is in range for 0..MaxPlayerNo-1 regardless of
// whether a controller is actually connected. Only an out-of-range index takes
// the blank-value guard; an unset in-range slot is filtered by
// IsJoystickPresent instead.
func TestJoystickAccessorsReturnBlankWhenIndexOutOfRange(t *testing.T) {
	var in Input

	if got := in.GetMaxJoystickCount(); got != MaxPlayerNo {
		t.Fatalf("GetMaxJoystickCount = %d, want %d", got, MaxPlayerNo)
	}

	for _, joy := range []int{-1, MaxPlayerNo, MaxPlayerNo + 5} {
		if got := in.GetJoystickName(joy); got != "" {
			t.Fatalf("GetJoystickName(%d) = %q, want blank", joy, got)
		}
		if got := in.GetJoystickPath(joy); got != "" {
			t.Fatalf("GetJoystickPath(%d) = %q, want blank", joy, got)
		}
		if got := in.GetJoystickGUID(joy); got != "" {
			t.Fatalf("GetJoystickGUID(%d) = %q, want blank", joy, got)
		}
	}
}

func TestIsJoystickPresentRejectsUnsetAndOutOfRangeSlots(t *testing.T) {
	var in Input

	// Unset in-range slot: the nil check is what keeps callers off a nil
	// controller, since the index itself is valid.
	if in.IsJoystickPresent(0) {
		t.Fatalf("IsJoystickPresent(0) = true on a zero Input, want false")
	}
	for _, joy := range []int{-1, MaxPlayerNo} {
		if in.IsJoystickPresent(joy) {
			t.Fatalf("IsJoystickPresent(%d) = true, want false", joy)
		}
	}
}
