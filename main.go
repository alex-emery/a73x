package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

//go:embed public/index.html
var content embed.FS

func main() {
	logger, _ := zap.NewProduction()
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

	staticFs, err := fs.Sub(content, "public")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.FileServer(http.FS(staticFs)))

	server := http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
