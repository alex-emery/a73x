package html

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func MDToHTML(md []byte, hasToc bool) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes
	p := parser.NewWithExtensions(extensions)

	renderer := newRenderer(hasToc)
	return markdown.ToHTML(md, p, renderer)
}
