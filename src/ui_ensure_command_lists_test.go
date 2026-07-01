package main

import "testing"

func TestUIEnsureCommandListsTrimsGrowsAndRepairsLists(t *testing.T) {
	s := &System{}
	s.commandLists = []*CommandList{
		nil,
		{Buffer: nil, Names: map[string]int{}},
	}

	if err := s.uiEnsureCommandLists(1); err != nil {
		t.Fatalf("uiEnsureCommandLists() trim = %v", err)
	}
	if len(s.commandLists) != 1 {
		t.Fatalf("uiEnsureCommandLists() trim length = %d, want 1", len(s.commandLists))
	}

	if err := s.uiEnsureCommandLists(3); err != nil {
		t.Fatalf("uiEnsureCommandLists() grow = %v", err)
	}
	if len(s.commandLists) != 3 {
		t.Fatalf("uiEnsureCommandLists() grow length = %d, want 3", len(s.commandLists))
	}
	for i, cl := range s.commandLists {
		if cl == nil || cl.Buffer == nil {
			t.Fatalf("uiEnsureCommandLists() should repair nil/bufferless entries, got index %d => %#v", i, cl)
		}
	}
	if _, ok := s.uiCommandRegistry["D"]; !ok {
		t.Fatal("uiEnsureCommandLists() should seed default UI commands")
	}
	if _, ok := s.commandLists[0].Names["D"]; !ok {
		t.Fatal("uiEnsureCommandLists() should apply seeded defaults to new command lists")
	}
	if _, ok := s.commandLists[0].Names["a"]; !ok {
		t.Fatal("uiEnsureCommandLists() should apply action commands to new command lists")
	}
}
