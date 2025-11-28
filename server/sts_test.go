package server

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestStrictTransportSecurityPolicy_AddToResponse(t *testing.T) {
	type fields struct {
		Disabled          bool
		IncludeSubDomains bool
		MaxAge            time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Policy/Disabled",
			fields: fields{Disabled: true},
			want:   "",
		},
		{
			name:   "Policy/MaxAgeOnly",
			fields: fields{MaxAge: time.Second * 15},
			want:   "max-age=15",
		},
		{
			name:   "Policy/With-includeSubDomains",
			fields: fields{MaxAge: time.Second * 45, IncludeSubDomains: true},
			want:   "max-age=45; includeSubDomains",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := StrictTransportSecurityPolicy{
				Disabled:          tt.fields.Disabled,
				IncludeSubDomains: tt.fields.IncludeSubDomains,
				MaxAge:            tt.fields.MaxAge,
			}
			w := httptest.NewRecorder()
			p.AddToResponse(w)

			got := w.Result().Header.Get("Strict-Transport-Security")
			if tt.want != got {
				t.Errorf("header mismatch: got = %s, want = %s", got, tt.want)
			}
		})
	}
}
