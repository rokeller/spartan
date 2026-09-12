package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

type server struct {
	config ServerConfig

	fs http.FileSystem
}

func Serve(
	ctx context.Context,
	config ServerConfig,
) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	s := &server{
		config: config,
	}
	srv, err := s.startHttpServer(wg)
	if nil != err {
		return err
	}

	// Wait for the context to be cancelled or done.
	<-ctx.Done()
	klog.V(2).Info("Initiated web server shutdown")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); nil != err {
		klog.Error(err, "Failed to shutdown web server")
	}
	wg.Wait()
	klog.Info("Shut down web server")
	return nil
}

func (s *server) startHttpServer(wg *sync.WaitGroup) (*http.Server, error) {
	mux := http.NewServeMux()
	s.addHealthEndpoints(mux)

	// Normalize the server root path to start with a single slash and without a trailing slash.
	normalizedPath := "/" + strings.TrimLeft(strings.TrimRight(s.config.PathRoot, "/"), "/")
	if normalizedPath == "/" {
		s.config.PathRoot = normalizedPath
	} else {
		s.config.PathRoot = normalizedPath + "/"
	}

	s.fs = http.Dir(s.config.StaticContentDir)
	handler := http.FileServer(s.fs)
	p := fmt.Sprintf("GET %s", s.config.PathRoot)
	mux.Handle(p,
		withLoggingMiddleware(
			http.StripPrefix(s.config.PathRoot,
				withSecurityMiddleware(s.config.Security,
					withCachingMiddleware(s.config.Cache,
						withNotFoundMiddleware(
							s.config,
							handler,
						),
					),
				),
			),
		),
	)
	klog.V(1).InfoS("Registered handler", "pattern", p)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: mux,

		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		// Let the caller know that we're done.
		defer wg.Done()

		klog.InfoS("Starting web server",
			"port", s.config.Port,
			"staticContentDir", s.config.StaticContentDir,
			"serverPathRoot", s.config.PathRoot)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			klog.ErrorS(err, "Failed to start web server", "port", s.config.Port)
			os.Exit(2)
		}
	}()

	return srv, nil
}

func (s *server) addHealthEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /_spartan/live", s.healthEndpointLive)
	mux.HandleFunc("GET /_spartan/runtime", s.healthEndpointRuntime)
}

func (s *server) healthEndpointLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
	})
}

func (s *server) healthEndpointRuntime(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := map[string]any{
		"goroutines": runtime.NumGoroutine(),
		"memory": map[string]uint64{
			"totalAlloc": memStats.TotalAlloc,
			"sys":        memStats.Sys,
			"heapAlloc":  memStats.HeapAlloc,
			"heapSys":    memStats.HeapSys,
			"stackInuse": memStats.StackInuse,
			"stackSys":   memStats.StackSys,
		},
		"gc": map[string]uint32{
			"numGC": memStats.NumGC,
		},
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
