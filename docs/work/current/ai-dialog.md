# AI Dialog

Candidate: `developmentconexus-ops/marketplace-central` PR #65 @ `b254e05fe33ccb50efc161f064c17152490c42c1`
Round: R1
Methodology: `developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535`
Transport: temporary review-branch Evidence only; never merge to candidate/main.

## Findings

| ID | Severity | Claim attacked | Status |
|---|---|---|---|
| R1-F1 | IMPORTANT | verify-ci-policy proves lane policy | OPEN |
| R1-F2 | IMPORTANT | stable `required` check cannot be satisfied by a quick pass on a Ready candidate | OPEN |
| R1-F3 | MINOR | quick matcher selects "directly affected" verifiers | OPEN |
| R1-F4 | MINOR | quick-map entries are shown to fire | OPEN |
| R1-F5 | MINOR | pr-title `synchronize` removal is loss-free | CONFIRMED SAFE (recorded basis) |

No MATERIAL finding. Central decision (proportional lanes, full aggregate proof at Ready/main, stable single `required` context) survives the strongest attacks I could construct on every normal event path.

## Dialogue

### R1-H0 — ROUTING

Perform independent adversarial review of the exact candidate. Revalidate repository authority and candidate identity first. Attack whether proportional Draft quick verification can create a blind spot, whether every integration candidate/main path still receives full aggregate proof, whether change-aware verifier selection makes only claims it actually proves, whether the stable required check remains trustworthy, and whether the policy guard adds less ceremony than it removes. Do not modify the candidate branch.

### R1-C1 — CHALLENGER

Independent reconstruction performed from remote authority. Identity revalidated: PR #65 head `b254e05fe33ccb50efc161f064c17152490c42c1` (base `main` @ `689dab34`), review branch bootstrap `9b9e4567`, methodology pin `9c7210d1504bef01c0d134a6c3ae8627deebb535` (ROUTER → METHOD v1.1.0 + ADVERSARIAL-REVIEW-METHOD v1.0.0). Candidate branch not modified.

Evidence base examined: exact diff `origin/main...b254e05f`; full current `scripts/gate.ps1`, `ci.yml`, `pr-title.yml`, `package.json`, `scripts/verify-ci-policy.mjs`; `origin/main` gate lane semantics; main ruleset via API; live CI runs for #65 and #66; local execution of `verify-ci-policy.mjs`.

Confirmed load-bearing facts (attack-surface groundwork):

- Ruleset 20614327 on `main`: required status check context = `required` only; `conventional-title` is advisory; `pull_request` rule blocks direct pushes; linear history + non-fast-forward + deletion protection on. The stable `required` name IS enforced by the ruleset — not merely named "required".
- On `origin/main`, `-Lane` was label-only (used solely in a `Write-Host`) and `gate`/`gate:full` ran identical verifier chains. This PR makes the lane distinction real for the first time; the prior "quick" was a fiction. The root cause (accidental ceremony: every Draft push paid the full aggregate) is genuinely addressed, not patched.
- Go is consumed only by `verify-product-oad*.mjs` (they spawn `go`); the quick lane structurally excludes the whole `verify-product-oad` prefix and skips the product proof, so skipping `setup-go` on Draft hides no currently relevant failure. Quick-mapped verifiers spawn only node/npx (`verify-oad-source-reachability` runs redocly via npx).
- Event coverage: `opened|synchronize|reopened` while Draft → quick; `ready_for_review` → full on current head; `synchronize` after Ready → full; `converted_to_draft` → quick; `push` to main → full. `github.event.pull_request.draft` is null on push, so both `== true` and `== false` are false and the `push` disjunct covers it. Both lane steps live under the single `required` job, so exactly one gate step runs per event and the context stays stable.
- Live proof of both lanes: #65 run 32774632903 (`full gate`: `gate lane: full`, `product_proof: FULL`, `gate: PASS`, ~1m56s) and #66 run 32775056328 (`quick gate`: `gate lane: quick`, `review_mode: True`, `quick_verifiers: NONE`, `product_proof: SKIPPED_QUICK`, `gate: PASS`, ~18s).
- Review isolation: review-mode logic untouched; on a review branch the changed-file set is forced to `docs/work/current/ai-dialog.md` only, which maps to no quick verifier; candidate-tree contamination and diff-name checks still fire; negative-control census unchanged 6/6. `.gitattributes` forces `eol=lf` on yml/json/ps1, so the byte-level string asserts in `verify-ci-policy.mjs` are CRLF-safe on Windows checkouts (verified by local run: `ci_policy=PASS`).
- `gate:full` retains the entire prior chain (hygiene + all 14 verifiers + product proof) plus the new policy guard; nothing was weakened on the integration path.

**R1-F1 — IMPORTANT — proof weaker than claim in `verify-ci-policy.mjs`.**
Claim attacked: the guard's asserts claim lane-to-condition association ("CI must select full verification for non-Draft candidates", "Go setup must be reserved for full verification"). Counterexample: swap the two step `if:` expressions in `ci.yml` (Draft → `npm run gate:full`, Ready → `npm run gate`), or move the `if:` off `setup-go` to any other step — every assert still passes, because each assert only checks that tokens exist somewhere in the file, never that a condition is attached to a given step. The guard is a presence tripwire, not a policy proof; under the Structural Inversion Test the central protected property (Ready/main always full) is not held by this control at all — it is held only by unguarded ci.yml text. Why it matters: this is the exact instrument the increment adds to defend the lane policy against future drift, and today it cannot detect the most dangerous drift (silent lane inversion). Smallest scope: `scripts/verify-ci-policy.mjs` only. Recommended disposition: assert the exact contiguous step blocks (name + `if:` + `run:` as one multi-line includes each, LF-safe), which is deterministic, adds no machinery, and makes the swap counterexample fail; alternatively weaken the assert messages to what is actually proven. A deterministic negative path for this guard would also let it join the falsifier discipline the repository requires of material guards.

**R1-F2 — IMPORTANT — stale quick success can satisfy `required` on a Ready candidate.**
Claim attacked: "every Ready candidate path really receives full aggregate verification" as enforced, not merely as normally triggered. The ruleset requires context `required`; GitHub satisfies it with the latest check run of that name on the head SHA, and quick and full deliberately share the name. Two reachable windows: (a) an operator re-runs the Draft-era workflow run after the PR is Ready — re-runs replay the original event payload (`draft: true`), so the quick lane executes again and its fresh success becomes the newest `required` conclusion on that SHA, unblocking merge without full proof having the last word; (b) if the `ready_for_review` trigger is lost (Actions outage), the existing Draft quick success on the same head SHA already satisfies `required`. No normal path is affected — every organic flip/push produces a full run — and both windows are operator-visible (`product_proof: SKIPPED_QUICK` in the log, ~18s runtime), but "visible in a log nobody must read" is weaker than the enforcement claim. Smallest scope: quick-gate step in `ci.yml`. Recommended disposition: either record this explicitly as an accepted residual risk with a reopen trigger (cheapest, honest), or close it with a few lines — the quick step re-checks the live PR draft state via the API and fails when the PR is no longer Draft, which kills both windows deterministically. Not MATERIAL: exploitation requires a deliberate abnormal operator action or infra failure plus an unread log, and Draft PRs are themselves unmergeable; but accept-vs-harden is the Lead's decision to make on the record.

**R1-F3 — MINOR — "directly affected" is under-approximated and the map has no completeness discipline.**
A Draft touching `contracts/api/product/` runs only `verify-oad-source-reachability`, while at least seven full-lane verifiers consume the OAD (authorization-request oad/fable-fixes/d7r/d8r/w1-carrier, notification-oad, operational-read-contract); doc-driven verifiers (p9-screen-contracts, performance-evidence-knowledge) have no mapping at all. The roadmap wording is careful ("does not claim whole-Product non-regression") and full at Ready bounds every miss, so this is proportionality erosion, not a correctness hole: the promised early signal is silently absent for the heaviest Draft iteration class (OAD editing), and nothing forces a map entry when a future material verifier lands — the map rots silently. Recommended disposition: either state in the map comment that under-selection is intentional and Ready is the completeness boundary, or add a cheap census assert (every `verify-*.mjs` named in `gate:full` is either quick-mapped or on an explicit exclusion list).

**R1-F4 — MINOR — quick-map controls cannot be shown to fire.**
`Add-QuickVerifier` silently drops any mapped verifier path that does not exist (`Test-Path`). This is deliberate forward wiring for the unmerged PR #64 `verify-d6-r-b10-preparation.mjs`, and as a seam it is acceptable; but it also means a typo'd map path is forever indistinguishable from "no verifier applicable" — the control never demonstrably fires. `quick_verifiers:` reporting mitigates (an expected verifier absent from the summary is observable), which keeps this MINOR. No correction required now; if the map grows, distinguish "mapped but missing" in the summary line.

**R1-F5 — MINOR — pr-title `synchronize` removal confirmed safe; basis recorded.**
The title can only change through `opened`/`edited`/`reopened` (all retained), and `conventional-title` is not a required ruleset context, so no merge can be blocked waiting on a check that never runs for a new SHA. Recorded because the safety depends on `conventional-title` staying advisory — if it is ever promoted to a required context, `synchronize` must return, or pushes after the last title event leave the new head SHA with no `conventional-title` run.

Smallest-sustainable-solution verdict: yes. The increment deletes real accidental ceremony (identical double lane, Go setup and 14 verifiers on every Draft push, pr-title on every push), adds one 31-line tripwire and one bounded gate.ps1 branch, and weakens nothing on the integration path. `verify-ci-policy` is a useful safeguard in intent; R1-F1 is precisely that its current strength does not yet match its claim.

Verdict: no MATERIAL finding. R1-F1 and R1-F2 require Lead adjudication (ACCEPT bounded correction, or REFINE/DEFER with recorded residual-risk basis). LEAD RESPONSE REQUIRED.
