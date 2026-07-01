package main

import "testing"

func TestAnimation_GetLengthAndCurrentFrame_HandleEmptyAndMixedDurations(t *testing.T) {
	var a *Animation
	if got := a.GetLength(); got != 0 {
		t.Fatalf("nil animation length = %d, want 0", got)
	}

	a = &Animation{}
	if got := a.GetLength(); got != 0 {
		t.Fatalf("empty animation length = %d, want 0", got)
	}
	if got := a.CurrentFrame(); got != nil {
		t.Fatalf("empty animation CurrentFrame() = %#v, want nil", got)
	}

	a.frames = []AnimFrame{{Time: -1}, {Time: 3}, {Time: 0}}
	if got := a.GetLength(); got != 4 {
		t.Fatalf("mixed animation length = %d, want 4", got)
	}
	if got := a.CurrentFrame(); got != &a.frames[0] {
		t.Fatalf("CurrentFrame() = %#v, want first frame", got)
	}
}

func TestAnimation_SetAnimElem_ClampsAndUpdatesTime(t *testing.T) {
	a := &Animation{
		frames: []AnimFrame{{Time: 2}, {Time: -1}},
	}

	a.SetAnimElem(2, 1)
	if a.curelem != 1 || a.drawidx != 1 {
		t.Fatalf("SetAnimElem(2,1) indexes = curelem %d drawidx %d", a.curelem, a.drawidx)
	}
	if a.curelemtime != 1 {
		t.Fatalf("SetAnimElem(2,1) curelemtime = %d, want 1", a.curelemtime)
	}
	if a.curtime != -1 {
		t.Fatalf("SetAnimElem(2,1) curtime = %d, want -1", a.curtime)
	}

	a.SetAnimElem(9, -5)
	if a.curelem != 0 || a.drawidx != 0 {
		t.Fatalf("SetAnimElem clamp indexes = curelem %d drawidx %d", a.curelem, a.drawidx)
	}
	if a.curelemtime != 0 {
		t.Fatalf("SetAnimElem clamp curelemtime = %d, want 0", a.curelemtime)
	}
	if a.curtime != 0 {
		t.Fatalf("SetAnimElem clamp curtime = %d, want 0", a.curtime)
	}
}

func TestAnimation_Action_HandlesNilEmptyAndSkipsZeroDurationFrames(t *testing.T) {
	var a *Animation
	a.Action()

	a = &Animation{}
	a.Action()
	if !a.loopend {
		t.Fatal("empty animation should set loopend")
	}

	a = &Animation{
		frames: []AnimFrame{{Time: 0}, {Time: 2}},
	}
	a.Action()
	if a.curelem != 1 || a.curelemtime != 1 || a.curtime != 1 {
		t.Fatalf("Action advance = curelem %d curelemtime %d curtime %d", a.curelem, a.curelemtime, a.curtime)
	}
	if a.loopend {
		t.Fatal("non-empty animation should not end after one tick")
	}
}

func TestAnimation_AnimElemTimeAndNo_HandleBoundsAndLooping(t *testing.T) {
	a := &Animation{
		frames:      []AnimFrame{{Time: 2}, {Time: 3}, {Time: -1}},
		curelem:     1,
		curelemtime: 1,
		loopstart:   1,
		curtime:     6,
		totaltime:   5,
	}

	if got := a.AnimElemTime(1); got != 5 {
		t.Fatalf("AnimElemTime(1) = %d, want 5", got)
	}
	if got := a.AnimElemTime(4); got != 0 {
		t.Fatalf("AnimElemTime(4) = %d, want 0", got)
	}
	if got := a.AnimElemNo(0); got != 2 {
		t.Fatalf("AnimElemNo(0) = %d, want 2", got)
	}
	if got := a.AnimElemNo(10); got != 3 {
		t.Fatalf("AnimElemNo(10) = %d, want 3", got)
	}
}

func TestAnimation_AlphaToBlend_ResolvesDefaultsAndBrightness(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.brightness = 0.5

	a := &Animation{
		transType:                  TT_default,
		curtrans:                   TT_add,
		interpolate_blend_srcalpha: 200,
		interpolate_blend_dstalpha: 100,
	}
	mode, alpha := a.alphaToBlend()
	if mode != TT_add || alpha != [2]int32{100, 50} {
		t.Fatalf("alphaToBlend default path = %v %v", mode, alpha)
	}

	a.transType = TT_default
	a.curtrans = TT_default
	mode, alpha = a.alphaToBlend()
	if mode != TT_none || alpha != [2]int32{127, 0} {
		t.Fatalf("alphaToBlend fallback path = %v %v", mode, alpha)
	}

	a.transType = TT_sub
	a.srcAlpha = 10
	a.dstAlpha = 20
	mode, alpha = a.alphaToBlend()
	if mode != TT_sub || alpha != [2]int32{5, 20} {
		t.Fatalf("alphaToBlend explicit path = %v %v", mode, alpha)
	}
}

func TestAnimation_UpdateInterpolation_UsesConfiguredNextFrame(t *testing.T) {
	a := &Animation{
		frames: []AnimFrame{
			{Time: 4, Xoffset: 10, Yoffset: 20, Xscale: 1, Yscale: 2, Angle: 3, SrcAlpha: 200, DstAlpha: 100, TransType: TT_add},
			{Time: 4, Xoffset: 18, Yoffset: 32, Xscale: 2, Yscale: 4, Angle: 11, SrcAlpha: 120, DstAlpha: 60, TransType: TT_add},
		},
		drawidx:     0,
		curelemtime: 2,
		loopstart:   0,
		curtrans:    TT_add,
	}
	a.interpolate_offset = []int32{1}
	a.interpolate_scale = []int32{1}
	a.interpolate_angle = []int32{1}
	a.interpolate_blend = []int32{1}

	a.UpdateInterpolation()

	if a.interpolate_offset_x != 4 || a.interpolate_offset_y != 6 {
		t.Fatalf("offset interpolation = %v %v", a.interpolate_offset_x, a.interpolate_offset_y)
	}
	if a.scale_x != 1.5 || a.scale_y != 3 {
		t.Fatalf("scale interpolation = %v %v", a.scale_x, a.scale_y)
	}
	if a.rot.angle != 7 {
		t.Fatalf("angle interpolation = %v", a.rot.angle)
	}
	if a.interpolate_blend_srcalpha != 160 || a.interpolate_blend_dstalpha != 80 {
		t.Fatalf("blend interpolation = %v %v", a.interpolate_blend_srcalpha, a.interpolate_blend_dstalpha)
	}
}
