package main

import "testing"

func TestSelectConstructorsInitializeDefaults(t *testing.T) {
	sc := newSelectChar()
	if sc == nil {
		t.Fatal("newSelectChar returned nil")
	}
	if sc.localcoord != [2]int32{320, 240} || sc.portraitscale != 1 {
		t.Fatalf("unexpected SelectChar defaults: %#v", sc)
	}
	if sc.cns_scale != [2]float32{1, 1} {
		t.Fatalf("unexpected SelectChar scale defaults: %#v", sc.cns_scale)
	}
	if sc.anims == nil || sc.scp == nil {
		t.Fatalf("newSelectChar should initialize anims and params: %#v", sc)
	}

	ss := newSelectStage()
	if ss == nil {
		t.Fatal("newSelectStage returned nil")
	}
	if ss.localcoord != [2]int32{320, 240} || ss.portraitscale != 1 {
		t.Fatalf("unexpected SelectStage defaults: %#v", ss)
	}
	if ss.anims == nil || ss.ssp == nil {
		t.Fatalf("newSelectStage should initialize anims and params: %#v", ss)
	}
}
