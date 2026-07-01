package main

import (
	"reflect"
	"testing"
)

func TestSplitAndTrim_SplitsAndTrimsEachSegmentAgain(t *testing.T) {
	got := SplitAndTrim("  a , b ,, c  ", ",")
	want := []string{"a", "b", "", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SplitAndTrim() = %#v, want %#v", got, want)
	}
}

func TestSplitAndTrim_PreservesSingleFieldAfterTrimAgain(t *testing.T) {
	got := SplitAndTrim("   lone value   ", ",")
	want := []string{"lone value"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SplitAndTrim(single) = %#v, want %#v", got, want)
	}
}
