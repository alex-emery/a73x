package html

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

type Highlighter struct {
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
}

func (h Highlighter) htmlHighlight(w io.Writer, source, lang, defaultLang string) error {

	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return h.htmlFormatter.Format(w, h.highlightStyle, it)
}

func (h Highlighter) renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := ""
	lang := string(codeBlock.Info)
	h.htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
}

func (h Highlighter) myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		h.renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func newRenderer() *mdhtml.Renderer {
	htmlFormatter := html.New(html.WithClasses(true), html.TabWidth(2))
	if htmlFormatter == nil {
		panic("couldn't create html formatter")
	}
	styleName := "monokailight"
	highlightStyle := styles.Get(styleName)
	if highlightStyle == nil {
		panic(fmt.Sprintf("didn't find style '%s'", styleName))
	}
	h := Highlighter{
		htmlFormatter:  htmlFormatter,
		highlightStyle: highlightStyle,
	}
	opts := mdhtml.RendererOptions{
		Flags:          mdhtml.CommonFlags,
		RenderNodeHook: h.myRenderHook,
	}
	return mdhtml.NewRenderer(opts)
}
