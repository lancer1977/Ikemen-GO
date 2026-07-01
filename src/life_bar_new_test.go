package main

import "testing"

func TestNewLifeBar(t *testing.T) {
	t.Parallel()

	lb := newLifeBar()
	if lb == nil {
		t.Fatal("newLifeBar returned nil")
	}
	if lb.red == nil || lb.front == nil || lb.value == nil || lb.red_value == nil {
		t.Fatalf("newLifeBar should allocate maps: %#v", lb)
	}
	if lb.oldlife != 1 || lb.midlife != 1 || lb.midlifeMin != 1 {
		t.Fatalf("unexpected life defaults: %#v", lb)
	}
	if !lb.mid_freeze || lb.mid_delay != 30 || lb.mid_mult != 1.0 || lb.mid_steps != 8.0 {
		t.Fatalf("unexpected life defaults: %#v", lb)
	}
}
