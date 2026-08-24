# NOTIF-01 D3-R1 — Presentation Feed-Forward Revalidation

> **Status:** PASS / NO COMMUNICATION-TOPOLOGY REOPEN REQUIRED
> **Accepted parents:** [D3-R Communication & Propagation Ratification](D6-R2-NOTIF-01-D3-R-RATIFICATION.md) + [D2-R4 Ratification](D6-R2-NOTIF-01-D2-R4-RATIFICATION.md)
> **Scope:** prove that the one accepted immutable `subject_display_label` can cross the existing source-owner committed-fact `E` boundary without creating a generic payload or moving authority
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Question

D2-R4 requires each materialized Notification to retain one immutable, notification-safe, human-readable `subject_display_label` supplied from source-owner meaning. D3-R already permits event payloads to contain stable immutable occurrence facts materially required by the consumer and forbids whole mutable aggregate mirrors or free-form payloads.

The bounded revalidation question is therefore:

> Can the existing fourteen source-owner committed-fact `E` contracts carry exactly one additional typed immutable presentation atom without changing communication form, ownership, recoverability or replay semantics?

**Result: YES / PASS.**

---

## 2. Revised minimum immutable occurrence contract

For NOTIF-01 only, each admitted source-owner `E` now preserves:

```text
organization_id
source_ref
source_occurrence_key
source_occurred_at
source_committed_at
subject_display_label
```

Audience-specific fields remain unchanged:

```text
DIRECT_SOURCE
  + recipient_principal_id

OWNER_DERIVED
  + recipient_principal_id

ORG_ROUTED
  + no recipient field
```

The two previously accepted preferred-awareness contracts may still carry only their bounded reverse replacement basis.

No other common field is admitted by symmetry.

---

## 3. Ownership and immutability law

The source owner produces `subject_display_label` from source-owned meaning at the committed occurrence boundary.

The field crossing `E` is:

- immutable for that occurrence;
- notification-safe and data-minimized under D2-R4;
- presentation-only;
- not a source identity, source-currentness claim, routing input, dedup key, authorization proof or command parameter.

Personal Notifications copies the supplied atom into the durable Notification presentation snapshot without reinterpreting business meaning.

Replay/redelivery of the same source occurrence must carry semantically equivalent label content for that occurrence. A producer must not mutate the label across redelivery merely because the current source title changed later.

---

## 4. Currentness Q remains independent

For the eleven current-attention families, accepted D3-R `E → Q` currentness reconciliation remains unchanged.

`subject_display_label` does not answer whether the occurrence is still relevant.

```text
E carries historical presentation atom
        ↓
owner Q proves STILL_RELEVANT / NO_LONGER_RELEVANT / UNKNOWN_OR_UNAVAILABLE
        ↓
only then may current awareness materialize where required
```

If current source presentation changed while the occurrence remains relevant, the Notification still retains the immutable occurrence-time presentation label. Navigation later rereads current source truth.

---

## 5. Recovery / replay / reorder

No new failure semantic is required:

- duplicate/replay of the same occurrence remains semantically idempotent under the accepted source occurrence key;
- delayed delivery preserves the original immutable label rather than regenerating from current source state;
- ORG_ROUTED historical route selection remains based on `source_committed_at`;
- D2-R2/D2-R3 eligibility epochs, route revisions and bounded supersession rules remain unchanged;
- preferred-first/generic-first replacement basis does not use the presentation label.

A recovery mechanism that can reconstruct the occurrence but cannot reconstruct the accepted immutable presentation atom is incomplete for this Notification consumer. D7 must preserve/recover the atom together with the source occurrence handoff or from another accepted durable occurrence authority without requiring source-per-row rereads as the Inbox baseline.

---

## 6. Negative controls

This revalidation fails if realization introduces any of the following:

```text
payload: {}
metadata: {}
template_variables: {}
source DTO snapshot
current-source presentation reread per Inbox row
generic event envelope as business authority
label-based dedup/routing/navigation/authorization
```

It also fails if the producer changes the occurrence label on replay based on later mutable source state.

---

## 7. Coherence result

D2-R4 is fully expressible inside the accepted D3 topology:

```text
source owner commits admitted occurrence
+ one immutable safe subject_display_label
        ↓ E
Personal Notifications
        ↓
accepted currentness/audience/replay rules
        ↓
durable Notification presentation snapshot
```

No new D1 edge, communication form, public event API, broker, generic payload, event identity or source authority is needed.

Therefore **D3-R remains ACCEPTED with this bounded feed-forward revalidation PASS**.

## 8. Gate

```text
D2-R4 presentation snapshot       ACCEPTED / OPERATOR-RATIFIED
D3-R1 presentation feed-forward  PASS / NO TOPOLOGY REOPEN
D5-R3 operation admission        OPEN / NEXT
canonical Product OAD            UNCHANGED
D6 / D7 / D8                     BLOCKED for NOTIF-01
Product implementation           BLOCKED UNTIL D9
```

**Exact next action:** adjudicate the frozen four-operation NOTIF-01 D5-R3 operation-admission table before any canonical Product OpenAPI edit.