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

func walk(filename string, linkDirname string, walkFn filepath.WalkFunc) error {
	symWalkFunc := func(path string, info os.FileInfo, err error) error {

		if fname, err := filepath.Rel(filename, path); err == nil {
			path = filepath.Join(linkDirname, fname)
		} else {
			return err
		}

		if err == nil && info.Mode()&os.ModeSymlink == os.ModeSymlink {
			finalPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			info, err := os.Lstat(finalPath)
			if err != nil {
				return walkFn(path, info, err)
			}
			if info.IsDir() {
				return walk(finalPath, path, walkFn)
			}
		}

		return walkFn(path, info, err)
	}
	return filepath.Walk(filename, symWalkFunc)
}

// Walk extends filepath.Walk to also follow symlinks
func Walk(path string, walkFn filepath.WalkFunc) error {
	return walk(path, path, walkFn)
}

// glob returns relative file paths that match the extension
func glob(dir string, ext string) ([]string, error) {
	files := []string{}
	err := Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			// just diff
			relPath, _ := filepath.Rel(dir, path)

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

	path := strings.Replace(contentPath, ".md", ".html", 1)

	mc := Content{
		Body: string(rest),
		Path: path,
		Meta: matter,
	}

	return mc, nil
}
