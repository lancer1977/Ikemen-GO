package main

import "testing"

func TestNewCharCompiler(t *testing.T) {
	t.Parallel()

	c := newCharCompiler()
	if c == nil {
		t.Fatal("newCharCompiler returned nil")
	}
	if c.funcs == nil || c.scmap == nil {
		t.Fatalf("newCharCompiler should allocate maps: %#v", c)
	}
	if len(c.scmap) == 0 {
		t.Fatal("newCharCompiler should register state controllers")
	}
}
