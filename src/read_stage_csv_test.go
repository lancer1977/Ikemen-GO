package main

import "testing"

func TestIniSectionReadStageCsv_ParsesValuesAndStopsOnWhitespace(t *testing.T) {
	is := IniSection{
		"ints":   " 10, -20, 30 extra, 40",
		"floats": " 1.5, -2.25, 3.75 extra, 4.0",
	}

	ints := is.readI32CsvForStage("ints")
	if len(ints) != 3 || ints[0] != 10 || ints[1] != -20 || ints[2] != 30 {
		t.Fatalf("readI32CsvForStage = %#v, want [10 -20 30]", ints)
	}

	floats := is.readF32CsvForStage("floats")
	if len(floats) != 3 || floats[0] != 1.5 || floats[1] != -2.25 || floats[2] != 3.75 {
		t.Fatalf("readF32CsvForStage = %#v, want [1.5 -2.25 3.75]", floats)
	}
}

func TestIniSectionReadStageCsv_ReturnsEmptyForMissingKeys(t *testing.T) {
	is := IniSection{}
	if ints := is.readI32CsvForStage("missing"); len(ints) != 0 {
		t.Fatalf("expected empty int slice, got %#v", ints)
	}
	if floats := is.readF32CsvForStage("missing"); len(floats) != 0 {
		t.Fatalf("expected empty float slice, got %#v", floats)
	}
}
