package server

import (
	"context"
	"strings"
)

type FetchDirectiveValue interface {
	DirectiveValue
	FetchDirectiveValueMarker()
}

// MultiFetchDirectiveValue combines multiple fetch directive values.
type MultiFetchDirectiveValue []FetchDirectiveValue

func (mv MultiFetchDirectiveValue) Value(ctx context.Context) string {
	parts := make([]string, len(mv))
	for i, val := range mv {
		parts[i] = val.Value(ctx)
	}

	return strings.Join(parts, " ")
}
func (MultiFetchDirectiveValue) FetchDirectiveValueMarker() {}

var _ FetchDirectiveValue = MultiFetchDirectiveValue([]FetchDirectiveValue{})
var _ FetchDirectiveValue = NoneDirectiveValue{}
var _ FetchDirectiveValue = SelfDirectiveValue{}
var _ FetchDirectiveValue = UnsafeEvalDirectiveValue{}
var _ FetchDirectiveValue = WasmUnsafeEvalDirectiveValue{}
var _ FetchDirectiveValue = UnsafeInlineDirectiveValue{}
var _ FetchDirectiveValue = UnsafeHashesDirectiveValue{}
var _ FetchDirectiveValue = InlineSpeculationRulesDirectiveValue{}
var _ FetchDirectiveValue = StrictDynamicDirectiveValue{}
var _ FetchDirectiveValue = ReportSampleDirectiveValue{}
var _ FetchDirectiveValue = NonceValue{}
var _ FetchDirectiveValue = HostSourceDirectiveValue{}
var _ FetchDirectiveValue = SchemeSourceDirectiveValue("")
var _ FetchDirectiveValue = SubResourceHashDirectiveValue{}
