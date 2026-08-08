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

// CountFilesSuffix returns how many files the detectors would parse under root.
// It uses the same walk, so the number is the detectors' file set and not an
// estimate of it. The caller that prints findings must also print this and fail
// on zero: until 2026-08-08 archscan reported nothing but the findings, so a
// root with no Go files in it and a clean tree were byte-identical -- the
// default root never covered cmd/ (11 files) or tests/ (35 files), and nothing
// could say so.
func CountFilesSuffix(root, suffix string) (int, error) {
	n := 0
	err := walk(root, suffix, func(string, *token.FileSet, *ast.File) error {
		n++
		return nil
	})
	return n, err
}

// CountFiles counts real Go files.
func CountFiles(root string) (int, error) {
	return CountFilesSuffix(root, ".go")
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

// vendorAwareLayers is the declared list of path segments whose files may name a
// vendor, and it is a declaration rather than a guess. Until 2026-08-08 the only
// exemption was adapters/, re-derived from the path by substring, and the gap
// between that guess and the architecture produced 46 findings in production
// code that were all the design working as intended. The layers that know
// vendors are three, not one:
//
//   - adapters/: the rule's own definition. Vendor payloads live here.
//   - composition/: a composition root exists to name concrete adapters and wire
//     them; internal/composition and the modules' per-module composition/
//     schedulers both qualify. 31 of the 46 were here.
//   - marketplaces/registry/: the plugin registry, one file per marketplace,
//     whose entire job is to say "mercado_livre is this plugin".
//
// The two instrument trees are exempt for the same reason as before: a closed
// token list is the detector, not a violation, and archguard's testdata holds
// fabricated violations for a different checker.
var vendorAwareLayers = []string{
	"/adapters/",
	"/composition/",
	"/marketplaces/registry/",
	"/internal/arch/",
	"/platform/archguard/",
	// A probe utility whose identity is the vendor: it exists to call the
	// Mercado Livre API and record what came back. Same class as a registry
	// file named after its marketplace.
	"/cmd/mlprobe/",
}

// vendorRuleApplies reports whether Regra 2.3 governs this file.
//
// _test.go is exempt for THIS rule, and only this one, and the package comment
// explaining why walk() skips nothing still stands -- for imports. A test that
// imports another context's internals is the defect itself, which is why the
// cross-context detector reads tests. A test that says "mercado_livre" is
// fabricating a fixture: vendor identity as a value fed through a port, which
// is the system working. Measured 2026-08-08 with segment matching already in
// force: 470 of 478 vendor findings were in _test.go, every one a fixture
// value, zero of them control flow. A rule whose report is 98% fixtures is an
// instrument nobody reads; the production site a bad test would mirror is the
// one this rule reports.
func vendorRuleApplies(path string) bool {
	p := "/" + strings.TrimPrefix(filepath.ToSlash(path), "/")
	// TrimSuffix first, so a fixture named *_test.go.txt exercises this branch
	// the same way a real test file does.
	if strings.HasSuffix(strings.TrimSuffix(p, ".txt"), "_test.go") {
		return false
	}
	for _, layer := range vendorAwareLayers {
		if strings.Contains(p, layer) {
			return false
		}
	}
	return true
}

// wordSegments splits text into lowercase word segments: underscores and any
// other non-letter rune are boundaries, and so is a camel-case hump. "Timeline"
// is one segment, "MercadoLivrePlugin" is three, "shopee_partner" is two.
func wordSegments(text string) []string {
	var segs []string
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			segs = append(segs, strings.ToLower(string(cur)))
			cur = nil
		}
	}
	runes := []rune(text)
	for i, r := range runes {
		isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		if !isLetter {
			flush()
			continue
		}
		upper := r >= 'A' && r <= 'Z'
		if upper && len(cur) > 0 {
			prev := runes[i-1]
			prevLower := prev >= 'a' && prev <= 'z'
			nextLower := i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z'
			// aB is a boundary; so is the B of ABc (the acronym's last letter
			// starts the next word). ABC stays one segment.
			if prevLower || nextLower {
				flush()
			}
		}
		cur = append(cur, r)
	}
	flush()
	return segs
}

// matchVendorToken reports the first token that appears in text as a whole word
// or a run of consecutive words, never as a substring of one. The distinction
// is the difference between an instrument and noise, and it was measured before
// it was written: on 2026-08-08 substring matching reported 64 findings for
// "meli", and 64 of 64 were the letters m-e-l-i inside "Timeline". Zero real.
//
// Tokens are normalised the same way ("mercado_livre" and "mercadolivre" are one
// token), so "MercadoLivrePlugin", "mercado_livre" and "mercadolivre" all match
// while "Timeline" and "read listing timeline: %w" never do.
func matchVendorToken(tokens []string, text string) (string, bool) {
	segs := wordSegments(text)
	for _, t := range tokens {
		want := strings.Join(wordSegments(t), "")
		for i, seg := range segs {
			// A vendor name concatenated into one conventional lowercase
			// identifier -- shopeeclient, mercadolivreconnector -- is a single
			// segment longer than the token, and run-equality below walks past
			// it (found in review of PR #22, reproduced: two such vars scanned
			// as findings=0). Prefix or suffix of the segment, never interior,
			// so "timeline" stays silent; and only for tokens of 6+ letters,
			// so "meli" cannot claim "melissa" -- a short token inside a longer
			// word is the Timeline false positive wearing a new spelling.
			if len(want) >= 6 && len(seg) > len(want) &&
				(strings.HasPrefix(seg, want) || strings.HasSuffix(seg, want)) {
				return t, true
			}
			acc := ""
			for j := i; j < len(segs); j++ {
				acc += segs[j]
				if len(acc) > len(want) {
					break
				}
				if acc == want {
					return t, true
				}
			}
		}
	}
	return "", false
}

// ScanVendorTokensSuffix reports every vendor name appearing in an identifier
// or a string literal under root. Both, because the measured defect is a
// literal and an identifier-only detector walks past it. Literals are matched
// wherever they appear, not only in comparisons: the same defect ships as
// `providerCode != "mercado_livre"` in ingest_service.go:334 and as a hoisted
// `const defaultChannel = "mercado_livre"` compared by name, and a detector
// keyed on the comparison sees the first and walks past the second.
func ScanVendorTokensSuffix(root string, tokens []string, suffix string) (Findings, error) {
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		if !vendorRuleApplies(path) {
			return nil
		}
		report := func(pos token.Pos, where, text string) {
			if t, ok := matchVendorToken(tokens, text); ok {
				out = append(out, Finding{
					File:   filepath.ToSlash(path),
					Line:   fset.Position(pos).Line,
					Rule:   RuleVendorToken,
					Detail: t + " in " + where,
				})
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
