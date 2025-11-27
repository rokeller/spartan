package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
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
		// fs      http.FileSystem
		// handler http.Handler
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
			name: "Serve at root",
			fields: fields{
				config: ServerConfig{Port: 9091, PathRoot: ""},
			},
			wantUrl: "/",
			wantErr: false,
		},
		{
			name: "Serve at path with single segment",
			fields: fields{
				config: ServerConfig{Port: 9092, PathRoot: "/simple-path"},
			},
			wantUrl: "/simple-path/",
			wantErr: false,
		},
		{
			name: "Serve at path with multiple segments",
			fields: fields{
				config: ServerConfig{Port: 9093, PathRoot: "/path/a/b/c/"},
			},
			wantUrl: "/path/a/b/c/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				config: tt.fields.config,
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(r.URL.String()))
				}),
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
