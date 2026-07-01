package main

import "testing"

func TestNewAnimFrame_InitializesExpectedDefaults(t *testing.T) {
	af := newAnimFrame()
	if af == nil {
		t.Fatal("newAnimFrame returned nil")
	}
	if af.Time != -1 || af.Group != -1 || af.TransType != TT_none {
		t.Fatalf("newAnimFrame basics = %#v", af)
	}
	if af.SrcAlpha != 255 || af.DstAlpha != 0 || af.Hscale != 1 || af.Vscale != 1 {
		t.Fatalf("newAnimFrame alpha/scale defaults = %#v", af)
	}
	if af.Xscale != 1 || af.Yscale != 1 || af.Angle != 0 {
		t.Fatalf("newAnimFrame scale/angle defaults = %#v", af)
	}
}
