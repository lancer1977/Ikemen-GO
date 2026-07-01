package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOverlaySprite(t *testing.T) {
	if _, err := loadOverlaySprite("   "); err == nil {
		t.Fatal("expected blank asset path to fail")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "overlay.png")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create png: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.SetRGBA(0, 0, color.RGBA{R: 255, A: 255})
	img.SetRGBA(1, 0, color.RGBA{G: 255, A: 255})
	if err := png.Encode(f, img); err != nil {
		_ = f.Close()
		t.Fatalf("encode png: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close png: %v", err)
	}

	sprite, err := loadOverlaySprite(path)
	if err != nil {
		t.Fatalf("loadOverlaySprite returned error: %v", err)
	}
	if sprite == nil {
		t.Fatal("expected sprite to be created")
	}
	if sprite.Size != [2]uint16{2, 1} {
		t.Fatalf("unexpected sprite size: %#v", sprite.Size)
	}
}
