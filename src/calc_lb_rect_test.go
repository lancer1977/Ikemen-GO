package main

import "testing"

func TestLayoutCalcLBRect_ReturnsScreenRectForFullWindowAndScalesOtherwise(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}
	sys.fightScreen.localcoord = [2]int32{320, 240}

	l := &Layout{}
	if got := l.calcLBRect([4]int32{0, 0, 640, 480}); got != sys.scrrect {
		t.Fatalf("calcLBRect(full window) = %#v, want %#v", got, sys.scrrect)
	}

	got := l.calcLBRect([4]int32{10, 20, 100, 50})
	want := [4]int32{20, 40, 200, 100}
	if got != want {
		t.Fatalf("calcLBRect(scaled) = %#v, want %#v", got, want)
	}
}
