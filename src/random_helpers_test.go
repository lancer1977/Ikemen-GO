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

	// DEFECT: RandF32 reuses the integer bucket arithmetic from Rand, where the
	// "+1" counts an inclusive endpoint. For a continuous range that term simply
	// widens the span, so RandF32(min, max) actually spans [min, min+(max-min+1)]
	// and roughly half of all draws land above the requested maximum. RandF, just
	// below it in common.go, has the correct formula.
	// Tracked as lancer1977/Ikemen-GO#18.
	//
	// The assertions below pin the real behaviour: the widened bound holds, and
	// overshoot is actually produced. Both fail once the formula is fixed, which
	// is the point -- this test should then become a plain range assertion.
	Srand(5)
	const rfMin, rfMax float32 = 1.5, 2.5
	widened := rfMin + (rfMax - rfMin + 1.0)
	overshoot := 0
	const draws = 2000
	for i := 0; i < draws; i++ {
		got := RandF32(rfMin, rfMax)
		if got < rfMin || got > widened {
			t.Fatalf("RandF32() = %v, outside even the widened span [%v,%v]", got, rfMin, widened)
		}
		if got > rfMax {
			overshoot++
		}
	}
	if overshoot == 0 {
		t.Fatalf("RandF32 produced no values above the requested max of %v in %d draws; "+
			"#18 appears fixed, so this test should assert the real range instead", rfMax, draws)
	}
	if overshoot < draws/4 {
		t.Fatalf("RandF32 overshoot rate dropped to %d/%d; the span arithmetic changed, "+
			"re-check #18", overshoot, draws)
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
