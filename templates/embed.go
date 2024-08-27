package templates

import "embed"

//go:embed *.html
//go:embed layouts/*.html
var FS embed.FS
