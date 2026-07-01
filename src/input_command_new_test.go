package main

import "testing"

func TestNewCommand(t *testing.T) {
	t.Parallel()

	c := newCommand()
	if c == nil {
		t.Fatal("newCommand returned nil")
	}
	if c.maxtime != 1 || c.maxbuftime != 1 {
		t.Fatalf("unexpected defaults: maxtime=%d maxbuftime=%d", c.maxtime, c.maxbuftime)
	}
	if c.curtime != 0 || c.curbuftime != 0 {
		t.Fatalf("unexpected current timers: curtime=%d curbuftime=%d", c.curtime, c.curbuftime)
	}
	if c.steps != nil || c.completed != nil || c.stepTimers != nil || c.loopOrder != nil {
		t.Fatalf("newCommand should start with nil slices: %#v", c)
	}
}
