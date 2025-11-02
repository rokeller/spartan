package server

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

func ReferrerPolicyDecodeHook() mapstructure.DecodeHookFunc {
	return StringToReferrerPolicyValueHookFunc()
}

func StringToReferrerPolicyValueHookFunc() mapstructure.DecodeHookFunc {
	p := reflect.TypeOf(ReferrerPolicy{})
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != p {
			return data, nil
		}

		str := data.(string)
		switch str {
		case "no-referrer", "no-referrer-when-downgrade", "origin",
			"origin-when-cross-origin", "same-origin", "strict-origin",
			"strict-origin-when-cross-origin", "unsafe-url":
			return ReferrerPolicy{referrerPolicyStringValue(str)}, nil
		}

		return nil, fmt.Errorf("%v (%T) is not a valid referrer policy value.", data, data)
	}
}
