package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

// DEFECT: SearchFile's rootBases helper handles an extension-less lookup root
// differently inside a zip archive than on the plain filesystem. The plain
// branch adds the root itself as a second base, so a root naming a directory
// resolves files directly within it. The zip branch returns only the parent of
// the in-zip path, so the subdirectory is dropped and the file is never found.
// SearchFile then falls through and returns its input verbatim, which is its
// not-found signal. Tracked as lancer1977/Ikemen-GO#17.
//
// The two subtests below are deliberately the same layout, once in an archive
// and once on disk, so the asymmetry is what fails if either side changes.
func TestSearchFile_ZipDirectoryRootDropsSubdirectory(t *testing.T) {
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

	t.Run("directory_root_inside_zip_is_not_resolved", func(t *testing.T) {
		// Only chars.zip/ryu.def is tried, never chars.zip/chars/ryu.def.
		got := SearchFile("ryu.def", []string{zipPath + "/chars"})
		if got != "ryu.def" {
			t.Fatalf("SearchFile(zip dir root) = %q; want the unresolved input %q. "+
				"If this now returns the real path, #17 is fixed and this test "+
				"should assert the resolved path instead.", got, "ryu.def")
		}
	})

	t.Run("file_root_inside_zip_is_resolved", func(t *testing.T) {
		// A root naming a file inside the archive works, because its parent is
		// exactly the directory the lookup needs.
		got := SearchFile("ryu.def", []string{zipPath + "/chars/kfm.def"})
		want := filepath.ToSlash(zipPath + "/chars/ryu.def")
		if got != want {
			t.Fatalf("SearchFile(zip file root) = %q, want %q", got, want)
		}
	})

	t.Run("equivalent_plain_directory_root_is_resolved", func(t *testing.T) {
		// The same shape on the plain filesystem resolves, which is what makes
		// the zip result above a defect rather than the intended contract.
		plainRoot := filepath.Join(dir, "chars")
		if err := os.MkdirAll(plainRoot, 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		defPath := filepath.Join(plainRoot, "ryu.def")
		if err := os.WriteFile(defPath, []byte("def"), 0o644); err != nil {
			t.Fatalf("write def: %v", err)
		}
		got := SearchFile("ryu.def", []string{plainRoot})
		if got != filepath.ToSlash(defPath) {
			t.Fatalf("SearchFile(plain dir root) = %q, want %q", got, filepath.ToSlash(defPath))
		}
	})
}
