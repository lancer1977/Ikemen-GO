package main

import "testing"

func TestGameParamsEnsureOverride(t *testing.T) {
	p := &GameParams{}

	// When team < 0, ensureOverride returns a fresh default override without storing it
	got := p.ensureOverride(-1, 0)
	if got == nil || got.life != -1 || got.maps == nil {
		t.Fatalf("ensureOverride(-1, 0) should return a fresh default override, got %#v", got)
	}
	// ocd is a fixed [3][]OverrideCharData, so always has len 3, but slices should remain empty
	if len(p.ocd[0]) != 0 || len(p.ocd[1]) != 0 || len(p.ocd[2]) != 0 {
		t.Fatalf("ensureOverride(-1, 0) should not mutate storage, got %#v", p.ocd)
	}

	// When team >= 0 and within bounds, ensureOverride expands the slice and stores
	got = p.ensureOverride(1, 2)
	if got == nil || got.life != -1 || got.maps == nil {
		t.Fatalf("ensureOverride(1, 2) should return initialized override, got %#v", got)
	}
	// ocd is a fixed array of length 3, p.ocd[1] is expanded to have at least 3 elements
	if len(p.ocd) != 3 || len(p.ocd[1]) != 3 {
		t.Fatalf("ensureOverride(1, 2) should expand storage, got %#v", p.ocd)
	}
	if p.ocd[1][2].life != -1 || p.ocd[1][2].maps == nil {
		t.Fatalf("ensureOverride should initialize stored override, got %#v", p.ocd[1][2])
	}
}
