package main

import "testing"

func TestSpriteLz5DecodeEmptyInput(t *testing.T) {
	s := &Sprite{Size: [2]uint16{1, 1}}
	if got := s.Lz5Decode(nil); got != nil {
		t.Fatalf("expected nil passthrough, got %#v", got)
	}
	if got := s.Lz5Decode([]byte{}); len(got) != 0 {
		t.Fatalf("expected empty passthrough, got %#v", got)
	}
}
