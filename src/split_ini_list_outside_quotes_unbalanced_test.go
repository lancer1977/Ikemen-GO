package main

import (
	"reflect"
	"testing"
)

func TestSplitIniListOutsideQuotesTreatsUnbalancedQuotesAsScalar(t *testing.T) {
	got := splitIniListOutsideQuotes(`alpha, "bravo, charlie`)
	want := []string{`alpha`, `"bravo, charlie`}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitIniListOutsideQuotes() = %#v, want %#v", got, want)
	}
}
