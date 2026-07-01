package main

import "testing"

func TestNewProjectile(t *testing.T) {
	t.Parallel()

	p := newProjectile()
	if p == nil {
		t.Fatal("newProjectile returned nil")
	}
	if p.id != 0 || p.status != 0 || p.playerno != 0 {
		t.Fatalf("unexpected projectile defaults: %#v", p)
	}
	if p.shader != "" || p.shaderParams != [16]float32{} {
		t.Fatalf("unexpected projectile defaults: %#v", p)
	}
}
