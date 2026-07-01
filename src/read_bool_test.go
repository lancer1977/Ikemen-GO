package main

import "testing"

func TestIniSectionReadBool_ParsesCommaSeparatedBooleanValues(t *testing.T) {
	is := IniSection{
		"flags": "1, 0, 2",
	}
	a, b, c := false, true, false
	if !is.ReadBool("flags", &a, &b, &c) {
		t.Fatal("expected ReadBool to report success")
	}
	if !a || b || !c {
		t.Fatalf("ReadBool outputs = %v %v %v, want true false true", a, b, c)
	}
}

func TestIniSectionReadBool_ReturnsFalseForMissingOrEmptyValues(t *testing.T) {
	is := IniSection{}
	v := true
	if is.ReadBool("missing", &v) {
		t.Fatal("expected missing key to return false")
	}
	if !v {
		t.Fatal("missing key should not modify output")
	}

	is["empty"] = ""
	if is.ReadBool("empty", &v) {
		t.Fatal("expected empty value to return false")
	}
}
