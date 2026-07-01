package main

import "testing"

func TestSelectedCharKeys_UsesP1SelectionOrderAndBasenames(t *testing.T) {
	prevSys := sys
	defer func() { sys = prevSys }()

	sys = System{}
	sys.sel.charlist = []SelectChar{
		{def: "chars/ryu/ryu.def"},
		{def: "chars/ken/ken.def"},
	}
	sys.sel.selected[0] = [][2]int{{1, 0}, {0, 2}}

	if got := selectedCharKeys(); len(got) != 2 || got[0] != "ken" || got[1] != "ryu" {
		t.Fatalf("unexpected selected char keys: %#v", got)
	}
}

func TestLiveCharacterForSelection_RequiresMatchingSideMemberAndSlot(t *testing.T) {
	s := &System{}
	live := &Char{teamside: 1, memberNo: 2, selectNo: 3, name: "Ken"}
	s.chars[0] = []*Char{live}

	if got := s.liveCharacterForSelection(1, 2, 3); got != live {
		t.Fatalf("expected matching live character to be returned")
	}
	if got := s.liveCharacterForSelection(0, 2, 3); got != nil {
		t.Fatalf("expected mismatched side to return nil, got %#v", got)
	}
	if got := s.liveCharacterForSelection(1, 9, 3); got != nil {
		t.Fatalf("expected mismatched member to return nil, got %#v", got)
	}
}
