package templates

import "embed"

//go:embed *.gohtml **/*.gohtml **/**/*.gohtml
var FS embed.FS
