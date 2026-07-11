package migrations

import (
	"embed"
	"io/fs"
)

// embedded contains the canonical, forward-only migration set.
//
//go:embed *.sql
var embedded embed.FS

// Source returns the canonical migration filesystem. Callers must identify
// migrations by full filename; numeric prefixes are intentionally not unique.
func Source() fs.FS {
	return embedded
}
