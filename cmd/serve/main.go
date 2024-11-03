package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"git.sr.ht/~a73x/home/public"
	"go.uber.org/zap"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request received",
				zap.String("url", r.URL.Path),
				zap.String("method", r.Method),
				zap.Duration("duration", time.Since(start)),
				zap.String("user-agent", r.UserAgent()),
			)
		})
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", serveFile)

	server := http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	return server.ListenAndServe()
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	fsPath := strings.TrimRight(r.URL.Path, "/")

	if fsPath == "" {
		fsPath = "index"
	}

	if ext := filepath.Ext(fsPath); ext == "" {
		fsPath += ".html"
	}

	http.ServeFileFS(w, r, public.FS, fsPath)
}
