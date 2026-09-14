package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/rokeller/spartan/internal/test"
)

func TestServe(t *testing.T) {
	tests := []struct {
		name    string
		config  ServerConfig
		wantErr bool
	}{
		{
			name:    "Success",
			config:  ServerConfig{Port: 9080},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond*20)
			defer cancel()
			if err := Serve(ctx, tt.config); (err != nil) != tt.wantErr {
				t.Errorf("Serve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_startHttpServer(t *testing.T) {
	type fields struct {
		config ServerConfig
	}
	type args struct {
		wg *sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantUrl string
		wantErr bool
	}{
		{
			name: "HTTP/Serve at root",
			fields: fields{
				config: ServerConfig{Port: 9091, PathRoot: ""},
			},
			wantUrl: "/",
			wantErr: false,
		},
		{
			name: "HTTP/Serve at path with single segment",
			fields: fields{
				config: ServerConfig{Port: 9092, PathRoot: "/simple-path"},
			},
			wantUrl: "/simple-path/",
			wantErr: false,
		},
		{
			name: "HTTP/Serve at path with multiple segments",
			fields: fields{
				config: ServerConfig{Port: 9093, PathRoot: "/path/a/b/c/"},
			},
			wantUrl: "/path/a/b/c/",
			wantErr: false,
		},
		{
			name: "HTTPS/Serve at path with multiple segments",
			fields: fields{
				config: ServerConfig{Port: 9491, PathRoot: "/path/a/b/c/",
					TLSConfig: &TLSConfig{
						CertPath: test.RepoRelPath(t, "example/tls/cert.pem"),
						KeyPath:  test.RepoRelPath(t, "example/tls/key.pem"),
					},
				},
			},
			wantUrl: "/path/a/b/c/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				config: tt.fields.config,
			}
			wg := &sync.WaitGroup{}
			wg.Add(1)
			got, err := s.startHttpServer(wg)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.startHttpServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			req := httptest.NewRequest("GET", tt.wantUrl, nil)
			w := httptest.NewRecorder()
			got.Handler.ServeHTTP(w, req)
			if w.Result().StatusCode != 200 {
				t.Errorf("server responded with status %d, want 200", w.Result().StatusCode)
			}

			if err := got.Shutdown(t.Context()); nil != err {
				t.Errorf("server.Shutdown failed: %v", err)
			}
			wg.Wait()
		})
	}
}

func Test_server_addHealthEndpoints(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "Endpoint/Live",
			url:        "/_spartan/live",
			wantStatus: 200,
		},
		{
			name:       "Endpoint/Runtime",
			url:        "/_spartan/runtime",
			wantStatus: 200,
		},
		{
			name:       "Endpoint/Other",
			url:        "/_spartan/other",
			wantStatus: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			s := &server{}
			s.addHealthEndpoints(mux)

			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("server responded with status %d, want %d", w.Result().StatusCode, tt.wantStatus)
			}
		})
	}
}

func Test_server_getTLSConfig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		c    ServerConfig
		want *tls.Config
	}{
		{
			name: "Empty Config",
			c: ServerConfig{
				TLSConfig: nil,
			},
			want: nil,
		},
		{
			name: "Config/WithPaths",
			c: ServerConfig{
				TLSConfig: &TLSConfig{
					CertPath: "tls.crt",
					KeyPath:  "tls.key",
				},
			},
			want: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server{
				config: tt.c,
			}
			got := s.getTLSConfig()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTLSConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
