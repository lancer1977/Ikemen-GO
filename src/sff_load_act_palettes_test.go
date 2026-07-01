package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSffLoadActPalettes(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	dir := t.TempDir()
	actPath := filepath.Join(dir, "sample.act")
	data := make([]byte, 256*3)
	data[0], data[1], data[2] = 1, 2, 3
	if err := os.WriteFile(actPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	sys.cfg.Config.PaletteMax = 6
	sys.sel.charlist = []Char{
		{
			def: dir,
			pal_files: []string{
				"",
				"sample.act",
				"sample.act",
			},
			pal: []int32{
				1,
				2,
				9,
			},
		},
	}

	s := &Sff{}
	s.palList.init()
	if err := s.loadActPalettes(0); err != nil {
		t.Fatal(err)
	}
	if len(s.palList.palettes) != 1 {
		t.Fatalf("expected one loaded ACT palette, got %d", len(s.palList.palettes))
	}
	if got := s.palList.PalTable[[2]uint16{1, 2}]; got != 1 {
		t.Fatalf("expected palette slot 2 to map to index 1, got %d", got)
	}
	if _, ok := s.palList.PalTable[[2]uint16{1, 9}]; ok {
		t.Fatal("out-of-range palette slot should be ignored")
	}
}
