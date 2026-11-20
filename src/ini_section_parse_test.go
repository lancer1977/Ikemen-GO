package main

import "testing"

func TestNewIniSection_CreatesWritableEmptyMap(t *testing.T) {
	is := NewIniSection()
	if is == nil || len(is) != 0 {
		t.Fatalf("NewIniSection() = %#v, want empty map", is)
	}
	is["value"] = "10"
	if is["value"] != "10" {
		t.Fatalf("NewIniSection map did not accept writes: %#v", is)
	}
}

func TestIniSectionParse_StripsCommentsIgnoresInvalidLinesAndKeepsFirstValue(t *testing.T) {
	is := NewIniSection()
	lines := []string{
		"  ; comment",
		"invalid line",
		"foo = 1 ; trailing comment",
		"foo = 2",
		"bar 3",
		"[Next]",
	}
	i := 0
	is.Parse(lines, &i)
	if is["foo"] != "1" {
		t.Fatalf("Parse kept foo = %q, want %q", is["foo"], "1")
	}
	// "bar 3" has no "=", so Parse registers the key with empty data: the name
	// is taken up to the first "= \t" run, but data is only read after an "=".
	if v, ok := is["bar"]; !ok || v != "" {
		t.Fatalf("Parse kept bar = %q (present=%v), want an empty value", v, ok)
	}
	if i != 5 {
		t.Fatalf("Parse cursor = %d, want 5", i)
	}
}
