package main

import "testing"

func TestPalFXConstructors(t *testing.T) {
	def := newPalFXDef()
	if def == nil {
		t.Fatal("newPalFXDef returned nil")
	}
	if def.color != 1 || def.icolor != [2]float32{1, 1} {
		t.Fatalf("unexpected PalFXDef defaults: %#v", def)
	}
	if def.mul != [3]int32{256, 256, 256} || def.imul != [6]int32{256, 256, 256, 256, 256, 256} {
		t.Fatalf("unexpected PalFXDef mul defaults: %#v", def)
	}

	pfx := newPalFX()
	if pfx == nil {
		t.Fatal("newPalFX returned nil")
	}
	if pfx.enable || pfx.eColor != 0 || pfx.eHue != 0 {
		t.Fatalf("newPalFX should start zeroed, got %#v", pfx)
	}
}
