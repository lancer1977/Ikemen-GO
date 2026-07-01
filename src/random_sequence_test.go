package main

import "testing"

func TestRandomAndSrand_ProduceDeterministicSequence(t *testing.T) {
	prevSeed := sys.randseed
	t.Cleanup(func() { sys.randseed = prevSeed })

	Srand(1)
	if got := Random(); got != 16807 {
		t.Fatalf("unexpected first random value: %d", got)
	}
	if got := sys.randseed; got != 16807 {
		t.Fatalf("unexpected seed after first random: %d", got)
	}
	if got := Random(); got != 282475249 {
		t.Fatalf("unexpected second random value: %d", got)
	}
}
