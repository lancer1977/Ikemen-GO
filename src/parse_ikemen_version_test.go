package main

import "testing"

func TestParseIkemenVersion_ArrayParserDoesNotStripNonNumericSuffixes(t *testing.T) {
	// DEFECT: The array parser stops when it encounters a non-numeric character in a segment,
	// so "1.02.3-beta+build.7" yields [1, 2, 0] instead of [1, 2, 3].
	// The float parser handles suffixes correctly via regex stripping.
	// Tracked as lancer1977/Ikemen-GO#14.
	ver, verF := ParseIkemenVersion("1.02.3-beta+build.7")
	if ver != [3]uint16{1, 2, 0} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 2, 0})
	}
	if verF != 1.0237 {
		t.Fatalf("verF = %v, want %v", verF, 1.0237)
	}
}

func TestParseIkemenVersion_StopsIntegerParseOnInvalidMiddleSegment(t *testing.T) {
	ver, verF := ParseIkemenVersion("1.bad.3")
	if ver != [3]uint16{1, 0, 0} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 0, 0})
	}
	if verF != 1.3 {
		t.Fatalf("verF = %v, want %v", verF, 1.3)
	}
}
