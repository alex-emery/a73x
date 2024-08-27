package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"git.sr.ht/~a73x/home"
	"git.sr.ht/~a73x/home/pages"
	"go.uber.org/zap"
)

func New(logger *zap.Logger) (*http.Server, error) {
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
	pages, err := pages.Collect()
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		mux.HandleFunc("GET /"+page.Path, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(page.Content))
		})
	}

	staticFs, err := fs.Sub(home.Content, "public/static")
	if err != nil {
		return nil, fmt.Errorf("failed to setup static handler: %v", err)
	}

	mux.Handle("GET /static", http.FileServer(http.FS(staticFs)))

	server := http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	return &server, nil
}
