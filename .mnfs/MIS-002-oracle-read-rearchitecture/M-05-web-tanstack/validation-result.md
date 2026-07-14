# M-05-web-tanstack — Validation Result

**Verdict: PASS**
**Frozen SHA:** `c2aea877` (branch `claude/epic-pare-2eaed9`)
**Milestone diff range:** `d413fd78..c2aea877`
**Date:** 2026-07-14

Validated under the **End-of-Milestone Dual Review** protocol (`execution-plan.md`, introduced at `3f48ca91`), which replaced the Claude-dispatched `mpc-verifier`. Both reviewers ran independently at the frozen SHA. Fold rule: a blocking finding from either fails. Neither returned one.

| Reviewer | Model | Verdict |
|---|---|---|
| Codex cross-model review (read-only, fixed-SHA) | `gpt-5.6-sol`, medium | **NO BLOCKING FINDINGS** |
| Claude QA validator (proportional QA vs M-05-C01..C05) | Claude, orchestrator-dispatched | **PASS** |

The cross-model gate is load-bearing here: every line of F-01 and F-02 was written by `gpt-5.6-luna`, so a Luna-family self-review would have been the weakest link. Sol reviewing Luna's output is genuinely independent.

## Criteria

| Criterion | Status | Blocking failure observed |
|---|---|---|
| M-05-C01 — infinite pagination over cursor envelope | **Passed** | No |
| M-05-C02 — staleTime prevents redundant refetch | **Passed** | No |
| M-05-C03 — as_of + manual refresh with no-cache | **Passed** | No |
| M-05-C04 — mutations invalidate correct namespaces; linkage never cached | **Passed** | No |
| M-05-C05 — build green, no bypassing fetches | **Passed** | No |

Key independent confirmations (QA verified premises rather than accepting claims):
- **C01** — the SDK has *no offset parameter at all* (`sdk-runtime/src/index.ts:1055-1057`), so offset-style paging is not expressible; `CatalogPage.test.tsx:56` asserts exactly 2 fetches for 2 pages, proving no page-1 refetch per scroll.
- **C03** — no-cache is proven, not claimed: the test reads the *second* fetch's headers (`CatalogPage.test.tsx:90-91`). Achieved via the apps/web-owned `fetchImpl`; `sdk-runtime` untouched.
- **C04** — linkage `staleTime: 0`/`gcTime: 0` asserted from the *live cached query options*, not from source text. Every invalidation site pairs `toHaveBeenCalledTimes(N)` with exact-namespace assertions; three failure-path tests assert `not.toHaveBeenCalled()`.
- **C05** — `fetch(`/`XMLHttpRequest`/`axios` grep across feature sources: no matches. Inline queryKey literals: no matches. Namespaces are exactly the four reserved roots; none invented.
- **Seams** — `git diff --name-only d413fd78..c2aea877 | grep -E '\.go$|contracts/api/.*openapi|packages/sdk-runtime/src/'` → no matches.

## Evidence at frozen SHA (real exit codes; no `tail`/`head` pipe masking)

| Command | Exit | Result |
|---|---|---|
| `npx vitest run --root . packages/feature-{product-links,inventory,orders,products} --testTimeout=30000` | 0 | 32/32 passed |
| `npm test` (repo root, canonical) | 0 | 11/11 passed |
| `npm run build` (apps/web) | 0 | 1832 modules |

Independently reproduced by QA. QA's 2.69s runtime (vs a discarded 538s contended run) confirms the earlier 5000ms-timeout failures were machine contention, validating that discard.

## Contract interaction — IC-01 "product edit → ['catalog']" is MOOT

F-01 implemented the user-approved Option B and deleted the legacy `ProductsPage`, which held the only product-edit write. QA verified the premise rather than accepting it: `grep` for product-edit SDK methods returns **only the definitions, zero callers**.

Both reviewers judged this acceptable, independently:
- The row is a **conditional** — with no reachable product-edit write, it cannot be violated. Vacuous truth here is real safety.
- **M-05-C04 does not name product edit**; it names linkage confirm and margin-input import only.
- Inventing a product-edit page to satisfy the row would be **fabricating scope to fit a contract** — a validator must reject that, not reward it.

**Forward obligation (routed to Mission Strategist, IC-01's owner):** IC-01 should record the row as *dormant with no client implementation target as of M-05*. The row must **stay** in the contract — if a product-edit surface is reintroduced, the `['catalog']` invalidation obligation binds it. Without that note a future implementer may not read the row as live.

## Non-blocking findings carried forward (do NOT reopen M-05)

1. **`packages/web-query/src/index.ts:72` — `noCacheDepth` is transport-wide** (Sol). While one `withNoCache` operation awaits, any unrelated concurrent GET also receives `Cache-Control: no-cache`. IC-01 scopes forced refresh to the refetched GET; this can unnecessarily bypass L2. *Found only by the cross-model reviewer.*
2. **`CatalogPage.test.tsx:92` under-asserts C03's as_of update** (flagged independently by **both** reviewers). The final `waitFor` matches the same `/dados de \d{2}:\d{2}:\d{2}/` pattern already satisfied before refresh, so an indicator stuck on the original timestamp would pass. Fixtures do differ (`10:11:12Z` → `11:12:13Z`); the assertion should compare distinct rendered values. C03 passes on the substantive condition (the no-cache header assertion).
3. **C02's behavioral remount proof exists only for catalog**; stock (45s) and pricecost (120s) are verified by constant + call-site inspection. Acceptable at QA-2; coverage is asymmetric.
4. **`CatalogPage.test.tsx:73` double-unmounts** `first`; the assertion holds for the right reason but the test is sloppier than it reads.
5. **`listIntegrationInstallations` stays on plain `useEffect`** — installations are Postgres-backed config, not Oracle-backed, outside IC-01's namespaces. Consistent with declared scope.
6. **Product enrichment editing is now unreachable in the UI.** SDK methods survive, so the capability is dormant rather than destroyed — but it is a real user-facing capability removal that shipped under Option B.

## Mission-level hygiene (routed to Mission Strategist, not blocking M-05)

- **Repo typecheck is broken.** F-01 recorded `npx tsc --noEmit` as Blocked (`TS2688: Cannot find type definition file for 'node'`), accepted as a pre-existing baseline. Standing consequence: **`vite build` is esbuild-only and does not typecheck — a green build is not a type-safe build.** This is the exact blindness that let the stale-field ProductsPage survive undetected. Restoring a working typecheck is worth mission attention.
- **`AGENTS.md` drifts from config** — it text-pins `mpc-verifier` to `gpt-5.6-luna`/high, while `.codex/agents/mpc-verifier.toml` is now `gpt-5.6-sol`/medium (user's own change at `6e992d3c`). User's call to reconcile.
- **Seam greps were scoped narrowly** (`\.go$|sdk-runtime|openapi`). They were clean here, but QA notes that is fortunate rather than proven; a broader seam check would be stronger for future milestones.
