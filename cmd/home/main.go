package main

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"git.sr.ht/~a73x/home"
	"go.uber.org/zap"
)

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

	staticFs, err := fs.Sub(home.Content, "public")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("GET /", http.FileServer(http.FS(staticFs)))

	server := http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
