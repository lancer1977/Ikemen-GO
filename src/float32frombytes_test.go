package main

import (
	"math"
	"testing"
)

func TestFloat32frombytes_DecodesBitPatternsExactly(t *testing.T) {
	if got := Float32frombytes([4]byte{0x00, 0x00, 0x00, 0x00}); got != 0 {
		t.Fatalf("Float32frombytes(zero) = %v, want 0", got)
	}
	if got := Float32frombytes([4]byte{0x00, 0x00, 0x80, 0x3f}); got != 1 {
		t.Fatalf("Float32frombytes(one) = %v, want 1", got)
	}
	if got := Float32frombytes([4]byte{0x00, 0x00, 0x00, 0xc0}); got != -2 {
		t.Fatalf("Float32frombytes(neg) = %v, want -2", got)
	}
	if got := Float32frombytes([4]byte{0x00, 0x00, 0xc0, 0x7f}); !math.IsNaN(float64(got)) {
		t.Fatalf("Float32frombytes(nan) = %v, want NaN", got)
	}
}
