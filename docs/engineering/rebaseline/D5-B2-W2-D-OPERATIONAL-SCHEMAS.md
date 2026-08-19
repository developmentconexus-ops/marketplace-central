# D5-B2 — W2-D Operational Lifecycle Schema Grammar

> **Status:** ACCEPTED IN-STAGE  
> **Parent W2:** `D5-B2-W2-SCHEMA-GRAMMAR.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + ratified D5-B2 Whole-Matrix + Wire W1 + W2-A/B/C  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Accepted:** 2026-08-18

## 1. Governing invariant

> **Operational schemas expose durable owner state, source-qualified external results and explicit owner-specific checkpoints without constructing a shared workflow/status authority. A client may request an owner capability; it never writes another owner's fact, provider result, approval, physical evidence or closure by field mutation.**

This section covers Controlled Action Governance, Marketplace Sales, Business-System Materialization, Fulfillment Lifecycle, Post-Sale Resolution and Operational Work.

It does not create a generic OrderWorkflow, Action/Operation, Task/Case, provider resource graph, ERP ontology or cross-owner aggregate.

## 2. Controlled Action Governance

### 2.1 AuthorizationDecision is an immutable occurrence

`AuthorizationDecision` is a durable Governance-owned occurrence with proportionately:

- opaque MPC Decision ID;
- typed target reference;
- exact target material revision/context decided;
- decision outcome;
- intended-scope snapshot;
- authorized-scope snapshot;
- authority/delegation context;
- server-attributed decision Principal;
- decision time.

A Decision is never PATCHed into a different historical decision. Reapproval/rejection after an earlier decision creates a later occurrence and preserves prior history.

### 2.2 Typed Authorization target; no universal IntentRef

Baseline `AuthorizationTarget` is a closed union of concrete Product 1.0 owner intents actually subject to Governance, including proportionately:

- ListingIntent;
- PriceIntent;
- BusinessOrderIntent;
- InvoicingIntent.

A future real action class adds the smallest concrete variant. No `{entity_type, entity_id}` / universal BusinessIntent reference is introduced.

### 2.3 Referenced-resource revision proof

Creating an AuthorizationDecision is a POST to a new Decision resource, while the stale-state safety requirement belongs to the referenced target Intent. Therefore the request preserves the target resource's opaque strong validator as a typed **referenced-resource precondition**; `If-Match` remains reserved for conditioning the selected HTTP request resource.

Conceptually:

```json
{
  "target": {
    "kind": "price_intent",
    "price_intent_id": "...",
    "target_etag": "\"opaque-validator\""
  },
  "decision": "authorize"
}
```

The exact reusable spelling of referenced-resource preconditions is finalized in W2-E. The property frozen here is that the exact target material revision is mechanically proven, not reconstructed from mutable target state after the fact.

### 2.4 Authorized scope baseline

For Product 1.0 baseline, `authorize` applies to the exact current intended target scope. This satisfies `authorized_scope ⊆ intended_scope` by equality and avoids inventing a generic subset-selection DSL without a proven partial-approval consumer.

A future real multi-target/partial approval need adds an owner/target-specific scope-selection contract; it does not create a generic authorization-expression language.

### 2.5 AuthorizationDelegation gains stable wire identity

D2 intentionally did not freeze Grant/Delegation physical identity by convenience. W2-D now has a concrete update/revoke consumer and a real stability property: delegate/action/scope may change without the standing delegation becoming a different addressed resource.

Therefore `AuthorizationDelegation` receives a stable opaque MPC-owned ID for the Product wire contract.

This does not create a generic Grant domain or IAM rules engine.

Delegation schema carries proportionately:

- delegation ID;
- delegate Principal;
- bounded action class;
- bounded scope;
- server-owned establishment/update/revocation provenance.

No arbitrary condition/expression/query DSL is admitted.

Delegation update uses typed PATCH + `If-Match`. Delegation revoke uses an explicit `:revoke` capability and retains the ratified fail-safe/monotonic rule: stale client snapshots do not leave the targeted standing authority active; a later re-grant is a new explicit authority action.

## 3. Marketplace Sales

### 3.1 Sale remains source-qualified external identity

`MarketplaceSale` has no synthetic MPC Sale ID. Identity is the Marketplace Installation-qualified native Sale/Order key.

The Product representation contains only Sales-owned interpretation/context needed by downstream Product 1.0 consumers, proportionately including:

- source-qualified Sale reference;
- current interpreted sale meaning;
- transaction-specific Selling Entity attribution;
- sale lines / quantities / commercial transaction facts required downstream;
- materially post-sale-relevant committed facts;
- source observation/provenance/freshness.

Provider Pack/Shipment/payment DTO topology is not absorbed into Sales.

### 3.2 Sale-line selector key

Post-Sale/Fulfillment scopes need stable line/quantity selection without exposing provider array position or minting a canonical SaleLine entity.

Each interpreted Sale line therefore exposes an opaque `sale_line_key` stable within that Sale interpretation. It is a selector/correlation key scoped to one Sale, not a global MPC identity.

### 3.3 Selling Entity attribution

Attribution is an owner-specific state such as resolved / ambiguous / unresolved / unavailable where materially reachable.

A resolve request selects one eligible same-Organization Selling Entity ID against the current Sales attribution state and uses the current Sale/attribution validator. The caller does not submit legal/company/provider facts as replacement truth.

## 4. Business-System Materialization

### 4.1 BusinessOrderIntent is owner-triggered read/tracking state

Product clients never create BusinessOrderIntent. The read/tracking resource carries proportionately:

- MPC Intent ID;
- source-qualified Sale reference;
- target SourceInstance;
- Party Resolution;
- Destination Realization;
- current prerequisites;
- external-effect state/evidence;
- source-qualified native business-order result when established;
- convergence;
- material Work references/history.

No Product schema exposes Sankhya TOP/NUNOTA/status/choreography as MPC semantics. A native result is represented by SourceInstance-qualified native business-document identity.

### 4.2 PartyResolution is a contained singleton meaning

Business-System Party Resolution has durable meaning but no independent global identity is required. It is naturally identified as the Party Resolution of one BusinessOrderIntent and may therefore use containment such as conceptually:

```text
/business-order-intents/{id}/party-resolution
```

This is legitimate identity containment, not process-order nesting.

Party Resolution states distinguish proportionately unresolved, ambiguous, resolved-existing, resolved-new and unavailable where materially reachable.

Ambiguous candidate entries use opaque Materialization-local `candidate_key` plus the minimum PII-safe evidence required for human disambiguation. Raw `CODPARC`/provider DTOs are not Product vocabulary.

### 4.3 ResolveBusinessSystemPartyResolution request

The request is a closed union:

```text
use_existing(candidate_key)
establish_new_from_current_sale_identity
```

The second variant does not accept arbitrary Customer/Party master fields. The human decision is only that no compatible match exists and a new native party may be established from currently authoritative Sale/source identity facts. D4 owns the native creation protocol.

The operation retains mandatory `Idempotency-Key` plus current PartyResolution/candidate-set precondition. Intake deduplication never authorizes blind replay of an ambiguous native effect.

### 4.4 DestinationRealization remains read-only baseline

Destination Realization is a BusinessOrderIntent-contained singleton meaning with honest owner states such as realized, external-required, unsupported, ambiguous, unknown or unavailable where material.

No Product Address/Contact/Partner CRUD exists. No destination resolve/write capability is admitted until D8 proves the selected safe lane.

### 4.5 InvoicingIntent is owner-triggered read/tracking state

Product clients do not create InvoicingIntent. The resource carries proportionately:

- MPC Intent ID;
- BusinessOrderIntent reference;
- current Fulfillment physical-readiness reference/evidence;
- prerequisites;
- external-effect state;
- zero/one/many source-qualified fiscal/native results where material;
- convergence/history.

TOP 306 or provider-native status is not the canonical Product identity/state.

### 4.6 No giant Materialization status

Materialization keeps distinct owner axes such as prerequisites, external-effect state, native-result knowledge and convergence. These do not become a platform/global `OperationState` or writable status field.

## 5. Fulfillment Lifecycle

### 5.1 FulfillmentExecution wire identity

The admitted physical checkpoint operations need one stable Fulfillment-owned resource on which to act. Using the Sale resource would place Fulfillment capabilities on Sales-owned identity and would make future split physical scopes expensive.

D2 already permits durable Fulfillment routing/dispatch intent identity where material. W2-D therefore crystallizes **FulfillmentExecution** as the MPC-owned durable identity for the accepted physical fulfillment lifecycle.

This is not a new D1 domain, WMS/TMS model or generic Workflow. It is the concrete Fulfillment-owned identity required by the already admitted checkpoint/query surface.

A Sale may have 0..N FulfillmentExecutions when real physical scope requires it; Product 1.0 may normally have one without collapsing the identity model.

### 5.2 FulfillmentExecution scope

Scope is explicit and Sale-relative:

```text
whole_sale
sale_lines[{sale_line_key, quantity}]
```

No Product IDs or provider item positions substitute for the Sale-specific physical scope.

### 5.3 FulfillmentExecution representation

Carries proportionately:

- FulfillmentExecution ID;
- source-qualified Sale reference;
- scope;
- selected FulfillmentNode where established;
- separation checkpoint;
- physical-conference checkpoint and physical-readiness conclusion;
- provider-readiness/requirement-closure state;
- packing checkpoint;
- dispatch-handoff checkpoint;
- correlated Shipment references;
- applicable provider deadlines and MPC internal target evidence;
- material Work/history references.

No single `status` replaces these axes.

### 5.4 Physical checkpoints are occurrences, not writable state enums

Separation, physical conference, packing and dispatch handoff are explicit owner capabilities/checkpoints. Their responses preserve server-attributed Principal/source/time and material checkpoint evidence.

The request cannot claim effective actor or `trusted_physical_evidence=true`.

### 5.5 PhysicalConference distinguishes occurrence from readiness

A physical conference may be recorded and reveal a discrepancy. Therefore:

```text
conference recorded != physical readiness confirmed
```

The conference request/result supports at least confirmed versus discrepancy semantics with bounded discrepancy detail when material. The resulting Fulfillment-owned physical-readiness conclusion remains separate.

### 5.6 Client-class physical-evidence fence

`RecordPhysicalConference` and other checkpoint operations whose fact requires physical establishment accept only the admitted client classes from B2:

- human Principal; or
- explicitly provisioned/proven system Principal/source capable of establishing that physical fact.

An ordinary automation Principal does not acquire this epistemic capability merely through `fulfillment.execute`. Client class cannot be escalated through the request body.

### 5.7 Provider readiness and artifacts

Provider-readiness state distinguishes Fulfillment-owned closure conclusion from source-qualified provider requirement/deadline evidence. Provider checklist/status vocabulary remains D4-local.

Fulfillment artifacts are scoped to one FulfillmentExecution using an opaque local `artifact_key`; no global Artifact/ProductAsset identity is created. Artifact descriptors are PII-minimized and source-qualified where material. Blob/content delivery/storage remains D7 mechanism.

### 5.8 Shipment remains external

Shipment has Marketplace Installation-qualified native identity and no synthetic MPC Shipment ID. The Product response exposes only Fulfillment-interpreted shipment/delivery progress, material deadlines/outcomes/exceptions, freshness/provenance and bounded provider enrichment needed by the Product 1.0 consumer.

Provider Pack/status/substatus ontology is not mirrored.

## 6. Post-Sale Resolution

### 6.1 Explicit scoped obligation

`PostSaleResolution` remains an MPC-owned canonical identity. Its scope is explicit:

```text
whole_sale
sale_lines[{sale_line_key, quantity}]
```

A Product-created Resolution creates a coordination obligation, not immediate provider/ERP/refund effects.

### 6.2 Concerns are a set, not one mutually exclusive type

Initiating concerns may include cancellation, return and refund simultaneously where applicable. This is deliberately not a single `type = RETURN|REFUND|CANCELLATION` enum because consequences may coexist or be causally related.

The concern set describes why the Resolution exists; it does not command every consequence owner to execute those actions.

### 6.3 Independent consequence tracks

The Resolution response tracks consequences independently, proportionately including:

- Sales cancellation consequence;
- physical return/reverse-logistics consequence;
- refund/payment consequence;
- business-system/fiscal consequence;
- economic adjustment/reconciliation consequence.

Each track contains owner/result references and current evidence state rather than provider `available_actions` strings.

### 6.4 Closure is server-owned

Resolution lifecycle may expose open/closed meaning, but no direct Product `:close`/PATCH-close operation exists. Post-Sale closes only after the applicable consequence owners provide sufficient committed evidence.

## 7. Operational Work

### 7.1 Work remains canonical obligation identity

`Work` carries proportionately:

- Work ID;
- typed originating condition;
- responsibility role;
- optional assignment;
- hold/resume lifecycle;
- escalation target/state;
- source-condition reference/current reconciliation state;
- closure provenance/history/time.

Work is not a free-form user Task/Case.

### 7.2 Work origin is a Work-local closed union

No universal `{entity_type, entity_id}` graph is introduced. Work uses only concrete origin variants proven by Product 1.0, such as readiness, listing/price intent, sale, materialization, fulfillment, shipment, post-sale, economic attribution or authorization conditions, each with the appropriate typed subject/reference.

New condition classes add the smallest Work-local variant when real.

### 7.3 Responsibility role is not AccessRole

Work uses a `responsibility_role_key` or equivalent Work-owned semantic role. It is not an AccessRole/Permission identity. Assignment to a Principal does not alter ordinary access state.

### 7.4 Assignment / hold / resume / escalation

Assignment is unassigned or assigned to one same-Organization Principal and uses current Work precondition.

Hold/resume are explicit owner capabilities, not writable `status`.

Escalation is declarative to an explicit target responsibility/state. Increment/occurrence-style escalation is not baseline because it would change the ratified structural-idempotency assumption.

### 7.5 No generic Work resolution/close capability

No `SubmitWorkResolution` / direct close is admitted. Source owners retain the concrete closure meaning and Work reconciles source-owner committed resolution/coverage/supersession evidence.

## 8. Work closure-path audit — PASS for current Product 1.0

The known Work-producing classes have at least one legitimate closure path without a generic Work command bus:

| Condition class | Legitimate closure path |
|---|---|
| Readiness missing/conflict | new source evidence and/or Resolve/Clear ProductChannelCorrespondence |
| Listing authoring blocker | UpdateListingIntentDraft and/or current Readiness change |
| Listing/Price external ambiguity | owner D4 authoritative reread/reconciliation; new owner Intent only when a new decision is genuinely required |
| Availability divergence/staleness | Availability owner automatic reconciliation / owner policy/config change |
| Sales attribution ambiguity | ResolveSaleSellingEntityAttribution |
| Party ambiguity | ResolveBusinessSystemPartyResolution |
| Destination unsupported/unavailable | Materialization state becomes realized when safely established or remains explicit external-required; Work cannot fabricate success |
| Governance pending | CreateAuthorizationDecision |
| Fulfillment physical/deadline exception | Fulfillment checkpoint/current owner state and, when the problem becomes post-sale, PostSaleResolution |
| Shipment delivery exception | authoritative Shipment evolution and/or PostSaleResolution coordination |
| Economic attribution/divergence | ResolveEconomicAttribution / Economics re-evaluation |
| Post-Sale consequence | committed result from the owning consequence domain; any later human-only semantic decision gets the smallest Post-Sale-specific capability |

**Current result:** no proven Product 1.0 Work condition requires generic `SubmitWorkResolution`.

Reopen trigger: a real condition is closable only by human-held evidence, no admitted owner-specific capability can carry that evidence, and source auto-resolution is impossible. Then admit the smallest bounded Work→source evidence-submission capability without transferring source truth to Work.

## 9. W2-D negative controls

Later executable contract/conformance proof must reject or make unreachable at least:

1. PATCHing AuthorizationDecision into another historical decision;
2. generic `{entity_type, entity_id}` Authorization/Work target graphs;
3. authorization of a target revision different from the reviewed revision;
4. authorized scope wider than intended scope;
5. stale revoke leaving a targeted standing delegation alive solely because the client snapshot is old;
6. synthetic MPC Sale/Shipment/native business-order/fiscal IDs created only for normalization;
7. Product request containing TOP/NUNOTA/CODPARC/provider status as canonical MPC semantics;
8. PartyResolution create-new accepting arbitrary Customer master fields;
9. direct Destination realization write before D8 proof;
10. Product client creating BusinessOrderIntent/InvoicingIntent;
11. Fulfillment checkpoint operation hosted on Sales-owned resource instead of Fulfillment-owned identity;
12. physical conference recorded as automatic readiness when discrepancy was observed;
13. ordinary automation Principal establishing physical fact without explicit proven system capability;
14. global Artifact/ProductAsset identity created by Fulfillment artifact mechanics;
15. Post-Sale collapsing return/refund/cancellation into one mutually exclusive type;
16. direct Post-Sale close or Work close fabricating source/consequence resolution;
17. Work responsibility encoded as AccessRole/Permission;
18. generic Work resolution becoming a command bus for source-owner decisions;
19. `If-Match` being misapplied to a different referenced resource than the selected HTTP request resource.

## 10. Method outcome

**Parent D0→D5-B1 / ratified B2 / W1 / W2-A/B/C:** `CURRENT STRUCTURE CONFIRMED`.

W2-local identity crystallizations justified by real wire consumers:

- `AuthorizationDelegationId` — stable addressed identity now required for update/revoke while delegation attributes may change;
- `FulfillmentExecutionId` — stable Fulfillment-owned lifecycle identity required for physical checkpoint addressing and already prepared by D2's durable Fulfillment routing/dispatch-intent seam.

Neither introduces a new D1 domain or generic framework.

> **Operational resources are independently addressable only where their meaning/history genuinely needs identity; external facts remain source-qualified; checkpoints, consequence tracks, authorization, physical evidence and Work responsibility remain separate axes; no client can make a provider result, physical fact, approval or source resolution true by mutating a workflow status.**

No parent-stage reopen is required.

## 11. Exact next W2 work

**W2-E — transversal/final schema consistency.**

W2-E must close the cross-owner schema laws that should be decided once after A–D:

1. policy/config default/inheritance/explicit-override grammar for Availability, Commercial Economics and Fulfillment without a generic Rules/SLA platform;
2. capability/business outcome consistency and owner-local external-effect/convergence semantics without universal Result/Operation state;
3. referenced-resource precondition grammar for create/capability operations whose safety depends on another resource revision;
4. exact relationship among `ETag`/`If-Match`, `Idempotency-Key`, semantic prerequisites and referenced-resource validators;
5. RFC 9457 Problem Details extensions/problem types required by the admitted contract;
6. response/status/body rules for create/read/PATCH/`:verb`, including why valid business rejection/pending/ambiguity is not transport failure and why `202` is not a generic async marker;
7. final request/body closure and authority-field negative controls across owners;
8. final Whole-W2 consistency audit before independent review.

After W2-E converges, run a Whole-W2 coherence pass, prepare a disposable non-authoritative review candidate and invoke one Fable review following the canonical Standard Fable review workflow before operator ratification/consolidation.

Implementation remains blocked until D9.
