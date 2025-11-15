package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

type routeMatchingCacheControlHeaderValue struct {
	RouteMatcher
	headerValue string
}
type routeMatching struct {
	routes       []routeMatchingCacheControlHeaderValue
	defaultValue string
}

func withCachingMiddleware(c Cache, next http.Handler) http.Handler {
	defaultPolicyHeaderValue := ""
	if nil != c.DefaultPolicy {
		v := c.DefaultPolicy.CacheControlHeaderValue()
		if v != "" {
			defaultPolicyHeaderValue = v
		}
	}

	routes := routeMatching{
		routes:       make([]routeMatchingCacheControlHeaderValue, len(c.RouteMatches)),
		defaultValue: defaultPolicyHeaderValue,
	}
	pos := 0
	for i, m := range c.RouteMatches {
		if nil == m.Match {
			klog.ErrorS(nil, "Misconfigured match for cache route match will be ignored; check configuration", "routeMatchIndex", i)
			continue
		}
		routes.routes[pos] = routeMatchingCacheControlHeaderValue{
			c.RouteMatches[i].Match,
			c.RouteMatches[i].Policy.CacheControlHeaderValue(),
		}
		pos++
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := routes.determinePolicyHeaderValue(r); v != "" {
			w.Header().Add("cache-control", v)
		}

		next.ServeHTTP(w, r)
	})
}

func (m routeMatching) determinePolicyHeaderValue(r *http.Request) string {
	res := m.defaultValue

	// Find out which matching route has the longest match and set the resulting
	// header value accordingly.
	if len(m.routes) > 0 {
		longest := -1
		for _, matcher := range m.routes {
			if l, ok := matcher.Match(r); !ok {
				continue
			} else if l > longest {
				res = matcher.headerValue
				longest = l
			}
		}
	}
	return res
}

func (p *CachePolicy) CacheControlHeaderValue() string {
	directives := []string{}
	if nil != p.MaxAge && time.Duration(p.MaxAge.Seconds()) > 0 {
		directives = append(directives, fmt.Sprintf("max-age=%d",
			int64(p.MaxAge.Truncate(time.Second).Seconds())))
	}
	if nil != p.SharedMaxAge && time.Duration(p.SharedMaxAge.Seconds()) > 0 {
		directives = append(directives, fmt.Sprintf("s-maxage=%d",
			int64(p.SharedMaxAge.Truncate(time.Second).Seconds())))
	}
	if nil != p.StaleIfError && time.Duration(p.StaleIfError.Seconds()) > 0 {
		directives = append(directives, fmt.Sprintf("stale-if-error=%d",
			int64(p.StaleIfError.Truncate(time.Second).Seconds())))
	}
	if nil != p.StaleWhileRevalidate && time.Duration(p.StaleWhileRevalidate.Seconds()) > 0 {
		directives = append(directives, fmt.Sprintf("stale-while-revalidate=%d",
			int64(p.StaleWhileRevalidate.Truncate(time.Second).Seconds())))
	}

	if p.Immutable {
		directives = append(directives, "immutable")
	}
	if p.MustRevalidate {
		directives = append(directives, "must-revalidate")
	}
	if p.MustUnderstand {
		directives = append(directives, "must-understand")
	}
	if p.NoCache {
		directives = append(directives, "no-cache")
	}
	if p.NoStore {
		directives = append(directives, "no-store")
	}
	if p.NoTransform {
		directives = append(directives, "no-transform")
	}
	if p.Private {
		directives = append(directives, "private")
	}
	if p.ProxyRevalidate {
		directives = append(directives, "proxy-revalidate")
	}
	if p.Public {
		directives = append(directives, "public")
	}

	if len(directives) > 0 {
		return strings.Join(directives, ", ")
	}

	return ""
}
