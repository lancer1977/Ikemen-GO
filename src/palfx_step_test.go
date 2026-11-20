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
		pfx.invertblend = -1
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
		// invertblend above -2 is copied through untouched.
		if pfx.eInvertblend != -1 || !pfx.eAllowNeg {
			t.Fatalf("step should copy invert/neg settings, got %#v", pfx)
		}
	})

	// step() clamps an invertblend of -2 or below to 3, but only while invertall
	// is set. Both halves of that guard need exercising or the clamp branch goes
	// uncovered.
	t.Run("invertblend_clamped_when_inverting_all", func(t *testing.T) {
		pfx := newPalFX()
		pfx.time = 2
		pfx.interpolate = false
		pfx.invertblend = -3
		pfx.invertall = true
		pfx.step()
		if pfx.eInvertblend != 3 {
			t.Fatalf("invertblend -3 with invertall should clamp to 3, got %d", pfx.eInvertblend)
		}
	})

	t.Run("invertblend_not_clamped_without_invertall", func(t *testing.T) {
		pfx := newPalFX()
		pfx.time = 2
		pfx.interpolate = false
		pfx.invertblend = -3
		pfx.invertall = false
		pfx.step()
		if pfx.eInvertblend != -3 {
			t.Fatalf("invertblend -3 without invertall should pass through, got %d", pfx.eInvertblend)
		}
	})
}
