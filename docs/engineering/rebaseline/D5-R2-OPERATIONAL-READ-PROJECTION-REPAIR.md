# D5-R2 — Operational Read Projection Repair

> **Status:** AUTHORIZED / OPERATOR-APPROVED FOR EXECUTION — proof pending
> **Trigger:** D6-R2 `OP-READ-01`
> **Parent authority:** accepted D5 API + W2/W3 + canonical Product OAD
> **Consumer evidence:** [D6-R2 P4-R1 Global IA / Operational Mass Reopen](D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Purpose

D6-R2 proved a real human operational-triage consumer that the accepted Product owners can answer semantically, but several current ListItem/query contracts cannot answer efficiently without N+1 point reads or frontend-authored workflow projection.

This amendment repairs **read projection/query expressibility only**. It does not create a new business capability, owner, lifecycle or operation.

## 2. Governing invariant

> **When an accepted owner already owns the state required by an evidenced collection consumer, the Product collection may expose the smallest semantic subset and typed owner-local narrowing needed to scan and triage that owner truth; the client must not reconstruct a new cross-owner workflow truth.**

The existing W2/W3 laws remain binding:

- ListItem fields are semantic subsets of the same owner point meaning;
- filters exist only for a real consumer and combine by W3 AND semantics;
- no generic projection/filter DSL;
- pagination does not prove total population or knowledge completeness;
- no caller-selectable sort or universal count by convenience.

## 3. Bounded repair

### 3.1 Business-System Materialization — BusinessOrderIntent

`BusinessOrderIntentListItem` additionally exposes required `convergence`, reusing the exact point-resource enum/meaning.

`ListBusinessOrderIntents` additionally admits optional owner-native filters:

```text
external_effect_state?
convergence?
```

No Party/Destination detail is copied into the collection merely for convenience.

### 3.2 Business-System Materialization — InvoicingIntent

`InvoicingIntent` additionally carries required source-qualified `sale` correlation.

`InvoicingIntentListItem` exposes:

```text
sale                       required
convergence                required
fulfillment_execution_id?  optional
```

`ListInvoicingIntents` additionally admits:

```text
external_effect_state?
convergence?
```

The Sale reference is correlation/navigation evidence only; Sales ownership does not move to Materialization.

### 3.3 Fulfillment

`FulfillmentExecutionListItem` additionally exposes the owner-native checkpoint representations already present on `FulfillmentExecution`:

```text
separation
physical_conference
packing
dispatch_handoff
provider_dispatch_deadline?
```

`ListFulfillmentExecutions` additionally admits typed owner-native narrowing:

```text
physical_readiness?                  ready | blocked | unknown
separation_state?                    pending | recorded
physical_conference_state?           pending | recorded
packing_state?                       pending | recorded
dispatch_handoff_state?              pending | recorded
provider_dispatch_deadline_before?   date-time
```

No `stage`, `next_action`, `priority`, severity or synthetic queue status is created.

### 3.4 Shipment

`ShipmentListItem` additionally exposes optional source-qualified `sale` and optional `dispatch_deadline`, both already owned by the point resource.

`ListShipments` additionally admits optional `state?` using the existing Shipment state vocabulary:

```text
pending | ready | dispatched | delivered | exception | unknown
```

No additional Shipment deadline filter is admitted now; the evidenced consumer does not yet require it.

### 3.5 No repair required

No D5-R2 change is required for MarketplaceSales, PostSaleResolution or OperationalWork. Their current collection contracts already satisfy the present bounded consumer or remain intentionally owner-local specialist views.

## 4. Product-surface conservation

D5-R2 must preserve exactly:

```text
Product operations       99
ordinary Permissions     30
Principal kinds          H / A / S only
semantic owners          unchanged
Product write surface    unchanged
Technical Ingress        unchanged
```

D5-R2 admits **zero** new Product operations and **zero** new Permissions.

## 5. Hard negative controls

D5-R2 explicitly rejects:

```text
/operational-dashboard
OperationalWorkflow owner
operational_stage
next_action
priority / severity / urgency score
total_count
kanban_column
generic filter DSL
cross-owner synthetic lifecycle
frontend N+1 detail fan-out as baseline
```

A user-facing phrase such as `A embalar` may be a frontend presentation of current Fulfillment predicates; it is not a new Product business state.

## 6. Proof and downstream revalidation

The canonical OAD remains the machine-readable wire authority. `scripts/verify-operational-read-contract.mjs` proves this amendment mechanically and includes negative controls.

After exact-head proof passes:

1. Product OAD remains 99/30/H-A-S and generated projections remain valid;
2. GF-02 is revalidated only for affected read/composition properties;
3. `OP-READ-01` may close;
4. D6-R2 returns to the corrected global-frame/B00 cycle before any dependent wireframe progresses.

No Product implementation is authorized by this amendment.
