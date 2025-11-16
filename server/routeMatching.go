package server

import (
	"fmt"
	"net/http"
	"strings"
)

type RouteMatcher interface {
	Match(r *http.Request) (int, bool)
	RouteString() string
}

type PathPrefixMatcher struct {
	PathPrefix string
}

var _ RouteMatcher = PathPrefixMatcher{}

// Match implements RouteMatcher.
func (p PathPrefixMatcher) Match(r *http.Request) (int, bool) {
	if strings.HasPrefix(r.URL.Path, p.PathPrefix) {
		return len(p.PathPrefix), true
	}
	return 0, false
}

// RouteString implements RouteMatcher.
func (p PathPrefixMatcher) RouteString() string {
	return fmt.Sprintf("PathPrefix: %s", p.PathPrefix)
}
