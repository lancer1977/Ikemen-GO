package main

import "testing"

func TestApplyMapOverridesCopiesCharacterMapEntries(t *testing.T) {
	c := &Char{}
	c.ocd().maps = map[string]float32{"alpha": 1.5, "beta": 2.5}

	c.applyMapOverrides()
	if len(c.mapArray) != 2 || c.mapArray["alpha"] != 1.5 || c.mapArray["beta"] != 2.5 {
		t.Fatalf("applyMapOverrides() = %#v", c.mapArray)
	}

	c.ocd().maps["alpha"] = 3.5
	if c.mapArray["alpha"] != 1.5 {
		t.Fatalf("applyMapOverrides() should copy values, got %#v", c.mapArray)
	}
}

func TestApplyMapOverridesNoOpsWithoutMapData(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	// applyMapOverrides lazily creates sys.sel.gameParams via the global sys, so
	// this test mutates shared state. The deferred restore above already puts the
	// whole System back; forcing gameParams to nil in a t.Cleanup would run after
	// that restore and leak a nil into the tests that follow.
	sys = *newTestSystem()

	c := &Char{}
	c.applyMapOverrides()
	if c.mapArray != nil {
		t.Fatalf("applyMapOverrides() without map data should leave mapArray nil, got %#v", c.mapArray)
	}
}
