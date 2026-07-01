package main

import "testing"

func TestNewSprite(t *testing.T) {
	t.Parallel()

	s := newSprite()
	if s == nil {
		t.Fatal("newSprite returned nil")
	}
	if s.palidx != -1 {
		t.Fatalf("palidx = %d, want -1", s.palidx)
	}
	if !s.isBlank() {
		t.Fatal("expected new sprite to be blank")
	}
}
