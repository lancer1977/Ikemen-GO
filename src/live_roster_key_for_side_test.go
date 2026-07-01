package main

import "testing"

func TestLiveRosterKeyForSide(t *testing.T) {
	s := &System{}

	if got := s.liveRosterKeyForSide(-1); got != "" {
		t.Fatalf("liveRosterKeyForSide(-1) = %q, want blank", got)
	}
	if got := s.liveRosterKeyForSide(0); got != "" {
		t.Fatalf("liveRosterKeyForSide(0) with empty state = %q, want blank", got)
	}

	s.chars[0] = []*Char{{name: " Ryu "}}
	if got := s.liveRosterKeyForSide(0); got != "ryu" {
		t.Fatalf("liveRosterKeyForSide(0) = %q, want ryu", got)
	}
}
