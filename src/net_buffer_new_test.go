package main

import "testing"

func TestNewNetBuffer(t *testing.T) {
	t.Parallel()

	nb := NewNetBuffer()
	if nb.InputReader == nil {
		t.Fatal("expected InputReader to be initialized")
	}
	if nb.curT != 0 || nb.inpT != 0 || nb.senT != 0 {
		t.Fatalf("unexpected timer defaults: %#v", nb)
	}
}
