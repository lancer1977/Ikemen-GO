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
	// Production interpolates between imul[i+3] (start) and imul[i] (end).
	// With t=0.75, eiMul[i] = (1-0.75)*imul[i+3] + 0.75*imul[i]
	// Then eMul[i] = eiMul[i] * mul[i] / 256
	if pfx.eiMul != [3]int32{35, 70, 105} || pfx.eMul != [3]int32{35, 35, 26} {
		t.Fatalf("unexpected interpolated muls: ei=%#v e=%#v", pfx.eiMul, pfx.eMul)
	}
	// eiAdd[i] = (1-0.75)*iadd[i+3] + 0.75*iadd[i], then eAdd[i] = eiAdd[i] + add[i]
	if pfx.eiAdd != [3]int32{3, 4, 5} || pfx.eAdd != [3]int32{7, 9, 11} {
		t.Fatalf("unexpected interpolated adds: ei=%#v e=%#v", pfx.eiAdd, pfx.eAdd)
	}
	// eiColor = Lerp(icolor[1], icolor[0], 0.75) = 0.25*1 + 0.75*0.5 = 0.25 + 0.375 = 0.625
	// eColor = eiColor * color = 0.625 * 0.25 = 0.15625
	if pfx.eiColor != 0.625 || pfx.eColor != 0.15625 {
		t.Fatalf("unexpected interpolated color: ei=%v e=%v", pfx.eiColor, pfx.eColor)
	}
	// eiHue = Lerp(ihue[1], ihue[0], 0.75) = 0.25*6 + 0.75*2 = 1.5 + 1.5 = 3
	// eHue = eiHue + hue = 3 + 1 = 4
	if pfx.eiHue != 3 || pfx.eHue != 4 {
		t.Fatalf("unexpected interpolated hue: ei=%v e=%v", pfx.eiHue, pfx.eHue)
	}
}
