# HUB EVENT — BLOCKED — M-01-listings-read-spine

```yaml
event: BLOCKED
from_chip: M-01-listings-read-spine
to: hub (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
date: 2026-07-15
review_sha: e2cde3648a6bdd534afc0ae076a08a93d06d7c7a
milestone_status: Blocked (pending C10) — STATUS NOT FLIPPED
```

## Summary

M-01 passes every code-verifiable criterion **C01–C09** across P5 ladder + P6 dual gate +
P7 milestone-validate (fresh execution evidence, cold 5-member crew, adjudicated C07 under
harness truth order). The **single** open item is **M01-C10 (live-provider-read, `Required: Yes`)**.
This session drove the C10 live lane against the real running platform and it is
`could-not-drive` for **environment** reasons, not code.

## C10 live-drive result (concrete, this session)

Real M-01 server brought up from the worktree (`:8080`, main `.env` loaded into process env,
values never printed), live `/listings` surface driven directly (single-config tenant):

- `GET /listings` · `GET /listings/summary` · `POST /listings/refresh {}` → `400 installation_required`.
- `POST /listings/refresh {"installation_id":"probe-nonexistent"}` → **`503 source_unavailable`**
  (not `404 installation_not_found`) → installation lookup fails at the DB layer.
- Server startup: `.env` `MC_DATABASE_URL` **not migrated** — `SQLSTATE 42P01` for
  `integration_provider_definitions`, `marketplace_definitions`, `marketplace_fee_schedules`;
  Oracle/assisted-linkage readers unavailable. → bare/partial local Postgres, **no connected
  ML installation**.

Endpoints, route-class, and the full error matrix behave per contract. **No code defect.**
Evidence (sanitized): `_gate-evidence/round-1/c10-live-attempt.md`; ledger `P7b`.

## Two operator-only prerequisites block C10 (chip cannot self-provision)

1. A migrated DB with a **real connected Mercado Livre installation** (valid OAuth tokens).
2. The **OAuth connect** itself (provider login/consent) — a credential-entry action the
   milestone session must never perform.

## Hub ruling requested — pick close path

- **(a) HOLD Blocked** until operator provisions the live env (migrated dev DB + connected ML
  installation via operator-run OAuth), then chip re-drives `POST /listings/refresh` and asserts
  run succeeded / tenant-scoped `count(*)>0` / real `MLB…` ids / `<20%` unknown status → re-gate → P8 close.
- **(b) RATIFY `Required: Yes`→deferred** amendment of C10 in `validation-contract.md` (route C10
  to a dedicated post-merge live lane), then M-01 closes on C01–C09. Requires operator/hub sign-off.

Doc reconciliation queued regardless (non-blocking): `validation-contract.md:127` C07
`below_margin`→`below_margin_worst_case` to match ratified D-22 + binding OpenAPI.

---

## Retro addendum (→ hub, process, applies M-02+): cut the plugin crew from close

Operator flagged review fatigue during M-01 close ("toda hora review depois mais review depois
outro review"). Diagnosis: **over-gating**, root-caused to the mission plan mapping "P7" to the
mnfs-workflow plugin `/milestone-validate` (structural precondition + qa-validator re-run + a
5-member ★1–★7 cold crew incl. 2 adversarial). Under harness truth order the BINDING
`docs/HARNESS.md` P7 is **fresh browser/live QA**, not that plugin crew — which largely
RE-DERIVES what the P4 per-slice reviews + P6 dual gate already established. Net on M-01: same
88-file diff cold-reviewed ~11× (slice) + 2× (P6) + 5× + qa-rerun (P7 crew) ≈ 18 passes before
the real harness P7 (live QA) even ran.

**Recommendation for M-02+ close:** `P4 (light, 1 reviewer/slice) → P5 verify ladder → P6 dual
gate (the ONE cold whole-diff layer) → P7 fresh browser/live QA (this is what passes a milestone,
incl. the live-provider read) → P8 close`. Do **not** stack `/milestone-validate`'s 5-member crew
as an extra layer; fold its single new value (a fresh integrated execution re-run) into P6/P7 as
ONE run. Full detail: `RETROSPECTIVE.md`.
