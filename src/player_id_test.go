package main

import "testing"

func TestInitPlayerIDAssignsAndPreservesCharacterIds(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.charList.idMap = make(map[int32]*Char)
	sys.chars = [MaxPlayerNo][]*Char{}
	sys.cfg.Config.HelperMax = 2
	sys.round = 1

	p0 := &Char{}
	p1 := &Char{}
	sys.chars[0] = []*Char{p0}
	sys.chars[1] = []*Char{p1}
	sys.initPlayerID()
	if p0.id != 2 || p1.id != 3 {
		t.Fatalf("initPlayerID() round 1 ids = %d %d, want 2 3", p0.id, p1.id)
	}

	sys.round = 2
	sys.lastCharId = 5
	sys.charList.idMap[6] = p0
	if got := sys.newCharId(); got != 7 {
		t.Fatalf("newCharId() should skip conflicting id, got %d", got)
	}
}

func TestInitPlayerIDPreservesTurnsLoadingPromotedIds(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.charList.idMap = make(map[int32]*Char)
	sys.chars = [MaxPlayerNo][]*Char{}
	sys.cfg.Config.HelperMax = 2
	sys.cfg.Config.TurnsLoading = true
	sys.tmode[0] = TM_Turns
	sys.round = 1

	p0 := &Char{id: 9}
	sys.chars[0] = []*Char{p0}
	sys.initPlayerID()
	if p0.id != 9 {
		t.Fatalf("initPlayerID() should preserve promoted Turns id, got %d", p0.id)
	}
}
