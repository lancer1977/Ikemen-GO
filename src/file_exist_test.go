package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestFileExist_ReturnsPlainFilePathAndBlankForMissing(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "select.def")
	if err := os.WriteFile(filePath, []byte("root"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if got := FileExist(filePath); got != filepath.ToSlash(filePath) {
		t.Fatalf("FileExist(plain) = %q, want %q", got, filepath.ToSlash(filePath))
	}

	if got := FileExist(filepath.Join(dir, "missing.def")); got != "" {
		t.Fatalf("FileExist(missing) = %q, want empty", got)
	}
}

func TestFileExist_FindsFileInsideZipArchive(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "chars.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("ryu/ryu.def")
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

	got := FileExist(zipPath + "/ryu/ryu.def")
	if got != filepath.ToSlash(zipPath+"/ryu/ryu.def") {
		t.Fatalf("FileExist(zip entry) = %q, want %q", got, filepath.ToSlash(zipPath+"/ryu/ryu.def"))
	}
}
