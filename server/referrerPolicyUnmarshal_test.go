package server

import (
	"reflect"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

func Test_referrerPolicyUnmarshal_ValidYaml(t *testing.T) {
	v := viper.New()
	var c SecurityConfig
	hookOpt := viper.DecodeHook(ReferrerPolicyDecodeHook())

	y := `
ReferrerPolicy: no-referrer
`
	var tempMap map[string]any
	err := yaml.Unmarshal([]byte(y), &tempMap)
	if err != nil {
		panic(err)
	}
	v.MergeConfigMap(tempMap)
	if err := v.Unmarshal(&c, hookOpt); err != nil {
		t.Errorf("failed to unmarshal YAML: %v", err)
	}

	rp := c.ReferrerPolicy
	if nil == rp {
		t.Error("ReferrerPolicy: got nil, want value")
	}
	if "no-referrer" != rp.Value.Value() {
		t.Errorf("ReferrerPolicy: got = %q, want \"no-referrer\"", rp.Value.Value())
	}
}

func TestStringToReferrerPolicyValueHookFunc(t *testing.T) {
	p := reflect.TypeOf(ReferrerPolicy{})
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Valid/no-referrer",
			from: "no-referrer",
			want: ReferrerPolicy{referrerPolicyStringValue("no-referrer")},
		},
		{
			name: "Valid/no-referrer-when-downgrade",
			from: "no-referrer-when-downgrade",
			want: ReferrerPolicy{referrerPolicyStringValue("no-referrer-when-downgrade")},
		},
		{
			name: "Valid/origin",
			from: "origin",
			want: ReferrerPolicy{referrerPolicyStringValue("origin")},
		},
		{
			name: "Valid/origin-when-cross-origin",
			from: "origin-when-cross-origin",
			want: ReferrerPolicy{referrerPolicyStringValue("origin-when-cross-origin")},
		},
		{
			name: "Valid/same-origin",
			from: "same-origin",
			want: ReferrerPolicy{referrerPolicyStringValue("same-origin")},
		},
		{
			name: "Valid/strict-origin",
			from: "strict-origin",
			want: ReferrerPolicy{referrerPolicyStringValue("strict-origin")},
		},
		{
			name: "Valid/strict-origin-when-cross-origin",
			from: "strict-origin-when-cross-origin",
			want: ReferrerPolicy{referrerPolicyStringValue("strict-origin-when-cross-origin")},
		},
		{
			name: "Valid/unsafe-url",
			from: "unsafe-url",
			want: ReferrerPolicy{referrerPolicyStringValue("unsafe-url")},
		},
		{
			name:    "Invalid",
			from:    "foo",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := StringToReferrerPolicyValueHookFunc()
			from := reflect.ValueOf(tt.from)
			to := reflect.New(p).Elem()
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
