package home

import "embed"

//go:embed public/static/*.woff2
//go:embed public/*.html
var Content embed.FS
