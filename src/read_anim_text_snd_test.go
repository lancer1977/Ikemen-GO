package main

import "testing"

func TestReadAnimTextSnd_PopulatesSoundTextLayoutAndDisplaytime(t *testing.T) {
	is := IniSection{
		"intro.snd":         " 3, 4 ",
		"intro.text":        `"Hello\nWorld"`,
		"intro.displaytime": " 60 ",
	}

	ats := ReadAnimTextSnd("intro.", is, &Sff{}, AnimationTable{}, 7, map[int]*Fnt{})
	if ats == nil {
		t.Fatal("expected ReadAnimTextSnd to return a value")
	}
	if ats.snd != [2]int32{3, 4} {
		t.Fatalf("snd = %#v, want [3 4]", ats.snd)
	}
	if ats.text.text != "Hello\\nWorld" {
		t.Fatalf("text = %q, want %q", ats.text.text, "Hello\\nWorld")
	}
	// Layout.Read caps layerno to Min(2, ln), so 7 becomes 2
	if ats.animLayout.lay.layerno != 2 {
		t.Fatalf("layerno = %d, want 2", ats.animLayout.lay.layerno)
	}
	if ats.displaytime != 60 {
		t.Fatalf("displaytime = %d, want 60", ats.displaytime)
	}
}
