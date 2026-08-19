# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 + Technical Ingress ACCEPTED / CANONICAL; Final Problem/media lead + Fable + GPT review CONVERGED; operator final ratification = NEXT**  
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

`D5-B2-FINAL-PROBLEM-MEDIA-CONSISTENCY-REVIEW-CANDIDATE.md` and the current `AI-DIALOG.md` cycle are **NON-AUTHORITATIVE review evidence**. The converged package does not amend W1/W2/W3/W4/Technical Ingress until explicit operator ratification and canonical filing.

`docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md` is a derived portable production-research guide. It is not part of this architecture authority path and does not redefine the organizational Method.

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
        Whole-Ingress review/adjudication                 COMPLETE / FILED
      Final Problem/media consistency                    OPEN / ACTIVE
        lead coherent review                             COMPLETE
        independent Fable review                         COMPLETE
        GPT final adjudication                           CONVERGED
        Round 2                                           NOT REQUIRED
        operator final ratification                      NEXT
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Final Problem/media converged package — evidence pending ratification

The consolidated candidate incorporates lead analysis, independent Fable challenge and GPT adjudication. The material corrections are:

1. standard status-only Product failures use RFC 9457 `about:blank`, while status-specific HTTP headers remain binding (`405 + Allow`);
2. media failure grammar distinguishes malformed multipart, unsupported representation, excess size, contract validation, revision conflict, current-state conflict, idempotency and internal failure;
3. canonical W2 idempotency mappings are preserved: reused-key remains `422 idempotency-key-reused`, in-progress remains `409 idempotency-request-in-progress`, and missing-key receives an explicit `400 idempotency-key-required` decision;
4. one bounded `409 resource-state-conflict` family distinguishes current lifecycle inadmissibility from stale revision after exact-repeat handling;
5. excess-size/media-type transport guards may fail before AuthN without full-body buffering or tenant/resource disclosure;
6. successful authored-media creation advances the parent ListingIntent validator and returns the new validator as typed custom-method result data;
7. authored-media identity/provenance and volatile presentation access use separate descriptor shapes;
8. authored-media byte delivery is a separately justified technical presentation surface requiring current Organization + `offering.read`, not a Product operation, ingress lane or anonymous durable locator;
9. source and authored media retain distinct locator/access-reference trust types;
10. no client erasure operation or universal retention duration is invented; retention/erasure remains an explicit D2/D7 Unknown with reopen triggers;
11. Product and technical-protocol Problem type namespaces remain disjoint;
12. stale sequencing text is removed substitutively from active canonical artifacts at filing.

Dispositions preserved:

```text
Product operations                              95 unchanged
ordinary Permissions                            29 unchanged
standalone Product media GET/Delete/Update       not admitted
Technical Ingress A/B                            unchanged
D0→D5-B1 semantic parent reopen                  none
implementation                                   blocked until D9
```

The Fable review found real defects in the lead candidate. GPT accepted the two findings that could have required Round 2 — canonical idempotency transcription and byte-access authority — and resolved them without material disagreement. The remaining findings were absorbed or adjudicated locally; no real contradiction survives, so Round 2 is not required.

---

## 4. Allowed now

Only this exact operator action is open:

> **Explicitly ratify or reject/request changes to the converged Final Problem/media package.**

The package to ratify is the consolidated non-authoritative candidate:

`docs/engineering/rebaseline/D5-B2-FINAL-PROBLEM-MEDIA-CONSISTENCY-REVIEW-CANDIDATE.md`

A valid ratification may be expressed plainly, for example `Ratifico` or `Aprovado`, after reviewing the converged summary.

After ratification, the next repository write must:

1. incorporate the converged corrections substitutively into their existing canonical homes, primarily W2 with bounded W1/W4/Technical-Ingress cross-references;
2. preserve 95 Product operations / 29 Permissions;
3. remove stale sequencing text instead of appending another superseding layer;
4. remove the review candidate from the active tree;
5. reset `AI-DIALOG.md` to protocol-only;
6. update this router to mark Final Problem/media accepted/canonical;
7. synchronize the non-authoritative cockpit;
8. advance to the single OpenAPI wire authority/tooling decision.

Implementation remains blocked until D9.

---

## 5. Prohibited now

- do not amend canonical W1/W2/W3/W4/Technical Ingress before operator ratification;
- do not treat the candidate, reviewer severity or this summary as architecture authority;
- do not add a 96th Product operation, 30th Permission, standalone media GET/Delete/Update or generic media library;
- do not place the technical media-delivery surface or provider/OAuth ingress in Product OpenAPI/SDK;
- do not introduce a durable anonymous media locator as baseline;
- do not create a global Product `code` parallel to Problem `type`;
- do not leak object-store, CDN, scanner, transformer or provider detail into Product Problems;
- do not choose D7 blob, proxy, CDN, scanner, transaction or deployment technology now;
- do not reopen D0→D4, D4-R1 or D5-B1 without material evidence;
- do not begin Product OpenAPI/tooling, D6–D9 or implementation out of sequence;
- do not update the cockpit from an unratified candidate.

---

## 6. Exact next action

**Operator final ratification of the converged Final Problem/media package.**

No additional Fable round is required unless the operator requests a material change that reintroduces unresolved contradiction.

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4 and Technical Ingress remain accepted/canonical and unchanged by the candidate;
- Product surface remains 95 operations / 29 Permissions;
- lead review, independent Fable review and GPT adjudication are complete;
- Round 2 is not required;
- the consolidated Final Problem/media candidate is non-authoritative;
- operator final ratification is the exact next action;
- OpenAPI/tooling, D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.
