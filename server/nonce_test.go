package server

import (
	"context"
	"testing"
)

func TestWithNonceMap(t *testing.T) {
	ctx := context.TODO()
	n := NonceMap{}
	ctx2 := WithNonceMap(ctx, n)

	if _, ok := NonceMapFromContext(ctx); ok {
		t.Error("expected NonceMap missing in context")
	}

	n["foo"] = "bar"
	if n2, ok := NonceMapFromContext(ctx2); !ok {
		t.Error("expected NonceMap in context")
	} else if v, ok := n2["foo"]; !ok {
		t.Error("expected 'foo' in NonceMap")
	} else if v != "bar" {
		t.Errorf("expected 'bar' for 'foo' in NonceMap, got %q", v)
	}
}

func TestNonceMap_ValueForPlaceholder(t *testing.T) {
	n := NonceMap{}

	got := n.ValueForPlaceholder("something")
	if got == "" {
		t.Error("expected random value in NonceMap")
	}

	if got2 := n.ValueForPlaceholder("something"); got != got2 {
		t.Errorf("NonceMap.ValueForPlaceholder() = %v, want %v", got2, got)
	}
}
