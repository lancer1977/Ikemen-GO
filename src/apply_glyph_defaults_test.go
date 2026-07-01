package main

import "testing"

func TestApplyGlyphDefaultsFromMovelist(t *testing.T) {
	t.Parallel()

	m := &Motif{
		PauseMenu: map[string]*MenuInfoProperties{
			"pause_menu": {
				Movelist: struct {
					Pos   [2]float32 `ini:"pos"`
					Title struct {
						TextProperties
						Uppercase bool `ini:"uppercase"`
					} `ini:"title"`
					Text struct {
						TextProperties
						Spacing [2]float32 `ini:"spacing"`
					} `ini:"text"`
					Glyphs MovelistGlyphsProperties `ini:"glyphs"`
					Window struct {
						Margins struct {
							Y [2]float32 `ini:"y"`
						} `ini:"margins"`
						VisibleItems int32 `ini:"visibleitems"`
						Width        int32 `ini:"width"`
					} `ini:"window"`
				}{
					Glyphs: MovelistGlyphsProperties{
						Offset:     [2]float32{1, 2},
						Scale:      [2]float32{3, 4},
						Layerno:    5,
						Localcoord: [2]int32{6, 7},
						Spacing:    [2]float32{8, 9},
					},
				},
			},
		},
		Glyphs: map[string]*GlyphProperties{
			"a": {Offset: [2]float32{0, 0}},
			"b": nil,
		},
	}

	m.applyGlyphDefaultsFromMovelist()

	if got := m.Glyphs["a"]; got == nil || got.Offset != [2]float32{1, 2} || got.Scale != [2]float32{3, 4} || got.Layerno != 5 || got.Localcoord != [2]int32{6, 7} {
		t.Fatalf("glyph defaults not applied: %#v", got)
	}
	if m.Glyphs["b"] != nil {
		t.Fatalf("nil glyph entry should remain nil placeholder, got %#v", m.Glyphs["b"])
	}
}
