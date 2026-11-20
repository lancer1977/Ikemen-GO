package main

import "testing"

func TestParseIkemenVersion_ArrayParserStopsOnNonNumericSegment(t *testing.T) {
	// FIXED: The array parser now strips non-numeric characters from each segment before parsing,
	// matching the float parser behavior. "1.02.3-beta+build.7" now yields [1, 2, 3] (fixed from [1, 2, 0]).
	// The float parser also correctly handles it as 1.0237.
	// Tracked as lancer1977/Ikemen-GO#14.
	ver, verF := ParseIkemenVersion("1.02.3-beta+build.7")
	if ver != [3]uint16{1, 2, 3} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 2, 3})
	}
	if verF != 1.0237 {
		t.Fatalf("verF = %v, want %v", verF, 1.0237)
	}

	// When the second segment has no digits, it becomes 0, but parsing continues for subsequent segments
	ver, verF = ParseIkemenVersion("1.bad.3")
	if ver != [3]uint16{1, 0, 3} {
		t.Fatalf("ver = %#v, want %#v", ver, [3]uint16{1, 0, 3})
	}
	if verF != 1.3 {
		t.Fatalf("verF = %v, want %v", verF, 1.3)
	}
}
