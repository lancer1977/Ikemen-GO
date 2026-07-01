package main

import "testing"

func TestConstructorsInitializeEmptyState(t *testing.T) {
	p := newPaldata()
	if p == nil {
		t.Fatal("newPaldata returned nil")
	}
	if p.palList.palettes != nil || p.palList.paletteMap != nil || p.palList.PalTex != nil {
		t.Fatalf("newPaldata should leave palette storage empty, got %#v", p.palList)
	}
	if p.palList.PalTable == nil || p.palList.numcols == nil {
		t.Fatalf("newPaldata should initialize palette maps, got %#v", p.palList)
	}

	s := newSff()
	if s == nil {
		t.Fatal("newSff returned nil")
	}
	if s.sprites == nil {
		t.Fatal("newSff should initialize sprite map")
	}
	if s.palList.PalTable == nil || s.palList.numcols == nil {
		t.Fatalf("newSff should initialize palette storage, got %#v", s.palList)
	}
}
