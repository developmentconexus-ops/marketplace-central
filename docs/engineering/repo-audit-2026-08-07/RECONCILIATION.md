# Phase 2 — Reconciliation of the two synthesis arms

Inputs: `synthesis-arm-a.md` (9 axes), `synthesis-arm-b.md` (8 axes). Neither arm saw the other.
Both received the same brief and the same three mid-flight amendments.

Playbook: *convergence between independent arms is weak evidence of correctness — divergence is
strong evidence of an unsettled question.* This file records the divergences and resolves them.

---

## Where they converged

Convergence here is worth little on its own, but its **shape** is informative: the arms cut the same
material nine and eight ways and produced near-identical totals from independent groupings.

- **Total effort:** A ~33–46 agent-days, B ~33–48. Independently derived.
- **First axis:** both put verification first, and both justify it as the multiplier that makes every
  later axis cheaper. Neither was told to.
- **Root causes:** both name the same underlying set — the verifier is off the change path; the
  contract is hand-copied with no arbiter; identity is a boot constant; the boundary instruments are
  anchored to the tree being abandoned; money has no compulsory representation; failures leave no
  trace; evidence certifies itself.
- **Both independently refuted `cicd` F2** on the remote, matching the measurement taken in-session.
- **Both placed the money axis late**, and for the same stated reason: it is ~78 files of conversion
  work that must not run while 91% of the test suite is unreachable.

---

## Divergence 1 — the public flip versus the open door *(the one that matters)*

**Arm A** schedules an edge control on the unauthenticated PII endpoint at position **0**, measured in
hours, and calls it a **precondition of making the repository public**.
**Arm B** does not, and sequences the flip immediately after its verification axis.

**Resolved in favour of A, and it is not a close call.**

The repository documents, in files that go public with it:

- `deploy/Caddyfile` — the exact path predicate that routes `/orders` with a non-HTML `Accept` header
  to `backend:8080`
- `docker-compose.yml` profile `oauth` — the ngrok tunnel that publishes the frontend
- `apps/web/vite.config.ts` — the full proxy table
- `orders/transport/http_handler.go:608-618` — the handler that serialises buyer name, CPF/CNPJ and
  address
- `composition/root.go:994` — proof that the chain is `CORSMiddleware(apierror.Recover(mux))` and
  contains no identity check

Publishing that set is publishing a map to an open door. The flip does not create the exposure — the
exposure is already live — but it removes the only thing currently limiting who knows the route.
B's sequence is internally coherent and would be correct for a repository that stays private.

**Binding consequence:** the edge control lands before the flip, and the pre-public sweep's verdict
is necessary but no longer sufficient on its own.

---

## Divergence 2 — granularity of the verification axis

**B** merges it into one axis: *the verifier is not on the path a change takes.*
**A** splits it into three: a delivery path nothing is required to traverse (A1); a verifier with no
single entry point (A2); a green result that does not prove the check ran (A3).

**Resolved in favour of A.** The playbook's test is whether two axes want the same finding. They do
not — the three have disjoint findings, disjoint fixes, and different failure modes:

| Cause | Fix | Fails as |
|---|---|---|
| nothing is *required* to pass | ruleset + required check | a bad change merges unexamined |
| no single entry point | one command that runs everything | a good check is never invoked |
| green does not prove execution | fixtures, `failure_token` | a check reports success without running |

Fixing any one of these leaves the other two fully live. B's merge would have produced one issue
whose "done" is unfalsifiable.

---

## Divergence 3 — where "the guard and its proof are written together" belongs

**B** gives it its own axis (G). **A** folds it into A3, *a green result does not prove the check ran*.

**Resolved in favour of A — they are one cause.** Both are *evidence that certifies itself*: a
detector with no fixture proving it can fail, an all-skipped run byte-identical to a green one, a
negative fixture asserting an error code the engine never implemented, and a guard whose only test
was authored in the same commit as the guard (catalog RLS at `0098`/`7e3dcc47`; the tenant guard at
`47a76837`). One axis, named for the cause.

---

## Divergence 4 — browser failure containment

**B** gives it its own axis (H). **A** folds it into A5, *every request boundary is re-invented per
module*.

**Resolved in favour of B.** A React error boundary is not a request boundary; the conflation is
verbal. Distinct cause, distinct fix, cheap and independently schedulable.

---

## Divergence 5 — the routing tables

Both arms received the finding that `deploy/Caddyfile` and `apps/web/vite.config.ts` maintain
duplicate route tables by hand, with a comment instructing a human to keep them in sync, and the Go
router as a third source of truth.

**A** counts it as part of the contract axis — six hand copies, four for shape and two for routing.
**B** treats it as a verification concern.

**Resolved in favour of A.** The cause is identical to the four-copy shape seam: several transcriptions
of one truth with nothing comparing them. Six copies, one axis.

---

## Refutations of Phase 1 lanes — adjudicated

| Refutation | Ruling |
|---|---|
| **A-D1** — `delivery` F14 claims 10 OpenAPI operations unreachable from the SDK; 4 are reachable (`market.ts:179,182,185,190` call `postJson`/`getJson`) and 3 more are OAuth browser redirects where an SDK method would be a defect. **Real gap: 3.** | **Upheld.** The refuting fact sat in `lanes/frontend.md:39` the whole time. Cross-lane contradiction neither lane checked. |
| **A-D2** — `frontend` F3 calls the missing `onRetry` a crash that downs the SPA. React ignores an undefined handler. | **Upheld, verified in-session.** `packages/ui/src/ErrorState.tsx` passes `onClick={onRetry}` with no invocation guard; undefined yields an inert button, not a throw. Severity re-rated: a permanently dead retry control. Its value as live proof that `tsc` never runs is undiminished. **This error was also repeated verbally to the operator and has been corrected.** |
| **B-2** — `layering` L-03 claims `internal/composition` is unreachable by every instrument. `scripts/arch-gate.sh:30` does scan it, and `cicd` F7 measured 42 findings from that root. | **Upheld.** Verdict survives only under *a control that does not fire is absent*. Remediation changes from **build a detector** to **wire a script** — materially cheaper. |
| **B-3** — `duplication` DUP-6 proposes replacing 82 hand-written `rows.Next()` loops with `RowToStructByName`. B says keep them forever. | **Upheld, and promoted.** The swap trades a compile-time arity and type check for name-based runtime reflection, in a codebase whose named failure mode is *the wrong noun used correctly everywhere*. Rejecting it is consistent with the program's own goal. Recorded on the do-not-touch list. |
| **B-4** — `testing` T-2 (the ADR-023 detector has no fixtures) versus `layering`'s "actually fine" verdict on the same detector. | **T-2 wins.** No `testdata/` exists under `internal/composition`. |
| **B-new** — the `migration-prefix-0021-duplicate` governance exception (2026-07-11) predates by three weeks the second collision it does not cover (2026-08-03). | **Accepted.** No lane found it. Promotes ledger decay from one stale entry to a full axis. |
| **A-self** — A withdrew its own proposed `127.0.0.1` mitigation after determining it would read as done and be inert. | **Noted as method working as designed.** Same failure class as the findings it was meant to fix. |

---

## Merged axis set

Ten axes. Every Phase 1 finding belongs to exactly one.

| # | Axis (named for cause) | Source |
|---|---|---|
| **V1** | Nothing is *required* to pass through the delivery path | A1 |
| **V2** | The verifier is scripts to remember, not one product with one entry point | A2 ≡ B-A |
| **V3** | Evidence certifies itself | A3 + B-G |
| **C1** | The published contract is six hand copies with no arbiter | A4 + B-B |
| **I1** | Identity is a boot-time constant, so every downstream enforcement enforces a constant | A8 ≡ B-C |
| **R1** | Every request boundary is re-invented per module instead of provided once | A5 ≡ B-partial |
| **M1** | There is no single compiled vocabulary for a value | A6 ≡ B-F |
| **B1** | The boundary instruments are prefix-anchored to the tree being abandoned | A7 ≡ B-D |
| **O1** | A failure leaves no durable trace | A9 ≡ B-E |
| **F1** | The browser has no failure containment | B-H |

---

## Sequence

| # | Step | Size | Why here |
|---|---|---|---|
| **0** | **I1-edge** — close the unauthenticated PII route at the edge (Caddy `@orders_api`, ngrok scope) | hours | **Precondition of the public flip.** Divergence 1. |
| 1 | **V2** — one entry point; wire what already exists | 4–6d | The multiplier. Built green locally, before anything is required. |
| 2 | **V3** — fixtures, `failure_token`, no self-certifying evidence | 3–4d | Must precede V1, or the required check certifies itself. |
| 3 | **V1** — flip public, push the 98, enforce the ruleset | 1d | Gates land before the push, per operator sequencing. |
| 4 | **I1-structural** — authentication and per-request tenant | 5–8d + auth unsized | Auth itself is a product decision, not an audit output. |
| 5 | **C1** — one arbiter for the six copies | 4–6d | `oapi-codegen` + `openapi-typescript` already approved. |
| 6 | **R1** — provide the request boundary once | 5–7d | |
| 7 | **M1** — one compiled money vocabulary | 7–10d | ~78 files. Must not run while the suite is dark. |
| 8 | **B1** — re-anchor the instruments to both trees | 4–5d | |
| 9 | **O1** — durable failure traces | 4–6d | |
| 10 | **F1** — browser failure containment | 1–2d | Independent; schedulable any time after V2. |

**~34–49 agent-days.** Both arms landed within three days of this independently.

---

## Open, and not resolvable by synthesis

- **Authentication is a product decision.** Both arms declined to size it. Who are the principals,
  what is a session, is there more than one human user ever? Until answered, I1-structural is a
  range, not an estimate.
- **Is a production host currently running?** `MPC_DOMAIN` is external. Changes I1-edge from urgent
  to immediate.
- **D-2 / D-50** — is `removal_owner=HARNESS-D-N` ratified practice, or do the four existing uses
  need real owners? Feeds B1.
