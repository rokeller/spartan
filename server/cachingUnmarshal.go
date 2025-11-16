package server

import (
	"github.com/go-viper/mapstructure/v2"
)

func CachingDecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		MapToRouteMatcherHookFunc(),
	)
}
