package public

import "embed"

//go:embed static/*.svg
//go:embed static/*.woff2
//go:embed static/*.css
//go:embed static/*.js
var FS embed.FS
