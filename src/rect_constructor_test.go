package main

import "testing"

func TestNewRect(t *testing.T) {
	t.Parallel()

	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	r := NewRect()
	if r == nil {
		t.Fatal("NewRect returned nil")
	}
	if r.window != sys.scrrect {
		t.Fatalf("window = %#v, want %#v", r.window, sys.scrrect)
	}
	if r.alpha != [2]int32{255, 0} || r.localScale != 1 {
		t.Fatalf("unexpected defaults: %#v", r)
	}
}
