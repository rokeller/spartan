package server

import (
	"context"
	"testing"
)

func TestNoneDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'none'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n NoneDirectiveValue
			got := n.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelfDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'self'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s SelfDirectiveValue
			got := s.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsafeEvalDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'unsafe-eval'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u UnsafeEvalDirectiveValue
			got := u.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWasmUnsafeEvalDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'wasm-unsafe-eval'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w WasmUnsafeEvalDirectiveValue
			got := w.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsafeInlineDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'unsafe-inline'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u UnsafeInlineDirectiveValue
			got := u.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsafeHashesDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'unsafe-hashes'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u UnsafeHashesDirectiveValue
			got := u.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInlineSpeculationRulesDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'inline-speculation-rules'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i InlineSpeculationRulesDirectiveValue
			got := i.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrictDynamicDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'strict-dynamic'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s StrictDynamicDirectiveValue
			got := s.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReportSampleDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Works",
			want: "'report-sample'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r ReportSampleDirectiveValue
			got := r.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestNonceValue_Value(t *testing.T) {
// 	tests := []struct {
// 		name string // description of this test case
// 		want string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// TODO: construct the receiver type.
// 			var n NonceValue
// 			got := n.Value(context.Background())
// 			// TODO: update the condition below to compare got with tt.want.
// 			if true {
// 				t.Errorf("Value() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestHostSourceDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		input HostSourceDirectiveValue
		want  string
	}{
		{
			name:  "HostOnly",
			input: HostSourceDirectiveValue{Host: "a.b.c"},
			want:  "a.b.c",
		},
		{
			name:  "SchemeAndHost",
			input: HostSourceDirectiveValue{Host: "ftp.foo.bar", Scheme: toPtr("ftps")},
			want:  "ftps://ftp.foo.bar",
		},
		{
			name:  "SchemeWithColonAndHost",
			input: HostSourceDirectiveValue{Host: "ftp.foo.bar", Scheme: toPtr("ftps:")},
			want:  "ftps://ftp.foo.bar",
		},
		{
			name:  "HostAndPort",
			input: HostSourceDirectiveValue{Host: "localhost", Port: toPtr(uint16(8080))},
			want:  "localhost:8080",
		},
		{
			name:  "HostAndPath",
			input: HostSourceDirectiveValue{Host: "localhost", Path: toPtr("a/b/c")},
			want:  "localhost/a/b/c",
		},
		{
			name:  "HostAndPathWithLeadingSlash",
			input: HostSourceDirectiveValue{Host: "localhost", Path: toPtr("/a/b/c/")},
			want:  "localhost/a/b/c/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchemeSourceDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		input SchemeSourceDirectiveValue
		want  string
	}{
		{
			name:  "SchemeWithColon",
			input: SchemeSourceDirectiveValue("https:"),
			want:  "https:",
		},
		{
			name:  "SchemeWithoutColon",
			input: SchemeSourceDirectiveValue("ws"),
			want:  "ws:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubResourceHashDirectiveValue_Value(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		input SubResourceHashDirectiveValue
		want  string
	}{
		{
			name: "Works",
			input: SubResourceHashDirectiveValue{
				Alg:  "test",
				Hash: []byte{0xde, 0xad, 0xbe, 0xef, 0x01, 0x23},
			},
			want: "'test-3q2+7wEj'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Value(context.Background())
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
