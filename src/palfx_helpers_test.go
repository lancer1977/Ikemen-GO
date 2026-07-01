package main

import "testing"

func TestPalFXClearAndGetSynFx(t *testing.T) {
	pfx := newPalFX()
	pfx.allowNeg = true
	pfx.sintime = [4]int32{1, 2, 3, 4}

	pfx.clear()
	if pfx.allowNeg {
		t.Fatalf("clear() should disable negative color math")
	}
	if pfx.sintime != [4]int32{} {
		t.Fatalf("clear() should reset sintime, got %#v", pfx.sintime)
	}
	if pfx.color != 1 || pfx.mul != [3]int32{256, 256, 256} {
		t.Fatalf("clear() should restore defaults, got %#v", pfx.PalFXDef)
	}

	pfx.clearWithNeg(true)
	if !pfx.allowNeg {
		t.Fatalf("clearWithNeg(true) should preserve negative math")
	}

	orig := sys.allPalFX
	defer func() { sys.allPalFX = orig }()

	sys.allPalFX = newPalFX()
	if got := pfx.getSynFx(TT_add, [2]int32{0, 0}); got != sys.allPalFX {
		t.Fatalf("disabled pfx should defer to sys.allPalFX")
	}

	pfx.enable = true
	sys.allPalFX.enable = false
	if got := pfx.getSynFx(TT_add, [2]int32{0, 0}); got != pfx {
		t.Fatalf("enabled pfx without global FX should return itself")
	}
}

func TestPalFXGetFxPalReturnsInputWhenDisabled(t *testing.T) {
	orig := sys.workpal
	defer func() { sys.workpal = orig }()

	sys.workpal = make([]uint32, 1)
	pfx := newPalFX()
	pal := []uint32{0x11223344}
	if got := pfx.getFxPal(TT_add, pal, false); len(got) != 1 || got[0] != pal[0] {
		t.Fatalf("disabled getFxPal should return input palette, got %#v", got)
	}
}

func TestPalFXGetFinalPalFx(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		orig := sys.allPalFX
		defer func() { sys.allPalFX = orig }()
		sys.allPalFX = newPalFX()

		pfx := newPalFX()
		got := pfx.getFinalPalFx(TT_add, [2]int32{0, 0})
		if got != (ShaderPalFX{neg: false, add: [3]float32{0, 0, 0}, mult: [3]float32{1, 1, 1}, gray: 0, hue: 0, invblend: 0}) {
			t.Fatalf("disabled getFinalPalFx = %#v", got)
		}
	})

	t.Run("enabled", func(t *testing.T) {
		orig := sys.allPalFX
		defer func() { sys.allPalFX = orig }()
		sys.allPalFX = newPalFX()

		pfx := newPalFX()
		pfx.enable = true
		pfx.eInvertall = true
		pfx.eColor = 0.5
		pfx.eHue = 2
		pfx.eAdd = [3]int32{255, 128, 64}
		pfx.eMul = [3]int32{256, 128, 64}

		got := pfx.getFinalPalFx(TT_add, [2]int32{0, 0})
		if !got.neg || got.gray != 0.5 || got.hue != 2 || got.invblend != 0 {
			t.Fatalf("enabled getFinalPalFx unexpected header fields: %#v", got)
		}
		if got.add != [3]float32{1, 128.0 / 255.0, 64.0 / 255.0} {
			t.Fatalf("enabled getFinalPalFx add = %#v", got.add)
		}
		if got.mult != [3]float32{1, 0.5, 0.25} {
			t.Fatalf("enabled getFinalPalFx mult = %#v", got.mult)
		}
	})
}
