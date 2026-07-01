package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadText_StripsUTF8BOMFromPlainFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bom.txt")
	if err := os.WriteFile(path, append([]byte{0xef, 0xbb, 0xbf}, []byte("hello")...), 0o644); err != nil {
		t.Fatalf("write bom file: %v", err)
	}

	got, err := LoadText(path)
	if err != nil {
		t.Fatalf("LoadText() error = %v", err)
	}
	if got != "hello" {
		t.Fatalf("LoadText() = %q, want %q", got, "hello")
	}
}

func TestLoadText_ReadsEntryFromZipArchive(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "data.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("notes/readme.txt")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := w.Write([]byte("zip text")); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	got, err := LoadText(zipPath + "/notes/readme.txt")
	if err != nil {
		t.Fatalf("LoadText(zip) error = %v", err)
	}
	if got != "zip text" {
		t.Fatalf("LoadText(zip) = %q, want %q", got, "zip text")
	}
}

func TestLoadText_StripsUTF8BOMFromZipEntry(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "data.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("notes/bom.txt")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := w.Write(append([]byte{0xef, 0xbb, 0xbf}, []byte("zip hello")...)); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	got, err := LoadText(zipPath + "/notes/bom.txt")
	if err != nil {
		t.Fatalf("LoadText(zip BOM) error = %v", err)
	}
	if got != "zip hello" {
		t.Fatalf("LoadText(zip BOM) = %q, want %q", got, "zip hello")
	}
}
