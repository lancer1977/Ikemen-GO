package main

import "testing"

func TestPaletteListSetSourceAndNewPal_ExpandAndMapIndices(t *testing.T) {
	pl := &PaletteList{}
	pl.init()

	existing := []uint32{9}
	pl.SetSource(0, existing)
	pl.SetSource(-1, []uint32{1})
	if len(pl.palettes) != 1 || len(pl.paletteMap) != 1 || len(pl.PalTex) != 1 {
		t.Fatalf("SetSource(-1) should not change storage, got %#v", pl)
	}
	if len(pl.palettes[0]) != len(existing) || (len(pl.palettes[0]) > 0 && &pl.palettes[0][0] != &existing[0]) || pl.paletteMap[0] != 0 {
		t.Fatalf("SetSource(-1) should not mutate existing palette state, got %#v", pl)
	}

	i, pal := pl.NewPal()
	// Production returns the index of the newly created palette, which is the old length.
	// Since we already have 1 palette (at index 0), the new one is at index 1.
	if i != 1 || pal == nil {
		t.Fatalf("NewPal() = %d, %#v", i, pal)
	}
	if len(pl.palettes) != 2 || len(pl.paletteMap) != 2 || len(pl.PalTex) != 2 {
		t.Fatalf("NewPal() did not initialize storage: %#v", pl)
	}
	if len(pl.palettes[1]) != len(pal) || (len(pl.palettes[1]) > 0 && &pl.palettes[1][0] != &pal[0]) || pl.paletteMap[1] != 1 {
		t.Fatalf("NewPal() did not wire the new palette source: %#v", pl)
	}

	other := make([]uint32, 256)
	pl.SetSource(2, other)
	if len(pl.palettes) != 3 || len(pl.paletteMap) != 3 || len(pl.PalTex) != 3 {
		t.Fatalf("SetSource(2) did not expand storage: %#v", pl)
	}
	if len(pl.palettes[2]) != len(other) || (len(pl.palettes[2]) > 0 && &pl.palettes[2][0] != &other[0]) || pl.paletteMap[2] != 2 {
		t.Fatalf("SetSource(2) did not assign palette source: %#v", pl)
	}
}
