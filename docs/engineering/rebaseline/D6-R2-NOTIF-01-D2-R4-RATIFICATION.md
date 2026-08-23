# NOTIF-01 D2-R4 — Notification Presentation Snapshot Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED
> **Accepted artifact:** [D2-R4 Notification Presentation Snapshot](D6-R2-NOTIF-01-D2-R4-PRESENTATION-SNAPSHOT.md), blob `e2c0c073223d05a7d716d6f9ee673e370297e779`
> **Accepted base:** D2-R + D2-R2 + D2-R3 remain ACCEPTED / OPERATOR-RATIFIED
> **Triggering consumer:** [D5-F2 Inbox Presentation Context Finding](D6-R2-NOTIF-01-D5-F2-INBOX-PRESENTATION-CONTEXT-FINDING.md)
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## Ratified result

The operator ratified only the bounded immutable human-presentation snapshot required for a useful personal Inbox:

1. canonical Notification state gains one non-empty bounded `subject_display_label`;
2. the source owner owns the business-derived meaning used to produce the label; Personal Notifications owns only the retained immutable snapshot after materialization;
3. the label is presentation-only and never identity, routing, deduplication, source-currentness, authorization, navigation or machine-automation authority;
4. final localized Notification title/copy remains derived from `NotificationKind`; no second summary/detail field, generic payload, metadata, template variables or template engine is admitted;
5. the label is immutable for historical awareness even if current source presentation later changes;
6. the label must be notification-safe and data-minimized: buyer/customer PII, address, payment/fiscal detail, credentials, provider raw payload and arbitrary source JSON remain excluded by default;
7. source access remains re-authorized on navigation through canonical typed `source_ref`; possession of the retained label grants nothing.

No D3 event naming, D5 wire spelling, D6 component layout or D7 persistence mechanism is accepted through this ratification.

## Gate

```text
D0-R / D1-R                    ACCEPTED
D2-R / D2-R2 / D2-R3          ACCEPTED
D2-R4 presentation snapshot   ACCEPTED / OPERATOR-RATIFIED
D3 presentation feed-forward  REVALIDATION NEXT
D5-R3 operation table         BLOCKED until feed-forward PASS
canonical Product OAD         UNCHANGED
Product implementation        BLOCKED UNTIL D9
```

**Exact next action:** revalidate only the already-accepted D3-R occurrence contract for propagation of the one approved immutable `subject_display_label`; do not reopen communication form or edit the Product OAD first.