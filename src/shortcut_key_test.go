package main

import "testing"

func TestNewShortcutKey(t *testing.T) {
	t.Parallel()

	prevAlt := ModAlt
	prevCtrlAlt := ModCtrlAlt
	prevCtrlAltShift := ModCtrlAltShift
	t.Cleanup(func() {
		ModAlt = prevAlt
		ModCtrlAlt = prevCtrlAlt
		ModCtrlAltShift = prevCtrlAltShift
	})

	ModAlt = 0
	ModCtrlAlt = 0
	ModCtrlAltShift = 0

	sk := NewShortcutKey(KeyEnter, true, false, true)
	if sk == nil {
		t.Fatal("NewShortcutKey returned nil")
	}
	if sk.Key != KeyEnter {
		t.Fatalf("Key = %v, want %v", sk.Key, KeyEnter)
	}
	expectedMod := NewModifierKey(true, false, true)
	if sk.Mod != expectedMod {
		t.Fatalf("Mod = %v, want %v (ctrl+shift)", sk.Mod, expectedMod)
	}
	if ModAlt == 0 || ModCtrlAlt == 0 || ModCtrlAltShift == 0 {
		t.Fatalf("expected modifier cache to be initialized, got %v %v %v", ModAlt, ModCtrlAlt, ModCtrlAltShift)
	}
}
