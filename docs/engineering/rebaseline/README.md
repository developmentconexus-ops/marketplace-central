# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 CANONICAL; Technical Ingress A+B ACCEPTED IN-STAGE; Whole-Ingress lead review COMPLETE / RESTRUCTURE INGRESS-LOCAL; IG-G1…G6 operator-ratified for independent challenge only; Fable Whole-Ingress independent review = NEXT**  
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

`D5-B2-TECHNICAL-INGRESS-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` is **NON-AUTHORITATIVE lead review evidence** and is deliberately outside the authority path. The operator ratified IG-G1…G6 on 2026-08-19 **only as correction direction for independent challenge**. That ratification does not yet amend accepted Technical Ingress A+B.

`AI-DIALOG.md` is the non-authoritative working review channel. Cockpit remains non-authoritative and is synchronized only after canonical status changes.

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
        IG-G1…G6 lead direction                          RATIFIED FOR INDEPENDENT CHALLENGE
        Fable Whole-Ingress independent review           NEXT
        GPT adjudication / final convergence             BLOCKED BY FABLE REVIEW
      Final Problem/media consistency                    BLOCKED BY INGRESS
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

## 3. Accepted Technical Ingress authority that remains current during review

### Acquisition-A

- provider-specific inbound transport adapter → MPC-native recoverable acquisition ingress → closed typed acquisition request → authoritative D4 reread → consumer-owned translation → owner commit;
- seven accepted acquisition families: Listing, Price, Sale, Shipment, Payment, Post-Sale Claim, Competitive Position;
- exact namespace binding determines Organization; provider payload never does;
- positive acknowledgement means technical custody/quarantine only;
- delivery dedup/coalescing never replaces owner semantic idempotency;
- no generic Webhook/ExternalEvent/Integration authority.

### OAuth-B

- begin requires current authenticated human + Principal eligibility + Membership + `portfolio.manage` + target Installation;
- one current server-bound Authorization Attempt per Installation/provider app;
- opaque/single-use/finite-lived state; current Mercado Livre lane uses Authorization Code + PKCE where supported;
- callback revalidates current initiator authority and Installation eligibility;
- provider seller/account identity is proven authoritatively;
- initial bind or same-seller reauthorization only; different seller fails closed;
- complete credential generation activation; stale refresh/older attempt cannot overwrite newer generation;
- OAuth remains outside acquisition inbox and Product OpenAPI/SDK business operation surface.

## 4. Ratified Whole-Ingress lead direction — REVIEW INPUT, NOT YET CANONICAL

The operator ratified these six Technical-Ingress-local corrections for independent challenge:

### IG-G1 — signal disposition / quarantine scope

```text
unverified/malformed
→ protocol reject; no MPC custody

verified but non-admitted topic/resource
→ terminal technical non-processing; provider ack may be positive as protocol requires; no quarantine/future-work obligation

admitted + exact namespace binding
→ recoverable attributed custody

admitted + missing/ambiguous/contradictory binding
→ bounded technical quarantine
```

### IG-G2 — namespace attribution is independent from active posture and credential usability

A retained unambiguous Installation↔seller/account namespace binding may attribute late historical Payment/Post-Sale/Shipment signals after business deactivation or while credentials are auth-invalid. Auth-blocked attributed acquisition remains Organization-scoped recoverable work, not unbound quarantine. Deactivation does not erase historical namespace correlation.

### IG-G3 — path convergence does not invent discovery/recovery capability

When D4 admits multiple discovery mechanisms for one acquisition family, they converge on the same typed path. The shared path never promises polling/enumeration/recovery/completeness for a family whose provider contract does not establish it.

### IG-G4 — successful OAuth activation awakens technical recovery

Successful initial/same-seller credential activation may awaken Installation-scoped capability revalidation, bounded bootstrap and pending attributed acquisition/reconciliation blocked on auth. This is technical recovery, not Product `Sync`/`Refresh`, D3 domain event or owner-state proof.

### IG-G5 — durable non-secret authorization/binding lineage

Preserve the smallest durable trust-boundary provenance needed to explain Installation/provider app/initiating Principal/proven seller, initial-bind vs same-seller reauthorization/supersession/failure, generation correlation and material timestamps. Never preserve authorization code, client secret, access token, refresh token or PKCE verifier in that lineage.

### IG-G6 — Product-authenticated technical OAuth begin is not a 96th Product operation

OAuth begin consumes current D2/W4 access facts (`H`, Principal eligibility, explicit Organization Membership, `portfolio.manage`, exact Installation) as a technical initiation gate but remains outside the 95 Product operations/Product OpenAPI. Callback/acquisition routes use provider protocol trust. Technical routes are collision-separated from Product roots; exact technical prefix/host is deferred to final ingress wire closure.

These statements are not canonical amendments yet.

## 5. Prohibited now

Until independent Whole-Ingress review is adjudicated:

- do not silently apply IG-G1…G6 to accepted A+B;
- do not begin final Problem/media consistency, Product OpenAPI/tooling, D6–D9 or implementation;
- do not create generic Integration/Ingress/Webhook/ExternalEvent/OAuth authority;
- do not treat Installation deactivation/auth failure as loss of historical namespace attribution;
- do not quarantine unsupported topics as future-feature backlog;
- do not invent generic scan/recovery capability;
- do not add OAuth begin as a Product operation merely to reuse W4 middleware;
- do not select D7 queue/storage/credential/audit implementation.

## 6. Exact next action

**Run one coherent independent Fable Whole-Ingress review over accepted Technical Ingress A+B plus the operator-ratified-for-challenge IG-G1…G6 direction.**

The reviewer must reconstruct authority independently and challenge the complete external trust-boundary system, including at minimum:

1. provider signal verification, closed admission and quarantine boundary;
2. historical namespace attribution versus Installation active/auth posture;
3. seven typed acquisition families and Mercado Livre ADMIT/DEFER/REJECT matrix;
4. recoverable custody/ack semantics without choosing D7 mechanism;
5. push/poll/recovery path convergence without fabricated coverage;
6. OAuth begin trust: H + current W4 facts without becoming Product operation;
7. state/PKCE/single-use/supersession/callback replay;
8. callback-time current-authority revalidation;
9. initial bind versus same-seller reauthorization and different-seller mismatch;
10. credential-generation concurrency/stale-refresh protection;
11. non-secret durable binding lineage and secret minimization;
12. OAuth-success technical recovery wake-up without Sync/domain-event authority;
13. Product OpenAPI/SDK separation and technical route collision separation;
14. new findings beyond IG-G1…G6 and smallest parent reopen only if materially unavoidable.

Reviewer output is evidence only. GPT adjudicates every material finding. Round 2 occurs only if a real material contradiction survives.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4 are canonical;
- Technical Ingress A+B are accepted in-stage/operator-ratified;
- Whole-Ingress lead review found IG-G1…G6 and `RESTRUCTURE INGRESS-LOCAL`;
- IG-G1…G6 were operator-ratified **only as direction for independent challenge**;
- review candidate and `AI-DIALOG.md` remain non-authoritative evidence;
- **one coherent Fable Whole-Ingress independent review is the exact next action**;
- no parent reopen is currently proven;
- final Problem/media/OpenAPI/D6–D9/implementation remain blocked.

If not, the active authority tree is inconsistent.
