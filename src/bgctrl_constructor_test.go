package main

import (
	"math"
	"testing"
)

func TestNewBgCtrl(t *testing.T) {
	t.Parallel()

	bgc := newBgCtrl()
	if bgc == nil {
		t.Fatal("newBgCtrl returned nil")
	}
	if bgc.looptime != -1 {
		t.Fatalf("looptime = %d, want -1", bgc.looptime)
	}
	if !math.IsNaN(float64(bgc.x)) || !math.IsNaN(float64(bgc.y)) {
		t.Fatalf("expected x/y to start as NaN, got %v %v", bgc.x, bgc.y)
	}
}
