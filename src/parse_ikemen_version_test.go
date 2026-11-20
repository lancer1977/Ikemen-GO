package main

import "testing"

func TestParseIkemenVersion_ArrayParserDoesNotStripNonNumericSuffixes(t *testing.T) {
	// FIXED: The array parser now strips non-numeric characters from each segment before parsing,
	// so "1.02.3-beta+build.7" now correctly yields [1, 2, 3] instead of [1, 2, 0].
	// Both the array and float paths now handle suffixes correctly.
	// Tracked as lancer1977/Ikemen-GO#14.
	ver, verF := ParseIkemenVersion("1.02.3-beta+build.7")
	if ver != [3]uint16{1, 2, 3} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 2, 3})
	}
	if verF != 1.0237 {
		t.Fatalf("verF = %v, want %v", verF, 1.0237)
	}
}

func TestParseIkemenVersion_StopsIntegerParseOnInvalidMiddleSegment(t *testing.T) {
	// When a segment has no digits, it becomes 0, but parsing continues for subsequent segments
	ver, verF := ParseIkemenVersion("1.bad.3")
	if ver != [3]uint16{1, 0, 3} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 0, 3})
	}
	if verF != 1.3 {
		t.Fatalf("verF = %v, want %v", verF, 1.3)
	}
}
