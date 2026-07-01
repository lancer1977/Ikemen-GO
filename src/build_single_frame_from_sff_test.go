package main

import "testing"

func TestBuildSingleFrameFromSFF(t *testing.T) {
	t.Parallel()

	sff := newSff()
	sff.sprites[[2]uint16{1, 2}] = &Sprite{Group: 1, Number: 2}

	anim := buildSingleFrameFromSFF(sff, 1, 2)
	if anim == nil {
		t.Fatal("expected animation for existing sprite")
	}
	if anim.mask != 0 {
		t.Fatalf("anim.mask = %d, want 0", anim.mask)
	}
	if len(anim.frames) != 1 {
		t.Fatalf("frames len = %d, want 1", len(anim.frames))
	}
	if anim.frames[0].Group != 1 || anim.frames[0].Number != 2 || anim.frames[0].Time != 1 {
		t.Fatalf("unexpected frame: %#v", anim.frames[0])
	}

	if got := buildSingleFrameFromSFF(nil, 1, 2); got != nil {
		t.Fatal("nil sff should return nil")
	}
	if got := buildSingleFrameFromSFF(sff, 9, 9); got != nil {
		t.Fatal("missing sprite should return nil")
	}
}
