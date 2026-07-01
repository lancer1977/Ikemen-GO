package main

import "testing"

func TestIniSectionGetText_ParsesQuotedStringsAndRejectsUnquotedValues(t *testing.T) {
	is := IniSection{
		"quoted":   `"hello world"`,
		"unquoted": `hello world`,
	}

	got, ok, err := is.getText("quoted")
	if !ok || err != nil || got != "hello world" {
		t.Fatalf("getText(quoted) = %q, %v, %v; want hello world, true, nil", got, ok, err)
	}

	got, ok, err = is.getText("unquoted")
	if !ok || err == nil || got != `hello world` {
		t.Fatalf("getText(unquoted) = %q, %v, %v; want original, true, error", got, ok, err)
	}

	got, ok, err = is.getText("missing")
	if ok || err != nil || got != "" {
		t.Fatalf("getText(missing) = %q, %v, %v; want empty, false, nil", got, ok, err)
	}
}
