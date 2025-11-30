package server

import (
	"context"
	"testing"
)

func TestFetchDirectiveValueMarkers(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		v    FetchDirectiveValue
	}{
		{"MultiFetchDirectiveValue", MultiFetchDirectiveValue([]FetchDirectiveValue{})},
		{"NoneDirectiveValue", NoneDirectiveValue{}},
		{"SelfDirectiveValue", SelfDirectiveValue{}},
		{"UnsafeEvalDirectiveValue", UnsafeEvalDirectiveValue{}},
		{"WasmUnsafeEvalDirectiveValue", WasmUnsafeEvalDirectiveValue{}},
		{"UnsafeInlineDirectiveValue", UnsafeInlineDirectiveValue{}},
		{"UnsafeHashesDirectiveValue", UnsafeHashesDirectiveValue{}},
		{"InlineSpeculationRulesDirectiveValue", InlineSpeculationRulesDirectiveValue{}},
		{"StrictDynamicDirectiveValue", StrictDynamicDirectiveValue{}},
		{"ReportSampleDirectiveValue", ReportSampleDirectiveValue{}},
		{"NonceValue", NonceValue{}},
		{"HostSourceDirectiveValue", HostSourceDirectiveValue{}},
		{"SchemeSourceDirectiveValue", SchemeSourceDirectiveValue("")},
		{"SubResourceHashDirectiveValue", SubResourceHashDirectiveValue{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.FetchDirectiveValueMarker()
		})
	}
}

func TestMultiFetchDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		mv   MultiFetchDirectiveValue
		want string
	}{
		{
			name: "SingleValue",
			mv: MultiFetchDirectiveValue{
				NoneDirectiveValue{},
			},
			want: "'none'",
		},
		{
			name: "MultipleValues",
			mv: MultiFetchDirectiveValue{
				SelfDirectiveValue{},
				HostSourceDirectiveValue{Host: "test.com"},
			},
			want: "'self' test.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mv.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
