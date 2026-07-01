package main

import "testing"

func TestNewPreloadedAnims(t *testing.T) {
	t.Parallel()

	pa := NewPreloadedAnims()
	if pa == nil {
		t.Fatal("NewPreloadedAnims returned nil")
	}
	if len(pa) != 0 {
		t.Fatalf("expected empty preloaded anim map, got %d entries", len(pa))
	}
}
