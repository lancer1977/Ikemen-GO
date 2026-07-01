package main

import "testing"

func TestParseMugenVersion_ParsesKnownAndFallbackVersions(t *testing.T) {
	ver, verF := ParseMugenVersion("1.1")
	if ver != [2]uint16{1, 1} || verF != 1.1 {
		t.Fatalf("ParseMugenVersion(1.1) = %v, %v; want [1 1], 1.1", ver, verF)
	}

	ver, verF = ParseMugenVersion("1.0")
	if ver != [2]uint16{1, 0} || verF != 1.0 {
		t.Fatalf("ParseMugenVersion(1.0) = %v, %v; want [1 0], 1.0", ver, verF)
	}

	ver, verF = ParseMugenVersion("2.5")
	if ver != [2]uint16{2, 5} || verF != 0.5 {
		t.Fatalf("ParseMugenVersion(2.5) = %v, %v; want [2 5], 0.5", ver, verF)
	}

	ver, verF = ParseMugenVersion("1.bad")
	if ver != [2]uint16{} || verF != 0 {
		t.Fatalf("ParseMugenVersion invalid input = %v, %v; want zero values", ver, verF)
	}
}
