# ADR-016: `sdk-runtime` is a manual client — OpenAPI and SDK change in the same commit

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-004's contract discipline and was checked
in review across every feature that touched the API surface, but no document was ever
written under this number. It is reconstructed here from the 12 live citations of
`ADR-12` that assert this meaning, harvested at
`docs/architecture/decisions/_citations/adr-012-citations.md` (Assertion A3). The same
`ADR-12` label was also used, unrelated, for two other decisions in two other missions —
see the registry at `docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`.

## Context

`packages/sdk-runtime` is a hand-written TypeScript client for the API — there is no code
generator. That means nothing mechanically keeps it in sync with `OpenAPI` when the
handler shape changes: a developer can edit the OpenAPI spec, or the handler, without
touching the SDK, and the drift is invisible until a frontend caller hits a field that
does not exist, or a type that no longer matches. The risk named in planning was literal:
"regenerar" fantasma (there is nothing to regenerate) and collision inside the
hand-maintained client when two features touch it in parallel.

## Decision

**Every commit that changes the OpenAPI contract must land the corresponding
`sdk-runtime` type changes in that same commit.** There is no separate "regenerate the
SDK" step, because there is no generator — the SDK is maintained by hand, one client file
per milestone, folded into a hub-adjudicated barrel.

### The clauses

**§1 — OpenAPI and SDK changes land in the same commit.** This is the rule as ratified
in the mission record.
> `.mnfs/MIS-004-mvp-demo/mission.md:96` — "ADR-12 sdk-runtime manual: OpenAPI+SDK mesmo
> commit; arquivo de client por milestone + barrel hub-adjudicado | ratificada"

**§2 — The SDK has no generator; "regenerate" is not a step that exists.** The risk this
rule closes is named explicitly as the "regenerar" fantasma — there is no tool to run
that would resync a drifted SDK, so drift can only be prevented at commit time, not
repaired mechanically afterward.
> `.mnfs/MIS-004-mvp-demo/mission.md:96` — "\"regenerar\" fantasma; colisão no client" (the
> named risk this decision closes)
> `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/PLAN.md:98` — "No generated output:
> `packages/sdk-runtime/package.json:7-10` has only manual TypeScript build/test scripts
> and ADR-12 explicitly declares the SDK manual"

**§3 — Parity between OpenAPI, SDK, and handler is checked as a smoke-level gate.** The
mission's validation contract runs an explicit OpenAPI↔SDK↔handler parity check as part
of its L2 smoke lane, on main, post-merge of every wave.
> `.mnfs/MIS-004-mvp-demo/validation-contract.md:141` — "rodar lanes L0–L2 (profile §5) no
> main pós-merge de todas as waves, `GOCACHE=.gocache`; incluir parity check
> OpenAPI↔SDK↔handler (smoke L2, ADR-12)"

**§4 — Feature validation records the same-commit discipline per surface changed.** Every
feature that added or changed an API surface in MIS-004 recorded its OpenAPI+SDK change
as a same-commit pair in its own validation evidence.
> `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-01-identity-semantics-fix/validation.md:15`
> — "OpenAPI catalog schemas + manual SDK parity (ADR-12 same commit)"
> `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-02-erp-import-module/validation.md:17`
> — "OpenAPI /erp/imports* paths + ErpImport* schemas + sdk-runtime types (ADR-12)"
> `.mnfs/MIS-004-mvp-demo/M-04-vinculos-import-ui/evidence/d02-investigator-gaps.md:42` —
> "Update **OpenAPI spec in the same commit** (ADR-12 discipline; confirmed by a2c0698
> commit message: \"OpenAPI spec and the manual sdk-runtime types land in the same
> commit\")."
> `.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/evidence/p2-batch-plan.md:422` — "SDK
> contract file | `packages/sdk-runtime/src/market.ts` | Manual client/types must land
> with OpenAPI under ADR-12 | F-04-S7, same commit as OpenAPI"

**§5 — An owner can be granted an additive-only exception to this rule.** F-01 of M-01
was granted an additive-lock exception on `sdk-runtime/src/index.ts` catalog types: it
may add fields for identity, under the same-commit rule, but the grant is scoped to
additive changes only, not free rewrite rights over that file.
> `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md:55`
> — "tipos catalog em `packages/sdk-runtime/src/index.ts` (**additive-lock grant**
> registrado na matriz da missão — aditivo apenas, ADR-12 mesmo commit)."

## Rationale

Same-commit atomicity is the cheapest available substitute for a generator: if the SDK
cannot be produced mechanically, the next best guarantee is that the two files are never
in a state where one has moved and the other has not, because they are never committed
separately. This is a discipline about *when* things change together, not about
*whether* they agree once changed.

## Consequences

- **This rule enforces atomicity, not agreement.** `GOV_API_SDK_SPLIT` (the same-commit
  check named in this rule) only verifies that the OpenAPI spec and the SDK types changed
  in the same commit — it does not verify that the resulting SDK types actually match the
  shape the OpenAPI spec now declares. A commit that changes both files, but changes them
  to disagree with each other, satisfies this rule. Closing that gap — checking
  *agreement*, not just *same commit* — is a separate, known piece of work being addressed
  outside this decision.
- Because there is no generator, a hand-maintained SDK is permanently exposed to human
  transcription error between OpenAPI and TypeScript; the smoke-level parity check (§3)
  is the mission's mitigation for that risk, not a guarantee against it.
- Any feature touching an API surface must budget the SDK edit as part of the same
  change, not as follow-up work — a same-commit rule that is violated blocks merge under
  this discipline.
- Exceptions to file ownership (like the additive-lock grant in §5) must be recorded in
  the mission's ownership matrix; they do not suspend the same-commit rule itself, only
  narrow who may make additive edits under it.

## Alternatives Considered

**Generate the SDK from OpenAPI.** Not adopted for this mission: `sdk-runtime` has no
generator and none was built; the manual client with same-commit discipline was the
mechanism actually in place.

**Allow SDK updates as follow-up commits with a tracking issue.** Rejected: this is
exactly the "regenerar fantasma" risk the rule was written against — a follow-up that
does not happen leaves the SDK silently stale with no mechanical signal that it drifted.

**Check OpenAPI/SDK agreement instead of, or in addition to, same-commit atomicity.** Not
built in this mission — the smoke-level check that exists (§3) verifies parity as part of
L2, but the harvested citations only assert the same-commit gate (`GOV_API_SDK_SPLIT`) as
ratified; a dedicated agreement check is recorded above as a known, separately-addressed
gap rather than as part of this decision.
