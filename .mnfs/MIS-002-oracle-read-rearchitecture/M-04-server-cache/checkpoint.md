# M-04-server-cache — Milestone Checkpoint

```yaml
id: M-04-checkpoint
type: milestone-checkpoint
parent: M-04-server-cache
mission: MIS-002-oracle-read-rearchitecture
owner: Milestone Orchestrator
created: 2026-07-14
updated: 2026-07-14
status: complete
verdict: PASS
reviewed_sha: 56d5be9ecb02c361b6df889d1baf49ec3b1c48c2
base_sha: d413fd78
branch: claude/focused-borg-9c9811
merged_to_main: false
dual_review: complete (Codex gpt-5.6-sol PASS + qa-validator PASS, both at 56d5be9e)
```

## Verdict

**PASS.** Both halves of the End-of-Milestone Dual Review passed at frozen SHA `56d5be9e`.
All five criteria pass. No blockers. Hub owns the main merge — this milestone did **not** self-merge.

| Criterion | Verdict | QA mutation result |
| --- | --- | --- |
| M-04-C01 TTL hit/expiry per class | PASS | TTL check removed → `TTLPerClass` FAIL |
| M-04-C02 Singleflight collapses misses | PASS (residual) | collapsing defeated → `ErrorNotCached` FAIL 20/20 |
| M-04-C03 no-cache bypass + linkage exclusion | PASS | 3 mutations → all FAIL |
| M-04-C04 Bounded memory + observability | PASS | LRU disabled → FAIL; raw key leak → FAIL |
| M-04-C05 Evict-on-mutation | PASS | `InvalidateClass` deleted → composed AND direct FAIL |

Neither contract blocking failure is observed: linkage is never served from cache; bypass is never
ignored; no stale L2 entry survives a successful mutation.

## Process Failure — read this before the technical content

**The first terminal callback (at `df2cac6a`) was wrong and the hub was right to bounce it.**
The milestone required a dual review — Codex gpt-5.6-sol **and** Claude qa-validator. I ran only the
QA half and reported terminal on it. The hub ran the missing half, which returned **FAIL with 3
blocking findings**. I independently verified all 3 against the repo: all real. One half of a
two-half gate is not the gate.

Worse, two of the three were the exact defect class this milestone had already spent four rounds
fighting, and my own checkpoint had *named* the seam in question as "load-bearing and easy to miss"
— then shipped without a test crossing it. Writing the warning is not the same as acting on it.

## The Three Findings (all real, all now closed)

**F1 — contract retry policy stale.** `validation-contract.md` said `correction_attempts: 0`,
`last_validation_result: none` against five real rounds. Now reconciled: 6 attempts, max 2, with an
explicit `overrun:` line recording that rounds 3-6 exceeded the max and were user-authorized and
disclosed. Claims no result at the new SHA that had not yet happened.

**F2 — M-04-C03 had no composed HTTP test.** The contract (line 66) literally requires
`httptest: warm cache -> Cache-Control: no-cache`. Only the two seams existed separately:
handler→MaxAge=0 (against a *fake* reader that merely recorded the policy) and MaxAge=0→bypass
(cache unit test, no HTTP). Nothing crossed the seam.

This was substance, not ceremony. My own mutation proved it: neutering the bypass branch in
`load()` — simulating the transport/domain context-key regression this milestone created the risk
of — made the new composed test **FAIL** (`bypass Oracle calls=1, want 2`) while the old
`TestCatalogSearchPageEnvelopeAndNoCachePolicy` **passed green**. A completely dead
`Cache-Control: no-cache` header would have shipped invisible to the old evidence.

**F3 — M-04-C05 success path entirely unguarded.** The only orders invalidation test asserted
*zero* invalidations on the **failure** path. Deleting the production `InvalidateClass("catalog")`
left the whole suite green. A guard that covers only the negative case is a hollow test wearing a
disguise — the fifth of this milestone.

## Fix (correction 6, commit `56d5be9e`) — test/evidence/contract only

Production `.go` diff vs `df2cac6a` is **empty**. Verified independently by me, by Codex sol, and by QA.

- NEW `apps/server_core/tests/unit/cache_composed_test.go`:
  - `TestComposedCatalogHTTPNoCacheBypassesAndRepopulates` — real catalog transport handler → real
    cache decorator → fake Oracle reader. Warm (1 call, identical `as_of`) → `no-cache` (2 calls,
    strictly newer `as_of`) → ordinary GET (still 2 calls, returns the **bypass-refreshed** `as_of`,
    read from the HTTP response body).
  - `TestComposedLinkageConfirmEvictsCatalogButPreservesPriceCost` — warm catalog via real handler →
    real successful `AssistedSankhyaLinkageService.Confirm` wired with the **real** `*cache.Cache` as
    Invalidator → repeat GET inside TTL ⇒ fresh Oracle call + newer `as_of`, while a warm `pricecost`
    key survives (proves class-scoped, not blanket, invalidation).
- NEW `TestAssistedSankhyaConfirmInvalidatesCatalogAfterSuccessfulPersistence` — asserts recorded
  classes equal **exactly** `["catalog"]`.
- `validation-contract.md` Retry Policy reconciled; `validation.md` updated and two overstatements
  **retracted** rather than defended.

## Verification

Three independent parties, all mutation-based. Nobody trusted anyone's word.

- **Orchestrator (me):** 2 mutations. `InvalidateClass` deleted → both new C05 tests FAIL. Bypass
  branch neutered → composed C03 FAILS while the old test PASSES. Restored to empty diff each time.
- **Codex gpt-5.6-sol:** PASS. F1/F2/F3 FIXED with file:line evidence. Hollow-test sweep: all three
  new tests behavior-sensitive. Confirmed eviction-not-TTL (1s advance vs 5m TTL). No new defects.
- **qa-validator:** PASS. **9 mutations, all 9 caught.** Its sharpest result: making bypass fetch
  fresh but skip `storeIfCurrent` caught step 3 serving the stale original
  (`repopulated as_of=16:00:00, want 16:00:01`) while the old cache-only test stayed green —
  empirically confirming `cache_test.go:585` is insensitive and superseded. Under the C05 mutation
  the post-confirm log emitted `cache=hit key_class=catalog`, proving the entry was still live and
  the refresh is genuine eviction.

Suite at `56d5be9e`: full `go test ./... -count=1` green (63 packages), `go vet ./...` clean,
production diff empty at handoff.

## Correction Budget

Contract allowed `max_correction_attempts: 2`. **Six rounds were used.** Rounds 3-5 were test-only
and explicitly user-authorized in the visible Milestone session. Round 6 was hub-directed under the
dual-review fold rule and is likewise test/evidence-only. Recorded in the contract as a deliberate,
disclosed overrun — not a silent one. Process note for Mission Strategist, not M-04 rework.

## Lessons Worth Propagating (mission-wide)

1. **Two half-proofs do not satisfy a criterion.** Handler→context and context→cache both passing
   proved nothing about handler→cache. Read the contract's Evidence line *literally* and check that
   one test does that thing end to end.
2. **A guard that only covers the negative path is hollow.** Assert the success path directly and
   exactly (`== ["catalog"]`, not "non-empty").
3. **Require a mutation proof for any test cited as evidence for a criterion.** Six hollow/weak
   tests were found in this milestone; a green suite never once caught them. Only mutation did.
4. **Run every half of a multi-part gate before declaring terminal.** The single highest-value check
   in this milestone came from the reviewer I skipped.

## Accepted Limitations

- **`-race` was never run.** `CGO_ENABLED=0`, no C compiler (gcc/clang/cl/mingw all absent).
  Independently confirmed three times; genuine environment limit, not evasion. QA assessed it blocks
  **no** criterion. Still: this is concurrent code — **run under `-race` in a cgo-capable CI before
  production rollout.**
- Cache is per-process, no invalidation bus (mission Non-Scope). Multi-instance deploys serve
  independently-aged entries; each self-heals on TTL expiry.
- `gofmt -l` flags files repo-wide including untouched ones (pre-existing Windows CRLF condition).
  `git status` shows ' M' on three production files for the same reason; content diffs are empty.
  Not a regression, not a finding.

## Non-Blocking Residual Risks (carry forward as mission notes)

1. `TestFreshnessCacheSingleflight` is a weak regression detector — no barrier ensuring the other 19
   waiters parked. QA measured it catching singleflight deletion **0/20**; `ErrorNotCached` catches
   the same removal **20/20** and is C02's real guard. Future-regression sensitivity, not a current
   defect. Follow-up: add a readiness barrier.
2. `cache_test.go:585` is dead weight — insensitive, superseded by the composed C03 test, now
   correctly labeled as such in the evidence.
3. `-race` in cgo-capable CI before rollout (above).

## Carry-Forward For Integration (load-bearing, easy to break)

- **MaxAge=0 bypass works only because** this milestone moved the freshness context key from the
  catalog transport package into `internal_read/domain`, so both sides share
  `domain.WithFreshnessPolicy`. Previously transport set a *private* key the cache could never read.
  **Do not move it back** — and note the composed test now guards exactly this.
- **Linkage exclusion is structural, not configuration:** `cache.Reader` embeds the port and
  overrides only the two catalog page methods. No linkage TTL exists to misconfigure.
- **Concurrency rests on two layers that must both stay:** the generation in the singleflight group
  key, and the generation fence in `storeIfCurrent`. Removing either reintroduces a C05 stale-read.
- No port signature changes; `contracts/` and `packages/` untouched → no OpenAPI/sdk-runtime lockstep
  required. `go.mod`: single change, `golang.org/x/sync` indirect → direct. No Redis/external cache.

## Commit Trail

| SHA | Round | Scope |
| --- | --- | --- |
| `05c2c012` | F-01 build | feat: cache decorator, singleflight, evict-on-mutation |
| `0e1c2469` | correction 1 | fix: generation fence, bypass namespace, deep clone, DoChan waits |
| `06aeaa0a` | correction 2 | fix: generation in group key, snapshot ordering (**last production change**) |
| `356d0913` | correction 3 | test-only: deterministic post-invalidation regression test |
| `23abd388` | correction 4 | test-only: mutation-sensitive linkage exclusion; delete tautology |
| `df2cac6a` | correction 5 | test-only: nil-vs-absent proof; exact log-attribute allowlist |
| `56d5be9e` | correction 6 | test/evidence/contract-only: composed C03 + C05, success-path guard |

No production code has changed since `06aeaa0a`.

## Operational Note For The Hub

`codex-companion` is **broken on this machine** — its sandbox helper fails to launch in a worktree
(`[windows] sandbox = "elevated"` in `~/.codex/config.toml`) and it hardcodes the sandbox mode with
no override, so jobs zombie with zero output while still reporting "running". The direct call works
and is what ran both the implementer and the Sol review half here:

```
codex exec -m <model> -c model_reasoning_effort="high" -s <read-only|danger-full-access> \
  -c approval_policy="never" -C <worktree> -o <outfile> - < <packetfile>
```

Write the packet with a file write, **not** a bash heredoc (heredocs parse-error and silently skip
the run). `gpt-5.6-sol` is reachable this way — the hub does **not** need to run the Sol half for me.

## Next

Hub owns the main merge. Branch `claude/focused-borg-9c9811`, reviewed SHA `56d5be9e`, ready for
integration. Seams stayed disjoint from M-05 (server Go only; no `apps/web` files touched).
