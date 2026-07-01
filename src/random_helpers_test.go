package main

import "testing"

func TestRandomHelpers_StayWithinConfiguredRangesAgain(t *testing.T) {
	Srand(1)
	for i := 0; i < 20; i++ {
		if got := Rand(3, 7); got < 3 || got > 7 {
			t.Fatalf("Rand() = %v, want within [3,7]", got)
		}
		if got := RandF32(1.5, 2.5); got < 1.5 || got > 2.5 {
			t.Fatalf("RandF32() = %v, want within [1.5,2.5]", got)
		}
		if got := RandI(10, 15); got < 10 || got > 15 {
			t.Fatalf("RandI() = %v, want within [10,15]", got)
		}
		if got := RandF(4.0, 6.0); got < 4.0 || got > 6.0 {
			t.Fatalf("RandF() = %v, want within [4,6]", got)
		}
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
