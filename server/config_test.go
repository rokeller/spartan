package server

import (
	"reflect"
	"testing"
	"time"
)

func TestSecurityConfig_GetContentSecurityPolicy(t *testing.T) {
	csp := &ContentSecurityPolicy{
		DefaultSrc: NonceValue{Placeholder: "foo"},
	}
	tests := []struct {
		name string
		csp  *ContentSecurityPolicy
		want *ContentSecurityPolicy
	}{
		{
			name: "Null CSP results in default CSP returned",
			want: &DefaultContentSecurityPolicy,
		},
		{
			name: "Non-null CSP results in CSP returned",
			csp:  csp,
			want: csp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SecurityConfig{
				ContentSecurityPolicy: tt.csp,
			}
			if got := c.GetContentSecurityPolicy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecurityConfig.GetContentSecurityPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurityConfig_GetContentTypeOptionsNoSniff(t *testing.T) {
	bt, bf := true, false
	type fields struct {
		ContentTypeOptionsNoSniff *bool
	}
	tests := []struct {
		name    string
		nosniff *bool
		want    bool
	}{
		{
			name: "Null value results in true",
			want: true,
		},
		{
			name:    "Non-null value results in value/true",
			nosniff: &bt,
			want:    true,
		},
		{
			name:    "Non-null value results in value/false",
			nosniff: &bf,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SecurityConfig{
				ContentTypeOptionsNoSniff: tt.nosniff,
			}
			if got := c.GetContentTypeOptionsNoSniff(); got != tt.want {
				t.Errorf("SecurityConfig.GetContentTypeOptionsNoSniff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurityConfig_GetPermissionsPolicy(t *testing.T) {
	pp := &PermissionsPolicy{
		Accelerometer: AllowWildcardPermission{},
	}
	tests := []struct {
		name string
		pp   *PermissionsPolicy
		want *PermissionsPolicy
	}{
		{
			name: "Null PermissionsPolicy results in default PermissionsPolicy returned",
			want: &DefaultPermissionsPolicy,
		},
		{
			name: "Non-null PermissionsPolicy results in PermissionsPolicy returned",
			pp:   pp,
			want: pp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SecurityConfig{
				PermissionsPolicy: tt.pp,
			}
			if got := c.GetPermissionsPolicy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecurityConfig.GetPermissionsPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurityConfig_GetReferrerPolicy(t *testing.T) {
	rp := &ReferrerPolicy{
		Value: UnsafeUrlReferrerPolicyValue,
	}
	tests := []struct {
		name string
		rp   *ReferrerPolicy
		want *ReferrerPolicy
	}{
		{
			name: "Null ReferrerPolicy results in default ReferrerPolicy returned",
			want: &DefaultReferrerPolicy,
		},
		{
			name: "Non-null ReferrerPolicy results in ReferrerPolicy returned",
			rp:   rp,
			want: rp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SecurityConfig{
				ReferrerPolicy: tt.rp,
			}
			if got := c.GetReferrerPolicy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecurityConfig.GetReferrerPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurityConfig_GetStrictTransportSecurityPolicy(t *testing.T) {
	stsp := &StrictTransportSecurityPolicy{
		MaxAge: time.Microsecond * 123,
	}
	tests := []struct {
		name string
		stsp *StrictTransportSecurityPolicy
		want *StrictTransportSecurityPolicy
	}{
		{
			name: "Null StrictTransportSecurityPolicy results in default StrictTransportSecurityPolicy returned",
			want: &DefaultStrictTransportSecurityPolicy,
		},
		{
			name: "Non-null StrictTransportSecurityPolicy results in StrictTransportSecurityPolicy returned",
			stsp: stsp,
			want: stsp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SecurityConfig{
				StrictTransportSecurityPolicy: tt.stsp,
			}
			if got := c.GetStrictTransportSecurityPolicy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecurityConfig.GetStrictTransportSecurityPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
