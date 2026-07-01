package main

import "testing"

func TestBytecodeParamClassifiers(t *testing.T) {
	t.Parallel()

	if !isPalFXParam(0) || !isPalFXParam(255) {
		t.Fatal("expected palfx params to accept low and high ids")
	}
	if isPalFXParam(256) {
		t.Fatal("expected palfx params to reject out-of-range ids")
	}

	if !isAfterImageParam(0) || !isAfterImageParam(255) {
		t.Fatal("expected afterimage params to accept low and high ids")
	}
	if isAfterImageParam(256) {
		t.Fatal("expected afterimage params to reject out-of-range ids")
	}

	if !isHitDefParam(0) || !isHitDefParam(255) {
		t.Fatal("expected hitdef params to accept low and high ids")
	}
	if isHitDefParam(256) {
		t.Fatal("expected hitdef params to reject out-of-range ids")
	}
}
