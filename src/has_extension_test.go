package main

import "testing"

func TestHasExtension_MatchesLowercasedFileExtensionAgain(t *testing.T) {
	if !HasExtension("select.DEF", `(?i)^\.def$`) {
		t.Fatal("expected .def extension to match case-insensitively")
	}
	if HasExtension("select.txt", `(?i)^\.def$`) {
		t.Fatal("expected .txt extension to not match .def")
	}
	if HasExtension("select", `(?i)^\.def$`) {
		t.Fatal("expected file with no extension to not match .def")
	}
}
