package public

import "embed"

//go:embed static/*.svg
//go:embed static/*.woff2

var FS embed.FS
