package main

import "testing"

func TestSelectEnsurePreloadSlots(t *testing.T) {
	sel := &Select{}

	sel.ensureCharPreloadSlot(-1)
	if len(sel.charPreload) != 0 {
		t.Fatalf("ensureCharPreloadSlot(-1) should no-op, got %#v", sel.charPreload)
	}
	sel.ensureCharPreloadSlot(2)
	if len(sel.charPreload) != 3 {
		t.Fatalf("ensureCharPreloadSlot(2) should expand to 3, got %d", len(sel.charPreload))
	}

	sel.ensureStagePreloadSlot(0)
	if len(sel.stagePreload) != 0 {
		t.Fatalf("ensureStagePreloadSlot(0) should no-op, got %#v", sel.stagePreload)
	}
	sel.ensureStagePreloadSlot(3)
	if len(sel.stagePreload) != 3 {
		t.Fatalf("ensureStagePreloadSlot(3) should expand to 3, got %d", len(sel.stagePreload))
	}
}
