package templates

import "embed"

// FS is a variable of type FS (filesystem). The following line defines that all files ending with
// .gohtml will be embedded into the binary.
//go:embed *.gohtml
var FS embed.FS
