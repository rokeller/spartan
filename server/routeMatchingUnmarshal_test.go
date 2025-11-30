package server

import (
	"reflect"
	"testing"

	"github.com/go-viper/mapstructure/v2"
)

func TestMapToRouteMatcherHookFunc(t *testing.T) {
	m := reflect.TypeOf((*RouteMatcher)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Valid/PathPrefix",
			from: map[string]any{
				"pathprefix": "/test",
			},
			want: PathPrefixMatcher{PathPrefix: "/test"},
		},
		{
			name: "Invalid/Array",
			from: []any{123},
			want: []any{123},
		},
		{
			name: "Invalid/IntegerPathPrefix",
			from: map[string]any{
				"pathprefix": 123,
			},
			wantErr: true,
		},
		{
			name: "Invalid/MapWithoutPathPrefix",
			from: map[string]any{
				"something-else": "blah",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := MapToRouteMatcherHookFunc()
			from := reflect.ValueOf(tt.from)
			to := reflect.New(m).Elem()
			got, err := mapstructure.DecodeHookExec(hook, from, to)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeHookExec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeHookExec(): got = %v (%T), want %v", got, got, tt.want)
			}
		})
	}
}
