package main

import "testing"

func TestNewLoaderInitializesReadyState(t *testing.T) {
	l := newLoader()
	if l == nil {
		t.Fatal("newLoader returned nil")
	}
	if l.state != LS_NotYet {
		t.Fatalf("newLoader state = %v, want %v", l.state, LS_NotYet)
	}
	if l.loadExit == nil {
		t.Fatal("newLoader should allocate loadExit channel")
	}
	select {
	case l.loadExit <- LS_Loading:
	default:
		t.Fatal("newLoader loadExit channel should be buffered")
	}
}
