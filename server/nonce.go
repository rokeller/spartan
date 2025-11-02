package server

import (
	"context"
	"crypto/rand"
	"fmt"
)

const (
	KeyNonces = "SpartanNonces"
)

type NonceMap map[string]string

func WithNonceMap(ctx context.Context, n NonceMap) context.Context {
	return context.WithValue(ctx, KeyNonces, n)
}

func NonceMapFromContext(ctx context.Context) (NonceMap, bool) {
	val, ok := ctx.Value(KeyNonces).(NonceMap)
	return val, ok
}

func (c NonceMap) ValueForPlaceholder(placeholder string) string {
	val, ok := c[placeholder]
	if !ok {
		key := make([]byte, 32)
		rand.Read(key)
		val = fmt.Sprintf("%x", key)
		c[placeholder] = val
	}
	return val
}
