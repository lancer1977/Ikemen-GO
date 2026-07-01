package main

import (
	"math"
	"testing"
)

func TestF64toI32_ClampsToInt32BoundsLegacy(t *testing.T) {
	if got := F64toI32(float64(math.MaxInt32) + 1000); got != math.MaxInt32 {
		t.Fatalf("F64toI32(high) = %v, want %v", got, math.MaxInt32)
	}
	if got := F64toI32(float64(math.MinInt32) - 1000); got != math.MinInt32 {
		t.Fatalf("F64toI32(low) = %v, want %v", got, math.MinInt32)
	}
	if got := F64toI32(1234.75); got != 1234 {
		t.Fatalf("F64toI32(mid) = %v, want %v", got, 1234)
	}
}
