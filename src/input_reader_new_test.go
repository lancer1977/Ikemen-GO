package main

import "testing"

func TestNewInputReader(t *testing.T) {
	t.Parallel()

	ir := NewInputReader()
	if ir == nil {
		t.Fatal("NewInputReader returned nil")
	}
	if ir.SocdFirst != [4]bool{} || ir.SocdAllow != [4]bool{} {
		t.Fatalf("expected zeroed SOCD state, got %#v", ir)
	}
}
