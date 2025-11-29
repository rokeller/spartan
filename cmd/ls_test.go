//go:build !minimal

package cmd

import (
	"bytes"
	"context"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func Test_runLs(t *testing.T) {
	repoRootDir := func() string {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("failed to get caller information")
		}
		dirname := path.Join(filepath.Dir(filename), "..")
		rootPath, err := filepath.Abs(dirname)
		if nil != err {
			panic(err)
		}
		return rootPath
	}()
	tests := []struct {
		name       string // description of this test case
		prepare    func(t *testing.T, cmd *cobra.Command)
		wantErr    bool
		wantStdOut string
	}{
		{
			name: "DirectoryNotFound",
			prepare: func(t *testing.T, cmd *cobra.Command) {
				vpr.Set("server.staticcontentdir", "/does/not/exist")
			},
			wantErr: true,
		},
		{
			name: "EmptyDir",
			prepare: func(t *testing.T, cmd *cobra.Command) {
				dir := path.Join(repoRootDir, "example/content")
				vpr.Set("server.staticcontentdir", dir)
			},
			wantStdOut: strings.ReplaceAll(
				"$ROOT/example/content\n"+
					"$ROOT/example/content/index.html\n"+
					"$ROOT/example/content/static\n"+
					"$ROOT/example/content/static/spartan.webp\n"+
					"$ROOT/example/content/static/styles.css\n",
				"$ROOT", repoRootDir),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			buf := &bytes.Buffer{}
			ctx := context.Background()
			cmd.SetOut(buf)
			cmd.SetContext(ctx)
			if nil != tt.prepare {
				tt.prepare(t, cmd)
			}
			gotErr := runLs(cmd, nil)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("runLs() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("runLs() succeeded unexpectedly")
			}
			stdout := buf.String()
			if stdout != tt.wantStdOut {
				t.Errorf("stdout mismatch: got = %q, want = %q", stdout, tt.wantStdOut)
			}
		})
	}
}
