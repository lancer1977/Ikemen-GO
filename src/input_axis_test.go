package main

import "testing"

func TestCheckAxisForDpad(t *testing.T) {
	// Not parallel: this test mutates the package-level sys, which every test in
	// this package shares.

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
	// Not parallel: this test mutates the package-level sys, which every test in
	// this package shares.

	prevLUT := ButtonToStringLUT
	defer func() { ButtonToStringLUT = prevLUT }()

	// Not a defect: CheckAxisForTrigger indexes ButtonToStringLUT at 15+i for
	// SDL axis index i (4 for TRIGGERLEFT, 5 for TRIGGERRIGHT), landing on
	// keys 19 and 20. The real LUT built by initLUTs() (input_sdl.go:202-212)
	// defines exactly those keys as "LT"/"RT" -- the formula and the LUT
	// agree. The stub below mirrors those real key numbers, not an
	// arbitrarily different pair, so this test exercises the real contract.
	ButtonToStringLUT = map[int]string{
		19: "LT",
		20: "RT",
	}

	axes := [6]float32{0, 0, 0, 0, 1, 0}
	if got := CheckAxisForTrigger(&axes); got != "LT" {
		t.Fatalf("left trigger = %q, want %q", got, "LT")
	}

	axes = [6]float32{0, 0, 0, 0, 0, 2}
	if got := CheckAxisForTrigger(&axes); got != "RT" {
		t.Fatalf("right trigger = %q, want %q", got, "RT")
	}
}
