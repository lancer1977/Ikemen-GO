package main

import "testing"

func TestAnimLayout_ActionAndReset_HandleNilAndForwardReset(t *testing.T) {
	al := &AnimLayout{}
	al.Action()
	al.Reset()

	al.anim = &Animation{
		frames:      []AnimFrame{{Time: 1}},
		curelem:     0,
		drawidx:     0,
		curelemtime: 1,
		curtime:     2,
		newframe:    false,
		loopend:     true,
	}
	al.Reset()
	if al.anim.curelem != 0 || al.anim.drawidx != 0 || al.anim.curelemtime != 0 || al.anim.curtime != 0 {
		t.Fatalf("AnimLayout.Reset() did not reset animation: %#v", al.anim)
	}
	if !al.anim.newframe || al.anim.loopend {
		t.Fatalf("AnimLayout.Reset() did not restore animation flags: %#v", al.anim)
	}
}
