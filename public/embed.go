package public

import "embed"

//go:embed static/*.svg
//go:embed static/*.woff2
//go:embed static/*.css
var FS embed.FS
