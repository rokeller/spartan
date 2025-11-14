package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/rokeller/spartan/server"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

const (
	EnvPrefix            = "SPARTAN"
	EnvScopeSeparator    = "__"
	ConfigScopeSeparator = "."
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
	Args: cobra.NoArgs,

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
	vpr.BindPFlag("server.port", rootCmd.Flags().Lookup("port"))

	rootCmd.PersistentFlags().StringP("static-content-dir", "d", "/content",
		"The path to the directory holding the static content to serve.")
	vpr.BindPFlag("server.staticContentDir", rootCmd.PersistentFlags().Lookup("static-content-dir"))

	rootCmd.Flags().StringP("server-path-root", "r", "",
		"The absolute path on the server where the static content is exposed.")
	vpr.BindPFlag("server.pathRoot", rootCmd.Flags().Lookup("server-path-root"))
}

// initConfig reads in config file and ENV variables, if set.
func initConfig() {
	// Give prefix for environment variables and read them automatically.
	vpr.SetEnvPrefix(EnvPrefix)
	vpr.SetEnvKeyReplacer(strings.NewReplacer(ConfigScopeSeparator, EnvScopeSeparator))
	vpr.AutomaticEnv()
	// Bind environment variables matching the prefix because viper doesn't do
	// this by itself (see https://github.com/spf13/viper/issues/761).
	bindEnv()
	vpr.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		vpr.SetConfigFile(cfgFile)
	} else {
		// Find current working directory.
		cwd, err := os.Getwd()
		if nil != err {
			klog.Error(err, "Failed to get current working directory")
			os.Exit(1)
		}

		vpr.AddConfigPath(cwd)
		vpr.SetConfigName("config")
	}

	// If a config file is found, read it in.
	if err := vpr.ReadInConfig(); nil == err {
		klog.V(1).InfoS("Loaded config file", "file", vpr.ConfigFileUsed())
	} else {
		klog.ErrorS(err, "Failed to load config file", "file", vpr.ConfigFileUsed())
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	if config, err := getConfig(); nil != err {
		return err
	} else {
		return server.Serve(cmd.Context(), config.Server)
	}
}

func bindEnv() {
	replacer := strings.NewReplacer(
		EnvPrefix+"_" /* added automatically by viper */, "",
		EnvScopeSeparator, ConfigScopeSeparator,
	)
	for _, e := range os.Environ() {
		k, _, ok := strings.Cut(e, "=")
		if !ok || !strings.HasPrefix(k, EnvPrefix+"_") {
			continue
		}

		k = replacer.Replace(k)
		vpr.BindEnv(k)
	}
}
