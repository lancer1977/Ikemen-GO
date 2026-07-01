package main

import "testing"

func TestJoystickAccessorsReturnBlankWhenUnset(t *testing.T) {
	var in Input

	if got := in.GetMaxJoystickCount(); got != 0 {
		t.Fatalf("GetMaxJoystickCount = %d, want 0", got)
	}
	if got := in.GetJoystickName(0); got != "" {
		t.Fatalf("GetJoystickName = %q, want blank", got)
	}
	if got := in.GetJoystickPath(0); got != "" {
		t.Fatalf("GetJoystickPath = %q, want blank", got)
	}
	if got := in.GetJoystickGUID(0); got != "" {
		t.Fatalf("GetJoystickGUID = %q, want blank", got)
	}
}
