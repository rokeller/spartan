package server

import (
	"net/http"
)

func withCachingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "max-age=31536000")
		next.ServeHTTP(w, r)
	})
}
