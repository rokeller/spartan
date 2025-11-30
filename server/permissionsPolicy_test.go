package server

import (
	"reflect"
	"testing"
)

func TestPermissionAllowValueMarkers(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		v    PermissionAllowValue
	}{
		{"AllowWildcardPermission", AllowWildcardPermission{}},
		{"AllowNonePermission", AllowNonePermission{}},
		{"AllowMultiplePermission", AllowMultiplePermission{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.PermissionAllowValueMarker()
		})
	}
}

func TestPermissionAllowListItemMarkers(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		v    PermissionAllowListItem
	}{
		{"AllowSelfPermission", AllowSelfPermission{}},
		{"AllowSrcPermission", AllowSrcPermission{}},
		{"AllowOriginPermission", AllowOriginPermission("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.PermissionAllowListItemMarker()
		})
	}
}

func Test_appendPermissionDirectiveValue(t *testing.T) {
	tests := []struct {
		name       string // description of this test case
		directives []string
		n          string
		v          string
		want       []string
	}{
		{
			name: "EmptyInputSlice",
			n:    "key",
			v:    "val",
			want: []string{"key=val"},
		},
		{
			name:       "NonEmptyInputSlice",
			directives: []string{"one"},
			n:          "test",
			v:          "yes",
			want:       []string{"one", "test=yes"},
		},
		{
			name:       "KeyOnly",
			directives: []string{"two"},
			n:          "key-only",
			want:       []string{"two", "key-only"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendPermissionDirectiveValue(tt.directives, tt.n, tt.v)
			// TODO: update the condition below to compare got with tt.want.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendPermissionDirectiveValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllowWildcardPermission_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "Value",
			want: "*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AllowWildcardPermission
			got := a.Value()
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllowMultiplePermission_Value(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		a    AllowMultiplePermission
		want string
	}{
		{
			name: "SingleAllowed",
			a:    AllowMultiplePermission{AllowSelfPermission{}},
			want: "(self)",
		},
		{
			name: "MultipleAllowed",
			a:    AllowMultiplePermission{AllowOriginPermission("test.origin"), AllowSrcPermission{}},
			want: "(\"test.origin\" src)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Value()
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
