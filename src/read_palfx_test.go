package main

import "testing"

func TestReadPalFX_PopulatesFieldsAndReturnsInitTime(t *testing.T) {
	is := IniSection{
		"fx.time":        " 5 ",
		"fx.add":         "1, 2, 3",
		"fx.mul":         "4, 5, 6",
		"fx.sinadd":      "7, 8, 9, -10",
		"fx.sinmul":      "11, 12, 13, 14",
		"fx.sincolor":    "15, -16",
		"fx.sinhue":      "17, 18",
		"fx.invertall":   "1",
		"fx.invertblend": "2",
		"fx.color":       "256",
		"fx.hue":         "512",
	}
	pfx := newPalFX()
	initTime := ReadPalFX("fx.", is, pfx)

	if initTime != 5 || pfx.time != 5 {
		t.Fatalf("ReadPalFX time = %d/%d, want 5/5", initTime, pfx.time)
	}
	if pfx.add != [3]int32{1, 2, 3} || pfx.mul != [3]int32{4, 5, 6} {
		t.Fatalf("ReadPalFX add/mul = %#v %#v", pfx.add, pfx.mul)
	}
	if pfx.sinadd != [3]int32{-7, -8, -9} || pfx.cycletime[0] != 10 {
		t.Fatalf("ReadPalFX sinadd/cycletime = %#v %#v", pfx.sinadd, pfx.cycletime)
	}
	if pfx.sinmul != [3]int32{11, 12, 13} || pfx.cycletime[1] != 14 {
		t.Fatalf("ReadPalFX sinmul/cycletime = %#v %#v", pfx.sinmul, pfx.cycletime)
	}
	if pfx.sincolor != -15 || pfx.cycletime[2] != 16 {
		t.Fatalf("ReadPalFX sincolor/cycletime = %d %#v", pfx.sincolor, pfx.cycletime)
	}
	if pfx.sinhue != 17 || pfx.cycletime[3] != 18 {
		t.Fatalf("ReadPalFX sinhue/cycletime = %d %#v", pfx.sinhue, pfx.cycletime)
	}
	if !pfx.invertall || pfx.invertblend != 2 {
		t.Fatalf("ReadPalFX invert flags = %#v %d", pfx.invertall, pfx.invertblend)
	}
	if pfx.color != 1 || pfx.hue != 1 {
		t.Fatalf("ReadPalFX color/hue = %v %v, want 1 1", pfx.color, pfx.hue)
	}
}

func TestReadPalFX_DefaultsWhenInputsAreMissing(t *testing.T) {
	pfx := newPalFX()
	initTime := ReadPalFX("missing.", IniSection{}, pfx)
	if initTime != -1 || pfx.time != -1 {
		t.Fatalf("missing ReadPalFX time = %d/%d, want -1/-1", initTime, pfx.time)
	}
	if pfx.color != 1 || pfx.mul != [3]int32{256, 256, 256} {
		t.Fatalf("missing ReadPalFX defaults = %#v", pfx)
	}
}

func TestPalFXSetColorClampsInputs(t *testing.T) {
	pfx := newPalFX()
	pfx.setColor(-10, 128, 999)

	if !pfx.enable {
		t.Fatal("setColor should enable the effect")
	}
	if pfx.eColor != 1 || pfx.eHue != 0 {
		t.Fatalf("setColor should reset color/hue, got %v/%v", pfx.eColor, pfx.eHue)
	}
	if pfx.eMul != [3]int32{0, 128, 255} {
		t.Fatalf("setColor should clamp and scale channels, got %#v", pfx.eMul)
	}
}
