package main

import (
	"arena"
	"testing"
)

func TestCloneDialogueToken(t *testing.T) {
	t.Parallel()

	src := DialogueToken{param: "p", side: 2, redirection: "self", pn: 3, value: []interface{}{"a", 1}}
	dst := cloneDialogueToken(src)

	if dst.param != src.param || dst.side != src.side || dst.redirection != src.redirection || dst.pn != src.pn {
		t.Fatalf("scalar fields not cloned: %#v", dst)
	}
	if len(dst.value) != len(src.value) || dst.value[0] != src.value[0] || dst.value[1] != src.value[1] {
		t.Fatalf("value slice not cloned correctly: %#v", dst.value)
	}
	dst.value[0] = "changed"
	if src.value[0] != "a" {
		t.Fatalf("original token mutated through clone: %#v", src.value)
	}
}

func TestCloneDialogueParsedLine(t *testing.T) {
	t.Parallel()

	a := arena.NewArena()
	defer a.Free()

	src := DialogueParsedLine{
		side:     1,
		text:     "line",
		typedCnt: 4,
		tokens: map[int][]DialogueToken{
			2: {
				{param: "x", value: []interface{}{"x"}},
				{param: "y", value: []interface{}{"y"}},
			},
		},
	}

	dst := cloneDialogueParsedLine(a, src)
	if dst.side != src.side || dst.text != src.text || dst.typedCnt != src.typedCnt {
		t.Fatalf("scalar fields not cloned: %#v", dst)
	}
	if len(dst.tokens) != len(src.tokens) {
		t.Fatalf("tokens map size = %d, want %d", len(dst.tokens), len(src.tokens))
	}

	clonedTokens := dst.tokens[2]
	if len(clonedTokens) != len(src.tokens[2]) {
		t.Fatalf("tokens slice length = %d, want %d", len(clonedTokens), len(src.tokens[2]))
	}
	if clonedTokens[0].value[0] != "x" || clonedTokens[1].value[0] != "y" {
		t.Fatalf("nested token values not cloned: %#v", clonedTokens)
	}

	clonedTokens[0].value[0] = "changed"
	if src.tokens[2][0].value[0] != "x" {
		t.Fatalf("original parsed line mutated through clone: %#v", src.tokens[2])
	}
}
