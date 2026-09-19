package main

import "testing"

func TestIsZipPath_HandlesUppercaseExtensionAndNestedEntry(t *testing.T) {
	// Forward slashes only: filepath.ToSlash is a no-op for backslashes on
	// non-Windows hosts, so a backslash path would not be recognised here.
	isZip, zipFile, inside := IsZipPath("C:/data/Packs/ARCHIVE.ZIP/Chars/Ryu.DEF")
	if !isZip || zipFile != "C:/data/Packs/ARCHIVE.ZIP" || inside != "Chars/Ryu.DEF" {
		t.Fatalf("IsZipPath(uppercase internal) = %v, %q, %q", isZip, zipFile, inside)
	}

	isZip, zipFile, inside = IsZipPath("STAGES/COMMON.ZIP")
	if !isZip || zipFile != "STAGES/COMMON.ZIP" || inside != "" {
		t.Fatalf("IsZipPath(uppercase archive) = %v, %q, %q", isZip, zipFile, inside)
	}
}
