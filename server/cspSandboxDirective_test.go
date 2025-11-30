package server

import (
	"context"
	"testing"
)

func TestSandboxDirectiveValueMarkers(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		v    SandboxDirectiveValue
	}{
		{"SandboxAll", SandboxAll{}},
		{"SandboxWithAllowed", SandboxWithAllowed{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.SandboxDirectiveValueMarker()
		})
	}
}

func TestSandboxWithAllowed_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		s    SandboxWithAllowed
		want string
	}{
		{
			name: "SingleAllowed",
			s:    SandboxWithAllowed{SandboxAllow("single")},
			want: "single",
		},
		{
			name: "MultipleAllowed",
			s:    SandboxWithAllowed{SandboxAllow("one"), SandboxAllow("two")},
			want: "one two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
