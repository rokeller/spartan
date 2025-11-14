package cmd

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/rokeller/spartan/server"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	vpr = viper.New()
)

func getConfig() (*server.Config, error) {
	var config server.Config
	hookOpt := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		server.CspDecodeHook(),
		server.PermissionPolicyDecodeHook(),
		server.ReferrerPolicyDecodeHook(),
	))
	if err := vpr.Unmarshal(&config, hookOpt); err != nil {
		return nil, err
	} else {
		klog.InfoS("Loaded configuration", "config", config)
	}
	return &config, nil
}
