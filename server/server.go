package server

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

type server struct {
	port             uint16
	staticContentDir string
	serverRootPath   string

	fs      http.FileSystem
	handler http.Handler
}

func Serve(
	ctx context.Context,
	port uint16,
	staticContentDir string,
	serverRootPath string,
) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	s := &server{
		port:             port,
		staticContentDir: staticContentDir,
		serverRootPath:   serverRootPath,
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

	// Normalize the server root path to start with a single slash and without a trailing slash.
	normalizedPath := "/" + strings.TrimLeft(strings.TrimRight(s.serverRootPath, "/"), "/")
	if normalizedPath == "/" {
		s.serverRootPath = normalizedPath
	} else {
		s.serverRootPath = normalizedPath + "/"
	}

	s.fs = http.Dir(s.staticContentDir)
	s.handler = http.FileServer(s.fs)

	h := http.StripPrefix(
		s.serverRootPath,
		http.HandlerFunc(s.handleStaticFiles))
	p := fmt.Sprintf("GET %s", s.serverRootPath)
	mux.Handle(p, withLoggingMiddleware(withCachingMiddleware(h)))
	klog.V(1).InfoS("Registered handler", "pattern", p)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,

		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		// Let the caller know that we're done.
		defer wg.Done()

		klog.InfoS("Starting web server",
			"port", s.port,
			"staticContentDir", s.staticContentDir,
			"serverRootPath", s.serverRootPath)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			klog.ErrorS(err, "Failed to start web server", "port", s.port)
		}
	}()

	return srv, nil
}

func (s *server) handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	upath := path.Clean(r.URL.Path)
	if upath == "." {
		upath = "/"
	} else if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	r.URL.Path = upath
	klog.V(5).InfoS("Client requested", "path", upath)

	f, err := s.fs.Open(upath)
	if nil != err {
		klog.V(4).ErrorS(err, "Failed")
		klog.V(5).Info("Serving index.html from staticContentDir instead", "staticContentDir", s.staticContentDir)
		r.URL.Path = "/"
		// Recurse so we can leverage the same logic as for all other files.
		s.handleStaticFiles(w, r)
		return
	}
	defer f.Close()

	// Let the configured handler serve the request.
	s.handler.ServeHTTP(w, r)
}
