package main

import "testing"

func TestScreenHeightAndWidthReflectGameSize(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys.gameHeight = 240
	sys.gameWidth = 320
	if got := sys.screenHeight(); got != 240 {
		t.Fatalf("screenHeight() = %v, want 240", got)
	}
	if got := sys.screenWidth(); got != 320 {
		t.Fatalf("screenWidth() = %v, want 320", got)
	}
}
