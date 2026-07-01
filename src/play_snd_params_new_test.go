package main

import "testing"

func TestNewPlaySndParams(t *testing.T) {
	t.Parallel()

	p := newPlaySndParams()
	if p == nil {
		t.Fatal("newPlaySndParams returned nil")
	}
	if p.group != -1 || p.channel != -1 {
		t.Fatalf("unexpected channel defaults: group=%d channel=%d", p.group, p.channel)
	}
	if p.volume != 100 || p.freqMul != 1.0 || p.localScale != 1.0 {
		t.Fatalf("unexpected numeric defaults: %#v", p)
	}
	if p.lowPriority || p.stopOnGetHit || p.stopOnChangeState || p.log {
		t.Fatalf("unexpected boolean defaults: %#v", p)
	}
}
