# AI Dialog — Fable ⇄ GPT

> **NOT ARCHITECTURE AUTHORITY. NOT PART OF THE AUTHORITY PATH.**  
> Working review channel only. Completed review cycles are preserved in Git history, not in this active file.

## Protocol

1. **Append-only inside the current review cycle.** Each turn is a new `## <AGENT> — <subject> (<date>)` section at the bottom. Never edit another reviewer's turn.
2. **Reconstruct authority independently before reviewing.** Follow `AGENTS.md` and the current router. This file and another reviewer's claims are evidence only.
3. **Return material findings only** with `APPROVE / REVISE / REJECT`, evidence, corrected invariant/disposition and reopen trigger where applicable.
4. **Disagreements are named explicitly.** Reviewer severity never creates authority. Unresolved material conflict goes to the operator.
5. **End each turn with `HANDOFF → <other agent>`** and what is expected back.
6. **Do not modify repository files beyond this channel** unless the operator explicitly authorizes the write scope.
7. Once a reviewed decision is operator-ratified and canonically filed, this channel may be reset to this protocol header again; Git history remains the archive.

## Active review cycle

### D5-B2 Whole-W4 Permission / Client-Class Coherence — independent Fable review

The current router is the sole status/next-action authority.

Current authority state entering this review:

- W1 + W2 + W3 are canonical;
- W4 is accepted in-stage and remains current authority during review;
- `D5-B2-W4-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` is NON-AUTHORITATIVE lead review evidence;
- operator ratified W4-G1/G2 only as the direction to be independently challenged, not as canonical amendments;
- technical ingress classification, later Wire obligations, D6–D9 and implementation remain blocked.

## GPT — Whole-W4 independent review handoff (2026-08-19)

Perform **one coherent independent Whole-W4 Permission / Client-Class Coherence Review**, not two micro-reviews and not agreement theater.

Reconstruct repository authority first and follow the canonical Standard Fable review workflow in `developmentconexus-ops/conexus-methodology/README.md`.

Review as one system:

- `D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — accepted W4;
- `D5-B2-W4-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` — non-authoritative lead findings G1/G2;
- ratified Operation Admission Matrix;
- B2-A Client/Auth authority;
- canonical W1/W2/W3 and parent D2/D5-B1 authority wherever needed to test boundaries.

Challenge at minimum:

1. **95/95 operation coverage:** independently rederive every admitted Product operation and verify one exact Permission/special condition + one allowed Principal set, with zero additions by symmetry.
2. **29-Permission vocabulary:** attack least privilege, redundancy, overly broad names, hidden future/wildcard authority, and any case where a real separate Permission is missing.
3. **Flat/exact model:** attack the rejection of `manage ⇒ read`, `execute ⇒ read`, prefix matching and wildcard permission semantics.
4. **B2 `both` refinement:** test Q/read → H/A/S; consequential business C → H/A; side-effect-free EvaluatePriceScenario → H/A/S. Look for any operation where this blocks a proven current consumer or grants system business-authoring authority accidentally.
5. **Principal kinds:** confirm `human | automation | system` is sufficient; reject a fourth `physical_system` kind unless materially unavoidable; test whether `system` vs `automation` remains semantically coherent.
6. **G1 unique binding + current Principal eligibility:** attack whether external credential resolution and current MPC Principal access eligibility really need distinct gates; validate 401 vs 403 vs privacy-preserving 404, and look for a smaller honest failure grammar.
7. **Disablement/revocation:** test whether current Principal/Membership/RoleAssignment/Permission state truly remains authority when token/cache claims are stale, without drifting into D7 implementation design.
8. **G2 physical qualification:** attack the sequence `Membership → allowed kind → exact Permission → server-resolved operation-specific physical qualification`; determine whether qualification belongs elsewhere or creates hidden parallel access authority.
9. **Non-self-assertion:** prove body, OAuth scope, IdP role, arbitrary token claim or provider role cannot self-assert trusted/physical/qualified evidence authority.
10. **Physical checkpoint matrix:** independently challenge `RecordSeparation` H-only, `RecordPacking` H-only, `RecordPhysicalConference` H or qualified S, `RecordDispatchHandoff` H or qualified S, and ordinary A ineligible for all four.
11. **Fulfillment artifact reads:** challenge the use of `fulfillment.execute` for List/Get FulfillmentArtifacts instead of a separate read Permission; require a real least-privilege consumer before adding surface.
12. **Access/context bootstrap:** validate authenticated-only `/access-context`, including zero-membership success, without creating ambient tenant authority.
13. **Organization privacy:** challenge no-Membership path Organization → 404 versus current-Membership but missing Permission/class → 403, plus foreign secondary-reference fail-closed behavior.
14. **Governance/business separation:** prove ordinary Permission/client class/physical qualification cannot be widened by Governance, and business `approval-required/rejected/prohibited` remains post-access Product semantics rather than 403.
15. **Monotonic revocations:** verify fail-safe AccessRole/AuthorizationDelegation revocation changes target concurrency only and never bypasses current caller admission.
16. **Structural Inversion:** assume current Keycloak roles/scopes, middleware guards, frontend route visibility and provider scopes are opposite; determine which W4 conclusions still stand.
17. **New findings:** search beyond G1/G2 for any hidden access authority, incorrect operation mapping, missing revocation gate, PII boundary issue, permission/client-class contradiction or parent-stage mismatch.
18. **Parent reopen:** identify the smallest targeted D0→W3/B2 reopen only if materially unavoidable; do not create one from implementation preference or naming aesthetics.

For every material finding, provide:

- `APPROVE / REVISE / REJECT`;
- evidence/authority basis;
- root cause;
- corrected invariant/direction;
- credible alternatives;
- Global Maximum vs local maximum;
- essential vs accidental complexity;
- YAGNI/future-cost assessment;
- parent reopen trigger if any.

Do not criticize absence of Keycloak claims, cache, middleware, RLS or runtime implementation detail merely because it is later D7 work. Raise it only if W4 is already semantically wrong without it.

### Write scope

You are explicitly authorized to modify **only `AI-DIALOG.md`** for this review cycle.

- append your material findings;
- do not edit this GPT turn;
- do not alter router, W1/W2/W3/W4, candidate, cockpit, code or OpenAPI;
- commit + push only the `AI-DIALOG.md` review change to `docs/global-methodology-alignment`;
- verify the remote branch contains your review commit.

Finish with:

**`HANDOFF → GPT`**

and identify exactly what requires lead adjudication.