package main

import "testing"

func TestIniSectionReadF32ForStage_ParsesValuesAndStopsOnWhitespace(t *testing.T) {
	is := IniSection{
		"coords": " 1.5, -2.25 extra, 3.75",
	}
	a, b, c := float32(0), float32(0), float32(0)
	if !is.readF32ForStage("coords", &a, &b, &c) {
		t.Fatal("expected readF32ForStage to report success")
	}
	if a != 1.5 || b != -2.25 || c != 0 {
		t.Fatalf("readF32ForStage outputs = %v %v %v, want 1.5 -2.25 0", a, b, c)
	}
}

func TestIniSectionReadF32ForStage_ReturnsFalseForMissingValues(t *testing.T) {
	is := IniSection{}
	v := float32(7.25)
	if is.readF32ForStage("missing", &v) {
		t.Fatal("expected missing key to return false")
	}
	if v != 7.25 {
		t.Fatal("missing key should not modify output")
	}
}
