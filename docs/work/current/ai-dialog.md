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

### 1. Verdict

`ACCEPT WITH BOUNDED FIXES`

### 2. Revalidation record

```text
remote main HEAD             8a9bcdb2c0333b94eecef8fb18cfa54b824b2881   (= expected)
candidate HEAD               b3469258348289865a036bc7a946077f79d61faf   (= expected; NOT STALE)
PR #59                       OPEN / Draft / MERGEABLE / base main / head stage/d8-golden-flows @ b3469258
changed files vs main        3 (docs/engineering/rebaseline/D8-GOLDEN-FLOWS.md, docs/roadmap.md, scripts/gate.ps1)
CI on candidate HEAD         required: SUCCESS; conventional-title: SUCCESS
review ancestry              review/d8-fable @ 83d8b936; merge-base with candidate = candidate HEAD
review isolation             candidate...review differs ONLY by docs/work/current/ai-dialog.md
local gate (review mode)     GATE_BASE_REF=stage/d8-golden-flows -> gate: PASS, review_mode: True,
                             changed_files: 1, legacy_runtime_population: 0, negative_controls: 1/1,
                             bootstrap_bytes 18418/20480 (~90%)
```

### 3. Executive coherence assessment

The candidate is a faithful composition of accepted D0–D7 authority, independently re-traced owner by owner:

- **GF-01** reproduces D4-R1 exactly: §5.2 choreography matches the R1-G1 joint-realization diagram; acceptance properties 1–11 map 1:1 onto D4-R1 §4/§5/§9/§10/§13 plus D7-B §11 (fingerprint conflict) and D7-C §7–§9 (marker/ambiguity classes). §5.5 carries the D4-R1 §13 deferral ("D8 owns the first controlled real Mercado Livre creation/write proof") verbatim in meaning.
- **GF-02** matches D4-B3: TOP 313 → 306 via `SelecaoDocumentoSP.faturar`, Destination Realization CONDITIONED/NOT-YET-CLAIMED (D4 §5.7), G3 "DEFER SAFELY → D8", Gateway-only/no-Direct-Oracle. The candidate's negative claim "no Product CreateBusinessOrderIntent / CreateInvoicingIntent" verifies against the canonical OAD: `business-order-intents` and `invoicing-intents` expose GET-only collection/item plus party-resolution/destination-realization reads and one `party-resolution:resolve` POST; no create/command operation exists.
- **GF-03** reproduces D6-R1 §9–§11 with no drift: all four operationIds exist exactly as named in `paths-performance.yaml`; properties 1–10 are D6-R1 negative controls 4, 8–12 plus the W3 survivorship law and the D2-R2 custody law.
- **SR-01** is the accepted D7-R1 §4 continuity law re-stated as a falsifier, including the F-1-accepted arming law (witness absent/unverifiable ⇒ fence arms, no manual flag), the complementary continuous-lineage release control, and the "pre-engaged fence proof is insufficient" clause. Nothing in SR-01 weakens or re-litigates D7-R1.
- **D8-F1** characterizes D6 accurately: "39 representative route/screen states" is literal D6-B1 text (line 349: "39 estados de tela/rota (S00–S127 não contíguos) para 99 operações"), and the finding correctly claims no D6 reopen inside D8.

Non-regression fence independently verified on the exact candidate: the diff touches no contract file; the bundled OAD surface remains 99 referenced operations (104 source pathItem operations minus the 5 duplicated access pathItems in `paths-identity-portfolio-readiness.yaml` that `openapi.yaml` does not reference), 30 ordinary Permissions plus the `authenticated` pseudo-requirement, Principal kinds only H/A/S (58×[H,A,S], 14×[H,A], 30×[H], 2×[H,S]), exactly 2 physically-qualified operations (`RecordPhysicalConference`, `RecordDispatchHandoff`, both [H,S] — A excluded, matching GF-02 falsifiers). Stable origin and `active runtime NONE` unchanged in roadmap.

Method alignment is good: selection by defect class, not domain symmetry; proof horizons explicitly split so architecture evidence cannot masquerade as runtime conformance; YAGNI rejections are argued, not asserted. The candidate is ready for GPT adjudication after the bounded fixes below.

### 4. Material findings

**F-1 — D8 probe-closure law is undefined and the D4-deferred probe enumeration is a proper subset of D4's own list**
- **classification:** `D8_FIX` · **severity:** Important
- **location:** D8-GOLDEN-FLOWS §4.1, §5.5, §6.6, §13/§14; roadmap "D8 boundary" paragraph.
- **governing authority:** D4 §7 ("D8 still owns the first controlled real effects…", six enumerated items) and D4 §4.7 residuals R1/R2; D4-R1 §13.
- **counterexample:** D4 §7 enumerates six deferred probes: (1) first ML Price/Availability effect + authoritative reread; (2) selected-lane fiscal/invoice/**label** progression; (3) first irreversible 313→306; (4) first consequential native party create/update; (5) first controlled alternate-destination/contact realization; (6) any unexercised fiscal branch material to a selected flow. The candidate names (3)+(4)+(5) in §6.6 and the ML creation/write in §5.5; the standalone Price/Availability write (R1) is only obliquely covered if the creation probe is agreed to discharge it, and label progression and the unexercised-fiscal-branch clause are not carried at all. Worse, no clause states the closure law: if the operator never grants live-write authorization, does D8 stall indefinitely, or close with the obligations silently evaporating? Both failure modes are real and the second is the classic declared-risk-becomes-debt defect.
- **smallest correction:** one disposition table in D8 listing every D4-deferred probe with exactly one of: `EXECUTE + RECORD BEFORE D8 CLOSE` or `EXPLICIT OPERATOR-RATIFIED RE-DEFER to <named later gate>`; plus one sentence: D8 cannot close while any row is undispositioned.
- **owner:** D8 is the named receiving gate for these obligations; no D4 reopen — D4's list is unchanged, D8 must mirror and disposition it.

**F-2 — no revalidation law if D6-R2 or the Pre-D9 contract materially changes accepted surface/owner/invariant**
- **classification:** `D8_FIX` · **severity:** Important
- **location:** D8-GOLDEN-FLOWS §12.2–§12.6; roadmap "Approved pre-D9 realization sequence".
- **governing authority:** D8 §9 conservation control ("a required operation 100 or Permission 31 is a material finding"); Method global-coherence rule.
- **counterexample:** precedent exists — D6's own realization work surfaced D6-R1 and legitimately moved the surface 95→99/29→30. If D6-R2 similarly surfaces a missing material operation/Permission via the smallest owning authority, the accepted D8 flows and the §9 "exactly 99/30" conservation row become silently stale, and D9 would adjudicate golden flows validated against a surface that no longer exists. The candidate is silent on this ordering (Challenge H sequencing question).
- **smallest correction:** one clause in §12.2: any D6-R2/Pre-D9-contract change that alters an accepted operation, owner, Permission or material invariant (only via the smallest owning authority) triggers bounded revalidation of the affected D8 flow(s) before D9 entry, and the §9 conservation control re-baselines to the then-current accepted surface rather than the frozen numerals.
- **owner:** D8 owns flow validity across the operator-approved post-D8 sequence it itself introduces.

**F-3 — Governance is bound as a cross-cutting control but exercised by zero variants**
- **classification:** `D8_FIX` · **severity:** Minor
- **location:** D8-GOLDEN-FLOWS §9 (Governance/Work row), §5.4, §6.5.
- **governing authority:** D1 Controlled Action Governance ("Governance owns authorization when required; execution stays with the action owner"); D7-C §7.1 step 5 (pre-dispatch Governance currentness revalidation); D5-B2 pipeline ("Governance when required").
- **counterexample:** every other cross-cutting row (isolation, idempotency, physical qualification, knowledge honesty, effect safety) has at least one concrete falsifier row in §5.4/§6.5/§7.4; Governance has none. An implementation that treats a stale/revoked AuthorizationDecision as still valid at dispatch time, or lets Governance approval itself trigger execution, would pass the entire current minimum matrix.
- **smallest correction:** one variant row in GF-01 or GF-02: "disposition requires current Governance authorization; authorization revoked/expired between approval and dispatch → fails closed pre-dispatch; Governance approval never executes the target." No new flow.
- **owner:** variant addition inside an existing selected flow — exactly the candidate's own rule for cross-cutting properties.

**F-4 — §12.5 partially restates the ratified D6 dependency graph, dropping the `app/routes → ui` edge**
- **classification:** `D8_FIX` · **severity:** Minor
- **location:** D8-GOLDEN-FLOWS §12.5 diagram.
- **governing authority:** D6 §9.4 (allowed direction is closed-world and includes both `features ──→ ui` and `app/routes ─→ ui`).
- **counterexample:** the Pre-D9 readiness contract is told to "convert material dependency directions into enforceable acceptance rules" starting from this §12.5 sketch. Under a closed-world/default-deny import rule, the omitted edge makes a D6-legal edge (`app/routes → ui`) forbidden by omission — a disguised amendment of a ratified mapping via incomplete restatement.
- **smallest correction:** add the missing edge or replace the diagram with a bare citation of D6 §9.4 as the single source.
- **owner:** D8 text fix; D6 unchanged.

**F-5 — `scripts/gate.ps1` lost its trailing final newline in the candidate**
- **classification:** `REPOSITORY_FIX` · **severity:** Minor
- **location:** candidate diff, last hunk of `scripts/gate.ps1` (`\ No newline at end of file`).
- **counterexample:** none functional (gate runs green); POSIX-tool and future-diff hygiene only. Prettier's lane does not cover `.ps1`, so nothing will ever flag it.
- **smallest correction:** restore the final newline in the candidate branch.

No other material finding exists. I did not manufacture style findings.

### 5. Golden-flow coverage assessment

- **GF-01 — necessary, sufficient, non-redundant.** Sole representative of the consequential-marketplace-write class: joint realization, fail-closed missing owner input, blast radius, idempotency fingerprint, override supersession, convergence-vs-transport. Removing it leaves D4-R1's entire §13 falsifier family unrepresented.
- **GF-02 — necessary, sufficient with the F-3 variant, non-redundant.** Sole representative of party-ambiguity, destination safety, sanctioned-progression, physical-authority and cross-consequence-closure classes. The Post-Sale branch is correctly a branch: its distinct defect class (one consequence falsely closing another; Work close fabricating source closure) is present in the §6.5 table, so a standalone Post-Sale journey would add localization ceremony, not a new class. Not a hidden mega-flow: splitting it would sever exactly the cross-owner seams whose composition is the thing under test.
- **GF-03 — necessary, sufficient, non-redundant.** Not redundant with Market/Economics reads: measurement-basis comparability, the survivorship-bias population law, scope-collapse prohibition, custody-vs-authorship and `performance.read` non-implication exist only in the Performance contract, and the 13th boundary is the least battle-tested accepted authority. No missing distinct Performance invariant requiring another path; the measure-by-scope obligation is D6-R1's explicitly assigned post-D9 real-dependency proof.
- **SR-01 — necessary, sufficient, non-redundant, correctly systemic.** GF-01/GF-02 crash variants prove continuous-timeline safety; only SR-01 falsifies the rollback case where marker absence stops being evidence. It restates accepted D7-R1 §4 without weakening it, covers both false-negative continuity (fence must arm) and false-positive permanence (lineage-proved failover need not fence), and remains meaningfully executable post-D9.
- **Missing distinct defect class:** none requiring a new flow. Candidates examined and rejected: channel onboarding/OAuth ingress (composition already whole-reviewed at the D5-B2 whole-ingress restructure; ingress is non-Product surface with its own falsifiers; implementation conformance re-executes them), Availability configuration (owner-separation class already in GF-01), access administration (revocation currency exercised via SR-01 + pre-dispatch gates), partial-acquisition absence≠closure (representable as an existing-flow variant under the knowledge-honesty control; D4-R1 proof obligation 3 already owns it).

### 6. Cross-cutting control assessment

Exercised, not merely mentioned — with one exception:

- **Surface conservation:** independently recounted (99/30/H-A-S, 2 physical-qualification ops); the §9 row is checkable against the bundled OAD.
- **Organization isolation:** exercised by the GF-01 "Organization/reference mismatch" variant; structural falsification (RLS/composite-FK/pool leakage) correctly assigned to post-D9 real PostgreSQL per D7-B §13.
- **Auth carriers:** distinct-carrier and CSRF laws bound; dual-carrier ambiguity and carrier/scheme mismatch remain D7-R1 F-4 negative controls at implementation — acceptable as cross-cutting.
- **Permission/kind/physical-qualification/Governance non-equivalence:** kind and physical qualification exercised by GF-02 variants; **Governance currency not exercised anywhere → F-3.**
- **Wire authority / frontend:** GF-01 property 11 + the GF-02 §6.2 negative list verified against the OAD (no screen-shaped or command operations exist to leak).
- **Idempotency/revision:** exercised by three GF-01 variants; 412-vs-409 distinctness is D7-B §13.12.
- **Durable effects / reconciliation:** exercised by timeout/crash variants in both business flows plus SR-01.
- **Work non-authority:** exercised by the GF-02 "Work close" variant.
- **Knowledge honesty:** exercised by the whole GF-03 matrix plus GF-02 party-ambiguity.
- **Technical exclusion:** the Direct-Oracle prohibition is exercised by GF-02 property 3 and D7-C falsifier 16.

### 7. Proof-horizon assessment

- **Provable now (D8, repository reasoning):** everything in §3 above — flow/authority correspondence, OAD surface conservation, negative claims about nonexistent operations, falsifier-table representativeness, SR-01/D7-R1 identity.
- **Requires separately authorized real external probes before D8 close (or explicit ratified re-defer — F-1):** the D4 §7 list — controlled ML creation/write with authoritative reread + shared-UP blast-radius verification (discharging the Price/Availability effect if so recorded), controlled Destination Realization, first 313→306 with reread/correlation, selected-lane label/artifact progression, unexercised-fiscal-branch check. These are external-contract evidence and do not violate the implementation gate: they exercise provider/ERP contracts, not a Product runtime.
- **Must wait for post-D9 implemented conformance:** all D7-B §13 PostgreSQL/RLS proofs, D7-C §15 River proofs, D7-R1 §4 recovery falsifiers (both unarmed-restore and continuous-lineage probes), D7-R1 F-4 validator/carrier negative controls, browser/CSRF behavior, object-store integrity, and the D6-R1 measure-by-scope proof. No D8 language claims any of these currently execute; §4.2 and roadmap "active runtime NONE" are honest.

### 8. Implementation-readiness assessment

- **D8-F1 is material and correctly owned.** The architecture-closure ≠ implementation-readiness distinction is real (D6 proved 39 representative states for 99 operations — representative, not exhaustive; a coding agent would still decide material screen states). Surfacing it during whole-composition review is exactly where it becomes visible; filing the work as D6-R2 + a pre-D9 gate rather than doing it inside D8 avoids a D6 reopen by convenience. Not ceremony: neither D6 (closed as architecture) nor D9 (adversarial review, not authoring) implies this completion.
- **D6-R2 obligations are proportionate:** derivation from 99/30+D7 rather than visual-first, shared annotated wireframes over combinatorial images, and the fabricate-nothing/return-to-smallest-owner stop rule are the right guards against both under- and over-specification.
- **Realization Contract:** the field set (outcome, authorities, operations, permissions, state ownership, runtime properties, dependency directions, negative controls, executable acceptance, real-dependency proof, exit state) passes the two-materially-different-implementations test for authority/contract/dependency/acceptance dimensions; residual freedom (styling, private decomposition, naming) is exactly the freedom §12.4 intends to leave. The D9 entry challenge in §12.6 is the correct adversarial closure for whatever ambiguity survives.
- **Import/dependency enforcement:** converting D6 §9.4 into default-deny acceptance rules is right, and D6 itself already deferred the mechanism until a real source tree exists — but see F-4: the conversion must start from the complete ratified graph, not §12.5's partial sketch.
- **Sequencing:** correct except for the missing revalidation law — F-2.

### 9. YAGNI / Global Maximum assessment

The candidate is the smallest sustainable structure. 3+1 is genuinely irreducible: each row is the sole representative of a disjoint defect-class family (consequential marketplace write; cross-owner business/fiscal/physical composition; knowledge honesty; acknowledged-state rollback), and every proposed expansion in §11 was checked and correctly rejected — no rejected row names a defect class the selected set cannot falsify. Compression into mega-flows is correctly rejected for localization. D8-F1 does not create speculative machinery: it adds one completion pass and one contract seam, both consumed by the already-required D9; the Realization Contract is a checklist, not a framework. No over-engineering found; under-engineering is limited to F-1/F-2/F-3.

### 10. Required corrections before D8 closeout

**Candidate fixes (block closeout, not adjudication):**

1. F-1 — probe disposition table + closure law.
2. F-2 — post-D8 material-change revalidation clause.
3. F-3 — one Governance-currency variant row.
4. F-4 — complete or cite the D6 §9.4 graph.
5. F-5 — restore `gate.ps1` trailing newline.

**Later proof obligations (not candidate text):** execution/recording of the D4-deferred probes under explicit operator authorization per the F-1 disposition table; all post-D9 conformance items in §7.

**Rejected suggestions (do not adopt):** standalone flows for onboarding, Post-Sale, Work lifecycle, Availability config, access administration or per-operation coverage; any exhaustive scenario catalog; any second recovery flow.

### 11. Reopen decision

**NO D0–D7 REOPEN.** Every candidate claim traced to an accepted owner without contradiction; all five findings are D8-local text/bookkeeping corrections plus one repository hygiene item. No falsifier implicates any accepted D0–D7 semantic, ownership, safety or proof decision.
