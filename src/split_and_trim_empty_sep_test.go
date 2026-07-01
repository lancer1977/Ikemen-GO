package main

import (
	"reflect"
	"testing"
)

func TestSplitAndTrim_SplitsIntoRunesWhenSeparatorIsEmpty(t *testing.T) {
	got := SplitAndTrim(" a ", "")
	want := []string{"", "a", ""}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SplitAndTrim(empty sep) = %#v, want %#v", got, want)
	}
}
