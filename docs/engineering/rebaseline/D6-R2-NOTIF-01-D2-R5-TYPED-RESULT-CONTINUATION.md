# NOTIF-01 D2-R5 — Typed Result & Requester Continuation Repair

> **Status:** OPEN — TARGETED D2 REOPEN / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** [D5-F3 Global Maximum Operation / Permission Review](D6-R2-NOTIF-01-D5-F3-GLOBAL-MAXIMUM-OPERATION-PERMISSION-REVIEW.md)
> **Accepted base:** D2-R + D2-R2 + D2-R3 + D2-R4 remain ACCEPTED / OPERATOR-RATIFIED; unaffected clauses remain unchanged
> **Scope law:** repair only the immutable result presentation of F02/F14 and the requester continuation identity of F14
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Why D2 reopens narrowly

The operator-approved semantic family contracts require two result-bearing Notifications to communicate an immutable result, not merely that “something happened”:

```text
F02 OFFERING_ASYNC_ACTION_RESULT
→ converged | rejected | ambiguous | divergent

F14 AUTHORIZATION_DECISION_RESULT
→ authorize | reject
```

Accepted D2-R4 adds only `subject_display_label`, which identifies the human subject but intentionally carries no result/status meaning. Therefore `NotificationKind + subject_display_label` is insufficient to render the already-approved human meaning of F02/F14 without a source reread.

D5-F3 also proves that F14's current `AuthorizationDecisionRef` source continuation is authority-centric rather than requester-centric: canonical Product `GetAuthorizationDecision` requires `governance.read`, which the exact requester/initiator need not possess.

This repair adds only the smallest typed immutable result state and corrects that one continuation reference.

---

## 2. Two typed result atoms — CANDIDATE

Canonical Notification state gains two **kind-constrained optional fields**:

```text
offering_async_result_outcome?
authorization_decision_outcome?
```

They are not generic metadata. Each is legal/required only for its exact NotificationKind.

### 2.1 `offering_async_result_outcome`

Required iff:

```text
kind = OFFERING_ASYNC_ACTION_RESULT
```

Closed values:

```text
converged
rejected
ambiguous
divergent
```

Meaning is the immutable Offering-owned material outcome that caused this exact Notification occurrence.

This lets Product/client presentation distinguish, for example:

```text
converged → “Alteração de preço concluída”
rejected  → “Alteração de preço rejeitada”
ambiguous → “Resultado da alteração de preço está incerto”
divergent → “Alteração de preço não convergiu”
```

without rereading current mutable Offering state solely to know the historical result occurrence.

### 2.2 `authorization_decision_outcome`

Required iff:

```text
kind = AUTHORIZATION_DECISION_RESULT
```

Closed values reuse canonical Governance decision meaning:

```text
authorize
reject
```

This lets Product/client presentation distinguish the operator-approved human messages:

```text
authorize → “Sua solicitação foi aprovada”
reject    → “Sua solicitação foi rejeitada”
```

The field is the immutable decision disposition for this exact requester-result occurrence. It is not current target execution state and does not claim the approved action later executed/converged.

---

## 3. Result-state ownership and immutability — CANDIDATE

For both fields:

- source owner owns the result meaning;
- Personal Notifications retains the immutable typed value after materialization;
- replay/redelivery of the same `source_occurrence_key` must carry the same result value;
- result value is never recomputed from later current source state;
- the result is presentation/history meaning only and never drives source mutation, authorization, routing or deduplication;
- a later distinct source outcome that legitimately creates another Notification is a new occurrence with its own result atom.

No generic `result:string`, `reason:string`, `status`, `summary`, payload or template variables are admitted.

---

## 4. F14 requester continuation correction — CANDIDATE

Current accepted D2-R uses:

```text
AUTHORIZATION_DECISION_RESULT
→ AuthorizationDecisionRef
```

D5-F3 proves this is not the smallest truthful continuation for the exact requester/initiator consumer.

The corrected compatibility is:

```text
AUTHORIZATION_DECISION_RESULT
→ AuthorizationTargetRef
```

where the already-accepted closed target union remains:

```text
AuthorizationTargetRef =
    ListingIntentRef
  | PriceIntentRef
  | BusinessOrderIntentRef
  | InvoicingIntentRef
```

Rationale:

- the requester initiated/owns a job in the action-owning target workspace;
- `GetAuthorizationDecision` requires `governance.read`, which the requester need not have;
- Governance read authority must not be broadened merely to repair Notification navigation;
- the target is the natural continuation where the requester can see the action's current authorized/pending/rejected progression under current target access;
- Governance still owns the decision occurrence and its audit/history.

The Notification's `source_occurrence_key` remains Governance-owned and distinguishes the exact committed Authorization Decision occurrence. The canonical AuthorizationDecision identity may supply that stable occurrence meaning without becoming the Notification's navigation/source subject.

---

## 5. Source-union consequence — CANDIDATE

After the F14 correction, no admitted Launch-V1 NotificationKind requires `AuthorizationDecisionRef` as its `source_ref` variant.

Therefore the closed Notification source union removes that unused variant rather than preserving it by symmetry.

The Governance-related variants become:

```text
AuthorizationTargetRef
```

used by:

```text
AUTHORIZATION_ACTION_REQUIRED
AUTHORIZATION_DECISION_RESULT
```

with distinct kinds and occurrence keys preserving distinct semantics.

Governance canonical `AuthorizationDecisionID` remains valid Governance identity; it simply no longer needs to be a Personal Notifications source-subject variant.

---

## 6. Revised relevant canonical Notification state — CANDIDATE

Conceptually:

```text
notification_id
organization_id
recipient_principal_id
kind
source_ref
source_occurrence_key
source_occurred_at
source_committed_at
subject_display_label
offering_async_result_outcome?      // F02 only; required for F02
authorization_decision_outcome?     // F14 only; required for F14
created_at
read_at?
archived_at?
superseded_at?
superseded_by_notification_id?
supersession_reason?
work_replacement_basis?
sale_attention_replacement_basis?
revision
```

Invalid state examples:

```text
NEW_MARKETPLACE_SALE + offering_async_result_outcome
WORK_ASSIGNMENT + authorization_decision_outcome
OFFERING_ASYNC_ACTION_RESULT without offering_async_result_outcome
AUTHORIZATION_DECISION_RESULT without authorization_decision_outcome
AUTHORIZATION_DECISION_RESULT + AuthorizationDecisionRef source_ref
```

---

## 7. Access and navigation law — CANDIDATE

The typed result atoms do not grant source access.

For F14:

```text
Notification
→ AuthorizationTargetRef
→ current target Product read
→ current target Permission/access rechecked
```

The requester may learn the immutable authorization disposition from their own Notification while current target navigation remains separately authorized.

A requester is not granted `governance.read` merely because Governance produced the decision result.

If current target access was later revoked, the Notification can still preserve the non-sensitive immutable decision result/history while target navigation fails closed, exactly like accepted Notification history after source access loss.

---

## 8. D3 feed-forward obligation — CANDIDATE

If D2-R5 is accepted, D3 needs only bounded revalidation:

```text
OfferingAsyncActionResultCommitted E
+ offering_async_result_outcome

AuthorizationDecisionResultCommitted E
+ AuthorizationTargetRef
+ authorization_decision_outcome
```

All existing Organization scope, occurrence identity, `source_committed_at`, currentness profile, recipient derivation, recoverability, replay and suppression semantics remain unchanged.

The F14 event continues to be Governance-owned committed-fact `E`; carrying the governed target reference does not transfer decision authority to the target owner or Personal Notifications.

No new communication form or generic envelope is justified.

---

## 9. Negative controls — CANDIDATE

This repair fails if it:

1. adds a generic `result`, `reason`, `status`, metadata or payload field;
2. stores final localized sentences as source truth;
3. permits result fields on unrelated NotificationKinds;
4. uses result fields for routing/dedup/authorization/source mutation;
5. preserves `AuthorizationDecisionRef` in the Notification union after no admitted kind uses it merely for symmetry;
6. grants requester `governance.read` as a navigation workaround;
7. treats Governance `authorize` as proof that the target action later executed/converged;
8. changes any of the fourteen family meanings, audiences or suppression rules;
9. selects D5 path/schema spelling, D7 storage or D6 rendering mechanics inside D2.

---

## 10. Global-Maximum result

This repair closes the two remaining semantic/presentation holes found by the operator-requested Global Maximum review without adding operations or Permissions by mechanism.

```text
new Product operations through D2-R5     0
new ordinary Permissions through D2-R5   0
new typed immutable result atoms          2
source-ref variants                       -1 unused AuthorizationDecisionRef
```

Counts remain non-goals.

## 11. Gate

```text
D2-R / R2 / R3 / R4              ACCEPTED / OPERATOR-RATIFIED
D5-F3 Global Maximum review       COMPLETE / FALSIFIERS FOUND
D2-R5 typed result/continuation   READY FOR OPERATOR REVIEW
D3-R2 bounded feed-forward        BLOCKED BY D2-R5
D5-R3 final table                 BLOCKED BY D2-R5/D3-R2
canonical Product OAD             UNCHANGED — 99/30
D6 / D7 / D8                      BLOCKED for NOTIF-01
Product implementation            BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this D2-R5 repair. Do not grant extra Permissions or edit the canonical Product OpenAPI as a workaround.