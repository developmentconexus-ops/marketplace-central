# D6-R1 Marketplace Performance Intelligence — Independent Fable Review

> Review branch only: `review/d6-fable`
> Candidate branch: `stage/d6-frontend`
> Candidate HEAD at handoff: `43bd8b2b9d126931ac488a11ebc311ef0f12b0eb`
> Candidate PR: #54 — `docs(d6): open frontend authority stage`
> Base `main` at handoff: `9d2c81e175bc39ac388c9d8924ddad21f2a86480`
> Current Product candidate: 99 Product operations / 30 ordinary Permissions / Principal kinds H-A-S only
> Active runtime baseline: NONE
> D7–D9: BLOCKED
> Product implementation: BLOCKED UNTIL D9

## Review protocol

This is a **material independent architecture / Product-contract / frontend-coherence review** requested by the operator before any further D6 frontend topology or dependency work.

Do **not** optimize for agreement with GPT, the operator, prior chat reasoning, CodeRabbit, prior Fable reviews, or the current candidate. Reconstruct the reasoning independently from repository authority, executable evidence, current primary external evidence and the canonical DevelopmentConexus Method.

Reviewer output is **Evidence, not authority**. Do not change the candidate branch. Write only below `## Fable response` in this file on `review/d6-fable`.

Before analysis, revalidate:

1. remote `main`;
2. `stage/d6-frontend` exact HEAD;
3. PR #54 base/head/state/changed files;
4. CI on the exact candidate;
5. that this review branch differs from the candidate by **only** `docs/work/current/ai-dialog.md`.

If the candidate moved from the SHA above, stop reviewing the stale SHA and record the new exact candidate SHA first.

## Reading discipline

Start strictly:

1. `AGENTS.md`
2. `docs/index.md`
3. `docs/roadmap.md`
4. `docs/engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md`
5. `ARCHITECTURE.md`

Then switch **bounded packs** as concrete questions require them. Do not recursively ingest the entire rebaseline at once.

Recommended review sequence:

### Product / boundary pack

- `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
- `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`

### Identity / historical-evidence / communication pack

- `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
- `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`

### External evidence pack

- `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
- only the concrete current provider evidence needed to adjudicate Mercado Livre Visits / Product Ads / advertiser binding / retention / attribution / contractual admissibility

### Product API pack

- canonical `contracts/api/product/openapi.yaml`
- `contracts/api/product/paths-access-performance.yaml`
- `contracts/api/product/paths-performance.yaml`
- switch to D5 W1/W2/W3/W4 narrative authority only if the OAD cannot answer a concrete semantic question
- inspect `scripts/verify-product-oad.mjs`, `scripts/verify-product-oad-baseline.mjs` and `scripts/gate.ps1` when challenging proof quality

### Frontend-consumption pack

- `docs/engineering/rebaseline/D6-FRONTEND.md`
- `docs/engineering/rebaseline/D6-B1-INTERACTION-MAP.md`
- `qualification/d6-wireframes/index.html`

For method / evidence discipline use:

- `docs/development/evidence-grounded-production-engineering-for-llm-agents.md` (derived, non-authoritative);
- canonical `developmentconexus-ops/conexus-methodology/METHOD.md`;
- canonical `developmentconexus-ops/conexus-methodology/REPOSITORY-STANDARD.md`.

Repository authority outranks this handoff. If this file conflicts with current routed authority, follow the repository and report the conflict.

---

## Candidate under review

D6 frontend falsification exposed a material gap: the accepted 95-operation / 29-Permission Product could operate listings, price, availability, sales, fulfillment, market comparison and economics, but it could not honestly support the strategy team's first-party marketplace-performance questions without the frontend inventing analytics authority.

The operator approved a bounded D0→D5 repair, now mechanically represented in the current candidate.

### Intended semantic change

Add one read/derive business boundary:

`Marketplace Performance Intelligence`

Its intended question is:

> **How is our own participation in an exact marketplace installation performing over an explicit period, with enough evidence semantics to support human investigation without stealing authority from the domains that own listing, price, availability, sales, competitive market meaning or economics?**

Performance may interpret:

- first-party exposure / traffic;
- engagement;
- provider-defined conversion evidence where semantically available;
- retail-media performance and efficiency;
- period-over-period performance when measurement bases are comparable;
- evidence sufficiency, coverage, freshness and provenance;
- bounded human-investigation meaning.

Performance must **not** own:

- Listing / ListingIntent / PriceIntent → Offering;
- sellable availability → Availability;
- canonical Sale truth → Marketplace Sales;
- competitive comparability / market position → Market Intelligence;
- profitability / margin / economic attribution → Commercial Economics;
- authorization → Governance;
- responsibility/work lifecycle → Operational Work;
- provider protocol/auth/DTOs → D4;
- campaign/bid/budget/targeting/creative execution → deferred;
- AI/MCP → future mechanism only, not current business authority.

The strategy workspace in D6 is intended to remain a read-only composition, **not** a `Strategy` business domain or `/dashboard` / `/analytics` Product API.

### Intended D2 historical-evidence rule

Provider retention must not silently become MPC historical retention.

Where external performance evidence can expire but is materially required for future trend/comparison/explanation, Marketplace Performance Intelligence may preserve the smallest sufficient **source-qualified historical observation/evidence** needed for its own claims.

This custody must not convert provider-reported evidence into MPC-authored external fact, create a generic Metric/Data Lake/Data Warehouse business authority, or justify indefinite raw-provider payload retention by default.

### Intended D3 communication rule

Baseline is Q/P-oriented:

- Performance consumes current/public meanings from Portfolio, Offering, Sales, Availability, Market Intelligence and Commercial Economics as needed;
- frontend Strategy Workspace is a read composition;
- no KPI event bus (`CtrDropped`, `RoasChanged`, etc.);
- no automatic Work generation merely because a metric changed;
- no Performance capability request that mutates another owner.

### Intended D4 Mercado Livre proof lane

Current admitted first-provider evidence is bounded to current official surfaces sufficient for Product 1.0 Performance, especially:

- Listing Visits;
- current Product Ads / retail-media performance evidence;
- advertiser namespace discovery/binding;
- current provider retention / refresh / attribution semantics;
- provider scope distinctions such as campaign / listing(item) / marketplace catalog grouping / marketplace family grouping.

Advertiser identity must not be equated with Marketplace Installation/seller identity by name, position or convenience. Multi-candidate binding must fail closed until an exact current candidate is explicitly selected/revalidated under the accepted technical ceremony.

Technical API access must not be confused with contractual permission to derive/publish every statistic. Review the current Mercado Livre developer terms / official documentation and decide whether the candidate's `contract_restricted` / capability-gate treatment is sufficient at architecture stage.

Do not introduce Ads management.

### Intended D5 Product surface

Current candidate adds exactly four H/A/S Q operations under one new ordinary Permission `performance.read`:

1. `GetMarketplacePerformanceSummary`
2. `ListMarketplaceListingPerformance`
3. `GetMarketplaceListingPerformance`
4. `ListRetailMediaPerformance`

Expected totals:

- 99 Product operations;
- 30 ordinary Permissions;
- 28 List/Search operations;
- Principal kinds H/A/S only;
- zero new Performance mutations;
- zero generic Product P/dashboard operations;
- no generic `/analytics`, `/metrics`, `/strategy`, metric-selector DSL, `group_by`, caller-sort, time-series baseline, `signals[]`, recommendation API, refresh/sync/collect API or Ads mutation.

Performance reads are scoped to an **exact Marketplace Installation** and explicit calendar reporting periods. Future providers must not be forced into false metric equivalence.

Retail Media evidence preserves closed scope variants instead of forcing campaign/family/catalog evidence into Listing identity.

Knowledge/coverage must preserve distinctions such as `complete`, `partial`, `unknown`, `unavailable`, `unsupported`; known zero/empty must remain distinct from unknown/unavailable.

Historical reads may be satisfied by preserved evidence and are not required to be live passthroughs of the provider.

### Intended frontend result

Primary navigation remains task-oriented. The strategic group is now:

```text
ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia
```

Performance requires exact Installation context. Market and Economics retain their own authorities.

The revised interaction map claims:

- 99/99 operation coverage;
- 40 route/screen states;
- 12 user flows;
- 32 negative frontend falsifiers.

The revised HTML proof is `pt-BR` and includes Performance Resumo / Publicações / Mídia plus Mercado and Economia. It is low-fidelity structural evidence, not final design or browser/runtime proof.

---

## Current executable evidence at handoff

Exact candidate `43bd8b2b9d126931ac488a11ebc311ef0f12b0eb` has current `gate:full` PASS.

The proof deliberately runs the accepted D5 baseline first and the D6-R1 repair second.

### Accepted D5 baseline non-regression

- OpenAPI 3.1.2;
- 95/95 Product operations;
- 29/29 ordinary Permissions;
- H/A/S only;
- stable origin `https://conexus.fun`;
- 14/14 idempotency carriers;
- 26/26 List/Search operations;
- 12/12 baseline OAD negative controls;
- generated projection semantics PASS;
- 405 / `Allow`, colon-suffix generated mux invocation and Organization privacy 404 proofs PASS;
- legacy runtime population 0.

### Current D6-R1 candidate

- 99/99 Product operations;
- 30/30 ordinary Permissions;
- H/A/S only;
- 28/28 List/Search operations;
- 7/7 Performance-specific negative controls;
- full TypeScript + Go generated projection semantics PASS;
- `product_oad_performance_repair=PASS`.

### Repository gate

- bootstrap 18825 / 20480 bytes;
- durable docs reachable 43;
- machine route files/selectors 0;
- legacy runtime population 0;
- repository negative controls 5/5;
- exact base...candidate diff has 14 files;
- gate PASS;
- PR title check PASS.

Do **not** accept these claims merely because CI is green. Challenge whether each proof is capable of falsifying the property it reports and whether preserving the 95/29 baseline verifier actually guards regressions rather than merely replaying a detached fixture.

---

# Required independent review

## 1. Global Product / platform coherence

Review D6-R1 against the already accepted Marketplace Central product direction, not in isolation.

Explicitly answer:

- Does Performance Intelligence materially belong in `Marketplace Central = Marketplace Operations Control Plane + Commercial Intelligence`, or has scope drift occurred?
- Is a **13th semantic boundary** genuinely justified by independent meaning/lifecycle, or should the meaning live in Market Intelligence / Commercial Economics / a read projection instead?
- Does the candidate preserve the distinction between own-participation Performance, competitive Market Intelligence and Commercial Economics?
- Does the strategy workspace remain frontend composition rather than a disguised new `Strategy` authority?
- Has the repair accidentally contradicted any accepted D0–D5 invariant, Product non-goal or authority edge?
- Is the bounded reopen truly the **smallest responsible reopen**, or did the candidate spread changes too widely?

If a contradiction exists, identify the **smallest owning authority** that must reopen. Do not call for a general rebaseline by preference.

## 2. DevelopmentConexus Method / Root Cause / Global Maximum / YAGNI

Perform an explicit method review.

Challenge whether the sequence actually follows:

`evidence → user need → invariant → alternatives → falsifier → decision → executable proof`

Review the rejected alternatives, especially:

- expand Market Intelligence into generic analytics;
- frontend-only dashboard composition over raw provider data;
- generic Analytics/Metric API;
- generic Metric entity/store;
- Strategy domain/API;
- Ads-management domain now;
- Data Lake / Data Warehouse now;
- event-per-KPI;
- time-series DSL now;
- `signals[]` / recommendation API now;
- AI/MCP now.

Decide whether the selected solution is the **global maximum under current evidence** rather than simply a locally convenient compromise.

Explicitly classify any issue as:

- UNDER-ENGINEERED;
- PROPORTIONATE / YAGNI;
- OVER-ENGINEERED;
- SPECULATIVE FUTURE ABSTRACTION.

Do not suggest extra layers/tools merely because they are common in analytics systems.

## 3. D1 semantic-boundary review

Attack `Marketplace Performance Intelligence` as a boundary.

Use concrete counterexamples:

- visits up, sales down;
- ROAS up, profitability down;
- ads active while availability is low;
- competitive price worsens while traffic rises;
- FAMILY/CATALOG media evidence with multiple Listings;
- a provider metric definition changes between periods;
- future Amazon `sessions` vs Mercado Livre `visits`.

Determine whether ownership remains singular and whether Performance avoids recomputing producer-owned facts/meaning.

## 4. D2 ownership / historical evidence / storage review

This is a material area and must receive a first-principles pass.

Challenge whether D2-R2 creates any accidental second authority or generic analytics store.

Review at least:

- why provider expiration makes durable custody necessary;
- minimum sufficient granularity / provenance / measurement basis;
- source qualification;
- historical-vs-current authority;
- retention without raw payload hoarding;
- PII/data-minimization impact;
- whether preserved evidence can honestly answer later periods after provider expiry;
- whether basis changes invalidate comparisons;
- whether the contract leaves D7 enough freedom to choose physical storage without under-specifying correctness.

Explicitly decide whether the statement **“MPC has memory, not an indiscriminate API dump”** is actually preserved by the authority and OAD.

Do not select Postgres schema, warehouse, ClickHouse, Snowflake, BigQuery, event store, Kafka or another D7 mechanism.

## 5. D3 communication review

Challenge Q/P-only baseline.

Determine whether any real Product 1.0 correctness requirement actually demands:

- a committed Performance event;
- autonomous Work creation;
- another owner reaction;
- durable occurrence communication beyond preserved Performance evidence.

If no real consumer exists, confirm that no KPI event bus is the correct YAGNI choice.

## 6. D4 provider / Retail Media review

Independently verify current official primary evidence, with publication/update dates where possible.

Challenge:

- Mercado Livre Visits semantics and coverage;
- Product Ads current API topology and 2026 migrations/deprecations;
- metric semantics (`prints`/impressions, clicks, CTR, CPC, CVR, ROAS, ACOS, TACOS, SOV, impression-share variants, organic/direct/indirect evidence where actually current);
- provider retention horizon;
- refresh/update timing;
- attribution window;
- `ITEM` / `CATALOG` / `FAMILY` or current equivalent scope semantics;
- advertiser discovery and whether `advertiser_id` can ever safely be inferred from seller/Installation;
- contractual terms around derived/performance statistics and whether architecture needs another gate.

Use official Mercado Livre developer documentation/terms as primary authority. Secondary sources may only supplement.

Also use current Amazon primary API documentation only as a **falsifier of false cross-provider equivalence**, not as a provider we are admitting now.

Explicitly decide whether the current seam makes a later Amazon/Shopee integration easier **without** pretending those providers are currently supported or creating a universal marketplace analytics interface.

## 7. D5 operation / Permission / wire-contract review

Challenge whether exactly four Qs and one Permission are the smallest sufficient Product API.

For each operation decide whether its consumer and meaning are real:

- `GetMarketplacePerformanceSummary`;
- `ListMarketplaceListingPerformance`;
- `GetMarketplaceListingPerformance`;
- `ListRetailMediaPerformance`.

Attack alternatives:

- 3 operations only;
- one giant `GetMarketplacePerformance`;
- `/dashboard` / `GetStrategyWorkspace`;
- N+1 listing-point reads instead of list summary;
- forcing Retail Media through Listing performance.

Review path grammar, exact Installation scoping, period semantics, comparison constraints, coverage unions, list population, retail-media scope union and absence of caller-controlled metric/sort/group DSL.

Check whether `performance.read` is the correct access boundary and whether `portfolio.manage` remains the correct authority for the **Technical Non-Product** advertiser-binding ceremony without inventing `performance.manage`.

Explicitly inspect the Product access views to confirm all clients can actually receive the new Permission while the accepted baseline 29-Permission proof remains meaningful.

## 8. OAD / executable-proof quality

Review the current 99/30 proof as a falsifier, not a test-count ceremony.

Challenge at least:

- exact 99 operation IDs, not count only;
- exact 30 Permission vocabulary;
- no fourth Principal kind;
- exact four Performance Qs / owner / Permission / H-A-S mapping;
- exact Installation-scoped paths;
- period parameters required on all four reads;
- comparison all-or-none semantics;
- collection paging + explicit semantic query repetition;
- no Performance POST/PATCH/PUT/DELETE;
- no generic analytics/metrics/strategy paths;
- no metric-selector/group-by/sort DSL;
- no baseline time-series or `signals[]`;
- no Ads mutations;
- retail-media scope branches remain closed/typed;
- exact-listing Performance cannot silently absorb FAMILY/CATALOG/CAMPAIGN evidence;
- historical custody does not become live-provider-only passthrough;
- generated TypeScript and Go preserve the new meaning;
- baseline non-regression proof actually demonstrates old D5 invariants still hold after subtracting only the approved repair.

Identify self-referential fixtures or negative controls that can pass while the protected property is broken.

Do not turn D5 generated-code proof into a claim of D7 runtime schema enforcement/router support.

## 9. Frontend / IA coherence review

Review the revised Portuguese D6 artifacts **only as consumers of accepted Product authority**.

Explicitly decide whether:

```text
ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia
```

is understandable and coherent for a real strategy/marketplace team without exposing bounded-context jargon.

Challenge:

- exact account/channel context for Performance;
- no misleading global cross-marketplace aggregated metric;
- analysis of an individual Listing is discoverable;
- Retail Media analysis is discoverable without looking like Ads management;
- Market and Economics remain visible as distinct strategic questions;
- known/unknown/partial/unavailable states are representable;
- source/basis warnings are proportional rather than noisy;
- the interface does not imply causal claims the Product cannot make;
- Portuguese copy is suitable as the primary operator language;
- route/button visibility is not treated as authorization;
- the wireframe does not imply time-series/signals/recommendations not admitted by D5.

Do not recommend final visual design, design-system selection, router/form library or frontend package topology in this review. Those are deliberately paused pending your verdict.

## 10. AI/MCP future-readiness without present scope creep

Review whether the accepted evidence model genuinely leaves a clean future path for AI to summarize/compare/investigate **without making AI current authority**.

Explicitly decide whether any current choice prematurely binds the system to MCP, agents, embeddings, vectors or model-specific schemas. Such coupling should be rejected unless already required by Product 1.0.

## 11. Reconstruction / continuation decision

Give an explicit answer:

> **Is there any material reason, under current repository authority and current primary evidence, to reconstruct the Marketplace Performance Intelligence approach or reopen broader accepted architecture before continuing D6 frontend topology?**

Allowed evidence-backed answers:

- `NO — candidate is coherent; continue D6 after operator adjudication`;
- `ACCEPT WITH BOUNDED FIXES`;
- `REOPEN <smallest authority>`;
- `RECONSTRUCT` only if a material falsifier proves the approach cannot preserve accepted Product meaning.

Preference for a different analytics product, dashboard style, newer technology, warehouse or fashionable AI architecture is not a reconstruction reason.

---

# Targeted falsifiers

At minimum, attempt to falsify these claims:

1. Performance is distinct from Market Intelligence.
2. Performance is distinct from Commercial Economics.
3. Strategy Workspace is not a business authority.
4. Performance cannot mutate Listing/Price/Availability/Ads.
5. `performance.read` does not imply Market/Economics/write Permissions.
6. advertiser namespace cannot be silently inferred from seller/name/first candidate.
7. campaign/family/catalog evidence cannot be attributed to one Listing without proof.
8. provider CVR/ROAS semantics cannot be reconstructed/relabelled by frontend convenience.
9. zero/empty cannot substitute for unknown/unavailable/partial.
10. cross-provider metric names do not imply equivalence.
11. provider retention does not erase MPC-required historical evidence.
12. historical custody does not make MPC the source of the original external fact.
13. raw provider payload retention is not the default.
14. no Data Warehouse/Data Lake is selected as architecture by anticipation.
15. no KPI event bus exists without a concrete autonomous consumer.
16. no automatic Work item arises merely from a metric change.
17. no generic Product Analytics/Metric/Strategy API exists.
18. Summary does not steal Sales/Market/Economics/Availability meaning.
19. Retail Media collection has a real consumer independent of Listing performance.
20. four Qs are necessary and sufficient; operation #100 is not hidden by symmetry.
21. exact periods and reporting timezone avoid ambiguous metric windows.
22. partial coverage prevents invalid period comparison.
23. measurement-basis change can prevent false comparison.
24. no time-series contract is accidentally implied by the wireframe.
25. no signals/recommendation semantics are accidentally implied by the wireframe.
26. AI/MCP remains future mechanism, not authority.
27. the new boundary makes a later provider easier without a generic provider analytics interface.
28. the 95/29 baseline verifier plus 99/30 wrapper is non-vacuous.
29. current Product implementation remains blocked until D9.
30. no D7 server/runtime/database/storage mechanism was selected by the repair.
31. Portuguese frontend information architecture remains task-oriented rather than domain-taxonomy-oriented.
32. the whole repair remains proportional to the actual strategic user need.

---

# Output contract

Append below `## Fable response` using this structure.

## 1. Verdict

Choose exactly one:

- `ACCEPT`
- `ACCEPT WITH BOUNDED FIXES`
- `REOPEN SMALLEST AUTHORITY`
- `REJECT / RECONSTRUCT`

## 2. Executive coherence assessment

In concise form answer:

- Is D6-R1 coherent with the previously accepted Marketplace Central platform?
- Is it consistent with the DevelopmentConexus Method?
- Is it proportionate/YAGNI?
- Does it preserve future strategic/AI value without current scope creep?

## 3. Material findings

Number findings highest severity first.

For each finding provide:

- **classification:** `D0_FIX`, `D1_FIX`, `D2_FIX`, `D3_FIX`, `D4_FIX`, `D5_FIX`, `D6_FIX`, `D7_OBLIGATION`, `LATER_NON_BLOCKING`, `AUTHORITY_CONTRADICTION`, or `REVIEW_FALSE_POSITIVE`;
- **severity:** Critical / Important / Minor;
- exact candidate location;
- governing repository authority;
- primary external evidence when provider/technology-dependent;
- concrete counterexample/failure;
- smallest correction;
- why it belongs in that exact authority/stage.

If there are no material findings, state that explicitly. Do not manufacture style findings.

## 4. Boundary/ownership assessment

Explicitly adjudicate Performance vs Market vs Economics vs frontend Strategy composition.

## 5. Historical-evidence / storage assessment

Explicitly adjudicate D2-R2 and whether the design is neither lossy nor a generic analytics warehouse.

## 6. Mercado Livre / Retail Media evidence assessment

List current primary sources checked and conclude on Visits, Product Ads, advertiser binding, scope, retention, refresh, attribution and contractual-use gates.

## 7. D5/OAD proof assessment

State what 99/30 proof genuinely establishes, what the baseline non-regression proof establishes, and any blind spots.

## 8. Frontend coherence assessment

Explicitly conclude whether the revised Portuguese IA/wireframes are coherent enough to proceed to topology after findings are adjudicated.

## 9. Method / YAGNI / Global Maximum assessment

Explicitly classify the architecture and rejected alternatives.

## 10. Reconstruction decision

Answer yes/no with evidence.

## 11. Continuation recommendation

State exactly what must happen before D6 frontend topology/dependency work resumes.

---

## Interaction rule

Fable writes **only** to this file on `review/d6-fable`. Do not edit PR #54 or candidate files. GPT will independently adjudicate every finding against repository authority and executable proof before any candidate change. No review output directly authorizes merge, phase progression, D7 work or Product implementation.

Round 2 occurs only if a material contradiction survives GPT adjudication and bounded fixes.

---

## Fable response

<!-- Fable: append independent review here. -->
