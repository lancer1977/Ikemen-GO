package main

import "testing"

func TestAnimLayout_HasAnim_RequiresNonNilAnimationWithFrames(t *testing.T) {
	var al *AnimLayout
	if al.HasAnim() {
		t.Fatal("nil AnimLayout should not report an animation")
	}

	al = &AnimLayout{}
	if al.HasAnim() {
		t.Fatal("AnimLayout without animation should not report an animation")
	}

	al.anim = &Animation{}
	if al.HasAnim() {
		t.Fatal("AnimLayout with empty animation should not report an animation")
	}

	al.anim.frames = []AnimFrame{{}}
	if !al.HasAnim() {
		t.Fatal("AnimLayout with frames should report an animation")
	}
}
