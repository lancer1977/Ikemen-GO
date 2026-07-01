package main

import (
	"arena"
	"testing"
)

func TestCloneTextSprite(t *testing.T) {
	t.Parallel()

	a := arena.NewArena()
	defer a.Free()

	ts := NewTextSprite()
	ts.text = "hello"
	ts.params = []interface{}{"x", 2}
	ts.palfx.setColor(4, 5, 6)

	cp := cloneTextSprite(a, ts)
	if cp == nil {
		t.Fatal("cloneTextSprite returned nil")
	}
	if cp == ts {
		t.Fatal("expected clone to allocate a new object")
	}
	if cp.text != ts.text {
		t.Fatalf("text = %q, want %q", cp.text, ts.text)
	}
	if len(cp.params) != len(ts.params) || cp.params[0] != ts.params[0] || cp.params[1] != ts.params[1] {
		t.Fatalf("params not cloned correctly: %#v", cp.params)
	}
	if cp.palfx == nil || cp.palfx == ts.palfx {
		t.Fatal("expected palfx to be cloned independently")
	}
	cp.params[0] = "changed"
	if ts.params[0] != "x" {
		t.Fatalf("original params mutated through clone: %#v", ts.params)
	}
}
