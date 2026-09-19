package main

import "testing"

func TestRandomHelpers_StayWithinConfiguredRangesAgain(t *testing.T) {
	// Test integer ranges with deterministic seed
	Srand(1)
	for i := 0; i < 20; i++ {
		if got := Rand(3, 7); got < 3 || got > 7 {
			t.Fatalf("Rand() = %v, want within [3,7]", got)
		}
		if got := RandI(10, 15); got < 10 || got > 15 {
			t.Fatalf("RandI() = %v, want within [10,15]", got)
		}
		if got := RandF(4.0, 6.0); got < 4.0 || got > 6.0 {
			t.Fatalf("RandF() = %v, want within [4,6]", got)
		}
	}

	// Test RandF32 with deterministic seed - statistically verify it stays within range
	// Fixed in #18: removed the incorrect "+1.0" term that was copied from the integer
	// Rand function. For continuous float ranges, the divisor should scale Random()
	// across the span (max-min), not (max-min+1).
	Srand(5)
	const rfMin, rfMax float32 = 1.5, 2.5
	const draws = 20000
	var aboveMax, belowMin int
	for i := 0; i < draws; i++ {
		got := RandF32(rfMin, rfMax)
		if got < rfMin {
			belowMin++
		}
		if got > rfMax {
			aboveMax++
		}
	}
	if belowMin > 0 {
		t.Fatalf("RandF32 produced %d/%d values below the requested min of %v", belowMin, draws, rfMin)
	}
	if aboveMax > 0 {
		t.Fatalf("RandF32 produced %d/%d values above the requested max of %v", aboveMax, draws, rfMax)
	}
}

func TestRandomAndSrand_ProduceDeterministicSequenceAgain(t *testing.T) {
	Srand(123)
	a := []int32{Random(), Random(), Random()}
	Srand(123)
	b := []int32{Random(), Random(), Random()}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("Random sequence mismatch at %d: %v vs %v", i, a, b)
		}
	}
}
