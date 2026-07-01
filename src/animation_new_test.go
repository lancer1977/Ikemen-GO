package main

import "testing"

func TestNewAnimation(t *testing.T) {
	t.Parallel()

	a := newAnimation(nil, nil)
	if a == nil {
		t.Fatal("newAnimation returned nil")
	}
	if a.sff != nil || a.palettedata != nil {
		t.Fatalf("unexpected source pointers: %#v", a)
	}
	if a.mask != -1 || a.transType != TT_default || a.srcAlpha != 255 || a.dstAlpha != 0 || a.curtrans != TT_none {
		t.Fatalf("unexpected animation defaults: %#v", a)
	}
	if a.interpolate_blend_srcalpha != 255 || a.interpolate_blend_dstalpha != 0 || !a.newframe || a.copyAction != -1 || !a.warnMissing {
		t.Fatalf("unexpected animation state: %#v", a)
	}
	if a.remap == nil || len(a.remap) != 0 {
		t.Fatalf("expected empty remap map, got %#v", a.remap)
	}
	if a.start_scale != [2]float32{1, 1} {
		t.Fatalf("start_scale = %#v, want [1 1]", a.start_scale)
	}
}
