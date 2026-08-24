# NOTIF-01 D8-R — AuthorizationRequest Golden-Flow Revalidation

> **Status:** OPERATOR-RATIFIED / ACCEPTED
> **Operator ratification:** 2026-08-24 — `Aprovado`
> **Trigger:** [D7-R AuthorizationRequest Runtime Repair](D6-R2-NOTIF-01-D7-R-AUTHORIZATION-REQUEST-RUNTIME-REPAIR.md)
> **Parent D8 authority:** [D8 Golden Flows](D8-GOLDEN-FLOWS.md) + [D8-R1 Proof Closure & Implementation-Readiness Coherence](D8-R1-PROOF-CLOSURE-COHERENCE.md)
> **Canonical Product:** 106 Product operations · 31 ordinary Permissions · H/A/S only
> **Active runtime:** NONE
> **Implementation:** BLOCKED UNTIL accepted D9

```text
D8R_GLOBAL_MAXIMUM:ACCEPTED
D8R_GOLDEN_FLOW_SET:THREE_BUSINESS_PLUS_SR01_UNCHANGED
D8R_PRODUCT_SURFACE:106_OPERATIONS_31_PERMISSIONS_HAS
D8R_GF01:REVALIDATED
D8R_GF02:REVALIDATED
D8R_GF03:NOT_MATERIALLY_AFFECTED
D8R_SR01:REVALIDATED_CONTINUOUS_TIMELINE_VS_PITR
D8R_REVIEW_BASIS:LISTING_INTENT_PRICE_INTENT_BUSINESS_ORDER_INTENT_INVOICING_INTENT
D8R_TYPED_503:KNOWN_NO_EFFECT_ONLY_FOR_EXACT_PROBLEM_TYPE
D8R_NOTIFICATIONS:F13_CURRENTNESS_F14_DECISION_OCCURRENCE
D8R_ZERO_DECIDER:WORK_NOT_AUTHORITY
D8R_PITR:IDEMPOTENCY_IS_NOT_CONTINUITY_ORACLE
D8R_LIVE_PROBES:NOT_REOPENED
D8R_P2:PRESERVED_FUTURE_PROOF
D8R_NEW_PRODUCT_OPERATIONS:0
D8R_NEW_PERMISSIONS:0
D8R_NEW_GOLDEN_FLOWS:0
D8R_NEW_INFRASTRUCTURE:0
D8R_RUNTIME_PROOF:POST_D9_IMPLEMENTATION
```

## 1. Purpose and bounded trigger

D8-R applies the accepted D8-R1 material-change revalidation law to the AuthorizationRequest redesign that D6-R2 legitimately carried through D1-R2, D2-R6, D3-R3, D5-R6, final P9, independent Fable closure and D7-R.

The original D8 closeout was derived against the then-current 99-operation / 30-Permission Product surface. Current canonical Product authority is 106 operations / 31 ordinary Permissions / H/A/S, with the delta concentrated in the bounded AuthorizationRequest/Notification/Work repair.

D8-R therefore revalidates only the D8 flow/control assumptions that materially changed. It does **not** reconstruct D8, reopen D4 provider contracts, create another business golden flow, repeat live probes for ceremony, add Product surface, select infrastructure or begin Product implementation.

## 2. Global-Maximum alternatives challenged

### A — no D8 revalidation

**REJECTED.** D8-R1 explicitly requires bounded revalidation when a legitimate post-D8 repair changes an invariant relied upon by a selected flow/control. Ignoring 99/30 → 106/31 plus the new Governance runtime semantics would leave D9 reviewing superseded assumptions.

### B — bounded GF-01 + GF-02 + cross-cutting/SR-01 revalidation

**SELECTED / OPERATOR-RATIFIED — GLOBAL MAXIMUM.** It covers every changed invariant while preserving D8's defect-class selection rule and keeping unrelated evidence closed.

### C — add a fourth business flow for approvals / Governance

**REJECTED / YAGNI.** Governance remains a cross-cutting authority. The selected consequential paths already contain the exact ListingIntent, PriceIntent, BusinessOrderIntent and InvoicingIntent meanings whose authorization episodes need falsification. A `GF-04` would duplicate defect classes without adding a distinct accepted invariant.

### D — rerun all D8 flows and controlled external probes

**REJECTED.** D7-R changes Governance/runtime composition, not Mercado Livre or Sankhya protocol meaning. Repeating P1–P6 without a changed external contract would add risk/ceremony rather than evidence.

## 3. Impact matrix

| D8 authority | D7-R impact | D8-R disposition |
| --- | --- | --- |
| GF-01 — Publication & Marketplace Convergence | material — `listing_intent` + `price_intent` authorization basis and decision/retry/currentness laws | **REVALIDATED** |
| GF-02 — Sale → Business System → Fiscal → Fulfillment | material — `business_order_intent` + `invoicing_intent` authorization basis and decision/retry/currentness laws | **REVALIDATED** |
| GF-03 — Performance Evidence Honesty | no new write, authorization basis, Permission implication or evidence meaning | **NOT MATERIALLY AFFECTED** |
| SR-01 — PITR / Timeline Continuity Recovery | material composition seam — D7-R replay evidence is continuous-timeline evidence, not rollback continuity proof | **REVALIDATED** |
| Product-surface conservation | original 99/30 snapshot superseded by legitimate bounded reopen | **REBASELINED TO 106/31/H-A-S** |
| Governance / Work / Notification controls | materially enriched by AuthorizationRequest, F13/F14 and zero-decider laws | **REVALIDATED** |
| D8 real probe ledger P1–P6 | no provider/business-system effect contract changed | **NOT REOPENED** |

The D8 set remains exactly three business golden flows plus SR-01.

## 4. Review-basis mapping into existing golden flows

D5-R6 admits exactly four closed AuthorizationRequest review-basis kinds:

```text
GF-01
  listing_intent
  price_intent

GF-02
  business_order_intent
  invoicing_intent
```

This mapping is why a separate Governance golden flow is unnecessary. The authorization semantics are exercised where the governed action already belongs.

No Performance review-basis kind exists, and D8-R does not invent one.

## 5. GF-01 revalidation — Listing / Price authorization

When current policy/disposition says a ListingIntent or PriceIntent requires Governance, the representative path composes as:

```text
exact Organization + exact Marketplace Installation
→ action-owner ListingIntent / PriceIntent
→ action owner establishes one authorization episode when required
→ Controlled Action Governance AuthorizationRequest
→ exact current eligible human decision path
→ immutable AuthorizationDecision
→ recoverable Decision propagation to action owner
→ action-owner execution-time currentness/material-validity revalidation
→ external dispatch only if the action is still admissible
→ provider reread / convergence remains authoritative
```

Binding revalidation laws:

1. `AuthorizationDecision(outcome=authorize)` is authorization evidence only. It never invokes Mercado Livre, submits a ListingIntent, changes price or executes the governed target by itself.
2. The governed action remains Offering/action-owner authority. Ordinary Permission does not substitute for a required Governance decision.
3. A historical Decision is immutable. If governing state changes after Decision commit and the action is no longer materially admissible, execution fails closed; the Decision is not rewritten into a rejection.
4. If the action owner still needs the action after material reauthorization conditions change, it creates a genuinely new authorization episode/request under D2-R6/D3-R3 rather than reopening the terminal historical request.
5. Pre-dispatch currentness remains required. A stale historical authorization cannot waive Availability, provider blast-radius, revision or other GF-01 execution-time laws.
6. Provider `2xx` remains weaker than authoritative convergence reread. Authorization does not strengthen provider transport evidence.

GF-01's external-effect laws are unchanged; D8-R adds the now-canonical authorization path as a cross-cutting variant.

## 6. GF-02 revalidation — Business Order / Invoicing authorization

When current disposition says BusinessOrderIntent or InvoicingIntent requires Governance:

```text
Sale / materialization or invoicing intent
→ action-owner authorization episode
→ AuthorizationRequest
→ exact current eligible-human decision
→ immutable AuthorizationDecision
→ recoverable propagation to exact action owner
→ revalidate current owner truth and material basis
→ continue Sankhya / Fulfillment progression only if still admissible
```

Binding revalidation laws:

1. Governance never becomes a Business-System Materialization, Fiscal or Fulfillment workflow owner.
2. An `authorize` Decision cannot waive Party ambiguity, destination safety, sanctioned Sankhya Gateway-only access, current source truth or authoritative reread.
3. Invoicing authorization cannot waive Fulfillment-owned physical readiness/conference. An action that becomes physically invalid after approval fails closed at execution time.
4. Assignment/escalation of a zero-decider Work item never grants `governance.decide`, physical qualification or business disposition.
5. Request/Decision state does not fabricate source-domain closure. Work close likewise cannot make a Sale, invoice, Shipment, refund or return converged.
6. No Direct Oracle, generic Sankhya command or approval-triggered provider shortcut is introduced.

The P2 future Mercado Livre invoice→label proof remains exactly where D8-R1 placed it; AuthorizationRequest planning does not fabricate a qualifying Sale or discharge that external proof.

## 7. Decision retry / concurrency revalidation

D7-R specializes `CreateAuthorizationDecision` idempotency to:

```text
organization_id
+ effective PrincipalID
+ CreateAuthorizationDecision operation identity
+ digest(Idempotency-Key)
```

Within one Principal namespace, fingerprint material includes request + supplied `If-Match` + outcome.

### 7.1 Lost success response

```text
Decision commits
→ AuthorizationRequest revision advances / becomes terminal
→ 201 response is lost
→ same Principal + same key + same request + same If-Match + same outcome
→ recover original committed Decision
→ BEFORE current request lock / stale If-Match / current eligibility re-evaluation
```

The replay must not become `412` merely because the request revision advanced as a consequence of the already-committed Decision.

### 7.2 Same raw key under another Principal

The binding D7-R law is preserved verbatim in meaning: **same raw Idempotency-Key under a different effective Principal** belongs to an independent namespace.

Therefore another eligible human using the same raw key string never receives replay/disclosure of the first human's Decision and never receives an artificial reused-key collision solely because the string matches. That human performs a genuinely new attempt, and current request state/revision/eligibility determines whether it can proceed.

### 7.3 Same Principal, changed semantic command

Same Principal + same key with changed request / `If-Match` / outcome remains a reused-key semantic mismatch. Idempotency-Key never becomes business identity or a way to bypass request concurrency.

## 8. Semantic 503 and ambiguous-result revalidation

Only this exact Product Problem type is the accepted semantic known-no-effect result:

```text
https://conexus.fun/marketplace-central/problems/product/authorization-validity-unavailable
```

For that exact typed result:

```text
material validity = UNKNOWN_OR_UNAVAILABLE
Decision           = none
request lifecycle  = unchanged / PENDING
new success intake = not retained as a false effect
```

Any bodyless, unparsable, proxy/infrastructure or otherwise non-matching `503` remains ambiguous potentially accepted. The client/frontend may not treat every `503` as `known failed; generate a fresh key and retry`.

Recovery of an ambiguous decision attempt preserves the same effective Principal, Idempotency-Key and semantic command until the established intake can be converged or the system has explicit evidence that a new attempt is semantically appropriate.

## 9. F13 / F14 Notification revalidation

### F13 — action required

A current eligibility occurrence may create actionable awareness only after PersonalNotifications revalidates:

```text
AuthorizationRequest still PENDING?
+ exact human still eligible now?
```

Delayed/replayed F13 after terminal request or eligibility loss creates no new actionable awareness. Notification possession remains awareness only and never grants `governance.decide`.

### F14 — decision result

A committed immutable Decision may independently feed the exact requester when accepted lineage permits it. Semantic duplicate identity is anchored by immutable Decision occurrence + exact recipient, so replay/rescue cannot create duplicate result Notifications.

F14 remains target-oriented and does not require the requester to possess `governance.read` merely to receive awareness of the result.

## 10. Zero-decider Work revalidation

```text
PENDING AuthorizationRequest
+ known-empty current eligible-human set
→ explicit Governance blocking condition
→ durable Operational Work materialization/reconciliation
```

There is no fallback approver, default admin, role inference or grace-period authority.

Work owns the responsibility/obligation lifecycle only. Work assignment, escalation or ownership never grants `governance.decide` and never changes AuthorizationRequest truth by itself.

When a valid decider becomes available again or the request becomes terminal, Governance/Work reconciliation removes/closes the no-longer-applicable obligation without rewriting historical Work or Notification evidence.

## 11. Invalidation and missed-wakeup recovery

Action-owner material invalidation may wake Governance through accepted event propagation, but event delivery is not sole current truth.

D7-R's fast wake-up lane is paired with durable recovery over PENDING Governance truth. A missed eligibility/invalidation event or scheduler tick cannot permanently leave a materially invalid request silently PENDING if current owner evidence can later prove invalidity.

The recovery mechanism re-Qs current Governance/request truth and action-owner validity; it does not infer currentness from River job existence or Notification state.

## 12. SR-01 revalidation — continuous timeline versus PITR

D7-R's exact replay guarantees and long-lived replay correlation operate on the normal **continuous durable database timeline**.

D7-R1 already establishes that a PITR/acknowledged-state rollback can erase database evidence that existed before the restore point, including dispatch markers, idempotency mappings, sessions/access changes and owner state. Therefore restored presence/absence of an AuthorizationRequest, AuthorizationDecision or idempotency record is not proof that no later acknowledged state ever existed.

Binding composition law:

```text
continuous positively established lineage
→ ordinary D7-R idempotency/replay/currentness laws apply

PITR / acknowledged-state rollback / lineage cannot be positively established
→ out-of-rollback-domain continuity witness fails or is unverifiable
→ recovery fence arms automatically
→ consequential external dispatch remains disabled
→ restored dispatchable/ambiguous work is reconciliation-only
→ current provider/business-system/access/integrity truth is reacquired
→ fence releases only for scopes whose safety is positively established
```

`Idempotency-Key` or its restored mapping is **not a continuity oracle**. The system may not say “this key/record exists or is absent, therefore the current database necessarily contains the latest acknowledged authorization/effect history.”

If required authorization/history/currentness cannot be positively established after rollback, the affected scope remains fenced for consequential dispatch and is surfaced for safe operator resolution/current reauthorization rather than guessed into safety.

No second idempotency database, external AuthorizationDecision mirror or new business authority is introduced by D8-R.

## 13. GF-03 disposition

GF-03 remains materially unchanged.

The AuthorizationRequest repair adds no Performance write, Performance authorization basis, measurement-basis rule, evidence aggregation, comparison rule, provider fact or `performance.read` implication.

Therefore:

```text
GF-03 choreography          UNCHANGED
GF-03 evidence semantics    UNCHANGED
GF-03 Permission boundary   UNCHANGED
```

Only the global current Product census used by cross-cutting Product-surface conservation changes to 106/31/H-A-S.

## 14. D8 live-probe ledger disposition

D7-R/D8-R does not change Mercado Livre or Sankhya provider/business-system effect contracts, so no live probe is reopened merely because Governance composition became more precise.

The accepted D8-R1 probe ledger remains:

```text
P1  Mercado Livre Price/Availability        PASS_CONVERGED
P2  invoice/fiscal/label progression        OPERATOR_RATIFIED_REDEFER
                                             first real open Sale / beta-flagged implementation drive
P3  Sankhya 313→306                         PASS_CONVERGED
P4  native Party create/update              NOT_TRIGGERED
P5  alternate destination/contact           CAPABILITY_NOT_PROVEN for full destination override
P6  additional fiscal branch/component      NOT_TRIGGERED
```

P2 remains explicit future proof and may not disappear from Pre-D9/D9/implementation readiness.

## 15. Required falsifiers

D8-R acceptance preserves a future proof contract capable of falsifying at least:

1. lost committed Decision success being returned as stale `412` instead of exact replay for same Principal/key/command;
2. another Principal receiving replay/disclosure or reused-key collision solely because the raw key string matches;
3. same Principal/key with changed request/`If-Match`/outcome avoiding reused-key mismatch semantics;
4. exact typed validity-unavailable `503` creating a Decision or mutating the PENDING request;
5. ambiguous/non-semantic `503` being treated as known failure and retried as a fresh semantic command;
6. delayed F13 creating new actionable awareness after terminal request or current eligibility loss;
7. duplicate F14 Decision delivery creating duplicate semantic result Notifications;
8. zero-decider materialization inventing a fallback approver or Work assignment granting `governance.decide`;
9. `AuthorizationDecision(authorize)` executing Listing/Price/BusinessOrder/Invoicing directly or bypassing action-owner execution-time validity;
10. a post-Decision material governing change still allowing consequential dispatch solely because historical authorization exists;
11. missed invalidation wake-up permanently stranding a request PENDING despite recoverable current invalidity evidence;
12. restored idempotency/AuthorizationRequest/Decision state being treated as proof of database timeline continuity after PITR;
13. recovery fence allowing blind consequential dispatch before rollback safety is positively re-established;
14. GF-03 being coupled to Governance, or a new `GF-04`, Product operation, Permission, workflow engine or infrastructure service appearing without a real accepted consumer.

Real PostgreSQL/River/browser/router/provider execution remains post-D9 implementation acceptance. This planning-stage proof may prove authority structure and falsifier presence only; it does not claim an active runtime.

## 16. Accepted result

```text
Golden-flow set                 3 business + SR-01 — UNCHANGED
GF-01                           REVALIDATED
  listing_intent                PASS
  price_intent                  PASS
GF-02                           REVALIDATED
  business_order_intent         PASS
  invoicing_intent              PASS
GF-03                           NOT MATERIALLY AFFECTED
SR-01                           REVALIDATED — continuous timeline vs PITR explicit
Product surface                 106 operations / 31 Permissions / H-A-S
Governance                      AuthorizationRequest + immutable AuthorizationDecision
Notifications                   awareness only
Work                            obligation only
new Product operations          0
new ordinary Permissions        0
new Golden Flows                0
new infrastructure              0
D8 live probes reopened         NO
P2 future proof                 PRESERVED
runtime real proof              POST-D9 IMPLEMENTATION
```

**Operator-ratified outcome:** D8-R is accepted as the smallest Global-Maximum revalidation of D8 against current AuthorizationRequest/D7-R authority. The D8 golden-flow set remains unchanged. NOTIF-01's D6-R → D7-R → D8-R closure chain is complete; B10 may return to the D6-R2 frontend realization sequence. Product implementation remains blocked until accepted D9, and PR #61 remains unmerged unless separately authorized.
