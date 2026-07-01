package main

import "testing"

func TestGameParamsEnsureOverride(t *testing.T) {
	p := &GameParams{}

	got := p.ensureOverride(-1, 0)
	if got == nil || got.life != -1 || got.maps == nil {
		t.Fatalf("ensureOverride(-1, 0) should return a fresh default override, got %#v", got)
	}
	if len(p.ocd) != 0 {
		t.Fatalf("ensureOverride(-1, 0) should not mutate storage, got %#v", p.ocd)
	}

	got = p.ensureOverride(1, 2)
	if got == nil || got.life != -1 || got.maps == nil {
		t.Fatalf("ensureOverride(1, 2) should return initialized override, got %#v", got)
	}
	if len(p.ocd) != 3 || len(p.ocd[1]) != 3 {
		t.Fatalf("ensureOverride(1, 2) should expand storage, got %#v", p.ocd)
	}
	if p.ocd[1][2].life != -1 || p.ocd[1][2].maps == nil {
		t.Fatalf("ensureOverride should initialize stored override, got %#v", p.ocd[1][2])
	}
}
