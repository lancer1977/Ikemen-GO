package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenFile_ReturnsErrorForMissingZipEntry(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "assets.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(f)
	if _, err := zw.Create("chars/ryu.def"); err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	rc, err := OpenFile(zipPath + "/chars/ken.def")
	if err == nil {
		rc.Close()
		t.Fatal("expected error for missing zip entry")
	}
	if !strings.Contains(err.Error(), "not found in zip archive") {
		t.Fatalf("OpenFile missing-entry error = %v, want not found in zip archive", err)
	}
}
