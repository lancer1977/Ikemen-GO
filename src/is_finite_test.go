package main

import (
	"math"
	"testing"
)

func TestIsFinite_RejectsNaNAndInfinities(t *testing.T) {
	if !IsFinite(0) || !IsFinite(123.5) {
		t.Fatal("IsFinite should accept finite values")
	}
	if IsFinite(float32(math.Inf(1))) || IsFinite(float32(math.Inf(-1))) {
		t.Fatal("IsFinite should reject infinities")
	}
	if IsFinite(float32(math.NaN())) {
		t.Fatal("IsFinite should reject NaN")
	}
}
