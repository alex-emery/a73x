package content

import "embed"

//go:embed *.md
//go:embed posts/*.md
var FS embed.FS
