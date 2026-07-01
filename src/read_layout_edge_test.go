package main

import "testing"

func TestReadLayout_ClampsLayerNormalizesWindowAndProjection(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	l := ReadLayout("lay.", IniSection{
		"lay.offset":      "10, 20",
		"lay.facing":      "-1",
		"lay.vfacing":     "0",
		"lay.layerno":     "9",
		"lay.scale":       "1.5, 2.5",
		"lay.angle":       "3.5",
		"lay.xangle":      "4.5",
		"lay.yangle":      "5.5",
		"lay.xshear":      "0.5",
		"lay.focallength": "1234",
		"lay.projection":  "perspective2",
		"lay.window":      "10, 20, 110, 220",
	}, 6)

	if l.layerno != 2 {
		t.Fatalf("ReadLayout layerno = %d, want 2", l.layerno)
	}
	if l.facing != -1 || l.vfacing != 1 {
		t.Fatalf("ReadLayout facing flags = %#v", l)
	}
	if l.offset != [2]float32{10, 20} || l.scale != [2]float32{1.5, 2.5} {
		t.Fatalf("ReadLayout offset/scale = %#v", l)
	}
	if l.rot.angle != 3.5 || l.rot.xangle != 4.5 || l.rot.yangle != 5.5 || l.xshear != 0.5 {
		t.Fatalf("ReadLayout rotation/shear = %#v", l)
	}
	if l.fLength != 1234 || l.projection != Projection_Perspective2 {
		t.Fatalf("ReadLayout projection/fLength = %#v", l)
	}
	if l.window != [4]int32{10, 20, 100, 200} {
		t.Fatalf("ReadLayout window = %#v", l.window)
	}
}

func TestReadLayout_DefaultsToScreenRectWhenWindowMissing(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{1, 2, 3, 4}

	l := ReadLayout("lay.", IniSection{}, 8)
	if l.layerno != 8 {
		t.Fatalf("ReadLayout missing layerno = %d, want 8", l.layerno)
	}
	if l.window != sys.scrrect {
		t.Fatalf("ReadLayout missing window = %#v, want %#v", l.window, sys.scrrect)
	}
}
