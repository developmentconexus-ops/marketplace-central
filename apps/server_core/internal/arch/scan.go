// Package arch holds the level-2 instruments: the rules the Go compiler cannot
// express by itself. It is deliberately NOT a _test.go file.
//
// The instrument that opened Wave 1 lived in a test, skipped every _test.go it
// walked, and hid 100 cross-module imports across 67 files. Living in normal
// code means these detectors are themselves testable against fabricated
// violations, which is the only way to know that a green run means "nothing to
// find" and not "nothing was looked at".
package arch

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Rule names, used in findings and in the gate's output.
const (
	RuleCrossContextInternal = "context/internal-import"
	RuleFloatInContracts     = "numbers/float-in-contracts"
	RuleVendorToken          = "adapters/vendor-token-outside-adapters"
	RuleFactValueDiscard     = "facts/value-discarded"
)

// VendorTokens is the closed list of marketplace names that may not appear
// outside adapters/. A list and not a regex, because a regex over an import
// path is blind in ways a list is not.
var VendorTokens = []string{
	"mercado_livre", "mercadolivre", "mercadolibre", "meli",
	"shopee", "amazon", "magalu", "americanas",
}

const modulePrefix = "marketplace-central/apps/server_core/"

// Finding is one violation, at a file and a line.
type Finding struct {
	File   string
	Line   int
	Rule   string
	Detail string
}

// Findings is a list of violations, sorted by file then line.
type Findings []Finding

func (f Findings) sortInPlace() {
	sort.Slice(f, func(i, j int) bool {
		if f[i].File != f[j].File {
			return f[i].File < f[j].File
		}
		return f[i].Line < f[j].Line
	})
}

// walk parses every file under root whose name ends in suffix. Nothing is
// skipped by name: a _test.go file is code and lives under the same rules.
func walk(root, suffix string, visit func(path string, fset *token.FileSet, file *ast.File) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "node_modules", ".git", ".gocache":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), suffix) {
			return nil
		}
		src, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		fset := token.NewFileSet()
		parsed, parseErr := parser.ParseFile(fset, path, src, parser.ParseComments)
		if parseErr != nil {
			return parseErr
		}
		return visit(path, fset, parsed)
	})
}

// contextOf returns the context name for a path under .../contexts/NAME/... .
func contextOf(path string) (string, bool) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i, p := range parts {
		if p == "contexts" && i+1 < len(parts) {
			return parts[i+1], true
		}
	}
	return "", false
}

// importedContextInternal returns the context whose internals an import path
// reaches into, and whether it reaches into any.
func importedContextInternal(importPath string) (string, bool) {
	rest, ok := strings.CutPrefix(importPath, modulePrefix+"internal/contexts/")
	if !ok {
		return "", false
	}
	parts := strings.Split(rest, "/")
	if len(parts) < 2 {
		return "", false
	}
	for _, p := range parts[1:] {
		if p == "internal" {
			return parts[0], true
		}
	}
	return "", false
}

// ScanCrossContextInternalSuffix reports every file inside one context that
// imports another context's internal packages.
//
// The Go toolchain already refuses this at compile time. This detector exists
// for the window in which a context is being moved, and for the gate's report:
// a number a human can read is worth more than a build error nobody records.
// It must never be the only enforcement.
func ScanCrossContextInternalSuffix(root, suffix string) (Findings, error) {
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		// here is "" for a file that lives outside any context — the composition
		// root, an adapter, a cmd. Those are NOT skipped: the one import that
		// ever broke this rule came from exactly there, and a detector that
		// starts by skipping them can only ever report zero.
		here, _ := contextOf(path)
		for _, imp := range file.Imports {
			value, uErr := strconv.Unquote(imp.Path.Value)
			if uErr != nil {
				continue
			}
			target, reaches := importedContextInternal(value)
			if !reaches {
				continue
			}
			// A context reaching into its OWN internal is the design; anything
			// else, including here == "", is a finding.
			if here != "" && target == here {
				continue
			}
			from := here
			if from == "" {
				from = "outside any context"
			}
			out = append(out, Finding{
				File:   filepath.ToSlash(path),
				Line:   fset.Position(imp.Pos()).Line,
				Rule:   RuleCrossContextInternal,
				Detail: from + " imports " + value,
			})
		}
		return nil
	})
	out.sortInPlace()
	return out, err
}

// ScanCrossContextInternal scans real Go files.
func ScanCrossContextInternal(root string) (Findings, error) {
	return ScanCrossContextInternalSuffix(root, ".go")
}

// ScanFloatInContractsSuffix reports every float type in a contracts package.
//
// A published contract carrying float64 promises the number survives the round
// trip, and it does not: 0.1 has no binary representation, and a tax base built
// from one is wrong before any rule is applied. 51 such fields were measured in
// the module tree this replaces.
func ScanFloatInContractsSuffix(root, suffix string) (Findings, error) {
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		if !strings.Contains(filepath.ToSlash(path), "/contracts/") {
			return nil
		}
		ast.Inspect(file, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if ident.Name != "float64" && ident.Name != "float32" {
				return true
			}
			out = append(out, Finding{
				File:   filepath.ToSlash(path),
				Line:   fset.Position(ident.Pos()).Line,
				Rule:   RuleFloatInContracts,
				Detail: ident.Name + " in a published contract",
			})
			return true
		})
		return nil
	})
	out.sortInPlace()
	return out, err
}

// ScanFloatInContracts scans real Go files.
func ScanFloatInContracts(root string) (Findings, error) {
	return ScanFloatInContractsSuffix(root, ".go")
}

// vendorRuleApplies reports whether Regra 2.3 governs this file.
//
// The rule is "a vendor name does not appear OUTSIDE adapters/", so adapters are
// exempt by definition. The scanner package is exempt too: its closed token list
// is the instrument, not a violation, and a detector that permanently accuses
// itself is a detector nobody can ever act on.
func vendorRuleApplies(path string) bool {
	p := filepath.ToSlash(path)
	if strings.Contains(p, "/adapters/") || strings.HasPrefix(p, "adapters/") {
		return false
	}
	if strings.Contains(p, "/internal/arch/") || strings.HasPrefix(p, "internal/arch/") {
		return false
	}
	return true
}

// ScanVendorTokensSuffix reports every vendor name appearing in an identifier
// or a string literal under root. Both, because the measured defect is a
// literal and an identifier-only detector walks past it.
func ScanVendorTokensSuffix(root string, tokens []string, suffix string) (Findings, error) {
	lower := make([]string, len(tokens))
	for i, t := range tokens {
		lower[i] = strings.ToLower(t)
	}
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		if !vendorRuleApplies(path) {
			return nil
		}
		report := func(pos token.Pos, where, text string) {
			hay := strings.ToLower(text)
			for _, t := range lower {
				if strings.Contains(hay, t) {
					out = append(out, Finding{
						File:   filepath.ToSlash(path),
						Line:   fset.Position(pos).Line,
						Rule:   RuleVendorToken,
						Detail: t + " in " + where,
					})
					return
				}
			}
		}
		ast.Inspect(file, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.Ident:
				report(v.Pos(), "identifier", v.Name)
			case *ast.BasicLit:
				if v.Kind == token.STRING {
					report(v.Pos(), "string literal", v.Value)
				}
			}
			return true
		})
		return nil
	})
	out.sortInPlace()
	return out, err
}

// ScanVendorTokens scans real Go files.
func ScanVendorTokens(root string, tokens []string) (Findings, error) {
	return ScanVendorTokensSuffix(root, tokens, ".go")
}

// ScanFactValueDiscardSuffix reports every `v, _ := something.Value()`.
//
// The blank is the whole defect. Fact.Value returns (T, bool) precisely so that
// an Unknown cannot be read as a value; discarding the bool hands back the zero
// value of T, which is the mistake the fact package exists to prevent. This is
// a syntactic detector on purpose: it cannot know the receiver's type without a
// type checker, so it reports every two-result .Value() call whose second result
// is discarded. A call site that legitimately does that renames its method.
func ScanFactValueDiscardSuffix(root, suffix string) (Findings, error) {
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		ast.Inspect(file, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok || len(assign.Lhs) != 2 || len(assign.Rhs) != 1 {
				return true
			}
			blank, ok := assign.Lhs[1].(*ast.Ident)
			if !ok || blank.Name != "_" {
				return true
			}
			call, ok := assign.Rhs[0].(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Value" {
				return true
			}
			out = append(out, Finding{
				File:   filepath.ToSlash(path),
				Line:   fset.Position(assign.Pos()).Line,
				Rule:   RuleFactValueDiscard,
				Detail: "the bool from .Value() is discarded: unknown would read as the zero value",
			})
			return true
		})
		return nil
	})
	out.sortInPlace()
	return out, err
}

// ScanFactValueDiscard scans real Go files.
func ScanFactValueDiscard(root string) (Findings, error) {
	return ScanFactValueDiscardSuffix(root, ".go")
}
