package main

import "testing"

func TestNewFightScreenCombo(t *testing.T) {
	t.Parallel()

	co := newFightScreenCombo()
	if co == nil {
		t.Fatal("newFightScreenCombo returned nil")
	}
	if co.counter == nil || co.text == nil {
		t.Fatalf("newFightScreenCombo should allocate maps: %#v", co)
	}
	if co.displaytime != 90 || co.showspeed != 8 || co.hidespeed != 4 {
		t.Fatalf("unexpected timing defaults: %#v", co)
	}
	if !co.autoalign {
		t.Fatal("newFightScreenCombo should default autoalign to true")
	}
	if co.counterShake.freq != 60 || co.counterShake.decay != 1.0 || co.counterShake.scale != 1.0 {
		t.Fatalf("unexpected counter shake defaults: %#v", co.counterShake)
	}
	if co.textShake.freq != 60 || co.textShake.decay != 1.0 || co.textShake.scale != 1.0 {
		t.Fatalf("unexpected text shake defaults: %#v", co.textShake)
	}
}
