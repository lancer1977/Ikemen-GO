package main

import "testing"

func TestPaletteListRemapHelpers_CopySwapAndResetMappings(t *testing.T) {
	pl := &PaletteList{}
	pl.init()
	p0 := make([]uint32, 256)
	p1 := make([]uint32, 256)
	pl.SetSource(0, p0)
	pl.SetSource(1, p1)

	pl.Remap(0, 1)
	pl.Remap(1, 0)
	if got := pl.GetPalMap(); len(got) != 2 || got[0] != 1 || got[1] != 0 {
		t.Fatalf("GetPalMap() = %#v", got)
	}
	copyMap := pl.GetPalMap()
	copyMap[0] = 7
	if got := pl.GetPalMap(); len(got) != 2 || got[0] != 1 || got[1] != 0 {
		t.Fatalf("GetPalMap() should return a copy, got %#v after mutation", got)
	}

	swap := []int{9, 8}
	if !pl.SwapPalMap(&swap) {
		t.Fatal("SwapPalMap should succeed on same-length slices")
	}
	if len(swap) != 2 || swap[0] != 1 || swap[1] != 0 {
		t.Fatalf("SwapPalMap did not swap out current map: %#v", swap)
	}
	if got := pl.GetPalMap(); len(got) != 2 || got[0] != 9 || got[1] != 8 {
		t.Fatalf("SwapPalMap did not install incoming map: %#v", got)
	}

	bad := []int{1}
	if pl.SwapPalMap(&bad) {
		t.Fatal("SwapPalMap should fail on mismatched lengths")
	}
	if len(bad) != 1 || bad[0] != 1 {
		t.Fatalf("SwapPalMap should not mutate failed input: %#v", bad)
	}

	pl.ResetRemap()
	if got := pl.GetPalMap(); len(got) != 2 || got[0] != 0 || got[1] != 1 {
		t.Fatalf("ResetRemap() = %#v, want identity map", got)
	}
}
