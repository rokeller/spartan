package server

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

func withSecurityMiddleware(c SecurityConfig, next http.Handler) http.Handler {
	csp := c.GetContentSecurityPolicy()
	pp := c.GetPermissionsPolicy()
	contentTypeOptionsNoSniff := c.GetContentTypeOptionsNoSniff()
	rp := c.GetReferrerPolicy()
	sts := c.GetStrictTransportSecurityPolicy()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if contentTypeOptionsNoSniff {
			w.Header().Add("x-content-type-options", "nosniff")
		}
		pp.AddToResponse(w)
		rp.AddToResponse(w)
		AddReportingEndpointsToResponse(c.ReportingEndpoints, w)
		sts.AddToResponse(w)

		var niw *nonceInjectingResponseWriter
		if shouldAddContentSecurityPolicyHeader(r) {
			nm := NonceMap{}
			r = r.WithContext(WithNonceMap(r.Context(), nm))
			v := csp.HeaderValue(r.Context())
			if csp.ReportOnly {
				w.Header().Add("content-security-policy-report-only", v)
			} else {
				w.Header().Add("content-security-policy", v)
			}

			if len(nm) > 0 {
				niw = &nonceInjectingResponseWriter{
					ResponseWriter: w,
					statusCode:     200,
					buf:            &bytes.Buffer{},
					nonces:         nm,
				}
			}
		}

		if nil != niw {
			next.ServeHTTP(niw, r)
			if n, err := niw.replaceNoncePlaceholdersAndWrite(); nil != err {
				klog.ErrorS(err, "Failed to replace nonce placeholders", "n", n)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func shouldAddContentSecurityPolicyHeader(r *http.Request) bool {
	secFetchDest := r.Header.Get("Sec-Fetch-Dest")
	switch secFetchDest {
	case "document", "script", "style",
		"fencedFrame", "frame", "iframe", "":
		return true
	}
	return false
}

type nonceInjectingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        *bytes.Buffer
	nonces     NonceMap
}

var _ http.ResponseWriter = &nonceInjectingResponseWriter{}

// Write implements http.ResponseWriter.
func (n *nonceInjectingResponseWriter) Write(b []byte) (int, error) {
	// Write to the buffer so we can replace nonce placeholders later.
	return n.buf.Write(b)
}

// WriteHeader implements http.ResponseWriter.
func (n *nonceInjectingResponseWriter) WriteHeader(statusCode int) {
	n.statusCode = statusCode
}

func (n *nonceInjectingResponseWriter) Unwrap() http.ResponseWriter {
	return n.ResponseWriter
}

func (n *nonceInjectingResponseWriter) replaceNoncePlaceholdersAndWrite() (int, error) {
	// For HTTP 304 (Not Modified) responses, we need to remove the CSP header
	// because otherwise the browser is instructed to trust only a new nonce
	// value, rather than the previously cached one.
	if n.statusCode == 304 {
		n.ResponseWriter.Header().Del("content-security-policy")
		n.ResponseWriter.Header().Del("content-security-policy-report-only")
		n.ResponseWriter.WriteHeader(n.statusCode)
		return 0, nil
	}

	oldNew := []string{}
	for key, val := range n.nonces {
		oldNew = append(oldNew, key, val)
	}

	r := strings.NewReplacer(oldNew...)
	replaced := []byte(r.Replace(n.buf.String()))
	n.ResponseWriter.Header().Set("content-length", strconv.FormatInt(int64(len(replaced)), 10))
	n.ResponseWriter.WriteHeader(n.statusCode)
	return n.ResponseWriter.Write(replaced)
}
