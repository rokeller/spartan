package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecute(t *testing.T) {
	var lastErr error
	origRunE := rootCmd.RunE
	tests := []struct {
		name    string // description of this test case
		prepare func(t *testing.T, ctx context.Context, cancel context.CancelFunc)
		wantErr bool
	}{
		{
			name: "CancelledContext",
			prepare: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) {
				cancel()
			},
		},
		{
			name: "InjectedError",
			prepare: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) {
				rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
					return errors.New("injected error")
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() { rootCmd.RunE = origRunE }()
			lastErr = nil
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if nil != tt.prepare {
				tt.prepare(t, ctx, cancel)
			}
			Execute(ctx, func(err error) {
				lastErr = err
			})
			if lastErr != nil {
				if !tt.wantErr {
					t.Errorf("Execute() failed: %v", lastErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Execute() succeeded unexpectedly")
			}
		})
	}
}
