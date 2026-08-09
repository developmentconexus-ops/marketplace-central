package composition

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// ADR-023 §2: a module is importable by another module ONLY at X/ports.
// X/domain, X/application, X/adapters, X/transport, X/wiring are private.
//
// The two carve-outs, as amended 2026-08-05, are encoded below and nowhere else:
//
//  1. internal/composition — the composition root — may import any layer of any module.
//     Enforced by construction: the root does not live under internal/modules, so this
//     walk never visits it.
//  2. The shared core (§5) — today sourcekind — is importable by any module at any layer.
//
// Two things are deliberately NOT carve-outs:
//
//   - wiring is public to the composition root only. A module receives another module's
//     assembled object from the root; it does not reach for it.
//   - adapters/ has no cross-module licence. An adapter translates between our domain and a
//     FOREIGN system — provider HTTP, Oracle, postgres — and that has nothing to do with the
//     module boundary. An adapter reaching into a sibling module's domain is coupling with
//     an extra step, not translation. inventory/adapters/internalread consuming
//     internal_read/ports is the shape that proves no carve-out is needed.
//
// This test is the RED that opens Wave 1. It reports counts broken down by origin layer,
// because the wave's acceptance is a measurement and not a green test.
//
// Corrigido em 2026-08-06: esta caminhada saltava todo o ficheiro _test.go e por
// isso reportava um número que não era o do repositório — 100 imports cruzados
// em 67 ficheiros de teste ficavam de fora. Um teste é código e vive sob as
// mesmas regras. A contagem sobe porque o instrumento passou a ver, não porque
// o código piorou.
const modulesImportPrefix = "marketplace-central/apps/server_core/internal/modules/"

// sharedCore is carve-out 3. It is a closed list; ADR-023 §5 requires an amendment to grow
// it, so a new entry here without a matching ADR amendment is itself the defect.
var sharedCore = map[string]bool{
	"sourcekind": true,
}

type boundaryViolation struct {
	file       string
	line       int
	fromModule string
	fromLayer  string
	toModule   string
	toLayer    string
}

// moduleAndLayer splits a path relative to internal/modules into its module and its layer
// directory. A file sitting directly in the module root has an empty layer — which is not a
// sanctioned shape (ADR-023 §8) but is reported honestly rather than skipped.
func moduleAndLayer(rel string) (module, layer string) {
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) == 0 {
		return "", ""
	}
	module = parts[0]
	if len(parts) >= 3 {
		layer = parts[1]
	}
	return module, layer
}

func importedModuleAndLayer(importPath string) (module, layer string, ok bool) {
	if !strings.HasPrefix(importPath, modulesImportPrefix) {
		return "", "", false
	}
	parts := strings.Split(strings.TrimPrefix(importPath, modulesImportPrefix), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", "", false
	}
	module = parts[0]
	if len(parts) >= 2 {
		layer = parts[1]
	}
	return module, layer, true
}

func TestModuleBoundaryADR023(t *testing.T) {
	modulesDir := filepath.Join("..", "modules")
	fset := token.NewFileSet()

	var violations []boundaryViolation
	filesParsed := 0

	err := filepath.WalkDir(modulesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		filesParsed++

		rel, err := filepath.Rel(modulesDir, path)
		if err != nil {
			return err
		}
		fromModule, fromLayer := moduleAndLayer(rel)

		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			toModule, toLayer, ok := importedModuleAndLayer(importPath)
			if !ok || toModule == fromModule {
				continue
			}
			// Carve-out 2: the shared core is importable by anyone, at any layer.
			if sharedCore[toModule] {
				continue
			}
			// The rule itself: ports is the only public surface.
			if toLayer == "ports" {
				continue
			}
			violations = append(violations, boundaryViolation{
				file:       filepath.ToSlash(filepath.Join("internal/modules", rel)),
				line:       fset.Position(spec.Pos()).Line,
				fromModule: fromModule,
				fromLayer:  fromLayer,
				toModule:   toModule,
				toLayer:    toLayer,
			})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", modulesDir, err)
	}

	// Logged on every outcome, not only on failure: the gate's boundary lane
	// reads this count, and a run that parsed zero files must be distinguishable
	// from a tree with zero violations. WalkDir over a missing root fails loudly
	// above, but an empty or hollowed-out root does not.
	t.Logf("boundary_files=%d", filesParsed)
	if filesParsed == 0 {
		t.Fatalf("parsed zero Go files under %s; a verdict over the empty set is not a verdict", modulesDir)
	}

	if len(violations) == 0 {
		return
	}

	byTarget := map[string]int{}
	byOrigin := map[string]int{}
	for _, v := range violations {
		target := v.toModule
		if v.toLayer != "" {
			target += "/" + v.toLayer
		}
		byTarget[target]++

		origin := v.fromLayer
		if origin == "" {
			origin = "(module root, no layer)"
		}
		byOrigin[origin]++
	}
	origins := make([]string, 0, len(byOrigin))
	for origin := range byOrigin {
		origins = append(origins, origin)
	}
	sort.Slice(origins, func(i, j int) bool {
		if byOrigin[origins[i]] != byOrigin[origins[j]] {
			return byOrigin[origins[i]] > byOrigin[origins[j]]
		}
		return origins[i] < origins[j]
	})
	targets := make([]string, 0, len(byTarget))
	for target := range byTarget {
		targets = append(targets, target)
	}
	sort.Slice(targets, func(i, j int) bool {
		if byTarget[targets[i]] != byTarget[targets[j]] {
			return byTarget[targets[i]] > byTarget[targets[j]]
		}
		return targets[i] < targets[j]
	})

	sort.Slice(violations, func(i, j int) bool {
		if violations[i].file != violations[j].file {
			return violations[i].file < violations[j].file
		}
		return violations[i].line < violations[j].line
	})

	var b strings.Builder
	b.WriteString("ADR-023 §2 module boundary violated\n\n")
	b.WriteString("by origin layer:\n")
	for _, origin := range origins {
		b.WriteString("  " + strconv.Itoa(byOrigin[origin]) + "\t" + origin + "\n")
	}
	b.WriteString("\nby target:\n")
	for _, target := range targets {
		b.WriteString("  " + strconv.Itoa(byTarget[target]) + "\t" + target + "\n")
	}
	b.WriteString("\nsites:\n")
	for _, v := range violations {
		from := v.fromModule
		if v.fromLayer != "" {
			from += "/" + v.fromLayer
		}
		to := v.toModule
		if v.toLayer != "" {
			to += "/" + v.toLayer
		}
		b.WriteString("  " + v.file + ":" + strconv.Itoa(v.line) + "\t" + from + " -> " + to + "\n")
	}

	t.Errorf("%d violation(s)\n%s", len(violations), b.String())
}
