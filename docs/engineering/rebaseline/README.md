# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 CANONICAL; W3-A/B/C ACCEPTED IN-STAGE; Whole-W3 independent review + GPT final adjudication CONVERGED / RESTRUCTURE W3-LOCAL; operator final Whole-W3 ratification = NEXT**  
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

`D5-B2-W3-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` are **NON-AUTHORITATIVE review evidence**, deliberately outside the authority path. The converged Whole-W3 corrections do not amend W2/W3 until final operator ratification and canonical filing.

`D5-B2-W3-C-CURSOR-SAFETY.md` remains a bounded staging artifact. After final operator ratification, W3-A/B/C + converged corrections must be consolidated into one canonical W3 authority; W3-C staging and Whole-W3 candidate are removed, and Git history remains the archive.

`docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection only. It is updated after canonical filing, never used as status authority.

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
        W3-B Filter / Search / Ordering                  ACCEPTED IN-STAGE
        W3-C Cursor Safety / Population / Limits         ACCEPTED IN-STAGE
        Whole-W3 lead review                             COMPLETE / RESTRUCTURE W3-LOCAL
        Fable Whole-W3 independent review                COMPLETE
        GPT final adjudication                           CONVERGED
        operator final Whole-W3 ratification             NEXT
      Remaining wire obligations                         BLOCKED BY W3
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Load-bearing current authority

### 3.1 W1/W2 laws W3 may not weaken

- semantic-owner Product API under explicit Organization scope;
- source-qualified external Product/Listing/Sale/Shipment identity; no mirror IDs/bare native correlation keys;
- owner capabilities remain explicit; no generic Action/Command/Workflow owner;
- one opaque revision authority; true same-resource HTTP conditional vs typed custom/reference ETag carriers remain distinct;
- request/read authority separation, closed schemas and honest unknown/unavailable/partial semantics;
- no universal Fact/Evidence/Result/Subject/Scope/Policy/Workflow wrappers;
- ListingIntent / PriceIntent / Availability remain distinct;
- Market coverage != evidence sufficiency; Economics never fabricates from missing evidence;
- Work never becomes source truth;
- implementation remains blocked until D9.

### 3.2 Accepted W3-A/B/C baseline that survived Whole-W3

- owner-named collection responses; no generic Page/PagedResult/data/metadata wrapper;
- forward-only `limit?` + opaque `cursor?`; no offset/page/previous cursor baseline;
- cursor exhaustion never proves source/provider/market/all-time completeness;
- no universal total count, caller-selectable sort or snapshot resource/promise;
- typed operation-specific filters only; no OData/SQL/GraphQL-like query language;
- SearchSourceProductsForMarketplace remains the only baseline Search;
- source qualification and same-Organization safety remain explicit;
- stable MPC/source-qualified member identity appears at most once per traversal;
- mutable population may affect not-yet-traversed members; pagination is not a change feed;
- bounded owner-specific ListItem projections are allowed, with no generic projection DSL;
- `invalid-cursor` and `cursor-expired` remain the two collection cursor problems, both HTTP 400;
- cursor is ephemeral continuation, not a durable bookmark; no public TTL baseline;
- positive finite `limit`, no silent clamp; exact default/max remains DEFER SAFELY to final OpenAPI contract after concrete ListItem payloads exist;
- D4 owns provider paging/source coverage; D7 owns cursor persistence/signing/index/cache realization.

---

## 4. Converged Whole-W3 corrections — AWAITING FINAL OPERATOR RATIFICATION

Independent Fable review confirmed G1…G6 and added F-W3-1/F-W3-2. GPT final adjudication converged them as follows:

1. **Continuation query carrier:** cursor carries continuation only. Every continuation request repeats required semantic subject/search fields and the same effective optional filters; only `limit` may vary. Missing required fields → ordinary `422 validation-error`; a well-formed repeated query that mismatches cursor binding → `400 invalid-cursor`.
2. **CompetitivePosition / ExpectedEconomics populations:** Lists enumerate currently known existing marketplace Listing subjects only. Pre-listing contexts remain explicit Search → selected SourceProductRef → point Get/evaluation. Provider-universe completeness is bounded by Listing acquisition coverage; no duplicate coverage field is added.
3. **Deterministic ordering:** close all placeholder tuples. `ListEconomicAttributions` total order is exactly `economic_attribution_id ASC`. Availability, Competitive/Expected and Shipment use the explicit stable target/ref/key orders recorded by the converged package.
4. **ComparableOffer continuity:** one cursor chain stays on one Market Intelligence owner-local evaluation/acquisition basis. No shared EvaluationBasis/TraversalBasis type, no public Snapshot/Traversal resource, no fabricated ComparableOffer ID. Same basis unavailable → `cursor-expired`.
5. **Source Product Search:** required non-empty opaque query + explicit SourceInstance/Marketplace context; exact native/legitimate identifiers get precedence where sanctioned evidence establishes them. No universal tokenizer, case-folding, diacritics/accent-folding, stemming, fuzzy/edit-distance, locale/collation, vector/embedding/AI algorithm or public relevance score. If materially required Search cannot be satisfied through sanctioned reads/projections without a new Product-search data authority, STOP/re-adjudicate rather than silently create a Product mirror/index.
6. **One Problem Details catalog:** final W2 catalog gains `invalid-cursor` + `cursor-expired`; W3 owns applicability. Keep these names; no extra cursor taxonomy.
7. **SellableAvailability population:** current population = `new_listing` ListingIntents still genuinely pre-creation in Intent lifecycle `draft | submitted` plus currently known existing Listing subjects. `discarded` is excluded. Once a provider Listing is established, current Availability is addressed by `existing_listing`; unknown/unavailable Availability remains visible rather than disappearing because no persistence row exists. Existing-Listing completeness remains Listing-acquisition-bounded, with no new coverage field.
8. **ListItem semantic-subset law:** where a member has a point resource/Q, ListItem may omit detail but cannot introduce a list-only derived business conclusion. Reused fields preserve the same schema/name/meaning except qualifiers already made unambiguous by operation/path scope. Collection-only owner values remain their owner-native W2 schemas. No `fields/select/expand`, generic Summary/View or projection DSL.
9. **Surviving fences:** at-most-once stable identity, coverage/exhaustion/dedup independence, owner-local basis wording, 26/26 admitted collection coverage, zero list-by-symmetry, zero caller-sort baseline and bounded numeric-limit deferral remain intact.

No D0→D5-B1/W1 parent reopen is required. W2 requires only the bounded cursor-problem catalog amendment during final canonical consolidation.

---

## 5. Prohibited now

Until final operator Whole-W3 ratification:

- do not modify canonical W2 or accepted W3-A/B/C from the converged review package;
- do not begin later Wire obligations, D6–D9 target design or implementation;
- do not treat review candidate, AI-DIALOG or cockpit as authority;
- do not add generic query/filter/sort/projection/snapshot/traversal frameworks, total counts or fabricated evidence identities;
- do not choose exact OpenAPI minor/generator or numeric `limit` defaults yet;
- do not create a Round 2 absent a newly demonstrated material contradiction.

---

## 6. Exact next action

**Operator final-ratifies or revises the complete converged Whole-W3 package above.**

If ratified:

1. revalidate branch HEAD;
2. amend canonical W2 Problem Details catalog with `invalid-cursor` + `cursor-expired`;
3. consolidate W3-A/B/C + converged corrections into one canonical `D5-B2-W3-COLLECTION-GRAMMAR.md`;
4. remove `D5-B2-W3-C-CURSOR-SAFETY.md` and `D5-B2-W3-WHOLE-COHERENCE-REVIEW-CANDIDATE.md`; Git history remains archive;
5. reset `AI-DIALOG.md` to protocol-only after canonical filing;
6. update the non-authoritative cockpit projection;
7. update this router to W3 ACCEPTED/CANONICAL and advance to the next router-ordered Wire Contract obligation.

Implementation remains blocked until D9.

---

## 7. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- W1/W2 are canonical;
- W3-A/B/C remain current accepted in-stage authority until final ratification/canonical filing;
- Whole-W3 lead + independent Fable + GPT final adjudication converged with G1…G6 + F-W3-1/F-W3-2 and no Round 2;
- **operator final Whole-W3 ratification is the exact next action**;
- no parent-stage reopen is currently required;
- later Wire obligations, D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.
