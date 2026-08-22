# D8 — Independent Fable Challenge

> Review branch only: `review/d8-fable`  
> Candidate branch: `stage/d8-golden-flows`  
> Candidate HEAD expected: `b3469258348289865a036bc7a946077f79d61faf`  
> Candidate Draft PR: #59 — `docs(d8): derive golden flows stage`  
> Base `main` expected: `8a9bcdb2c0333b94eecef8fb18cfa54b824b2881`  
> Product candidate: 99 Product operations / 30 ordinary Permissions / Principal kinds H-A-S only  
> D0–D7: ACCEPTED / CLOSED  
> D8: OPEN / ACTIVE — DERIVED CANDIDATE  
> D9: BLOCKED  
> Product implementation: BLOCKED UNTIL D9

## Purpose

Run an isolated adversarial challenge against the exact D8 candidate before D8 closeout.

The candidate deliberately selects **3 business golden flows + 1 systemic recovery falsifier** rather than an exhaustive 99-operation catalog:

```text
GF-01  Publication & Marketplace Convergence
GF-02  Sale → Business System → Fiscal → Fulfillment → Outcome
GF-03  Performance Evidence Honesty
SR-01  PITR / Timeline Continuity Recovery
```

Your job is not to maximize scenario count. Attempt to prove that this set is insufficient, redundant, internally contradictory, impossible to execute as later proof, or unable to falsify a materially different accepted invariant. Also challenge D8-F1: architecture closure versus implementation-readiness closure and the proposed post-D8 sequence `D6-R2 → Pre-D9 Implementation Readiness Contract → D9`.

Reviewer output is **evidence, not authority**. Do not edit `stage/d8-golden-flows` or PR #59. Write only below `## Fable response` in this file on `review/d8-fable`.

## Mandatory revalidation before analysis

Independently record:

1. remote `main` HEAD;
2. `stage/d8-golden-flows` HEAD;
3. PR #59 state/base/head/draft/mergeability;
4. changed files and count relative to `main`;
5. required CI on the exact candidate HEAD;
6. this review branch ancestry/tree relation to the candidate;
7. that `candidate...review` differs by **only** `docs/work/current/ai-dialog.md`.

If `stage/d8-golden-flows` differs from:

```text
b3469258348289865a036bc7a946077f79d61faf
```

**STOP and report `STALE_REVIEW_CANDIDATE`.** Do not review a moved target.

Expected candidate state at handoff:

```text
main                         8a9bcdb2c0333b94eecef8fb18cfa54b824b2881
candidate                    b3469258348289865a036bc7a946077f79d61faf
PR #59                       Draft / open / mergeable
changed files                3
ci                           SUCCESS
pr-title                     SUCCESS
active runtime               NONE
Product                      99 operations / 30 Permissions / H-A-S
```

## Strict reading discipline

Start exactly:

1. `AGENTS.md`
2. `docs/index.md`
3. `docs/roadmap.md`
4. `docs/engineering/rebaseline/D8-GOLDEN-FLOWS.md`

Then switch only to the exact prior owner/OAD needed for a concrete counterexample. Likely bounded references are:

- canonical `contracts/api/product/openapi.yaml` for exact Product surface;
- D4-R1 for publication authoring/joint realization;
- D4 for Mercado Livre/Sankhya external effect and fiscal/materialization evidence;
- D6-R1 for Performance semantics;
- D6 for frontend interaction/topology;
- D7-B for Organization/RLS/idempotency/concurrency;
- D7-C for durable work/ambiguity/reconciliation;
- D7-D for H/A/S carrier realization;
- D7-R1 for PITR continuity/recovery fence.

Do not recursively ingest D0–D7 history. Reopen an older authority only if a material falsifier actually implicates it.

## DevelopmentConexus Method challenge

Evaluate the candidate against the canonical Method, especially:

- smallest sustainable solution;
- proof strategy before implementation;
- control must be demonstrably falsifiable;
- YAGNI must not remove isolation/recovery/audit/knowledge correctness;
- Global Maximum does not mean maximum abstraction;
- accepted authority reopens only on material evidence;
- global coherence across accepted stages.

Do not request ceremony merely for completeness.

## Challenge A — selection minimality and completeness

Try to falsify the central rule:

> removing any selected row leaves a material accepted invariant without representative cross-boundary falsification; adding another business journey currently adds coverage rather than a distinct defect class.

Attempt to identify:

1. a material D0–D7 invariant class not represented by GF-01/GF-02/GF-03/SR-01 or their cross-cutting controls;
2. a selected flow that is actually redundant because another flow can falsify the same class without losing localization;
3. a hidden mega-flow where the candidate compresses materially distinct failure classes too aggressively;
4. a scenario included only because it is important operationally, not because it falsifies architecture composition;
5. an exhaustive-test impulse disguised as a golden-flow requirement.

A proposed new flow must name the exact defect class that cannot be falsified by the current set.

## Challenge B — GF-01 Publication & Marketplace Convergence

Attack whether GF-01 can really preserve:

```text
Readiness meaning
≠ Offering ListingIntent meaning
≠ Availability meaning
≠ provider protocol/effect
≠ provider actual Listing state
```

Attempt counterexamples involving:

- `FOLLOW_SOURCE` vs `EXPLICIT_OVERRIDE`;
- missing Availability input before provider create;
- shared User Product blast radius;
- PriceIntent remaining Offering-owned;
- stale ETag/revision;
- same idempotency key same vs changed semantic request;
- H vs A authoring and standing human override;
- timeout/crash after possible provider acceptance;
- provider 2xx before convergence;
- cross-Organization secondary reference;
- frontend inventing a screen-shaped write or raw provider shortcut.

Challenge whether the D4-deferred controlled real Mercado Livre write is necessary/sufficient D8 evidence and whether its proof can remain architecture/external-contract evidence without violating the Product implementation gate.

## Challenge C — GF-02 Sale → Business System → Fiscal → Fulfillment → Outcome

Attack the composition:

```text
Marketplace Sale
→ Party Resolution / Destination Realization
→ Business Order Intent / Sankhya TOP 313
→ Fulfillment physical readiness
→ Invoicing Intent / 313→306 fiscal result
→ provider artifacts / packing / dispatch
→ Shipment
→ bounded Post-Sale / Economics / Work consequences
```

Attempt counterexamples involving:

- multiple active Sankhya Party candidates;
- destination not safely representable;
- accidental Partner master-address overwrite or duplicate Party creation;
- Direct Oracle fallback;
- possible-acceptance timeout creating duplicate native order;
- invoicing before physical conference/readiness;
- A Principal recording a physical fact;
- unqualified S self-asserting physical authority;
- 313→306 2xx without authoritative result correlation;
- refund evidence falsely closing physical return or vice versa;
- Work close falsely resolving source truth;
- Sale detail becoming a cross-owner command/workflow authority.

Challenge whether Post-Sale as a branch of GF-02 is sufficient or whether it hides a materially distinct architecture defect class.

Challenge the D4-deferred controlled Destination Realization and first 313→306 fiscal proof: identify what must be proven before D8 can close and what must wait for the post-D9 Product runtime.

## Challenge D — GF-03 Performance Evidence Honesty

Attack whether a read-heavy flow is genuinely required and sufficient to falsify the Commercial Intelligence half of Product 1.0.

Attempt counterexamples involving:

- known zero vs unknown/unavailable/unsupported;
- partial coverage shown as complete;
- Listing survivorship bias;
- incompatible measurement basis producing a numeric delta;
- preserved historical evidence presented as live/current or MPC-authored;
- frontend-reconstructed CVR/ROAS;
- campaign/catalog/family scope collapsed into Listing;
- `performance.read` implying another Permission;
- cross-provider/all-marketplace aggregation without proven equivalence;
- Ads/optimization/AI write authority appearing by convenience.

Try to prove that GF-03 is redundant with Market/Economics reads or, conversely, that it omits a distinct Performance invariant that requires another representative path.

## Challenge E — SR-01 PITR / Timeline Continuity Recovery

Attack whether SR-01 is correctly systemic rather than another business golden flow.

Required counterexample:

```text
external effect survives
+ database restored to before acknowledged dispatch/idempotency/access state
+ ordinary boot with no manual recovery flag
```

Try to find a path where the candidate would:

- redispatch the already-possible effect;
- resurrect a revoked/expired human session or stale machine authority;
- trust absence of a restored marker as proof of no dispatch;
- release the fence without affirmative continuity/reconciliation;
- remain permanently fenced even when continuous lineage is positively established.

Do not demand implementation execution during D8 when D7-R1 explicitly assigns real Product-runtime conformance post-D9; do flag any D8 claim that cannot later be executed meaningfully.

## Challenge F — cross-cutting controls

Challenge whether the selected flows actually exercise, rather than merely mention:

- exact 99-operation / 30-Permission / H-A-S surface conservation;
- explicit Organization scope and secondary-reference isolation;
- H session+CSRF vs A/S bearer separation;
- exact Permission vs Principal kind vs physical qualification vs Governance;
- canonical OAD as the one Product wire authority;
- idempotency and opaque revision semantics;
- River/durable handoff and possible-acceptance reconciliation;
- Governance authorization without execution authority;
- Work responsibility without source-truth authority;
- honest unknown/partial/unavailable/provenance semantics;
- no Direct Oracle/provider/plugin/workflow/screen-shaped API escape hatch.

If a cross-cutting property requires its own flow, explain why it cannot be credibly falsified as a variant of the current set.

## Challenge G — proof horizons / implementation gate

D8 distinguishes:

```text
D8 architecture + bounded external-contract proof
!= post-D9 implemented Product-runtime conformance
```

Attack that boundary.

Determine whether any D4 statement genuinely requires a live external proof **before D8 close**, and if so, define the smallest safely authorized evidence required. Conversely, identify any D8 language that incorrectly claims current PostgreSQL/River/Keycloak/browser/router/object-store execution while active runtime is NONE.

Mocks/spikes may support local claims but cannot substitute for a real dependency where the eventual implementation claim depends on one.

## Challenge H — D8-F1 implementation-readiness finding

The candidate adds this pre-D9 sequence:

```text
D8 close
→ D6-R2 Complete Frontend Realization Closure
→ Pre-D9 Implementation Readiness Contract
→ D9
→ implementation only after accepted D9
```

Challenge whether this is:

- a material and necessary correction to implementation readiness;
- merely process ceremony already implied by D6/D9;
- incorrectly owned by D8;
- an accidental reopen of accepted D6 technology/topology;
- insufficient because implementation ambiguity also exists outside frontend;
- over-specified because it tries to pre-design private/local implementation mechanics.

The intended boundary is:

> pre-decide material behavior, authority, UX state, API/owner connection, dependency direction and acceptance outcome; leave local/private mechanics that do not alter those properties to implementation.

Challenge the proposed **Realization Contract** fields and import/dependency acceptance rule. Specifically test whether a coding agent could still make two materially different choices and both claim conformance.

Also challenge sequencing: if D6-R2 or the Pre-D9 readiness contract materially changes an accepted operation/owner/invariant, should the affected D8 flow be revalidated before D9 rather than silently relying on this candidate?

## YAGNI / explicit exclusions

Challenge the decision not to create separate golden flows for:

- channel onboarding;
- Availability config;
- Market/Economics analysis;
- Governance administration;
- access administration;
- Work lifecycle;
- Post-Sale standalone;
- every D6 interaction flow;
- every Product operation;
- every D7 proof obligation.

Recommend an addition only when a concrete distinct falsifier exists.

Also challenge whether D8-F1 creates speculative framework/process machinery rather than the smallest implementation-readiness seam.

## Product/non-regression fence

Independently verify the candidate has not changed:

```text
99 Product operations
30 ordinary Permissions
Principal kinds H / A / S only
stable origin https://conexus.fun
active Product runtime NONE
```

Flag any hidden Product operation, Permission, Principal kind, business authority, provider-shaped Product vocabulary, Direct Oracle path or Product implementation.

## Output contract

Append below `## Fable response` using this exact structure.

### 1. Verdict

Choose exactly one:

- `ACCEPT`
- `ACCEPT WITH BOUNDED FIXES`
- `REOPEN SMALLEST AUTHORITY`
- `REJECT / RECONSTRUCT`

### 2. Revalidation record

Record exact main SHA, candidate SHA, PR state/base/head, changed-file count, candidate CI status and review-isolation result.

### 3. Executive coherence assessment

Assess D8 selection quality, Method alignment, YAGNI/proportionality, global D0–D7 composition and readiness for GPT adjudication.

### 4. Material findings

Number findings highest severity first. For each include:

- **classification:** `D8_FIX`, `D0_D7_REOPEN`, `D6_R2_OBLIGATION`, `D9_ENTRY`, `IMPLEMENTATION_PROOF`, `REPOSITORY_FIX`, `LATER_NON_BLOCKING`, or `REVIEW_FALSE_POSITIVE`;
- **severity:** Critical / Important / Minor;
- exact candidate location;
- governing repository authority;
- current primary external evidence if externally dependent;
- concrete counterexample/failure;
- smallest correction;
- why it belongs to that exact owner/gate.

If no material finding exists, say so explicitly. Do not manufacture style findings.

### 5. Golden-flow coverage assessment

For GF-01, GF-02, GF-03 and SR-01 state whether each is necessary, sufficient and non-redundant. Name any missing distinct defect class.

### 6. Cross-cutting control assessment

Adjudicate Organization/auth/access/wire/idempotency/effect/Governance/Work/knowledge controls and whether they are truly exercised by the selected set.

### 7. Proof-horizon assessment

Separate:

- what D8 can prove now through repository/authority reasoning;
- what D4 still requires as separately authorized real external probe before D8 close;
- what must wait for post-D9 implemented real-dependency conformance.

### 8. Implementation-readiness assessment

Adjudicate D8-F1, D6-R2, the Pre-D9 Implementation Readiness Contract, import/dependency enforcement and the two-materially-different-implementations adversarial test.

### 9. YAGNI / Global Maximum assessment

State whether the candidate is the smallest sustainable proof/readiness structure or whether it is under/over-engineered.

### 10. Required corrections before D8 closeout

List only material corrections. Distinguish candidate fixes from later proof obligations and from rejected suggestions.

### 11. Reopen decision

State explicitly whether any D0–D7 authority must reopen now. If yes, name the **smallest owning authority and exact falsifier**. Otherwise say `NO D0–D7 REOPEN`.

---

## Fable response

<!-- Fable: append the independent review below this line. Do not alter the prompt above. -->
