package main

import "testing"

func TestParseMapKeyRejectsEmptyMapName(t *testing.T) {
	if _, ok := parseMapKey("map."); ok {
		t.Fatal("parseMapKey should reject an empty map name")
	}
}
