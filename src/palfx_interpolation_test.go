package main

import "testing"

func TestPalFXInterpolationUpdate(t *testing.T) {
	pfx := newPalFX()
	pfx.itime = 4
	pfx.eiTime = 2
	pfx.imul = [6]int32{10, 20, 30, 110, 220, 330}
	pfx.mul = [3]int32{256, 128, 64}
	pfx.iadd = [6]int32{1, 2, 3, 11, 12, 13}
	pfx.add = [3]int32{4, 5, 6}
	pfx.icolor = [2]float32{0.5, 1}
	pfx.color = 0.25
	pfx.ihue = [2]float32{2, 6}
	pfx.hue = 1

	pfx.interpolationUpdate()

	if pfx.eiTime != 3 {
		t.Fatalf("expected eiTime to advance to 3, got %d", pfx.eiTime)
	}
	if pfx.eiMul != [3]int32{60, 120, 180} || pfx.eMul != [3]int32{60, 60, 45} {
		t.Fatalf("unexpected interpolated muls: ei=%#v e=%#v", pfx.eiMul, pfx.eMul)
	}
	if pfx.eiAdd != [3]int32{6, 7, 8} || pfx.eAdd != [3]int32{10, 12, 14} {
		t.Fatalf("unexpected interpolated adds: ei=%#v e=%#v", pfx.eiAdd, pfx.eAdd)
	}
	if pfx.eiColor != 0.625 || pfx.eColor != 0.15625 {
		t.Fatalf("unexpected interpolated color: ei=%v e=%v", pfx.eiColor, pfx.eColor)
	}
	if pfx.eiHue != 5 || pfx.eHue != 6 {
		t.Fatalf("unexpected interpolated hue: ei=%v e=%v", pfx.eiHue, pfx.eHue)
	}
}
