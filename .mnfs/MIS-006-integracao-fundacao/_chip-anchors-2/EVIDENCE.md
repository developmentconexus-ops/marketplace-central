# CHIP-ANCHORS-2 — evidence pack

```yaml
chip: CHIP-ANCHORS-2
branch: chip/anchors-2
base_sha: d51a27b665e91132427a0efa1877e9a5df2f11bf   # floor, not ceiling (hub 8fcc8dcf)
rebased_onto: e98d8193eabacc4bc18a35e45fd35960c64422ec  # hub A2-R1/A2-R2 rulings
tip: ae6c6525a3db3ad4c4b278e587040442ec153d22   # code tip; this pack lands as one further commit
level: QA-0
status: closed-pending-hub
```

BASE-SHA: `d51a27b665e91132427a0efa1877e9a5df2f11bf` — declared by the hub as a FLOOR. This branch
was rebased onto `e98d8193` (the commit carrying rulings A2-R1/A2-R2), which has `d51a27b6` as an
ancestor. `git merge-base HEAD main` = `e98d8193eabacc4bc18a35e45fd35960c64422ec`, i.e. the branch
is exactly `main` plus this chip's commits — the eight listed below, plus the commit that files this
pack (`git rev-list --count e98d8193..HEAD` = 8 at the code tip).

CONTRATO: `.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/validation-contract.md` (present on
this branch; read in place, not via `git show`).

HUB-SESSION: `local_99feb041-a5b3-4161-b6dc-bd38e65b6156`

EXEMPLO-IO: protocol `#003-E` · `GET /erp/imports/{id}/chain` →
`{"protocol","importados","vinculados","enfileirados","queue_read_at"}`. The seeded proof uses
protocol `#610-E` / `#630-E` against a real PostgreSQL (this chip does not touch the live stack).
The five-field freeze is in this chip's `chip.md` §F-04 — decision D-D decides only that chain-read
is backend-owned and names no field (correction raised by the F-04 reviewer and verified here).

## Commits

```
ae6c6525 chip(CHIP-ANCHORS-2): F-04 corrective — earn the vinculados DISTINCT, stop promising a drain
37d6b7cc chip(CHIP-ANCHORS-2): F-04-S2 — GET /erp/imports/{id}/chain, contract and SDK in one commit
1bd6ff55 feat(erp-import): read an import's chain from the live tables
e095685a fix(links): give UNAVAILABLE one meaning again
23f82b67 feat(links): stable dimension ordering with equivalence display
655f5d16 fix(links): preserve anchor reason evidence
9c030154 feat(links): add incomparable anchor reasons
b8ca9550 chip(CHIP-ANCHORS-2): F-01 — refforn leaves the cross-side anchor vocabulary
e98d8193 hub(CHIP-ANCHORS-2): rule A2-R1 and A2-R2 — grant the reclassification, narrow C2   <- base
```

## Dispatch ledger

Rows written at dispatch time (core §1). Codex workers ran as OS processes with the prompt read
from a file, stdin closed, effort explicit, output teed to `agent__<id>.log` and captured verbatim
by `-o agent__<id>.last.md`. All artifacts live in this session's scratchpad
`…\6fc58a7a-34f3-402c-9cd2-ac7530a0f59d\scratchpad\`.

| # | Phase | Role | Model / effort | Path | Prompt artifact | Output artifact | Status |
|---|-------|------|----------------|------|-----------------|-----------------|--------|
| D1 | P3 | F-01 implementation | **none — authored by the chip owner (Claude)** | in-session | — | commit `b8ca9550` | GREEN. Declared, not hidden: F-01 was written before the per-slice codex lane was adopted for this chip. See FINDING F-6. |
| D2 | P3 | F-02 implementation (whole feature, not per slice) | gpt-5.6-luna / high | OS-process codex | `prompt-f02.md` | `agent__f02.last.md` + `.log` | committed `9c030154`. See FINDING F-6: dispatched per FEATURE, contrary to core §4. |
| D3 | P4 | Adversarial reviewer, feature F-02 | Claude sonnet subagent (independent of the implementer) | Agent tool, async | inline brief | `tasks/a4a524054a2636477.output` | completed; findings accepted, corrective dispatched as D4 |
| D4 | P3 | F-02 corrective (reason evidence preserved) | gpt-5.6-luna / high | OS-process codex | `prompt-f02-corr.md` | `agent__f02corr.last.md` + `.log` | committed `655f5d16` |
| D5 | P2 | Slice planner for F-03 + F-04 | gpt-5.6-sol / medium | OS-process codex | `prompt-plan-f03-f04.md` | `agent__planf0304.last.md` + `.log`; cards in `slice-cards-f03-f04.md` | completed — from here on, one dispatch per SLICE |
| D6 | P3 | F-03-S1 (stable ordering + equivalence display) | gpt-5.6-luna / high | OS-process codex | `prompt-f03s1.md` | `agent__f03s1.last.md` + `.log` | code GREEN; **evidence rejected** (FINDING F-5) — corrective D7 |
| D7 | P3 | F-03-S1 corrective (evidence re-run) | gpt-5.6-luna / high | OS-process codex | `prompt-f03s1-corr.md` | `agent__f03s1corr.last.md` + `.log` | committed `23f82b67` after the chip owner re-ran the must-fail itself |
| D8 | P3 | C11-S1 (A2-R1 reclassification) | gpt-5.6-sol / low (complex slice) | OS-process codex | `prompt-c11s1.md` | `agent__c11s1.last.md` + `.log` | committed `e095685a`, chip-verified independently |
| D9 | P4 | Adversarial reviewer, F-03-S1 + C11-S1 (`23f82b67` + `e095685a`) | Claude sonnet subagent | Agent tool, async | inline brief | `tasks/a69561993e72b5065.output` | completed — **no BLOCKING, no MAJOR**; one MINOR (dead `product == nil` disjunct at the three `applySingleAnchorScore` sites), accepted as-is |
| D10 | P3 | F-04-S1 (domain + port + SQL + fixture) | gpt-5.6-sol / low (complex slice) | OS-process codex | `prompt-f04s1.md` | `agent__f04s1.last.md` + `.log` | SQL correct; **fixture did not execute** against a real database (FINDING F-3) |
| D11 | P3 | F-04-S1 corrective (split the multi-statement seeds) | gpt-5.6-sol / low | OS-process codex | `prompt-f04s1-corr.md` | `agent__f04s1corr.last.md` + `.log` | second blind round still failed (`42P18`); fixed in the owner lane, committed `1bd6ff55` (FINDING F-3) |
| D12 | P3 | F-04-S2 (service + handler + OpenAPI + SDK) | gpt-5.6-luna / high | OS-process codex | `prompt-f04s2.md` | `agent__f04s2.last.md` + `.log` | GREEN for Go; SDK vitest **could not run** in its sandbox and it said so; owner ran the DB proof and the SDK suite, fixed two defects, committed `37d6b7cc` |
| D13 | P4 | Adversarial reviewer, feature F-04 (`1bd6ff55` + `37d6b7cc`) | Claude sonnet subagent | Agent tool, async | inline brief | `tasks/aaa3a19adafbfc380.output` | completed — **no BLOCKING, no MAJOR**; 2 MINOR, both accepted and fixed in `ae6c6525`; one citation correction accepted. See §"F-04 adversarial review" |
| D14 | P3 | F-04 corrective (both MINOR findings) | **chip owner (Claude)** — two-line scope, no dispatch | in-session | — | commit `ae6c6525` | GREEN; both fixes proven by must-fail / by grep |

The three reviewer `.output` files are 0 bytes on disk (harness artifact of this session — the
transcripts were streamed, not persisted). Their verdicts are transcribed here from the completion
notifications; the transcription, not the file, is the record.

---

## C1 — `refforn` out of the cross-side vocabulary (F-01)

`git grep -n "IdentityAnchorRefforn"` at the tip:

```
.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/chip.md:59:Remove `IdentityAnchorRefforn` de `knownIdentityAnchors` em
.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/validation-contract.md:8:`IdentityAnchorRefforn` não aparece mais em `knownIdentityAnchors`. Provado por **string**:
.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/validation-contract.md:9:`git grep -n "IdentityAnchorRefforn"` na tip, com a saída colada. As ocorrências que
.mnfs/MIS-006-integracao-fundacao/_chip-anchors/cite-audit.txt:66:pack:1377  `marketplace_capability.go:25-29`                          | IdentityAnchorSellerSKU IdentityAnchor = "seller_sku"  ..  ->29: IdentityAnchorRefforn   IdentityAnchor = "refforn"
```

Four occurrences remain, named one by one: three are pack prose (this chip's `chip.md` and
`validation-contract.md`), one is a citation table inside the CLOSED CHIP-ANCHORS pack. **Zero are
Go code**, so the identifier no longer exists in `knownIdentityAnchors` or anywhere the generator
reads.

`refforn` is intact on the ERP side — the removal is from the anchor list, not from the data:

```
apps/server_core/migrations/0046_create_erp_import_products.sql:10:    refforn TEXT,
```

## C2 — a declared anchor whose VALUE is missing on either side does not vanish silently (F-02)

Narrowed by ruling A2-R2: the case "declared and present on both sides" still emits nothing,
because the classifier compares PRESENCE and never established that the two sides agree. This pack
claims only what the table verifies.

`TestGenerateLinkCandidatesKeepsEveryDeclaredAnchorVisible` covers the four cases and asserts the
anchor **is present** in `reasons[]` in all four; `TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide`
covers the per-site `side` values; `TestProviderDeclaredUnmodelledAnchorIsIncomparableWithoutSide`
covers the provider-declares-but-model-has-no-such-field case.

```
--- PASS: TestProviderDeclaredUnmodelledAnchorIsIncomparableWithoutSide (0.00s)
--- PASS: TestGenerateLinkCandidatesKeepsEveryDeclaredAnchorVisible (0.00s)
--- PASS: TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide (0.00s)
```

The `side=both` emission survives the C11 reclassification — the classifier still returns it:

```
generation_service.go:723: return domain.LinkCandidateReasonDirectionIncomparable, domain.LinkCandidateReasonSideBoth, fmt.Sprintf("anúncio e produto ERP sem %s", anchor.Anchor), true
```

**Gap recorded, by design (FINDING F-1):** an anchor declared and present on BOTH sides emits no
reason at all. Narrowing the criterion was the hub's remedy under R-24; widening the code to emit an
unverified `FOR` would be the defect R-24 exists to prevent.

## C3 — `side` exists only on `INCOMPARABLE`

`TestProductLinkReasonSideJSONIsOnlyPresentForIncomparable` serialises a `FOR`, an `AGAINST`, an
`UNAVAILABLE` and an `INCOMPARABLE` and asserts against the produced JSON, not the struct.

```
--- PASS: TestProductLinkReasonSideJSONIsOnlyPresentForIncomparable (0.00s)
```

## C4 — D-121 auto-approval policy intact

```
--- PASS: TestConcordantAnchorsAreTheOnlyAutomaticPath (0.00s)
--- PASS: TestIncomparableReasonsDoNotChangeAutomaticOrConfirmationPolicy (0.00s)
--- PASS: TestSingleAnchorGoesToConfirmationNeverAutoApproved (0.00s)
```

String proof that no confidence/score/band path reads `INCOMPARABLE` — every non-test occurrence in
`product_links`, listed exhaustively:

```
generation_service.go:636   reason.Direction = domain.LinkCandidateReasonDirectionIncomparable        (classifier assignment)
generation_service.go:639   reason.Direction = domain.LinkCandidateReasonDirectionIncomparable        (classifier assignment)
generation_service.go:659   reason.Direction != domain.LinkCandidateReasonDirectionIncomparable       (precedence guard, not a score)
generation_service.go:690   if direction == domain.LinkCandidateReasonDirectionIncomparable           (precedence guard, not a score)
generation_service.go:711   return ... Incomparable, "", "não foi possível comparar a âncora %s"      (classifier return)
generation_service.go:715   return ... Incomparable, ...SideProvider, "anúncio sem %s"                (classifier return)
generation_service.go:723   return ... Incomparable, ...SideBoth, "anúncio e produto ERP sem %s"      (classifier return)
generation_service.go:726   return ... Incomparable, ...SideProvider, "anúncio sem %s"                (classifier return)
generation_service.go:728   return ... Incomparable, ...SideERP, "produto ERP sem %s cadastrado"      (classifier return)
domain/link_candidate.go:57 LinkCandidateReasonDirectionIncomparable LinkCandidateReasonDirection = "INCOMPARABLE"   (the constant)
```

None is a confidence sum, a band decision or a status transition.

## C5 — stable ordering and equivalence display (F-03)

**(a) Determinism.** `TestHardNegativeDimensionStableEquivalentDisplay` at `-count=10` — full output
in `scratchpad/c5a.txt`, ten `--- PASS` lines:

```
--- PASS: TestHardNegativeDimensionStableEquivalentDisplay (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/product_links/application	1.492s
```

**(b) Load-bearing proven.** `slices.SortStableFunc` reverted to `slices.SortFunc`, same guard at
`-count=10`: **ten `--- FAIL` lines**, exit 1 (`scratchpad/c5b.txt`). The signature of the defect is
the display flipping orientation while the canonical key stays identical:

```
    generation_service_test.go:79: hardNegativeDimension("Produto 50cm 500MM 1m 1000mm …")=("mm:100|mm:1000|…", "100MM ≡ 10cm|1m ≡ 1000mm|…|500MM ≡ 50cm|…", true), want ("mm:100|mm:1000|…", "10cm ≡ 100MM|1m ≡ 1000mm|…|50cm ≡ 500MM|…", true)
--- FAIL: TestHardNegativeDimensionStableEquivalentDisplay (0.00s)
```

(The two long tuples are elided at `…` here; `scratchpad/c5b.txt` holds them in full, unelided.)
`SortStableFunc` restored immediately after; `git status --porcelain` and `git diff --stat` are both
empty at the tip.

**(c) Display.** `50cm` × `500MM` → `50cm ≡ 500MM` is asserted in the same guard (see the `want`
string above). `50cm` × `50cm` → `50cm` — the set-wide dedup that prevents `50cm ≡ 50cm` is
`if !slices.Contains(last.displays, pair.display)` in `generation_service.go`, and the independent
reviewer (D9) confirmed the dedup is set-wide rather than adjacent-only.

## C6 — the chain read answers the three numbers (F-04)

Proven against a REAL PostgreSQL 16 with the repository migrations applied (`applied 69
migration(s)`), provisioned by this chip as an ephemeral container — no dev stack, no `:8080`, no
`.env*` read:

```
=== RUN   TestGetImportChainCountsCurrentQueueAcrossInstallations
--- PASS: TestGetImportChainCountsCurrentQueueAcrossInstallations (2.18s)
=== RUN   TestGetImportChainHTTPIntegration
--- PASS: TestGetImportChainHTTPIntegration (2.44s)
=== RUN   TestGetImportChainMissingAndProtocolWithoutSyncState
--- PASS: TestGetImportChainMissingAndProtocolWithoutSyncState (3.33s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	11.303s
```

The fixture carries every element C6 names, and each is load-bearing:

- `accepted_count=99` on the protocol row while `Importados` asserts `4` — proves the count comes
  from the rows, not from the frozen counter column.
- two installations with different queues (`installation-a`, `installation-b`) and `101` present in
  both — proves `DISTINCT` rather than a summed duplicate. Removing `DISTINCT` from the statement
  yields `Enfileirados:4` and the test FAILS (must-fail run during F-04-S1).
- `103` queued only in installation B — proves aggregation across installations.
- `OUTSIDE` in a queue — proves the queue is intersected with the protocol's codprods.
- `104` queued nowhere — proves `K ≠ N`, so the filter is exercised.
- the second test's other-tenant rows — proves tenant isolation.
- codprod `101` resolved TWICE, one listing per installation (added in `ae6c6525` after the F-04
  review found the `vinculados` `DISTINCT` unexercised) — proves `vinculados` counts a codprod once.
  Dropping that `DISTINCT` yields `Vinculados:3` and the test FAILS; output in §"F-04 adversarial
  review".

Run at the tip `ae6c6525`, after the corrective:

```
--- PASS: TestGetImportChainCountsCurrentQueueAcrossInstallations (1.99s)
--- PASS: TestGetImportChainHTTPIntegration (1.62s)
--- PASS: TestGetImportChainMissingAndProtocolWithoutSyncState (2.25s)
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	9.012s
```

## C7 — `enfileirados` semantics declared in the contract

```
contracts/api/marketplace-central.openapi.yaml:8080:      required: [protocol, importados, vinculados, enfileirados, queue_read_at]
contracts/api/marketplace-central.openapi.yaml:8091:            Current market queue at queue_read_at, never an import-history total:
contracts/api/marketplace-central.openapi.yaml:8092:            the number falls as the queue is consumed. No consumer is wired yet,
contracts/api/marketplace-central.openapi.yaml:8093:            so today it only grows — a later drop is drainage, never data loss.
contracts/api/marketplace-central.openapi.yaml:8094:        queue_read_at:
contracts/api/marketplace-central.openapi.yaml:6562:          description: Emitted only when direction is INCOMPARABLE.
```

`queue_read_at` is in the response and in the schema, typed `format: date-time`, and required. Line
6562 closes the gap F-02 left on `ProductLinkCandidateReason.side`.

The description shipped in `37d6b7cc` said *"the number falls as the scheduler drains it"*. The F-04
reviewer showed nothing drains it — the only registered job is `products`
(`_ = scheduler.RegisterJob(domain.EntityProducts, …)`), and the only writer of the market queue is
the enqueuer, which appends. The sentence was replaced in `ae6c6525` rather than annotated (R-25).
What C7 exists to prevent — a consumer reading the field as an import-history total and reporting
data loss — is still stated, and now stated truthfully.

The same semantics are where a Go reader meets them, in `domain.ImportChain`:
`"Enfileirados is the CURRENT market queue, not a historical total: it falls as the queue is
consumed. No consumer is registered today — the scheduler runs a products job only — so the count
only grows until one is wired, and a later drop is drainage rather than data loss."`

## C8 — negative cases of the endpoint

Both asserted inside `TestGetImportChainHTTPIntegration` (PASS above), against the RAW response
body — a decoded struct cannot tell `0` from `null`, and that distinction IS this criterion:

- missing protocol → `404` (`if missing.Code != http.StatusNotFound`);
- real protocol with no `sync_state` row → `200` containing the literal `"enfileirados":0`.

**Load-bearing proven.** Marking the field `json:"enfileirados,omitempty"` drops the key and the
test fails:

```
--- FAIL: TestGetImportChainHTTPIntegration (1.28s)
    chain_query_repository_integration_test.go:172: empty status=200 body={"protocol":"#631-E","importados":1,"vinculados":0,"queue_read_at":"2026-07-28T00:45:14Z"}
        , want numeric enfileirados zero
FAIL
```

Tag restored; the tip carries `json:"enfileirados"`.

## C9 — OpenAPI + SDK in the same commit as the behaviour

F-02, commit `9c030154`:

```
 .../application/auto_link_policy_test.go           |  44 ++++-
 .../application/generation_service.go              | 168 ++++++++++++++-----
 .../application/generation_service_test.go         | 181 +++++++++++++++++----
 .../modules/product_links/domain/link_candidate.go |  15 +-
 contracts/api/marketplace-central.openapi.yaml     |   5 +-
 packages/sdk-runtime/src/index.ts                  |   4 +-
```

F-04, commit `37d6b7cc`:

```
 apps/server_core/internal/composition/root_test.go |  6 ++
 .../chain_query_repository_integration_test.go     | 98 ++++++++++++++++++++++
 .../adapters/postgres/query_repository.go          |  5 ++
 .../erp_import/application/query_service.go        | 23 ++++-
 .../erp_import/application/query_service_test.go   | 42 ++++++++++
 .../modules/erp_import/transport/http_handler.go   | 29 +++++++
 .../erp_import/transport/http_handler_test.go      | 53 +++++++++++-
 contracts/api/marketplace-central.openapi.yaml     | 51 +++++++++++
 packages/sdk-runtime/src/erpImport.test.ts         | 13 ++-
 packages/sdk-runtime/src/erpImport.ts              |  8 ++
 packages/sdk-runtime/src/index.test.ts             | 23 +++++
 packages/sdk-runtime/src/index.ts                  |  3 +-
```

Go behaviour, `marketplace-central.openapi.yaml` and `packages/sdk-runtime/src/*` in one changeset
each (profile §7). The F-02 corrective `655f5d16` touches only Go — it changed no wire shape, so it
carries no contract file. The F-04 corrective `ae6c6525` touches the OpenAPI `description`, the
matching Go doc comment and one test fixture — no wire shape, no SDK surface, so no SDK file: field
names, types and required-ness are byte-identical to `37d6b7cc`.

## C10 — disjoint write-set, `apps/web` untouched

`git diff --name-only d51a27b6 HEAD` (the base_sha floor, so the hub's own pack commits appear):

```
.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/chip.md
.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/hub-rulings.md
.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/validation-contract.md
apps/server_core/internal/composition/root_test.go
apps/server_core/internal/modules/connectors/ports/marketplace_capability.go
apps/server_core/internal/modules/erp_import/adapters/postgres/chain_query_repository_integration_test.go
apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go
apps/server_core/internal/modules/erp_import/application/query_service.go
apps/server_core/internal/modules/erp_import/application/query_service_test.go
apps/server_core/internal/modules/erp_import/domain/import.go
apps/server_core/internal/modules/erp_import/ports/repository.go
apps/server_core/internal/modules/erp_import/transport/http_handler.go
apps/server_core/internal/modules/erp_import/transport/http_handler_test.go
apps/server_core/internal/modules/product_links/adapters/connectors/identity_anchor_adapter_test.go
apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go
apps/server_core/internal/modules/product_links/application/generation_service.go
apps/server_core/internal/modules/product_links/application/generation_service_test.go
apps/server_core/internal/modules/product_links/domain/link_candidate.go
contracts/api/marketplace-central.openapi.yaml
packages/sdk-runtime/src/erpImport.test.ts
packages/sdk-runtime/src/erpImport.ts
packages/sdk-runtime/src/index.test.ts
packages/sdk-runtime/src/index.ts
```

Taken at the code tip `ae6c6525`; the commit that files this pack adds one more path,
`.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/EVIDENCE.md`, and no other.

Paths under `apps/web/`: **0**. New migrations: **0**.

**Declared, not fixed — the `apps/web` `tsc` red this chip causes (FINDING F-2).** The wave-2 chip
**CHIP-VINC-NEUTRO** owns it. Three errors are this chip's, caused by `INCOMPARABLE` joining
`ProductLinkReasonDirection`:

```
apps/web/src/pages/vinculos/QueueRow.tsx(34,7): error TS2741: Property 'INCOMPARABLE' is missing in type '{ FOR: string; AGAINST: string; UNAVAILABLE: string; }' but required in type 'Record<ProductLinkReasonDirection, string>'.
apps/web/src/pages/vinculos/QueueRow.tsx(75,7): error TS2741: Property 'INCOMPARABLE' is missing in type '{ FOR: string; AGAINST: string; UNAVAILABLE: string; }' but required in type 'Record<ProductLinkReasonDirection, string>'.
apps/web/src/pages/vinculos/VinculoDrawer.tsx(118,17): error TS7053: Element implicitly has an 'any' type because expression of type 'ProductLinkReasonDirection' can't be used to index type '{ readonly FOR: "bg-accent-soft text-accent-ink"; readonly AGAINST: "bg-warn-soft text-warn"; readonly UNAVAILABLE: "bg-surface-2 text-faint"; }'.
```

The same run reports 15 `error TS` lines in total. The other 12 are pre-existing baseline and belong
to nobody in this chip — `anunciosQueries.ts`, `anunciosQueryState.test.ts` ×2,
`AnunciosTable.test.tsx`, `ListingsRefreshControl.test.tsx`, `MutationPreviewModal.tsx`,
`MutationResultSummary.tsx`, `ProdutoPage.partialFailure.test.tsx` ×4, `ProdutoPage.test.tsx`. Full
output in `scratchpad/web-tsc-tip.txt`. This chip's diff touches no file listed above, so declaring
is all it may do.

## C11 — `UNAVAILABLE` has one meaning again (A2-R1)

**(a) Reclassification.** `missingMatchedAnchorReason` classifies by which side is missing and is
called at all four sites named in A2-R1; every site's `detail` string is unchanged (the classifier
overwrites direction and side only).

**(b) Absence proven by string.** `git grep -n "DirectionUnavailable"` in `generation_service.go` —
three occurrences, each named:

```
generation_service.go:642:		reason.Direction = domain.LinkCandidateReasonDirectionUnavailable
generation_service.go:658:		if reason.Direction != domain.LinkCandidateReasonDirectionUnavailable &&
generation_service.go:704:		return domain.LinkCandidateReasonDirectionUnavailable, "", fmt.Sprintf("provider não fornece a âncora %s", anchor.Anchor), true
```

- `:704` is the `!anchor.Supplied` path — the one meaning `UNAVAILABLE` is allowed to keep.
- `:642` is the default branch of the classifier: reached only when BOTH values are present, i.e.
  the branch A2-R1 explicitly excluded (see (c)).
- `:658` is a comparison inside the precedence guard, not an emission.

No `UNAVAILABLE` is emitted outside `!anchor.Supplied` except the excluded branch, which the ruling
requires to stay as it is.

**(c) Excluded branch untouched, and it is a FINDING (F-1b).**
`TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason` pins it: listing `EAN: "EAN-LISTING"`,
ERP product `EAN: "EAN-ERP"` — both non-empty and different — still yields
`direction=UNAVAILABLE`, `side=""`, and the detail string
`"sem EAN para corroborar o CODPROD: o EAN do anúncio não casa nenhum produto"`.

```
--- PASS: TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason (0.00s)
```

Recorded for the hub to take to the operator: this reason says "unavailable" about a value that IS
available on both sides and merely does not match. Per A2-R1 it does NOT become `INCOMPARABLE`
(swapping one wrong statement for another) and does NOT become `AGAINST` (that is D-121, the
operator's decision).

**(d) C4 revalidated AFTER the reclassification** — the three C4 tests above were run at the tip
`37d6b7cc`, i.e. after `e095685a`, not before.

---

## Ladder

**L0**

```
go build ./...    → exit 0, no output
go vet ./...      → exit 0, no output
go vet -tags integration ./internal/modules/erp_import/...  → exit 0
```

Governance lane: `pwsh scripts/harness.ps1 -Command governance -BaseSha e98d8193eabacc4bc18a35e45fd35960c64422ec`
exits 1 with `status=failed` and 53 finding lines. **This chip causes none of them.** Proven by
running the same command in a detached worktree checked out at the base commit `e98d8193` and
diffing the two outputs:

```
diff gov-base.txt gov.txt  →  IDENTICAL
```

Byte-identical at base and at tip. The findings are pre-existing tree state (`GOV_MODULE_COVERAGE`
for `sourcekind`/`tenant_config`, module-dependency edges into `erp_import`, `RCFG_*` on deploy
scripts and provider adapters); not one cites a path in this chip's write-set. Both captures are in
`scratchpad/gov.txt` and `scratchpad/gov-base.txt`. The baseline worktree was removed afterwards
(`git worktree remove --force`); `git worktree list` shows only the pre-existing entries.

Note for the hub: invoked without `-BaseSha`, the lane fails immediately with
`GOV_SEMANTIC_DRIFT / base-sha-invalid` — `scripts/harness.ps1:113` requires a full 40-hex sha and
does not derive one.

**L1**

```
go test ./... -count=1    → no FAIL lines; every package ok or [no test files]
```

Plus the F-03 guard at `-count=10` green (C5a) and failing under the reverted sort (C5b), and the
three real-database integration tests (C6/C8).

SDK, run by the chip owner because the worker's sandbox could not
(`Cannot read directory "../../../../../../..": Access is denied.`):

```
 Test Files  5 passed (5)
      Tests  76 passed (76)
```

`tsc --noEmit` for `packages/sdk-runtime` exits 0.

---

## FINDINGS

**F-1 — C2 gap, by design.** An anchor declared and present on both sides emits no reason. The
classifier compares presence, never agreement. Recorded under A2-R2; not closed by this chip.

**F-1b — the excluded `UNAVAILABLE` branch (C11c).** `generation_service.go:642` reports
`UNAVAILABLE` for a value present on both sides that simply does not match. Pinned by test, left
unchanged per A2-R1, escalated here for the operator.

**F-2 — `apps/web` `tsc` red.** Three errors caused by this chip; owner CHIP-VINC-NEUTRO. Strings
and the 12-error baseline are in C10 above.

**F-3 — a fixture that cannot execute is not evidence.** F-04-S1's seeds packed four `;`-separated
INSERTs into one parameterised `pool.Exec`; pgx sends parameterised queries over the extended
protocol, which accepts one command per statement. Every number those tests claimed to prove was
unproven, and the worker could not see it because its sandbox has no database. The corrective round
split the statements but still passed an unused `importID` as `$1` (`could not determine data type
of parameter $1`, SQLSTATE `42P18`). The owner fixed the parameter lists rather than spend a third
blind round. **Lesson for the harness: a slice whose validation_kind is `integration` cannot be
discharged by a worker with no database — the owner must run the proof, and the worker must be told
that a SKIP is not a pass.** Both workers on this chip were told exactly that and both complied.

**F-4 — the F-04-S2 worker's two defects, found by the owner running what it could not.**
(a) `packages/sdk-runtime/src/index.test.ts` called `createMarketplaceCentralClient("http://localhost", fetcher)`
while the function takes a single options object — the test would have failed at runtime and the
request URL would have begun `undefined/`. `tsc --noEmit` did not catch it. (b) `NewQueryService`
obtained the optional chain capability with `chainRepo, _ := repo.(ports.ImportChainRepository)` and
called it unguarded: a decorator that wraps the repository without re-exposing the port erases the
capability and the handler panics on a nil interface. This is the D-120 catalog-503 class exactly.
Both fixed in the owner lane before the commit: the test now uses the options form, the adapter
carries `var _ ports.ImportChainRepository = (*Repository)(nil)`, and `GetImportChain` returns
`ErrImportChainUnavailable` instead of panicking (`TestQueryServiceGetImportChainFailsHonestlyWithoutTheCapability`).

**F-5 — reconstructed evidence, rejected.** The F-03-S1 worker submitted a must-fail block that was
internally impossible for the code it described. It was rejected, the corrective prompt made the
VERBATIM requirement explicit, and the owner re-ran the must-fail himself (C5b above is the owner's
run, not the worker's). Reconstructed evidence reads as proof while proving nothing.

**F-6 — dispatch shape, declared.** F-01 was authored by the chip owner with no codex dispatch, and
F-02 was dispatched as a whole FEATURE rather than per slice; both predate the operator's correction
mid-chip. From D5 onward every dispatch is one slice (F-03-S1, C11-S1, F-04-S1, F-04-S2), which is
what core §4 requires. Declaring the deviation is the point; the commits stand on their tests.

**F-7 — owner error: `git checkout -- <path>` used to revert a must-fail edit.** While proving the
SDK encoding assertion load-bearing, the owner reverted the temporary edit with
`git checkout -- packages/sdk-runtime/src/index.ts`, which discarded the whole uncommitted slice
change in that file, not just the temporary edit. Detected immediately by grep, restored by
re-authoring the two lines, and verified by `git diff` (exactly the two intended lines) plus a
re-run of the SDK suite (76/76) and `tsc` (exit 0) before the commit. **Lesson: revert a must-fail
edit with the inverse edit; never with a git discard on a file that holds uncommitted work.**

**F-8 — a malformed `{id}` probably answers 500, not 404 (pre-existing, not a regression).** The
handler builds `domain.ImportID(r.PathValue("id"))` with no format check, and the repository binds
that string to a `uuid`-typed parameter, so a non-UUID id is a driver-level error rather than
`pgx.ErrNoRows` — the 404 branch is never reached. Raised by the F-04 reviewer and NOT verified by
either of us: no test in this chip drives a malformed id, and the owner did not add one, because the
same shape already governs `handleGetImport` on the pre-existing `/erp/imports/{id}` route. The new
route inherits the module's behaviour rather than introducing a new one. Recorded for the hub; no
contract criterion covers it.

**F-9 — a DISTINCT nobody can break is not proven, only inspected.** The F-04 reviewer found that
`vinculados` used `SELECT DISTINCT` while every fixture resolved each codprod through exactly one
`product_links` row: deleting the keyword would have changed no assertion. The same pack asserted
`enfileirados`'s DISTINCT with a codprod queued by two installations, so the gap was a lapse in
coverage, not in belief. **Lesson: a de-duplicating clause earns its place only when some fixture
duplicates; otherwise the guard is decoration and the next refactor deletes it silently.** Fixed in
`ae6c6525`; the must-fail is in the review section below.

**F-10 — the contract described a drain no code performs.** The F-04-S2 OpenAPI text said the
`enfileirados` count "falls as the scheduler drains it". Nothing in this repository consumes the
market queue: the only registered scheduler job is `products`, and the only writer of
`sync_state.cursor -> 'pending'` for `entity='market'` appends to it. The sentence was true of the
system we intend and false of the system we have. Under R-25 that is a falsehood, not a gap, so it
was deleted and replaced with what holds today (§C7) — the Go doc comment on `domain.ImportChain`
carries the same correction, so the two cannot drift apart. **Lesson: a description that narrates a
consumer must name the consumer, or it is a promise wearing the clothes of a fact.**

---

## F-04 adversarial review

Independent Claude sonnet subagent (implementer ≠ reviewer, D-51), read-only, over `1bd6ff55` and
`37d6b7cc`. **No BLOCKING, no MAJOR.** What it verified and could NOT refute, in its own terms: the
three counts and their tenant scoping in every CTE; the absence of any `ORDER BY` (so the
`::text`-cast trap does not apply — the only cast lives in a `JOIN ON`); the cross-tenant leak test;
the `enfileirados` DISTINCT counter-example; C8 asserted against the raw body; C9 field-by-field
between the Go struct, the OpenAPI schema and the TS interface; C10 zero `apps/web` and zero
migrations; the optional-capability seam including that `root.go` passes `erpRepo` **unwrapped**, so
the assertion succeeds live; and that no database detail escapes into a response body.

Two MINOR findings, both accepted and both fixed in `ae6c6525`:

1. **The `vinculados` DISTINCT was never exercised.** Every fixture resolved each codprod through
   exactly one `product_links` row, so removing `DISTINCT` would have left `Vinculados` at 2 and both
   tests would still have passed — the correctness claim rested on inspection, not on a guard. Fixed:
   codprod `101` is now resolved twice, one listing per installation. Must-fail proof:

   ```
   --- FAIL: TestGetImportChainCountsCurrentQueueAcrossInstallations (0.93s)
       chain_query_repository_integration_test.go:81: chain=domain.ImportChain{Protocol:"#610-E", Importados:4, Vinculados:3, Enfileirados:3, …}
   ```

2. **The contract promised a drain that nothing performs.** `description: Current queue count; the
   number falls as the scheduler drains it.` describes behaviour absent from this repository: the
   only registered scheduler job is `products` (`_ = scheduler.RegisterJob(domain.EntityProducts, …)`
   in `sync/composition/products_job.go`), and the only writer of the market queue is the import
   enqueuer, which appends. I verified both greps myself before accepting. Today the count only
   grows. Under R-25 a false sentence is deleted, not annotated, so the text now states what is true
   (see §C7).

One citation correction, accepted: the five-field freeze lives in this chip's `chip.md` §F-04, not
in decision D-D — D-D decides only that chain-read is backend-owned and closes an M-06 planning gap,
and names no field. The dispatch prompt for F-04-S2 cited D-D for the freeze; the code matches the
frozen list either way, but the authority pointer in this pack is corrected here and in §EXEMPLO-IO.

Its "could not verify" list, kept as-is because a named gap is worth more than a guess:

- all three `//go:build integration` tests SKIP in the reviewer's sandbox (no database), so it
  reasoned through the SQL by hand. **The chip owner ran them against a real PostgreSQL — §C6.**
- a malformed, non-UUID `{id}` on this route probably produces 500 rather than 404, since
  `domain.ImportID(r.PathValue("id"))` does no format validation before a `uuid`-typed parameter.
  The reviewer notes this is the pre-existing pattern of `handleGetImport`, so it is not a regression
  from this chip. Recorded as FINDING F-8 for the hub.
- it observed an uncommitted, line-ending-only modification to `generation_service.go` mid-review.
  That was the chip owner's C5b must-fail in flight; the tip is clean (`git status --porcelain`
  empty before each commit).

---

## Discharge ledger

| Criterion | Discharged by | Where |
|---|---|---|
| C1 | `git grep` at tip, four occurrences named, ERP field shown intact | §C1 |
| C2 | 3 tests PASS + `side=both` string at `generation_service.go:723`; gap declared as F-1 | §C2 |
| C3 | `TestProductLinkReasonSideJSONIsOnlyPresentForIncomparable` over produced JSON | §C3 |
| C4 | 3 policy tests PASS at the tip + exhaustive non-test `Incomparable` grep | §C4 |
| C5 | `-count=10` green; `SortFunc` must-fail 10× FAIL; equivalence display in the same guard | §C5 |
| C6 | 3 integration tests PASS on a real PostgreSQL at `ae6c6525`; every fixture element mapped to its claim; both DISTINCTs must-fail-proven (F-9) | §C6 |
| C7 | 4 strings in `marketplace-central.openapi.yaml` at `ae6c6525`, the drain sentence deleted per R-25 (F-10) | §C7 |
| C8 | 404 + literal `"enfileirados":0` asserted on the raw body; `omitempty` must-fail | §C8 |
| C9 | `git show --stat` of `9c030154` and `37d6b7cc`; `ae6c6525` carries no wire change, so no SDK file | §C9 |
| C10 | `git diff --name-only`, 0 `apps/web`, 0 migrations; the red `tsc` declared with its owner | §C10 |
| C11 | classifier at four sites, 3-occurrence grep each named, excluded branch pinned by test, C4 re-run after | §C11 |
| L0 | build 0, vet 0, governance byte-identical to base | §Ladder |
| L1 | `go test ./...` no FAIL, guard `-count=10`, SDK 76/76, `tsc` 0 | §Ladder |

Housekeeping: the ephemeral PostgreSQL container `mpc-chip-anchors2-pg` (host port 55432) was
created by this chip for the C6/C8 proof and removed at close. No dev stack was started, no `:8080`
binding, no `.env*` read, no migration authored, no push.

AGREEMENT — P6 discharged
