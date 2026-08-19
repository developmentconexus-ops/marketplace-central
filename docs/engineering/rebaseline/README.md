# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 CANONICAL; Technical Ingress A+B ACCEPTED IN-STAGE; lead Whole-Ingress review COMPLETE / RESTRUCTURE INGRESS-LOCAL; operator ratification of IG-G1…G6 = NEXT**  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Last updated:** 2026-08-19

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this router
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md`
6. `docs/architecture/decisions/README.md`
7. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
8. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
9. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
10. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
11. `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
12. `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
13. `docs/engineering/rebaseline/D5-API.md`
14. `docs/engineering/rebaseline/D5-B2-PRODUCT-OPERATION-SURFACE.md`
15. `docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md`
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md` — canonical W1
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical W2
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — canonical W3
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — canonical W4
20. `docs/engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md` — accepted Technical Ingress A+B design home
21. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
22. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in accepted artifacts.

`D5-B2-TECHNICAL-INGRESS-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` is **NON-AUTHORITATIVE lead review evidence** and is deliberately outside the authority path. IG-G1…G6 do not modify accepted Technical Ingress A+B until operator ratification and later canonical filing.

`AI-DIALOG.md` remains protocol-only with no active review cycle. Cockpit remains non-authoritative and is synchronized only after canonical status changes.

---

## 2. Program state

```text
D0 — Product / System Definition                         CLOSED / ACCEPTED
D1 — Domains / Boundaries                                CLOSED / ACCEPTED
D2 — Identity / Tenant / Data Ownership                  CLOSED / ACCEPTED
D3 — Communication / Events                              CLOSED / ACCEPTED
D4 — External Integrations                               CLOSED / ACCEPTED AS A WHOLE
D4-R1 — Publication Input & Listing Authoring            ACCEPTED / CANONICAL
Decision Reconciliation                                  ACCEPTED / CANONICAL
D5 — API                                                  OPEN / ACTIVE
  B1 Semantic API Model                                  ACCEPTED / CANONICAL
  B2 Product Operation / Resource Surface                 OPEN / ACTIVE
    Operation Admission Matrix                           ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                        ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar                  ACCEPTED / CANONICAL
      W2 Request / Response Schema Grammar               ACCEPTED / CANONICAL
      W3 Collections / Query / Cursor Grammar            ACCEPTED / CANONICAL
      W4 Permission / Client-Class Enforcement           ACCEPTED / CANONICAL
      Technical non-Product ingress                      OPEN / ACTIVE
        A External Acquisition Ingress                   ACCEPTED IN-STAGE / OPERATOR-RATIFIED
        B OAuth / Authorization Ceremony                 ACCEPTED IN-STAGE / OPERATOR-RATIFIED
        Whole-Ingress lead coherence                     COMPLETE / RESTRUCTURE INGRESS-LOCAL
        IG-G1…G6 operator direction                      NEXT
        Fable Whole-Ingress independent review           BLOCKED BY OPERATOR DIRECTION
      Final Problem/media consistency                    BLOCKED BY INGRESS
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Accepted Technical Ingress authority that remains current during review

### Acquisition-A

- provider-specific inbound transport adapter → MPC-native recoverable acquisition ingress → closed typed acquisition request → authoritative D4 reread → consumer-owned translation → owner commit;
- seven accepted acquisition families: Listing, Price, Sale, Shipment, Payment, Post-Sale Claim, Competitive Position;
- exact namespace binding determines Organization; provider payload never does;
- positive acknowledgement means technical custody/quarantine only;
- push/poll/recovery converge on the same typed path where provider contracts admit them;
- delivery dedup/coalescing never replaces owner semantic idempotency;
- quarantine is technical and cannot dispatch D1 work before attribution;
- no generic Webhook/ExternalEvent/Integration authority.

### OAuth-B

- begin requires current authenticated human + Principal eligibility + Membership + `portfolio.manage` + target Installation;
- one current server-bound Authorization Attempt per Installation/provider app;
- opaque/single-use/finite-lived `state`; selected current Mercado Livre lane uses Authorization Code + PKCE where supported;
- callback revalidates current initiator authority and Installation eligibility;
- provider seller/account identity is proven authoritatively;
- initial bind or same-seller reauthorization only; different seller fails closed;
- complete credential generation activation; stale refresh/older attempt cannot overwrite newer generation;
- OAuth remains outside acquisition inbox and Product OpenAPI/SDK business operation surface.

---

## 4. Lead Whole-Ingress findings — NON-AUTHORITATIVE UNTIL OPERATOR DECISION

Lead review found six Technical-Ingress-local corrections and **no parent-stage semantic reopen**:

### IG-G1 — signal disposition / quarantine scope

Distinguish:

```text
unverified/malformed → protocol reject; no custody
verified but non-admitted topic → terminal technical non-processing; may ack per provider contract; no quarantine obligation
admitted + exact binding → recoverable attributed custody
admitted + unresolved/ambiguous binding → bounded technical quarantine
```

Unknown/DEFER/REJECT topics never become latent future work merely because they arrived.

### IG-G2 — namespace attribution != business activation != credential availability

Late Payment/Post-Sale/Shipment evidence may still belong to a deactivated or temporarily auth-invalid Installation. Stable seller/account namespace correlation may attribute the signal even when current business participation/credentials are unavailable. Auth-unavailable attributed acquisitions remain Organization-scoped recoverable work, not unbound quarantine.

### IG-G3 — path convergence != discovery-source equivalence

Push/poll/recovery use the same typed path **only where D4 actually admits each discovery source**. Shared handlers never invent enumeration/recovery/completeness capability for Payment/Claim/competition or another family.

### IG-G4 — OAuth success awakens technical recovery

Successful initial/same-seller credential activation may awaken Installation-scoped capability revalidation, bounded bootstrap and pending acquisition recovery. This is technical recovery, not Product `Sync`/`Refresh`, D3 domain event or proof that owner meaning changed.

### IG-G5 — durable non-secret OAuth binding lineage

Preserve the smallest non-secret trust-boundary provenance needed to explain who initiated the binding, which Installation/provider application/seller were involved, whether it was initial bind/reauthorization/supersession/failure, generation correlation and material timestamps. Never preserve authorization codes/client secrets/access tokens/refresh tokens/PKCE verifier in that lineage.

### IG-G6 — Product-authenticated technical begin is not a 96th Product operation

OAuth begin consumes existing current D2/W4 access facts (`H`, Principal eligibility, explicit Organization Membership, `portfolio.manage`, exact Installation) as a technical initiation gate, but remains outside the 95 Product operations and Product OpenAPI. Callback/acquisition routes use provider protocol trust. Technical routes must be collision-separated from Product roots; exact technical prefix/host spelling is deferred to final ingress wire closure.

---

## 5. Prohibited now

Until operator ratifies/revises IG-G1…G6:

- do not silently modify accepted Technical Ingress A+B;
- do not open Fable review yet;
- do not begin final Problem/media consistency, Product OpenAPI/tooling, D6–D9 or implementation;
- do not create generic Integration/Ingress/Webhook/ExternalEvent/OAuth authority;
- do not treat Installation deactivation/auth failure as loss of historical namespace attribution by convenience;
- do not quarantine unsupported topics as future-feature backlog;
- do not invent generic scan/recovery capability;
- do not add OAuth begin as a Product operation merely to reuse W4 middleware;
- do not select D7 queue/storage/credential/audit implementation.

---

## 6. Exact next action

**Operator reviews and ratifies/revises IG-G1…IG-G6 as the Whole-Ingress lead correction direction.**

If ratified, because Technical Ingress crosses provider authentication, tenant attribution, recoverable custody and credential trust boundaries, run **one coherent independent Fable Whole-Ingress review** over A+B+G1…G6 before canonical filing.

After independent review/adjudication and final operator ratification:

1. consolidate corrections into the single `D5-B2-TECHNICAL-INGRESS.md` authority;
2. remove review candidate; Git history remains archive;
3. reset review channel;
4. update cockpit;
5. advance to final Problem/media consistency;
6. then decide Product OpenAPI authority/tooling/minor version.

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4 are canonical;
- Technical Ingress A+B are accepted in-stage/operator-ratified;
- Whole-Ingress lead review found only IG-G1…G6 and `RESTRUCTURE INGRESS-LOCAL`;
- review candidate is evidence only;
- **operator ratification/revision of IG-G1…G6 is the exact next action**;
- no parent reopen is currently proven;
- later Wire/D6–D9/implementation remain blocked.

If not, the active authority tree is inconsistent.
