# D5-R2 — Operational Read Projection Repair

> **Status:** ACCEPTED / CANONICAL / CURRENT CONSOLIDATED AUTHORITY  
> **Trigger:** D6-R2 `OP-READ-01`  
> **Parent authority:** current D5 + W2/W3 + canonical Product OAD  
> **Current consumer evidence:** [D6-R2 P5 Current Screen / Material-Surface Inventory](D6-R2-P5-SCREEN-SURFACE-INVENTORY.md) — operational homes R70–R101 and state grammar  
> **Current Product:** 106 operations / 31 ordinary Permissions / H-A-S; this repair's delta remains **0 operations / 0 Permissions**  
> **Implementation:** blocked until D9

## 1. Purpose

D6-R2 proved a real human operational-triage consumer that accepted Product owners could answer semantically, while several earlier collection contracts would have required N+1 point reads or frontend-authored workflow conclusions.

D5-R2 repairs **owner-local read projection/query expressibility only**. It creates no new business capability, owner, lifecycle, Product operation or write authority.

## 2. Governing invariant

> **When an accepted owner already owns the state required by an evidenced collection consumer, the Product collection may expose the smallest semantic subset and typed owner-local narrowing needed to scan and triage that owner truth; the client must not reconstruct a new cross-owner workflow truth.**

Binding W2/W3 laws remain:

- ListItem fields are semantic subsets of the same owner point meaning;
- filters exist only for a real consumer and combine under typed owner-local query semantics;
- pagination does not prove total population or knowledge completeness;
- no generic projection/filter DSL, caller-selectable sort or universal count;
- a production baseline must not use N+1 detail fan-out to compensate for a materially deficient owner collection.

## 3. Accepted bounded repair

### 3.1 Business-System Materialization — BusinessOrderIntent

`BusinessOrderIntentListItem` exposes owner-native `convergence`.

`ListBusinessOrderIntents` admits optional:

```text
external_effect_state?
convergence?
```

No Party/Destination detail is copied into the collection merely for convenience.

### 3.2 Business-System Materialization — InvoicingIntent

`InvoicingIntent` carries source-qualified `sale` correlation.

`InvoicingIntentListItem` exposes:

```text
sale                       required
convergence                required
fulfillment_execution_id?  optional
```

`ListInvoicingIntents` admits:

```text
external_effect_state?
convergence?
```

Sale correlation/navigation does not transfer Sales authority to Materialization.

### 3.3 Fulfillment

`FulfillmentExecutionListItem` exposes owner-native checkpoint meaning already present on the point owner:

```text
separation
physical_conference
packing
dispatch_handoff
provider_dispatch_deadline?
```

`ListFulfillmentExecutions` admits typed owner-native narrowing:

```text
physical_readiness?                  ready | blocked | unknown
separation_state?                    pending | recorded
physical_conference_state?           pending | recorded
packing_state?                       pending | recorded
dispatch_handoff_state?              pending | recorded
provider_dispatch_deadline_before?   date-time
```

No synthetic `stage`, `next_action`, `priority`, severity or queue status is created.

### 3.4 Shipment

`ShipmentListItem` exposes optional source-qualified `sale` and optional `dispatch_deadline`, both already owner point meaning.

`ListShipments` admits optional owner-native:

```text
state? = pending | ready | dispatched | delivered | exception | unknown
```

No additional Shipment deadline filter is admitted without another real consumer.

### 3.5 No repair required

MarketplaceSales, PostSaleResolution and OperationalWork did not require D5-R2 schema/query expansion. Their accepted owner-local collections remain authoritative for their jobs.

## 4. Product-surface conservation

D5-R2 itself adds:

```text
new Product operations       0
new ordinary Permissions     0
new Principal kinds          0
new semantic owners          0
new Product writes           0
Technical Ingress changes    0
```

The global Product later evolved through separate accepted increments to **106/31/H-A-S**. Those later additions do not change D5-R2's zero-surface-delta meaning.

## 5. Hard negative controls

Reject:

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
frontend N+1 detail fan-out as production baseline
```

A human phrase such as `A embalar` may present current Fulfillment predicates; it is not a new Product business state.

## 6. Executable proof

Canonical OAD remains machine-readable wire authority. `scripts/verify-operational-read-contract.mjs` is the targeted executable proof for these owner-local projection/filter invariants and their negative controls.

Historical ratification/CI checkpoints remain in Git history; they are not parallel current status authority. The current aggregate Product proof/gate protects the canonical Product contract under current 106/31 authority.

## 7. Downstream disposition

The required GF-02 operational-read revalidation is already accepted as D8-R2. `OP-READ-01` is closed.

Current frontend consumer homes and state grammar now live in D6-R2 P5. Future operational UX must preserve owner separation and may reopen D5 only when new human evidence proves an existing owner collection materially insufficient.
