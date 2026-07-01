package main

import "testing"

func TestNewAnim(t *testing.T) {
	t.Parallel()

	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	a := NewAnim(nil, "")
	if a == nil {
		t.Fatal("NewAnim returned nil for empty action")
	}
	if a.window != sys.scrrect || a.xscl != 1 || a.yscl != 1 || a.localScale != 1 || a.facing != 1 || a.lastUpdateFrame != -1 {
		t.Fatalf("unexpected defaults: %#v", a)
	}
	if a.palfx == nil {
		t.Fatal("expected palfx to be initialized")
	}
}
