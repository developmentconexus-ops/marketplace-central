# CHIP-ANCHORS-3 — impl R2: the false universal in `generation_service_test.go`

Scope honored: **exactly one file edited** —
`apps/server_core/internal/modules/product_links/application/generation_service_test.go`.

```
$ git diff HEAD --name-only
apps/server_core/internal/modules/product_links/application/generation_service_test.go
```

No other file, no `apps/web/`, no migrations, no `platform/httpx`. No commit, no push.

---

## 1. The comment — old and new, verbatim

Site: `apps/server_core/internal/modules/product_links/application/generation_service_test.go`,
immediately above the table case `"exact EAN listing seller SKU empty"` inside
`TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide`.

### OLD (verbatim, deleted)

```go
		// The ERP side of seller_sku is the CODPROD, which a resolved product
		// always has — so side=erp is reachable for seller_sku only through the
		// listing, never through the product. The mirror case (listing HAS a
		// seller_sku) is pinned by
		// TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference.
```

### NEW (verbatim, in the file now)

```go
		// The ERP side of seller_sku is the CODPROD, and findProducts drops any
		// candidate without one — so with a product present, side=erp cannot
		// arise for seller_sku. It arises on the nil-product (unresolved) path
		// instead, pinned by
		// TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide.
```

Per R-25 the false half was **deleted / narrowed**, not annotated. The universal
("only through the listing, never through the product") is gone; what remains is a
claim scoped to the fixture family this table case actually governs — *a product is
present* — plus a pointer to where the other side is now pinned. The old
cross-reference to `TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference` was
dropped because that test's fixtures also always carry a product, so it is not where
`side=erp` is pinned; the new test is.

The line below the comment (`"exact EAN listing seller SKU empty"`) was **not touched**
and is not narrated.

---

## 2. What I verified in `generation_service.go` — both false halves

All line numbers are current `apps/server_core/internal/modules/product_links/application/generation_service.go`
(the file is byte-identical to HEAD; see §5).

### Half A — "side=erp is reachable for seller_sku **only through the listing**" → FALSE

`missingMatchedAnchorReason`, `generation_service.go:636–650`:

```go
func missingMatchedAnchorReason(anchor, detail string, listing domain.ListingSnapshot, product *internalreaddomain.ProductCandidate) domain.LinkCandidateReason {
	listingValue, productValue, _ := identityAnchorValues(anchor, listing, product)
	reason := domain.LinkCandidateReason{Anchor: anchor, Detail: detail}
	switch {
	case listingValue == "":
		reason.Direction = domain.LinkCandidateReasonDirectionIncomparable
		reason.Side = domain.LinkCandidateReasonSideProvider
	case product == nil || productValue == "":
		reason.Direction = domain.LinkCandidateReasonDirectionIncomparable
		reason.Side = domain.LinkCandidateReasonSideERP
	default:
		reason.Direction = domain.LinkCandidateReasonDirectionUnavailable
	}
	return reason
}
```

- `generation_service.go:640` — exact string `case listingValue == "":`
- `generation_service.go:642` — exact string `reason.Side = domain.LinkCandidateReasonSideProvider`

The **listing-empty** branch is the *first* arm of the switch and sets
`SideProvider`. It can never produce `side=erp`. So `side=erp` for `seller_sku` is not
reachable "through the listing" at all — it is reachable *only* through the second arm.
The claim is inverted.

### Half B — "**never through the product**" → FALSE

- `generation_service.go:643` — exact string `case product == nil || productValue == "":`
- `generation_service.go:645` — exact string `reason.Side = domain.LinkCandidateReasonSideERP`

This arm sets `SideERP`, and `product == nil` **is reached in production**:

- `generation_service.go:216`:
  ```go
	applyUnresolvedScore(&unresolved, newProviderIdentityAnchorComparison(snapshot, identityAnchors, nil))
  ```
  (the `nil` third argument is the ERP product; this is the unresolved path taken when
  seller_sku, EAN and title all match no product — `generation_service.go:215–217`)
- `applyUnresolvedScore` then seeds, at `generation_service.go:630`:
  ```go
		missingMatchedAnchorReason("seller_sku", "seller_sku sem correspondência", comparison.listing, nil),
  ```
  — an unconditional literal `nil` product.

So a listing that carries a `seller_sku` matching nothing yields
`seller_sku / INCOMPARABLE / side=erp / "seller_sku sem correspondência"` in production.
"Never through the product" is false.

Two further call sites pass a nil product on the same unresolved path:
`generation_service.go:341` and `generation_service.go:380`.

### Survival through `appendProviderDeclaredUnavailableReasons` (derived before running)

`applyUnresolvedScore` finishes with `candidate.Reasons = appendProviderDeclaredUnavailableReasons(reasons, comparison)`
(`generation_service.go:633`). The seeded `seller_sku` reason is `INCOMPARABLE`, and no
`FOR`/`AGAINST` reason carries the `seller_sku` anchor, so it is kept and its index
recorded in `absenceIndexes` (`generation_service.go:662–675`). Then the declared-anchor
loop (`:676–703`) calls `classifyProviderIdentityAnchor`. With
`comparison.product == nil` and a non-empty listing value, `generation_service.go:718–723`:

```go
	if comparison.product == nil {
		if listingValue == "" {
			return domain.LinkCandidateReasonDirectionIncomparable, domain.LinkCandidateReasonSideProvider, fmt.Sprintf("anúncio sem %s", anchor.Anchor), true
		}
		return "", "", "", false
	}
```

returns `emit == false`, so the loop `continue`s and the seeded `side=erp` reason
**survives untouched**. Derived first, then confirmed by the run in §6.

Caveat I checked and am recording: this survival depends on the provider declaring
`seller_sku` with `Supplied: true`. If a provider declared `Supplied: false`,
`classifyProviderIdentityAnchor` (`:708–710`) returns `UNAVAILABLE` with `emit == true`
and — because the direction is not `INCOMPARABLE` — the promotion at `:692–699`
**replaces** the seeded reason wholesale, losing `side=erp`. The new test uses a
provider that declares `seller_sku` supplied (`mercadoLivreIdentityAnchorReader`,
test file `:189–198`), which is the real Mercado Livre shape.

---

## 3. Coverage form chosen: a **separate test function**, not a table case

The table `TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide` *could*
mechanically express a scorer-level version of this case (leave `state` at its zero
value so `applySingleAnchorScore` falls through to `applyUnresolvedScore`, leave
`product` nil, set `snapshot.SellerSKU`). I did **not** use it, for two reasons:

1. **Its runner hardcodes `nil` declarations.** Test file line (runner):
   `applySingleAnchorScore(&candidate, tc.state, newProviderIdentityAnchorComparison(tc.snapshot, nil, tc.product))`.
   With `declared == nil` the declared-anchor loop in
   `appendProviderDeclaredUnavailableReasons` never executes, so a table case would
   prove the seed and prove nothing about the pass that can overwrite it. That pass is
   exactly the risk surface identified in §2. Bending the runner to take declarations
   was explicitly out of scope.
2. **The table calls the internal scorer, not the production call site.** The defect
   being repaired is a false claim about *production reachability*. A test that invokes
   `applySingleAnchorScore` directly cannot witness `generation_service.go:216`. The new
   test drives `GenerateLinkCandidates` end to end via the existing `generateSingle`
   helper, so the nil product comes from production code, not from the fixture.

Placement: immediately before `TestProviderDeclaredUnmodelledAnchorIsIncomparableWithoutSide`,
i.e. directly after `TestUnresolvedDeclaredAnchorReadsProviderValueWithoutERPProduct`,
its nearest sibling (same unresolved path, `ean` anchor).

---

## 4. The new test, verbatim

```go
// TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide drives the whole
// generator down the unresolved path, whose production call site passes a nil
// ERP product (generation_service.go: `applyUnresolvedScore(&unresolved,
// newProviderIdentityAnchorComparison(snapshot, identityAnchors, nil))`). A
// listing that carries a seller_sku matching no product has a provider value
// and no ERP counterpart at all, so the honest side is erp — the one direction
// no other seller_sku assertion in this file covers.
//
// The provider here DOES declare seller_sku (mercadoLivreIdentityAnchorReader),
// so this also pins that the declared-anchor pass leaves the seeded side
// intact: classifyProviderIdentityAnchor emits nothing when the product is nil
// and the listing value is present.
func TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 28, 11, 0, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-unres-sku", ProviderCode: "mercado_livre", ProviderItemID: "item-unres-sku",
		SellerSKU: "SKU-SEM-PRODUTO", EAN: "EAN-SEM-PRODUTO", Title: "Produto sem correspondência", FetchedAt: now,
	}
	candidate := generateSingle(t, snapshot, &stubProductMatcher{
		results: map[string][]internalreaddomain.ProductCandidate{},
	}, now)

	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusNoCandidate {
		t.Fatalf("candidate=%#v, want the unresolved (NO_CANDIDATE) path", candidate)
	}
	reason, ok := findReasonByAnchor(candidate.Reasons, "seller_sku")
	if !ok {
		t.Fatalf("reasons=%#v, want a seller_sku reason", candidate.Reasons)
	}
	if reason.Direction != productlinksdomain.LinkCandidateReasonDirectionIncomparable ||
		reason.Side != productlinksdomain.LinkCandidateReasonSideERP ||
		reason.Detail != "seller_sku sem correspondência" {
		t.Fatalf("seller_sku reason=%#v, want direction=%q side=%q detail=%q",
			reason,
			productlinksdomain.LinkCandidateReasonDirectionIncomparable,
			productlinksdomain.LinkCandidateReasonSideERP,
			"seller_sku sem correspondência")
	}
}
```

Notes on the fixture:
- `stubProductMatcher{results: map[...]{}}` returns no product for sku, ean or title, so
  `generateSingles` falls through to `generation_service.go:215–217` — the unresolved path.
- `findReasonByAnchor` (not `findReason`) is used deliberately: it looks the anchor up
  **without** presuming a direction, so a regression reports the direction that actually
  came out instead of an empty "not found". §5 shows that paying off.
- The `MatchStatus` assertion is not decoration: it fails loudly if a future change makes
  this fixture take a different scoring path, so the reason assertion can never pass
  vacuously from some other branch.

First green run, before any mutation:

```
=== RUN   TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide
=== PAUSE TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide
=== CONT  TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide
--- PASS: TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/product_links/application	0.630s
```

---

## 5. Must-fail

### Mutation applied

`generation_service.go:643–645`, the `product == nil` arm of
`missingMatchedAnchorReason`, made to fall through to a different side:

```diff
 	case product == nil || productValue == "":
 		reason.Direction = domain.LinkCandidateReasonDirectionIncomparable
-		reason.Side = domain.LinkCandidateReasonSideERP
+		reason.Side = domain.LinkCandidateReasonSideBoth
```

### VERBATIM failure output (mutated tree)

```
$ cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache \
    go test -count=1 -run 'TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide' -v ./internal/modules/product_links/application/

=== RUN   TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide
=== PAUSE TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide
=== CONT  TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide
    generation_service_test.go:649: seller_sku reason=domain.LinkCandidateReason{Anchor:"seller_sku", Direction:"INCOMPARABLE", Side:"both", Detail:"seller_sku sem correspondência"}, want direction="INCOMPARABLE" side="erp" detail="seller_sku sem correspondência"
--- FAIL: TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide (0.00s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/product_links/application	1.471s
FAIL
```

Real values printed: the actual reason came out `Side:"both"` against the wanted
`side="erp"`, with `Direction` and `Detail` unchanged — i.e. the assertion discriminates
the *side* specifically, not merely the presence of a reason. That is the property the
deleted table case used to hold and the false universal was used to retire.

### Restore, by editing forward

The mutation was reverted with a forward `Edit` (`SideBoth` → `SideERP`). **No**
`git checkout`, `reset`, `revert`, `stash` or `clean` was used at any point.

### Empty-diff confirmation

```
$ git diff HEAD -- apps/server_core/internal/modules/product_links/application/generation_service.go
EXIT=0
--- diff --stat ---
(end)
```

Both the patch body and `--stat` are empty: `generation_service.go` is byte-identical to
HEAD. Corroborated by the repo-wide listing at the top of this report — the test file is
the only entry.

---

## 6. Verification

All four run as `cd apps/server_core` in the same command with
`GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache`, no `-mod=mod`.

### 1. `go build ./...`

```
=== 1. go build ./... ===
build exit=0
```

### 2. `go vet ./...`

```
=== 2. go vet ./... ===
vet exit=0
```

### 3. `go test -count=1 ./internal/modules/product_links/...`

```
=== 3. go test product_links ===
ok  	marketplace-central/apps/server_core/internal/modules/product_links/adapters/connectors	2.407s
?   	marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/product_links/application	2.306s
ok  	marketplace-central/apps/server_core/internal/modules/product_links/composition	1.990s
ok  	marketplace-central/apps/server_core/internal/modules/product_links/domain	1.959s
?   	marketplace-central/apps/server_core/internal/modules/product_links/ports	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/product_links/transport	2.410s
test exit=0
```

### 4. D-121 policy family

```
$ go test -count=1 -run 'TestCase[0-9]+.*|TestNamedMissingAnchorSites.*|TestSellerSKUAnchor.*|TestExactSKU.*|TestUnresolved.*|TestConcordant.*' -v ./internal/modules/product_links/application/
exit=0
```

Top-level results — **17 PASS, 0 FAIL** (`grep -c "^--- PASS"` = 17, `grep -c "FAIL"` = 0):

```
--- PASS: TestCase7VoltageHardNegativeCapsBaixaReject (0.00s)
--- PASS: TestCase3EANAloneYieldsMediaConfirm (0.00s)
--- PASS: TestConcordantAnchorsAreTheOnlyAutomaticPath (0.00s)
--- PASS: TestCase8NoAnchorResolvedYieldsZeroConfidenceNoCandidateWithReasons (0.00s)
--- PASS: TestCase4TitleOnlyYieldsBaixaReview (0.00s)
--- PASS: TestCase6DokaKitHardNegativeCapsBaixaReject (0.00s)
--- PASS: TestCase1ConcordantSKUAndEANYieldsAltaAccept (0.00s)
--- PASS: TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide (0.00s)
--- PASS: TestCase9ProviderDeclaredUnavailableAnchorsOnConcordantPayload (0.00s)
--- PASS: TestCase5SKUEANConflictYieldsBaixaReviewBothSides (0.00s)
--- PASS: TestConcordantCandidateDoesNotDerefNilProduct (0.00s)
--- PASS: TestCase10DimensionHardNegativeCapsBaixaReject (0.00s)
--- PASS: TestCase2SellerSKUAloneWithoutEANYieldsMediaConfirm (0.00s)
--- PASS: TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason (0.00s)
--- PASS: TestUnresolvedDeclaredAnchorReadsProviderValueWithoutERPProduct (0.00s)
--- PASS: TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference (0.00s)
--- PASS: TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/product_links/application	1.000s
```

The count is 17 including the new test; 16 before it. Nothing in the family regressed.

Raw log retained at (scratchpad, session-local, not in-repo):
`…/scratchpad/family.txt`.

---

## 7. What I could not verify, and why

1. **That the new test observes the production path all the way to the database /
   HTTP surface.** It drives `GenerateLinkCandidates` with stubbed
   `Snapshots`/`Matcher`/`Store`/`IdentityAnchors`. What is proven is that the
   *application-layer* production call site (`generation_service.go:216`) passes a nil
   product and that the resulting reason carries `side=erp`. Whether a real Postgres
   matcher can return the empty result set that reaches that branch is **not** proven
   here — it is inferred from `findProducts` returning no products, which I read but did
   not exercise against a live DB. Chips do not boot a server or DB, so this could not
   be closed in this slice.

2. **That `side=erp` survives for `seller_sku` under every provider declaration.** I
   verified survival for `Supplied: true` (Mercado Livre's real shape) and derived, but
   did **not** add a test for, the `Supplied: false` case — where
   `classifyProviderIdentityAnchor` returns `UNAVAILABLE` and the promotion at
   `generation_service.go:692–699` replaces the seeded reason, dropping `side=erp`
   entirely. Whether that is correct behaviour or a second defect is a **question I am
   flagging, not answering** — deciding it would change production classification and is
   outside this slice. No provider currently in the repo declares `seller_sku` as not
   supplied, so it may be unreachable today; I did not sweep every provider registration
   to confirm that.

3. **That the new comment is true for anchors other than `seller_sku`.** The narrowing
   is stated about `seller_sku` only, and I verified it only for `seller_sku` (via
   `identityAnchorValues` `case "seller_sku"` at `generation_service.go:739–754` plus
   `canonicalProductID` at `:464–469`). For `ean`, a present product with a nil
   `product.EAN` *does* yield `side=erp`, which is why the sibling table cases exist. The
   comment does not claim otherwise, but I want it on record that I did not audit the
   other anchors.

4. **The claim "findProducts drops any candidate without a canonical id"** — which the
   new comment relies on — I took from the existing in-file comment at the head of the
   table (the B-01 note) and corroborated only by reading `canonicalProductID`
   (`generation_service.go:464–469`), which returns `ok == false` for a nil or non-positive
   id. I did **not** re-read `findProducts` line by line to confirm the filter is
   unconditional on every branch. If that filter has a hole, the new comment's "with a
   product present, side=erp cannot arise" would need narrowing again.

5. **Whether the earlier deletion of the table case `"exact EAN ERP seller SKU empty"`
   removed any assertion beyond the `side=erp` one.** I restored coverage for the
   reachable `side=erp` case as instructed, but I never saw the deleted case's source —
   only the description of it in the dispatch. If it also asserted a detail string or a
   state I have not reproduced, that part is still missing and I cannot tell.

6. **Anything outside the one file.** Per the slice, I did not touch or run anything in
   `apps/web/`, migrations, `platform/httpx`, or the governance lane. `go build`/`go vet`
   over `./...` passed, which is the only cross-module signal I have.
