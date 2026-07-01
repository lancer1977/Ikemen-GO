package main

import (
	"reflect"
	"testing"
)

func TestGetFieldFromHierarchy(t *testing.T) {
	type primary struct {
		Direct string
	}
	type parent struct {
		Fallback string
	}

	pv := reflect.ValueOf(&primary{Direct: "primary"}).Elem()
	pa := reflect.ValueOf(&parent{Fallback: "parent"}).Elem()

	if got, ok := getFieldFromHierarchy(pv, pa, "Direct"); !ok || got.String() != "primary" {
		t.Fatalf("getFieldFromHierarchy(primary) = %#v %v", got, ok)
	}
	if got, ok := getFieldFromHierarchy(reflect.Value{}, pa, "Fallback"); !ok || got.String() != "parent" {
		t.Fatalf("getFieldFromHierarchy(parent) = %#v %v", got, ok)
	}
	if got, ok := getFieldFromHierarchy(reflect.Value{}, reflect.Value{}, "Missing"); ok || got.IsValid() {
		t.Fatalf("getFieldFromHierarchy(missing) = %#v %v", got, ok)
	}
}
