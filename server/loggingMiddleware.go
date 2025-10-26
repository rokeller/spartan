package server

import (
	"net/http"
	"time"

	"k8s.io/klog/v2"
)

type loggingMiddlewareWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingMiddlewareWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingMiddlewareWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func withLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		ww := &loggingMiddlewareWriter{ResponseWriter: w}
		next.ServeHTTP(ww, r)

		d := time.Since(startTime)
		klog.V(2).InfoS("Responded",
			"method", r.Method,
			"uri", r.RequestURI,
			"resultCode", ww.status,
			"duration", d)
	})
}
