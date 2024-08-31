// style.go: generate chroma css
// go run contrib/style.go > syntax.css
package main

import (
	"os"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	style := styles.Get("xcode") // Choose your preferred style
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(true), html.TabWidth(2))
	formatter.WriteCSS(os.Stdout, style)
}
