package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
)

type DirectiveValue interface {
	Value(ctx context.Context) string
}

// NoneDirectiveValue represents the 'none' directive value.
type NoneDirectiveValue struct{}

func (NoneDirectiveValue) Value(ctx context.Context) string {
	return "'none'"
}
func (NoneDirectiveValue) FetchDirectiveValueMarker()        {}
func (NoneDirectiveValue) NoneOrSourceExpressionListMarker() {}

// SelfDirectiveValue represents the 'self' directive value.
type SelfDirectiveValue struct{}

func (SelfDirectiveValue) Value(ctx context.Context) string {
	return "'self'"
}
func (SelfDirectiveValue) FetchDirectiveValueMarker()        {}
func (SelfDirectiveValue) NoneOrSourceExpressionListMarker() {}
func (SelfDirectiveValue) SourceExpressionListItemMarker()   {}

// UnsafeEvalDirectiveValue represents the 'unsafe-eval' directive value.
type UnsafeEvalDirectiveValue struct{}

func (UnsafeEvalDirectiveValue) Value(ctx context.Context) string {
	return "'unsafe-eval'"
}
func (UnsafeEvalDirectiveValue) FetchDirectiveValueMarker() {}

// WasmUnsafeEvalDirectiveValue represents the 'wasm-unsafe-eval' directive value.
type WasmUnsafeEvalDirectiveValue struct{}

func (WasmUnsafeEvalDirectiveValue) Value(ctx context.Context) string {
	return "'wasm-unsafe-eval'"
}
func (WasmUnsafeEvalDirectiveValue) FetchDirectiveValueMarker() {}

// UnsafeInlineDirectiveValue represents the 'unsafe-inline' directive value.
type UnsafeInlineDirectiveValue struct{}

func (UnsafeInlineDirectiveValue) Value(ctx context.Context) string {
	return "'unsafe-inline'"
}
func (UnsafeInlineDirectiveValue) FetchDirectiveValueMarker() {}

// UnsafeHashesDirectiveValue represents the 'unsafe-hashes' directive value.
type UnsafeHashesDirectiveValue struct{}

func (UnsafeHashesDirectiveValue) Value(ctx context.Context) string {
	return "'unsafe-hashes'"
}
func (UnsafeHashesDirectiveValue) FetchDirectiveValueMarker() {}

// InlineSpeculationRulesDirectiveValue represents the 'inline-speculation-rules' directive value.
type InlineSpeculationRulesDirectiveValue struct{}

func (InlineSpeculationRulesDirectiveValue) Value(ctx context.Context) string {
	return "'inline-speculation-rules'"
}
func (InlineSpeculationRulesDirectiveValue) FetchDirectiveValueMarker() {}

// StrictDynamicDirectiveValue represents the 'strict-dynamic' directive value.
type StrictDynamicDirectiveValue struct{}

func (StrictDynamicDirectiveValue) Value(ctx context.Context) string {
	return "'strict-dynamic'"
}
func (StrictDynamicDirectiveValue) FetchDirectiveValueMarker() {}

// ReportSampleDirectiveValue represents the 'report-sample' directive value.
type ReportSampleDirectiveValue struct{}

func (ReportSampleDirectiveValue) Value(ctx context.Context) string {
	return "'report-sample'"
}
func (ReportSampleDirectiveValue) FetchDirectiveValueMarker() {}

// NonceValue represents a 'nonce-<nonce_value>' directive value.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#nonce-nonce_value
type NonceValue struct {
	Placeholder string
}

func (n NonceValue) Value(ctx context.Context) string {
	nm, ok := NonceMapFromContext(ctx)
	if !ok {
		klog.Error(errors.New("missing nonce context"))
		return ""
	}
	rand := nm.ValueForPlaceholder(n.Placeholder)
	return fmt.Sprintf("'nonce-%s'", rand)
}
func (NonceValue) FetchDirectiveValueMarker() {}

// HostSourceDirectiveValue represents a <host-source> directive value.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#host-source
type HostSourceDirectiveValue struct {
	Scheme *string
	Host   string
	Port   *uint16
	Path   *string
}

func (v HostSourceDirectiveValue) Value(ctx context.Context) string {
	val := ""
	if v.Scheme != nil {
		if !strings.HasSuffix(*v.Scheme, ":") {
			*v.Scheme += ":"
		}
		val = *v.Scheme + "//"
	}
	val += v.Host
	if v.Port != nil {
		val += fmt.Sprintf(":%d", *v.Port)
	}
	if v.Path != nil {
		if !strings.HasPrefix(*v.Path, "/") {
			*v.Path = "/" + *v.Path
		}
		val += *v.Path
	}
	return val
}
func (HostSourceDirectiveValue) FetchDirectiveValueMarker()        {}
func (HostSourceDirectiveValue) NoneOrSourceExpressionListMarker() {}
func (HostSourceDirectiveValue) SourceExpressionListItemMarker()   {}

// SchemeSourceDirectiveValue represents a <scheme-source> directive value.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#scheme-source
type SchemeSourceDirectiveValue string

func (v SchemeSourceDirectiveValue) Value(ctx context.Context) string {
	s := string(v)
	if !strings.HasSuffix(s, ":") {
		s += ":"
	}

	return s
}
func (SchemeSourceDirectiveValue) FetchDirectiveValueMarker()        {}
func (SchemeSourceDirectiveValue) NoneOrSourceExpressionListMarker() {}
func (SchemeSourceDirectiveValue) SourceExpressionListItemMarker()   {}

// SubResourceHashDirectiveValue represents a '<hash_algorithm>-<hash_value>' directive value.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#hash_algorithm-hash_value
type SubResourceHashDirectiveValue struct {
	Alg  string
	Hash []byte
}

func (v SubResourceHashDirectiveValue) Value(ctx context.Context) string {
	b64 := base64.StdEncoding.EncodeToString(v.Hash)
	return fmt.Sprintf("'%s-%s'", v.Alg, b64)
}
func (SubResourceHashDirectiveValue) FetchDirectiveValueMarker() {}
