package main

import "testing"

func TestHasExtension_ReturnsFalseForMissingExtensionAndMatchesUppercaseFileName(t *testing.T) {
	if HasExtension("select", `(?i)^\.def$`) {
		t.Fatal("expected file without extension to not match")
	}
	if !HasExtension("SELECT.DEF", `(?i)^\.def$`) {
		t.Fatal("expected uppercase file name to match case-insensitively")
	}
}
