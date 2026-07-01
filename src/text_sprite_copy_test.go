package main

import "testing"

func TestTextSpriteCopy_ReturnsIndependentCloneAndHandlesNil(t *testing.T) {
	if got := (*TextSprite)(nil).Copy(); got != nil {
		t.Fatalf("nil Copy() = %#v, want nil", got)
	}

	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	ts := NewTextSprite()
	ts.text = "hello"
	ts.x = 12
	ts.y = 34
	ts.SetColor(1, 2, 3, 4)
	ts.palfx.setColor(5, 6, 7)

	cp := ts.Copy()
	if cp == nil {
		t.Fatal("Copy() returned nil")
	}
	if cp == ts {
		t.Fatal("Copy() should allocate a distinct object")
	}
	if cp.text != ts.text || cp.x != ts.x || cp.y != ts.y || cp.fnt != ts.fnt {
		t.Fatalf("Copy() did not preserve fields: orig %#v copy %#v", ts, cp)
	}
	if cp.palfx == nil || cp.palfx == ts.palfx {
		t.Fatal("Copy() should deep-copy palfx")
	}

	cp.text = "bye"
	cp.palfx.setColor(9, 9, 9)
	if ts.text != "hello" {
		t.Fatalf("mutating copy should not change original text, got %q", ts.text)
	}
	if !ts.palfx.enable || ts.palfx.eMul != [3]int32{5, 6, 7} {
		t.Fatalf("mutating copy should not change original palfx, got %#v", ts.palfx)
	}
}
