# M-02 catalog-batch-cutover — Milestone Checkpoint

```yaml
id: M-02
mission: MIS-002-oracle-read-rearchitecture
type: milestone-checkpoint
verdict: PASS
frozen_sha: fc8342b4
integrated_sha: 648adbdf
base_sha: a714acf7
branch: claude/jovial-haibt-29f2f0
merged_to_main: false
closed: 2026-07-14
```

## Verdict

**PASS** — proportional QA (mpc-verifier, gpt-5.6-luna high) issued binary PASS against `validation-contract.md` M-02-C01..C06 on integrated SHA `648adbdf`; evidence committed at `fc8342b4` (`validation-result.md`). Per-lane fixed-SHA reviews PASS (F-01, F-02), integrated fixed-SHA review PASS with no findings.

## Outcome

Catalog listing cut from 1+3N sequential Oracle calls to exactly ONE set-based keyset-paginated query per page, served via the IC-01 cursor envelope on the existing `/catalog/*` namespace; OpenAPI + `sdk-runtime` updated in the same feature commit; old N+1 composition removed from the catalog listing flow.

## Integrated commits (base a714acf7)

| SHA | Scope |
| --- | --- |
| 68b6917b | F-01 step 1 — catalog page port interface + domain page types (`internal_read/ports/catalog_page.go`); interface-first handshake commit |
| 8b75fa4e | F-01 step 2 — Oracle adapter (single JOIN over base TGF* tables per M-01-C04 FULL_SCAN verdict, keyset by internal product id, FETCH FIRST :limit+1 peek), nullable facts + quality flags, TimingReader + Service forwarding, catalog listing switched to new port, fake-queryer tests |
| cbd4d009 | F-02 — transport routes (list + `/catalog/products/search`) IC-01 envelope, cursor/limit validation + error matrix, Cache-Control→FreshnessPolicy mapping, OpenAPI + `packages/sdk-runtime` typed client in SAME commit |
| 88c14677 | merge — integrate F-02 lane into F-01 lane (disjoint seams, clean) |
| 648adbdf | integration — wire `internalReadSvc` as catalog `Handler.PageReader` in `composition/root.go` (nil when Oracle unavailable → source_unavailable) |
| fc8342b4 | docs — QA PASS `validation-result.md` |

## Criteria (QA PASS)

- **C01** one Oracle query per page (fake counting queryer; sizes 1/50/100; no N+1) — PASS
- **C02** envelope + cursor chain conform to IC-01 (keys exactly items/next_cursor/page_size/as_of; null on last page; decimal strings; RFC3339 as_of; non-overlapping/gapless) — PASS
- **C03** error matrix (400 invalid_cursor no Oracle call; 400 invalid_limit for 0 and 101; 503 source_unavailable, no driver leak) — PASS
- **C04** unknown facts null + quality flags (missing_stock/missing_price/missing_cost; duplicate active price → ambiguous_price, page still 200; never zero) — PASS
- **C05** OpenAPI + sdk-runtime same commit (cbd4d009); sdk build green; old unpaginated listing deprecated per IC-01 — PASS
- **C06** search bounded 1..50, next_cursor null, sorted internal_product_id asc, 1 query; limit=51 → 400 invalid_limit — PASS

Whole suite `go build ./... && go test ./...` green (GOCACHE=.gocache absolute). Boundaries intact: SQL/driver types confined to `internal_read/adapters/oracle`; redaction (wrapOracleError/safeOracleCause) preserved; no forbidden seams touched (inventory/profitability/product_links/oraclebatch untouched).

## Seam ownership honored

Owned: catalog module, transport catalog routes, OpenAPI, sdk-runtime, new internal_read port/adapter, and the declared conflict point `composition/root.go` (wired mine; M-03 rebases later). No foreign seam touched.

## OPEN — merge to main deferred to Hub (user directive)

Milestone NOT merged to `main`. User directed the Portfolio Hub to handle the main merge.

- Default branch is **`main`** (tip a714acf7), not `master` as the goal packet stated.
- Milestone work is committed on branch `claude/jovial-haibt-29f2f0` @ `fc8342b4`, which descends cleanly from `main` (a714acf7) — trivial `--no-ff` merge, no conflicts expected.
- Primary `main` worktree currently dirty with 25 unrelated uncommitted files (MIS-001, ARCHITECTURE.md, governance configs, docker, docs, `.codex/agents/mpc-verifier.toml`) — none overlap the M-02 seam.
- M-02 was to merge FIRST, before sibling M-03. Hub owns sequencing.

## Next (Hub)

Hub merges `claude/jovial-haibt-29f2f0` @ `fc8342b4` `--no-ff` into `main` when the shared worktree is safe, ahead of M-03; then proceed with M-03/M-04 per roadmap.

## Lane-B worktree

`.claude/worktrees/m02-f02-lane` (branch `m02/f02-envelope`) still registered; fully merged into the milestone branch — Hub/user may `git worktree remove` it after the main merge.
