package server

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

func PermissionPolicyDecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		SliceToPermissionAllowValueHookFunc(),
		StringToPermissionAllowValueHookFunc(),
	)
}

func SliceToPermissionAllowValueHookFunc() mapstructure.DecodeHookFunc {
	v := reflect.TypeOf((*PermissionAllowValue)(nil)).Elem()
	i := reflect.TypeOf((*PermissionAllowListItem)(nil)).Elem()
	listItemDecodeHookFunc := StringToPermissionAllowListItemHookFunc()

	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.Slice {
			return data, nil
		}
		if t.Kind() != reflect.Interface || !t.Implements(v) {
			return data, nil
		}

		s := data.([]any)
		res := []PermissionAllowListItem{}
		for _, el := range s {
			f := reflect.ValueOf(el)
			t := reflect.New(i).Elem()
			item, err := mapstructure.DecodeHookExec(listItemDecodeHookFunc, f, t)
			if nil != err {
				return nil, err
			}
			if i, ok := item.(PermissionAllowListItem); !ok {
				return nil, fmt.Errorf("%v (%T) is not allowed in a permission policy allow list.", el, el)
			} else {
				res = append(res, i)
			}
		}
		return AllowMultiplePermission(res), nil
	}
}

func StringToPermissionAllowValueHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*PermissionAllowValue)(nil)).Elem()
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != i {
			return data, nil
		}

		str := data.(string)
		switch str {
		case "*", "all", "wildcard":
			return AllowWildcardPermission{}, nil
		case "none", "()":
			return AllowNonePermission{}, nil
		case "self":
			return AllowMultiplePermission{AllowSelfPermission{}}, nil
		case "src":
			return AllowMultiplePermission{AllowSrcPermission{}}, nil

		default:
			if o, err := parsePermissionPolicyOrigin(str); nil != err {
				return nil, err
			} else {
				return AllowMultiplePermission{o}, nil
			}
		}
	}
}

func StringToPermissionAllowListItemHookFunc() mapstructure.DecodeHookFunc {
	i := reflect.TypeOf((*PermissionAllowListItem)(nil)).Elem()
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Interface || !t.Implements(i) {
			return data, nil
		}

		str := data.(string)
		switch str {
		case "self":
			return AllowSelfPermission{}, nil
		case "src":
			return AllowSrcPermission{}, nil

		default:
			if o, err := parsePermissionPolicyOrigin(str); nil != err {
				return nil, err
			} else {
				return o, nil
			}
		}
	}
}

func parsePermissionPolicyOrigin(str string) (PermissionAllowListItem, error) {
	if u, err := url.Parse(str); nil != err {
		return nil, err
	} else {
		if (u.Scheme == "" || u.Host == "") ||
			(u.Fragment != "" || u.Path != "" || u.RawQuery != "" || u.User != nil) {
			return nil, fmt.Errorf("%q is not a valid origin. Origins must have a scheme and host, without any other URL components or trailing slash", str)
		}
		return AllowOriginPermission(u.String()), nil
	}
}
