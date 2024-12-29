package html

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/toc"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	chromastyles "github.com/alecthomas/chroma/v2/styles"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

func MDToHTML(md []byte, hasToc bool) []byte {
	style := chromastyles.Get("autumn")
	style, err := style.Builder().Add(chroma.Background, "bg:#f0f0f0").Build()
	if err != nil {
		return nil
	}

	extensions := []goldmark.Extender{
		extension.GFM,
		extension.Footnote,
		highlighting.NewHighlighting(
			highlighting.WithCustomStyle(style),
			highlighting.WithFormatOptions(
				chromahtml.WrapLongLines(true),
				chromahtml.InlineCode(true),
				chromahtml.TabWidth(4),
				chromahtml.PreventSurroundingPre(false),
			),
		),
	}

	if hasToc {
		extensions = append(extensions, &toc.Extender{
			Compact: true,
		})
	}

	converter := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	err = converter.Convert(md, &buf)
	if err != nil {
		return nil
	}

	return buf.Bytes()
}
