package main

import (
	"reflect"
	"testing"
)

func TestPaletteListGetFallbacks(t *testing.T) {
	var pl PaletteList
	if got := pl.Get(0); got != nil {
		t.Fatalf("expected nil from empty palette list, got %v", got)
	}

	pl.init()
	pl.SetSource(0, []uint32{10})
	pl.SetSource(1, []uint32{20})

	if got := pl.Get(1); !reflect.DeepEqual(got, []uint32{20}) {
		t.Fatalf("expected direct lookup to return slot 1, got %v", got)
	}
	if got := pl.Get(9); !reflect.DeepEqual(got, []uint32{10}) {
		t.Fatalf("expected out-of-range lookup to fall back to slot 0, got %v", got)
	}
}
