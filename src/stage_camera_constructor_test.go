package main

import (
	"math"
	"testing"
)

func TestNewStageCameraInitializesDefaults(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()
	sys.gameWidth = 640

	sc := newStageCamera()
	if sc == nil {
		t.Fatal("newStageCamera returned nil")
	}
	if sc.verticalfollow != 0.2 || sc.tensionvel != 1 || sc.tension != 50 {
		t.Fatalf("unexpected stage camera defaults: %#v", sc)
	}
	if sc.cutlow != math.MinInt32 || sc.localcoord != [2]int32{320, 240} {
		t.Fatalf("unexpected stage camera coordinate defaults: %#v", sc)
	}
	if sc.localscl != 2 {
		t.Fatalf("newStageCamera localscl = %v, want 2", sc.localscl)
	}
	if sc.startzoom != 1 || sc.zoomin != 1 || sc.zoomout != 1 {
		t.Fatalf("unexpected zoom defaults: %#v", sc)
	}
}
