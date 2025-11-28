package server

import (
	"reflect"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

func Test_permissionsPolicyUnmarshal_ValidYaml(t *testing.T) {
	v := viper.New()
	var c SecurityConfig
	hookOpt := viper.DecodeHook(PermissionPolicyDecodeHook())

	y := `
PermissionsPolicy:
  accelerometer: none
  autoplay: self
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

	pp := c.PermissionsPolicy
	if nil == pp {
		t.Error("PermissionsPolicy: got nil, want struct")
	}
	if !reflect.DeepEqual(pp.Accelerometer, AllowNonePermission{}) {
		t.Errorf("PermissionsPolicy Accelerometer: got = %v (%T), want ()", pp.Accelerometer, pp.Accelerometer)
	}
	if !reflect.DeepEqual(pp.Autoplay, AllowMultiplePermission{AllowSelfPermission{}}) {
		t.Errorf("PermissionsPolicy Autoplay: got = %v (%T), want [self]", pp.Autoplay, pp.Autoplay)
	}
	if !reflect.DeepEqual(pp.Camera, nil) {
		t.Errorf("PermissionsPolicy Camera: got = %v (%T), want nil", pp.Autoplay, pp.Autoplay)
	}
}

func TestSliceToPermissionAllowValueHookFunc(t *testing.T) {
	i := reflect.TypeOf((*PermissionAllowValue)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Valid/MultipleSupportedValue",
			from: []any{"src", "self", "https://foo.bar.com"},
			want: AllowMultiplePermission{AllowSrcPermission{}, AllowSelfPermission{}, AllowOriginPermission("https://foo.bar.com")},
		},
		{
			name:    "Invalid/Integer",
			from:    []any{1234},
			wantErr: true,
		},
		{
			name:    "Invalid/UrlWithTrailingSlash",
			from:    []any{"http://test.com/"},
			wantErr: true,
		},
		{
			name:    "Invalid/UrlWithPath",
			from:    []any{"http://test.com/foo/bar"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := SliceToPermissionAllowValueHookFunc()
			from := reflect.ValueOf(tt.from)
			to := reflect.New(i).Elem()
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

func TestStringToPermissionAllowValueHookFunc(t *testing.T) {
	i := reflect.TypeOf((*PermissionAllowValue)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Valid/Wildcard/*",
			from: "*",
			want: AllowWildcardPermission{},
		},
		{
			name: "Valid/Wildcard/all",
			from: "all",
			want: AllowWildcardPermission{},
		},
		{
			name: "Valid/Wildcard/wildcard",
			from: "wildcard",
			want: AllowWildcardPermission{},
		},
		{
			name: "Valid/None/none",
			from: "none",
			want: AllowNonePermission{},
		},
		{
			name: "Valid/None/()",
			from: "()",
			want: AllowNonePermission{},
		},
		{
			name: "Valid/Self",
			from: "self",
			want: AllowMultiplePermission{AllowSelfPermission{}},
		},
		{
			name: "Valid/Src",
			from: "src",
			want: AllowMultiplePermission{AllowSrcPermission{}},
		},
		{
			name: "Valid/Origin",
			from: "wss://baz.com",
			want: AllowMultiplePermission{AllowOriginPermission("wss://baz.com")},
		},
		{
			name:    "Invalid/UrlWithTrailingSlash",
			from:    "http://test.com/",
			wantErr: true,
		},
		{
			name:    "Invalid/UrlWithPath",
			from:    "http://test.com/foo/bar",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := StringToPermissionAllowValueHookFunc()
			from := reflect.ValueOf(tt.from)
			to := reflect.New(i).Elem()
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

func TestStringToPermissionAllowListItemHookFunc(t *testing.T) {
	i := reflect.TypeOf((*PermissionAllowListItem)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Valid/Self",
			from: "self",
			want: AllowSelfPermission{},
		},
		{
			name: "Valid/Src",
			from: "src",
			want: AllowSrcPermission{},
		},
		{
			name: "Valid/Origin",
			from: "wss://baz.com",
			want: AllowOriginPermission("wss://baz.com"),
		},
		{
			name:    "Invalid/UrlWithTrailingSlash",
			from:    "http://test.com/",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := StringToPermissionAllowListItemHookFunc()
			from := reflect.ValueOf(tt.from)
			to := reflect.New(i).Elem()
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

func Test_parsePermissionPolicyOrigin(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		str     string
		want    PermissionAllowListItem
		wantErr bool
	}{
		{
			name: "ValidOrigin",
			str:  "http://a.b.c",
			want: AllowOriginPermission("http://a.b.c"),
		},
		{
			name:    "InvalidOrigin/HasPath",
			str:     "http://a.b.c/path/to/resource",
			wantErr: true,
		},
		{
			name:    "InvalidOrigin/HasQuery",
			str:     "http://a.b.c?x=y",
			wantErr: true,
		},
		{
			name:    "InvalidOrigin/HasFragment",
			str:     "http://a.b.c#x/y/z",
			wantErr: true,
		},
		{
			name:    "InvalidOrigin/HasUser",
			str:     "http://user@a.b.c",
			wantErr: true,
		},
		{
			name:    "InvalidOrigin/InvalidUrl",
			str:     "abc\x09",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := parsePermissionPolicyOrigin(tt.str)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parsePermissionPolicyOrigin() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parsePermissionPolicyOrigin() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePermissionPolicyOrigin() = %v, want %v", got, tt.want)
			}
		})
	}
}
