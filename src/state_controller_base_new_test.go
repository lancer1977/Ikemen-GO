package main

import "testing"

func TestNewStateControllerBase(t *testing.T) {
	t.Parallel()

	scb := newStateControllerBase()
	if scb == nil {
		t.Fatal("newStateControllerBase returned nil")
	}
	if len(*scb) != 0 {
		t.Fatalf("expected empty bytecode slice, got len=%d", len(*scb))
	}
	if scb == nil {
		t.Fatal("unexpected nil state controller base")
	}
}
