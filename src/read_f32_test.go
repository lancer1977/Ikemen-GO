package main

import "testing"

func TestIniSectionReadF32_ParsesCommaSeparatedFloats(t *testing.T) {
	is := IniSection{
		"coords": " 1.5, -2.25, 3.75 ",
	}
	a, b, c := float32(0), float32(0), float32(0)
	if !is.ReadF32("coords", &a, &b, &c) {
		t.Fatal("expected ReadF32 to report success")
	}
	if a != 1.5 || b != -2.25 || c != 3.75 {
		t.Fatalf("ReadF32 outputs = %v %v %v, want 1.5 -2.25 3.75", a, b, c)
	}
}

func TestIniSectionReadF32_ReturnsFalseForMissingValues(t *testing.T) {
	is := IniSection{}
	v := float32(7.25)
	if is.ReadF32("missing", &v) {
		t.Fatal("expected missing key to return false")
	}
	if v != 7.25 {
		t.Fatal("missing key should not modify output")
	}
}
