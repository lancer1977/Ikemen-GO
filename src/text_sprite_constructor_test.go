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
	// loadDefaults reslices params to params[:0] so a recycled TextSprite keeps
	// its backing array across a Clear(). Slicing a nil slice is legal and yields
	// nil, so a freshly constructed sprite simply has no backing array yet. nil is
	// a usable empty slice here -- appending to it works -- so this is the intended
	// outcome, not a missing initialisation.
	if ts.params != nil {
		t.Fatalf("a fresh TextSprite has no params backing array yet, got %#v", ts.params)
	}
}

// The point of the params[:0] reslice in loadDefaults is reuse: clearing a
// sprite that already accumulated params must drop the contents while keeping
// the allocation. Without this the reslice would be pointless and a plain nil
// assignment would do.
func TestTextSpriteClearEmptiesParamsButKeepsBackingArray(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	ts := NewTextSprite()
	ts.params = make([]interface{}, 0, 8)
	ts.params = append(ts.params, "a", "b")
	before := &ts.params[:1][0]

	ts.Clear()

	if len(ts.params) != 0 {
		t.Fatalf("Clear() should empty params, got %#v", ts.params)
	}
	if cap(ts.params) != 8 {
		t.Fatalf("Clear() should keep the params capacity, got %d want 8", cap(ts.params))
	}
	ts.params = append(ts.params, "c")
	if &ts.params[0] != before {
		t.Fatal("Clear() should reuse the existing params backing array")
	}
}
