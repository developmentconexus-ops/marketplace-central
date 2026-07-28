## R7 — second false universal + degradation-parity repair

Scope: `apps/server_core/internal/modules/product_links/application/generation_service_test.go` only.
`generation_service.go` was mutated transiently for must-fail proofs and restored by
editing forward each time (never `git checkout`). No file outside the test file
carries a diff introduced by this dispatch.

### 1. FIX 1 — false universal string is gone

```
$ grep -n "unreachable for seller_sku" internal/modules/product_links/application/generation_service_test.go
(no output, exit=1)
```

Replacement text (verbatim, lines 505-511 of the file):

```go
		// Every case in this table runs with a product present (product
		// above always carries a canonical id). The ERP side of seller_sku
		// is the CODPROD, and findProducts drops any candidate without one
		// — so with a product present, side=erp cannot arise for
		// seller_sku here. It arises on the nil-product (unresolved) path
		// instead, pinned by
		// TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide.
		"listing has no seller_sku": {
```

### 2. FIX 2 — "degrades like both of its siblings" is gone

```
$ grep -n "degrades like both of its siblings" internal/modules/product_links/application/generation_service_test.go
(no output, exit=1)
```

Replacement text (verbatim, lines 543-550 of the file):

```go
// TestConcordantCandidateDoesNotDerefNilProduct pins the guard, not a
// production path: the single call site (buildExactCandidates) always passes a
// non-nil pointer, so production reachability is NOT proven. What is proven is
// that this scorer nil-checks the product like both of its siblings, so it no
// longer panics. It does NOT degrade like them: the siblings degrade into
// ABSENCE reasons, while this scorer degrades into CORROBORATION over a
// zeroed ProductCandidate{} — seller_sku and ean FOR reasons at
// Confidence 95, band ALTA, status ACCEPT.
```

The nil-check parity ("nil-checks the product like both of its siblings") is kept
because it is true (`applyUnresolvedScore` / `applySingleAnchorScore` both nil-check
`comparison.product` the same way `buildConcordantCandidate` now does). Only the
degradation-parity claim was narrowed/negated.

### 3. Full diff of `generation_service_test.go`

```diff
diff --git a/apps/server_core/internal/modules/product_links/application/generation_service_test.go b/apps/server_core/internal/modules/product_links/application/generation_service_test.go
index 846b6ac4..62129675 100644
--- a/apps/server_core/internal/modules/product_links/application/generation_service_test.go
+++ b/apps/server_core/internal/modules/product_links/application/generation_service_test.go
@@ -502,8 +502,13 @@ func TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference(t *testing.T)
 			wantSide:         "",
 			wantDetail:       "sem CODPROD para corroborar o EAN: o seller_sku do anúncio não casa nenhum produto",
 		},
-		// side=erp is unreachable for seller_sku now, so the only honest
-		// INCOMPARABLE is the provider side: the ANÚNCIO carries no SKU.
+		// Every case in this table runs with a product present (product
+		// above always carries a canonical id). The ERP side of seller_sku
+		// is the CODPROD, and findProducts drops any candidate without one
+		// — so with a product present, side=erp cannot arise for
+		// seller_sku here. It arises on the nil-product (unresolved) path
+		// instead, pinned by
+		// TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide.
 		"listing has no seller_sku": {
 			wantDirection: productlinksdomain.LinkCandidateReasonDirectionIncomparable,
 			wantSide:      productlinksdomain.LinkCandidateReasonSideProvider,
@@ -538,7 +543,11 @@ func TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference(t *testing.T)
 // TestConcordantCandidateDoesNotDerefNilProduct pins the guard, not a
 // production path: the single call site (buildExactCandidates) always passes a
 // non-nil pointer, so production reachability is NOT proven. What is proven is
-// that this scorer now degrades like both of its siblings instead of panicking.
+// that this scorer nil-checks the product like both of its siblings, so it no
+// longer panics. It does NOT degrade like them: the siblings degrade into
+// ABSENCE reasons, while this scorer degrades into CORROBORATION over a
+// zeroed ProductCandidate{} — seller_sku and ean FOR reasons at
+// Confidence 95, band ALTA, status ACCEPT.
 func TestConcordantCandidateDoesNotDerefNilProduct(t *testing.T) {
 	t.Parallel()
 	now := time.Date(2026, 7, 28, 9, 30, 0, 0, time.UTC)
@@ -561,6 +570,18 @@ func TestConcordantCandidateDoesNotDerefNilProduct(t *testing.T) {
 	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
 		t.Fatalf("reasons=%#v, want the concordant seller_sku FOR reason", candidate.Reasons)
 	}
+	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
+		t.Fatalf("reasons=%#v, want the concordant ean FOR reason", candidate.Reasons)
+	}
+	if candidate.Confidence != 95 {
+		t.Fatalf("confidence=%d, want 95 for a nil-product concordant candidate", candidate.Confidence)
+	}
+	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandAlta {
+		t.Fatalf("confidence_band=%q, want %q", candidate.ConfidenceBand, productlinksdomain.LinkCandidateConfidenceBandAlta)
+	}
+	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusAccept {
+		t.Fatalf("match_status=%q, want %q", candidate.MatchStatus, productlinksdomain.LinkCandidateMatchStatusAccept)
+	}
 }
 
 func TestUnresolvedDeclaredAnchorReadsProviderValueWithoutERPProduct(t *testing.T) {
```

### 4. Must-fail proofs (4 new assertion classes added in FIX 2)

Each mutation was applied to `generation_service.go` (`buildConcordantCandidate`),
run against `TestConcordantCandidateDoesNotDerefNilProduct` only, then reverted by
editing forward (Python string replace, exact inverse of the mutation — never
`git checkout`). `git diff --stat` on `generation_service.go` was captured
immediately after each restore.

#### Mutation A — drop the `ean` FOR reason (proves the new "ean FOR reason present" assertion)

Change: removed the `{Anchor: "ean", Direction: ..., Detail: "ean corrobora o mesmo codprod (unproved)"}`
entry from the `reasons` slice literal.

Raw failure output:
```
=== RUN   TestConcordantCandidateDoesNotDerefNilProduct
=== PAUSE TestConcordantCandidateDoesNotDerefNilProduct
=== CONT  TestConcordantCandidateDoesNotDerefNilProduct
    generation_service_test.go:574: reasons=[]domain.LinkCandidateReason{domain.LinkCandidateReason{Anchor:"seller_sku", Direction:"FOR", Side:"", Detail:"seller_sku resolve exato para codprod"}}, want the concordant ean FOR reason
--- FAIL: TestConcordantCandidateDoesNotDerefNilProduct (0.00s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/product_links/application	1.481s
FAIL
```

Restore: `git diff --stat` on `generation_service.go` → empty (only the pre-existing
warning about LF/CRLF normalization printed; zero content diff).

#### Mutation B — seeded Confidence 95 → 90 (proves the new "Confidence == 95" assertion)

Change: `candidate.Confidence = 95` → `candidate.Confidence = 90` in the non-hard-negative branch.

Raw failure output:
```
=== RUN   TestConcordantCandidateDoesNotDerefNilProduct
=== PAUSE TestConcordantCandidateDoesNotDerefNilProduct
=== CONT  TestConcordantCandidateDoesNotDerefNilProduct
    generation_service_test.go:577: confidence=90, want 95 for a nil-product concordant candidate
--- FAIL: TestConcordantCandidateDoesNotDerefNilProduct (0.00s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/product_links/application	1.104s
FAIL
```

Restore: `git diff --stat` on `generation_service.go` → empty.

#### Mutation C — seeded ConfidenceBand ALTA → MEDIA (proves the new "band == ALTA constant" assertion)

Change: `candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandAlta` →
`domain.LinkCandidateConfidenceBandMedia`.

Raw failure output:
```
=== RUN   TestConcordantCandidateDoesNotDerefNilProduct
=== PAUSE TestConcordantCandidateDoesNotDerefNilProduct
=== CONT  TestConcordantCandidateDoesNotDerefNilProduct
    generation_service_test.go:580: confidence_band="MEDIA", want "ALTA"
--- FAIL: TestConcordantCandidateDoesNotDerefNilProduct (0.00s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/product_links/application	1.085s
FAIL
```

Restore: `git diff --stat` on `generation_service.go` → empty.

#### Mutation D — seeded MatchStatus ACCEPT → CONFIRM (proves the new "status == ACCEPT constant" assertion)

Change: `candidate.MatchStatus = domain.LinkCandidateMatchStatusAccept` →
`domain.LinkCandidateMatchStatusConfirm`.

Raw failure output:
```
=== RUN   TestConcordantCandidateDoesNotDerefNilProduct
=== PAUSE TestConcordantCandidateDoesNotDerefNilProduct
=== CONT  TestConcordantCandidateDoesNotDerefNilProduct
    generation_service_test.go:583: match_status="CONFIRM", want "ACCEPT"
--- FAIL: TestConcordantCandidateDoesNotDerefNilProduct (0.00s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/product_links/application	1.242s
FAIL
```

Restore: `git diff --stat` on `generation_service.go` → empty.

All four mutations made the test fail; none was a no-op. Final `git diff --stat -- generation_service.go`
after all four rounds is empty (no net change from this dispatch's mutation/restore cycles).

### 5. Verification ladder (final state, run from `apps/server_core` with `GOCACHE` set to an absolute path under the worktree)

`go build ./...`:
```
(no output — success)
```

`go vet ./internal/modules/product_links/...`:
```
(no output — success)
```

`go test ./internal/modules/product_links/... -count=1 -v` (PASS lines for the two touched tests):
```
--- PASS: TestConcordantCandidateDoesNotDerefNilProduct (0.00s)
--- PASS: TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference (0.00s)
    --- PASS: TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference/listing_has_a_seller_sku,_ERP_product_has_no_refforn (0.00s)
    --- PASS: TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference/listing_has_no_seller_sku (0.00s)
    --- PASS: TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference/listing_has_a_seller_sku,_ERP_product_has_a_refforn (0.00s)
```
Full package results: `ok marketplace-central/apps/server_core/internal/modules/product_links/application`,
`.../composition`, `.../domain`, `.../transport`, `.../adapters/connectors` all `ok`
(0 skips, 0 fails). No test in the package failed.

`go test ./internal/modules/product_links/... -count=10` (flake guard):
```
ok  	marketplace-central/apps/server_core/internal/modules/product_links/adapters/connectors	2.117s
?   	marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/product_links/application	2.089s
ok  	marketplace-central/apps/server_core/internal/modules/product_links/composition	1.655s
ok  	marketplace-central/apps/server_core/internal/modules/product_links/domain	1.813s
?   	marketplace-central/apps/server_core/internal/modules/product_links/ports	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/product_links/transport	2.510s
```
10x count, no failures.

### 6. What I could not do, and why

- Nothing was blocked. Both fixes landed with only `generation_service_test.go`
  touched, as required. `generation_service.go` shows zero net diff from this
  dispatch (its pre-existing worktree modification, present before this dispatch
  started per the initial `git status`, is untouched by me — every mutation I
  applied to it for the must-fail proofs was reverted to its exact pre-mutation
  content).
- One note, not a blocker: `git status --short` reports `generation_service.go` as
  `M` and `git diff --stat` on it is empty of content changes — this is a
  pre-existing line-ending (LF/CRLF) normalization state in the worktree, present
  before this dispatch began (visible in the harness-supplied initial git status),
  not something introduced here. Git prints a "LF will be replaced by CRLF" warning
  on every command touching that file; this is cosmetic and orthogonal to the R7
  scope.
- I did not touch `apps/web/`, migrations, `platform/httpx`, or
  `generation_service.go`'s actual logic, per the hard rules.
- I did not commit; the changes remain in the working tree for the chip session to
  commit.
