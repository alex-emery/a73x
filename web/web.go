package web

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"git.sr.ht/~a73x/home"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"go.uber.org/zap"
)

type CommonData struct {
	Posts []PostData
}

type PostData struct {
	Title string
	Path  string
}

func GeneratePosts(mux *http.ServeMux) ([]PostData, error) {
	converter := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		))

	t, err := template.ParseFS(home.Content, "templates/layouts/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse layouts: %v", err)
	}

	final := []PostData{}
	posts, err := fs.Glob(home.Content, "posts/*.md")
	if err != nil {
		return nil, fmt.Errorf("failed to glob posts: %v", err)
	}

	for _, post := range posts {
		postName := filepath.Base(post)
		postName = strings.Split(postName, ".")[0] + ".html"

		// parse markdown
		input, err := fs.ReadFile(home.Content, post)
		if err != nil {
			return nil, fmt.Errorf("failed to read post: %v", err)
		}

		var b bytes.Buffer
		ctx := parser.NewContext()
		if err := converter.Convert(input, &b, parser.WithContext(ctx)); err != nil {
			return nil, err
		}

		d := frontmatter.Get(ctx)
		var meta map[string]string
		if err := d.Decode(&meta); err != nil {
			return nil, err
		}

		foo, err := t.ParseFS(home.Content, "templates/post.html")
		if err != nil {
			return nil, fmt.Errorf("failed to parse post template: %v", err)
		}

		content, err := io.ReadAll(&b)
		if err != nil {
			return nil, fmt.Errorf("failed to read post template: %v", err)
		}

		type IPostData struct {
			Title string
			Post  string
		}

		mux.HandleFunc("/posts/"+postName, func(w http.ResponseWriter, r *http.Request) {
			if err := foo.ExecuteTemplate(w, "post.html", IPostData{
				Title: meta["title"],
				Post:  string(content),
			}); err != nil {
				fmt.Println(err)
			}
		})

		final = append(final, PostData{
			Title: meta["title"],
			Path:  "posts/" + postName,
		})
	}

	return final, nil
}

func loadTemplates(mux *http.ServeMux, data any) error {
	tmplFiles, err := fs.ReadDir(home.Content, "templates")
	if err != nil {
		return fmt.Errorf("failed to parse template layouts")
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(home.Content, "templates/"+tmpl.Name(), "templates/layouts/*.html")
		if err != nil {
			return fmt.Errorf("failed to parse template "+tmpl.Name(), err)
		}

		if tmpl.Name() == "index.html" {
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				pt.ExecuteTemplate(w, tmpl.Name(), data)
			})
		} else if tmpl.Name() == "page.html" {
			continue
		} else {
			mux.HandleFunc("/"+tmpl.Name(), func(w http.ResponseWriter, r *http.Request) {
				pt, err := template.ParseFS(home.Content, "templates/"+tmpl.Name(), "templates/layouts/*.html")
				if err != nil {
					fmt.Printf("failed to parse template "+tmpl.Name(), err)
				}
				pt.ExecuteTemplate(w, tmpl.Name(), data)
			})
		}
	}
	return nil
}
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
	postData, err := GeneratePosts(mux)
	if err != nil {
		return nil, fmt.Errorf("failed to generate posts: %v", err)
	}

	data := CommonData{
		Posts: postData,
	}

	if err := loadTemplates(mux, data); err != nil {
		return nil, fmt.Errorf("failed to parse templates: %v", err)
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
