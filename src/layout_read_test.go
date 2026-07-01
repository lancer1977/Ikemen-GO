package main

import "testing"

func TestNewLayoutAndRead_DefaultsAndWindowNormalization(t *testing.T) {
	oldScrrect := sys.scrrect
	defer func() { sys.scrrect = oldScrrect }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	l := newLayout(3)
	if l.facing != 1 || l.vfacing != 1 || l.layerno != 3 {
		t.Fatalf("newLayout defaults = %#v", l)
	}
	if l.scale != [2]float32{1, 1} || l.projection != Projection_Orthographic || l.fLength != 2048 {
		t.Fatalf("newLayout defaults = %#v", l)
	}

	is := IniSection{
		"lay.offset":      "10, 20",
		"lay.facing":      "-1",
		"lay.vfacing":     "0",
		"lay.layerno":     "5",
		"lay.scale":       "1.5, 2.5",
		"lay.angle":       "45",
		"lay.xangle":      "10",
		"lay.yangle":      "20",
		"lay.xshear":      "3",
		"lay.focallength": "1024",
		"lay.projection":  "perspective2",
		"lay.window":      "50, 75, 20, 10",
	}

	l.Read("lay.", is)
	if l.offset != [2]float32{10, 20} || l.facing != -1 || l.vfacing != 1 {
		t.Fatalf("Layout.Read basic fields = %#v", l)
	}
	if l.layerno != 2 {
		t.Fatalf("Layout.Read layerno = %d, want 2", l.layerno)
	}
	if l.scale != [2]float32{1.5, 2.5} || l.rot.angle != 45 || l.rot.xangle != 10 || l.rot.yangle != 20 || l.xshear != 3 || l.fLength != 1024 {
		t.Fatalf("Layout.Read numeric fields = %#v", l)
	}
	if l.projection != Projection_Perspective2 {
		t.Fatalf("Layout.Read projection = %v, want perspective2", l.projection)
	}
	if l.window != [4]int32{20, 10, 30, 65} {
		t.Fatalf("Layout.Read normalized window = %#v", l.window)
	}

	l = *newLayout(4)
	l.Read("alt.", IniSection{
		"alt.projection": "perspective",
	})
	if l.projection != Projection_Perspective {
		t.Fatalf("Layout.Read perspective projection = %v, want perspective", l.projection)
	}

	l = *newLayout(1)
	l.Read("missing.", IniSection{})
	if l.window != sys.scrrect {
		t.Fatalf("Layout.Read missing window = %#v, want %#v", l.window, sys.scrrect)
	}

	l = *newLayout(2)
	l.Read("bad.", IniSection{
		"bad.projection": "sideways",
	})
	if l.projection != Projection_Orthographic {
		t.Fatalf("Layout.Read unknown projection = %v, want orthographic", l.projection)
	}
}
