package main

import "testing"

func TestNewGameStatePool(t *testing.T) {
	t.Parallel()

	p := NewGameStatePool()
	if p.curStateID != 0 {
		t.Fatalf("curStateID = %d, want 0", p.curStateID)
	}
	// NewGameStatePool initializes poolObjs as an empty map, not nil.
	// An initialized empty map and nil behave identically for reads but differently
	// for writes (nil would panic, empty map succeeds). See state.go:556.
	if len(p.poolObjs) != 0 {
		t.Fatalf("poolObjs should be empty, got %d entries", len(p.poolObjs))
	}
	if got := p.gameStatePool.New; got == nil {
		t.Fatal("expected gameStatePool.New to be initialized")
	}
	if got := p.stringIntMapPool.New; got == nil {
		t.Fatal("expected stringIntMapPool.New to be initialized")
	}
	if got := p.animationTablePool.New; got == nil {
		t.Fatal("expected animationTablePool.New to be initialized")
	}
}
