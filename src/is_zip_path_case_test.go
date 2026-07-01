package main

import "testing"

func TestIsZipPath_HandlesUppercaseExtensionAndNestedEntry(t *testing.T) {
	isZip, zipFile, inside := IsZipPath(`C:\data\Packs\ARCHIVE.ZIP\Chars\Ryu.DEF`)
	if !isZip || zipFile != "C:/data/Packs/ARCHIVE.ZIP" || inside != "Chars/Ryu.DEF" {
		t.Fatalf("IsZipPath(uppercase internal) = %v, %q, %q", isZip, zipFile, inside)
	}

	isZip, zipFile, inside = IsZipPath("STAGES/COMMON.ZIP")
	if !isZip || zipFile != "STAGES/COMMON.ZIP" || inside != "" {
		t.Fatalf("IsZipPath(uppercase archive) = %v, %q, %q", isZip, zipFile, inside)
	}
}
