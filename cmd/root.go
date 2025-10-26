package cmd

import (
	"context"

	"github.com/rokeller/spartan/server"
	"github.com/spf13/cobra"
)

var (
	version string

	port             uint16
	staticContentDir string
	serverRootPath   string
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
	rootCmd.Flags().SortFlags = false

	rootCmd.Flags().Uint16VarP(&port, "port", "p", 8080,
		"The local port to listen on for incoming requests.")

	rootCmd.Flags().StringVarP(&staticContentDir, "static-content-dir", "d", "",
		"The path to the directory holding the static content to serve.")
	rootCmd.MarkFlagRequired("static-content-dir")

	rootCmd.Flags().StringVarP(&serverRootPath, "server-root-path", "r", "",
		"The absolute path on the server where the static content is exposed.")
	rootCmd.MarkFlagRequired("server-root-path")
}

func runServer(cmd *cobra.Command, args []string) error {
	return server.Serve(cmd.Context(), port, staticContentDir, serverRootPath)
}
