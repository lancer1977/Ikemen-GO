package main

import "testing"

func TestParseIkemenVersion_ArrayParserStopsOnNonNumericSegment(t *testing.T) {
	// DEFECT: The array parser stops when it encounters a segment that fails strconv.ParseUint.
	// "1.02.3-beta+build.7" stops at "3-beta+build" (non-numeric), yielding [1, 2, 0].
	// The float parser correctly strips all non-numeric characters except dots.
	// Tracked as lancer1977/Ikemen-GO#14.
	ver, verF := ParseIkemenVersion("1.02.3-beta+build.7")
	if ver != [3]uint16{1, 2, 0} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 2, 0})
	}
	if verF != 1.0237 {
		t.Fatalf("verF = %v, want %v", verF, 1.0237)
	}

	// When the second segment is invalid, parsing stops early
	ver, verF = ParseIkemenVersion("1.bad.3")
	if ver != [3]uint16{1, 0, 0} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 0, 0})
	}
	if verF != 1.3 {
		t.Fatalf("verF = %v, want %v", verF, 1.3)
	}
}
