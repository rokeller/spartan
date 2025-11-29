package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	executeWithExit(t, "TestMain", func(t *testing.T) {
		origArgs := os.Args
		defer func() { os.Args = origArgs }()

		os.Args = []string{"", "unsupported-command"}
		main()
	}, 1)
}

// executeWithExit helps testing cases where os.Exit(...) is called. For more
// information see: https://go.dev/talks/2014/testing.slide#23
func executeWithExit(
	t *testing.T,
	name string,
	fnWithExit func(*testing.T),
	wantExitCode int) {
	t.Helper()

	if os.Getenv("ACTUALLY_EXECUTE") == "1" {
		fnWithExit(t)
		return
	}

	testArgs := append(os.Args[1:], "-test.run=^"+name+"$")
	cmd := exec.Command(os.Args[0], testArgs...)
	cmd.Env = append(os.Environ(), "ACTUALLY_EXECUTE=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok {
		if e.ExitCode() != wantExitCode {
			t.Errorf("exit code mismatch: got = %d, want = %d", e.ExitCode(), wantExitCode)
		}
		return
	}

	if wantExitCode != 0 {
		t.Errorf("exit code mismatch: got = %d, want = %d", 0, wantExitCode)
	}
	if nil != err {
		t.Errorf("Execute error: got = %v, want nil", err)
	}
}
