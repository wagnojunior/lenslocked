package migrations

import "embed"

// FS is a variable of type FS (filesystem). The following line defines that all files ending with
// .sql will be embedded into the binary.
//go:embed *.sql
var FS embed.FS
