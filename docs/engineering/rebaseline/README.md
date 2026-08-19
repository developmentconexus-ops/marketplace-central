# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 + Technical Ingress ACCEPTED / CANONICAL; final Problem/media consistency = NEXT**  
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
20. `docs/engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md` — canonical technical non-Product ingress
21. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
22. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in accepted artifacts.

`AI-DIALOG.md` is a protocol-only, non-authoritative review channel. Completed Technical Ingress review evidence and the removed Whole-Ingress candidate remain available only in Git history.

`docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection and is synchronized only after canonical status changes.

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
      Technical non-Product ingress                      ACCEPTED / CANONICAL
        A External Acquisition Ingress                   ACCEPTED / CANONICAL
        B OAuth / Authorization Ceremony                 ACCEPTED / CANONICAL
        Whole-Ingress lead review                        COMPLETE
        Fable independent review                         COMPLETE
        GPT final adjudication                           CONVERGED
        Round 2                                           NOT REQUIRED
        operator final Whole-Ingress ratification        COMPLETE / FILED
      Final Problem/media consistency                    NEXT
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Technical Ingress closure record

The operator's final Whole-Ingress ratification was materialized on 2026-08-19.

The single canonical authority is now `D5-B2-TECHNICAL-INGRESS.md`. It incorporates the converged TI-C1…TI-C13 package and preserves these dispositions:

- provider-specific transport feeds one MPC-native, recoverable, typed marketplace-acquisition seam;
- the seven Product 1.0 acquisition families remain Listing, Price, Sale, Shipment, Payment, Post-Sale Claim and Competitive Position;
- acquisition families arise from distinct authoritative read/coverage contracts plus consumer-owned claims, never provider topic symmetry;
- provider signal is neither D3 domain event nor current provider truth;
- Product 1.0 inbound acquisition is MarketplaceInstallation-scoped; no SourceInstance callback seam is fabricated for Sankhya;
- quarantine is bounded pre-attribution platform state, not feature backlog, attacker payload storage or D3 Organization communication;
- positive acknowledgement requires attributed recoverable custody, bounded quarantine disposition or explicit terminal non-processing; it never means business success;
- per-family recovery/completeness remains explicit, including the unproven cancellation-inclusive `orders_v2` recovery universe and Unknown Claim-action occurrence completeness;
- deactivation preserves retained namespace attribution and permits same-seller technical credential restoration for evidence recovery without business reactivation;
- OAuth begin uses current authenticated human + Principal eligibility + Membership + `portfolio.manage` + exact Installation but creates no Product operation;
- current Mercado Livre authorization uses transaction-bound anti-injection protection with PKCE where supported, administrator authorization and selling-account namespace proof;
- credential activation is generation-safe; current Mercado Livre refresh is single-use/consuming and therefore serialized per active generation by correctness, with D7 choosing mechanism;
- non-secret authorization/binding lineage is historical security explanation only; current Installation/D4 binding remains the sole current namespace authority;
- technical routes remain outside Product roots/OpenAPI/SDK; exact host/prefix is deferred, but collision with `/access-context` or `/organizations/{organization_id}/...` is forbidden;
- no generic Integration, Webhook, ExternalEvent, OAuth or provider-resource business authority was introduced;
- no D0→W4/D3/D4 semantic parent reopen occurred.

The Whole-Ingress review candidate was removed from the active tree; Git history is the archive. `AI-DIALOG.md` was reset to protocol-only. W1, the Operation Matrix and W4 received routing/cross-reference corrections only; Product operation and Permission semantics remain unchanged at **95 operations / 29 Permissions**.

---

## 4. Allowed now

Only the router's exact next action is open:

> **Perform final Problem/media consistency across the canonical D5-B2 package.**

The review must derive from the accepted authorities, not legacy OpenAPI/routes/controllers, and must close only D5-B2-local inconsistencies involving:

1. W2 Product Problem Details versus technical provider/OAuth protocol-local failures;
2. ListingIntent-scoped authored-media intake, Product representations and D7-deferred blob/storage/delivery mechanics;
3. W1/W2/W3/W4/Technical Ingress cross-references and terminology;
4. the absence of provider callback/webhook/OAuth operations or errors from the Product OpenAPI/SDK surface;
5. any remaining duplicate meaning, hidden second authority or contradictory wire claim.

Use Global Maximum reasoning without widening scope. Reopen a parent only with material evidence and only the smallest implicated scope.

After operator ratification and canonical filing of final Problem/media consistency, the router may advance to the **single OpenAPI wire authority/tooling decision**.

---

## 5. Prohibited now

- do not reopen D0→D4, D4-R1, D5-B1, the Operation Matrix or W1→W4 for naming/style/current implementation preference;
- do not create generic Integration/Ingress/Webhook/ExternalEvent/OAuth/provider-resource authority;
- do not add Product Sync/Refresh/OAuth/Webhook operations or provider protocol errors to the Product API;
- do not place technical ingress in the Product OpenAPI/SDK;
- do not fabricate SourceInstance/Sankhya callback capability;
- do not claim cancellation-inclusive Sales recovery or Claim-action occurrence completeness before D8 proof;
- do not choose D7 queue/storage/worker/secret/transaction/lock/CAS/refresh-serialization realization;
- do not begin Product OpenAPI/tooling, D6–D9 or implementation out of sequence;
- do not restore removed candidates, completed review dialogue or cockpit text as parallel authority.

Implementation remains blocked until D9.

---

## 6. Exact next action

**Final Problem/media consistency.**

Expected completion sequence:

1. produce one coherent D5-B2-local consistency review;
2. use independent review only if material contradiction warrants it under the repository method;
3. adjudicate evidence against repository authority;
4. obtain explicit operator ratification of the converged package;
5. consolidate into the existing canonical authority homes and remove staging evidence from the active tree;
6. synchronize this router and the non-authoritative cockpit;
7. advance to the single OpenAPI wire authority/tooling decision.

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4 are canonical;
- Technical Ingress A+B and Whole-Ingress are accepted/canonical;
- the lead review, Fable independent review and GPT adjudication are complete; Round 2 was not required;
- operator final Whole-Ingress ratification is complete and canonically filed;
- the review candidate is absent from the active tree and `AI-DIALOG.md` is protocol-only;
- the Product surface remains 95 operations / 29 Permissions;
- no D0→W4/D3/D4 semantic reopen occurred;
- **final Problem/media consistency is the exact next action**;
- Product OpenAPI/tooling, D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.
