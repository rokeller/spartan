package server

import (
	"reflect"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

func Test_cspUnmarshal_ValidYaml(t *testing.T) {
	v := viper.New()
	var c SecurityConfig
	hookOpt := viper.DecodeHook(CspDecodeHook())

	y := `
ContentSecurityPolicy:
  reportOnly: true
  reportTo: test
  defaultSrc: "'none'"
  imgSrc: self
  scriptSrc:
  - nonce: CSP_NONCE_PLACEHOLDER
  - host:
      host: test.com
      scheme: https
  scriptSrcElem:
  - host: foo.bar.com
  - scheme: ws
  - hash:
      alg:  sha384
      hash: QQJt2MfQO6NNZigaUfp4MyrsRYl1c7vCw89PXW1VbwiJRML9qDh5sY1QtbcyNeR4
  - nonce:
      placeholder: blah
  styleSrc: unsafe-eval
  styleSrcElem:
  - wasm-unsafe-eval
  - unsafe-inline
  - unsafe-hashes
  - inline-speculation-rules
  - strict-dynamic
  - report-sample

  sandbox:
    - allow-scripts

  formAction:
    - self
    - host: a.b.com
    - host:
        host: my.test.com
    - scheme: http
  frameAncestors: none
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

	csp := c.ContentSecurityPolicy
	if nil == csp {
		t.Error("CSP: got nil, want struct")
	}
	if !csp.ReportOnly {
		t.Errorf("CSP ReportOnly: got = %t, want true", csp.ReportOnly)
	}
	if csp.ReportTo != "test" {
		t.Errorf("CSP ReportTo: got = %q, want 'test'", csp.ReportTo)
	}
	if !reflect.DeepEqual(csp.DefaultSrc, NoneDirectiveValue{}) {
		t.Errorf("CSP DefaultSrc: got = %v (%T), want 'none'", csp.DefaultSrc, csp.DefaultSrc)
	}
	if !reflect.DeepEqual(csp.ImgSrc, SelfDirectiveValue{}) {
		t.Errorf("CSP ImgSrc: got = %v (%T), want 'self'", csp.ImgSrc, csp.ImgSrc)
	}
	if !reflect.DeepEqual(csp.ScriptSrc, MultiFetchDirectiveValue{
		NonceValue{Placeholder: "CSP_NONCE_PLACEHOLDER"},
		HostSourceDirectiveValue{Host: "test.com", Scheme: toPtr("https")},
	}) {
		t.Errorf("CSP ScriptSrc: got = %v (%T), want multiple", csp.ScriptSrc, csp.ScriptSrc)
	}
	if !reflect.DeepEqual(csp.ScriptSrcElem, MultiFetchDirectiveValue{
		HostSourceDirectiveValue{Host: "foo.bar.com"},
		SchemeSourceDirectiveValue("ws"),
		SubResourceHashDirectiveValue{Alg: "sha384", Hash: []byte{0x41, 0x02, 0x6d, 0xd8, 0xc7, 0xd0, 0x3b, 0xa3, 0x4d, 0x66, 0x28, 0x1a, 0x51, 0xfa, 0x78, 0x33, 0x2a, 0xec, 0x45, 0x89, 0x75, 0x73, 0xbb, 0xc2, 0xc3, 0xcf, 0x4f, 0x5d, 0x6d, 0x55, 0x6f, 0x08, 0x89, 0x44, 0xc2, 0xfd, 0xa8, 0x38, 0x79, 0xb1, 0x8d, 0x50, 0xb5, 0xb7, 0x32, 0x35, 0xe4, 0x78}},
		NonceValue{Placeholder: "blah"},
	}) {
		t.Errorf("CSP ScriptSrcElem: got = %v (%T), want multiple", csp.ScriptSrcElem, csp.ScriptSrcElem)
	}
	if !reflect.DeepEqual(csp.StyleSrc, UnsafeEvalDirectiveValue{}) {
		t.Errorf("CSP StyleSrc: got = %v (%T), want 'unsafe-eval'", csp.StyleSrc, csp.StyleSrc)
	}
	if !reflect.DeepEqual(csp.StyleSrcElem, MultiFetchDirectiveValue{
		WasmUnsafeEvalDirectiveValue{},
		UnsafeInlineDirectiveValue{},
		UnsafeHashesDirectiveValue{},
		InlineSpeculationRulesDirectiveValue{},
		StrictDynamicDirectiveValue{},
		ReportSampleDirectiveValue{},
	}) {
		t.Errorf("CSP StyleSrcElem: got = %v (%T), want multiple", csp.StyleSrcElem, csp.StyleSrcElem)
	}

	if !reflect.DeepEqual(csp.Sandbox, SandboxWithAllowed{
		SandboxAllow("allow-scripts"),
	}) {
		t.Errorf("CSP Sandbox: got = %v (%T), want multiple", csp.StyleSrcElem, csp.StyleSrcElem)
	}

	if !reflect.DeepEqual(csp.FormAction, SourceExpressionList{
		SelfDirectiveValue{},
		HostSourceDirectiveValue{Host: "a.b.com"},
		HostSourceDirectiveValue{Host: "my.test.com"},
		SchemeSourceDirectiveValue("http"),
	}) {
		t.Errorf("CSP FormAction: got = %v (%T), want multiple", csp.FormAction, csp.FormAction)
	}

	if !reflect.DeepEqual(csp.FrameAncestors, NoneDirectiveValue{}) {
		t.Errorf("CSP FrameAncestors: got = %v (%T), want 'none'", csp.FrameAncestors, csp.FrameAncestors)
	}
}

func TestSliceToFetchDirectiveValueSliceHookFunc(t *testing.T) {
	i := reflect.TypeOf((*FetchDirectiveValue)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "ValidInputSlice",
			from: []any{"self", "'none'"},
			want: MultiFetchDirectiveValue{SelfDirectiveValue{}, NoneDirectiveValue{}},
		},
		{
			name:    "InputSliceWithUnsupportedValue/Integer",
			from:    []any{123},
			wantErr: true,
		},
		{
			name:    "InputSliceWithUnsupportedValue/UnsupportedString",
			from:    []any{"test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := SliceToFetchDirectiveValueSliceHookFunc()
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

func TestSliceToNoneOrSourceExpressionListHookFunc(t *testing.T) {
	i := reflect.TypeOf((*NoneOrSourceExpressionList)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "ValidInputSlice",
			from: []any{"self", map[string]any{"host": "foo.bar.com"}},
			want: SourceExpressionList{SelfDirectiveValue{}, HostSourceDirectiveValue{Host: "foo.bar.com"}},
		},
		{
			name:    "InputSliceWithUnsupportedValue/Integer",
			from:    []any{123},
			wantErr: true,
		},
		{
			name:    "InputSliceWithUnsupportedValue/UnsupportedString",
			from:    []any{"test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := SliceToNoneOrSourceExpressionListHookFunc()
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

func TestSliceToSandboxWithAllowedHookFunc(t *testing.T) {
	i := reflect.TypeOf((*SandboxDirectiveValue)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "ValidInputSlice",
			from: []any{"allow-downloads", "allow-forms"},
			want: SandboxWithAllowed{SandboxAllow("allow-downloads"), SandboxAllow("allow-forms")},
		},
		{
			name:    "InputSliceWithUnsupportedValue/Integer",
			from:    []any{123},
			wantErr: true,
		},
		{
			name:    "InputSliceWithUnsupportedValue/UnsupportedString",
			from:    []any{"test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := SliceToSandboxWithAllowedHookFunc()
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

func TestStringToFetchDirectiveValueHookFunc(t *testing.T) {
	i := reflect.TypeOf((*FetchDirectiveValue)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name:    "Invalid",
			from:    "foo",
			wantErr: true,
		},
		{
			name: "Valid/None/none",
			from: "none",
			want: NoneDirectiveValue{},
		},
		{
			name: "Valid/None/None",
			from: "None",
			want: NoneDirectiveValue{},
		},
		{
			name: "Valid/None/'none'",
			from: "'none'",
			want: NoneDirectiveValue{},
		},
		{
			name: "Valid/Self/self",
			from: "self",
			want: SelfDirectiveValue{},
		},
		{
			name: "Valid/Self/Self",
			from: "Self",
			want: SelfDirectiveValue{},
		},
		{
			name: "Valid/Self/'self'",
			from: "'self'",
			want: SelfDirectiveValue{},
		},
		{
			name: "Valid/UnsafeEval/unsafe-eval",
			from: "unsafe-eval",
			want: UnsafeEvalDirectiveValue{},
		},
		{
			name: "Valid/UnsafeEval/UnsafeEval",
			from: "UnsafeEval",
			want: UnsafeEvalDirectiveValue{},
		},
		{
			name: "Valid/UnsafeEval/'unsafe-eval'",
			from: "'unsafe-eval'",
			want: UnsafeEvalDirectiveValue{},
		},
		{
			name: "Valid/WasmUnsafeEval/wasm-unsafe-eval",
			from: "wasm-unsafe-eval",
			want: WasmUnsafeEvalDirectiveValue{},
		},
		{
			name: "Valid/WasmUnsafeEval/WasmUnsafeEval",
			from: "WasmUnsafeEval",
			want: WasmUnsafeEvalDirectiveValue{},
		},
		{
			name: "Valid/WasmUnsafeEval/'wasm-unsafe-eval'",
			from: "'wasm-unsafe-eval'",
			want: WasmUnsafeEvalDirectiveValue{},
		},
		{
			name: "Valid/UnsafeInline/unsafe-inline",
			from: "unsafe-inline",
			want: UnsafeInlineDirectiveValue{},
		},
		{
			name: "Valid/UnsafeInline/UnsafeInline",
			from: "UnsafeInline",
			want: UnsafeInlineDirectiveValue{},
		},
		{
			name: "Valid/UnsafeInline/'unsafe-inline'",
			from: "'unsafe-inline'",
			want: UnsafeInlineDirectiveValue{},
		},
		{
			name: "Valid/UnsafeHashes/unsafe-hashes",
			from: "unsafe-hashes",
			want: UnsafeHashesDirectiveValue{},
		},
		{
			name: "Valid/UnsafeHashes/UnsafeHashes",
			from: "UnsafeHashes",
			want: UnsafeHashesDirectiveValue{},
		},
		{
			name: "Valid/UnsafeHashes/'unsafe-hashes'",
			from: "'unsafe-hashes'",
			want: UnsafeHashesDirectiveValue{},
		},
		{
			name: "Valid/InlineSpeculationRules/inline-speculation-rules",
			from: "inline-speculation-rules",
			want: InlineSpeculationRulesDirectiveValue{},
		},
		{
			name: "Valid/InlineSpeculationRules/InlineSpeculationRules",
			from: "InlineSpeculationRules",
			want: InlineSpeculationRulesDirectiveValue{},
		},
		{
			name: "Valid/InlineSpeculationRules/'inline-speculation-rules'",
			from: "'inline-speculation-rules'",
			want: InlineSpeculationRulesDirectiveValue{},
		},
		{
			name: "Valid/StrictDynamic/strict-dynamic",
			from: "strict-dynamic",
			want: StrictDynamicDirectiveValue{},
		},
		{
			name: "Valid/StrictDynamic/StrictDynamic",
			from: "StrictDynamic",
			want: StrictDynamicDirectiveValue{},
		},
		{
			name: "Valid/StrictDynamic/'strict-dynamic'",
			from: "'strict-dynamic'",
			want: StrictDynamicDirectiveValue{},
		},
		{
			name: "Valid/ReportSample/report-sample",
			from: "report-sample",
			want: ReportSampleDirectiveValue{},
		},
		{
			name: "Valid/ReportSample/ReportSample",
			from: "ReportSample",
			want: ReportSampleDirectiveValue{},
		},
		{
			name: "Valid/ReportSample/'report-sample'",
			from: "'report-sample'",
			want: ReportSampleDirectiveValue{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := StringToFetchDirectiveValueHookFunc()
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

func TestStringToNoneOrSourceExpressionListHookFunc(t *testing.T) {
	i := reflect.TypeOf((*SourceExpressionListItem)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name:    "Invalid",
			from:    "foo",
			wantErr: true,
		},
		{
			name: "Valid/None/none",
			from: "none",
			want: NoneDirectiveValue{},
		},
		{
			name: "Valid/None/None",
			from: "None",
			want: NoneDirectiveValue{},
		},
		{
			name: "Valid/None/'none'",
			from: "'none'",
			want: NoneDirectiveValue{},
		},
		{
			name: "Valid/Self/self",
			from: "self",
			want: SelfDirectiveValue{},
		},
		{
			name: "Valid/Self/Self",
			from: "Self",
			want: SelfDirectiveValue{},
		},
		{
			name: "Valid/Self/'self'",
			from: "'self'",
			want: SelfDirectiveValue{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := StringToNoneOrSourceExpressionListHookFunc()
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

func TestStringToSandboxDirectiveValueHookFunc(t *testing.T) {
	i := reflect.TypeOf((*SandboxDirectiveValue)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "All",
			from: "all",
			want: SandboxAll{},
		},
		{
			name:    "Invalid",
			from:    "some",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := StringToSandboxDirectiveValueHookFunc()
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

func TestStringToSandboxAllowHookFunc(t *testing.T) {
	i := reflect.TypeOf(SandboxAllow(""))
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Valid/allow-downloads",
			from: "allow-downloads",
			want: SandboxAllow("allow-downloads"),
		},
		{
			name: "Valid/allow-forms",
			from: "allow-forms",
			want: SandboxAllow("allow-forms"),
		},
		{
			name: "Valid/allow-modals",
			from: "allow-modals",
			want: SandboxAllow("allow-modals"),
		},
		{
			name: "Valid/allow-orientation-lock",
			from: "allow-orientation-lock",
			want: SandboxAllow("allow-orientation-lock"),
		},
		{
			name: "Valid/allow-pointer-lock",
			from: "allow-pointer-lock",
			want: SandboxAllow("allow-pointer-lock"),
		},
		{
			name: "Valid/allow-popups",
			from: "allow-popups",
			want: SandboxAllow("allow-popups"),
		},
		{
			name: "Valid/allow-popups-to-escape-sandbox",
			from: "allow-popups-to-escape-sandbox",
			want: SandboxAllow("allow-popups-to-escape-sandbox"),
		},
		{
			name: "Valid/allow-presentation",
			from: "allow-presentation",
			want: SandboxAllow("allow-presentation"),
		},
		{
			name: "Valid/allow-same-origin",
			from: "allow-same-origin",
			want: SandboxAllow("allow-same-origin"),
		},
		{
			name: "Valid/allow-scripts",
			from: "allow-scripts",
			want: SandboxAllow("allow-scripts"),
		},
		{
			name: "Valid/allow-storage-access-by-user-activation",
			from: "allow-storage-access-by-user-activation",
			want: SandboxAllow("allow-storage-access-by-user-activation"),
		},
		{
			name: "Valid/allow-top-navigation",
			from: "allow-top-navigation",
			want: SandboxAllow("allow-top-navigation"),
		},
		{
			name: "Valid/allow-top-navigation-by-user-activation",
			from: "allow-top-navigation-by-user-activation",
			want: SandboxAllow("allow-top-navigation-by-user-activation"),
		},
		{
			name: "Valid/allow-top-navigation-to-custom-protocols",
			from: "allow-top-navigation-to-custom-protocols",
			want: SandboxAllow("allow-top-navigation-to-custom-protocols"),
		},
		{
			name:    "Invalid",
			from:    "something-else",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := StringToSandboxAllowHookFunc()
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

func TestMapToFetchDirectiveValueHookFunc(t *testing.T) {
	i := reflect.TypeOf((*FetchDirectiveValue)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Hash/Valid",
			from: map[string]any{"hash": map[string]any{"alg": "the-alg", "hash": "AAECAw=="}},
			want: SubResourceHashDirectiveValue{Alg: "the-alg", Hash: []byte{0, 1, 2, 3}},
		},
		{
			name:    "Hash/Invalid/Integer",
			from:    map[string]any{"hash": 111},
			wantErr: true,
		},
		{
			name:    "Hash/Invalid/IntegerInsteadOfBase64Hash",
			from:    map[string]any{"hash": map[string]any{"hash": 123}},
			wantErr: true,
		},
		{
			name: "Host/Valid/HostnameString",
			from: map[string]any{"host": "unit.test.com"},
			want: HostSourceDirectiveValue{Host: "unit.test.com"},
		},
		{
			name: "Host/Valid/Map",
			from: map[string]any{"host": map[string]any{"host": "foo.com", "scheme": "https"}},
			want: HostSourceDirectiveValue{Host: "foo.com", Scheme: toPtr("https")},
		},
		{
			name:    "Host/Invalid/Integer",
			from:    map[string]any{"host": 111},
			wantErr: true,
		},
		{
			name:    "Host/Invalid/IntegerInsteadOfHostname",
			from:    map[string]any{"host": map[string]any{"host": 123}},
			wantErr: true,
		},
		{
			name: "Nonce/Valid/PlaceholderString",
			from: map[string]any{"nonce": "test-value"},
			want: NonceValue{Placeholder: "test-value"},
		},
		{
			name: "Nonce/Valid/Map",
			from: map[string]any{"nonce": map[string]any{"placeholder": "replace-me"}},
			want: NonceValue{Placeholder: "replace-me"},
		},
		{
			name:    "Nonce/Invalid/Integer",
			from:    map[string]any{"nonce": 111},
			wantErr: true,
		},
		{
			name:    "Nonce/Invalid/IntegerInsteadOfHostname",
			from:    map[string]any{"nonce": map[string]any{"placeholder": 123}},
			wantErr: true,
		},
		{
			name: "Scheme/Valid/String",
			from: map[string]any{"scheme": "ftp"},
			want: SchemeSourceDirectiveValue("ftp"),
		},
		{
			name:    "Scheme/Invalid/Integer",
			from:    map[string]any{"scheme": 111},
			wantErr: true,
		},
		{
			name:    "Scheme/Invalid/Boolean",
			from:    map[string]any{"scheme": true},
			wantErr: true,
		},
		{
			name:    "Unsupported",
			from:    map[string]any{"other": true},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := MapToFetchDirectiveValueHookFunc()
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

func TestMapToSourceExpressionListItemHookFunc(t *testing.T) {
	i := reflect.TypeOf((*SourceExpressionListItem)(nil)).Elem()
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Host/Valid/HostnameString",
			from: map[string]any{"host": "unit.test.com"},
			want: HostSourceDirectiveValue{Host: "unit.test.com"},
		},
		{
			name: "Host/Valid/Map",
			from: map[string]any{"host": map[string]any{"host": "foo.com", "scheme": "https"}},
			want: HostSourceDirectiveValue{Host: "foo.com", Scheme: toPtr("https")},
		},
		{
			name:    "Host/Invalid/Integer",
			from:    map[string]any{"host": 111},
			wantErr: true,
		},
		{
			name:    "Host/Invalid/IntegerInsteadOfHostname",
			from:    map[string]any{"host": map[string]any{"host": 123}},
			wantErr: true,
		},
		{
			name: "Scheme/Valid/String",
			from: map[string]any{"scheme": "ftp"},
			want: SchemeSourceDirectiveValue("ftp"),
		},
		{
			name:    "Scheme/Invalid/Integer",
			from:    map[string]any{"scheme": 111},
			wantErr: true,
		},
		{
			name:    "Scheme/Invalid/Boolean",
			from:    map[string]any{"scheme": true},
			wantErr: true,
		},
		{
			name:    "Unsupported",
			from:    map[string]any{"other": true},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := MapToSourceExpressionListItemHookFunc()
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

func TestBase64StringToByteSliceHookFunc(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		from    any
		want    any
		wantErr bool
	}{
		{
			name: "Valid",
			from: "3q2+7w==",
			want: []byte{0xde, 0xad, 0xbe, 0xef},
		},
		{
			name:    "Invalid/MissingPadding",
			from:    "3q2+7w",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Base64StringToByteSliceHookFunc()
			hook := Base64StringToByteSliceHookFunc()
			from := reflect.ValueOf(tt.from)
			to := reflect.ValueOf([]byte{})
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
