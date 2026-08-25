# D8 — Golden Flows

> **Status:** ACCEPTED / CLOSED / CURRENT CONSOLIDATED AUTHORITY  
> **Accepted prerequisites:** D0–D7  
> **Current Product:** 106 operations / 31 ordinary Permissions / H-A-S  
> **Active runtime:** NONE; implementation blocked until D9

## 1. Purpose / selection rule

D8 owns the smallest representative cross-boundary flow set needed to falsify accepted Product/ownership/API/frontend/runtime assumptions before D9.

> **A golden flow earns a place only when removing it would leave a material accepted invariant without a representative cross-boundary falsifier.**

Current irreducible set remains:

```text
GF-01 Publication & Marketplace Convergence
GF-02 Sale → Business System → Fiscal → Fulfillment → Outcome
GF-03 Performance Evidence Honesty
SR-01 systemic PITR/timeline-continuity recovery
```

AuthorizationRequest/Notifications are cross-cutting variants inside GF-01/GF-02/SR-01, not a fourth business flow.

D8 does not begin runtime implementation or add Product operations/Permissions by flow convenience.

## 2. Proof horizons

### Architecture / external-contract proof

May use accepted authority/OAD trace, adversarial falsifiers, deterministic repository proof and bounded operator-authorized external probes.

### Post-D9 implementation proof

Must exercise the same protected properties through real PostgreSQL/River/Keycloak/browser/router/OAD/object-store/provider/business-system seams. Mock-only tests cannot close those claims.

## 3. GF-01 — Publication & Marketplace Convergence

Representative choreography:

```text
Organization + exact Marketplace Installation
→ exact SourceInstance-qualified Product
→ Readiness search/requirements/correspondence
→ ListingIntent draft + FOLLOW_SOURCE / EXPLICIT_OVERRIDE
→ authored media if needed
→ separate PriceIntent when needed
→ Availability-owned quantity/current availability
→ Governance AuthorizationRequest/Decision when current disposition requires it
→ action-owner currentness/material-validity revalidation
→ owner-local durable intake / external-effect handoff
→ Mercado Livre realization
→ authoritative provider reread
→ Offering evaluates Listing convergence
→ Availability evaluates its own convergence
```

Binding falsifiers:

- source-following != human override;
- ListingIntent never owns Availability quantity;
- PriceIntent remains Offering actuation meaning even when Economics/Market/Performance informs it;
- required Availability input/blast-radius proof fails closed before dispatch;
- approval authorizes but never directly executes the target;
- stale authorization does not waive current source-owner validity/readiness;
- provider 2xx/201/202 != convergence;
- ambiguous possibly accepted effect is reconciled, not blindly replayed;
- stale owner revision != provider/business rejection;
- exact idempotent intake can replay safely; same key/different meaning fails;
- recurring automation cannot silently reverse standing human authored meaning;
- Organization/source mismatch fails closed.

Controlled D8 Mercado Livre Price/Availability probe converged on one bounded seller-owned lane and restored the fixture; this proves the selected D4 external contract, not a runtime implementation.

## 4. GF-02 — Sale → Materialization → Fulfillment → Outcome

Representative choreography:

```text
Mercado Livre Sale evidence + authoritative reread
→ current SellingEntity attribution when required
→ BusinessOrderIntent / Party Resolution / Destination Realization
→ Governance AuthorizationRequest/Decision when required
→ sanctioned Sankhya Gateway materialization + authoritative reread
→ Fulfillment separation / physical conference
→ InvoicingIntent + Governance when required
→ Sankhya fiscal progression + authoritative reread/correlation
→ provider-required artifacts/readiness
→ packing / dispatch handoff
→ source-qualified Shipment observation
→ Economics/reconciliation
→ bounded PostSaleResolution / Work for material exception consequences
```

Binding falsifiers:

- multiple native party candidates remain ambiguous; never first-result-wins;
- unsafe/unproven destination is explicit external-required/Work, not destructive customer-master mutation;
- Sankhya stays sanctioned API Gateway only; no Direct Oracle fallback;
- 2xx != materialization/fiscal convergence without authoritative reread;
- Fulfillment physical readiness gates normal invoicing where accepted;
- A cannot self-author physical checkpoints; S only under accepted physically-qualified path;
- Authorization cannot waive Party/destination/fiscal/physical/source-currentness conditions;
- InvoicingIntent and native fiscal result remain distinct;
- refund/payment and physical return consequences do not fabricate each other's closure;
- Work close does not close source business truth;
- ambiguous consequential effect is reconciled, not blind-retried.

D8 Sankhya `313→306` controlled probe converged on a test order. It also proved sanctioned `Contato` is not a full alternate street-address override, so full alternate destination remains explicit unsupported/external-required on that SourceInstance.

Mercado Livre invoice→label P2 remains **operator-ratified future proof** for the first qualifying real open Sale or beta implementation drive; D8 did not fabricate a Sale to satisfy ceremony.

## 5. GF-03 — Performance Evidence Honesty

Representative choreography:

```text
Organization + exact Marketplace Installation
→ explicit primary period
→ optional valid comparison period
→ Performance summary
→ Listing Performance collection/detail
→ Retail Media Performance
→ owner-separated frontend composition
```

Binding falsifiers:

- known zero != unknown/unavailable;
- period/comparison meaning remains explicit;
- provider measurement basis/coverage/provenance stays source-qualified;
- same metric name across scopes/providers != automatic equivalence;
- Listing population does not silently exclude unknown-performance Listings;
- provider-retention expiry does not rewrite preserved historical evidence as MPC-authored truth;
- Performance never mutates Offering/Ads or creates causal recommendations by measurement change;
- `performance.read` stays independent from Market/Economics/Offering authority.

AuthorizationRequest redesign does not materially affect GF-03 beyond the global current Product census.

## 6. AuthorizationRequest cross-cutting revalidation

Closed review-basis mapping:

```text
GF-01: listing_intent, price_intent
GF-02: business_order_intent, invoicing_intent
```

No Performance review basis exists.

For a governed action:

```text
action owner determines approval-required disposition
→ one AuthorizationRequest episode
→ exact current eligible-human decision
→ immutable AuthorizationDecision
→ recoverable propagation to action owner/requester awareness
→ action owner revalidates current business/material validity
→ external dispatch only if still admissible
```

Decision never substitutes for source-owner execution/currentness.

### 6.1 Current decision retry/concurrency

Current carrier/fingerprint:

```text
organization_id
+ effective PrincipalID
+ CreateAuthorizationDecision operation
+ Idempotency-Key digest

fingerprint:
authorization_request_id + body.etag + outcome
```

Lost-success exact replay returns the committed Decision **before current Request revision-precondition evaluation**. Same raw key under a different Principal is a different namespace. Same Principal/key with changed Request ETag/outcome is semantic mismatch.

No `If-Match`/412/428 remains on the custom `:decide` operation; stale typed Request revision is 409.

### 6.2 Semantic 503

Only the exact Product `authorization-validity-unavailable` Problem is known-no-effect:

```text
Decision none
Request remains PENDING
```

Bodyless/proxy/unparseable/non-matching 503 remains ambiguous potentially accepted and must not cause a fresh blind retry.

### 6.3 F13 / F14

F13 materializes only after current Request PENDING + exact-human eligibility revalidation; delayed stale action-required awareness is suppressed.

F14 is anchored to immutable Decision occurrence + exact requester and remains target-oriented; it does not grant Governance history/source authority.

### 6.4 Zero-decider Work

PENDING Request + known-empty current decider set becomes explicit Work/reconciliation, not a fallback admin/approver. Work assignment/escalation never grants `governance.decide` or changes Request truth.

## 7. SR-01 — PITR / timeline continuity

Normal D7 idempotency/replay laws assume positively established continuous durable timeline.

After PITR/acknowledged-state rollback or unverifiable lineage:

```text
continuity witness fails/unverifiable
→ recovery fence arms automatically
→ consequential external dispatch disabled
→ restored ambiguous/dispatchable work becomes reconciliation-only
→ current provider/business-system/access/integrity truth reacquired
→ release only scopes whose safety is positively established
```

Restored presence/absence of AuthorizationRequest/Decision/idempotency row is **not** proof no later acknowledged history existed. Idempotency storage is not a continuity oracle. No second Decision/idempotency database is introduced.

## 8. D8 live-probe disposition

Current accepted probe ledger remains:

```text
P1 Mercado Livre Price/Availability   PASS_CONVERGED
P2 invoice/fiscal/label progression   OPERATOR_RATIFIED_REDEFER
P3 Sankhya 313→306                    PASS_CONVERGED
P4 native Party create/update         NOT_TRIGGERED
P5 alternate destination/contact      full override CAPABILITY_NOT_PROVEN / external-required
P6 additional fiscal branch           NOT_TRIGGERED
```

AuthorizationRequest/Notification consolidation changed no provider/business-system effect contract, so P1–P6 are not rerun by ceremony.

Exact protocol/evidence details remain in `D8-LIVE-PROBE-PROTOCOL.md` and `D8-LIVE-PROBE-EVIDENCE.md` while those proof records have current consumers.

## 9. Implementation proof contract

When implementation opens after D9, the flow set must falsify at least:

- cross-Organization isolation/access confusion;
- stale/duplicate/ambiguous consequential intake/effect behavior;
- provider/business-system current reread/convergence;
- authorization vs execution-time validity separation;
- zero-decider/F13/F14 recovery;
- PITR external-effect continuity fencing;
- frontend/OAD route/permission/currentness composition;
- Performance knowledge/measurement-basis honesty;
- provider/business-system safety boundaries including Gateway-only Sankhya.

D8 architecture acceptance is not runtime implementation acceptance.

## 10. Reopen triggers

Reopen D8 only when a material post-D8 repair changes an invariant relied on by one selected flow/systemic control or when a distinct defect class has no representative falsifier. Do not add a flow by domain/screen count or rerun live external probes when their governing external contract did not change.
