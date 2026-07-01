package main

import "testing"

func TestPaletteListSwapPalMap(t *testing.T) {
	pl := PaletteList{}
	pl.init()
	pl.SetSource(0, []uint32{1})
	pl.SetSource(1, []uint32{2})
	pl.Remap(0, 1)
	pl.Remap(1, 0)

	other := []int{9, 8}
	if !pl.SwapPalMap(&other) {
		t.Fatal("expected swap to succeed")
	}
	if other[0] != 1 || other[1] != 0 {
		t.Fatalf("unexpected swapped-out map: %#v", other)
	}
	if got := pl.GetPalMap(); len(got) != 2 || got[0] != 9 || got[1] != 8 {
		t.Fatalf("unexpected palette map after swap: %#v", got)
	}

	bad := []int{1}
	if pl.SwapPalMap(&bad) {
		t.Fatal("expected length mismatch to fail")
	}
}
