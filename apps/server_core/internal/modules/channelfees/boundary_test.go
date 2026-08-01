// Package channelfees hosts only the Q6 architectural boundary test: the
// core (domain/ports/adapters/postgres) must never import a Mercado Livre /
// provider-specific type (feature.md Constraints). This is a flat
// zero-forbidden-imports assertion — no allowlist to shrink, unlike F-04's
// archguard, so a small AST-based scan is enough.
package channelfees

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// forbiddenImportSubstrings marks any import path reaching into a
// provider-specific adapter tree (Mercado Livre, and the connectors/
// integrations adapter layers that host provider code generally — Q6 reads
// "provider-specific", not "ML-only").
var forbiddenImportSubstrings = []string{
	"mercado_livre",
	"mercadolivre",
	"/internal/modules/connectors/adapters/",
	"/internal/modules/integrations/adapters/",
}

// TestChannelFeesCoreNeverImportsProviderSpecificTypes walks every .go file
// under this package's own tree (domain/ports/adapters/postgres) and fails
// if any import path matches a forbidden provider-specific substring.
func TestChannelFeesCoreNeverImportsProviderSpecificTypes(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	var offenders []string
	fset := token.NewFileSet()
	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if parseErr != nil {
			return parseErr
		}
		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			for _, forbidden := range forbiddenImportSubstrings {
				if strings.Contains(importPath, forbidden) {
					offenders = append(offenders, path+" imports "+importPath)
				}
			}
		}
		return nil
	})
	if walkErr != nil {
		t.Fatalf("walk %s: %v", root, walkErr)
	}

	if len(offenders) > 0 {
		t.Fatalf("channelfees core must never import a provider-specific type (Q6):\n%s", strings.Join(offenders, "\n"))
	}
}
