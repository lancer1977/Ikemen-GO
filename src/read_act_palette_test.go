package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadActPalette_ReversesEntriesAndMarksIndexZeroTransparent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.act")
	data := make([]byte, 256*3)
	data[0], data[1], data[2] = 1, 2, 3
	data[3], data[4], data[5] = 4, 5, 6
	data[765], data[766], data[767] = 7, 8, 9
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write act file: %v", err)
	}

	pal, err := readActPalette(path)
	if err != nil {
		t.Fatalf("readActPalette returned error: %v", err)
	}
	if len(pal) != 256 {
		t.Fatalf("readActPalette length = %d, want 256", len(pal))
	}

	if got := pal[255]; got != 0xFF030201 {
		t.Fatalf("readActPalette highest entry = %#x, want 0xFF030201", got)
	}
	if got := pal[254]; got != 0xFF060504 {
		t.Fatalf("readActPalette next entry = %#x, want 0xFF060504", got)
	}
	if got := pal[0]; got != 0x00090807 {
		t.Fatalf("readActPalette index zero = %#x, want 0x00090807", got)
	}
}
