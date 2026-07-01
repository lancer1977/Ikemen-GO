package main

import (
	"path/filepath"
	"testing"
)

func TestSplitPath_NormalizesSeparatorsAndSplitsAtLastSlash(t *testing.T) {
	dir, file := SplitPath(`C:\games\ikemen\data\chars\ryu.def`)
	if dir != "C:/games/ikemen/data/chars/" {
		t.Fatalf("dir = %q, want %q", dir, "C:/games/ikemen/data/chars/")
	}
	if file != "ryu.def" {
		t.Fatalf("file = %q, want %q", file, "ryu.def")
	}
}

func TestStripComment_RemovesInlineCommentAndTrimsSpace(t *testing.T) {
	if got := StripComment("value ; trailing comment  "); got != "value" {
		t.Fatalf("StripComment() = %q, want %q", got, "value")
	}
}

func TestSectionName_ParsesBracketedSectionAndIgnoresComment(t *testing.T) {
	name, body := SectionName("[State  100]; comment")
	if name != "state " {
		t.Fatalf("name = %q, want %q", name, "state ")
	}
	if body != " 100" {
		t.Fatalf("body = %q, want %q", body, " 100")
	}
}

func TestHasExtension_MatchesLowercasedFileExtension(t *testing.T) {
	if !HasExtension("select.DEF", `(?i)^\.def$`) {
		t.Fatal("expected .def extension to match case-insensitively")
	}
	if HasExtension("select.txt", `(?i)^\.def$`) {
		t.Fatal("expected .txt extension to not match .def")
	}
}

func TestLowercaseNoExtension_StripsExtensionAndLowercasesBaseName(t *testing.T) {
	got := LowercaseNoExtension(filepath.Join("data", "Chars", "RYU.CNS"))
	if got != "ryu" {
		t.Fatalf("LowercaseNoExtension() = %q, want %q", got, "ryu")
	}
}

func TestParseIkemenVersion_ParsesPreciseAndFloatForms(t *testing.T) {
	ver, verF := ParseIkemenVersion("1.10.3-beta")
	if ver != [3]uint16{1, 10, 3} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 10, 3})
	}
	if verF != 1.103 {
		t.Fatalf("verF = %v, want %v", verF, 1.103)
	}
}

func TestParseMugenVersion_MapsKnownVersionsAndRejectsInvalid(t *testing.T) {
	ver, verF := ParseMugenVersion("1.1")
	if ver != [2]uint16{1, 1} || verF != 1.1 {
		t.Fatalf("ParseMugenVersion(1.1) = %#v, %v, want [1 1], 1.1", ver, verF)
	}

	ver, verF = ParseMugenVersion("2.0")
	if ver != [2]uint16{2, 0} || verF != 0.5 {
		t.Fatalf("ParseMugenVersion(2.0) = %#v, %v, want [2 0], 0.5", ver, verF)
	}

	ver, verF = ParseMugenVersion("bad.version")
	if ver != [2]uint16{} || verF != 0 {
		t.Fatalf("ParseMugenVersion(invalid) = %#v, %v, want zero values", ver, verF)
	}
}
