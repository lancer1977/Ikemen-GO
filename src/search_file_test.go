package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchFile_TrimsCommentsAndUsesDefaultDirectory(t *testing.T) {
	dir := t.TempDir()
	defaultDir := filepath.Join(dir, "font")
	if err := os.MkdirAll(defaultDir, 0o755); err != nil {
		t.Fatalf("mkdir default dir: %v", err)
	}
	wantPath := filepath.Join(defaultDir, "select.fnt")
	if err := os.WriteFile(wantPath, []byte("font"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got := SearchFile(`  select.fnt ; comment`, []string{dir}, "font")
	if got != filepath.ToSlash(wantPath) {
		t.Fatalf("SearchFile() = %q, want %q", got, filepath.ToSlash(wantPath))
	}
}

// Quote stripping lives in LoadFile, not SearchFile: LoadFile peels a matching
// leading and trailing quote off the path before it delegates. SearchFile trims
// whitespace and strips a trailing comment, but treats a quote as part of the
// filename, so a quoted name does not resolve and comes back unchanged. Pinning
// both halves keeps the split of responsibility from drifting silently.
func TestSearchFile_DoesNotStripQuotesThatLoadFileHandles(t *testing.T) {
	dir := t.TempDir()
	defaultDir := filepath.Join(dir, "font")
	if err := os.MkdirAll(defaultDir, 0o755); err != nil {
		t.Fatalf("mkdir default dir: %v", err)
	}
	wantPath := filepath.Join(defaultDir, "select.fnt")
	if err := os.WriteFile(wantPath, []byte("font"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	quoted := `  "select.fnt" ; comment`
	if got := SearchFile(quoted, []string{dir}, "font"); got != `"select.fnt"` {
		t.Fatalf("SearchFile(quoted) = %q; want the unresolved input %q, since "+
			"quote stripping is LoadFile's job", got, `"select.fnt"`)
	}

	// The same name does resolve once LoadFile has peeled the quotes. LoadFile
	// only strips a quote pair that wraps the whole trimmed string, so the
	// trailing comment has to be gone before the quotes are its outermost
	// characters.
	loaded := `  "select.fnt"  `
	var seen string
	if err := LoadFile(&loaded, []string{dir}, "font", func(p string) error {
		seen = p
		return nil
	}); err != nil {
		t.Fatalf("LoadFile(quoted) error: %v", err)
	}
	if seen != filepath.ToSlash(wantPath) {
		t.Fatalf("LoadFile(quoted) resolved to %q, want %q", seen, filepath.ToSlash(wantPath))
	}
}

func TestSearchFile_ReturnsAbsoluteAndRootFallbackMatches(t *testing.T) {
	dir := t.TempDir()
	absPath := filepath.Join(dir, "data.def")
	if err := os.WriteFile(absPath, []byte("def"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if got := SearchFile(absPath, nil); got != filepath.ToSlash(absPath) {
		t.Fatalf("SearchFile(abs) = %q, want %q", got, filepath.ToSlash(absPath))
	}

	if got := SearchFile(filepath.Base(absPath), []string{dir}); got != filepath.ToSlash(absPath) {
		t.Fatalf("SearchFile(root fallback) = %q, want %q", got, filepath.ToSlash(absPath))
	}
}
