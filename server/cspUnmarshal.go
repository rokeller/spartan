package server

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

func CspDecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		SliceToFetchDirectiveValueSliceHookFunc(),
		SliceToNoneOrSourceExpressionListHookFunc(),
		SliceToSandboxWithAllowedHookFunc(),

		StringToFetchDirectiveValueHookFunc(),
		StringToNoneOrSourceExpressionListHookFunc(),
		StringToSandboxDirectiveValueHookFunc(),
		StringToSandboxAllowHookFunc(),

		MapToFetchDirectiveValueHookFunc(),
		MapToSourceExpressionListItemHookFunc(),
	)
}

func SliceToFetchDirectiveValueSliceHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*FetchDirectiveValue)(nil)).Elem()

	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		decodeHook := CspDecodeHook()
		if f.Kind() != reflect.Slice {
			return data, nil
		}
		if t.Kind() != reflect.Interface || !t.Implements(i) {
			return data, nil
		}

		s := data.([]any)
		res := []FetchDirectiveValue{}
		for _, el := range s {
			f := reflect.ValueOf(el)
			t := reflect.New(i).Elem()
			item, err := mapstructure.DecodeHookExec(decodeHook, f, t)
			if nil != err {
				return nil, err
			}
			if i, ok := item.(FetchDirectiveValue); !ok {
				return nil, fmt.Errorf("%v (%T) is not allowed in a fetch directive value list", el, el)
			} else {
				res = append(res, i)
			}
		}
		return MultiFetchDirectiveValue(res), nil
	}
}

func SliceToNoneOrSourceExpressionListHookFunc() mapstructure.DecodeHookFunc {
	l := reflect.TypeOf((*NoneOrSourceExpressionList)(nil)).Elem()
	i := reflect.TypeOf((*SourceExpressionListItem)(nil)).Elem()

	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		decodeHook := CspDecodeHook()
		if f.Kind() != reflect.Slice {
			return data, nil
		}
		if t.Kind() != reflect.Interface || !t.Implements(l) {
			return data, nil
		}

		s := data.([]any)
		res := []SourceExpressionListItem{}
		for _, el := range s {
			f := reflect.ValueOf(el)
			t := reflect.New(i).Elem()
			item, err := mapstructure.DecodeHookExec(decodeHook, f, t)
			if nil != err {
				return nil, err
			}
			if i, ok := item.(SourceExpressionListItem); !ok {
				return nil, fmt.Errorf("%v (%T) is not allowed in a source expression list.", el, el)
			} else {
				res = append(res, i)
			}
		}
		return SourceExpressionList(res), nil
	}
}

func SliceToSandboxWithAllowedHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*SandboxDirectiveValue)(nil)).Elem()

	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		decodeHook := CspDecodeHook()
		if f.Kind() != reflect.Slice {
			return data, nil
		}
		if t.Kind() != reflect.Interface || !t.Implements(i) {
			return data, nil
		}

		s := data.([]any)
		res := []SandboxAllow{}
		for _, el := range s {
			f := reflect.ValueOf(el)
			t := reflect.ValueOf(SandboxAllow(""))
			item, err := mapstructure.DecodeHookExec(decodeHook, f, t)
			if nil != err {
				return nil, err
			}
			if i, ok := item.(SandboxAllow); !ok {
				return nil, fmt.Errorf("%v (%T) is not allowed in a sandbox.", el, el)
			} else {
				res = append(res, i)
			}
		}
		return SandboxWithAllowed(res), nil
	}
}

func StringToFetchDirectiveValueHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*FetchDirectiveValue)(nil)).Elem()
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != i {
			return data, nil
		}

		str := data.(string)
		switch str {
		case "none", "None", "'none'":
			return NoneDirectiveValue{}, nil

		case "self", "Self", "'self'":
			return SelfDirectiveValue{}, nil

		case "unsafe-eval", "UnsafeEval", "'unsafe-eval'":
			return UnsafeEvalDirectiveValue{}, nil

		case "wasm-unsafe-eval", "WasmUnsafeEval", "'wasm-unsafe-eval'":
			return WasmUnsafeEvalDirectiveValue{}, nil

		case "unsafe-inline", "UnsafeInline", "'unsafe-inline'":
			return UnsafeInlineDirectiveValue{}, nil

		case "unsafe-hashes", "UnsafeHashes", "'unsafe-hashes'":
			return UnsafeHashesDirectiveValue{}, nil

		case "inline-speculation-rules", "InlineSpeculationRules", "'inline-speculation-rules'":
			return InlineSpeculationRulesDirectiveValue{}, nil

		case "strict-dynamic", "StrictDynamic", "'strict-dynamic'":
			return StrictDynamicDirectiveValue{}, nil

		case "report-sample", "ReportSample", "'report-sample'":
			return ReportSampleDirectiveValue{}, nil
		}

		return nil, fmt.Errorf("%v (%T) is not a valid directive value.", data, data)
	}
}

func StringToNoneOrSourceExpressionListHookFunc() mapstructure.DecodeHookFunc {
	l := reflect.TypeOf((*NoneOrSourceExpressionList)(nil)).Elem()
	i := reflect.TypeOf((*SourceExpressionListItem)(nil)).Elem()
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != l && t != i {
			return data, nil
		}

		str := data.(string)
		switch str {
		case "none", "None", "'none'":
			return NoneDirectiveValue{}, nil

		case "self", "Self", "'self'":
			return SelfDirectiveValue{}, nil
		}

		return nil, fmt.Errorf("%v (%T) is not a valid directive value.", data, data)
	}
}

func StringToSandboxDirectiveValueHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*SandboxDirectiveValue)(nil)).Elem()
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != i {
			return data, nil
		}

		str := data.(string)
		switch str {
		case "all":
			return SandboxAll{}, nil
		}

		return nil, fmt.Errorf("%v (%T) is not a valid sandbox value.", data, data)
	}
}

func StringToSandboxAllowHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf(SandboxAllow(""))
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != i {
			return data, nil
		}

		str := data.(string)
		switch str {
		case "allow-downloads", "allow-forms", "allow-modals",
			"allow-orientation-lock", "allow-pointer-lock", "allow-popups",
			"allow-popups-to-escape-sandbox", "allow-presentation",
			"allow-same-origin", "allow-scripts",
			"allow-storage-access-by-user-activation", "allow-top-navigation",
			"allow-top-navigation-by-user-activation",
			"allow-top-navigation-to-custom-protocols":
			return SandboxAllow(str), nil
		}

		return nil, fmt.Errorf("%v (%T) is not a valid sandbox allow value.", data, data)
	}
}

func MapToFetchDirectiveValueHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*FetchDirectiveValue)(nil)).Elem()
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}
		if t != i {
			return data, nil
		}

		m := data.(map[string]any)
		if hash, ok := m["hash"]; ok {
			// Always a map with `alg` and `hash` properties.
			var sri SubResourceHashDirectiveValue
			config := &mapstructure.DecoderConfig{
				Metadata:   nil,
				DecodeHook: Base64StringToByteSliceHookFunc(),
				Result:     &sri,
			}
			if d, err := mapstructure.NewDecoder(config); nil != err {
				return nil, err
			} else if err := d.Decode(hash); nil != err {
				return nil, err
			}
			return sri, nil
		} else if host, ok := m["host"]; ok {
			// When set with just a string value, interpret that as the host.
			if v, ok := host.(string); ok {
				return HostSourceDirectiveValue{Host: v}, nil
			}
			// Otherwise, must be a map with the `scheme`, `port`, `path`
			// properties optional and the `host` property mandatory.
			var h HostSourceDirectiveValue
			if err := mapstructure.Decode(host, &h); nil != err {
				return nil, err
			}
			return h, nil
		} else if nonce, ok := m["nonce"]; ok {
			// When set with just a string value, interpret that as the placeholder
			// string in the resource.
			if v, ok := nonce.(string); ok {
				return NonceValue{Placeholder: v}, nil
			}
			// Otherwise, must be a map with the `placeholder` property set.
			var n NonceValue
			if err := mapstructure.Decode(nonce, &n); nil != err {
				return nil, err
			}
			return n, nil
		} else if scheme, ok := m["scheme"]; ok {
			// When set with just a string value, interpret that as the scheme.
			if v, ok := scheme.(string); ok {
				return SchemeSourceDirectiveValue(v), nil
			}
			return nil, fmt.Errorf("%v (%T) is not allowed in a scheme directive value", scheme, scheme)
		}

		return nil, fmt.Errorf("%v (%T) is not a valid fetch directive value.", data, data)
	}
}

func MapToSourceExpressionListItemHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*SourceExpressionListItem)(nil)).Elem()
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}
		if t != i {
			return data, nil
		}

		m := data.(map[string]any)
		if host, ok := m["host"]; ok {
			// When set with just a string value, interpret that as the host.
			if v, ok := host.(string); ok {
				return HostSourceDirectiveValue{Host: v}, nil
			}
			// Otherwise, must be a map with the `scheme`, `port`, `path`
			// properties optional and the `host` property mandatory.
			var h HostSourceDirectiveValue
			if err := mapstructure.Decode(host, &h); nil != err {
				return nil, err
			}
			return h, nil
		} else if scheme, ok := m["scheme"]; ok {
			// When set with just a string value, interpret that as the scheme.
			if v, ok := scheme.(string); ok {
				return SchemeSourceDirectiveValue(v), nil
			}
			return nil, fmt.Errorf("%v (%T) is not allowed in a scheme directive value", scheme, scheme)
		}

		return nil, fmt.Errorf("%v (%T) is not a valid source expression list item value.", data, data)
	}
}

func Base64StringToByteSliceHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Slice || t.Elem().Kind() != reflect.Uint8 {
			return data, nil
		}

		b64 := data.(string)
		if b, err := base64.StdEncoding.DecodeString(b64); nil != err {
			return nil, err
		} else {
			return b, nil
		}
	}
}
