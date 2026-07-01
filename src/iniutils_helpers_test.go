package main

import (
	"reflect"
	"testing"
)

func TestIniutilsSplitLangPrefix_ParsesOnlyTwoLetterPrefixes(t *testing.T) {
	if lang, base, ok := splitLangPrefix("pt.select"); !ok || lang != "pt" || base != "select" {
		t.Fatalf("splitLangPrefix(pt.select) = %q %q %v", lang, base, ok)
	}
	if lang, base, ok := splitLangPrefix("EN.select"); !ok || lang != "en" || base != "select" {
		t.Fatalf("splitLangPrefix(EN.select) = %q %q %v", lang, base, ok)
	}
	if lang, base, ok := splitLangPrefix("select"); ok || lang != "" || base != "select" {
		t.Fatalf("splitLangPrefix(select) = %q %q %v", lang, base, ok)
	}
}

func TestIniutilsTagHelpers_RespectExplicitStructTags(t *testing.T) {
	type tagged struct {
		Insensitive string `insensitivekeys:"false"`
		KeyFirst    string `keyfirst:"yes"`
		Literal     string `ini:"select" literal:"true"`
	}
	typ := reflect.TypeOf(tagged{})

	if tagInsensitiveKeys(typ.Field(0)) {
		t.Fatal("tagInsensitiveKeys should honor insensitivekeys:false")
	}
	if !tagKeyFirstMap(typ.Field(1)) {
		t.Fatal("tagKeyFirstMap should honor keyfirst:yes")
	}
	if !isLiteralSectionFor(&tagged{}, "select") {
		t.Fatal("isLiteralSectionFor should detect matching literal section tag")
	}
	if isLiteralSectionFor(&tagged{}, "other") {
		t.Fatal("isLiteralSectionFor should reject non-matching sections")
	}
}

func TestIniutilsParseQueryPath_SeparatesIndexesAndNames(t *testing.T) {
	parts := parseQueryPath("root.child[2].leaf")
	if len(parts) != 3 {
		t.Fatalf("parseQueryPath length = %d, want 3", len(parts))
	}
	if parts[0].name != "root" || parts[0].index != nil {
		t.Fatalf("parseQueryPath part 0 = %#v", parts[0])
	}
	if parts[1].name != "child" || parts[1].index == nil || *parts[1].index != "2" {
		t.Fatalf("parseQueryPath part 1 = %#v", parts[1])
	}
	if parts[2].name != "leaf" || parts[2].index != nil {
		t.Fatalf("parseQueryPath part 2 = %#v", parts[2])
	}
}
