package server

import (
	"net/http/httptest"
	"testing"
)

func TestPathPrefixMatcher_Match(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		p       PathPrefixMatcher
		path    string
		wantLen int
		wantOk  bool
	}{
		{
			name:    "NoMatch",
			p:       PathPrefixMatcher{PathPrefix: "/abc"},
			path:    "/def",
			wantLen: 0,
			wantOk:  false,
		},
		{
			name:    "PartialMatch",
			p:       PathPrefixMatcher{PathPrefix: "/abc"},
			path:    "/abc/def/ghi",
			wantLen: 4,
			wantOk:  true,
		},
		{
			name:    "FullMatch",
			p:       PathPrefixMatcher{PathPrefix: "/123/456"},
			path:    "/123/456",
			wantLen: 8,
			wantOk:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			gotLen, gotOk := tt.p.Match(req)
			if gotLen != tt.wantLen {
				t.Errorf("Match() length got = %d, want %d", gotLen, tt.wantLen)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Match() ok got = %t, want %t", gotOk, tt.wantOk)
			}
		})
	}
}

func TestPathPrefixMatcher_RouteString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		p    PathPrefixMatcher
		want string
	}{
		{
			name: "PathPrefix/abc",
			p:    PathPrefixMatcher{PathPrefix: "/abc"},
			want: "PathPrefix: /abc",
		},
		{
			name: "PathPrefix/123/456",
			p:    PathPrefixMatcher{PathPrefix: "/123/456"},
			want: "PathPrefix: /123/456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.RouteString()
			if got != tt.want {
				t.Errorf("RouteString() = %q, want %q", got, tt.want)
			}
		})
	}
}
