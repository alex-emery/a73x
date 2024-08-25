package home

import "embed"

//go:embed public/static/*.svg
//go:embed public/static/*.woff2
//go:embed templates/*.html
//go:embed templates/layouts/*.html
//go:embed posts/*.md
var Content embed.FS
