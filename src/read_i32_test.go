package main

import "testing"

func TestIniSectionReadI32_ParsesCommaSeparatedIntegers(t *testing.T) {
	is := IniSection{
		"coords": " 10, -20, 30 ",
	}
	a, b, c := int32(0), int32(0), int32(0)
	if !is.ReadI32("coords", &a, &b, &c) {
		t.Fatal("expected ReadI32 to report success")
	}
	if a != 10 || b != -20 || c != 30 {
		t.Fatalf("ReadI32 outputs = %v %v %v, want 10 -20 30", a, b, c)
	}
}

func TestIniSectionReadI32_ReturnsFalseForMissingValues(t *testing.T) {
	is := IniSection{}
	v := int32(7)
	if is.ReadI32("missing", &v) {
		t.Fatal("expected missing key to return false")
	}
	if v != 7 {
		t.Fatal("missing key should not modify output")
	}
}
