package main

import (
	"reflect"
	"testing"
)

func TestSplitIniListOutsideQuotesPreservesEscapedQuotesInsideTokens(t *testing.T) {
	got := splitIniListOutsideQuotes(`alpha, "bravo \"charlie\", delta", echo`)
	want := []string{`alpha`, `"bravo \"charlie\", delta"`, `echo`}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitIniListOutsideQuotes() = %#v, want %#v", got, want)
	}
}
