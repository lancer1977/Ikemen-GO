package main

import "testing"

func TestIniSectionReadI32ForStage_ParsesValuesAndStopsOnWhitespace(t *testing.T) {
	is := IniSection{
		"coords": " 10, -20 extra, 30",
	}
	a, b, c := int32(0), int32(0), int32(0)
	if !is.readI32ForStage("coords", &a, &b, &c) {
		t.Fatal("expected readI32ForStage to report success")
	}
	if a != 10 || b != -20 || c != 0 {
		t.Fatalf("readI32ForStage outputs = %v %v %v, want 10 -20 0", a, b, c)
	}
}

func TestIniSectionReadI32ForStage_ReturnsFalseForMissingValues(t *testing.T) {
	is := IniSection{}
	v := int32(7)
	if is.readI32ForStage("missing", &v) {
		t.Fatal("expected missing key to return false")
	}
	if v != 7 {
		t.Fatal("missing key should not modify output")
	}
}
