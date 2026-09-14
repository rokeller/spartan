package test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	// Stores the full path to the cloned repository's root directory.
	RepoRootDir = repoRootDir()
)

func repoRootDir() string {
	return func() string {
		_, filename, _, ok := runtime.Caller(0) // the initializer of this module
		if !ok {
			panic("failed to get caller information")
		}
		// navigate 2 directories up (internal/test -> root of repo)
		dirname := path.Join(filepath.Dir(filename), "../..")
		rootPath, err := filepath.Abs(dirname)
		if nil != err {
			panic(err)
		}
		return rootPath
	}()
}

func RepoRelPath(t *testing.T, pathParts ...string) string {
	t.Helper()
	return path.Join(append([]string{RepoRootDir}, pathParts...)...)
}

func RepoFileContent(t *testing.T, pathParts ...string) []byte {
	t.Helper()
	p := RepoRelPath(t, pathParts...)
	data, err := os.ReadFile(p)
	if nil != err {
		t.Fatalf("failed to read repo file %q: %v", p, err)
	}
	return data
}

func RepoFileContentString(t *testing.T, pathParts ...string) string {
	t.Helper()
	return string(string(RepoFileContent(t, pathParts...)))
}
