package main

import "testing"

func TestAtob_ParsesCommonBooleanFormsAndRejectsInvalid(t *testing.T) {
	if !Atob("true") || !Atob("1") || Atob("false") || Atob("0") {
		t.Fatal("Atob did not parse common boolean forms as expected")
	}
	if Atob("not-a-bool") {
		t.Fatal("Atob should reject invalid input")
	}
}

func TestItoa_FormatsIntegersAsDecimalStrings(t *testing.T) {
	if got := Itoa(-42); got != "-42" {
		t.Fatalf("Itoa() = %q, want %q", got, "-42")
	}
	if got := Itoa(0); got != "0" {
		t.Fatalf("Itoa() = %q, want %q", got, "0")
	}
}
