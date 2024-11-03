package pages

import (
	"bytes"
	"fmt"
	"path"
	"sort"
	"text/template"

	"git.sr.ht/~a73x/home/html"
	"git.sr.ht/~a73x/home/markdown"
)

type Navigation struct {
	Title string
	Path  string
}
type GlobalState struct {
	Collections map[string][]markdown.Content
	Nav         []Navigation
}

type ParserPair struct {
	GlobalState
	markdown.Content
}

func renderTemplate(config GlobalState, content markdown.Content) (string, error) {
	tmpl := content.Meta["template"]
	chosenTemplate := fmt.Sprintf("%s.html", tmpl)
	t, err := template.ParseGlob("templates/layouts/*.html")
	if err != nil {
		return "", fmt.Errorf("failed to parse layouts: %v", err)
	}

	t, err = t.ParseFiles(path.Join("templates", chosenTemplate))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
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

	hasToc := false
	_, ok := content.Meta["toc"]
	if ok {
		hasToc = content.Meta["toc"].(bool)
	}

	data.Body = string(html.MDToHTML(newContent.Bytes(), hasToc))

	b := &bytes.Buffer{}
	if err := t.ExecuteTemplate(b, chosenTemplate, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

type Page struct {
	Path    string
	Content string
}

func Collect(contentPath string) ([]Page, error) {
	contents, err := markdown.ParseContents(contentPath)
	if err != nil {
		return nil, err
	}

	gs := GlobalState{
		Collections: map[string][]markdown.Content{
			"all": contents,
		},
		Nav: make([]Navigation, 0),
	}

	for _, content := range contents {
		tags, ok := content.Meta["tags"]
		if ok {
			switch tags := tags.(type) {
			case string:
				gs.Collections[tags] = append(gs.Collections[tags], content)
			case []string:
				for _, tag := range tags {
					gs.Collections[tag] = append(gs.Collections[tag], content)
				}
			}
		}

		nav, ok := content.Meta["nav"]
		if ok {
			gs.Nav = append(gs.Nav, Navigation{content.Meta["title"].(string), nav.(string)})
		}
	}

	sortNavBar(gs.Nav)

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

func sortNavBar(nav []Navigation) {
	// see no evil, speak no evil
	sort.Slice(nav, func(i, j int) bool {
		order := map[string]int{
			"home":  1,
			"posts": 2,
			"about": 3,
		}

		// Get the order value, default to a high number if not in the predefined order
		orderI, okI := order[nav[i].Title]
		if !okI {
			orderI = len(nav) + 1
		}

		orderJ, okJ := order[nav[j].Title]
		if !okJ {
			orderJ = len(nav) + 1
		}

		// If both strings have the same order value, sort lexicographically
		if orderI == orderJ {
			return nav[i].Title < nav[j].Title
		}

		return orderI < orderJ
	})
}
