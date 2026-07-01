package main

import "testing"

func TestReadAnimLayout_UsesSpriteAndTableAnimationPaths(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.scrrect = [4]int32{0, 0, 640, 480}

	sff := &Sff{}
	at := NewAnimationTable()
	at.anims[7] = &Animation{frames: []AnimFrame{{Time: 1}}}

	bySpr := ReadAnimLayout("lay.", IniSection{
		"lay.spr":         "3, 4",
		"lay.palfx.time":  "5",
		"lay.offset":      "9, 10",
		"lay.layerno":     "2",
		"lay.focallength": "1024",
	}, sff, at, 6)
	if bySpr.anim == nil || len(bySpr.anim.frames) != 1 {
		t.Fatalf("ReadAnimLayout(spr) anim = %#v", bySpr.anim)
	}
	if bySpr.anim.frames[0].Group != 3 || bySpr.anim.frames[0].Number != 4 {
		t.Fatalf("ReadAnimLayout(spr) frame = %#v", bySpr.anim.frames[0])
	}
	if bySpr.lay.layerno != 2 || bySpr.lay.offset != [2]float32{9, 10} || bySpr.lay.fLength != 1024 {
		t.Fatalf("ReadAnimLayout(spr) layout = %#v", bySpr.lay)
	}

	byAnim := ReadAnimLayout("lay.", IniSection{
		"lay.anim": "7",
	}, sff, at, 6)
	if byAnim.anim != at.anims[7] {
		t.Fatalf("ReadAnimLayout(anim) = %#v, want table animation", byAnim.anim)
	}
	if byAnim.lay.layerno != 6 {
		t.Fatalf("ReadAnimLayout(anim) layerno = %d, want 6", byAnim.lay.layerno)
	}

	missing := ReadAnimLayout("lay.", IniSection{
		"lay.anim": "99",
	}, sff, at, 8)
	if missing.anim == nil || missing.anim == at.anims[7] {
		t.Fatalf("ReadAnimLayout(missing anim) = %#v, want default animation", missing.anim)
	}
	if missing.lay.layerno != 8 {
		t.Fatalf("ReadAnimLayout(missing anim) layerno = %d, want 8", missing.lay.layerno)
	}
}
