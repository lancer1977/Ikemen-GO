package main

import "testing"

func TestReadPalFX_NormalizesNegativeCycleTimingAcrossFields(t *testing.T) {
	is := IniSection{
		"fx.sinadd":   "1, 2, 3, -4",
		"fx.sinmul":   "5, 6, 7, -8",
		"fx.sincolor": "9, -10",
		"fx.sinhue":   "11, -12",
		"fx.color":    "128",
		"fx.hue":      "256",
	}

	pfx := newPalFX()
	initTime := ReadPalFX("fx.", is, pfx)

	if initTime != -1 || pfx.time != -1 {
		t.Fatalf("ReadPalFX time = %d/%d, want -1/-1", initTime, pfx.time)
	}
	if pfx.sinadd != [3]int32{-1, -2, -3} || pfx.cycletime[0] != 4 {
		t.Fatalf("ReadPalFX sinadd/cycletime[0] = %#v/%d, want [-1 -2 -3]/4", pfx.sinadd, pfx.cycletime[0])
	}
	if pfx.sinmul != [3]int32{-5, -6, -7} || pfx.cycletime[1] != 8 {
		t.Fatalf("ReadPalFX sinmul/cycletime[1] = %#v/%d, want [-5 -6 -7]/8", pfx.sinmul, pfx.cycletime[1])
	}
	if pfx.sincolor != -9 || pfx.cycletime[2] != 10 {
		t.Fatalf("ReadPalFX sincolor/cycletime[2] = %d/%d, want -9/10", pfx.sincolor, pfx.cycletime[2])
	}
	if pfx.sinhue != -11 || pfx.cycletime[3] != 12 {
		t.Fatalf("ReadPalFX sinhue/cycletime[3] = %d/%d, want -11/12", pfx.sinhue, pfx.cycletime[3])
	}
	if pfx.color != 0.5 || pfx.hue != 0.5 {
		t.Fatalf("ReadPalFX color/hue = %v/%v, want 0.5/0.5", pfx.color, pfx.hue)
	}
}
