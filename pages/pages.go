package pages

import (
	"bytes"
	"fmt"
	"text/template"

	"git.sr.ht/~a73x/home"
	"git.sr.ht/~a73x/home/html"
	"git.sr.ht/~a73x/home/markdown"
)

type GlobalState struct {
	Collections map[string][]markdown.Content
}

type ParserPair struct {
	GlobalState
	markdown.Content
}

func renderTemplate(config GlobalState, content markdown.Content) (string, error) {
	tmpl := content.Meta["template"]
	chosenTemplate := fmt.Sprintf("templates/%s.html", tmpl)
	t, err := template.ParseFS(home.Content, chosenTemplate, "templates/layouts/*.html")
	if err != nil {
		return "", fmt.Errorf("failed to parse layouts: %v", err)
	}

	contentParser, err := template.New("content").Parse(string(content.Body))
	if err != nil {
		return "", fmt.Errorf("failed parsing content: %v", err)
	}

	data := ParserPair{
		config,
		content,
	}

	newContent := &bytes.Buffer{}
	if err := contentParser.Execute(newContent, data); err != nil {
		return "", fmt.Errorf("failed to execute content template: %v", err)
	}

	data.Body = string(html.MDToHTML(newContent.Bytes()))

	b := &bytes.Buffer{}
	if err := t.Execute(b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

type Page struct {
	Path    string
	Content string
}

func Collect() ([]Page, error) {
	contents, err := markdown.ParseContents()
	if err != nil {
		return nil, err
	}

	gs := GlobalState{
		Collections: map[string][]markdown.Content{
			"all": contents,
		},
	}

	for _, content := range contents {
		tags, ok := content.Meta["tags"]
		if !ok {
			continue
		}

		switch tags := tags.(type) {
		case string:
			gs.Collections[tags] = append(gs.Collections[tags], content)
		case []string:
			for _, tag := range tags {
				gs.Collections[tag] = append(gs.Collections[tag], content)
			}
		}
	}

	pages := []Page{}
	for _, content := range contents {
		page, err := renderTemplate(gs, content)
		if err != nil {
			return nil, fmt.Errorf("failed to build site: %v", err)
		}

		path := content.Path
		if path == "index" {
			path = ""
		}

		pages = append(pages, Page{
			Path:    path,
			Content: page,
		})
	}

	return pages, nil
}
