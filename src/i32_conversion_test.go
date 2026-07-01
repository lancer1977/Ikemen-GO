package main

import (
	"math"
	"testing"
)

func TestI32ToI16_ClampsToInt16Bounds(t *testing.T) {
	if got := I32ToI16(int32(math.MaxInt16) + 1); got != math.MaxInt16 {
		t.Fatalf("I32ToI16(high) = %v, want %v", got, math.MaxInt16)
	}
	if got := I32ToI16(int32(math.MinInt16) - 1); got != math.MinInt16 {
		t.Fatalf("I32ToI16(low) = %v, want %v", got, math.MinInt16)
	}
	if got := I32ToI16(1234); got != 1234 {
		t.Fatalf("I32ToI16(mid) = %v, want %v", got, 1234)
	}
}

func TestI32ToU16_ClampsToUint16Bounds(t *testing.T) {
	if got := I32ToU16(-1); got != 0 {
		t.Fatalf("I32ToU16(low) = %v, want 0", got)
	}
	if got := I32ToU16(int32(math.MaxUint16) + 1); got != math.MaxUint16 {
		t.Fatalf("I32ToU16(high) = %v, want %v", got, math.MaxUint16)
	}
	if got := I32ToU16(1234); got != 1234 {
		t.Fatalf("I32ToU16(mid) = %v, want %v", got, 1234)
	}
}
