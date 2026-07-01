package main

import (
	"math"
	"testing"
)

func TestPalFXSineHelpers(t *testing.T) {
	pfx := &PalFX{
		sintime:   [4]int32{1, 1, 1, 1},
		cycletime: [4]int32{2, 2, 2, 2},
		sinadd:    [3]int32{10, 20, 30},
		sinmul:    [3]int32{40, 50, 60},
		sincolor:  256,
		sinhue:    512,
	}

	color := [3]int32{1, 2, 3}
	pfx.sinAdd(&color)
	if color != [3]int32{-6, -12, -18} {
		t.Fatalf("sinAdd = %#v", color)
	}

	color = [3]int32{1, 2, 3}
	pfx.sinMul(&color)
	if color != [3]int32{-27, -33, -39} {
		t.Fatalf("sinMul = %#v", color)
	}

	c := float32(1)
	pfx.sinColor(&c)
	wantColor := float32(1 + math.Sin(5*math.Pi/4))
	if math.Abs(float64(c-wantColor)) > 1e-6 {
		t.Fatalf("sinColor = %v, want %v", c, wantColor)
	}

	h := float32(1)
	pfx.sinHueshift(&h)
	wantHue := float32(1 + math.Sin(5*math.Pi/4))
	if math.Abs(float64(h-wantHue)) > 1e-6 {
		t.Fatalf("sinHueshift = %v, want %v", h, wantHue)
	}
}
