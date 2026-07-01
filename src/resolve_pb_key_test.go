package main

import "testing"

func TestResolvePBKey(t *testing.T) {
	t.Parallel()

	keys := map[int32]string{
		0:  "low",
		10: "mid",
		20: "high",
	}

	if got := resolvePBKey(keys, 0, 30); got != 0 {
		t.Fatalf("pbval=0 -> %d, want 0", got)
	}
	if got := resolvePBKey(keys, 15, 30); got != 10 {
		t.Fatalf("pbval=15 -> %d, want 10", got)
	}
	if got := resolvePBKey(keys, 25, 30); got != 20 {
		t.Fatalf("pbval=25 -> %d, want 20", got)
	}

	keys[-1] = "max"
	if got := resolvePBKey(keys, 30, 30); got != -1 {
		t.Fatalf("pbval=max -> %d, want -1", got)
	}
	if got := resolvePBKey(keys, 29, 30); got != 20 {
		t.Fatalf("pbval below max -> %d, want 20", got)
	}
}
