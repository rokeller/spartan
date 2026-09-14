package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"

	"github.com/rokeller/spartan/internal/test"
)

func TestNotFoundBehavior_getStatusCode(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		b    NotFoundBehavior
		want int
	}{
		{
			name: "Defaults",
			b:    DefaultNotFoundBehavior,
			want: 200,
		},
		{
			name: "FallbackToIndex=nil/StatusCode=nil",
			b:    NotFoundBehavior{},
			want: 200,
		},
		{
			name: "FallbackToIndex=nil/StatusCode ignored",
			b:    NotFoundBehavior{StatusCode: new(123)},
			want: 200,
		},
		{
			name: "FallbackToIndex=true/StatusCode=nil",
			b:    NotFoundBehavior{FallbackToIndex: new(true)},
			want: 200,
		},
		{
			name: "FallbackToIndex=true/StatusCode ignored",
			b:    NotFoundBehavior{FallbackToIndex: new(true), StatusCode: new(123)},
			want: 200,
		},
		{
			name: "FallbackToIndex=false/StatusCode=nil",
			b:    NotFoundBehavior{FallbackToIndex: new(false)},
			want: 404,
		},
		{
			name: "FallbackToIndex=false/StatusCode set",
			b:    NotFoundBehavior{FallbackToIndex: new(false), StatusCode: new(456)},
			want: 456,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.getStatusCode()
			if got != tt.want {
				t.Errorf("getStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFoundBehavior_getResponseBody(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		c    ServerConfig
		want []byte
	}{
		{
			name: "Defaults",
			c: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
			},
			want: test.RepoFileContent(t, "example/content/index.html"),
		},
		{
			name: "FallbackToIndex=nil/BodyPath ignored",
			c: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: nil,
					BodyPath:        new("path/is/irrelevant"),
				},
			},
			want: test.RepoFileContent(t, "example/content/index.html"),
		},
		{
			name: "FallbackToIndex=true/BodyPath ignored",
			c: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(true),
					BodyPath:        new("path/is/irrelevant"),
				},
			},
			want: test.RepoFileContent(t, "example/content/index.html"),
		},
		{
			name: "FallbackToIndex=false/BodyPath=nil",
			c: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(false),
				},
			},
			want: []byte("404 page not found"),
		},
		{
			name: "FallbackToIndex=false/BodyPath does not exist",
			c: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(false),
					BodyPath:        new("/tmp/file/does/not/exist"),
				},
			},
			want: []byte("404 page not found (/tmp/file/does/not/exist)"),
		},
		{
			name: "FallbackToIndex=false/BodyPath set",
			c: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(false),
					BodyPath:        new(test.RepoRelPath(t, "example/content/static/styles.css")),
				},
			},
			want: test.RepoFileContent(t, "example/content/static/styles.css"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.c.GetNotFoundBehavior()
			got := b.getResponseBody(tt.c)
			if !slices.Equal(got, tt.want) {
				t.Errorf("getResponseBody() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNotFoundBehavior_getContentType(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		b    NotFoundBehavior
		want string
	}{
		{
			name: "Defaults",
			b:    DefaultNotFoundBehavior,
			want: "text/html; charset=utf-8",
		},
		{
			name: "FallbackToIndex=nil/ContentType ignored",
			b:    NotFoundBehavior{ContentType: new("test/type")},
			want: "text/html; charset=utf-8",
		},
		{
			name: "FallbackToIndex=true/ContentType ignored",
			b: NotFoundBehavior{
				FallbackToIndex: new(true),
				ContentType:     new("test/type"),
			},
			want: "text/html; charset=utf-8",
		},
		{
			name: "FallbackToIndex=false/ContentType=nil",
			b: NotFoundBehavior{
				FallbackToIndex: new(false),
				ContentType:     nil,
			},
			want: "text/plain; charset=utf-8",
		},
		{
			name: "FallbackToIndex=false/ContentType set",
			b: NotFoundBehavior{
				FallbackToIndex: new(false),
				ContentType:     new("test/type"),
			},
			want: "test/type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.getContentType()
			if got != tt.want {
				t.Errorf("getContentType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNotFoundBehavior_Middleware(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		config          ServerConfig
		wantStatus      int
		wantBody        string
		wantContentType string
	}{
		{
			name: "FallbackToIndex=nil/Request root",
			url:  "/",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
			},
			wantStatus:      200,
			wantContentType: "text/html; charset=utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/index.html"),
		},
		{
			name: "FallbackToIndex=nil/Request styles.css",
			url:  "/static/styles.css",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
			},
			wantStatus:      200,
			wantContentType: "text/css; charset=utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/static/styles.css"),
		},
		{
			name: "FallbackToIndex=nil/Request missing resource",
			url:  "/foo.css",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
			},
			wantStatus:      200,
			wantContentType: "text/html; charset=utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/index.html"),
		},
		{
			name: "FallbackToIndex=true/Request root",
			url:  "/",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(true),
				},
			},
			wantStatus:      200,
			wantContentType: "text/html; charset=utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/index.html"),
		},
		{
			name: "FallbackToIndex=true/Request styles.css",
			url:  "/static/styles.css",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(true),
				},
			},
			wantStatus:      200,
			wantContentType: "text/css; charset=utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/static/styles.css"),
		},
		{
			name: "FallbackToIndex=true/Request missing resource",
			url:  "/foo.css",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(true),
				},
			},
			wantStatus:      200,
			wantContentType: "text/html; charset=utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/index.html"),
		},
		{
			name: "FallbackToIndex=false/Request root",
			url:  "/",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(false),
				},
			},
			wantStatus:      200,
			wantContentType: "text/html; charset=utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/index.html"),
		},
		{
			name: "FallbackToIndex=false/Request styles.css",
			url:  "/static/styles.css",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(false),
				},
			},
			wantStatus:      200,
			wantContentType: "text/css; charset=utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/static/styles.css"),
		},
		{
			name: "FallbackToIndex=false/Request missing resource",
			url:  "/foo.css",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(false),
				},
			},
			wantStatus:      404,
			wantContentType: "text/plain; charset=utf-8",
			wantBody:        "404 page not found",
		},
		{
			name: "FallbackToIndex=false/BodyPath set/ContentType set/Request missing resource",
			url:  "/foo",
			config: ServerConfig{
				StaticContentDir: test.RepoRelPath(t, "example/content"),
				NotFoundBehavior: &NotFoundBehavior{
					FallbackToIndex: new(false),
					BodyPath:        new(test.RepoRelPath(t, "example/content/static/styles.css")),
					ContentType:     new("text/css; charset-utf-8"),
				},
			},
			wantStatus:      404,
			wantContentType: "text/css; charset-utf-8",
			wantBody:        test.RepoFileContentString(t, "example/content/static/styles.css"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{config: tt.config}
			s.fs = http.Dir(s.config.StaticContentDir)
			b := s.config.GetNotFoundBehavior()
			handler :=
				b.Middleware(s.config, http.FileServer(s.fs))

			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("server responded with status %d, want %d", w.Result().StatusCode, tt.wantStatus)
			}

			if w.Result().Header.Get("content-type") != tt.wantContentType {
				t.Errorf("server responded with content type %q, want %q",
					w.Result().Header.Get("content-type"), tt.wantContentType)
			}

			if w.Body.String() != tt.wantBody {
				t.Errorf("server responded with body %q, want %q", w.Body.String(), tt.wantBody)
			}
		})
	}
}

func Test_customNotFoundResponseWriter_Unwrap(t *testing.T) {
	w := httptest.NewRecorder()
	wrapped := customNotFoundResponseWriter{
		ResponseWriter: w,
	}

	got := wrapped.Unwrap()
	if !reflect.DeepEqual(got, w) {
		t.Errorf("Unwrap() = %v, want %v", got, w)
	}
}
