package main

import "testing"

func TestParseOverrideKey(t *testing.T) {
	t.Parallel()

	team, member, field, ok := parseOverrideKey(" p2.3.life ")
	if !ok || team != 1 || member != 2 || field != "life" {
		t.Fatalf("unexpected parse: team=%d member=%d field=%q ok=%v", team, member, field, ok)
	}

	if _, _, _, ok := parseOverrideKey("p3.1.life"); ok {
		t.Fatalf("expected invalid team to fail")
	}
	if _, _, _, ok := parseOverrideKey("p1.0.life"); ok {
		t.Fatalf("expected non-positive member to fail")
	}
}

func TestParseMapKVAndBoolLoose(t *testing.T) {
	t.Parallel()

	name, ok := parseMapKey(" map.   stage ")
	if !ok || name != "stage" {
		t.Fatalf("unexpected map parse: name=%q ok=%v", name, ok)
	}
	if _, ok := parseMapKey("stage"); ok {
		t.Fatalf("expected non-map key to fail")
	}

	key, value, ok := parseKV("enabled = on")
	if !ok || key != "enabled" || value != "on" {
		t.Fatalf("unexpected kv parse: key=%q value=%q ok=%v", key, value, ok)
	}
	if _, _, ok := parseKV("=missing"); ok {
		t.Fatalf("expected empty key to fail")
	}

	if v, ok := parseBoolLoose(" YES "); !ok || !v {
		t.Fatalf("expected yes to parse true")
	}
	if v, ok := parseBoolLoose("off"); !ok || v {
		t.Fatalf("expected off to parse false")
	}
	if _, ok := parseBoolLoose("maybe"); ok {
		t.Fatalf("expected unknown value to fail")
	}
}
