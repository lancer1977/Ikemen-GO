package main

import "testing"

func TestCheckAxisForDpad(t *testing.T) {
	t.Parallel()

	prevSys := sys
	prevLUT := ButtonToStringLUT
	defer func() {
		sys = prevSys
		ButtonToStringLUT = prevLUT
	}()

	sys.cfg.Input.ControllerStickSensitivity = 0.5
	ButtonToStringLUT = map[int]string{
		0: "UP",
		1: "LEFT",
		2: "RIGHT",
		3: "DOWN",
		6: "RS_UP",
		7: "RS_LEFT",
		8: "RS_RIGHT",
		9: "RS_DOWN",
	}

	axes := [6]float32{0.75, 0, 0, 0}
	if got := CheckAxisForDpad(&axes, 0); got != "RIGHT" {
		t.Fatalf("right axis = %q", got)
	}

	axes = [6]float32{0, -0.75, 0, 0}
	if got := CheckAxisForDpad(&axes, 0); got != "UP" {
		t.Fatalf("up axis = %q", got)
	}

	axes = [6]float32{0, 0, 0.75, 0}
	if got := CheckAxisForDpad(&axes, 0); got != "RS_RIGHT" {
		t.Fatalf("right stick axis = %q", got)
	}
}

func TestCheckAxisForTrigger(t *testing.T) {
	t.Parallel()

	prevLUT := ButtonToStringLUT
	defer func() { ButtonToStringLUT = prevLUT }()

	ButtonToStringLUT = map[int]string{
		15: "LT",
		16: "RT",
	}

	axes := [6]float32{0, 0, 0, 0, 1, 0}
	if got := CheckAxisForTrigger(&axes); got != "LT" {
		t.Fatalf("left trigger = %q", got)
	}

	axes = [6]float32{0, 0, 0, 0, 0, 2}
	if got := CheckAxisForTrigger(&axes); got != "RT" {
		t.Fatalf("right trigger = %q", got)
	}
}
