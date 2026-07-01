package main

import "testing"

func TestPalFXStep(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		pfx := newPalFX()
		pfx.time = 0
		pfx.step()
		if pfx.enable {
			t.Fatal("step should keep disabled FX off")
		}
	})

	t.Run("enabled_without_interpolation", func(t *testing.T) {
		pfx := newPalFX()
		pfx.time = 2
		pfx.interpolate = false
		pfx.invertblend = -3
		pfx.invertall = true
		pfx.allowNeg = true
		pfx.mul = [3]int32{10, 20, 30}
		pfx.add = [3]int32{1, 2, 3}
		pfx.color = 0.75
		pfx.hue = 4
		pfx.step()
		if !pfx.enable || pfx.eInterpolate {
			t.Fatalf("step should enable without interpolation, got %#v", pfx)
		}
		if pfx.eMul != pfx.mul || pfx.eAdd != pfx.add || pfx.eColor != pfx.color || pfx.eHue != pfx.hue {
			t.Fatalf("step should copy base values, got %#v", pfx)
		}
		if pfx.eInvertblend != -3 || !pfx.eAllowNeg {
			t.Fatalf("step should copy invert/neg settings, got %#v", pfx)
		}
	})
}
