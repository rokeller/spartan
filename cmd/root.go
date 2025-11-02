package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/rokeller/spartan/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	version string

	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spartan",
	Short: "A simple secure web server for SPA or similar apps with static assets.",
	Long: `A simple and secure-by-default web server for SPA or similar web app hosting
of static assets.`,

	Version: version,
	RunE:    runServer,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context, onError func(error)) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		onError(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().SortFlags = false

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "The config file to use.")

	rootCmd.Flags().Uint16P("port", "p", 8080,
		"The local port to listen on for incoming requests.")
	viper.BindPFlag("server.port", rootCmd.Flags().Lookup("port"))

	rootCmd.Flags().StringP("static-content-dir", "d", "",
		"The path to the directory holding the static content to serve.")
	viper.BindPFlag("server.staticContentDir", rootCmd.Flags().Lookup("static-content-dir"))

	rootCmd.Flags().StringP("server-path-root", "r", "",
		"The absolute path on the server where the static content is exposed.")
	viper.BindPFlag("server.pathRoot", rootCmd.Flags().Lookup("server-path-root"))

	viper.BindPFlags(rootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find current working directory.
		cwd, err := os.Getwd()
		if nil != err {
			klog.Error(err, "Failed to get current working directory")
			os.Exit(1)
		}

		viper.AddConfigPath(cwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName("spartan")
	}

	// Give prefix for environment variables and read them automatically.
	viper.SetEnvPrefix("SPARTAN_")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); nil == err {
		klog.V(1).InfoS("Loaded config file", "file", viper.ConfigFileUsed())
	} else {
		klog.ErrorS(err, "Failed to load config file", "file", viper.ConfigFileUsed())
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	var config server.Config
	hookOpt := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		server.CspDecodeHook(),
		server.ReferrerPolicyDecodeHook(),
	))
	if err := viper.Unmarshal(&config, hookOpt); err != nil {
		return err
	} else {
		klog.InfoS("Server configuration", "config", config)
	}

	return server.Serve(cmd.Context(), config.Server)
}
