package main

import "testing"

func TestPaletteListSelectablePalIndex_RespectsConfiguredBoundsAndMappings(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.cfg.Config.PaletteMax = 6

	pl := &PaletteList{
		PalTable: map[[2]uint16]int{
			{1, 2}: 7,
		},
	}

	if got := pl.SelectablePalIndex(0); got != 0 {
		t.Fatalf("SelectablePalIndex(low) = %d, want 0", got)
	}
	if got := pl.SelectablePalIndex(9); got != 0 {
		t.Fatalf("SelectablePalIndex(high) = %d, want 0", got)
	}
	if got := pl.SelectablePalIndex(3); got != 0 {
		t.Fatalf("SelectablePalIndex(missing) = %d, want 0", got)
	}
	if got := pl.SelectablePalIndex(2); got != 7 {
		t.Fatalf("SelectablePalIndex(mapped) = %d, want 7", got)
	}
}
