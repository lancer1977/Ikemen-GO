package main

import "testing"

func TestAnimTextSnd_NoSoundAndHasDrawableReflectState(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.fightScreen.fnt = map[int]*Fnt{}

	ats := newAnimTextSnd(&Sff{}, 0)
	if !ats.NoSound() {
		t.Fatal("new AnimTextSnd should start with no sound")
	}
	if ats.HasDrawable() {
		t.Fatal("new AnimTextSnd should not be drawable")
	}

	ats.snd = [2]int32{1, 2}
	if ats.NoSound() {
		t.Fatal("AnimTextSnd with snd set should not report NoSound")
	}

	ats.text.font[0] = 3
	ats.text.text = "hello"
	sys.fightScreen.fnt[3] = &Fnt{}
	if !ats.HasDrawable() {
		t.Fatal("AnimTextSnd with text and font should be drawable")
	}

	ats.text.text = ""
	ats.animLayout.anim = &Animation{frames: []AnimFrame{{}}}
	if !ats.HasDrawable() {
		t.Fatal("AnimTextSnd with animation frames should be drawable")
	}
}

func TestAnimTextSnd_EndHonorsDisplayTimeAndAnimationState(t *testing.T) {
	ats := newAnimTextSnd(&Sff{}, 0)
	if !ats.End(0, false) {
		t.Fatal("empty AnimTextSnd with negative displaytime should end")
	}

	ats.displaytime = 2
	if ats.End(2, false) {
		t.Fatal("End should still be false when dt == displaytime")
	}
	if !ats.End(3, false) {
		t.Fatal("End should be true when dt exceeds displaytime")
	}

	ats.displaytime = -1
	ats.animLayout.anim = &Animation{frames: []AnimFrame{{Time: 1}}}
	ats.animLayout.anim.curelem = 0
	if ats.End(0, false) {
		t.Fatal("running animation should not end while frames remain")
	}
	ats.animLayout.anim.curelem = 0
	ats.animLayout.anim.frames[0].Time = 0
	if !ats.End(0, false) {
		t.Fatal("finished animation should end when inf is false")
	}
}
