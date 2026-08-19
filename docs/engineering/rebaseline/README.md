# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media consistency ACCEPTED / CANONICAL; single OpenAPI wire authority / tooling decision = NEXT**  
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
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical W2, including the canonical Problem Details and authored-media grammar
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — canonical W3
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — canonical W4
20. `docs/engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md` — canonical technical non-Product ingress
21. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
22. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in accepted artifacts.

`docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md` is a derived portable production-research guide. It is not part of this architecture authority path and does not redefine the organizational Method.

`AI-DIALOG.md` is a working review channel with no active cycle. It is never authority.

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
      Final Problem/media consistency                    ACCEPTED / CANONICAL
        lead + Fable + GPT review                        COMPLETE / CONVERGED
        operator final ratification                      COMPLETE / FILED
      Single OpenAPI wire authority/tooling decision     OPEN / NEXT
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Final Problem/media consistency — filed dispositions

The converged lead + Fable + GPT package was operator-ratified and filed substitutively into its existing canonical homes. Detailed semantics live only in those artifacts; this list is a routing index, not a second authority.

| Correction | Canonical home |
|---|---|
| RFC 9457 `about:blank` for status-only failures; `405 + Allow`, `415 Accept`, `413 Retry-After` obligations | W2 §§2.15.1–2.15.2 |
| exact `CreateListingIntentMedia` failure map (400/413/415/422/409/500) | W2 §3.9.2 |
| `400 idempotency-key-required`; `422 idempotency-key-reused` and `409 idempotency-request-in-progress` preserved unchanged | W2 §§3.9.2, 18, 19 |
| bounded `409 resource-state-conflict` family | W2 §§3.9.2, 19 |
| transport size/media-type guards before full-body buffering, without tenant/resource disclosure | W2 §§3.9.3, 18 |
| immutable ListingIntent-scoped authored-media identity | W2 §3.9.5 |
| successful media creation advances and returns the parent ListingIntent validator as typed result data | W2 §3.9.1 |
| exact idempotent retry resolves the prior intake before re-evaluating the now-stale validator | W2 §§3.9.1, 18 |
| stable `ListingIntentMediaDescriptor` vs volatile `ListingIntentMediaPresentationDescriptor` | W2 §3.9.7 |
| access reference never enters history, fingerprint, logs or Problem Details | W2 §3.9.7 |
| authored-media byte delivery = separately justified technical presentation surface reusing current Organization + `offering.read`; no durable anonymous locator; failures do not expand the Product Problem catalog | W2 §3.9.8, W1 §20, W4 §15.2, Technical Ingress §12.1 |
| retention/erasure remains Unknown / deferred to D2 + D7; no `DeleteMedia`, no invented lifetime | W2 §3.9.6 |
| Product and technical-protocol Problem namespaces remain disjoint | W2 §19 |
| source vs authored media stay distinct descriptor families | W2 §3.9.9 |
| converged negative controls / proof obligations | W2 §23 items 43–67 |
| stale sequencing text removed substitutively from active B2 artifacts | Product Operation Surface §5, Admission Matrix §§9–10, W1 §§19–20, W4 §§14–15, Technical Ingress §§12.1, 24 |

Dispositions preserved by this filing:

```text
Product operations                              95 unchanged
ordinary Permissions                            29 unchanged
standalone Product media GET/Update/Delete       not admitted
new Principal kind                               none
Technical Ingress A/B                            semantically unchanged
W3 collection/cursor semantics                   unchanged
D0→D5-B1 semantic parent reopen                  none
blob store / CDN / scanner / proxy / generator /
server framework / deployment selection          none
implementation                                   blocked until D9
```

---

## 4. Allowed now

> **Open the single OpenAPI wire authority / tooling decision for the Product API.**

That decision must derive from already-canonical W1/W2/W3/W4 and the Technical Ingress boundary, and must at minimum settle:

1. the one machine-readable Product API wire authority and where it lives;
2. OpenAPI minor selection justified by actual W1–W4 expression needs and real generator compatibility, not recency;
3. operation naming/spelling and `operationId` law consistent with W1 paths and custom methods;
4. Problem `type` URI host/namespace, keeping Product and technical-protocol namespaces disjoint;
5. how multipart `CreateListingIntentMedia` and typed ETag parts are expressed without inventing new semantics;
6. what is deliberately excluded — technical ingress lanes A/B and the authored-media delivery surface;
7. whether an SDK is generated from that single authority, with no second hand-authored representation;
8. the negative controls that prove the document cannot drift from W1–W4.

The decision does **not** choose server framework, runtime, deployment, storage, CDN or D6/D7 topology.

Implementation remains blocked until D9.

---

## 5. Prohibited now

- do not amend canonical W1/W2/W3/W4/Technical Ingress/Problem-media semantics as a side effect of tooling selection;
- do not add a 96th Product operation, 30th Permission, new Principal kind, standalone media GET/Delete/Update or generic media library;
- do not place the technical media-delivery surface or provider/OAuth ingress in Product OpenAPI/SDK;
- do not introduce a durable anonymous media locator as baseline;
- do not create a global Product `code` parallel to Problem `type`;
- do not leak object-store, CDN, scanner, transformer or provider detail into Product Problems;
- do not choose D7 blob, proxy, CDN, scanner, transaction or deployment technology now;
- do not reopen D0→D4, D4-R1 or D5-B1 without material evidence;
- do not begin D6–D9 or implementation out of sequence;
- do not treat `AI-DIALOG.md`, retired candidates, review dialogue or the cockpit as authority.

---

## 6. Exact next action

**Open the single OpenAPI wire authority / tooling decision as the next D5-B2 sub-batch.**

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4, Technical Ingress and Final Problem/media consistency are accepted/canonical;
- Product surface remains 95 operations / 29 Permissions, with no new Principal kind;
- authored-media byte delivery is a technical presentation surface, not a Product operation or ingress lane;
- retention/erasure remains an explicit D2/D7 Unknown;
- the Final Problem/media review candidate is gone from the active tree and `AI-DIALOG.md` has no active cycle;
- the single OpenAPI wire authority/tooling decision is the exact next action;
- D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.
