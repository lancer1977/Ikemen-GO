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
	// ReadAnimLayout calls at.get() which returns a copy of the animation, not the original pointer.
	// So byAnim.anim should be a copy with the same frame data.
	if byAnim.anim == nil || len(byAnim.anim.frames) != 1 || byAnim.anim.frames[0].Time != 1 {
		t.Fatalf("ReadAnimLayout(anim) = %#v, want copy of table animation with one frame", byAnim.anim)
	}
	// layerno is capped to Min(2, ln) by Layout.Read, so 6 becomes 2
	if byAnim.lay.layerno != 2 {
		t.Fatalf("ReadAnimLayout(anim) layerno = %d, want 2", byAnim.lay.layerno)
	}

	missing := ReadAnimLayout("lay.", IniSection{
		"lay.anim": "99",
	}, sff, at, 8)
	if missing.anim == nil || missing.anim == at.anims[7] {
		t.Fatalf("ReadAnimLayout(missing anim) = %#v, want default animation", missing.anim)
	}
	// layerno is capped to Min(2, ln) by Layout.Read, so 8 becomes 2
	if missing.lay.layerno != 2 {
		t.Fatalf("ReadAnimLayout(missing anim) layerno = %d, want 2", missing.lay.layerno)
	}
}
