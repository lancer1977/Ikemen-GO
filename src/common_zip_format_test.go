package main

import "testing"

func TestIsZipPath_RecognizesArchiveAndInternalEntryForms(t *testing.T) {
	isZip, zipFile, inside := IsZipPath(`C:\data\packs\chars.zip\ryu\ryu.def`)
	if !isZip || zipFile != "C:/data/packs/chars.zip" || inside != "ryu/ryu.def" {
		t.Fatalf("IsZipPath(internal) = %v, %q, %q", isZip, zipFile, inside)
	}

	isZip, zipFile, inside = IsZipPath("stages/common.zip")
	if !isZip || zipFile != "stages/common.zip" || inside != "" {
		t.Fatalf("IsZipPath(archive) = %v, %q, %q", isZip, zipFile, inside)
	}

	isZip, zipFile, inside = IsZipPath("stages/common.def")
	if isZip || zipFile != "" || inside != "" {
		t.Fatalf("IsZipPath(non-zip) = %v, %q, %q", isZip, zipFile, inside)
	}
}

func TestOldSprintf_NormalizesLegacyIntegerVerbs(t *testing.T) {
	got := OldSprintf("x=%i y=%u z=%d", 7, 8, 9)
	if got != "x=7 y=8 z=9" {
		t.Fatalf("OldSprintf() = %q, want %q", got, "x=7 y=8 z=9")
	}
}

func TestOldSprintf_TruncatesExtraArguments(t *testing.T) {
	got := OldSprintf("value=%d", 42, 99)
	if got != "value=42" {
		t.Fatalf("OldSprintf() = %q, want %q", got, "value=42")
	}
}
