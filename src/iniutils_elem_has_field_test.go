package main

import (
	"reflect"
	"testing"
)

func TestElemHasFieldWithINITag(t *testing.T) {
	t.Parallel()

	type child struct {
		Name string `ini:"name"`
	}

	if !elemHasFieldWithINITag(reflect.TypeOf(&child{}), "name") {
		t.Fatal("expected pointer element to expose ini tag")
	}
	if !elemHasFieldWithINITag(reflect.TypeOf(child{}), "name") {
		t.Fatal("expected struct element to expose ini tag")
	}
	if elemHasFieldWithINITag(reflect.TypeOf(child{}), "missing") {
		t.Fatal("unexpected match for missing ini tag")
	}
	if elemHasFieldWithINITag(reflect.TypeOf(0), "name") {
		t.Fatal("non-struct type should not match")
	}
}
