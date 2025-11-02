package cmd

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/rokeller/spartan/server"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

func getConfig() (*server.Config, error) {
	var config server.Config
	hookOpt := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		server.CspDecodeHook(),
		server.ReferrerPolicyDecodeHook(),
	))
	if err := viper.Unmarshal(&config, hookOpt); err != nil {
		return nil, err
	} else {
		klog.InfoS("Server configuration", "config", config)
	}
	return &config, nil
}
