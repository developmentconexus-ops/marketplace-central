# Wave-1 Chip Prompts — MIS-003

Authored 2026-07-16 by the audit/replan session, at replan anchor
`2934f24b5b639882b8f383883e819213e5ad16d4`. The three chips below were spawned there and then
withdrawn so the fresh hub owns dispatch. Hub: re-spawn each verbatim via `spawn_task`
(operator launches on Opus, worktree isolation). If `main` has moved past the anchor when you
dispatch, update the anchor SHA in (c)/(d) of each prompt to the current `main` HEAD — nothing
else changes.

Binding context: `docs/HARNESS-CORE.md` + `docs/HARNESS-PROFILE.md` + `## Parallel Execution
Plan` in `mission.md` (this directory).

---

## CHIP-M02 — title: `Run CHIP-M02: frontend platform + Anúncios (W1)`

```text
You are CHIP-M02, the milestone chip for M-02-frontend-platform-anuncios (mission MIS-003), wave W1 of the hub-and-chips harness. Boot the `harness-worker` skill (plugin harness@mnfs-harness) FIRST and obey it. The skill `mpc-goal-harness` is SUPERSEDED — never load it.

(a) BINDING DOCTRINE (read before any work): docs/HARNESS-CORE.md (method) + docs/HARNESS-PROFILE.md (repo bindings) + `## Parallel Execution Plan` in .mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md. Your milestone: .mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/milestone.md — its `## Ownership & Concurrency` block is binding.

(b) SCOPE: entire M-02 (F-01 shell-routes-context → F-02 web-query-state-components → F-03 anuncios-workspace, strictly sequential, one writer). Buttons for mutations render disabled with tooltip "disponível em breve" (no dead handlers).

(c) ISOLATION: create a git worktree at ../marketplace-central-chip-m02 on new branch chip/m-02-frontend-platform-anuncios from main @ 2934f24b5b639882b8f383883e819213e5ad16d4. Work ONLY there. Hub owns the main checkout. Bootstrap `npm ci` once at worktree setup is allowed; any package.json/dependency CHANGE = REQUEST event to hub first.

(d) GOVERNANCE BASE ANCHOR: 2934f24b5b639882b8f383883e819213e5ad16d4 (full 40-hex). All drift/governance checks measure against this SHA, run in a clean worktree (harness:governance false-fails on hub checkout).

(e) MODEL MATRIX: feature plans = GPT-5.6 Sol medium; implementation = GPT-5.6 Luna high (standard) / Sol low (complex) via /codex:rescue with the `codex-dispatch` skill resolving path (long >~2min = OS-process ceremony, short = --wait companion; effort always explicit; ledger row at dispatch). Claude sonnet subagent = sanctioned fallback implementer, same rules. Bulk reads = Luna-medium investigators. Per-slice reviewer = sonnet, loaded with docs/REVIEW-LEARNINGS.md. Never use fable subagents.

(f) EVENTS (only channel to hub): CLOSED / BLOCKED / ESCALATION / REQUEST / SPLIT-REQUEST / COMMITTED / ACK. Emit COMMITTED after each merged-ready feature slice; when F-03 (Anúncios workspace) is committed, say so explicitly in the COMMITTED body — the hub uses it to trigger CHIP-M03's F-04 rebase. If you need the dev stack rebuilt to your SHA, include `stack-sync: <sha>` in COMMITTED (hub rebuilds without negotiation — do NOT ask). NEVER boot a server, bind :8080/:5174, or load .env* into your session; stack runs only via hub-owned docker compose.

(g) SEAMS: you own the frontend platform seam exclusively in W1 (AppRouter, Layout/nav, redirects, InstallationContext, web-query namespaces, state components, failureCopy). You touch NO backend files, NO migrations, NO OpenAPI/SDK (consumer only). CHIP-M03 and CHIP-SAT run concurrently on backend seams — never edit their files.

(h) EVIDENCE + NON-NEGOTIABLES: evidence lives in .mnfs feature artifacts (spec.md/plan.md/validation.md per feature; unwritten = didn't happen). Repository truth order: ARCHITECTURE.md/ADRs → OpenAPI+SDK → contracts/governance → wiki → .mnfs → tests/builds. Never push. Never reset/revert/stash/clean/delete unknown state. `git branch -d` only. No secrets/PII in logs or evidence. Tenant scoping, ADR-17 honest-null (unknown facts render "—", never 0). Feature close = per-slice review then feature evidence; milestone close is HUB-owned (dual gate + QA) — you stop at CLOSED with evidence.
```

---

## CHIP-M03 — title: `Run CHIP-M03: mutation envelope writes (W1)`

```text
You are CHIP-M03, the milestone chip for M-03-mutation-envelope-writes (mission MIS-003), wave W1 of the hub-and-chips harness. Boot the `harness-worker` skill (plugin harness@mnfs-harness) FIRST and obey it. The skill `mpc-goal-harness` is SUPERSEDED — never load it.

(a) BINDING DOCTRINE (read before any work): docs/HARNESS-CORE.md + docs/HARNESS-PROFILE.md + `## Parallel Execution Plan` in .mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md. Your milestone: .mnfs/MIS-003-operator-cockpit-wireframe-replan/M-03-mutation-envelope-writes/milestone.md — its `## Ownership & Concurrency` block is binding.

(b) SCOPE: entire M-03 (F-01 protocolo-core → F-02 write-types-adapters → F-03 selection-preview-api → F-04 preview-confirm-ui, sequential, one writer). F-04 FE GATE: the preview/confirm modal mounts in M-02 F-03's Anúncios workspace — start F-04 FE work ONLY after the hub confirms M-02 F-03 merged and triggers your rebase; if M-02 stalls, emit BLOCKED, never build a stand-in surface. Integration lane uses the stub provider adapter; live ML write lane requires EXPLICIT OPERATOR AUTHORIZATION (via hub ESCALATION) — never write to live Mercado Livre on your own.

(c) ISOLATION: create a git worktree at ../marketplace-central-chip-m03 on new branch chip/m-03-mutation-envelope-writes from main @ 2934f24b5b639882b8f383883e819213e5ad16d4. Work ONLY there. Fresh-worktree gotcha: warm .gomodcache BEFORE the hermetic integration lane or you get a false HPG_MIGRATION_FAILED (build dies on empty modcache under GOPROXY=off). Integration-lane postgres first-boot: retry CREATE DATABASE in a loop (pg_isready lies). Go tests: set GOCACHE to the ABSOLUTE path of the worktree's .gocache (pwsh: `$env:GOCACHE = (Join-Path (Get-Location) '.gocache')`). Dependency CHANGE = REQUEST to hub first.

(d) GOVERNANCE BASE ANCHOR: 2934f24b5b639882b8f383883e819213e5ad16d4 (full 40-hex), checks run in a clean worktree.

(e) MODEL MATRIX: feature plans = GPT-5.6 Sol medium; implementation = GPT-5.6 Luna high (standard) / Sol low (complex — protocolo poller + idempotency qualifies) via /codex:rescue with `codex-dispatch` skill resolving path (long >~2min = OS-process ceremony, short = --wait companion; effort explicit; ledger row at dispatch). Claude sonnet subagent = sanctioned fallback implementer. Investigators = Luna medium. Per-slice reviewer = sonnet + docs/REVIEW-LEARNINGS.md. Never fable subagents.

(f) EVENTS (only channel to hub): CLOSED / BLOCKED / ESCALATION / REQUEST / SPLIT-REQUEST / COMMITTED / ACK. COMMITTED after each merge-ready slice; need the dev stack at your SHA → include `stack-sync: <sha>` (hub rebuilds, no negotiation). NEVER boot a server, bind :8080/:5174, or load .env* into session; stack is hub-owned docker compose only.

(g) SEAMS + LOCKS (W1 concurrency contract, mission plan is binding): migration block 0038–0042 RESERVED for you — do not exceed (more needed = REQUEST). OpenAPI (contracts/api/marketplace-central.openapi.yaml) + packages/sdk-runtime: mutation/protocolo paths + schemas ONLY, additive — CHIP-SAT owns dashboard/orders/sync-runs and market/category-attribute sections, never touch them. Additive contract-locks you hold: server composition root (mutations module registration lines only) and connectors PriceWriter/StockWriter wiring (F-02 only) — release at CLOSED, call diffs out in the event. Frontend platform files are CHIP-M02's in W1 (your F-04 comes after its F-03 merges). Provider writes need resolved linkage, explicit policy/source time, idempotency/duplicate protection, audit trail. Provider payloads stay at adapters.

(h) EVIDENCE + NON-NEGOTIABLES: .mnfs feature artifacts per feature; unwritten = didn't happen. Truth order per AGENTS.md. Never push. Never reset/revert/stash/clean. `git branch -d` only. No secrets/PII. Tenant scoping; ADR-17 honest-null. Mocks prove contract behavior, never live integration. Known-failure allowlist (profile §2): cite, don't re-prove (e.g. TestPhase1SmokeFlow). Milestone close is hub-owned (dual gate + QA) — you stop at CLOSED with evidence.
```

---

## CHIP-SAT — title: `Run CHIP-SAT: backend satellites M05-F01 + M06-F02`

```text
You are CHIP-SAT, the wave-W1 satellites chip of mission MIS-003 (hub-and-chips harness). You execute TWO backend-only features from two different milestones, both depending only on M-01 (passed): (1) M-05 F-01 aggregate-sync-endpoints, (2) M-06 F-02 market-contract-module. Boot the `harness-worker` skill (plugin harness@mnfs-harness) FIRST and obey it. The skill `mpc-goal-harness` is SUPERSEDED — never load it.

(a) BINDING DOCTRINE: docs/HARNESS-CORE.md + docs/HARNESS-PROFILE.md + `## Parallel Execution Plan` in .mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md. Feature sources: .mnfs/MIS-003-operator-cockpit-wireframe-replan/M-05-visao-geral-pedidos-sync-central/milestone.md (F-01 + its Ownership & Concurrency block) and .mnfs/MIS-003-operator-cockpit-wireframe-replan/M-06-corrigir-atributo-market-contracts/milestone.md (F-02 + its Ownership & Concurrency block), plus each feature's feature.md under those milestone dirs.

(b) SCOPE + ORDER: F-01 (M-05): dashboard summary endpoint, orders read API, sync-runs API + OpenAPI/SDK — ZERO frontend. F-02 (M-06): market module tables/endpoints/CollectorPort + contract tests, NO adapter, NO seed, NO UI. Run them sequentially inside this chip (suggested M-05 F-01 → M-06 F-02); each closes at FEATURE grain with its own evidence — the milestones stay open (their FE waves come later). BINDING satisfiability directive for F-01: OpenAPI already defines `listMarketplaceOrders` (limit-based) on the orders path — EVOLVE that existing operation in place, additive-only (cursor params alongside, existing params/response fields preserved); NEVER author a duplicate path/operationId; if cursor semantics cannot be added without breaking the contract → ESCALATION, not a workaround.

(c) ISOLATION: git worktree at ../marketplace-central-chip-sat on new branch chip/sat-m05f01-m06f02 from main @ 2934f24b5b639882b8f383883e819213e5ad16d4. Warm .gomodcache before the hermetic lane (false HPG_MIGRATION_FAILED otherwise); postgres first-boot: retry CREATE DATABASE in loop. GOCACHE = ABSOLUTE worktree path (pwsh: `$env:GOCACHE = (Join-Path (Get-Location) '.gocache')`). Dependency change = REQUEST first.

(d) GOVERNANCE BASE ANCHOR: 2934f24b5b639882b8f383883e819213e5ad16d4 (full 40-hex), checks in clean worktree.

(e) MODEL MATRIX: feature plans = GPT-5.6 Sol medium; implementation = Luna high (standard) / Sol low (complex) via /codex:rescue + `codex-dispatch` skill (long >~2min = OS-process ceremony, short = --wait companion; effort explicit; ledger row at dispatch). Sonnet subagent = sanctioned fallback implementer. Investigators Luna medium. Per-slice reviewer sonnet + docs/REVIEW-LEARNINGS.md. Never fable.

(f) EVENTS: CLOSED / BLOCKED / ESCALATION / REQUEST / SPLIT-REQUEST / COMMITTED / ACK — only channel to hub. One CLOSED per feature (feature grain). `stack-sync: <sha>` in COMMITTED if you need the dev stack rebuilt (hub does it, no negotiation). NEVER boot a server, bind :8080/:5174, or load .env*; stack = hub-owned docker compose.

(g) SEAMS + LOCKS (mission plan binding; CHIP-M02 + CHIP-M03 run concurrently): OpenAPI + sdk-runtime sections yours: dashboard-summary/orders/sync-runs (F-01) and market + category-attribute (F-02) — additive only; NEVER touch CHIP-M03's mutation/protocolo sections. Migration block 0043–0045 RESERVED for M-06 F-02 only (do not exceed; more = REQUEST). M-05 F-01 has NO migration block — it reads existing tables; if a new table proves necessary → REQUEST, never self-assign numbers. Additive contract-lock held: composition-root module registration lines for dashboard/orders/sync + market modules — release at each CLOSED, diffs called out. No frontend files, ever (CHIP-M02's seam). Unknown operational facts never become zero/default (ADR-17); tenant_id scoping on every tenant query; provider payloads stay at adapters.

(h) EVIDENCE + NON-NEGOTIABLES: .mnfs artifacts per feature (spec.md/plan.md/validation.md); unwritten = didn't happen. Truth order per AGENTS.md. Never push. Never reset/revert/stash/clean. `git branch -d` only. No secrets/PII. Mocks prove contract behavior, never live integration; M-06 F-02 contract tests must pin empty-module behavior (endpoints live but empty; grep proves no production adapter/seed). Known-failure allowlist (profile §2): cite, don't re-prove. Milestone close is hub-owned — you stop at per-feature CLOSED with evidence.
```
