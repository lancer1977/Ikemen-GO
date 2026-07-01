package main

import "testing"

func TestNewNormalizer(t *testing.T) {
	t.Parallel()

	n := NewNormalizer(nil)
	if n == nil {
		t.Fatal("NewNormalizer returned nil")
	}
	if n.streamer != nil {
		t.Fatalf("streamer = %#v, want nil", n.streamer)
	}
	if n.mul != 4 {
		t.Fatalf("mul = %v, want 4", n.mul)
	}
	if n.l == nil || n.r == nil {
		t.Fatal("expected both channel normalizers to be initialized")
	}
	if *n.l != (NormalizerLR{edge: 1, gain: 1, average: 1 / 32.0}) {
		t.Fatalf("unexpected left normalizer defaults: %#v", *n.l)
	}
	if *n.r != (NormalizerLR{edge: 1, gain: 1, average: 1 / 32.0}) {
		t.Fatalf("unexpected right normalizer defaults: %#v", *n.r)
	}
}
