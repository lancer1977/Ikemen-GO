package main

import "testing"

func TestNewGuardAndStunBars(t *testing.T) {
	t.Parallel()

	gb := newGuardBar()
	if gb == nil {
		t.Fatal("newGuardBar returned nil")
	}
	if gb.front == nil || gb.value == nil {
		t.Fatalf("newGuardBar should allocate maps: %#v", gb)
	}
	if gb.midpower != 0 || gb.midpowerMin != 0 || gb.invertfill || gb.scalefill || gb.leaderontop {
		t.Fatalf("unexpected guard bar defaults: %#v", gb)
	}

	sb := newStunBar()
	if sb == nil {
		t.Fatal("newStunBar returned nil")
	}
	if sb.front == nil || sb.value == nil {
		t.Fatalf("newStunBar should allocate maps: %#v", sb)
	}
	if sb.midpower != 0 || sb.midpowerMin != 0 || sb.invertfill || sb.scalefill || sb.leaderontop {
		t.Fatalf("unexpected stun bar defaults: %#v", sb)
	}
}
