// Package archguard is the F-04-read-guard-allowlist architectural guard
// (milestone M-02-sync-core-seam, mission MIS-007 ml-sync).
//
// It enumerates every interactive-request-time code path that can reach the
// sole Mercado Livre HTTP client
// (internal/modules/connectors/adapters/mercado_livre/capability_adapter.go,
// defaultBaseURL = "https://api.mercadolibre.com") as an explicit, reviewable
// allowlist. The allowlist lives as Go data in THIS test file (never an
// external json/yaml sidecar), so changing it is a normal, reviewable code
// diff.
//
// Scope: interactive request-time handlers only. Batch/background paths are
// explicitly out of scope and are never walked by this guard:
//
//   - POST /orders/import, POST /listings/refresh,
//     POST /product-links/listing-snapshots/imports,
//     /link-candidates/generations, /pricing/simulations/batch,
//     /admin/fee-schedules/* -- all registered httpx.BatchRouteClass via
//     registerBatchRoutes in internal/composition/root.go. Calling ML from a
//     batch job is by design, not a violation.
//   - POST /market/collections and other market/* wiring -- a different
//     mission's scope (MIS-008); see mlExcludedSymbols below for the one
//     adjacent call site this guard's raw detector structurally matches but
//     must not flag.
//   - /integrations/installations/{id}/probes/* integration probes.
//   - Background schedulers/pollers (FeeSyncScheduler, RefreshTicker,
//     mutationsbg.Poller).
//   - GET /listings does not call ML at all today (repo-backed only), so it
//     has nothing to allowlist.
//
// Mechanism: the sole ML client package is only ever imported directly by
// internal/composition (root.go, orders_adapters.go, market_adapters.go) --
// never by orders/ or pricing/ transport or application code, which depend
// only on ports (interfaces). That is verified separately by
// TestRealRepoTransportAndApplication_NeverImportMLDirectly below. Because of
// that, every interactive path to ML necessarily shows up as a composition
// wiring call in root.go: a locally-defined (bare-identifier, unqualified)
// "newXxxAdapter(...)" constructor call that receives one of the ML
// capability values (mercadoLivreCapabilities, or feeReader sourced from the
// marketplace capability registry keyed by the ML provider code) as a direct
// argument. scanWiringFile detects exactly that shape, by symbol name, not by
// a fixed list of expected names -- so a brand-new constructor introducing a
// brand-new 5th site is detected too, not just growth of the 4 known ones.
package archguard

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Site is one allowlisted interactive-request-time code path permitted to
// reach the Mercado Livre client.
type Site struct {
	// Name is a short human label: <handler-area>.<endpoint>.<role>.
	Name string
	// File is the composition wiring file (relative to apps/server_core)
	// where the ML-capability value is hand to the constructor.
	File string
	// Symbol is the local (unqualified) composition constructor function
	// invoked at that wiring site -- the literal identifier scanWiringFile
	// detects by parsing File's AST.
	Symbol string
}

// mlAllowlist is the CURRENT, exhaustive list of interactive-request-time
// code paths permitted to reach the Mercado Livre client. Verified
// file:line facts (2026-08, MIS-007 F-04-read-guard-allowlist research;
// A/B retired 2026-08, MIS-007 M-03 F-03-read-path-switch):
//
//	A - RETIRED (F-03). GET /orders (list) shipment enrichment used to reach
//	    ML live via composition/orders_adapters.go's ordersShipmentReaderAdapter
//	    (newOrdersShipmentReaderAdapter). It now reads order_shipments
//	    (migration 0088, populated by F-02's IngestOrder) through
//	    orders/adapters/postgres.ShipmentReader -- no ML capability argument at
//	    the wiring site, so the symbol no longer appears in root.go at all
//	    (deleted, not merely excluded).
//	B - RETIRED from this allowlist, RECLASSIFIED as batch-only (F-03).
//	    GET /orders/{id} (detail) buyer-fiscal enrichment used to reach ML live
//	    the same way (newOrdersBuyerFiscalReaderAdapter). The INTERACTIVE path
//	    now reads the buyer_* columns on orders_marketplace_orders (migration
//	    0089) through orders/adapters/postgres.BuyerFiscalReader. The
//	    newOrdersBuyerFiscalReaderAdapter constructor itself is NOT deleted --
//	    it still backs ordersIngestSvc.BuyerFiscal for POST /orders/import
//	    batch ingest (root.go, F-02) -- so the raw detector still finds it; it
//	    moves to mlExcludedSymbols below with that reason, replacing this
//	    allowlist entry (not simply added alongside it).
//	C - POST /pricing/decompose: pricing/transport/calc_handler.go
//	    handleDecompose -> live tariff resolver -> composition wiring
//	    (root.go, THIS site) -> ML catalog-match (category resolution).
//	D - POST /pricing/solve: calc_handler.go handleSolve -> live tariff
//	    resolver -> composition wiring (root.go, THIS site) -> ML fee-quote
//	    (commission quoting).
//
// This is the DONE criterion of F-04: if a 5th interactive site is ever
// wired the same way, TestRealRepoInteractiveMLSites_MatchesAllowlist fails
// and names the offending file:line:symbol (see
// TestFixture_FifthSiteIsDetectedAndNamed for the proof this is not
// vacuous). If C-D is later reduced, shrinking this slice is all that is
// required -- see TestFixture_ShrunkAllowlistStillPasses.
var mlAllowlist = []Site{
	{Name: "pricing.decompose.category_resolver", File: "internal/composition/root.go", Symbol: "newPricingCategoryResolverAdapter"},
	{Name: "pricing.solve.commission_quoter", File: "internal/composition/root.go", Symbol: "newPricingCommissionQuoterAdapter"},
}

// mlExcludedSymbols documents composition-local constructor calls that
// structurally match scanWiringFile's detector (a bare local
// "newXxxAdapter(...)" call receiving one of mlCapabilityVars as a direct
// argument) but are explicitly out of THIS guard's scope, each with a
// reviewable reason. An undocumented exclusion would defeat the guard, so
// every entry here must cite why, and
// TestRealRepoInteractiveMLSites_MatchesAllowlist asserts the raw detector
// really does find each excluded symbol (i.e. this list is not dead code).
var mlExcludedSymbols = map[string]string{
	"newMarketPriceIntelCollectorAdapter": "market/* competitor-price collection is MIS-008 scope, not MIS-007 F-04; tracked separately",
	// The 3 entries below are F-02/F-03 (MIS-007 M-03) batch-only ingest wiring:
	// POST /orders/import is httpx.BatchRouteClass (registerBatchRoutes), out of
	// this guard's interactive-only scope per the package doc comment.
	"newOrdersOrderDetailReaderAdapter":    "batch-only: feeds ordersIngestSvc.OrderDetail for POST /orders/import (F-02), not an interactive request-time path",
	"newOrdersShipmentDetailReaderAdapter": "batch-only: feeds ordersIngestSvc.ShipmentDetail for POST /orders/import (F-02), not an interactive request-time path",
	"newOrdersBuyerFiscalReaderAdapter":    "RECLASSIFIED by F-03 (was allowlist entry orders.detail.buyer_fiscal): the interactive GET /orders/{id} path no longer calls this constructor's result -- it now reads Postgres via orders/adapters/postgres.BuyerFiscalReader. The remaining root.go call site feeds ordersIngestSvc.BuyerFiscal for POST /orders/import batch ingest only (F-02), same batch-only reasoning as the two entries above",
}

// mlCapabilityVars are the local identifiers in root.go bound to a concrete
// Mercado Livre capability value: directly (mercadoLivreCapabilities, from
// mercadolivreconnector.NewCapabilityAdapter), or via the marketplace
// capability registry keyed by the ML provider code (feeReader, from
// marketplaceCapabilities.FeeQuoteReader(mercadoLivreProviderCode)). A
// composition constructor call that receives one of these as a direct
// argument is, definitionally, a path that reaches the ML client.
var mlCapabilityVars = map[string]bool{
	"mercadoLivreCapabilities": true,
	"feeReader":                true,
}

// mlImportPath is the sole ML HTTP client package's import path, as it
// appears (quoted) in an *ast.ImportSpec. Interactive transport/application
// code must never import it directly -- all reach must be mediated through
// composition, which is what makes an allowlist over composition wiring
// calls exhaustive.
const mlImportPath = `"marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"`

// callSite is one call detected inside a wiring file.
type callSite struct {
	Symbol string
	File   string
	Line   int
}

// resolveCapabilityAliases scans f for simple local variable assignments --
// `ident := capVar`, `ident = capVar`, and their element-wise multi-assign
// form `a, b := x, y` -- and returns capabilityVars extended with the
// transitive closure of every local alias that resolves back to a
// capability identifier through a chain of such assignments.
//
// This closes an adversarial-review counterexample: scanWiringFile used to
// match call arguments by literal identifier name only, so a one-line
// "extract variable" refactor --
//
//	mlAlias := mercadoLivreCapabilities
//	sneakyRefundReader := newSneakyRefundAdapter(mlAlias, installationSvc, tenantID)
//
// -- silently escaped detection: mlAlias is not itself the string
// "mercadoLivreCapabilities", so the raw hit count never changed and
// compareToAllowlist returned nil despite a brand-new interactive site
// reaching the ML capability. See
// TestFixture_AliasedSiteIsDetectedAndNamed / testdata/aliased_site/root.go
// for the reproduction and must-fail proof.
//
// This is deliberately NOT full dataflow/SSA analysis: the capability
// variables (mercadoLivreCapabilities, feeReader) are function-local to
// root.go and cannot be re-exported or threaded through anything other than
// a simple local assignment without going through composition's own
// constructor calls (which is what the allowlist already governs). A
// single-hop-per-edge fixpoint over plain-Ident-to-plain-Ident assignments
// is exactly what's needed to close that one hole, and no more.
func resolveCapabilityAliases(f *ast.File, capabilityVars map[string]bool) map[string]bool {
	edges := map[string]string{} // lhs ident name -> rhs ident name (single hop)
	ast.Inspect(f, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != len(assign.Rhs) {
			return true
		}
		for i := range assign.Lhs {
			lhs, ok := assign.Lhs[i].(*ast.Ident)
			if !ok || lhs.Name == "_" {
				continue
			}
			rhs, ok := assign.Rhs[i].(*ast.Ident)
			if !ok {
				continue
			}
			edges[lhs.Name] = rhs.Name
		}
		return true
	})

	resolved := make(map[string]bool, len(capabilityVars))
	for k, v := range capabilityVars {
		resolved[k] = v
	}

	// Fixpoint closure over the file-local edge set (small: composition
	// wiring functions, not arbitrary programs): repeatedly promote any
	// alias whose RHS is already known to be capability-bearing, so chains
	// like `a := mercadoLivreCapabilities; b := a` resolve transitively too.
	for changed := true; changed; {
		changed = false
		for lhs, rhs := range edges {
			if resolved[rhs] && !resolved[lhs] {
				resolved[lhs] = true
				changed = true
			}
		}
	}
	return resolved
}

// scanWiringFile parses a single Go source file (syntax only -- no
// type-checking, no go/packages, no module resolution) and returns every
// bare, locally-defined function call whose argument list directly contains
// one of capabilityVars, OR a local alias of one of capabilityVars per
// resolveCapabilityAliases.
//
// It works identically on real repo files and on testdata/ fixtures, since
// fixtures need not compile or type-check to be parsed.
//
// "Bare, locally-defined" (call.Fun is a plain *ast.Ident, not a qualified
// selector such as pkg.Func(...) or recv.Method(...)) is deliberate: it is
// exactly the convention this codebase uses for composition-root adapter
// constructors (newOrdersShipmentReaderAdapter, newPricingCommissionQuoterAdapter,
// ...). It also structurally EXCLUDES capability-var uses that are not
// adapter wiring. Two such uses are real, present-day occurrences in
// root.go, and neither is flagged by this rule:
//
//   - mercadoLivreCapabilities.ProviderCapabilitySet() -- a method call, so
//     call.Fun is a *ast.SelectorExpr, not a bare Ident.
//   - `CatalogMatch: mercadoLivreCapabilities` inside a
//     ProviderOperationServiceConfig{} literal passed to
//     integrationsapp.NewProviderOperationService(...) -- the OUTER call's
//     Fun is a qualified external-package selector, so the capability
//     identifier being nested two levels down inside a composite literal
//     never matters.
//
// TestRealRepoInteractiveMLSites_MatchesAllowlist below asserts the exact
// count this produces against the real repo, which is what proves those two
// adjacent uses are correctly excluded (not accidentally, via an unmatched
// heuristic).
func scanWiringFile(path string, capabilityVars map[string]bool) ([]callSite, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	resolvedVars := resolveCapabilityAliases(f, capabilityVars)

	var sites []callSite
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		fn, ok := call.Fun.(*ast.Ident)
		if !ok {
			return true // qualified/selector call: not a local constructor
		}
		for _, arg := range call.Args {
			id, ok := arg.(*ast.Ident)
			if !ok {
				continue
			}
			if resolvedVars[id.Name] {
				sites = append(sites, callSite{
					Symbol: fn.Name,
					File:   path,
					Line:   fset.Position(call.Pos()).Line,
				})
				break
			}
		}
		return true
	})
	return sites, nil
}

// filterExcluded drops any callSite whose Symbol is documented in excluded.
func filterExcluded(sites []callSite, excluded map[string]string) []callSite {
	var out []callSite
	for _, s := range sites {
		if _, skip := excluded[s.Symbol]; skip {
			continue
		}
		out = append(out, s)
	}
	return out
}

// compareToAllowlist checks that the discovered call sites are EXACTLY the
// allowlist, matched by Symbol (stable across line-number churn, unlike
// Line). It returns a descriptive error that names the exact offending
// file:line:symbol on any mismatch -- deliberately never a bare "count
// mismatch: expected 4 got 5", per the brief's warning against vacuous
// guards.
func compareToAllowlist(sites []callSite, allowlist []Site) error {
	bySymbol := make(map[string]callSite, len(sites))
	for _, s := range sites {
		if prev, dup := bySymbol[s.Symbol]; dup {
			return fmt.Errorf("archguard: symbol %q reaches the ML capability from more than one call site (%s:%d and %s:%d) -- allowlist entries must be 1:1 with call sites",
				s.Symbol, prev.File, prev.Line, s.File, s.Line)
		}
		bySymbol[s.Symbol] = s
	}

	allowedBySymbol := make(map[string]Site, len(allowlist))
	for _, a := range allowlist {
		allowedBySymbol[a.Symbol] = a
	}

	var unexpected []string
	for sym, s := range bySymbol {
		if _, ok := allowedBySymbol[sym]; !ok {
			unexpected = append(unexpected, fmt.Sprintf(
				"%s:%d: unallowlisted call to %s(...) reaches the Mercado Livre capability -- add it to mlAllowlist, or to mlExcludedSymbols with a documented reason if it is legitimately out of scope",
				s.File, s.Line, sym))
		}
	}
	sort.Strings(unexpected)

	var missing []string
	for sym, a := range allowedBySymbol {
		if _, ok := bySymbol[sym]; !ok {
			missing = append(missing, fmt.Sprintf(
				"allowlisted site %q (%s in %s) was not found -- either it was removed (shrink mlAllowlist to match) or the scanner regressed",
				a.Name, sym, a.File))
		}
	}
	sort.Strings(missing)

	if len(unexpected) == 0 && len(missing) == 0 {
		return nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "archguard: interactive-path Mercado Livre allowlist mismatch (found %d site(s), allowlist has %d entr(y/ies)):\n", len(bySymbol), len(allowlist))
	for _, m := range unexpected {
		fmt.Fprintf(&b, "  NEW SITE: %s\n", m)
	}
	for _, m := range missing {
		fmt.Fprintf(&b, "  MISSING:  %s\n", m)
	}
	return errors.New(b.String())
}

// scanDirectImports walks root (excluding testdata/ and _test.go files) and
// returns every .go file that directly imports targetImportPath (quoted, as
// it appears in an ast.ImportSpec.Path.Value).
func scanDirectImports(root string, targetImportPath string) ([]string, error) {
	var hits []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		f, perr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if perr != nil {
			return fmt.Errorf("parse %s: %w", path, perr)
		}
		for _, imp := range f.Imports {
			if imp.Path.Value == targetImportPath {
				hits = append(hits, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(hits)
	return hits, nil
}

// --- Baseline (state 1: must PASS today) ------------------------------------

// TestRealRepoInteractiveMLSites_MatchesAllowlist is the actual guard: it
// scans the real internal/composition/root.go and asserts the interactive
// Mercado Livre call sites it finds are exactly the 4-entry mlAllowlist --
// no more, no fewer.
func TestRealRepoInteractiveMLSites_MatchesAllowlist(t *testing.T) {
	const wiringFile = "../../composition/root.go"

	raw, err := scanWiringFile(wiringFile, mlCapabilityVars)
	if err != nil {
		t.Fatalf("scan real root.go: %v", err)
	}

	// Prove mlExcludedSymbols is not dead code: the raw (pre-exclusion)
	// detector must actually find every documented exclusion in the real
	// file, otherwise the exclusion is vacuous.
	rawBySymbol := map[string]bool{}
	for _, s := range raw {
		rawBySymbol[s.Symbol] = true
	}
	for sym, reason := range mlExcludedSymbols {
		if !rawBySymbol[sym] {
			t.Fatalf("mlExcludedSymbols documents %q (%s) as a raw detector hit that must be filtered, but the raw scan of %s no longer finds it -- remove the now-dead exclusion", sym, reason, wiringFile)
		}
	}
	if len(raw) != len(mlAllowlist)+len(mlExcludedSymbols) {
		t.Fatalf("raw detector found %d call site(s) in %s, expected exactly %d (allowlist) + %d (documented exclusions) = %d; a new, undocumented site appeared -- raw hits: %+v",
			len(raw), wiringFile, len(mlAllowlist), len(mlExcludedSymbols), len(mlAllowlist)+len(mlExcludedSymbols), raw)
	}

	sites := filterExcluded(raw, mlExcludedSymbols)
	if err := compareToAllowlist(sites, mlAllowlist); err != nil {
		t.Fatalf("baseline allowlist mismatch:\n%v", err)
	}

	t.Logf("baseline PASS: %d interactive Mercado Livre site(s) in %s match mlAllowlist exactly (%d documented exclusion(s) filtered):", len(sites), wiringFile, len(mlExcludedSymbols))
	for _, s := range sites {
		t.Logf("  %-40s %s:%d", s.Symbol, s.File, s.Line)
	}
}

// TestRealRepoTransportAndApplication_NeverImportMLDirectly is the structural
// premise the allowlist relies on: orders/ and pricing/ transport and
// application code depend only on ports (interfaces), never on the ML client
// package directly. If that ever stops being true -- a new handler or
// service importing mercado_livre straight, bypassing composition -- this
// fails immediately by naming the file, independent of (and in addition to)
// the composition wiring scan above.
func TestRealRepoTransportAndApplication_NeverImportMLDirectly(t *testing.T) {
	for _, dir := range []string{"../../modules/orders", "../../modules/pricing"} {
		hits, err := scanDirectImports(dir, mlImportPath)
		if err != nil {
			t.Fatalf("scan %s for direct ML imports: %v", dir, err)
		}
		if len(hits) != 0 {
			t.Fatalf("%s: found file(s) importing the Mercado Livre client directly, bypassing composition (new interactive site not mediated by an allowlisted wiring call): %v", dir, hits)
		}
		t.Logf("%s: 0 direct Mercado Livre imports (as expected)", dir)
	}
}

// --- Must-fail proof, state 2: a 5th site is introduced ---------------------

// TestFixture_FifthSiteIsDetectedAndNamed proves the guard can actually go
// red. testdata/five_sites/root.go mirrors the 4 real wiring calls plus one
// INJECTED 5th (newOrdersRefundReaderAdapter) that is not in mlAllowlist.
// This does not touch any real file (orders/, pricing/, and composition/
// root.go itself are forbidden paths for this feature) -- the injection is a
// fixture under testdata/, which go build/vet/test ignore for the real
// binary, and is only ever read here via scanWiringFile.
func TestFixture_FifthSiteIsDetectedAndNamed(t *testing.T) {
	const fixture = "testdata/five_sites/root.go"

	raw, err := scanWiringFile(fixture, mlCapabilityVars)
	if err != nil {
		t.Fatalf("scan fixture %s: %v", fixture, err)
	}
	sites := filterExcluded(raw, mlExcludedSymbols)

	err = compareToAllowlist(sites, mlAllowlist)
	if err == nil {
		t.Fatalf("expected the injected 5th site in %s to fail against mlAllowlist, but the guard passed -- it is not actually catching new sites", fixture)
	}

	const offendingSymbol = "newOrdersRefundReaderAdapter"
	if !strings.Contains(err.Error(), offendingSymbol) {
		t.Fatalf("guard failed as expected, but the failure message does not NAME the offending symbol %q (a generic 'count mismatch' is exactly what the brief forbids); got:\n%v", offendingSymbol, err)
	}
	if !strings.Contains(err.Error(), fixture) {
		t.Fatalf("guard failure message does not name the offending file %q; got:\n%v", fixture, err)
	}

	t.Logf("must-fail PASS (state 2): guard correctly rejected the injected 5th site and named it:\n%v", err)
}

// --- Must-fail proof, state 4: a locally-aliased capability var -------------

// TestFixture_AliasedSiteIsDetectedAndNamed closes the adversarial-review
// alias-bypass: testdata/aliased_site/root.go mirrors the 4 real wiring
// calls plus one INJECTED unallowlisted site (newSneakyRefundAdapter)
// reached via a local alias (`mlAlias := mercadoLivreCapabilities`) instead
// of the capability identifier passed directly. Pre-hardening,
// scanWiringFile matched call arguments by literal identifier name only, so
// this exact one-line "extract variable" shape produced a raw hit count of
// len(mlAllowlist) (the alias site was silently invisible) and
// compareToAllowlist returned nil -- the guard passed green despite the new
// site. resolveCapabilityAliases closes that hole by resolving simple local
// `ident := capVar` / `ident = capVar` assignment chains back to the
// capability identifier before matching, so the alias is treated exactly
// like the direct identifier.
func TestFixture_AliasedSiteIsDetectedAndNamed(t *testing.T) {
	const fixture = "testdata/aliased_site/root.go"

	raw, err := scanWiringFile(fixture, mlCapabilityVars)
	if err != nil {
		t.Fatalf("scan fixture %s: %v", fixture, err)
	}
	sites := filterExcluded(raw, mlExcludedSymbols)

	err = compareToAllowlist(sites, mlAllowlist)
	if err == nil {
		t.Fatalf("expected the injected alias-reached site in %s to fail against mlAllowlist, but the guard passed -- the alias bypass is not caught", fixture)
	}

	const offendingSymbol = "newSneakyRefundAdapter"
	if !strings.Contains(err.Error(), offendingSymbol) {
		t.Fatalf("guard failed as expected, but the failure message does not NAME the offending symbol %q (a generic 'count mismatch' is exactly what the brief forbids); got:\n%v", offendingSymbol, err)
	}
	if !strings.Contains(err.Error(), fixture) {
		t.Fatalf("guard failure message does not name the offending file %q; got:\n%v", fixture, err)
	}

	t.Logf("must-fail PASS (state 4): guard correctly resolved the local alias back to the ML capability, rejected the injected site, and named it:\n%v", err)
}

// --- Must-fail proof, state 3: an allowlist entry shrinks -------------------

// TestFixture_ShrunkAllowlistStillPasses proves the guard is not hardcoded
// to "always exactly 4": testdata/three_sites/root.go mirrors only 3 of the
// 4 real wiring calls (simulating a future milestone -- M-03/M-07 per
// AGENTS.md -- retiring site D, pricing.solve.commission_quoter /
// newPricingCommissionQuoterAdapter), and a 3-entry allowlist (mlAllowlist
// with that one entry removed) is checked against it. If this failed, the
// comparator would be structurally counting rather than matching by symbol.
func TestFixture_ShrunkAllowlistStillPasses(t *testing.T) {
	const fixture = "testdata/three_sites/root.go"
	const droppedSymbol = "newPricingCommissionQuoterAdapter"

	shrunk := make([]Site, 0, len(mlAllowlist)-1)
	for _, s := range mlAllowlist {
		if s.Symbol == droppedSymbol {
			continue
		}
		shrunk = append(shrunk, s)
	}
	if len(shrunk) != len(mlAllowlist)-1 {
		t.Fatalf("test bug: expected to drop exactly one entry (%s) from mlAllowlist, dropped %d", droppedSymbol, len(mlAllowlist)-len(shrunk))
	}

	raw, err := scanWiringFile(fixture, mlCapabilityVars)
	if err != nil {
		t.Fatalf("scan fixture %s: %v", fixture, err)
	}
	sites := filterExcluded(raw, mlExcludedSymbols)

	if err := compareToAllowlist(sites, shrunk); err != nil {
		t.Fatalf("a %d-entry allowlist should pass against a fixture with the corresponding site removed, but got:\n%v", len(shrunk), err)
	}

	t.Logf("must-fail PASS (state 3): shrunk %d-entry allowlist matches the %d-site fixture exactly (guard shrinks structurally, not hardcoded to 4)", len(shrunk), len(sites))
}
