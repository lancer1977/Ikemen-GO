package main

import "testing"

func TestSyncFontsMap(t *testing.T) {
	var dst map[string]*FontProperties
	fonts := map[int]*Fnt{
		2: {Type: 7, Size: 12, Spacing: 3, offset: [2]int32{4, 5}},
	}
	indexByKey := map[string]int{
		fontKey("font/select.fnt", 16): 2,
	}

	syncFontsMap(&dst, fonts, indexByKey)

	if dst == nil {
		t.Fatal("syncFontsMap should initialize destination map")
	}
	fp := dst["font2"]
	if fp == nil {
		t.Fatal("syncFontsMap should create font2 entry")
	}
	if fp.Font != "font/select.fnt" || fp.Height != 16 {
		t.Fatalf("syncFontsMap font metadata = %#v", fp)
	}
	if fp.Type != 7 || fp.Size != 12 || fp.Spacing != 3 || fp.Offset != [2]int32{4, 5} {
		t.Fatalf("syncFontsMap copied fields = %#v", fp)
	}
}
