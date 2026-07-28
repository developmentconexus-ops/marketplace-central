# P6 gate — Claude/Opus side, round 1 (CHIP-ANCHORS-2, post-merge)

Reviewer: Opus via `harness:gate-reviewer` (physically read-only — no Edit/Write/Bash).
Frozen input: `p6-input-r1.patch`, sha256 `0762d05f92aa2c620354f1ea4dcb5c2be99814e60c7b762dcb003da4e0b4dd32`.
Brief: identical to the Sol side (`p6-sol-gate-r1.md`), minus Bash, plus the instruction not to
guess at the other reviewer. The two reviewers ran concurrently and blind to each other.

Transport note: this text is the reviewer's final message as received. `&` and `<`/`>` arrived
HTML-escaped through the task channel and are left exactly as received rather than
"corrected" — verbatim is a claim about form.

---

## VERDICT

**REFUTED** — one BLOCKING finding (B-01). The change is already merged, so this becomes a corrective commit on `main`, not a merge block. D-121 itself is genuinely untouched (C4 holds under my own reading, not the pack's); the blocking defect is a NEW false operator-facing claim minted by the F-02 classifier.

## PART A — C1..C11

| Criterion | Verdict | Evidence (file:line + quoted line) |
|---|---|---|
| C1 — `refforn` out of vocabulary | PASS | `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go:40-45` — list is now `IdentityAnchorSellerSKU, IdentityAnchorEAN, IdentityAnchorTitle, IdentityAnchorMarca`; patch line 26 deletes `-\tIdentityAnchorRefforn   IdentityAnchor = "refforn"`. Repo-wide grep for `refforn` returns zero hits in the anchor vocabulary. Datum intact: `apps/server_core/migrations/0046_create_erp_import_products.sql:10` `    refforn TEXT,`. No live provider declared it, so the fail-closed path at `connectors/application/marketplace_capability_service.go` is not tripped: `connectors/adapters/mercado_livre/capability_adapter.go:90` `IdentityAnchors: []ports.IdentityAnchor{ports.IdentityAnchorSellerSKU, ports.IdentityAnchorEAN, ports.IdentityAnchorTitle},`. See B-03 for the doctrine comment. |
| C2 — declared anchor with missing VALUE does not vanish | PASS | `product_links/application/generation_service_test.go:292-365`, four cases, assertion via `findReason(result.Items[0].Reasons, "title", tc.direction)` at :356 — a missing reason fails. A2-R2 honoured: `generation_service.go:719-721` `if listingValue != "" && productValue != "" {` → `return "", "", "", false` — no `FOR` on the both-present path. |
| C3 — `side` only on `INCOMPARABLE` | PASS | `generation_service_test.go:536-543` marshals FOR/AGAINST/UNAVAILABLE and asserts `strings.Contains(string(payload), `"side"`)` is false; :551 asserts `"side":"erp"` present for INCOMPARABLE. Asserted on produced JSON, per contract. Tag: `domain/link_candidate.go:72` `Side      LinkCandidateReasonSide      \`json:"side,omitempty"\``. |
| C4 — D-121 intact | PASS | Sole auto-approval gate: `generation_service.go:239` `if candidate.MatchStatus != domain.LinkCandidateMatchStatusAccept {` → `continue`. ACCEPT is set only in `buildConcordantCandidate` **before** reasons exist. My own non-test grep of `LinkCandidateReasonDirectionIncomparable` returns only assignments (:636,:639), precedence guards (:659,:690), classifier returns (:711,:715,:723,:726,:728) and the constant — no confidence sum, band, or status transition. Test `auto_link_policy_test.go:124-162` proves ACCEPT survives an INCOMPARABLE title and single-anchor stays CONFIRM. Persistence cannot break either: `migrations/0065_product_link_candidate_confidence.sql:5` `reasons jsonb` has no direction CHECK. |
| C5(a) — `-count=10` determinism | NOT-PROVEN (by me) | I have no Bash; I cannot re-run. The pasted run is `EVIDENCE.md:160-167`. Structure is sound: `generation_service_test.go:38-81` drives 8 permutations of a 32-token corpus. Artifact that would settle it: the `-count=10` run re-executed at tip `ae6c6525`. |
| C5(b) — must-fail with `SortFunc` | NOT-PROVEN (by me) | Same reason. `EVIDENCE.md:169-180` pastes a FAIL with a display-orientation flip, which is the correct failure signature (a 32-element corpus clears Go's insertion-sort threshold, so unstable sort can reorder equal keys). Not independently verified. |
| C5(c) — equivalence display | PASS | `generation_service_test.go:28-30`: `{title: "Produto 50cm 500MM", wantDisplay: "50cm ≡ 500MM"}` and `{title: "Produto 50cm 50cm", wantDisplay: "50cm"}`. Dedup is set-wide: `generation_service.go:871` `if !slices.Contains(last.displays, pair.display) {`. |
| C6 — chain-read three numbers | PASS | `erp_import/adapters/postgres/chain_query_repository_integration_test.go:65-80`: two installations with different queues, `101` in both, `OUTSIDE` in one queue, `104` queued nowhere, `101` resolved twice, and `if got.Protocol != "#610-E" \|\| got.Importados != 4 \|\| got.Vinculados != 2 \|\| got.Enfileirados != 3`. Both DISTINCTs are exercised (drop either and the fixture's 2→3 / 3→4). Cross-tenant isolation proven separately at :180-237. |
| C7 — `enfileirados` semantics in contract | PASS | `contracts/api/marketplace-central.openapi.yaml:8090-8093` `Current market queue at queue_read_at, never an import-history total: / the number falls as the queue is consumed. No consumer is wired yet, / so today it only grows — a later drop is drainage, never data loss.` and `:8094 queue_read_at:` / `:8080 required: [protocol, importados, vinculados, enfileirados, queue_read_at]`. The "no consumer wired" clause is true: `EntityMarket` appears only in `sync/domain/sync_state.go:21`, `erp_import/adapters/sync/enqueuer.go:44` and tests — nothing consumes it. |
| C8 — negative cases | PASS | `chain_query_repository_integration_test.go:167-177`: `if missing.Code != http.StatusNotFound` and `!strings.Contains(empty.Body.String(), `"enfileirados":0`)` — raw body, so `0` vs `null` is actually distinguished. Handler: `erp_import/transport/http_handler.go:126-128` maps `ports.ErrImportNotFound` → 404; SQL returns `pgx.ErrNoRows` when `import_target` is empty (`query_repository.go:119-121`). Field carries no `omitempty` (`http_handler.go:215`). |
| C9 — OpenAPI + SDK in the same commit | NOT-PROVEN (by me) | Requires `git show --stat`, which I cannot run. The contract-required artifact is pasted at `EVIDENCE.md:277-303` and is internally consistent with the frozen diff (F-02 stat lists `generation_service.go` + `marketplace-central.openapi.yaml` + `packages/sdk-runtime/src/index.ts`; F-04 stat lists the Go files + yaml + 3 SDK files). Artifact that would settle it: `git show --stat 9c030154 37d6b7cc` re-run by the hub. |
| C10 — disjoint write-set, no `apps/web`, no migrations | PASS | Frozen diff `code-diff.patch`, 20 `diff --git` headers: 16 under `apps/server_core/internal/`, 1 `contracts/api/marketplace-central.openapi.yaml`, 4 under `packages/sdk-runtime/src/`. Zero paths under `apps/web/`, zero under `apps/server_core/migrations/` (which the diff's `apps` scope would have shown). |
| C11(a) — reclassification at the A2-R1 sites | PASS on placement, DEFECTIVE on content | Sites reclassified via `missingMatchedAnchorReason` at `generation_service.go:536, 546, 552, 553, 625, 626`; every `detail` string is preserved (patch lines 838/848/854-857 swap only the constructor). **But the `seller_sku` side is computed from the wrong ERP field — see B-01.** |
| C11(b) — no `UNAVAILABLE` outside `!anchor.Supplied` except the named branch | PASS with imprecision | Emissions: `generation_service.go:704` (`!anchor.Supplied`) and `:642` (default of `missingMatchedAnchorReason`). `:658` is a comparison. The pack names both (`EVIDENCE.md:373-381`). Imprecision recorded as B-06. |
| C11(c) — excluded branch pinned by test | PASS | `generation_service_test.go:381-387`: listing `EAN-LISTING` vs ERP `EAN-ERP`, asserts `LinkCandidateReasonDirectionUnavailable` **and** `reason.Side != "" \|\| reason.Detail != "sem EAN para corroborar o CODPROD: o EAN do anúncio não casa nenhum produto"` fails. It runs through the full service with `mercadoLivreIdentityAnchorReader()` (`:189-198`, `ean` Supplied), so the promotion path at `:687-694` is exercised, not just the seed. |
| C11(d) — C4 re-run AFTER reclassification | NOT-PROVEN | The contract says "Rode C4 de novo **no tip final**". `EVIDENCE.md:401-402` anchors the C4 run at `37d6b7cc`, not at the tip `ae6c6525`. The pack anchors C6 at the tip explicitly (`EVIDENCE.md:220` "Run at the tip `ae6c6525`, after the corrective"), so the omission is visible by contrast. `ae6c6525` touched a test fixture and doc/description text, which is very likely harmless — but "very likely" is not the artifact the criterion asks for. Artifact that would settle it: the three C4 tests re-run at `ae6c6525`. |

Marker check: `EXEMPLO-IO` present (`EVIDENCE.md:24`) and asserted by `TestGetImportChainHTTPIntegration`; `AGREEMENT — P6 discharged` present (`EVIDENCE.md:607`); `P6-DUAL-GATE:` correctly **absent** from the chip pack. No `LIVE-VERIFIED:` / `LIVE-WAIVED-BY-OPERATOR:` anywhere — the contract assigns live-drive to the hub (`validation-contract.md:147`), so the new HTTP surface has never been driven against the real stack. That is a stated hub obligation, not a chip gap, but it is still an open obligation.

## PART B — FINDINGS

### B-01. The `seller_sku` anchor's ERP-side value is read from the ERP *supplier reference* (`refforn`), not from CODPROD — the new `side` claim is false — BLOCKING

- **Locus**: `apps/server_core/internal/modules/product_links/application/generation_service.go:734-738` (new in this diff, patch lines 1058-1080):
  ```go
  case "seller_sku":
  	listingValue = listing.SellerSKU
  	if product != nil && product.ReferenceCode != nil {
  		productValue = *product.ReferenceCode
  	}
  ```
  Provenance of that field: `apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go:418` `... ReferenceCode: copyTrimmed(row.Referencia), ...` and `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_repository.go:92` `row.Codprod, row.Descrprod, row.Refforn, ...` — `Referencia` **is** `refforn`. The Oracle reader agrees: `internal_read/adapters/oracle/reader.go:92` `ReferenceCode:     nullableString(referenceValue),`.
- **Why it is wrong**: everywhere else in the system the ERP counterpart of a provider `seller_sku` is the **CODPROD**. `erp_import/adapters/internalread/reader.go:456` `if sku := matchableSellerSKU(input); sku != nil && strings.TrimSpace(row.CodigoProduto) == *sku {` — matching compares seller_sku to `CodigoProduto`, never to `Referencia`; the guard's own comment at :468-469 says "seller_sku carries the ERP CODPROD". So the new classifier decides "does the ERP side have this anchor's value?" by inspecting a field that has nothing to do with the anchor. Two concrete wrong outputs, both new (before this diff these reasons were a plain `UNAVAILABLE` with no side claim — patch line 847 `-{Anchor: "seller_sku", Direction: ...Unavailable, Detail: skuDetail},`):
  1. ERP product with an empty `refforn` → `generation_service.go:638-640` yields `INCOMPARABLE` + `side = erp`, attached to the detail `"sem CODPROD para corroborar o EAN"` (`:540`). The emitted reason therefore states *the ERP product has no CODPROD*. Every ERP product has a CODPROD — it is part of the primary key (`migrations/0046_create_erp_import_products.sql:13`) and `canonicalProductID` is a hard filter at `generation_service.go:279`. The reason is not merely imprecise, it is impossible.
  2. ERP product **with** a `refforn` and a listing with a `seller_sku` → `classifyProviderIdentityAnchor:719-721` returns `emit=false` and the anchor disappears from `reasons[]` entirely — the "present on both sides" branch A2-R2 sanctioned. But "present on both sides" was decided by comparing a marketplace SKU against a supplier reference, i.e. two unrelated fields. The silent-absence defect F-02 exists to kill is reintroduced through the wrong comparand.
  Secondary consequence: F-01 removed `refforn` from the cross-side vocabulary on the stated ground that no marketplace can answer it — and this function reinstates `refforn` as the ERP side of a cross-side comparison in the same changeset.
- **Reachability**: ordinary production data. Any Mercado Livre listing that resolves by EAN alone against an ERP product whose `refforn` is empty reaches case 1 via `generation_service.go:546`; the xlsx snapshot carries `refforn` as nullable (`migrations/0046...:10 refforn TEXT,`) and `validation.go:54 row.Refforn = normalizeOptionalString(row.Refforn)` accepts absence. The `TitleMatch` seed (`:552`) and `applyUnresolvedScore` (`:625`) have the same mapping. The guarding test cannot catch it because its fixture is production-impossible: `generation_service_test.go:417` `product:  &internalreaddomain.ProductCandidate{},` — a candidate with no `InternalProductID` at all, which `generation_service.go:279` filters out before any candidate is built. The assertion `side == erp` therefore encodes the defect rather than guarding against it.
- **Yes-if**: `identityAnchorValues`'s `seller_sku` case derives the ERP-side value from the product's canonical CODPROD (`product.InternalProductID` / `product.ProductID`) instead of `product.ReferenceCode`; and the pinning test uses a candidate that carries a canonical id, so the assertion is reachable. If instead the hub rules that `seller_sku`'s ERP counterpart really is `refforn`, then the details at `:531`, `:540`-`:542` that say "CODPROD" are the false half and must change — but the pair as shipped cannot both be true.

### B-02. `/vinculos` silently drops every `INCOMPARABLE` motivo from the collapsed row — and `tsc` does not flag the line that does it — NON-BLOCKING for this chip's write-set, BLOCKING for the wave-2 owner

- **Locus**: `apps/web/src/pages/vinculos/QueueRow.tsx:159`
  ```tsx
  const shown = [...byDirection("AGAINST"), ...byDirection("FOR"), ...byDirection("UNAVAILABLE")].slice(
  ```
- **Why it is wrong**: I am excluding the declared `tsc` breakage (the two `Record<ProductLinkReasonDirection, …>` literals at `:34` and `:75`) per the brief. This line is different: it is type-correct, so the compiler says nothing, and it enumerates directions by literal. A candidate whose reasons are all `INCOMPARABLE` — now the common shape for an unresolved listing (`generation_service.go:624-628`) — produces `shown = []`, i.e. no motivo chip at all. That is exactly the invariant the same file claims to hold, `QueueRow.tsx:154-156`: "Ranking (never filtering) is what keeps at least one motivo on screen even for a row whose only signals are UNAVAILABLE ones (ADR-17 — motivo sempre visível)". The glyph/class maps at `:93` and `:104` additionally resolve to `undefined`, rendering `"undefined SKU"` and a `className` containing the literal `undefined`.
- **Reachability**: any `/vinculos` render against candidates generated after this merge. `tsc` being red does not stop Vite serving the page.
- **Yes-if**: CHIP-VINC-NEUTRO's validation contract names `QueueRow.tsx:159` explicitly (not only the two `Record` maps that `tsc` reports) and pins the "at least one motivo visible for an all-INCOMPARABLE row" case with a test. Fixing only the compiler errors leaves this silent.

### B-03. The new doctrine comment asserts a universal that shipped code in this repo refutes — NON-BLOCKING

- **Locus**: `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go:33-34`
  ```go
  // refforn?" answers `no` for every provider present and future, and keeping it
  ```
- **Why it is wrong**: `apps/server_core/internal/modules/market/domain/identity_resolver.go:81` `candidateModel := normalizedOptional(candidateAttribute(candidate.Attrs, "MODEL"), false)` and `:90-92` `if agrees(localModel, candidateModel) { anchors = append(anchors, "refforn") }` — a live, wired resolver (`composition/market_adapters.go:149` `RefForn: candidate.ReferenceCode,`) compares the ERP `RefForn` against a marketplace-supplied `MODEL` attribute, and that anchor can carry a match to ACCEPT at `:101-102`. So a provider *does* answer a refforn-shaped question today. The removal from `knownIdentityAnchors` may still be right; the justification as written is false, and the pack's own R-24/R-25 doctrine says a claim an artifact cannot make totally is narrowed or deleted, not shipped.
- **Reachability**: no runtime effect — it is a comment. The reachability that matters is a future reader trusting it.
- **Yes-if**: the sentence is narrowed to the vocabulary this file governs (e.g. "no provider declares it through `MarketplaceCapabilityService`"), or the `market` module's `refforn` anchor is named as the known exception.

### B-04. `links.internal_product_id::text = products.codprod` undercounts `vinculados` silently for any non-canonical codprod spelling — NON-BLOCKING

- **Locus**: `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go:89`
  ```sql
  ON links.internal_product_id::text = products.codprod
  ```
- **Why it is wrong**: the two sides are produced by different pipelines. `erp_import_products.codprod` stores the raw string (`migrations/0046_create_erp_import_products.sql:13 PRIMARY KEY (tenant_id, protocol_id, codprod)`), while `product_links.internal_product_id` is `strconv.ParseInt(strings.TrimSpace(row.CodigoProduto), 10, 0)` (`erp_import/adapters/internalread/reader.go:162`). A codprod of `"0012345"` is accepted as valid — `internal_read/domain/seller_sku.go:25-31` only requires digits with at least one non-zero — and parses to `12345`, whose `::text` is `"12345"`. The join misses, `vinculados` silently reports a smaller number than the truth, and nothing anywhere fails. Note the queue side is immune: the enqueuer pushes the raw `row.Codprod` string (`erp_import/application/import_service.go:157`), so `queued_products` is a string↔string comparison. Secondarily, the cast makes the predicate unindexable — the only relevant index is `migrations/0025_product_link_workflows.sql:17-18` `product_links_installation_idx (tenant_id, installation_id, updated_at DESC)`, which cannot serve it; irrelevant at today's cardinality, and I cannot measure it from here.
- **Reachability**: not shown reachable — I have no access to real ERP data and cannot demonstrate a leading-zero or padded CODPROD in `erp_import_products`. Downgraded to NON-BLOCKING on that basis.
- **Yes-if**: either the join is symmetrised (compare canonical integers on both sides, or normalise `products.codprod` the same way the reader does), or an integration case with a leading-zero codprod proves the count is still right.

### B-05. Malformed `{id}` answers 500, not 400/404 — NON-BLOCKING

- **Locus**: `apps/server_core/internal/modules/erp_import/transport/http_handler.go:124` `chain, err := h.queries.GetImportChain(r.Context(), domain.ImportID(r.PathValue("id")))` with `:130` `writeError(w, http.StatusInternalServerError, "internal_error", "")`.
- **Why it is wrong**: `migrations/0045_create_erp_import_protocols.sql:2` `id UUID PRIMARY KEY`. `GET /erp/imports/not-a-uuid/chain` makes PostgreSQL raise `22P02 invalid input syntax for type uuid`, which is not `pgx.ErrNoRows`, so it falls to the 500 branch. A caller cannot distinguish "you sent garbage" from "we broke".
- **Reachability**: trivially reachable by any client; confirmed by construction from the column type, not executed. The class is pre-existing (`GET /erp/imports/{id}` behaves identically) and OpenAPI does declare a 500 for the path (`marketplace-central.openapi.yaml:3283-3288`), so this is inherited, not introduced.
- **Yes-if**: the id is parsed/validated before the query and a malformed value returns 404 (or 400), for both `{id}` routes.

### B-06. The pack describes the surviving `UNAVAILABLE` more narrowly than the code emits it — NON-BLOCKING

- **Locus**: `EVIDENCE.md:379-380` "`:642` is the default branch of the classifier: reached only when BOTH values are present, i.e. the branch A2-R1 explicitly excluded". Code: `generation_service.go:641-642` `default:` / `reason.Direction = domain.LinkCandidateReasonDirectionUnavailable`.
- **Why it is wrong**: A2-R1's carve-out is "valor **não-vazio e diferente**" (`hub-rulings.md:52`). The default branch fires on *both present*, which includes both present and **equal**. The set the code leaves as `UNAVAILABLE` is strictly larger than the set the ruling excused, and the "i.e." asserts they are the same set.
- **Reachability**: not shown reachable — reaching an equal-and-present pair at a seed site requires the anchor's matcher to have returned nothing for a value the ERP row actually carries. I traced the plausible route (invalid-GTIN EAN) and it is closed: `erp_import/adapters/internalread/reader.go:504-505` nulls out `ProductCandidate.EAN` when the mirror EAN is not a valid GTIN, so `productValue` is empty and the branch is `INCOMPARABLE/erp`, not the default. Remaining route is the `QualityMissingProduct` filter at `generation_service.go:279`, which I cannot show is reachable with equal anchor values.
- **Yes-if**: the pack's sentence says "reached when both values are present — which includes, but is wider than, the branch A2-R1 excluded", or the code narrows the default to the differing case.

### B-07. `jsonb_array_elements_text` over a non-array `pending` throws at query time — NON-BLOCKING

- **Locus**: `erp_import/adapters/postgres/query_repository.go:96-98`
  ```sql
  CROSS JOIN LATERAL jsonb_array_elements_text(
  	COALESCE(state.cursor -> 'pending', '[]'::jsonb)
  )
  ```
- **Why it is wrong**: `COALESCE` defends against `NULL` only. If any `sync_state` row for `entity='market'` ever holds a `pending` key that is an object, string or number, the function errors and the whole endpoint 500s for that tenant. `sync/adapters/postgres/sync_state_repo.go` exposes `RecordSuccess(cursor json.RawMessage)`, which writes a caller-supplied cursor verbatim for any entity.
- **Reachability**: not shown reachable. The only writer of an `entity='market'` cursor today is `AppendPendingCodigos` (`erp_import/adapters/sync/enqueuer.go:44`), which always writes an array, and no scheduler job is registered for `EntityMarket` (grep: only the domain constant, the enqueuer and tests).
- **Yes-if**: the CTE guards on `jsonb_typeof(state.cursor -> 'pending') = 'array'`, or `sync_state` constrains the shape.

### B-08. `RegisterRouteClass` is a no-op for every route because the lookup key includes the method — NON-BLOCKING, pre-existing

- **Locus**: `apps/server_core/internal/platform/httpx/route_deadline.go:46` `m.classes[pattern] = class` vs `:50` `m.mux.Handle(pattern, deadlineMiddleware(m.classes[pattern], handler))`. Registration stores `"/erp/imports/{id}/chain"`; `Handle` is called with `"GET /erp/imports/{id}/chain"` (`erp_import/transport/http_handler.go:55`), so the map never hits and `routeClassDeadline("")` returns `interactiveDeadline`.
- **Why it is wrong**: the class is silently ignored. For the new chain route the default happens to equal the intent (15s), so there is no behavioural harm here — but the same bug means `POST /erp/imports` never receives its declared 120s batch deadline (`http_handler.go:39`).
- **Reachability**: reachable today for the batch upload route; not this chip's write-set (`registerInteractiveRoute` is pre-existing; the chip added one call to it).
- **Yes-if**: the class map is keyed by the same string used for registration, with a test that a batch-classed route actually gets 120s.

### B-09. Latent nil-deref in `buildConcordantCandidate` — NON-BLOCKING

- **Locus**: `generation_service.go:490` `product := *comparison.product`
- **Why it is wrong**: `comparison.product` is a pointer that the sibling `applySingleAnchorScore` explicitly nil-checks (`:827-830`), and `classifyProviderIdentityAnchor:713` branches on `comparison.product == nil`. This function dereferences unconditionally.
- **Reachability**: not shown reachable — the only caller is `generation_service.go:302-303`, which always passes `&product`. It is a one-line invariant with no compiler or test protection.
- **Yes-if**: the field is non-nullable by construction (a value, or a constructor that rejects nil for the concordant path).

## WHAT I READ

Brief and binding context:
- `scratchpad\gate-a2\prompt-sol-gate.md`
- `scratchpad\gate-a2\code-diff.patch` (all 2036 lines)
- `.mnfs\MIS-006-integracao-fundacao\_chip-anchors-2\validation-contract.md`
- `...\_chip-anchors-2\hub-rulings.md`
- `...\_chip-anchors-2\chip.md`
- `...\_chip-anchors-2\EVIDENCE.md`

Go — `product_links`:
- `apps\server_core\internal\modules\product_links\application\generation_service.go`
- `apps\server_core\internal\modules\product_links\application\generation_service_test.go`
- `apps\server_core\internal\modules\product_links\application\auto_link_policy_test.go`
- `apps\server_core\internal\modules\product_links\domain\link_candidate.go`
- `apps\server_core\internal\modules\product_links\adapters\connectors\identity_anchor_adapter.go`

Go — `erp_import` / `sync` / `internal_read` / `market` / `connectors` / platform / composition:
- `apps\server_core\internal\modules\erp_import\adapters\postgres\query_repository.go`
- `apps\server_core\internal\modules\erp_import\adapters\postgres\chain_query_repository_integration_test.go`
- `apps\server_core\internal\modules\erp_import\adapters\postgres\mirror_repository.go` (grep)
- `apps\server_core\internal\modules\erp_import\adapters\internalread\reader.go`
- `apps\server_core\internal\modules\erp_import\adapters\sync\enqueuer.go` (grep)
- `apps\server_core\internal\modules\erp_import\application\query_service.go`
- `apps\server_core\internal\modules\erp_import\domain\validation.go`
- `apps\server_core\internal\modules\erp_import\ports\repository.go`
- `apps\server_core\internal\modules\erp_import\transport\http_handler.go`
- `apps\server_core\internal\modules\sync\domain\sync_state.go` (grep)
- `apps\server_core\internal\modules\internal_read\domain\seller_sku.go`
- `apps\server_core\internal\modules\internal_read\domain\internal_product.go`
- `apps\server_core\internal\modules\internal_read\adapters\oracle\reader.go`
- `apps\server_core\internal\modules\catalog\domain\canonical_product.go`
- `apps\server_core\internal\modules\market\domain\identity_resolver.go`
- `apps\server_core\internal\modules\connectors\ports\marketplace_capability.go`
- `apps\server_core\internal\modules\connectors\application\marketplace_capability_service.go`
- `apps\server_core\internal\modules\connectors\adapters\mercado_livre\capability_adapter.go`
- `apps\server_core\internal\platform\httpx\route_deadline.go`
- `apps\server_core\internal\composition\root.go`
- `apps\server_core\internal\composition\market_adapters.go`

Migrations (read for column types/keys only): `0025_product_link_workflows.sql`, `0045_create_erp_import_protocols.sql`, `0046_create_erp_import_products.sql`, `0065_product_link_candidate_confidence.sql`.

Contract and SDK: `contracts\api\marketplace-central.openapi.yaml` (paths ~3230-3290, schemas ~6550-6565 and ~8078-8097); `packages\sdk-runtime\src\index.ts`, `packages\sdk-runtime\src\erpImport.ts`, `...\erpImport.test.ts`, `...\index.test.ts`.

Frontend, read only to size the widening: `apps\web\src\pages\vinculos\QueueRow.tsx`.

Declared incomplete:
- I did **not** run any test, build, `go vet`, `tsc`, or database query — no Bash by instruction. Every execution-shaped criterion (C5a, C5b, and the run-anchoring half of C11d) rests on the pack's pasted output, which I treated as a claim and graded NOT-PROVEN.
- I did **not** run git, so C9's same-commit claim and C10's `--name-only` output are judged from the frozen patch plus the pack's pasted stats, not from commit objects.
- I did **not** open any `.env*` file and printed no environment variable value.
- I could not inspect live data, so B-04 (leading-zero codprod) and B-06 (equal-and-present) are stated as "not shown reachable" rather than dismissed.
