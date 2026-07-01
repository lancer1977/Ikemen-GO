package main

import "testing"

func TestNewGameState(t *testing.T) {
	t.Parallel()

	gs := NewGameState()
	if gs == nil {
		t.Fatal("NewGameState returned nil")
	}
	if gs.id <= 0 {
		t.Fatalf("id = %d, want positive timestamp-derived value", gs.id)
	}
	if gs.saved {
		t.Fatal("expected new game state to start unsaved")
	}
}
