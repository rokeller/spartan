package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func withSecurityMiddleware(c SecurityConfig, next http.Handler) http.Handler {
	csp := c.ContentSecurityPolicy
	if nil == csp {
		csp = &DefaultContentSecurityPolicy
	}

	if nil == c.ContentTypeOptionsNoSniff {
		c.ContentTypeOptionsNoSniff = new(bool)
		*c.ContentTypeOptionsNoSniff = true
	}

	if nil == c.StrictTransportSecurityPolicy {
		c.StrictTransportSecurityPolicy = &DefaultStrictTransportSecurityPolicy
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("strict-transport-security", c.StrictTransportSecurityPolicy.StsHeaderValue())

		if *c.ContentTypeOptionsNoSniff {
			w.Header().Add("x-content-type-options", "nosniff")
		}

		if shouldAddContentSecurityPolicyHeader(r) {
			v := csp.HeaderValue(r.Context())
			if csp.ReportOnly {
				w.Header().Add("content-security-policy-report-only", v)
			} else {
				w.Header().Add("content-security-policy", v)
			}
		}

		next.ServeHTTP(w, r)
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

func (p StrictTransportSecurityPolicy) StsHeaderValue() string {
	directives := []string{
		fmt.Sprintf("max-age=%d",
			int64(p.MaxAge.Truncate(time.Second).Seconds())),
	}

	if p.IncludeSubDomains {
		directives = append(directives, "includeSubDomains")
	}

	return strings.Join(directives, "; ")
}
