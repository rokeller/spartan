package server

import (
	"context"
	"testing"
)

func TestNoneOrSourceExpressionListMarkers(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		v    NoneOrSourceExpressionList
	}{
		{"NoneDirectiveValue", NoneDirectiveValue{}},
		{"SourceExpressionList", SourceExpressionList{}},
		{"SelfDirectiveValue", SelfDirectiveValue{}},
		{"HostSourceDirectiveValue", HostSourceDirectiveValue{}},
		{"SchemeSourceDirectiveValue", SchemeSourceDirectiveValue("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.NoneOrSourceExpressionListMarker()
		})
	}
}

func TestSourceExpressionListItemMarkers(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		v    SourceExpressionListItem
	}{
		{"SelfDirectiveValue", SelfDirectiveValue{}},
		{"HostSourceDirectiveValue", HostSourceDirectiveValue{}},
		{"SchemeSourceDirectiveValue", SchemeSourceDirectiveValue("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.SourceExpressionListItemMarker()
		})
	}
}

func TestSourceExpressionList_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		l    SourceExpressionList
		want string
	}{
		{
			name: "SingleValue",
			l: SourceExpressionList{
				SelfDirectiveValue{},
			},
			want: "'self'",
		},
		{
			name: "MultipleValues",
			l: SourceExpressionList{
				HostSourceDirectiveValue{Host: "unit.test.com", Scheme: toPtr("https")},
				SchemeSourceDirectiveValue("wss"),
			},
			want: "https://unit.test.com wss:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.l.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
