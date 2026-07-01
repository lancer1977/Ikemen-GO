package main

import "testing"

func TestNewBackGround(t *testing.T) {
	t.Parallel()

	sff := newSff()
	bg := newBackGround(sff)
	if bg == nil {
		t.Fatal("newBackGround returned nil")
	}
	if bg.palfx == nil || bg.anim == nil {
		t.Fatalf("newBackGround should allocate palfx and anim: %#v", bg)
	}
	if bg.delta != [2]float32{1, 1} || bg.xscale != [2]float32{1, 1} || bg.rasterx != [2]float32{1, 1} {
		t.Fatalf("unexpected background defaults: %#v", bg)
	}
	if bg.yscalestart != 100 || bg.fLength != 2048 || bg.projection != Projection_Orthographic {
		t.Fatalf("unexpected background defaults: %#v", bg)
	}
}
