package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"git.sr.ht/~a73x/home/pages"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func Build(directory string) error {
	pages, err := pages.Collect(directory)
	if err != nil {
		return err
	}

	var errs []error
	for _, page := range pages {
		fmt.Println("building", page.Path)
		err = writeFile(path.Join("public", page.Path), []byte(page.Content))
		if err != nil {
			errs = append(errs, err)
		}
	}

	if errs != nil {
		return errors.Join(errs...)
	}

	return nil
}

func Serve() error {
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

	http.ServeFile(w, r, path.Join("public", fsPath))
}

func watchDir(watchDir string) error {
	// Directory to watch

	// Ensure the directory exists
	if _, err := os.Stat(watchDir); os.IsNotExist(err) {
		return err
	}

	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Start a goroutine to process events
	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("Event: %s", event)

				// Trigger a command when a change is detected
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Remove == fsnotify.Remove {
					err := Build(watchDir)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("built")
					}

				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error: %v", err)
			}
		}
	}()

	// Add the directory to the watcher
	err = filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			log.Printf("Watching directory: %s", path)
			return watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
func Run() error {
	contentDir := "content"
	actualPath, err := filepath.EvalSymlinks(contentDir)
	if err != nil {
		return err
	}

	if err := Build(actualPath); err != nil {
		return err
	}

	// watcher
	if err := watchDir(actualPath); err != nil {
		return err
	}

	if err := Serve(); err != nil {
		return err
	}

	return nil
}

func writeFile(name string, contents []byte) error {
	folders := path.Dir(name)
	_, err := os.Stat(folders)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(folders, 0744); err != nil {
			return fmt.Errorf("failed to mkdir %s\n%w", folders, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat folder %s\n%w", folders, err)
	}

	err = os.WriteFile(name, contents, 0666)
	if err != nil {
		return fmt.Errorf("failed to write file %s\n%w", name, err)
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Panic(err)
	}
}
