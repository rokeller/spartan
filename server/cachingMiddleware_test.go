package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func toPtr[T any](x T) *T {
	return &x
}

func Test_withCachingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	tests := []struct {
		name   string
		config Cache
		want   string
	}{
		{
			name:   "Default config",
			config: Cache{},
			want:   "",
		},
		{
			name:   "Config/DefaultPolicy/Empty",
			config: Cache{DefaultPolicy: &CachePolicy{}},
			want:   "",
		},
		{
			name:   "Config/DefaultPolicy/With-max-age",
			config: Cache{DefaultPolicy: &CachePolicy{MaxAge: toPtr(1 * time.Hour)}},
			want:   "max-age=3600",
		},
		{
			name:   "Config/DefaultPolicy/With-s-maxage",
			config: Cache{DefaultPolicy: &CachePolicy{SharedMaxAge: toPtr(2 * time.Hour)}},
			want:   "s-maxage=7200",
		},
		{
			name:   "Config/DefaultPolicy/With-stale-if-error",
			config: Cache{DefaultPolicy: &CachePolicy{StaleIfError: toPtr(10 * time.Minute)}},
			want:   "stale-if-error=600",
		},
		{
			name:   "Config/DefaultPolicy/With-stale-while-revalidate",
			config: Cache{DefaultPolicy: &CachePolicy{StaleWhileRevalidate: toPtr(20 * time.Minute)}},
			want:   "stale-while-revalidate=1200",
		},
		{
			name:   "Config/DefaultPolicy/With-immutable",
			config: Cache{DefaultPolicy: &CachePolicy{Immutable: true}},
			want:   "immutable",
		},
		{
			name:   "Config/DefaultPolicy/With-must-revalidate",
			config: Cache{DefaultPolicy: &CachePolicy{MustRevalidate: true}},
			want:   "must-revalidate",
		},
		{
			name:   "Config/DefaultPolicy/With-must-understand",
			config: Cache{DefaultPolicy: &CachePolicy{MustUnderstand: true}},
			want:   "must-understand",
		},
		{
			name:   "Config/DefaultPolicy/With-no-cache",
			config: Cache{DefaultPolicy: &CachePolicy{NoCache: true}},
			want:   "no-cache",
		},
		{
			name:   "Config/DefaultPolicy/With-no-store",
			config: Cache{DefaultPolicy: &CachePolicy{NoStore: true}},
			want:   "no-store",
		},
		{
			name:   "Config/DefaultPolicy/With-no-transform",
			config: Cache{DefaultPolicy: &CachePolicy{NoTransform: true}},
			want:   "no-transform",
		},
		{
			name:   "Config/DefaultPolicy/With-private",
			config: Cache{DefaultPolicy: &CachePolicy{Private: true}},
			want:   "private",
		},
		{
			name:   "Config/DefaultPolicy/With-proxy-revalidate",
			config: Cache{DefaultPolicy: &CachePolicy{ProxyRevalidate: true}},
			want:   "proxy-revalidate",
		},
		{
			name:   "Config/DefaultPolicy/With-public",
			config: Cache{DefaultPolicy: &CachePolicy{Public: true}},
			want:   "public",
		},

		{
			name: "Config/NoRouteMatch",
			config: Cache{Routes: []RouteMatchingCachePolicy{
				{Match: PathPrefixMatcher{PathPrefix: "not-found"}},
			}},
			want: "",
		},
		{
			name: "Config/RouteWithoutMatcher",
			config: Cache{Routes: []RouteMatchingCachePolicy{
				{Match: nil, Policy: CachePolicy{Public: true}},
			}},
			want: "",
		},
		{
			name: "Config/SingleRouteWithMatch",
			config: Cache{Routes: []RouteMatchingCachePolicy{
				{
					Match:  PathPrefixMatcher{PathPrefix: "/cach"},
					Policy: CachePolicy{Public: true, MaxAge: toPtr(10 * time.Second)},
				},
			}},
			want: "max-age=10, public",
		},
		{
			name: "Config/MultipleRoutesWithMatch",
			config: Cache{Routes: []RouteMatchingCachePolicy{
				{
					Match:  PathPrefixMatcher{PathPrefix: "/cach"},
					Policy: CachePolicy{Public: true, MaxAge: toPtr(10 * time.Second)},
				},
			}},
			want: "max-age=10, public",
		},
		{
			name: "Config/MultipleRoutesWithMatch",
			config: Cache{Routes: []RouteMatchingCachePolicy{
				{
					Match:  PathPrefixMatcher{PathPrefix: "/cach"},
					Policy: CachePolicy{Public: true, MaxAge: toPtr(10 * time.Second)},
				},
				{
					Match:  PathPrefixMatcher{PathPrefix: "/cachable"},
					Policy: CachePolicy{Public: true, MaxAge: toPtr(20 * time.Second)},
				},
				{
					Match:  PathPrefixMatcher{PathPrefix: "/cachable-res"},
					Policy: CachePolicy{Public: true, MaxAge: toPtr(30 * time.Second)},
				},
			}},
			want: "max-age=30, public",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := withCachingMiddleware(tt.config, handler)

			req := httptest.NewRequest("GET", "/cachable-resource", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			if w.Result().StatusCode != 200 {
				t.Errorf("server responded with status %d, want 200", w.Result().StatusCode)
			}

			got := w.Result().Header.Get("Cache-Control")
			if tt.want != got {
				t.Errorf("header mismatch: got = %s, want = %s", got, tt.want)
			}
		})
	}
}
