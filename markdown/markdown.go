package markdown

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"git.sr.ht/~a73x/home/content"
	"github.com/adrg/frontmatter"
)

type Content struct {
	Meta map[string]any
	Body string
	Path string
}

func ParseContents() ([]Content, error) {
	contentFiles, err := glob(content.FS, ".", ".md")
	if err != nil {
		return nil, fmt.Errorf("failed to glob: %v", err)
	}

	res := make([]Content, 0, len(contentFiles))
	for _, contentFile := range contentFiles {
		c, err := parseMarkdownFile(content.FS, contentFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read markdown file: %v", err)
		}

		res = append(res, c)
	}

	return res, nil
}

func glob(embedded fs.FS, dir string, ext string) ([]string, error) {
	files := []string{}
	err := fs.WalkDir(embedded, dir, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func parseMarkdownFile(embedded fs.FS, path string) (Content, error) {
	input, err := fs.ReadFile(embedded, path)
	if err != nil {
		return Content{}, fmt.Errorf("failed to read post: %v", err)
	}

	matter := map[string]any{}
	matter["template"] = "default"
	rest, err := frontmatter.Parse(bytes.NewBuffer(input), &matter)
	if err != nil {
		return Content{}, fmt.Errorf("failed to parse frontmatter: %v", err)
	}
	path = strings.Replace(path, ".md", "", 1)
	if path == "index" {
		path = ""
	}

	path = "/" + path

	mc := Content{
		Body: string(rest),
		Path: path,
		Meta: matter,
	}

	return mc, nil
}
