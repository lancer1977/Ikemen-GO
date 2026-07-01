package main

import (
	"math"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaNumberToJsonFormatsNumericRangesAndRejectsInvalidValues(t *testing.T) {
	if got, err := luaNumberToJson(lua.LNumber(12)); err != nil || got != int64(12) {
		t.Fatalf("luaNumberToJson(int-like) = (%T) %#v, %v; want int64(12), nil", got, got, err)
	}
	if got, err := luaNumberToJson(lua.LNumber(12.5)); err != nil || got != float64(12.5) {
		t.Fatalf("luaNumberToJson(float) = (%T) %#v, %v; want float64(12.5), nil", got, got, err)
	}
	if got, err := luaNumberToJson(lua.LNumber(float64(^uint64(0)))); err != nil || got != uint64(^uint64(0)) {
		t.Fatalf("luaNumberToJson(uint64-range) = (%T) %#v, %v; want uint64(max), nil", got, got, err)
	}
	if _, err := luaNumberToJson(lua.LNumber(math.NaN())); err == nil {
		t.Fatal("luaNumberToJson should reject NaN")
	}
}
