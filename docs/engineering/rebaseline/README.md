# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 CANONICAL; Technical Ingress A+B ACCEPTED IN-STAGE; Whole-Ingress Fable review COMPLETE; GPT adjudication CONVERGED / RESTRUCTURE TECHNICAL-INGRESS-LOCAL; operator final Whole-Ingress ratification = NEXT**  
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

`D5-B2-TECHNICAL-INGRESS-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` are **NON-AUTHORITATIVE review evidence**. The converged corrections below do not amend Technical Ingress or canonical parents until final operator ratification and filing.

Cockpit remains non-authoritative and is synchronized only after canonical status changes.

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
        Whole-Ingress lead review                        COMPLETE / RESTRUCTURE INGRESS-LOCAL
        Fable Whole-Ingress independent review           COMPLETE
        GPT final adjudication                           CONVERGED
        operator final Whole-Ingress ratification        NEXT
      Final Problem/media consistency                    BLOCKED BY INGRESS
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Accepted Technical Ingress authority that remains current until final ratification

### Acquisition-A

- provider-specific inbound transport adapter → MPC-native recoverable acquisition ingress → closed typed acquisition request → authoritative D4 reread → consumer-owned translation → owner commit;
- seven accepted marketplace acquisition families: Listing, Price, Sale, Shipment, Payment, Post-Sale Claim, Competitive Position;
- provider signal is neither D3 domain event nor current provider truth;
- exact MarketplaceInstallation namespace binding determines Organization; provider payload never does;
- positive acknowledgement means technical custody/quarantine only under the current accepted text;
- delivery dedup/coalescing never replaces owner semantic correctness;
- no generic Webhook/ExternalEvent/Integration authority.

### OAuth-B

- technical begin uses current authenticated human + Principal eligibility + Membership + `portfolio.manage` + exact Installation;
- one current server-bound Authorization Attempt per Installation/provider app;
- opaque/single-use/finite-lived state; current ML lane selected Authorization Code + PKCE where supported;
- callback revalidates current initiator authority and Installation eligibility;
- initial bind or same-seller reauthorization only; different seller fails closed;
- complete credential generation activation; stale generation cannot overwrite newer generation;
- OAuth stays outside acquisition inbox and Product OpenAPI/SDK business operation surface.

---

## 4. Converged Whole-Ingress corrections — REVIEW RESULT, NOT YET CANONICAL

Independent review materially strengthened the package. GPT adjudication accepted the load-bearing findings; Round 2 is **NOT REQUIRED**. No D0→W4/D3/D4 semantic reopen is required.

### TI-C1 — deactivation preserves evidence-recovery capability

MarketplaceInstallation deactivation removes business participation authority but does not erase retained seller/account namespace attribution or historical-evidence recovery capability.

A deactivated Installation with an unambiguous retained seller binding may undergo **same-seller technical reauthorization for evidence-recovery only**, under current human/access checks. This is not business reactivation, creates no Product operation and authorizes no publication/write effect. Attributed acquisition blocked on credentials must remain detectable/recoverable rather than silently pending forever.

### TI-C2 — quarantine is much narrower

```text
A. malformed / provider-protocol verification failed
   → reject; no custody

B. protocol-admissible but non-admitted topic/resource,
   or admitted-looking signal with no plausible pending legitimate bind
   → explicit terminal technical non-processing

C. admitted + exact retained/current namespace binding
   → recoverable Organization-attributed custody
   → may remain auth-blocked until credentials are restored

D. admitted + unresolved binding
   → bounded platform quarantine ONLY when a legitimate binding is
      plausibly pending for a known current Installation/provider app
```

Multiple contradictory bindings, arbitrary unknown seller identities and unsupported future topics never become indefinite quarantine backlog. Quarantine capacity must be bounded; overflow is counted/observable and refused/fails honestly, never silently discarded.

For the current Mercado Livre HTTP notification lane, current official provider evidence establishes a real protocol-origin verification basis (HTTPS plus provider-documented notification-source validation). Concrete provider IPs remain adapter-local/revalidated facts rather than MPC constants.

### TI-C3 — transaction-bound OAuth anti-injection / CSRF protection

Initial seller binding may not rely on `state` secrecy alone.

Every provider authorization lane requires one transaction-bound anti-CSRF / authorization-code-injection control tied to the initiating transaction/user-agent.

For the current Mercado Livre lane, PKCE is selected and may be load-bearing for that proof when enabled/supported; `state` remains opaque/single-use/finite-lived MPC correlation. For a provider without usable PKCE, `state` must be securely bound to the initiating user-agent session or equivalent explicit human confirmation must occur before initial bind activation.

### TI-C4 — marketplace-only Product 1.0 acquisition seam

Remove the unused `SourceInstance` branch from native External Acquisition Ingress. Product 1.0 ingress is MarketplaceInstallation-scoped. Sankhya remains on accepted embedded/outbound acquisition paths; a future business-system callback must be admitted explicitly with its own families.

### TI-C5 — routing ownership cleanup

D4 remains authority for provider protocol/auth/source semantics. Technical Ingress is the D5-B2 wire/trust-boundary crystallization of that D4 protocol surface. At final filing, W1/Operation Matrix routing text is corrected to cross-reference Technical Ingress without changing any Product operation decision.

W4 receives only a non-normative cross-reference that technical OAuth begin performs a **W4-equivalent current-access evaluation** using `portfolio.manage`; this does not create a 96th Product operation or extend the Permission vocabulary.

### TI-C6 — acquisition-family criterion

A native acquisition family exists because a **distinct authoritative read/coverage contract is required to establish a distinct consumer-owned claim**, never because the provider exposes a topic of the same name.

One topic may awaken several families; one provider read may satisfy several families if authority/coverage is sufficient. Another provider is never forced to imitate Mercado Livre's topic/read decomposition.

### TI-C7 — coverage/recovery is explicit per family

Shared push/poll/recovery path never invents provider discovery/recovery/completeness. Each admitted family must carry its own accepted source/recovery/completeness statement.

First explicit residual: `AcquireMarketplaceSale` / `orders_v2` has **no proven cancellation-inclusive recovery universe**; seller Order enumeration cannot currently prove it. This remains a D8 proof obligation, not implied completeness.

### TI-C8 — three positive acknowledgement bases

Positive provider acknowledgement is allowed only after one explicit technical decision:

1. recoverable attributed custody;
2. bounded quarantine disposition;
3. explicit terminal non-admission/non-processing disposition where provider protocol should be acknowledged.

None means business success or convergence; silent discard is never a fourth basis.

### TI-C9 — quarantine is pre-attribution platform state

Quarantine is platform-scoped pre-attribution protocol state and is explicitly outside D3's Organization-scoped durable communication/recovery class. Once exact Installation attribution exists, the signal leaves quarantine and becomes Organization-scoped recoverable acquisition.

### TI-C10 — Claim action recoverability remains Unknown

`post_purchase:claims_actions` remains admitted as a trigger into Claim acquisition, but whether every materially relevant Claim-action occurrence can be reconstructed from authoritative current Claim evidence remains **Unknown** and a D4/D8 proof obligation.

### TI-C11 — Mercado Livre refresh is consuming

Current official Mercado Livre evidence establishes that only the latest refresh token is accepted and each refresh token is single-use. Therefore refresh of one active credential generation must be serialized/single-consumer by correctness; a stale refresh result cannot overwrite a newer OAuth/refresh generation. D7 chooses the locking/CAS mechanism.

### TI-C12 — selling-account binding subject

The Installation binding subject is the provider **selling-account namespace**, not an arbitrary authorizing operator identity. Current Mercado Livre authorization guidance requires an administrator and rejects operator/collaborator grants, so collaborator authorization is outside the current lane. A future delegated-operator provider requires explicit seller-account identity proof in D4.

### TI-C13 — non-secret lineage is historical explanation only

Authorization/binding lineage may preserve the smallest non-secret provenance necessary for security/history, but the current Installation/D4 binding state remains the sole current namespace authority. Historical lineage never becomes a second current binding source. Authorization code, client secret, access token, refresh token and PKCE verifier never enter that lineage.

### Whole-Ingress dispositions preserved

- native MPC acquisition seam = KEEP;
- seven current marketplace acquisition families = KEEP;
- no generic Integration/Webhook/ExternalEvent/OAuth authority;
- OAuth callback remains separate from acquisition inbox;
- current-authority revalidation, attempt supersession/replay protection and complete credential-generation activation remain;
- successful authorization may wake bounded technical recovery, never Product Sync/Refresh or D3 business event;
- exact technical route prefix/host remains `DEFER SAFELY`, but collision with Product roots is forbidden;
- Structural Inversion = PASS;
- no parent semantic reopen.

---

## 5. Prohibited now

Until final operator Whole-Ingress ratification:

- do not amend Technical Ingress, W1, Operation Matrix or W4 with the converged package;
- do not begin final Problem/media consistency, Product OpenAPI/tooling, D6–D9 or implementation;
- do not create generic Integration/Ingress/Webhook/ExternalEvent/OAuth authority;
- do not reactivate marketplace business participation merely to restore evidence-read credentials;
- do not use quarantine as an attacker-controlled or feature-backlog sink;
- do not fabricate SourceInstance/business-system ingress capability;
- do not claim cancellation-inclusive Sales recovery or Claim-action occurrence completeness before D8 proof;
- do not choose D7 queue/storage/credential/refresh serialization mechanism yet.

---

## 6. Exact next action

**Operator performs final Whole-Ingress ratification of the converged package in §4.**

If ratified:

1. consolidate TI-C1…C13 into the single `D5-B2-TECHNICAL-INGRESS.md` authority;
2. apply bounded W1 / Operation Matrix / W4 routing and cross-reference corrections without semantic parent reopen;
3. remove Whole-Ingress review candidate; Git history remains archive;
4. reset `AI-DIALOG.md` to protocol-only;
5. synchronize cockpit as non-authoritative projection;
6. advance router to **final Problem/media consistency**.

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4 are canonical;
- Technical Ingress A+B are accepted in-stage/operator-ratified;
- Whole-Ingress lead + Fable independent review are complete;
- GPT adjudication converged with TI-C1…C13 and no Round 2;
- no D0→W4/D3/D4 semantic reopen is required;
- **operator final Whole-Ingress ratification is the exact next action**;
- final Problem/media/OpenAPI/D6–D9/implementation remain blocked.

If not, the active authority tree is inconsistent.
