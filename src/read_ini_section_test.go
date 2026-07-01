package main

import "testing"

func TestReadIniSection_SkipsJunkAndParsesHeader(t *testing.T) {
	lines := []string{
		"; comment",
		"  ",
		"[State -2, Intro]",
		"value = 10",
		"[Next]",
	}
	i := 0
	is, name, subname := ReadIniSection(lines, &i)
	if name != "state " || subname != "-2, Intro" {
		t.Fatalf("ReadIniSection header = %q, %q, want %q, %q", name, subname, "state ", "-2, Intro")
	}
	if is["value"] != "10" {
		t.Fatalf("ReadIniSection map = %#v, want value=10", is)
	}
	if i != 4 {
		t.Fatalf("ReadIniSection cursor = %d, want 4", i)
	}
}

func TestReadIniSection_ReturnsZeroValuesWhenNoHeaderFound(t *testing.T) {
	lines := []string{"value = 10", "still not a header"}
	i := 0
	is, name, subname := ReadIniSection(lines, &i)
	if is != nil || name != "" || subname != "" {
		t.Fatalf("ReadIniSection = %#v, %q, %q; want nil, empty, empty", is, name, subname)
	}
	if i != len(lines) {
		t.Fatalf("ReadIniSection cursor = %d, want %d", i, len(lines))
	}
}
