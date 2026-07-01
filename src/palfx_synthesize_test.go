package main

import "testing"

func TestPalFXSynthesize(t *testing.T) {
	t.Run("add_mode", func(t *testing.T) {
		pf := &PalFX{
			PalFXDef:     PalFXDef{},
			eAdd:         [3]int32{10, 20, 30},
			eMul:         [3]int32{256, 128, 64},
			eHue:         2,
			eColor:       0.5,
			eInvertall:   true,
			eInvertblend: 1,
		}
		pfx := &PalFX{
			PalFXDef:   PalFXDef{invertall: false},
			eAdd:       [3]int32{1, 2, 3},
			eMul:       [3]int32{128, 64, 32},
			eHue:       4,
			eColor:     0.25,
			eInvertall: false,
		}
		pf.synthesize(pfx, TT_add, [2]int32{0, 0})
		if pf.eAdd != [3]int32{11, 22, 33} {
			t.Fatalf("add mode eAdd = %#v", pf.eAdd)
		}
		if pf.eMul != [3]int32{128, 32, 8} {
			t.Fatalf("add mode eMul = %#v", pf.eMul)
		}
		if pf.eHue != 6 || pf.eColor != 0.125 {
			t.Fatalf("add mode hue/color = %v/%v", pf.eHue, pf.eColor)
		}
		if !pf.eInvertall {
			t.Fatalf("add mode invertall should toggle true")
		}
	})

	t.Run("sub_mode", func(t *testing.T) {
		pf := &PalFX{
			eAdd:       [3]int32{100, 100, 100},
			eMul:       [3]int32{200, 150, 50},
			eHue:       1,
			eColor:     0.8,
			eInvertall: false,
		}
		pfx := &PalFX{
			eAdd:        [3]int32{10, -20, 300},
			eMul:        [3]int32{64, 32, 16},
			eHue:        3,
			eColor:      0.5,
			eInvertall:  true,
			invertall:   true,
			invertblend: 0,
		}
		pf.synthesize(pfx, TT_sub, [2]int32{0, 0})
		if pf.eAdd != [3]int32{90, 80, 0} {
			t.Fatalf("sub mode eAdd = %#v", pf.eAdd)
		}
		if pf.eMul != [3]int32{190, 170, 34} {
			t.Fatalf("sub mode eMul = %#v", pf.eMul)
		}
		if pf.eHue != 4 || pf.eColor != 0.4 {
			t.Fatalf("sub mode hue/color = %v/%v", pf.eHue, pf.eColor)
		}
	})

	t.Run("invertblend_remap", func(t *testing.T) {
		base := &PalFX{
			eAdd:        [3]int32{0, 0, 0},
			eMul:        [3]int32{256, 256, 256},
			eColor:      1,
			eHue:        0,
			eInvertall:  true,
			invertall:   true,
			invertblend: 1,
		}
		other := &PalFX{
			eAdd:         [3]int32{0, 0, 0},
			eMul:         [3]int32{256, 256, 256},
			eInvertall:   true,
			eInvertblend: 2,
			invertall:    true,
			invertblend:  1,
		}
		base.synthesize(other, TT_add, [2]int32{1, 1})
		if base.eInvertblend != 2 || !base.eInvertall {
			t.Fatalf("invertblend remap failed: %#v", base)
		}
	})
}
