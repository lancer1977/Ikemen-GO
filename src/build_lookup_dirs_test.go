package main

import (
	"reflect"
	"testing"
)

func TestBuildLookupDirs_ExpandsDefAndPreservesEmptyEntries(t *testing.T) {
	got := buildLookupDirs(" def , , data/ ", "base.def")
	want := []string{"base.def", "", "data/"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildLookupDirs() = %#v, want %#v", got, want)
	}

	if got := buildLookupDirs("", "base.def"); got != nil {
		t.Fatalf("buildLookupDirs(empty) = %#v, want nil", got)
	}
}
