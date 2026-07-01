package main

import "testing"

func TestParseKVRejectsMissingAndEmptyKeys(t *testing.T) {
	if _, _, ok := parseKV("value only"); ok {
		t.Fatal("parseKV should reject strings without an equals sign")
	}

	if _, _, ok := parseKV(" = missing"); ok {
		t.Fatal("parseKV should reject empty keys")
	}
}
