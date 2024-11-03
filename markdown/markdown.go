package markdown

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

type Content struct {
	Meta map[string]any
	Body string
	Path string
}

func ParseContents(contentPath string) ([]Content, error) {
	contentFiles, err := glob(contentPath, ".md")
	if err != nil {
		return nil, fmt.Errorf("failed to glob: %v", err)
	}

	res := make([]Content, 0, len(contentFiles))
	for _, contentFile := range contentFiles {
		c, err := parseMarkdownFile(contentPath, contentFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read markdown file: %v", err)
		}

		res = append(res, c)
	}

	return res, nil
}

// glob returns relative file paths that match the extension
func glob(dir string, ext string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ext {
			// just diff
			relPath, _ := filepath.Rel(dir, p)

			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}

func parseMarkdownFile(basePath, contentPath string) (Content, error) {
	input, err := os.ReadFile(path.Join(basePath, contentPath))
	if err != nil {
		return Content{}, fmt.Errorf("failed to read post: %v", err)
	}

	matter := map[string]any{}
	matter["template"] = "default"
	rest, err := frontmatter.Parse(bytes.NewBuffer(input), &matter)
	if err != nil {
		return Content{}, fmt.Errorf("failed to parse frontmatter: %v", err)
	}

	path := strings.Replace(contentPath, ".md", "", 1)
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
