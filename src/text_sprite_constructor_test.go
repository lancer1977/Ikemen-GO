package main

import "testing"

func TestNewTextSpriteInitializesDefaults(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	ts := NewTextSprite()
	if ts == nil {
		t.Fatal("NewTextSprite returned nil")
	}
	if ts.id != -1 || ts.align != 1 || ts.xscl != 1 || ts.yscl != 1 {
		t.Fatalf("unexpected defaults: %#v", ts)
	}
	if ts.window != sys.scrrect {
		t.Fatalf("window = %#v, want %#v", ts.window, sys.scrrect)
	}
	if ts.palfx == nil {
		t.Fatal("expected palfx to be initialized")
	}
	if ts.params == nil {
		t.Fatal("expected params slice to be initialized")
	}
}
