package server

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"k8s.io/klog/v2"
)

var (
	DefaultNotFoundBehavior = NotFoundBehavior{
		FallbackToIndex: new(true),
	}
)

// NotFoundBehavior defines the behavior for responses on requests for resources
// which are not found.
type NotFoundBehavior struct {
	// FallbackToIndex indicates whether spartan should respond with the
	// index.html of the configured root directory when a resource is not found.
	// Defaults to true.
	FallbackToIndex *bool
	// StatusCode defines the status code to use when responding to a request
	// for a resource that is not found. Defaults to [http.StatusNotFound] (404)
	// when [NotFoundBehavior.FallbackToIndex] is false and is ignored otherwise.
	StatusCode *int
	// BodyPath defines the path to a file that spartan should use to respond
	// to requests for resources which are not found. Ignored when
	// [NotFoundBehavior.FallbackToIndex] is true or nil and defaults to nil
	// otherwise, thus resulting in the default 404 response body "404 page not
	// found" to be returned. If there are any errors reading the contents of
	// this file, the response body for 404 responses will be set to "404 page
	// not found (<body-path-value>)".
	BodyPath *string
	// ContentType defines the media type (MIME type) of the resource defined
	// in the [NotFoundBehavior.BodyPath]. Defaults to "text/html; charset=utf-8"
	// when [NotFoundBehavior.FallbackToIndex] is true or nil and defaults to
	// "text/plain; charset=utf-8" otherwise.
	ContentType *string
}

func withNotFoundMiddleware(c ServerConfig, next http.Handler) http.Handler {
	b := c.GetNotFoundBehavior()
	return b.Middleware(c, next)
}

func (b *NotFoundBehavior) Middleware(c ServerConfig, next http.Handler) http.Handler {
	wrapResponseWriter := func(w http.ResponseWriter) http.ResponseWriter {
		wrapped := &customNotFoundResponseWriter{
			ResponseWriter:     w,
			notFoundStatusCode: b.getStatusCode(),
			body:               b.getResponseBody(c),
			contentType:        b.getContentType(),
		}
		return wrapped
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(w)
		next.ServeHTTP(w, r)
	})
}

func (b *NotFoundBehavior) getStatusCode() int {
	if b.FallbackToIndex == nil || *b.FallbackToIndex {
		return http.StatusOK
	}

	if b.StatusCode == nil {
		return http.StatusNotFound
	} else {
		return *b.StatusCode
	}
}

func (b *NotFoundBehavior) getResponseBody(c ServerConfig) []byte {
	var notFoundPath string
	if b.FallbackToIndex == nil || *b.FallbackToIndex {
		notFoundPath = path.Join(c.StaticContentDir, "index.html")
	} else if b.BodyPath == nil {
		return []byte("404 page not found")
	} else {
		notFoundPath = *b.BodyPath
	}

	if data, err := os.ReadFile(notFoundPath); err != nil {
		klog.ErrorS(err, "Failed to read file", "path", notFoundPath)
		return []byte(fmt.Sprintf("404 page not found (%s)", notFoundPath))
	} else {
		return data
	}
}

func (b *NotFoundBehavior) getContentType() string {
	if b.FallbackToIndex == nil || *b.FallbackToIndex {
		return "text/html; charset=utf-8"
	}

	if b.ContentType == nil {
		return "text/plain; charset=utf-8"
	} else {
		return *b.ContentType
	}
}

// customNotFoundResponseWriter implements a [http.ResponseWriter] that responds
// with a pre-defined response when a requested resource is not found.
type customNotFoundResponseWriter struct {
	http.ResponseWriter
	origStatusCode     int
	notFoundStatusCode int
	contentType        string
	body               []byte
}

var _ http.ResponseWriter = &customNotFoundResponseWriter{}

// WriteHeader implements [http.ResponseWriter.WriteHeader].
func (w *customNotFoundResponseWriter) WriteHeader(statusCode int) {
	w.origStatusCode = statusCode

	if statusCode == http.StatusNotFound {
		w.ResponseWriter.Header().Set("content-type", w.contentType)
		statusCode = w.notFoundStatusCode
	}
	w.ResponseWriter.WriteHeader(statusCode)

	if w.origStatusCode == http.StatusNotFound {
		w.ResponseWriter.Write(w.body)
	}
}

// Write implements [http.ResponseWriter.Write].
func (w *customNotFoundResponseWriter) Write(body []byte) (int, error) {
	if w.origStatusCode != http.StatusNotFound {
		return w.ResponseWriter.Write(body)
	}
	return 0, nil
}

// Unwrap returns the underlying ResponseWriter
func (w *customNotFoundResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
