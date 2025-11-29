//go:build !minimal

package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/rokeller/spartan/server"
	"github.com/spf13/cobra"
)

func Test_runPreviewHeaders(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		wantStdOut string
	}{
		{
			name: "Default",
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n\n" +
				"\x1b[1m\x1b[4mContent security headers\x1b[0m\n" +
				"Content-Security-Policy: default-src 'self'; object-src 'none'; base-uri 'none'; sandbox; form-action 'self'; frame-ancestors 'self'\n" +
				"X-Content-Type-Options: nosniff\n\n" +
				"\x1b[1m\x1b[4mPermissions policy headers\x1b[0m\n" +
				"Permissions-Policy: accelerometer=(), ambient-light-sensor=(), aria-notify=(), attribution-reporting=(), autoplay=(), bluetooth=(), browsing-topics=(), camera=(), captured-surface-control=(), compute-pressure=(), cross-origin-isolated=(), deferred-fetch=(), deferred-fetch-minimal=(), display-capture=(), encrypted-media=(), fullscreen=(), gamepad=(), geolocation=(), gyroscope=(), hid=(), identity-credential-get=(), idle-detection=(), language-detector=(), local-fonts=(), magnetometer=(), microphone=(), midi=(), on-device-speech-recognition=(), otp-credentials=(), payment=(), picture-in-picture=(), publickey-credentials-create=(), publickey-credentials-get=(), screen-wake-lock=(), serial=(), speaker-selection=(), storage-access=(), translator=(), summarizer=(), usb=(), web-share=(), window-management=(), xr-spatial-tracking=()\n\n" +
				"\x1b[1m\x1b[4mReferrer policy headers\x1b[0m\n" +
				"Referrer-Policy: same-origin\n\n" +
				"\x1b[1m\x1b[4mReporting endpoints headers\x1b[0m\n" +
				"\x1b[33mReporting endpoints not configured.\x1b[0m\n\n" +
				"\x1b[1m\x1b[4mTransport security headers\x1b[0m\n" +
				"Strict-Transport-Security: max-age=31536000; includeSubDomains\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			buf := &bytes.Buffer{}
			ctx := context.Background()
			cmd.SetOut(buf)
			cmd.SetContext(ctx)
			if err := runPreviewHeaders(cmd, nil); (err != nil) != tt.wantErr {
				t.Errorf("runPreviewHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
			stdout := buf.String()
			if stdout != tt.wantStdOut {
				t.Errorf("stdout mismatch: got = %q, want = %q", stdout, tt.wantStdOut)
			}
		})
	}
}

func Test_printCachePolicies(t *testing.T) {
	tests := []struct {
		name       string
		c          server.Cache
		wantStdOut string
	}{
		{
			name:       "None",
			c:          server.Cache{},
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n",
		},
		{
			name: "WithDefaultCachePolicy/NonEmpty",
			c: server.Cache{
				DefaultPolicy: &server.CachePolicy{
					Public: true,
					MaxAge: toPtr(time.Second * 20),
				},
			},
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n" +
				"\x1b[34m\x1b[3mDefault cache policy\x1b[0m\n" +
				"Cache-Control: max-age=20, public\n",
		},
		{
			name: "WithDefaultCachePolicy/Empty",
			c: server.Cache{
				DefaultPolicy: &server.CachePolicy{},
			},
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n" +
				"\x1b[34m\x1b[3mDefault cache policy\x1b[0m\n" +
				"\x1b[33mDefault cache policy is empty, no Cache-Control header generated/added.\x1b[0m\n",
		},
		{
			name: "WithRouteMatchedCachePolicies/NoRoutes",
			c: server.Cache{
				Routes: []server.RouteMatchingCachePolicy{},
			},
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n",
		},
		{
			name: "WithRouteMatchedCachePolicies/SingleRoute",
			c: server.Cache{
				Routes: []server.RouteMatchingCachePolicy{
					{
						Match: server.PathPrefixMatcher{PathPrefix: "/foo"},
						Policy: server.CachePolicy{
							Public: true,
						},
					},
				},
			},
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n" +
				"\x1b[34m\x1b[3mCache policy #0 for PathPrefix: /foo\x1b[0m\n" +
				"Cache-Control: public\n",
		},
		{
			name: "WithRouteMatchedCachePolicies/RouteWithEmptyPolicy",
			c: server.Cache{
				Routes: []server.RouteMatchingCachePolicy{
					{
						Match:  server.PathPrefixMatcher{PathPrefix: "/foo"},
						Policy: server.CachePolicy{},
					},
				},
			},
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n" +
				"\x1b[34m\x1b[3mCache policy #0 for PathPrefix: /foo\x1b[0m\n" +
				"\x1b[33mCache policy is empty, no Cache-Control header generated/added.\x1b[0m\n",
		},
		{
			name: "WithRouteMatchedCachePolicies/MultipleRoutes",
			c: server.Cache{
				Routes: []server.RouteMatchingCachePolicy{
					{
						Match: server.PathPrefixMatcher{PathPrefix: "/foo"},
						Policy: server.CachePolicy{
							Immutable: true,
						},
					},
					{
						Match: server.PathPrefixMatcher{PathPrefix: "/bar"},
						Policy: server.CachePolicy{
							Private: true,
						},
					},
				},
			},
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n" +
				"\x1b[34m\x1b[3mCache policy #0 for PathPrefix: /foo\x1b[0m\n" +
				"Cache-Control: immutable\n" +
				"\x1b[34m\x1b[3mCache policy #1 for PathPrefix: /bar\x1b[0m\n" +
				"Cache-Control: private\n",
		},
		{
			name: "WithRouteMatchedCachePolicies/RouteWithEmptyMatch",
			c: server.Cache{
				Routes: []server.RouteMatchingCachePolicy{
					{
						Policy: server.CachePolicy{Public: true},
					},
				},
			},
			wantStdOut: "\x1b[1m\x1b[4mCaching headers\x1b[0m\n" +
				"\x1b[33mCache policy route #0 has unknown/invalid match. Check configuration.\x1b[0m\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			buf := &bytes.Buffer{}
			ctx := context.Background()
			cmd.SetOut(buf)
			cmd.SetContext(ctx)
			printCachePolicies(cmd, tt.c)
			stdout := buf.String()
			if stdout != tt.wantStdOut {
				t.Errorf("stdout mismatch: got = %q, want = %q", stdout, tt.wantStdOut)
			}
		})
	}
}

func Test_printContentSecurityPolicy(t *testing.T) {
	type args struct {
		cmd     *cobra.Command
		p       *server.ContentSecurityPolicy
		noSniff bool
	}
	tests := []struct {
		name       string
		p          *server.ContentSecurityPolicy
		noSniff    bool
		wantStdOut string
		wantPrefix string
		wantSuffix string
	}{
		{
			name:    "EmptyCSP",
			p:       &server.ContentSecurityPolicy{},
			noSniff: true,
			wantStdOut: "\x1b[1m\x1b[4mContent security headers\x1b[0m\n" +
				"Content-Security-Policy: \n" +
				"X-Content-Type-Options: nosniff\n",
		},
		{
			name:    "BasicCSP",
			p:       &server.ContentSecurityPolicy{DefaultSrc: server.SelfDirectiveValue{}},
			noSniff: false,
			wantStdOut: "\x1b[1m\x1b[4mContent security headers\x1b[0m\n" +
				"Content-Security-Policy: default-src 'self'\n",
		},
		{
			name:    "ReportOnlyCSP",
			p:       &server.ContentSecurityPolicy{ReportOnly: true, ScriptSrc: server.NoneDirectiveValue{}},
			noSniff: false,
			wantStdOut: "\x1b[1m\x1b[4mContent security headers\x1b[0m\n" +
				"Content-Security-Policy-Report-Only: script-src 'none'\n",
		},
		{
			name:    "CSPWithNonce",
			p:       &server.ContentSecurityPolicy{StyleSrc: server.NonceValue{"PLACEHOLDER"}},
			noSniff: false,
			wantPrefix: "\x1b[1m\x1b[4mContent security headers\x1b[0m\n" +
				"\x1b[34m\x1b[3mNote: nonce values change every time.\x1b[0m\n" +
				"Content-Security-Policy: style-src 'nonce-",
			wantSuffix: "'\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			buf := &bytes.Buffer{}
			ctx := context.Background()
			cmd.SetOut(buf)
			cmd.SetContext(ctx)
			printContentSecurityPolicy(cmd, tt.p, tt.noSniff)
			stdout := buf.String()
			if tt.wantStdOut != "" && stdout != tt.wantStdOut {
				t.Errorf("stdout mismatch: got = %q, want = %q", stdout, tt.wantStdOut)
			}
			if tt.wantPrefix != "" && !strings.HasPrefix(stdout, tt.wantPrefix) {
				t.Errorf("stdout prefix mismatch: got = %q, want = %q", stdout, tt.wantPrefix)
			}
			if tt.wantSuffix != "" && !strings.HasSuffix(stdout, tt.wantSuffix) {
				t.Errorf("stdout suffix mismatch: got = %q, want = %q", stdout, tt.wantSuffix)
			}
		})
	}
}

func Test_printReportingEndpoints(t *testing.T) {
	tests := []struct {
		name       string
		e          *server.ReportingEndpoints
		wantStdOut string
	}{
		{
			name: "Empty",
			e:    &server.ReportingEndpoints{},
			wantStdOut: "\x1b[1m\x1b[4mReporting endpoints headers\x1b[0m\n" +
				"\x1b[33mReporting endpoints not configured.\x1b[0m\n",
		},
		{
			name: "SingleEndpoint",
			e: &server.ReportingEndpoints{
				server.ReportingEndpoint{Name: "test", Url: "https://foo.bar.com/reporting"},
			},
			wantStdOut: "\x1b[1m\x1b[4mReporting endpoints headers\x1b[0m\n" +
				"Reporting-Endpoints: test=\"https://foo.bar.com/reporting\"\n",
		},
		{
			name: "MultipleEndpoints",
			e: &server.ReportingEndpoints{
				server.ReportingEndpoint{Name: "one", Url: "https://one.com/reporting"},
				server.ReportingEndpoint{Name: "two", Url: "https://two.com/reporting"},
			},
			wantStdOut: "\x1b[1m\x1b[4mReporting endpoints headers\x1b[0m\n" +
				"Reporting-Endpoints: one=\"https://one.com/reporting\", two=\"https://two.com/reporting\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			buf := &bytes.Buffer{}
			ctx := context.Background()
			cmd.SetOut(buf)
			cmd.SetContext(ctx)
			printReportingEndpoints(cmd, tt.e)
			stdout := buf.String()
			if tt.wantStdOut != "" && stdout != tt.wantStdOut {
				t.Errorf("stdout mismatch: got = %q, want = %q", stdout, tt.wantStdOut)
			}
		})
	}
}
