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
	if pl.palettes[0] != existing || pl.paletteMap[0] != 0 {
		t.Fatalf("SetSource(-1) should not mutate existing palette state, got %#v", pl)
	}

	i, pal := pl.NewPal()
	if i != 0 || pal == nil {
		t.Fatalf("NewPal() = %d, %#v", i, pal)
	}
	if len(pl.palettes) != 1 || len(pl.paletteMap) != 1 || len(pl.PalTex) != 1 {
		t.Fatalf("NewPal() did not initialize storage: %#v", pl)
	}
	if pl.palettes[0] != pal || pl.paletteMap[0] != 0 {
		t.Fatalf("NewPal() did not wire the new palette source: %#v", pl)
	}

	other := make([]uint32, 256)
	pl.SetSource(2, other)
	if len(pl.palettes) != 3 || len(pl.paletteMap) != 3 || len(pl.PalTex) != 3 {
		t.Fatalf("SetSource(2) did not expand storage: %#v", pl)
	}
	if pl.palettes[2] != other || pl.paletteMap[2] != 2 {
		t.Fatalf("SetSource(2) did not assign palette source: %#v", pl)
	}
}
