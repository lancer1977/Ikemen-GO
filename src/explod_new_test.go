package main

import "testing"

func TestNewExplod(t *testing.T) {
	t.Parallel()

	e := newExplod()
	if e == nil {
		t.Fatal("newExplod returned nil")
	}
	if e.id != 0 || e.playerno != 0 || e.ownerId != 0 || e.layerno != 0 {
		t.Fatalf("unexpected explod defaults: %#v", e)
	}
	if e.shader != "" || e.shaderParams != [16]float32{} {
		t.Fatalf("unexpected explod defaults: %#v", e)
	}
}
