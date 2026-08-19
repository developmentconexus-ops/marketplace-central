# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 CANONICAL; W3-A/B/C ACCEPTED IN-STAGE; lead Whole-W3 Global Coherence COMPLETE / RESTRUCTURE W3-LOCAL; operator ratification of W3-G1…G6 = NEXT**  
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
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — W3-A/B accepted in-stage
19. `docs/engineering/rebaseline/D5-B2-W3-C-CURSOR-SAFETY.md` — W3-C accepted in-stage
20. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
21. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

`D5-B2-W3-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` is **NON-AUTHORITATIVE lead review evidence** and is deliberately excluded from the authority path. W3-G1…G6 do not modify W3/W2 until operator ratification and canonical filing.

W3-C remains a bounded staging artifact. After Whole-W3 convergence + operator ratification, W3-A/B/C must be consolidated into one canonical W3 authority and staging/review artifacts removed; Git history remains the archive.

`docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection only and never participates in the authority path.

Legacy/current code, routes, OpenAPI, SDK and frontend tables remain evidence only.

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
    B2-A Client/Auth                                     ACCEPTED IN-STAGE
    Operation Admission Matrix                           ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                        ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar                  ACCEPTED / CANONICAL
      W2 Request / Response Schema Grammar               ACCEPTED / CANONICAL
      W3 Collections / Pagination / Filter / Search /
         Cursor Grammar                                  OPEN / ACTIVE
        W3-A Collection + Cursor Core                    ACCEPTED IN-STAGE
        W3-B Per-Operation Filter / Search / Ordering    ACCEPTED IN-STAGE
        W3-C Cursor Safety / Population / Limits         ACCEPTED IN-STAGE
        Whole-W3 lead Global Coherence                   COMPLETE / RESTRUCTURE W3-LOCAL
        operator lead-direction ratification             NEXT
      Remaining wire obligations                         BLOCKED BY W3 SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

Whole-W2 independent review is complete and incorporated into canonical W1/W2. W3 is the only active Wire design surface.

---

## 3. Load-bearing current authority

### 3.1 W1/W2 laws that W3 may not weaken

- Product API is semantic-owner driven and Organization-scoped; protocol/provider topology is not Product ontology.
- external Product/Listing/Sale/Shipment identity remains source-qualified; no mirror IDs or bare native correlation keys.
- custom owner methods remain explicit capabilities; no generic Action/Command/Workflow surface.
- one opaque ETag revision authority per protected meaning; true same-resource standard mutation uses `If-Match`; custom/reference revision proofs use typed ETag request data.
- Idempotency-Key and revision/concurrency solve different failure classes.
- request/read schemas remain authority-separated and closed; `null` never means unknown/unavailable/partial/not-applicable.
- no universal Fact/Evidence/Result/Subject/Scope/Policy/Workflow wrappers.
- ListingIntent / PriceIntent / Availability remain distinct through publication.
- Market coverage is distinct from evidence sufficiency; Economics never fabricates conclusions from missing evidence.
- Work never becomes source-truth or generic resolution authority.

### 3.2 Accepted W3-A

- owner-specific collection responses + optional `next_cursor`; no universal Page/PagedResult/data/metadata wrapper;
- forward-only `limit?` + opaque `cursor?`; no page/offset/skip/previous cursor baseline;
- fewer than `limit` items does not prove exhaustion; only `next_cursor` indicates another page;
- cursor is opaque, query/operation/Organization-bound, not authorization and not raw provider/database state;
- cursor exhaustion never proves source/provider/market/all-time completeness;
- no universal total count, caller sort or snapshot-isolation promise;
- provider paging remains D4-local; cursor persistence/signing/index/cache remains D7.

### 3.3 Accepted W3-B

- typed operation-specific query parameters only; no generic query/filter expression language;
- admitted fields combine by AND; no OR/NOT/IN/functions/traversal DSL baseline;
- `SearchSourceProductsForMarketplace` is the only baseline Search;
- caller-selectable sort baseline count = 0;
- each collection has one owner-defined default order/tie-breaker contract;
- source-qualified identity and same-Organization safety remain explicit;
- collection coverage is owner-specific only where external enumeration can be incomplete.

The detailed accepted matrix remains in `D5-B2-W3-COLLECTION-GRAMMAR.md` until Whole-W3 correction is ratified.

### 3.4 Accepted W3-C

- owner-specific bounded list-item projection may be smaller than point Get; no generic `fields/select/expand` DSL;
- `400 invalid-cursor` and `400 cursor-expired`; no silent restart;
- cursor is ephemeral continuation, not durable bookmark; no public TTL baseline;
- stable MPC/source-qualified identity appears at most once per Product traversal;
- mutable population may affect not-yet-traversed members without implying snapshot/completeness;
- updates do not turn pagination into a change feed;
- non-identifiable evidence receives no fake ID merely for deduplication;
- provider continuation may be reconstructed only with identical Product continuation semantics, otherwise expire explicitly;
- transient provider outage is not cursor expiration by itself;
- no public snapshot/traversal/search-session resource;
- positive finite `limit`, server returns <= limit, zero/negative/above-max validation error, no silent clamp;
- exact numeric default/max remains `DEFER SAFELY` to final OpenAPI closure under fixed finite/no-clamp fences.

---

## 4. Lead Whole-W3 findings — NON-AUTHORITATIVE UNTIL OPERATOR RATIFICATION

The lead review found **six W3-local corrections; no current parent-stage reopen**:

1. **W3-G1 — continuation query carrier:** cursor must carry continuation only; every continuation request repeats all operation-required semantic subject/search fields and the same effective optional filters; only `limit` may vary.
2. **W3-G2 — keyed-Q list populations:** `ListCompetitivePositions` and `ListExpectedEconomics` must enumerate existing marketplace Listing subjects only; pre-listing source-product contexts remain explicit point Get/evaluation paths, never an implicit Product universe.
3. **W3-G3 — deterministic ordering:** close vague tuple placeholders; Availability uses fixed target-kind+identity order, Competitive/Expected use Listing ref after G2, EconomicAttribution uses its stable resource ID as final ordering key, Shipment uses source-qualified native key.
4. **W3-G4 — ComparableOffer evaluation basis:** one cursor chain remains on one Market evaluation/acquisition basis; if that basis cannot resume, `cursor-expired`; no fake ComparableOffer ID/public snapshot resource.
5. **W3-G5 — Source Product search:** keep required source-qualified opaque query, exact identifier matches ahead of textual results where supported, but do not freeze an unproven universal tokenizer/case algorithm or silently create an MPC Product mirror/index to satisfy it.
6. **W3-G6 — one Problem Details catalog:** final Whole-W3 consolidation must amend canonical W2 problem catalog with `invalid-cursor` and `cursor-expired`; no second W3-only error taxonomy.

Full reasoning/alternatives/reopen triggers are in the non-authoritative review candidate.

The following survived adversarial attack: 26/26 admitted List/Search coverage, no new list-by-symmetry, no generic Page/filter/sort/snapshot/total framework, at-most-once stable member identity, bounded list-item projection, and safely deferred exact numeric limit values.

---

## 5. Prohibited now

Until operator ratifies/revises the lead Whole-W3 direction:

- do not apply W3-G1…G6 silently to W3-A/B/C or canonical W2;
- do not begin remaining Wire Contract obligations, D6–D9 target design or implementation;
- do not treat the review candidate or cockpit as authority;
- do not introduce generic query/filter/sort/projection DSL, universal snapshot/traversal resources, total counts or fabricated evidence identities;
- do not use current code/OpenAPI/frontend/provider paging shape to overrule accepted W3 authority;
- do not choose exact OpenAPI minor/generator or numeric limit defaults yet.

---

## 6. Exact next action

**Operator reviews and ratifies/revises W3-G1…W3-G6 as the lead Whole-W3 correction direction.**

If ratified, the next step is to adjudicate whether proportional independent challenge adds material value before final canonical consolidation; no independent review is automatically authority and no extra round is created by ceremony.

After Whole-W3 final convergence + operator ratification:

1. amend W2 canonical cursor problem catalog as required;
2. consolidate W3-A/B/C + ratified corrections into one canonical `D5-B2-W3-COLLECTION-GRAMMAR.md`;
3. remove W3-C staging and Whole-W3 candidate; Git history remains archive;
4. update the cockpit projection separately without making it authority;
5. advance to the next router-ordered Wire Contract obligation.

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- W1 and W2 are canonical;
- W3-A/B/C remain accepted in-stage current authority;
- `D5-B2-W3-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` is non-authoritative evidence only;
- lead Whole-W3 review found W3-G1…G6 and `RESTRUCTURE NOW — W3-LOCAL`;
- **operator ratification/revision of W3-G1…G6 is the exact next action**;
- no parent-stage reopen is currently proven;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
