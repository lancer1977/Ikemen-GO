package main

import "testing"

func TestSectionName_ReturnsZeroValuesForMalformedInputs(t *testing.T) {
	name, body := SectionName("State 100]")
	if name != "" || body != "" {
		t.Fatalf("SectionName(no leading bracket) = %q, %q; want empty values", name, body)
	}

	name, body = SectionName("[State 100")
	if name != "" || body != "" {
		t.Fatalf("SectionName(no closing bracket) = %q, %q; want empty values", name, body)
	}
}
