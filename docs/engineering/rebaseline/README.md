# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 + Technical Ingress ACCEPTED / CANONICAL; Final Problem/media lead review COMPLETE; independent Fable review = NEXT**  
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

`D5-B2-FINAL-PROBLEM-MEDIA-CONSISTENCY-REVIEW-CANDIDATE.md` and the current `AI-DIALOG.md` cycle are **NON-AUTHORITATIVE review evidence**. Their findings do not amend canonical W1/W2/W3/W4/Technical Ingress until final adjudication, explicit operator ratification and filing.

`docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md` is a derived portable production-research guide. It is not part of this architecture authority path and does not redefine the organizational Method.

`docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection. It is synchronized after canonical architecture/status closure, not used to drive the active review.

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
      Final Problem/media consistency                    OPEN / ACTIVE
        lead coherent review                             COMPLETE / RESTRUCTURE D5-B2-LOCAL
        independent Fable review                         NEXT
        GPT adjudication                                 BLOCKED BY REVIEW
        operator final ratification                      BLOCKED BY ADJUDICATION
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

The single canonical authority is `D5-B2-TECHNICAL-INGRESS.md`. It incorporates TI-C1…TI-C13 and preserves:

- provider-specific transport → one MarketplaceInstallation-scoped recoverable acquisition seam → seven closed families → authoritative D4 acquisition → owner meaning;
- no generic Integration/Webhook/ExternalEvent/OAuth/provider-resource business authority;
- bounded pre-attribution quarantine and explicit acknowledgement bases;
- honest per-family coverage and the unproven Sales-cancellation / Claim-action residuals;
- same-seller evidence-recovery credential restoration without business reactivation;
- transaction-bound OAuth protection, selling-account/admin proof, generation-safe activation and serialized consuming refresh;
- Product OpenAPI/SDK/operation separation;
- 95 Product operations / 29 Permissions unchanged;
- no D0→W4/D3/D4 semantic parent reopen.

Whole-Ingress review evidence remains only in Git history; its former candidate is absent from the active tree.

---

## 4. Final Problem/media lead-review state — evidence only

The coherent lead review is recorded in:

`D5-B2-FINAL-PROBLEM-MEDIA-CONSISTENCY-REVIEW-CANDIDATE.md`

It is non-authoritative and proposes PM-C1…PM-C8 for independent challenge:

1. **standard status-only failures** use RFC 9457 `about:blank` instead of multiplying custom MPC problem types;
2. **one explicit media negative map** distinguishes malformed multipart, 413, 415, validation, revision, idempotency and internal failure without leaking storage/scanner/provider vocabulary;
3. **Problem `type` remains the sole global Product discriminator**; no duplicate global `code` taxonomy;
4. **Product versus technical-protocol error registries remain separate**, even where both reuse a standard representation format;
5. **ListingIntent authored-media identity is immutable and scoped**; selection removal is not deletion; no generic media CRUD;
6. **binary delivery/access is mechanism, not identity**; a bounded descriptor access seam may support authorized inspection without admitting a standalone Product media GET;
7. **source media and authored media remain distinct bounded descriptor families**; no generic Media/Asset authority;
8. **stale active sequencing text is removed at final filing** so only this router owns current next action.

Lead disposition:

```text
current semantic structure                     CONFIRMED
Problem/media consistency                      RESTRUCTURE NOW — D5-B2 LOCAL
Product operations                             95 unchanged
ordinary Permissions                           29 unchanged
parent semantic reopen                         NONE
independent review                             WARRANTED
```

Canonical W1/W2/W3/W4/Technical Ingress remain unchanged until the complete review cycle closes.

---

## 5. Allowed now

Only the exact independent-review action is open:

> **Fable performs one coherent independent adversarial review of the complete Final Problem/media candidate and appends findings only to `AI-DIALOG.md`.**

The reviewer must independently reconstruct authority and challenge at least:

- `about:blank` versus custom Product problem types;
- the complete 400/413/415/422/409/500 media map;
- `type` versus duplicate `code`;
- Product Problem versus technical provider/OAuth vocabulary;
- media identity, idempotency, revision, selection/history/deletion;
- binary inspectability/access without hidden Product media GET or storage identity leakage;
- untrusted-binary security properties without selecting D7 realization;
- source versus authored media separation;
- stale status/routing cleanup;
- operation/Permission counts and parent reopen.

The review scope and write authorization are recorded in `AI-DIALOG.md`.

After Fable returns:

1. GPT adjudicates every material finding against repository authority;
2. Round 2 occurs only if a material contradiction survives;
3. the operator explicitly ratifies the converged package;
4. only then are canonical artifacts amended substitutively, review evidence removed/reset and cockpit synchronized;
5. router advances to the single OpenAPI wire authority/tooling decision.

---

## 6. Prohibited now

- do not amend W1/W2/W3/W4/Technical Ingress from the lead candidate before review/adjudication/ratification;
- do not treat the candidate or reviewer severity as authority;
- do not add a 96th Product operation, 30th Permission, standalone Product media GET/Delete/Update or generic media library;
- do not add provider/OAuth callback operations or errors to Product OpenAPI/SDK;
- do not create a global Product `code` parallel to RFC 9457 `type` by tooling convention;
- do not leak object-store, CDN, scanner, transformer or provider vocabulary/details into Product Problems;
- do not choose blob store, CDN, signed-URL, scanner, image processing, retention, transaction or deployment realization;
- do not reopen D0→D4, D4-R1 or D5-B1 without material evidence;
- do not begin Product OpenAPI/tooling, D6–D9 or implementation out of sequence;
- do not synchronize cockpit from a non-ratified review candidate.

Implementation remains blocked until D9.

---

## 7. Exact next action

**Independent Fable Final Problem/media review in `AI-DIALOG.md`.**

No operator ratification is requested until GPT has adjudicated the independent evidence.

Implementation remains blocked until D9.

---

## 8. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4 and Technical Ingress remain accepted/canonical and unchanged by the candidate;
- Product surface remains 95 operations / 29 Permissions;
- Final Problem/media lead review is complete with PM-C1…PM-C8 as non-authoritative evidence;
- independent Fable review is the exact next action;
- AI-DIALOG contains the active bounded review handoff;
- no semantic parent reopen is currently proven;
- Product OpenAPI/tooling, D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.
