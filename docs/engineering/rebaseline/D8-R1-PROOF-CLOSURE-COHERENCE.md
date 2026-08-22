# D8-R1 — Proof Closure & Implementation-Readiness Coherence

> **Status:** OPEN / ACTIVE — INDEPENDENT REVIEW ADJUDICATED; LIVE-PROBE DISPOSITION REQUIRED  
> **Parent:** [D8 Golden Flows](D8-GOLDEN-FLOWS.md)  
> **Candidate reviewed:** `stage/d8-golden-flows @ b3469258348289865a036bc7a946077f79d61faf`  
> **Independent review:** `review/d8-fable` — `ACCEPT WITH BOUNDED FIXES`  
> **Accepted prerequisites:** D0–D7 — ACCEPTED / CLOSED  
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose

D8-R1 is the smallest bounded correction produced by the independent whole-D8 challenge. It does not change the selected `GF-01 / GF-02 / GF-03 / SR-01` set, Product operations, Permissions, Principal kinds, semantic owners, D6 technology/topology or D7 runtime authority.

It closes four D8-local coherence seams:

1. explicit disposition and closure law for D4-deferred real probes;
2. bounded D8 revalidation if D6-R2 / pre-D9 readiness work materially changes accepted authority;
3. one concrete Governance-currentness falsifier inside the existing flow set;
4. dependency/import authority must cite the complete ratified D6 graph rather than a partial D8 restatement.

A fifth Fable item — missing final newline in `scripts/gate.ps1` — is valid repository hygiene but has no functional or architectural consequence. It is **LATER_NON_BLOCKING**, not a D8 closeout property.

## 2. Independent-review disposition

| Finding | GPT disposition | D8 action |
| --- | --- | --- |
| F-1 — D4 probe closure law/list incomplete | **ACCEPT WITH REFINEMENT** | full D4 probe ledger + closure law in §3 |
| F-2 — no post-D8 material-change revalidation law | **ACCEPT** | bounded revalidation law in §4 |
| F-3 — Governance mentioned but not exercised | **ACCEPT** | cross-cutting Governance variant in §5 |
| F-4 — partial D6 import graph restatement | **ACCEPT** | D6 §9.4 is sole dependency-direction authority in §6 |
| F-5 — `gate.ps1` final newline | **VALID / NON-BLOCKING HYGIENE** | no architecture change; repair on next bounded gate-file touch |

No Fable finding requires D0–D7 semantic reopen.

## 3. D4-deferred probe ledger and D8 closure law

D4 §7 names the real-effect probes whose semantic contracts are accepted but whose first controlled execution was deferred to D8. D8 must not silently drop them, and architecture approval is not authorization to execute them.

### 3.1 Canonical probe ledger

| ID | D4-deferred probe | D8 relationship | Current disposition |
| --- | --- | --- | --- |
| P1 | first Mercado Livre **Price/Availability effect** + authoritative convergence reread | GF-01; may be discharged by one controlled publication experiment only if the experiment explicitly exercises and records both meanings | **UNRESOLVED — OPERATOR DISPOSITION REQUIRED** |
| P2 | selected-lane **fiscal / invoice / label progression** | GF-02 | **UNRESOLVED — OPERATOR DISPOSITION REQUIRED** |
| P3 | first irreversible Sankhya **`313→306` fiscal progression** | GF-02; may share one controlled lane with P2 when separately evidenced in the result | **UNRESOLVED — OPERATOR DISPOSITION REQUIRED** |
| P4 | first consequential native **Party create/update when needed** | GF-02 Materialization prerequisite | **CONDITIONAL — trigger only if the selected flow actually requires native Party mutation** |
| P5 | first controlled **alternate-destination/contact realization** before claiming that concrete Sankhya capability | GF-02 | **UNRESOLVED — OPERATOR DISPOSITION REQUIRED** |
| P6 | any currently unexercised fiscal **branch/component that becomes material** to a selected golden flow | GF-02 | **CONDITIONAL — trigger only when the selected flow makes the branch/component material** |

D4-R1 separately states that D8 owns the first controlled real Mercado Livre creation/write proof with authoritative reread and shared-User-Product blast-radius verification. That proof may satisfy P1 only when its recorded scope actually includes the Price/Availability effect required by D4; equivalence is never inferred merely because both involve marketplace creation.

### 3.2 Allowed probe disposition states

For an **unconditional** D8 probe row, exactly one of these must exist before D8 closes:

```text
EXECUTED_AND_RECORDED
or
OPERATOR_RATIFIED_REDEFER(<named later gate>)
```

For a D4 probe whose own authority says **`when needed`** or **`becomes material`**, D8 additionally permits:

```text
NOT_TRIGGERED(<explicit evidence that the selected D8 flow does not depend on that branch>)
```

This third state is not a convenience defer. It preserves D4's conditional authority instead of manufacturing a live effect merely to satisfy ceremony. If the condition becomes true, `NOT_TRIGGERED` is invalid and the row must become executed or explicitly operator-redeferred.

### 3.3 Closure law

> **D8 cannot close while any unconditional D4-deferred probe is undispositioned, or while any triggered conditional probe is undispositioned.**

No live Mercado Livre write, native Party mutation, destination mutation or irreversible Sankhya fiscal effect is authorized by D8 documentation, Fable review or GPT adjudication. Execution requires separate explicit operator authorization.

A failed real probe narrows/reopens only the smallest implicated capability/authority. It never authorizes Direct Oracle fallback, fabricated known values, customer-master corruption, blind retry or provider-ontology leakage.

## 4. Post-D8 material-change revalidation law

D8-F1 intentionally schedules:

```text
D8 close
→ D6-R2 Complete Frontend Realization Closure
→ Pre-D9 Implementation Readiness Contract
→ D9
```

D6-R2 and the readiness contract must normally consume accepted D0–D8 authority without changing it. However, D6 already proved that realization work can expose a real missing Product seam, as happened with the accepted D6-R1 `95→99 / 29→30` repair.

Therefore:

> **If D6-R2 or the Pre-D9 Implementation Readiness Contract produces material evidence that, through the smallest owning authority, changes an accepted Product operation, ordinary Permission, semantic owner, Principal/client admission, identity/isolation law, external-effect law, or another invariant relied upon by a D8 flow, every affected D8 flow/control must be boundedly revalidated against the new accepted authority before D9 may open.**

Rules:

- no D0–D7/D8 reopen occurs merely because a wireframe is inconvenient;
- a real missing meaning returns to the smallest owner first;
- only D8 flows/controls whose assumptions changed are revalidated;
- Product-surface conservation in D8 means **the then-current accepted canonical Product surface**, not permanently frozen numerals after a legitimate bounded reopen;
- D9 never reviews golden flows known to have been derived from superseded authority.

## 5. Governance-currentness variant

Governance is a cross-cutting authority, not a fourth golden flow. One selected consequential path in GF-01 or GF-02 whose accepted business disposition **actually requires Governance** must exercise this variant:

```text
owner intent/action is valid
→ current Governance decision/delegation authorizes the exact action/scope
→ before external dispatch, that authorization is revoked / expires / becomes invalid
→ pre-dispatch currentness revalidation fails closed
→ no external write occurs
```

Binding falsifiers:

1. a stale or revoked Authorization Decision/Delegation cannot remain dispatch authority;
2. ordinary Permission cannot substitute for required Governance;
3. Governance approval/decision creation never itself executes the target action;
4. execution stays with the semantic action owner and is revalidated at execution time.

If none of the selected D8 consequential fixtures requires Governance under current accepted disposition, the proof must record that fact rather than fabricate a Governance requirement. D9/implementation readiness may still use the same invariant when a real configured action makes it reachable.

## 6. Dependency/import authority correction

D8-GOLDEN-FLOWS §12.5 contains a representative dependency sketch. It is **not** a complete restatement and must not become an independent closed-world import authority.

The sole accepted dependency-direction authority remains **D6 Frontend §9.4**, including:

```text
app/routes
    ↓
features
    ↓
api/<owner-family>
    ↓
api/transport
    ↓
api/generated

features ──→ ui
app/routes ─→ ui
```

plus all binding prohibitions in D6 §9.4.

D6-R2 / the Pre-D9 Implementation Readiness Contract may convert that complete graph into mechanical default-deny checks once a real source tree exists, but may not infer allowed/forbidden edges from the abbreviated D8 sketch.

## 7. Golden-flow disposition after review

Independent review found the `3 + 1` set necessary, sufficient and non-redundant after the bounded Governance variant:

```text
GF-01  Publication & Marketplace Convergence                CONFIRMED
GF-02  Sale → Business System → Fiscal → Fulfillment        CONFIRMED
GF-03  Performance Evidence Honesty                         CONFIRMED
SR-01  PITR / Timeline Continuity Recovery                  CONFIRMED
```

No separate golden flow is added for onboarding, Availability configuration, access administration, Work lifecycle, Post-Sale, Governance, every Product operation or every D7 proof obligation. Their distinct properties are either variants/cross-cutting controls or later real-runtime conformance obligations.

## 8. Proof horizons remain unchanged

### D8 can prove now

- authority/OAD trace and non-regression;
- representative defect-class coverage;
- negative-contract claims and owner separation;
- D7-R1 recovery-falsifier design;
- explicit D4 probe disposition/closure law.

### D8 may execute only with separate operator authorization

The P1–P6 ledger rows that are selected for `EXECUTED_AND_RECORDED`.

### Post-D9 implementation conformance

Real PostgreSQL/RLS, River, Keycloak, browser/CSRF, OAD router/validator, object-storage and PITR/restart execution; D6-R1 measure-by-scope proof; and every other D7 real-dependency falsifier that requires an implemented Product runtime.

## 9. Candidate closeout state

```text
independent D8 Fable challenge       COMPLETE
Fable verdict                        ACCEPT WITH BOUNDED FIXES
GPT adjudication                      F-1..F-4 ACCEPTED / F-5 NON-BLOCKING
D0–D7 reopen                          NONE
Golden-flow set                       3 business + 1 systemic — CONFIRMED
Product surface                       99 operations / 30 Permissions / H-A-S unchanged
active Product runtime                NONE
D8 real-probe ledger                  OPEN — operator disposition required
D6-R2                                 BLOCKED UNTIL D8 CLOSE
Pre-D9 readiness contract             BLOCKED UNTIL D6-R2
D9                                    BLOCKED
Product implementation                BLOCKED UNTIL D9
```

D8 is **not yet closeable** because unconditional live-probe rows P1/P2/P3/P5 remain undispositioned.

## 10. Exact next action

**Obtain explicit operator disposition for each currently unresolved D4-deferred D8 probe: either authorize the bounded real probe for execution/recording, or explicitly ratify re-deferral to a named later gate. Resolve conditional P4/P6 only if their triggering condition is present in the selected flow. Do not begin D6-R2, D9 or Product implementation until D8 closes.**
