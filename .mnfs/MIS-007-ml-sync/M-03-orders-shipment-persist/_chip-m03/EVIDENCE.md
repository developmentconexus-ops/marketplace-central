# EVIDENCE — M-03-orders-shipment-persist

```yaml
milestone: M-03-orders-shipment-persist
parent: MIS-007-ml-sync
branch: claude/sleepy-perlman-d0d325
base_sha: 21ca3595  # M-03-specific base (git rev-parse cf2c09e1~1); contract yaml's dd89d4b3 is the
                     # MISSION-level design-ratification commit, an ANCESTOR of 21ca3595 (confirmed via
                     # git merge-base --is-ancestor) -- two different reference points, not a conflict.
tip_sha: d22d3d2076e93182c69a60cd429cfa1c78483105
commits:
  - cf2c09e1  feat(MIS-007/M-03/F-01): ML ingest readers for order+shipment detail
  - 65e31c03  feat(MIS-007/M-03/F-02): IngestOrder single writer, persists orders+shipment
  - fe2131e7  feat(MIS-007/M-03/F-03): read-path switch to Postgres-backed shipment+buyer-fiscal readers
  - d22d3d20  chore(MIS-007/M-03): remove ML connector shipment-reader dead code orphaned by F-03
review_model: LEAN (1 adversarial reviewer per feature at feature-end; 1 cold reviewer on full milestone diff)
generated_by: milestone orchestrator (independent verification of every claim below, not self-report)
status: P6 code-level gate PASS. P7 browser QA BLOCKED on hub REQUEST (dev stack repoint) -- see HUB-EVENT-REQUEST-repoint-dev-stack.md.
```

## P6-DUAL-GATE verdicts (per feature)

Per this milestone's LEAN review model (ratified at dispatch, carried in every feature's own
`spec.md`), each feature received ONE adversarial reviewer at feature-end rather than a dual
cold-gate, plus this ONE cold reviewer on the combined diff at milestone-end (immediately below).

| Feature | Adversarial reviewer verdict | Orchestrator independent re-verification |
|---|---|---|
| F-01 (ML ingest readers) | PASS (prior session) | build/vet/test re-run clean at F-02 dispatch time |
| F-02 (IngestOrder v1) | PASS (prior session) | build/vet/test re-run clean at F-03 dispatch time |
| F-03 (read-path switch) | PASS — 1 real finding (`GetOrderBucketCounts` missing default-case guard on persisted `bucket`), fixed + test added + independently re-verified by orchestrator (build/vet/full suite green) | Own diff read in full: root.go, archguard_test.go, enrich_service.go, order_repo.go, both new readers |
| **Milestone (combined diff)** | **PASS — 1 real finding (dead code: `CapabilityAdapter.GetShipmentInfo`/`GetFreeShippingCost` orphaned by F-03's own deletion), fixed (commit `d22d3d20`), independently re-verified** | Full whole-module test suite green post-fix; git status matches reported file list exactly |

Cold reviewer's full verdict is preserved verbatim below (Areas 1-6 map to the milestone contract's
cross-feature risk surface, not to individual criteria one-to-one).

## Criteria

### M03-C1 — GET /orders/{id} sem NENHUM GET vivo ML
**Status: PASS**

- Command: `go test ./internal/modules/orders/adapters/postgres/... -run TestShipmentAndBuyerFiscalReadersNeverImportMercadoLivreConnector -v` + independent grep sweep of `ordersEnrichShipmentReader`/`ordersEnrichBuyerFiscalReader` consumers in `root.go`.
- Expected: contador ML = 0 no read; enrich lê `order_shipments` + colunas fiscais do banco.
- Actual: `zero_live_ml_call_test.go` (read in full, firsthand, by the orchestrator) proves this
  structurally, not just empirically: `TestShipmentAndBuyerFiscalReadersNeverImportMercadoLivreConnector`
  AST-parses `shipment_reader.go`/`buyer_fiscal_reader.go`'s imports and fails if either ever imports
  the `mercado_livre` connector adapter package -- Go has no way to call a package's exported
  identifiers without importing it, so an absent import is a structural, compiler-enforced guarantee,
  not a runtime observation that could regress silently. `TestScanShipmentRowAndBuildBuyerFiscalInfoNeverProduceAnHTTPCall`
  and `TestShipmentReaderGetShipment_NilPoolDegradesWithoutPanicOrNetworkCall` are the runtime
  companions (pure functions, nil-pool degrade with no dial-out). All 3 tests re-run by the
  orchestrator: PASS. Cold review Area 2 independently confirmed (repo-wide grep, not just root.go)
  that `GET /orders/{id}`'s enrich path (`root.go:617`, `NewEnrichServiceWithReaders`) wires only
  `ordersEnrichShipmentReader`/`ordersEnrichBuyerFiscalReader` (both Postgres-backed), never the
  live-ML-backed `ordersBuyerFiscalReader` (which feeds only `ordersIngestSvc`, the batch path).
- Artifact: `apps/server_core/internal/modules/orders/adapters/postgres/zero_live_ml_call_test.go`;
  `internal/composition/root.go:604-621`.
- Blocking failure observed: No.

### M03-C2 — Persistência completa em uma passada
**Status: PASS**

- Command: F-02's own hermetic integration evidence (`IngestOrder` writes order+shipment+fiscal in
  one call) + F-03's `GetOrderBucketCounts` before/after fixture table.
- Expected: row em `order_shipments` p/ pedido com shipping_id; colunas fiscais tipadas preenchidas;
  `bucket` persistido == `DeriveOrderBucket` com shipment status REAL; pedido SEM fiscal (404 honesto)
  → colunas NULL, ingest NÃO falha.
- Actual: `order_shipments` PK is `(tenant_id, provider, provider_shipment_id)` (migration 0088);
  `IngestOrder` (F-02) writes it in the same transaction/call as the order row and the 9 buyer_*
  columns (migration 0089), computing `bucket` via `domain.DeriveOrderBucket` with the REAL shipment
  status available at ingest time (not the pre-M-03 hardcoded `shipmentStatus=""`). Honest-absence:
  `BuyerFiscalReader.GetBuyerFiscal`/`buildBuyerFiscalInfo` (read firsthand by the orchestrator, see
  `buyer_fiscal_reader.go:64-90`) treats a fiscal-data-absent row as `HasData()==false`, never an
  error -- an address is included ONLY when ≥1 of its 6 source fields is present, avoiding a
  fabricated non-nil-but-empty address. Cold review Area 5 independently traced the write side
  (`connectors/adapters/mercado_livre` trimmed-string helpers) and confirmed it never writes an
  empty-string sentinel for a genuinely-unknown field -- always `nil`/SQL `NULL` -- so the read-side
  `pgtype.Text.Valid` check cannot be fooled into treating absence as data. F-03's `validation.md`
  fixture table (3 rows: `persisted-enviado`, `legacy-faturar`, `persisted-novo`) demonstrates a
  real Faturar→Enviado correction for a row whose real shipment status the old code could never see.
- Artifact: `migrations/0088_order_shipments.sql`, `migrations/0089_orders_marketplace_orders_sync_fields.sql`,
  `orders/adapters/postgres/buyer_fiscal_reader.go`, `orders/adapters/postgres/order_repo.go:337-428`,
  F-02/validation.md, F-03/validation.md.
- Blocking failure observed: No.

### M03-C3 — Allowlist encolhe -2 com must-fail
**Status: PASS**

- Command: orchestrator personally re-ran (this session, not delegated) --
  ```
  cd apps/server_core
  # 1. temporarily re-inserted the retired call at root.go:611:
  #    _ = newOrdersShipmentReaderAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)
  GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gocache/mod go test ./internal/platform/archguard/... \
    -run TestRealRepoInteractiveMLSites_MatchesAllowlist -v -count=1
  # 2. reverted (git status confirmed byte-identical to committed tree)
  GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gocache/mod go test ./internal/platform/archguard/... -v -count=1
  ```
- Expected: entradas A/B removidas NO MESMO commit da troca dos readers; reintroduzir a chamada do
  site A quebra o guard NOMEANDO o sítio (output vermelho salvo).
- Actual (RED, captured verbatim by the orchestrator, not a subagent transcript):
  ```
  --- FAIL: TestRealRepoInteractiveMLSites_MatchesAllowlist (0.00s)
      archguard_test.go:426: raw detector found 7 call site(s) in ../../composition/root.go,
      expected exactly 2 (allowlist) + 4 (documented exclusions) = 6; a new, undocumented site
      appeared -- raw hits: [... {Symbol:newOrdersShipmentReaderAdapter File:../../composition/root.go Line:611} ...]
  FAIL
  ```
  Named the exact symbol (`newOrdersShipmentReaderAdapter`) and exact line (`root.go:611`).
  After revert: all 8 archguard tests (`TestRealRepoInteractiveMLSites_MatchesAllowlist`,
  `TestRealRepoTransportAndApplication_NeverImportMLDirectly`, `TestFixture_FifthSiteIsDetectedAndNamed`,
  `TestFixture_AliasedSiteIsDetectedAndNamed`, `TestFixture_ShrunkAllowlistStillPasses`, +3 more) GREEN.
  `mlAllowlist` confirmed 4→2 entries in the SAME commit (`fe2131e7`) as the reader wiring swap
  (`root.go` diff read in full by the orchestrator) -- no interim commit with a dangling entry.
  Site A's constructor is not merely excluded, it's deleted entirely from `orders_adapters.go`
  (confirmed: the must-fail above required temporarily re-inserting a call to a symbol that no
  longer has a definition in the tree, which `parser.ParseFile`'s syntax-only AST scan tolerates).
- Artifact: `internal/platform/archguard/archguard_test.go` (mlAllowlist lines ~here, 2 entries,
  both pricing); this session's terminal capture (above).
- Blocking failure observed: No.

### M03-C4 — Writer único: import path passa pelo IngestOrder
**Status: PASS, with one named non-blocking landmine**

- Command: cold review's repo-wide grep of every consumer of `ordersIngestSvc` in `root.go` +
  the batch-route classification chain (`registerBatchRoutes`/`httpx.BatchRouteClass`).
- Expected: import existente chama o MESMO `IngestOrder` (ADR-04); zero segundo caminho de escrita
  de orders.
- Actual: the only consumer of `ordersIngestSvc` is `ordersImportSvc` (`root.go:595-598`), wired
  into `orderstransport.NewHandlerWithSummary` (`root.go:621`), whose `Handler.Register` maps
  `Import` only to `/orders/import`, classified `httpx.BatchRouteClass` -- confirmed by the cold
  reviewer via repo-wide grep, not just root.go. **Named landmine, non-blocking**: `OrderRepository.UpsertOrders`
  (the pre-F-02 batch write method) still exists, still satisfies `ports.OrderStore`, and is NOT
  called by any production path -- but it IS still called by F-03's own `order_repo_bucket_counts_test.go`
  to seed legacy pre-F-02 fixture rows (a legitimate, understood use, not accidental orphaning).
  Nothing structurally prevents a future feature from calling it directly and reintroducing a second
  writer that skips shipment/0089 population -- this is a real gap in enforcement (no compile-time
  or archguard-style guard makes ADR-04 durable against a future regression), but it is not a CURRENT
  second writer in any reachable production path, so it does not trip the contract's blocking
  condition ("segundo writer de orders sobrevivendo" as an active, wired path). Recommend a
  HARNESS-DEBTS.md entry for a future archguard-style single-writer guard on `orders_marketplace_orders`
  INSERT/UPDATE statements, scoped to whichever milestone next touches this seam.
- Artifact: `internal/composition/root.go:578-621`, `orders/adapters/postgres/order_repo.go:499`
  (`UpsertOrders`), `orders/adapters/postgres/order_repo_bucket_counts_test.go:110-124` (legitimate caller).
- Blocking failure observed: No.

### M03-C5 — Truth table intocada
**Status: PASS**

- Command: orchestrator personally re-ran (this session) --
  ```
  git diff 21ca3595..d22d3d2 -- '*order_bucket*'
  cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gocache/mod go test \
    ./internal/modules/orders/domain/... -run TestDeriveOrderBucket -v -count=1
  ```
- Expected: `TestDeriveOrderBucket` byte-intocado e verde; `DeriveOrderBucket` REUSADO, não re-derivado.
- Actual: `git diff` produced **zero output** -- `order_bucket.go`/`order_bucket_test.go` are byte-identical
  across the entire milestone, base to tip. `go test -run TestDeriveOrderBucket -v` -- all 39 subtests
  PASS. Every read-side consumer this milestone touched (`GetOrderBucketCounts`'s two-tier fallback,
  `IngestOrder`'s ingest-time bucket computation) calls `ordersdomain.DeriveOrderBucket` -- confirmed
  by reading both call sites -- never reimplements bucket logic inline.
- Artifact: this session's terminal capture (above); `orders/domain/order_bucket.go` (unchanged).
- Blocking failure observed: No.

### M03-C6 — Q1: detalhe de pedido rápido (<2s, sem GET vivo no waterfall)
**Status: PENDING — blocked on hub dev-stack repoint**

Requires a browser drive against the actual F-03 binary. `docker inspect marketplace-central-backend-1
--format '{{json .Mounts}}'` (mounts-only, to avoid the documented full-`docker-inspect` secrets-dump
gotcha) confirmed the currently-running backend/frontend containers bind-mount the HUB's main checkout
(`C:\Users\leandro.theodoro\Documents\marketplace-central`), not this worktree -- so they predate this
entire milestone and cannot produce valid evidence for this criterion. REQUEST filed:
`HUB-EVENT-REQUEST-repoint-dev-stack.md` (this directory's parent). Not self-fixable: chip does not
boot servers or rebuild/re-point the dev stack (hub seam).

### M03-U1/U2/U3 — user-drive criteria (operator mandate)
**Status: PENDING — same blocker as M03-C6**

All three require the re-pointed stack. Will be driven and appended to this evidence pack once the
hub signals ENV-READY. Not skippable under this repo's user-drive validation mandate (memory:
"TODO contrato exige critérios M0X-U* dirigidos em browser real") regardless of the LEAN code-review
model this milestone otherwise used for adversarial-review headcount.

## Cold reviewer's full verbatim verdict (milestone-wide diff, 21ca3595..fe2131e7 -- before the
## d22d3d20 dead-code cleanup this verdict itself recommended)

Preserved in full for audit trail; see the reviewer's own text. Summary: no data-corrupting mismatch,
no broken zero-live-ML claim, no split ADR-05 violation, no ADR-17 sentinel bug. One concrete defect
(dead code, `CapabilityAdapter.GetShipmentInfo`/`shipping_reader.go`) -- fixed in `d22d3d20`, independently
re-verified by the orchestrator (build/vet/whole-module test suite green, diff matches the agent's
reported file list exactly). One non-blocking landmine (`UpsertOrders` as a callable-but-unreached
second writer) -- named above under M03-C4, recommend a `HARNESS-DEBTS.md` entry.

## Full-suite re-run (this session, post d22d3d2, GOCACHE absolute)

```
cd apps/server_core && GOCACHE=$(pwd)/.gocache go build ./...   # exit 0
GOCACHE=$(pwd)/.gocache go vet ./...                             # exit 0
GOCACHE=$(pwd)/.gocache go test ./...                            # 117 ok, 0 FAIL
```
Matches the hub-stated baseline ("117 pacotes ok / 0 FAIL") exactly, package-for-package.

## Handoff correction (2026-08-01, per hub message on session `local_c0c3c6c4-9f68-4e6d-ade5-50d23046b13c`)

Hub confirmed `orders_marketplace_orders` has **0 rows** in the (correctly-pointed) dev stack --
independent of the mount-mismatch this chip found. That means M03-C6 and M03-U1/U2/U3 are not
merely blocked on a repoint; they need a real F-02-ingested order that does not exist yet. Per the
hub's own instruction ("U* é drive MEU, não seu"), M03-U1/U2/U3 are EXPLICITLY the hub's drive, not
this chip's. M03-C6 shares the identical precondition (real order + live server) and is deferred to
the hub for the same reason, pending hub confirmation it agrees C6 falls under the same handoff.

## Open items before CLOSED can be sent

1. ~~Hub re-points dev stack~~ -- superseded: hub's own drive covers C6/U1-U3 regardless of repoint,
   since the blocking precondition is 0 ingested rows, not just which binary is mounted.
2. Hub drives M03-C6 + M03-U1/U2/U3 once an order exists; append results to this file (hub-owned append).
3. Recommend (not blocking) a `HARNESS-DEBTS.md` entry for the `UpsertOrders` single-writer
   enforcement gap noted under M03-C4.
