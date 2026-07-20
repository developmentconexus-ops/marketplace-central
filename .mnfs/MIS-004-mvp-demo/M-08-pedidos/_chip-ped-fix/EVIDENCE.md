# CHIP-PED-FIX — Evidence Pack

**Mission** MIS-004-mvp-demo · **Milestone** M-08-pedidos · **Demo day** 2026-07-20 (operator #1 priority: "o foco é pedido")
**Branch** `claude/jolly-bohr-5b0d5b` · **Base** `9e0beb41` · **Tip** `9ff1871d` (unpushed, NOT merged — hub owns merge/P7/live re-import)
**Chip status** CLOSED → hub (this pack) · **Live** LIVE-VERIFIED (partial, see §6) + one deploy-gap FINDING routed to hub

---

## 1. Deliverables (four goals, all implemented)

| Goal | What | Commit |
|---|---|---|
| **A** Perf | Serial per-order LIVE shipment GET N+1 in EnrichService → bounded-concurrency `errgroup.SetLimit(8)`, slice order preserved (distinct-index writes), per-order shipment error degrades honestly (nil, no batch abort/reorder). | `6ba9885f` |
| **C** Decomposição (ADR-17) | `decomposicao` computed from REAL persisted data: Comissão = Σ item `sale_fee_amount`; Custo = Σ ERP UnitCost×qty; Frete = shipment SenderCost. Unknowns stay `nil` → honest "—", never fabricated 0. Imposto explicitly `nil` (not persisted). | `6ba9885f` |
| **B** Faturado workflow | Migration `0074` (nullable `faturado_at timestamptz`); `DeriveOrderBucket` gains faturado dimension (paid+!faturado→**Faturar**, paid+faturado→**Enviar**, shipped/delivered ALWAYS→**Enviado**); `POST /orders/{provider_order_id}/faturado` (tenant+installation scoped, idempotent first-write-wins, OUR-DB write only, **ZERO ML write**); FE "Foi faturado" button gated to `bucket==="faturar"`. | `63e85257` |
| **contract** | Additive `POST /orders/{id}/faturado` (OpenAPI + `sdk-runtime.markOrderFaturado` in lockstep) + `faturado_at?` on OrderRead (see §4). | `937d3512`, `9ff1871d` |
| **D** Refresh | Header "Atualizar" reuses `/orders/import` (a READ from ML into our DB — never an ML write) + refetch; loading state; honest failure line. | `c0d674fb` |

### Commit list (`git log 9e0beb41..HEAD`)
```
9ff1871d fix(orders): document faturado_at on OrderRead per hub contract grant   <- FINAL C7 resolution
e3a0769a fix(orders): keep faturado_at off the OrderRead wire (P6 cold-gate C7)  <- superseded direction (see §4)
c0d674fb feat(web/pedidos): Atualizar refresh + working "Foi faturado" action
937d3512 feat(contract): additive POST /orders/{id}/faturado + sdk markOrderFaturado
63e85257 feat(orders): faturado workflow — bucket gate + POST /orders/{id}/faturado
6ba9885f perf(orders): parallelize shipment N+1 + real-data decomposição
```

---

## 2. Verification of record

| Lane | Result |
|---|---|
| `go build ./...` | OK |
| `go test ./internal/modules/orders/... ./internal/platform/migrate/...` | all `ok` (integrations, internalread, postgres, productlinks, application, domain, ports, transport, migrate) |
| `go vet ./internal/modules/orders/...` | clean |
| `sdk-runtime` `tsc --noEmit` | clean |
| `sdk-runtime` `vitest run` | **71 passed** (4 files) |
| OpenAPI parse | **90 paths**; `faturado_at` OPTIONAL on OrderRead (in properties, NOT in `required`); `markOrderFaturado` responses 204/400/404, requestBody `installation_id` required |
| `apps/web` `PedidosPage.test.tsx` | **28/28** (prior window; FE source unchanged by the §4 contract realign — that touched only Go/openapi/sdk) |
| `apps/web` `tsc` | 12 baseline errors, **ZERO in touched files** (prior window; baseline per memory `web-tsc-lane-cross-branch-resolution`) |

Env: hermetic `GOCACHE`/`GOMODCACHE` = absolute worktree paths; run from `apps/server_core` (gomodcache-root gitignore trap); `.gocache`/`.gomodcache` never staged (verified `git diff --cached` clean of both).
`-race` unavailable (no cgo toolchain on Windows) — concurrency (goal A) is race-free by construction (disjoint slice-index writes, errgroup, no shared mutable state); independently confirmed by the refuter (§5, R3).

---

## 3. Constraints honored (binding)

- **ZERO Mercado Livre writes.** `faturado` = a single Postgres `UPDATE orders_marketplace_orders SET faturado_at=now()`; `Atualizar` = existing `/orders/import` (READ from ML). Self-scan + cold gate C1 + refuter R1 all confirm no provider mutation on either path. Verified live: no provider write issued (§6).
- Did **NOT** merge to main / push / reset / revert / stash / clean.
- Did **NOT** install deps — `errgroup` (`golang.org/x/sync`) already vendored (no dep REQUEST needed).
- Did **NOT** rebuild or mutate the hub-owned dev stack. The live check (§6) issued **GET reads only** against `:8080`; no import/faturado POST was fired against shared infra.
- Did **NOT** self-stamp the hub's authoritative `P6-DUAL-GATE: AGREEMENT` marker (hub owns it + merge + P7 + live re-import). §5 records the CHIP-SIDE gate result only.
- No secrets/tokens/PII printed or committed.

---

## 4. C7 contract resolution (cold-gate finding → hub grant)

The cold P6 gate caught real drift: `enrichedOrderDTO` (http_handler.go:478) anonymously **embeds** `domain.OrderReadModel`, so Go's marshaller promoted the new `FaturadoAt` field onto every enriched `/orders` response — a wire field undocumented in OpenAPI/SDK `OrderRead`.

- First fix `e3a0769a` suppressed it (`json:"-"`) — drift-free, gates passed.
- **Superseded** by `9ff1871d` per an explicit HUB contract grant (seam owner): "faturado_at on OrderRead" in openapi + "OrderRead.faturado_at?: string" in sdk, "exactly as scoped". The original scope always had `faturado_at` as an OrderRead field; the true bug was emit-WITHOUT-document. Final resolution = **document + emit**:
  - `read_model.go`: `FaturadoAt *time.Time json:"faturado_at,omitempty"` (RFC3339 when set, key absent when nil = honest ADR-17 absence).
  - OpenAPI `OrderRead`: optional `faturado_at` (string date-time, not in `required`).
  - `sdk-runtime` `OrderRead`: `faturado_at?: string` — matches the openapi optional + Go omitempty semantics (absent-when-unknown, string-when-set).
- **wire == contract == SDK** verified across all three surfaces (cold gate round-3, §5).

---

## 5. P6 dual gate (chip-side; final SHA 9ff1871d)

Two independent reviewers over the **whole orders module surface** (not just the diff).

**Cold gate** — `harness:gate-reviewer`, Opus, physically read-only:
- Round-1: **FAIL** — sole blocker C7 (undocumented `faturado_at` on wire). C1–C6, C8 all **PASS** (incl. hardest constraint C1 ZERO-ML-write; C2 tenancy tenant+installation on UPDATE and every SELECT; C3 idempotency first-write-wins + EXISTS→404; C4 bucket precedence shipped/delivered-wins; C5 ADR-17; C6 order-preserving concurrency).
- Round-3 (final shape): **PASS** — C7 **TRUE**, all three surfaces consistent, guard proven load-bearing both directions.

**Adversarial refuter** — `general-purpose`, sonnet, prompted to REFUTE:
- Round-1 (R1–R9): no defect confirmed — ML-write leak, tenancy, concurrency/order-corruption, bucket regression, idempotency/404, contract drift on the endpoint, ADR-17, test-theater, migration all COULD-NOT-REFUTE. (It MISSED C7 — checked the endpoint triangle, not the embedded-struct promotion onto the response DTO; self-noted.)
- Round-2 (attack the fix): **COULD-NOT-REFUTE** + no new defect. Confirmed bucket derivation still fed by the Go field (tag affects only marshalling), no residual leak on any path, guard genuinely load-bearing, no consumer breakage.

**CHIP-SIDE VERDICT: AGREEMENT** — both reviewers concur the module is correct and C7 is resolved on the final committed shape. The refuter's substantive attack surface (tenancy / concurrency / idempotency / bucket / migration) is invariant to the `faturado_at` tag choice, so its round-2 endorsement carries to the final shape; the cold gate round-3 independently confirmed the document+emit surface. *(Authoritative P6-DUAL-GATE marker reserved to the hub.)*

**Load-bearing guards proven must-fail on revert** (`TestMapEnrichedOrderDerivesBucketFromFaturadoAt`):
- faturado order (non-nil) MUST emit `faturado_at == "2026-07-20T09:00:00Z"` → revert to `json:"-"` drops the key, test FAILS (observed: `json field faturado_at = <nil>`).
- not-invoiced order (nil) MUST OMIT the key → a non-omitempty tag emits `null`, test FAILS.

---

## 6. LIVE-VERIFIED (throwaway vite :5199 → shared `:8080` backend, GET-only)

Throwaway vite dev server (`:5199`, non-5174/8080 per charter), `/pedidos`, installation `inst-mercado_livre-d373dc64…`. Backend = shared hub-owned dev stack at `:8080` (vite proxy).

- **Atualizar button** — renders live; observed loading state `Atualizando…` → resolved `Atualizar`. ✓
- **List loads live** — 24 real orders, real buyers, honest "—" bucket chips; KPIs NOVOS 0 / A FATURAR 0 / A ENVIAR 0 / ENVIADOS 24 / DIFAL A PAGAR "—". ✓
- **Foi faturado gating** — golden order `2000017276984774` opened; bucket **`enviado`** (rastreio `delivered`) → **no "Foi faturado" button** (only inert `DIFAL agendar` / `Devolução…`). Correct negative gating. ✓ *(Positive faturar-case not exercisable live — 0 faturar-bucket orders in current data; covered by RTL 28/28.)*
- **Decomposição ADR-17 render** — all components honest "—", no fabricated 0. ✓
- **Additive-optional contract is backward-compatible** — the `:8080` (old-code) response carries NO `faturado_at`; the new SDK `faturado_at?: string` (optional) parses it and the page renders fine. ✓

### FINDING → HUB (hub point C, golden comissão=22,95)
The golden order's raw `/orders/2000017276984774` body has item **`sale_fee_amount: 22.95`** (the hub's exact number is present in the data), but `decomposicao.componentes_desconhecidos` lists `comissao` and `decomposicao.comissao` is **null** → FE shows "—".
- This is **NOT a silent key mismatch** (the hub's worry). Keys verified identical: DTO `json:"comissao"` (http_handler.go:422) ↔ FE `decomposicao.comissao` (PedidoDrawer.tsx:143), all 10 decomposição keys match 1:1.
- It is a **DEPLOY GAP**: the `:8080` stack runs **pre-goal-C** server code that does not sum `sale_fee_amount` into `decomposicao.comissao`. My goal-C code (`6ba9885f`, this worktree, NOT deployed to `:8080`) computes `Comissão = Σ sale_fee_amount`. Pinned by unit tests: `order_decomposition_test.go:71-77` (`Comissao == 22.95` surfaced verbatim) + `enrich_service_test.go` (sale_fee→Comissão summation).
- **Live comissão=22,95 render requires the hub to rebuild the `:8080` stack on the merged code (hub-owned).** The input data already carries `sale_fee_amount:22.95`, so once deployed it will compute and render.

---

## 7. FINDINGS (for hub ratification)

1. **Golden deploy-gap** (§6) — `:8080` runs pre-goal-C code; comissão renders "—" until the stack is rebuilt on merged code. Not a code defect; key-parity verified.
2. **Migration-count fixture** — `runner_test.go` asserts the real on-disk inventory of **62** `.sql` files (incl. `0074`), not the stale "58/59" in the original chip card. Honest count, not test theater.
3. **Go compiler flake** — one `go test` run hit an internal Go compiler segfault (`asm_amd64.s`) on Windows under parallel compile; a clean re-run passed. Non-deterministic tooling flake, not a code defect.

---

## 8. Handoff to hub

Merge (`--no-ff`), P7 browser QA, and the live re-import / `:8080` stack rebuild are hub-owned. After rebuild, run the golden live-drive: order `2000017276984774` → drawer → `decomposicao.comissao` should render **R$ 22,95** (input `sale_fee_amount:22.95` confirmed present). Contract is additive-only within the orders section (disjoint from CHIP-IMPORT-FIX's erp_import section).
