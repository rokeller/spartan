package server

import (
	"testing"
)

func TestTLSConfig_Empty(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		c    *TLSConfig
		want bool
	}{
		{
			name: "Nil",
			want: true,
		},
		{
			name: "NonNil/WithoutPaths",
			c:    &TLSConfig{},
			want: true,
		},
		{
			name: "NonNil/WithCertPath",
			c:    &TLSConfig{KeyPath: "foo"},
			want: true,
		},
		{
			name: "NonNil/WithKeyPath",
			c:    &TLSConfig{CertPath: "foo"},
			want: true,
		},
		{
			name: "NonNil/WithAll",
			c:    &TLSConfig{CertPath: "crt", KeyPath: "key"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.Empty()
			if got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}
