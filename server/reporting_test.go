package server

import (
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestReportingEndpoints_AddToResponse(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		input ReportingEndpoints
		want  string
	}{
		{
			name:  "Nil",
			input: nil,
			want:  "",
		},
		{
			name:  "Empty",
			input: []ReportingEndpoint{},
			want:  "",
		},
		{
			name: "Single",
			input: []ReportingEndpoint{
				{Name: "test", Url: "https://a.b.c"},
			},
			want: "test=\"https://a.b.c\"",
		},
		{
			name: "Multiple",
			input: []ReportingEndpoint{
				{Name: "first", Url: "https://a.com/first"},
				{Name: "second", Url: "https://a.com/second"},
			},
			want: "first=\"https://a.com/first\", second=\"https://a.com/second\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.input.AddToResponse(w)
			got := w.Result().Header.Get("Reporting-Endpoints")
			if got != tt.want {
				t.Errorf("Reporting-Endpoints value: got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReportingEndpoints_GetEndpoint(t *testing.T) {
	ep := ReportingEndpoints{
		ReportingEndpoint{"a", "https://a"},
		ReportingEndpoint{"b", "https://b"},
	}
	tests := []struct {
		name   string // description of this test case
		epName string
		want   *ReportingEndpoint
	}{
		{
			name:   "Found/First",
			epName: "a",
			want:   &ep[0],
		},
		{
			name:   "Found/Second",
			epName: "b",
			want:   &ep[1],
		},
		{
			name:   "NotFound",
			epName: "c",
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ep.GetEndpoint(tt.epName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
