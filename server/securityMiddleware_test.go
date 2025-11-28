package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func Test_withSecurityMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	tests := []struct {
		name                   string
		config                 SecurityConfig
		wantContentTypeOptions string
		wantCspReportOnly      string
		wantCsp                string
		wantCspPrefix          string
		wantPermissionsPolicy  string
		wantReferrerPolicy     string
		wantReportingEndpoints string
		wantBody               string
		wantRandomBody         bool
	}{
		{
			name:                   "Defaults",
			wantContentTypeOptions: "nosniff",
			wantCsp:                "default-src 'self'; object-src 'none'; base-uri 'none'; sandbox; form-action 'self'; frame-ancestors 'self'",
			wantPermissionsPolicy:  "accelerometer=(), ambient-light-sensor=(), aria-notify=(), attribution-reporting=(), autoplay=(), bluetooth=(), browsing-topics=(), camera=(), captured-surface-control=(), compute-pressure=(), cross-origin-isolated=(), deferred-fetch=(), deferred-fetch-minimal=(), display-capture=(), encrypted-media=(), fullscreen=(), gamepad=(), geolocation=(), gyroscope=(), hid=(), identity-credential-get=(), idle-detection=(), language-detector=(), local-fonts=(), magnetometer=(), microphone=(), midi=(), on-device-speech-recognition=(), otp-credentials=(), payment=(), picture-in-picture=(), publickey-credentials-create=(), publickey-credentials-get=(), screen-wake-lock=(), serial=(), speaker-selection=(), storage-access=(), translator=(), summarizer=(), usb=(), web-share=(), window-management=(), xr-spatial-tracking=()",
			wantReferrerPolicy:     "same-origin",
			wantBody:               "ok",
		},
		{
			name: "Csp/ReportOnly",
			config: SecurityConfig{
				ContentTypeOptionsNoSniff: toPtr(false),
				ContentSecurityPolicy: &ContentSecurityPolicy{
					ReportOnly: true,
					ImgSrc:     NoneDirectiveValue{},
				},
				PermissionsPolicy: &PermissionsPolicy{},
			},
			wantCspReportOnly:  "img-src 'none'",
			wantReferrerPolicy: "same-origin",
			wantBody:           "ok",
		},
		{
			name: "Csp/WithNonce",
			config: SecurityConfig{
				ContentTypeOptionsNoSniff: toPtr(false),
				ContentSecurityPolicy: &ContentSecurityPolicy{
					DefaultSrc: NonceValue{Placeholder: "ok"},
				},
				PermissionsPolicy: &PermissionsPolicy{},
			},
			wantCspPrefix:      "default-src 'nonce-",
			wantReferrerPolicy: "same-origin",
			wantBody:           "ok",
			wantRandomBody:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := withSecurityMiddleware(tt.config, handler)

			req := httptest.NewRequest("GET", "/endpoint", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			if w.Result().StatusCode != 200 {
				t.Errorf("server responded with status %d, want 200", w.Result().StatusCode)
			}

			got := w.Result().Header.Get("X-Content-Type-Options")
			if tt.wantContentTypeOptions != got {
				t.Errorf("X-Content-Type-Options mismatch: got = %s, want = %s", got, tt.wantContentTypeOptions)
			}

			got = w.Result().Header.Get("Content-Security-Policy-Report-Only")
			if tt.wantCspReportOnly != got {
				t.Errorf("Content-Security-Policy-Report-Only mismatch: got = %s, want = %s", got, tt.wantCspReportOnly)
			}
			got = w.Result().Header.Get("Content-Security-Policy")
			if tt.wantCsp != got && tt.wantCspPrefix == "" {
				t.Errorf("Content-Security-Policy mismatch: got = %s, want = %s", got, tt.wantCsp)
			} else if tt.wantCspPrefix != "" && !strings.HasPrefix(got, tt.wantCspPrefix) {
				t.Errorf("Content-Security-Policy mismatch: got = %s, want prefix = %s", got, tt.wantCspPrefix)
			}

			got = w.Result().Header.Get("Permissions-Policy")
			if tt.wantPermissionsPolicy != got {
				t.Errorf("Permissions-Policy mismatch: got = %s, want = %s", got, tt.wantPermissionsPolicy)
			}

			got = w.Result().Header.Get("Referrer-Policy")
			if tt.wantReferrerPolicy != got {
				t.Errorf("Referrer-Policy mismatch: got = %s, want = %s", got, tt.wantReferrerPolicy)
			}

			got = w.Result().Header.Get("Reporting-Endpoints")
			if tt.wantReportingEndpoints != got {
				t.Errorf("Reporting-Endpoints mismatch: got = %s, want = %s", got, tt.wantReportingEndpoints)
			}

			body, err := io.ReadAll(w.Result().Body)
			if nil != err {
				t.Errorf("failed to read recorded response body: %v", err)
			}
			if tt.wantBody != string(body) && !tt.wantRandomBody {
				t.Errorf("Body mismatch: got = %q, want = %q", body, tt.wantBody)
			} else if tt.wantRandomBody && tt.wantBody == string(body) {
				t.Errorf("Body mismatch: got = %q, want random", body)
			}
		})
	}
}

func Test_shouldAddContentSecurityPolicyHeader(t *testing.T) {
	tests := []struct {
		name       string
		addHeaders func(r *http.Request)
		want       bool
	}{
		{
			name: "Sec-Fetch-Dest/Missing",
			want: true,
		},
		{
			name:       "Sec-Fetch-Dest/document",
			addHeaders: func(r *http.Request) { r.Header.Add("Sec-fetch-dest", "document") },
			want:       true,
		},
		{
			name:       "Sec-Fetch-Dest/script",
			addHeaders: func(r *http.Request) { r.Header.Add("Sec-fetch-dest", "script") },
			want:       true,
		},
		{
			name:       "Sec-Fetch-Dest/style",
			addHeaders: func(r *http.Request) { r.Header.Add("Sec-fetch-dest", "style") },
			want:       true,
		},
		{
			name:       "Sec-Fetch-Dest/fencedFrame",
			addHeaders: func(r *http.Request) { r.Header.Add("Sec-fetch-dest", "fencedFrame") },
			want:       true,
		},
		{
			name:       "Sec-Fetch-Dest/frame",
			addHeaders: func(r *http.Request) { r.Header.Add("Sec-fetch-dest", "frame") },
			want:       true,
		},
		{
			name:       "Sec-Fetch-Dest/iframe",
			addHeaders: func(r *http.Request) { r.Header.Add("Sec-fetch-dest", "iframe") },
			want:       true,
		},
		{
			name:       "Sec-Fetch-Dest/empty",
			addHeaders: func(r *http.Request) { r.Header.Add("Sec-fetch-dest", "") },
			want:       true,
		},
		{
			name:       "Sec-Fetch-Dest/something",
			addHeaders: func(r *http.Request) { r.Header.Add("Sec-fetch-dest", "something") },
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/endpoint", nil)
			if nil != tt.addHeaders {
				tt.addHeaders(req)
			}

			if got := shouldAddContentSecurityPolicyHeader(req); got != tt.want {
				t.Errorf("shouldAddContentSecurityPolicyHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nonceInjectingResponseWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "Status/200",
			statusCode: 200,
		},
		{
			name:       "Status/304",
			statusCode: 304,
		},
		{
			name:       "Status/404",
			statusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			n := &nonceInjectingResponseWriter{
				ResponseWriter: w,
			}
			n.WriteHeader(tt.statusCode)

			if n.statusCode != tt.statusCode {
				t.Errorf("statusCode mismatch: got = %d, want = %d", n.statusCode, tt.statusCode)
			}
		})
	}
}

func Test_nonceInjectingResponseWriter_Unwrap(t *testing.T) {
	w := httptest.NewRecorder()
	n := &nonceInjectingResponseWriter{
		ResponseWriter: w,
	}
	if got := n.Unwrap(); !reflect.DeepEqual(got, w) {
		t.Errorf("nonceInjectingResponseWriter.Unwrap() = %v, want %v", got, w)
	}
}

func Test_nonceInjectingResponseWriter_replaceNoncePlaceholdersAndWrite(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		buf        *bytes.Buffer
		prepResp   func(w http.ResponseWriter)
		want       int
		wantErr    bool
		validate   func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:       "NoCspWithStatus304",
			statusCode: 304,
			prepResp: func(w http.ResponseWriter) {
				w.Header().Add("content-security-policy", "foo")
				w.Header().Add("content-security-policy-report-only", "foo")
			},
			want: 0,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				if "" != w.Result().Header.Get("content-security-policy") {
					t.Error("CSP header must be absent")
				}
				if "" != w.Result().Header.Get("content-security-policy-report-only") {
					t.Error("CSP (RO) header must be absent")
				}
			},
		},
		{
			name:       "ContentLengthPresentWithStatus200",
			statusCode: 200,
			buf:        bytes.NewBufferString("good"),
			want:       4,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				if val := w.Result().Header.Get("content-length"); val != "4" {
					t.Errorf("Content-Length header mismatch: got = %s, want 4", val)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if nil != tt.prepResp {
				tt.prepResp(w)
			}
			n := &nonceInjectingResponseWriter{
				ResponseWriter: w,
				statusCode:     tt.statusCode,
				buf:            tt.buf,
			}
			got, err := n.replaceNoncePlaceholdersAndWrite()
			if (err != nil) != tt.wantErr {
				t.Errorf("nonceInjectingResponseWriter.replaceNoncePlaceholdersAndWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nonceInjectingResponseWriter.replaceNoncePlaceholdersAndWrite() = %v, want %v", got, tt.want)
			}

			if nil != tt.validate {
				tt.validate(t, w)
			}
		})
	}
}
