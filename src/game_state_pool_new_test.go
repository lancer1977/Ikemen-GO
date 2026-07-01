package main

import "testing"

func TestNewGameStatePool(t *testing.T) {
	t.Parallel()

	p := NewGameStatePool()
	if p.curStateID != 0 {
		t.Fatalf("curStateID = %d, want 0", p.curStateID)
	}
	if p.poolObjs != nil {
		t.Fatalf("poolObjs = %#v, want nil", p.poolObjs)
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
