package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestSearchFile_ResolvesFilesInsideZipRoot(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "chars.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("chars/ryu.def")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := w.Write([]byte("def")); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	got := SearchFile("ryu.def", []string{zipPath + "/chars"})
	want := filepath.ToSlash(zipPath + "/chars/ryu.def")
	if got != want {
		t.Fatalf("SearchFile(zip root) = %q, want %q", got, want)
	}
}
