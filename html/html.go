package html

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func MDToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	// doc := p.Parse(md)

	// htmlFlags := html.CommonFlags
	// opts := html.RendererOptions{Flags: htmlFlags}
	// renderer := html.NewRenderer(opts)

	// return markdown.Render(doc, renderer)
	renderer := newRenderer()
	return markdown.ToHTML(md, p, renderer)
}
