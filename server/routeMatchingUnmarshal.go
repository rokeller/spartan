package server

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

func MapToRouteMatcherHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*RouteMatcher)(nil)).Elem()
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}
		if t.Kind() != reflect.Interface || !t.Implements(i) {
			return data, nil
		}

		m := data.(map[string]any)
		if _, ok := m["pathprefix"]; ok {
			var pp PathPrefixMatcher
			if err := mapstructure.Decode(data, &pp); nil != err {
				return nil, err
			}
			return pp, nil
		}

		return nil, fmt.Errorf("%v (%T) is not a valid route matcher value.", data, data)
	}
}
