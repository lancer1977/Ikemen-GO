package main

import (
	"reflect"
	"testing"
)

func TestSetDefaultFieldValue(t *testing.T) {
	t.Parallel()

	type nested struct {
		A int
		B string
	}

	var s string
	setDefaultFieldValue(reflect.ValueOf(&s).Elem(), reflect.StructField{}, "hello")
	if s != "hello" {
		t.Fatalf("string default = %q, want hello", s)
	}

	var n int
	setDefaultFieldValue(reflect.ValueOf(&n).Elem(), reflect.StructField{}, "42")
	if n != 42 {
		t.Fatalf("int default = %d, want 42", n)
	}

	var b bool
	setDefaultFieldValue(reflect.ValueOf(&b).Elem(), reflect.StructField{}, "true")
	if !b {
		t.Fatal("bool default = false, want true")
	}

	var arr [2]string
	setDefaultFieldValue(reflect.ValueOf(&arr).Elem(), reflect.StructField{}, "x,y,z")
	if arr != [2]string{"x", "y"} {
		t.Fatalf("array default = %#v, want [x y]", arr)
	}

	var sl []int
	setDefaultFieldValue(reflect.ValueOf(&sl).Elem(), reflect.StructField{}, "1,2,3")
	if len(sl) != 3 || sl[0] != 1 || sl[1] != 2 || sl[2] != 3 {
		t.Fatalf("slice default = %#v, want [1 2 3]", sl)
	}

	var st nested
	setDefaultFieldValue(reflect.ValueOf(&st).Elem(), reflect.StructField{}, "7,hi")
	if st.A != 7 || st.B != "hi" {
		t.Fatalf("struct default = %#v, want {7 hi}", st)
	}
}
