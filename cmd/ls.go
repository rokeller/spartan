//go:build !minimal

package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files in static content directory",
	Long: `Lists files recurively in the directory configured for static content.

This is a helper command that can be useful to inspect what files are present in
the container when debugging issues with static content.`,
	Args: cobra.NoArgs,
	RunE: runLs,
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func runLs(cmd *cobra.Command, args []string) error {
	if config, err := getConfig(); nil != err {
		return err
	} else {
		return listRecursive(cmd, config.Server.StaticContentDir)
	}
}

func listRecursive(cmd *cobra.Command, root string) error {
	w := cmd.OutOrStdout()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Fprintln(w, path)
		return nil
	})

	if nil != err {
		return err
	}
	return nil
}
