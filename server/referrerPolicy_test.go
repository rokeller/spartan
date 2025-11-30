package server

import (
	"net/http/httptest"
	"testing"
)

func TestReferrerPolicy_AddToResponse(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		input ReferrerPolicy
		want  string
	}{
		{
			name:  "Default",
			input: DefaultReferrerPolicy,
			want:  "same-origin",
		},
		{
			name:  "Policy/some-value",
			input: ReferrerPolicy{referrerPolicyStringValue("some-value")},
			want:  "some-value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.input.AddToResponse(w)
			got := w.Result().Header.Get("Referrer-Policy")
			if got != tt.want {
				t.Errorf("Referrer-Policy value: got = %q, want %q", got, tt.want)
			}
		})
	}
}
