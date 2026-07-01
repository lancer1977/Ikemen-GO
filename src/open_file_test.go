package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFile_ReadsPlainFileAndSupportsSeek(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plain.txt")
	if err := os.WriteFile(path, []byte("plain text"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	rc, err := OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}
	defer rc.Close()

	buf, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read plain file: %v", err)
	}
	if string(buf) != "plain text" {
		t.Fatalf("OpenFile() = %q, want %q", string(buf), "plain text")
	}
	if _, err := rc.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("seek plain file: %v", err)
	}
}

func TestOpenFile_ReadsZipEntryAndClosesArchive(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "assets.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("chars/ryu.def")
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

	rc, err := OpenFile(zipPath + "/chars/ryu.def")
	if err != nil {
		t.Fatalf("OpenFile(zip) error = %v", err)
	}

	buf, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read zip entry: %v", err)
	}
	if string(buf) != "zip text" {
		t.Fatalf("OpenFile(zip) = %q, want %q", string(buf), "zip text")
	}
	if err := rc.Close(); err != nil {
		t.Fatalf("close zip-backed reader: %v", err)
	}
}
