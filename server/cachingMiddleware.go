package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func withCachingMiddleware(c CachePolicy, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := c.CacheControlHeaderValue()
		if v != "" {
			w.Header().Add("cache-control", v)
		}

		next.ServeHTTP(w, r)
	})
}

func (c CachePolicy) CacheControlHeaderValue() string {
	directives := []string{}
	if nil != c.MaxAge && time.Duration(c.MaxAge.Seconds()) > 0 {
		directives = append(directives, fmt.Sprintf("max-age=%d",
			int64(c.MaxAge.Truncate(time.Second).Seconds())))
	}
	if nil != c.SharedMaxAge && time.Duration(c.SharedMaxAge.Seconds()) > 0 {
		directives = append(directives, fmt.Sprintf("s-maxage=%d",
			int64(c.SharedMaxAge.Truncate(time.Second).Seconds())))
	}
	if nil != c.StaleIfError && time.Duration(c.StaleIfError.Seconds()) > 0 {
		directives = append(directives, fmt.Sprintf("stale-if-error=%d",
			int64(c.StaleIfError.Truncate(time.Second).Seconds())))
	}
	if nil != c.StaleWhileRevalidate && time.Duration(c.StaleWhileRevalidate.Seconds()) > 0 {
		directives = append(directives, fmt.Sprintf("stale-while-revalidate=%d",
			int64(c.StaleWhileRevalidate.Truncate(time.Second).Seconds())))
	}

	if c.Immutable {
		directives = append(directives, "immutable")
	}
	if c.MustRevalidate {
		directives = append(directives, "must-revalidate")
	}
	if c.MustUnderstand {
		directives = append(directives, "must-understand")
	}
	if c.NoCache {
		directives = append(directives, "no-cache")
	}
	if c.NoStore {
		directives = append(directives, "no-store")
	}
	if c.NoTransform {
		directives = append(directives, "no-transform")
	}
	if c.Private {
		directives = append(directives, "private")
	}
	if c.ProxyRevalidate {
		directives = append(directives, "proxy-revalidate")
	}
	if c.Public {
		directives = append(directives, "public")
	}

	if len(directives) > 0 {
		return strings.Join(directives, ", ")
	}

	return ""
}
