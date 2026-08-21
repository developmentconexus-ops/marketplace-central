# D6-R1 Marketplace Performance Intelligence — Independent Fable Review

> Review branch only: `review/d6r1-fable`
> Review PR: #56 — `docs(d6): independent Fable performance intelligence review`
> Candidate branch: `stage/d6-frontend`
> Candidate HEAD at handoff: `48abb857791769e21971a9049af3d9fc943617f5`
> Candidate PR: #54 — `docs(d6): open frontend authority stage`
> Base `main` at handoff: `9d2c81e175bc39ac388c9d8924ddad21f2a86480`
> Product candidate: 99 Product operations / 30 ordinary Permissions / Principal kinds H-A-S only
> Active runtime baseline: NONE
> D7–D9: BLOCKED
> Product implementation: BLOCKED UNTIL D9

## Review protocol

This is a **material independent architecture / Product-contract / frontend-coherence review** requested by the operator before any further D6 frontend topology or dependency work.

Do **not** optimize for agreement with GPT, the operator, prior chat reasoning, CodeRabbit, prior Fable reviews, or the current candidate. Reconstruct the reasoning independently from repository authority, executable evidence, current primary external evidence and the canonical DevelopmentConexus Method.

Reviewer output is **Evidence, not authority**. Do not change the candidate branch. Write only below `## Fable response` in this file on `review/d6r1-fable`.

Before analysis, revalidate:

1. remote `main`;
2. `stage/d6-frontend` exact HEAD;
3. PR #54 base/head/state/changed files;
4. CI on the exact candidate;
5. PR #56 and that this review branch differs from the exact candidate by **only** `docs/work/current/ai-dialog.md`.

If the candidate moved from `48abb857791769e21971a9049af3d9fc943617f5`, stop reviewing the stale SHA and record the new exact candidate SHA first.

The bounded candidate change from the previously proved `43bd8b2...` to `48abb857...` is **repository review infrastructure only**: `scripts/gate.ps1` restores the already-ratified `review/*` isolation property after PR #55 exposed that it had been lost. It does not alter Product/Performance semantics.

## Reading discipline

Start strictly:

1. `AGENTS.md`
2. `docs/index.md`
3. `docs/roadmap.md`
4. `docs/engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md`
5. `ARCHITECTURE.md`

Then switch bounded packs as concrete questions require them. Do not recursively ingest the whole rebaseline at once.

### Product / boundary pack

- `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
- `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`

### Identity / historical evidence / communication pack

- `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
- `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`

### External evidence pack

- `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
- only concrete current provider evidence needed to adjudicate Mercado Livre Visits / Product Ads / advertiser binding / retention / attribution / contractual admissibility

### Product API pack

- `contracts/api/product/openapi.yaml`
- `contracts/api/product/paths-access-performance.yaml`
- `contracts/api/product/paths-performance.yaml`
- D5 W1/W2/W3/W4 narrative authority only when a concrete OAD question requires it
- `scripts/verify-product-oad.mjs`
- `scripts/verify-product-oad-baseline.mjs`
- `scripts/gate.ps1`

### Frontend-consumption pack

- `docs/engineering/rebaseline/D6-FRONTEND.md`
- `docs/engineering/rebaseline/D6-B1-INTERACTION-MAP.md`
- `qualification/d6-wireframes/index.html`

For method/evidence discipline use:

- `docs/development/evidence-grounded-production-engineering-for-llm-agents.md` — derived, non-authoritative;
- `developmentconexus-ops/conexus-methodology/METHOD.md` — canonical;
- `developmentconexus-ops/conexus-methodology/REPOSITORY-STANDARD.md` — canonical.

Repository authority outranks this handoff. If this file conflicts with current routed authority, follow the repository and report the conflict.

---

# Candidate under review

D6 frontend falsification exposed a material gap: the accepted 95-operation / 29-Permission Product could operate listings, price, availability, sales, fulfillment, competitive market analysis and economics, but could not honestly support a strategy team's first-party marketplace-performance questions without the frontend inventing analytics authority.

The operator approved a bounded D0→D5 repair, mechanically represented in the current candidate.

## Intended semantic change

Add one read/derive business boundary:

`Marketplace Performance Intelligence`

Its intended question is:

> **How is our own participation in an exact marketplace installation performing over an explicit period, with enough evidence semantics to support human investigation without stealing authority from the domains that own Listing, Price, Availability, Sales, competitive Market meaning or Economics?**

Performance may interpret:

- first-party exposure / traffic;
- engagement;
- provider-defined conversion evidence when semantically available;
- retail-media performance and efficiency;
- period-over-period performance when measurement bases are comparable;
- evidence sufficiency, coverage, freshness and provenance;
- bounded human-investigation meaning.

Performance must **not** own:

- Listing / ListingIntent / PriceIntent → Offering;
- Sellable Availability → Availability;
- canonical Sale truth → Marketplace Sales;
- competitive comparability / market position → Market Intelligence;
- profitability / margin / economic attribution → Commercial Economics;
- authorization → Governance;
- responsibility/work lifecycle → Operational Work;
- provider protocol/auth/DTOs → D4;
- campaign/bid/budget/targeting/creative execution → deferred;
- AI/MCP → future mechanism only, not current business authority.

The D6 Strategy Workspace must remain a read-only composition, not a `Strategy` business domain or `/dashboard` / `/analytics` Product API.

## Intended D2 historical-evidence rule

Provider retention must not silently become MPC historical retention.

Where external performance evidence can expire but is materially required for future trend/comparison/explanation, Marketplace Performance Intelligence may preserve the smallest sufficient **source-qualified historical observation/evidence** needed for its own claims.

This custody must not convert provider-reported evidence into MPC-authored external fact, create a generic Metric/Data Lake/Data Warehouse business authority, or justify indefinite raw-provider payload retention by default.

## Intended D3 communication rule

Baseline is Q/P-oriented:

- Performance consumes current/public meanings from Portfolio, Offering, Sales, Availability, Market Intelligence and Commercial Economics as needed;
- Strategy Workspace is read composition;
- no KPI event bus (`CtrDropped`, `RoasChanged`, etc.);
- no automatic Work generation merely because a metric changed;
- no Performance capability request that mutates another owner.

## Intended D4 Mercado Livre proof lane

Current admitted first-provider evidence is bounded to current official surfaces sufficient for Product 1.0 Performance, especially:

- Listing Visits;
- current Product Ads / retail-media performance evidence;
- advertiser namespace discovery/binding;
- current provider retention / refresh / attribution semantics;
- provider scope distinctions such as campaign / listing(item) / marketplace catalog grouping / marketplace family grouping.

Advertiser identity must not be equated with Marketplace Installation/seller identity by name, position or convenience. Multi-candidate binding must fail closed until an exact current candidate is explicitly selected/revalidated under the accepted technical ceremony.

Technical API access must not be confused with contractual permission to derive/publish every statistic. Review current Mercado Livre developer terms / official documentation and decide whether the candidate's contractual/capability-gate treatment is sufficient at architecture stage.

Do not introduce Ads management.

## Intended D5 Product surface

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

## Intended frontend result

Primary navigation remains task-oriented. The strategic group is:

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

The HTML proof is `pt-BR` and includes Performance Resumo / Publicações / Mídia plus Mercado and Economia. It is low-fidelity structural evidence, not final design or browser/runtime proof.

---

# Current executable evidence at handoff

Exact candidate `48abb857791769e21971a9049af3d9fc943617f5` has current `gate:full` PASS and `pr-title` PASS.

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

The current candidate also restores the review-isolation gate property:

- candidate/main may contain no `docs/work/**`;
- a `review/*` branch may contain only `docs/work/current/ai-dialog.md` under that root;
- review candidate ref must exist;
- exact candidate tree must be free of review material;
- candidate→review diff must be exactly the one admitted dialogue file;
- `docs/work/**` is temporary, not durable routed authority;
- repository negative controls are now 6/6, including review isolation.

Do **not** accept any claim merely because CI is green. Challenge whether each proof can falsify the property it reports.

---

# Required independent review

## 1. Global Product / platform coherence

Review D6-R1 against already accepted Marketplace Central direction, not in isolation.

Explicitly answer:

- Does Performance Intelligence materially belong in `Marketplace Central = Marketplace Operations Control Plane + Commercial Intelligence`, or has scope drift occurred?
- Is a **13th semantic boundary** justified by independent meaning/lifecycle, or should the meaning live in Market Intelligence / Commercial Economics / a read projection instead?
- Does the candidate preserve own-participation Performance vs competitive Market Intelligence vs Commercial Economics?
- Does Strategy Workspace remain frontend composition rather than disguised `Strategy` authority?
- Has the repair contradicted any accepted D0–D5 invariant/non-goal/authority edge?
- Is the bounded reopen the **smallest responsible reopen**, or did it spread too widely?

If contradiction exists, identify the smallest owning authority that must reopen. Do not request a general rebaseline by preference.

## 2. DevelopmentConexus Method / Root Cause / Global Maximum / YAGNI

Challenge whether the sequence truly follows:

`evidence → user need → invariant → alternatives → falsifier → decision → executable proof`

Review rejected alternatives:

- expand Market Intelligence into generic analytics;
- frontend-only dashboard over raw provider data;
- generic Analytics/Metric API;
- generic Metric entity/store;
- Strategy domain/API;
- Ads-management domain now;
- Data Lake / Data Warehouse now;
- event-per-KPI;
- time-series DSL now;
- `signals[]` / recommendation API now;
- AI/MCP now.

Classify any issue as UNDER-ENGINEERED, PROPORTIONATE/YAGNI, OVER-ENGINEERED or SPECULATIVE FUTURE ABSTRACTION.

Do not recommend extra layers/tools merely because analytics products commonly have them.

## 3. D1 semantic-boundary review

Attack `Marketplace Performance Intelligence` using concrete counterexamples:

- visits up, sales down;
- ROAS up, profitability down;
- ads active while availability is low;
- competitive price worsens while traffic rises;
- FAMILY/CATALOG media evidence with multiple Listings;
- provider metric definition changes between periods;
- future Amazon `sessions` vs Mercado Livre `visits`.

Determine whether ownership stays singular and Performance avoids recomputing producer-owned meaning.

## 4. D2 historical evidence / storage review

Challenge D2-R2 first-principles:

- why provider expiration makes custody necessary;
- minimum sufficient granularity/provenance/measurement basis;
- source qualification;
- historical vs current authority;
- retention without raw-payload hoarding;
- PII/data minimization;
- whether preserved evidence can answer later periods after provider expiry;
- measurement-basis changes and comparison validity;
- whether D7 retains physical-storage freedom without correctness becoming underspecified.

Explicitly decide whether the authority truly preserves **“MPC has memory, not an indiscriminate API dump.”**

Do not select Postgres schema, warehouse, ClickHouse, Snowflake, BigQuery, event store, Kafka or another D7 mechanism.

## 5. D3 communication review

Challenge the Q/P-only baseline.

Determine whether any Product 1.0 correctness need actually requires a Performance event, autonomous Work creation, another owner reaction or durable occurrence communication beyond preserved Performance evidence.

If no real consumer exists, confirm that no KPI event bus is the correct YAGNI decision.

## 6. D4 Mercado Livre / Retail Media review

Independently verify current official primary evidence, including dates/update status when material.

Challenge:

- Mercado Livre Visits semantics/coverage;
- Product Ads current API topology and 2026 migrations/deprecations;
- current metric semantics (prints/impressions, clicks, CTR, CPC, CVR, ROAS, ACOS, TACOS, SOV, impression-share variants, organic/direct/indirect evidence where actually current);
- provider retention horizon;
- refresh/update timing;
- attribution window;
- `ITEM` / `CATALOG` / `FAMILY` or current equivalent scope semantics;
- advertiser discovery and whether `advertiser_id` can safely be inferred from seller/Installation;
- contractual terms around derived/performance statistics and whether architecture needs another gate.

Use official Mercado Livre developer documentation/terms as primary authority.

Use current Amazon primary API documentation only as a **falsifier of false cross-provider equivalence**, not as a provider being admitted now.

Explicitly decide whether the seam makes later Amazon/Shopee integration easier **without** pretending they are currently supported or creating a universal marketplace analytics interface.

## 7. D5 operation / Permission / wire review

Challenge whether exactly four Qs + one Permission are the smallest sufficient Product API.

For each operation assess its real consumer/meaning:

- `GetMarketplacePerformanceSummary`;
- `ListMarketplaceListingPerformance`;
- `GetMarketplaceListingPerformance`;
- `ListRetailMediaPerformance`.

Attack alternatives:

- 3 operations only;
- one giant `GetMarketplacePerformance`;
- `/dashboard` / `GetStrategyWorkspace`;
- N+1 point reads instead of Listing-performance list;
- forcing Retail Media through Listing Performance.

Review path grammar, exact Installation scoping, reporting periods, comparison constraints, coverage unions, list population, Retail Media scope union and absence of caller metric/sort/group DSL.

Check `performance.read` access boundary and whether `portfolio.manage` remains correct for **Technical Non-Product** advertiser binding without inventing `performance.manage`.

Inspect Product access views to confirm clients can receive the new Permission while the 29-Permission baseline proof remains meaningful.

## 8. OAD / executable-proof quality

Review 99/30 proof as a falsifier, not test-count ceremony.

Challenge at least:

- exact 99 operation IDs, not count only;
- exact 30 Permission vocabulary;
- no fourth Principal kind;
- exact four Performance Qs / owner / Permission / H-A-S mapping;
- Installation-scoped paths;
- periods required on all four reads;
- comparison all-or-none semantics;
- collection paging + semantic query repetition;
- no Performance POST/PATCH/PUT/DELETE;
- no generic analytics/metrics/strategy paths;
- no metric-selector/group-by/sort DSL;
- no time-series baseline or `signals[]`;
- no Ads mutations;
- Retail Media scope branches remain closed/typed;
- exact-Listing Performance cannot absorb FAMILY/CATALOG/CAMPAIGN evidence silently;
- historical custody is not live-provider-only passthrough;
- generated TypeScript/Go preserve meaning;
- baseline non-regression proof really demonstrates old D5 invariants after subtracting only the approved repair;
- restored review-isolation gate is non-vacuous and does not permit `docs/work/**` contamination in candidate/main.

Identify self-referential fixtures or controls that pass while the protected property is broken.

Do not turn D5 proof into a D7 runtime schema/router claim.

## 9. Frontend / IA coherence review

Review revised Portuguese D6 artifacts only as consumers of accepted Product authority.

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

is understandable/coherent for a real marketplace strategy team without exposing bounded-context jargon.

Challenge:

- exact account/channel context for Performance;
- no misleading global cross-marketplace aggregate metric;
- individual Listing analysis discoverability;
- Retail Media analysis discoverable without looking like Ads management;
- Market and Economics visibly distinct strategic questions;
- known/unknown/partial/unavailable states representable;
- source/basis warnings proportional rather than noisy;
- no unsupported causal claims;
- Portuguese as primary operator language;
- route/button visibility not authorization;
- wireframe does not imply time-series/signals/recommendations not admitted by D5.

Do not select final visual design, design system, router/form library or frontend package topology. Those are intentionally paused pending this verdict.

## 10. AI/MCP future-readiness without present scope creep

Review whether evidence structure leaves a clean future path for AI to summarize/compare/investigate **without making AI current authority**.

Flag any premature binding to MCP, agents, embeddings, vectors or model-specific schemas.

## 11. Reconstruction / continuation decision

Explicitly answer:

> **Is there any material reason, under current repository authority and current primary evidence, to reconstruct Marketplace Performance Intelligence or reopen broader accepted architecture before continuing D6 frontend topology?**

Allowed evidence-backed outcomes:

- `NO — candidate is coherent; continue D6 after operator adjudication`;
- `ACCEPT WITH BOUNDED FIXES`;
- `REOPEN <smallest authority>`;
- `RECONSTRUCT` only if a material falsifier proves the approach cannot preserve accepted meaning.

Preference for another analytics product, dashboard style, newer technology, warehouse or fashionable AI architecture is not a reconstruction reason.

---

# Targeted falsifiers

Attempt to falsify at minimum:

1. Performance is distinct from Market Intelligence.
2. Performance is distinct from Commercial Economics.
3. Strategy Workspace is not business authority.
4. Performance cannot mutate Listing/Price/Availability/Ads.
5. `performance.read` does not imply Market/Economics/write Permissions.
6. advertiser namespace cannot be silently inferred.
7. campaign/family/catalog evidence cannot be attributed to one Listing without proof.
8. provider CVR/ROAS semantics cannot be reconstructed/relabelled by frontend convenience.
9. zero/empty cannot substitute for unknown/unavailable/partial.
10. cross-provider metric names do not imply equivalence.
11. provider retention does not erase MPC-required history.
12. historical custody does not make MPC source of the original external fact.
13. raw provider payload retention is not the default.
14. no warehouse/data lake is selected by anticipation.
15. no KPI event bus exists without autonomous consumer.
16. no automatic Work arises merely from metric change.
17. no generic Product Analytics/Metric/Strategy API exists.
18. Summary does not steal Sales/Market/Economics/Availability meaning.
19. Retail Media collection has a real consumer independent of Listing Performance.
20. four Qs are necessary/sufficient; no hidden operation #100 by symmetry.
21. exact periods/reporting timezone avoid ambiguous windows.
22. partial coverage prevents invalid comparison.
23. measurement-basis change can block false comparison.
24. no time-series contract is implied accidentally.
25. no signals/recommendation semantics are implied accidentally.
26. AI/MCP remains future mechanism, not authority.
27. later providers become easier without generic provider-analytics interface.
28. 95/29 baseline verifier + 99/30 repair proof is non-vacuous.
29. Product implementation remains blocked until D9.
30. no D7 server/runtime/database/storage mechanism was selected.
31. Portuguese IA stays task-oriented rather than domain-taxonomy-oriented.
32. the entire repair is proportional to the actual strategic user need.

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

Answer concisely:

- Is D6-R1 coherent with previously accepted Marketplace Central?
- Is it consistent with DevelopmentConexus Method?
- Is it proportionate/YAGNI?
- Does it preserve future strategic/AI value without current scope creep?

## 3. Material findings

Number findings highest severity first. For each provide:

- **classification:** `D0_FIX`, `D1_FIX`, `D2_FIX`, `D3_FIX`, `D4_FIX`, `D5_FIX`, `D6_FIX`, `REPOSITORY_FIX`, `D7_OBLIGATION`, `LATER_NON_BLOCKING`, `AUTHORITY_CONTRADICTION`, or `REVIEW_FALSE_POSITIVE`;
- **severity:** Critical / Important / Minor;
- exact candidate location;
- governing repository authority;
- primary external evidence when provider/technology-dependent;
- concrete counterexample/failure;
- smallest correction;
- why it belongs in that exact authority/stage.

If no material findings, state that explicitly. Do not manufacture style findings.

## 4. Boundary/ownership assessment

Adjudicate Performance vs Market vs Economics vs frontend Strategy composition.

## 5. Historical-evidence / storage assessment

Adjudicate D2-R2 and whether design is neither lossy nor generic analytics warehouse.

## 6. Mercado Livre / Retail Media evidence assessment

List current primary sources checked and conclude on Visits, Product Ads, advertiser binding, scope, retention, refresh, attribution and contractual-use gates.

## 7. D5/OAD proof assessment

State what 99/30 proof genuinely establishes, what baseline non-regression establishes and any blind spots.

## 8. Frontend coherence assessment

Conclude whether revised Portuguese IA/wireframes are coherent enough to proceed to topology after findings are adjudicated.

## 9. Method / YAGNI / Global Maximum assessment

Explicitly classify the selected architecture and rejected alternatives.

## 10. Reconstruction decision

Answer yes/no with evidence.

## 11. Continuation recommendation

State exactly what must happen before D6 frontend topology/dependency work resumes.

---

## Interaction rule

Fable writes **only** to this file on `review/d6r1-fable`. Do not edit PR #54 or candidate files. GPT independently adjudicates every finding against repository authority and executable proof before any candidate change. No review output directly authorizes merge, phase progression, D7 work or Product implementation.

Round 2 occurs only if a material contradiction survives GPT adjudication and bounded fixes.

---

## Fable response

> Reviewed candidate: `stage/d6-frontend @ 48abb857791769e21971a9049af3d9fc943617f5` (revalidated on origin at review time).
> Review branch delta vs exact candidate: only `docs/work/current/ai-dialog.md` (verified by `git diff --stat` and by the repository gate in review mode).
> PR #54: OPEN / Draft, base `main`, head `stage/d6-frontend`; changed files match the declared bounded scope. CI on `48abb857`: `conventional-title` success, `required` success.
> Executable evidence reproduced locally on this review branch: `gate:full` PASS with `product_oad_baseline_non_regression=PASS` (95/95 · 29/29 · 26/26 · 12/12), `product_oad_operations=99/99`, `product_oad_permissions=30/30`, `product_oad_collection_operations=28/28`, `product_oad_performance_negative_controls=7/7`, TypeScript + Go projections deterministic/compilable, `review_mode: True`, `negative_controls: 6/6`, `gate: PASS`.
> Additional live falsifiers executed by this review (not only re-running green): (a) running the gate with a wrong base produced a red `review branch must differ from exact candidate only by docs/work/current/ai-dialog.md` — the restored review-isolation property fires on the real enforcement path, not only in its fixture control; (b) mutating `paths-performance.yaml` (`performance.read` → `market.read`) turned `verify-product-oad.mjs` red with the attributable message `GetMarketplacePerformanceSummary permission mismatch`, then the tree was restored clean. The proof falsifies the real artifact, not only an in-memory clone.
> Baseline verifier byte-identity confirmed at blob level: `9d2c81e:scripts/verify-product-oad.mjs` and `48abb857:scripts/verify-product-oad-baseline.mjs` are the same Git blob (`8618df75…`). The earlier-looking file difference on Windows is CRLF checkout only.

### 1. Verdict

**ACCEPT WITH BOUNDED FIXES**

No finding reopens D0–D5 semantics, the 13th boundary, D2-R2 custody, the four-Q surface, or `performance.read`. The fixes below are contract-precision, artifact-labeling and repository-hygiene corrections inside the already-approved bounded scope.

### 2. Executive coherence assessment

- **Coherent with accepted Marketplace Central?** Yes. "Marketplace Operations Control Plane + **Commercial Intelligence**" already names the intelligence half of the product; first-party participation performance was a genuine hole between Market Intelligence (explicitly *external comparable-market* observations per D1 §catalog) and Commercial Economics (economic interpretation). Neither accepted owner could honestly answer "how is our own participation performing?" — the frontend would have had to invent analytics authority, which D6 invariant 1 forbids. The repair closes the hole at the correct altitude.
- **Consistent with the DevelopmentConexus Method?** Yes. The chain evidence → user need → invariant → alternatives → falsifier → decision → executable proof is genuinely present: the D6 falsifier exposed the gap, the D6-R1 document records invariants and rejected alternatives, and the proof executes against the real OAD with authority-anchored expectations (W4 matrix + D6-R1 matrix parsed from the owning documents, not from the artifact under test).
- **Proportionate / YAGNI?** Yes. Four Qs, one Permission, zero mutations, zero events, no Metric entity, no time-series contract, no warehouse, no Ads management, no AI/MCP. The one place the candidate is arguably *under*-engineered is schema-level enforcement of the zero-vs-unknown law (finding F-2); the one place it leaks forward is a wireframe card implying a signals surface that D5 deliberately does not admit (finding F-1).
- **Future strategic/AI value without current scope creep?** Yes. Closed, source-qualified, provenance-carrying evidence with explicit coverage/measurement-basis is exactly the substrate a future AI/MCP consumer needs, and nothing in the contract or docs binds to MCP, agents, embeddings or model schemas today.

### 3. Material findings

**F-1 — wireframe "Sinais para investigar" card implies a Product signals surface that D5 does not admit**
- classification: `D6_FIX` · severity: **Important**
- location: `qualification/d6-wireframes/index.html` (Performance / Resumo, card "Sinais para investigar": "Porcelanato A: visitas ↑, vendas ↓, preço comparável acima"; "Cuba B: ROAS alto e perda de impressão por orçamento")
- authority: D6-R1 §10 ("No `signals[]`, recommendation score or AI explanation schema is admitted yet"), D6-B1 falsifier 17, handoff falsifier 25; D6-FRONTEND §3.2 ("frontend never manufactures … causal claims").
- counterexample/failure: no admitted Product operation returns cross-owner investigation signals. The "Porcelanato A" line composes Performance traffic, Sales activity and Market price position into one ranked pointer. If D6-B1 is ratified with this card unqualified, frontend topology work inherits a screen whose data source does not exist, creating direct pressure to either invent client-side analytics authority or prematurely add a signals API — precisely the two rejected alternatives.
- smallest correction: annotate the card as client-side read composition of already-admitted owner reads under their own Permissions (no ranking, no derived score), or remove it until a real consumer justifies a bounded signals/series read. Rename away from a bare "Sinais" framing if kept.
- why this stage: the contract is untouched; only the D6 candidate proof artifact needs correction, before operator adjudication of D6-B1.

**F-2 — schema permits `coverage: complete` with all measures absent, weakening the zero-vs-unknown law at the wire**
- classification: `D5_FIX` · severity: **Important**
- location: `contracts/api/product/paths-performance.yaml` — `TrafficPerformance`, `SalesActivityPerformance`, `RetailMediaSummary`, `RetailMediaPerformance` (measures optional while `coverage` may be `complete`)
- authority: D6-R1 §9 ("Known zero/empty remains distinct from unknown/unavailable"), negative control 10; `ARCHITECTURE.md` proof bar ("invalid value combination → type/constructor/schema failure").
- counterexample/failure: a server can emit `{coverage: {state: complete}}` with no `visits` and remain contract-valid; generated TypeScript/Go clients then face an absent-but-complete value whose meaning (zero? omitted?) is undefined, which is exactly the collapse the law forbids. The candidate already solved this class structurally for `RetailMediaPerformanceCollection` (Available/Unavailable `oneOf` split), proving the pattern is available and cheap.
- smallest correction: apply the same split to the four evidence carriers (available branch: `complete|partial` coverage with required measures; unavailable branch: `unknown|unavailable|unsupported` with measures forbidden), or make the measure fields required alongside a `complete` coverage branch.
- why this stage: it is wire-contract expressiveness owned by the D6-R1 D5 amendment itself; fixing it later would be a D5 reopen instead of finishing the bounded repair correctly.

**F-3 — unmarked baseline-fixture duplicates of the access surface live inside the canonical contract source tree**
- classification: `REPOSITORY_FIX` · severity: Minor
- location: `contracts/api/product/paths-identity-portfolio-readiness.yaml` — retained `AccessContext`/`OrganizationMembers`/`AccessRoles`/`AssignAccessRole`/`RevokeAccessRole` pathItems plus `*View` schemas bound to the 29-value `components.yaml#/schemas/Permission`, no longer referenced by `openapi.yaml` (which now points at `paths-access-performance.yaml` with the 30-value `PermissionView`).
- authority: D5-B1/ARCHITECTURE ("OpenAPI is the single machine-readable Product API wire authority; a hand-written second wire authority is not target architecture"); prior D5 review lesson that Redocly bundling prunes orphans, making source-text presence misleading.
- counterexample/failure: the duplicates exist so `baselineProof()` can remap refs and rebuild the 95/29 projection — a legitimate purpose — but nothing in the file says so. They carry duplicate operationIds and a superseded Permission vocabulary; a future author re-referencing them (or a source-scanning tool counting them) silently resurrects a 29-Permission access surface.
- smallest correction: fence the retained access pathItems/schemas with explicit comments declaring them baseline-non-regression fixtures consumed only by `scripts/verify-product-oad.mjs`, never referenced by the canonical entrypoint.
- why this stage: repository hygiene of the current candidate; no semantics change.

**F-4 — D2-R1 §5 asserts "D5 retains the same 95 Product operations and 29 ordinary Permissions" inside a 99/30 candidate**
- classification: `D2_FIX` (document precision) · severity: Minor
- location: `docs/engineering/rebaseline/D2-R1-PRESENTATION-IDENTITY.md` §5
- counterexample/failure: both amendments land in the same candidate; a fresh actor routed to D2-R1 reads a false current count and can cite it against the proved 99/30 surface.
- smallest correction: rephrase to "D2-R1 itself admits no new Product operation, ordinary Permission, Principal kind or domain", removing the absolute counts.

**F-5 — repository negative controls #1–#5 are tautological fixtures; the 6/6 count overstates falsification power**
- classification: `REPOSITORY_FIX` · severity: Minor
- location: `scripts/gate.ps1` `Expect-Failure` blocks ('legacy runtime root' … 'bootstrap overflow'): each tests a literal against itself (e.g. `'apps/example/main.go' -match '^apps/'`), so they pass regardless of whether the corresponding gate predicate exists or works. Only the restored 'review isolation' control exercises real logic, and it exercises `Test-ReviewDiffNames` the function, not the review-mode enforcement path (that path is, however, genuinely exercised by CI env wiring and was demonstrated red in this review with a wrong base).
- authority: repository proof discipline — presence is not execution; material guards require a deterministic falsifier.
- smallest correction: either make each control mutate a real input consumed by the real predicate, or stop counting the five tautologies in the advertised negative-control total. Not blocking: the material controls for this candidate (the 12 baseline + 7 performance OAD controls and the review-isolation enforcement) are real.

**F-6 — Retail Media scope/metric availability is not uniform across provider aggregation levels; realization must not assume it**
- classification: `D7_OBLIGATION` (with one D4 naming note) · severity: Minor
- primary evidence (checked 2026-08-21, official Mercado Libre developer documentation, page last update 30/12/2025, `api-version: 2`): the impression-share family (`impression_share`, `top_impression_share`, `lost_impression_share_by_budget`, `lost_impression_share_by_ad_rank`, `acos_benchmark`) is exposed only on the **campaign-detail** metrics endpoint, not on campaign-search or ads(item)-level metrics; ads-level metrics omit `ctr`/`cvr`/`roas`/`sov` from the documented value list; metrics history is bounded to **90 days backward**; metrics refresh daily at **10:00 GMT-3**; `acos_target` remains visible in campaign metrics responses only until **2026-03-30** (provider deprecation in motion — a live example of `basis_revision`-class change).
- counterexample/failure: an acquisition design assuming every `RetailMediaMeasures` field exists at every scope would fabricate unsupported evidence; the contract's per-field-optional + coverage design already tolerates this, but the D7 lane must prove which measures materialize per scope instead of assuming symmetry. Naming note for D4 realization: the contract's owner-local `lost_impression_share_by_rank` maps to provider `lost_impression_share_by_ad_rank`; record the mapping explicitly so the semantic rename is deliberate, not drift.
- also: `marketplace_catalog_group`/`marketplace_family_group` scopes are grounded in real provider ontology (Catalog `parent_id`, User Product `family_id`, catalog listings participating in Product Ads) but the current metrics API documents only campaign and item aggregation; which grouped evidence actually materializes is a D7 acquisition proof, and until proven the honest answer for those scopes is `unsupported`/`unavailable`, which the contract already expresses.

**F-7 — wireframe never exercises `not_comparable` / `insufficient_evidence` comparison states**
- classification: `D6_FIX` · severity: Minor
- location: `qualification/d6-wireframes/index.html` — only the positive "comparável" outcome appears; D6-B1 §7.2 requires `insufficient_evidence` and `not_comparable` as visible states, and handoff falsifiers 22–23 target exactly the blocked-comparison path.
- smallest correction: add one example row/panel showing a blocked comparison before operator adjudication, so the low-fi proof demonstrates the honest-refusal state, not only the happy path.

**F-8 — routing freshness: the D1 route row and a dangling status pointer can hide the 13th boundary from a fresh actor**
- classification: `REPOSITORY_FIX` · severity: Minor
- location: `docs/index.md` route row "Domains, semantic owners, allowed ownership edges" → D1 only, while D1 §13 closes with "exactly these 12 business boundaries"; and `ARCHITECTURE.md` §Current stage still says "Read `docs/README.md` for the sole current status" — `docs/README.md` does not exist (pre-existing, not introduced by this candidate).
- counterexample/failure: a fresh actor asking "what are the semantic owners?" follows the route, reads "exactly 12 … No Product 1.0 semantic remains", and never reaches D6-R1's 13th boundary; the amendment-as-R-doc convention (D4-R1, D2-R1 precedent) is fine, but this specific route row answers the exact question the amendment changes.
- smallest correction: append "plus the D6-R1 bounded 13th boundary" to the D1 route row, and point ARCHITECTURE's status line at `docs/roadmap.md`.

No other material findings. I deliberately raise no style findings.

### 4. Boundary/ownership assessment

The 13th boundary is justified by independent meaning and lifecycle, not by symmetry or fashion:

- **Performance vs Market Intelligence:** D1 defines Market Intelligence as *external comparable-market* observation. First-party visits/engagement/retail-media performance is a different subject (our participation), different evidence custody (provider-reported first-party surfaces with finite retention), different failure modes (advertiser binding, attribution basis). Merging it into Market would have broken D1's own definition. Counterexamples hold: "visits up, sales down" needs no shared owner — traffic evidence is Performance, canonical Sale truth stays Sales; "competitive price worsens while traffic rises" composes Market + Performance in the frontend without either recomputing the other.
- **Performance vs Commercial Economics:** "ROAS up, profitability down" is the decisive case — ROAS is provider-reported media efficiency (Performance evidence, `provider_reported` basis), profitability is Economics' L0/L1/L2 meaning. The contract never lets Performance emit margin/profit, and Economics never reinterprets provider media metrics. The existing `GetEconomicPerformanceSummary` (Economics) and the new `GetMarketplacePerformanceSummary` (Performance) answer different questions with different owners; the Portuguese IA keeps them in visibly different homes (Economia vs Performance).
- **Sales-activity counts inside Performance** deserve one explicit note: `SalesActivityPerformance` is provider-derived activity evidence with `PerformanceProvenance`, not canonical Sale truth. This is the one seam where a careless frontend could present Performance `sales_count` as the Sales number. The contract carries the provenance to prevent it; D6 should keep labeling the basis (the interaction map's §7 laws cover this — keep them binding at topology time).
- **Strategy Workspace:** remains pure client composition. No `/strategy` path, schema, or Permission exists (verified in the bundled OAD and by negative control), and the interaction map's S40–S42 consume only owner-native Qs. Not a disguised domain.
- **Ownership edges:** all six feed-forward edges into Performance are Q-only; no reverse authority edge exists; `performance.read` appears in exactly the four Performance operations (verified against the bundle).

### 5. Historical-evidence / storage assessment

D2-R2 is necessary and correctly minimal. Provider evidence: Visits queries are limited to **150 days** and Product Ads metrics to **90 days backward** — so "MPC has memory" is not optional if the Product must answer "compare this quarter to last quarter" a year from now; without custody the contract's historical claims would silently decay into `unavailable` at provider whim.

The custody law is well-bounded: smallest sufficient source-qualified evidence, original authority preserved (`evidence_custody: preserved_source_evidence` vs `current_source_observation` is on the wire), no Data Lake/Warehouse/Metric store/raw-payload archive admitted, PII surface effectively nil (aggregated counts/rates carry no buyer identity), and physical persistence/granularity/retention/indexing explicitly left to D7. The contract exposes exactly what correctness requires of any future storage choice — coverage with covered periods, provenance with `acquired_at`/`observed_through`, measurement basis with optional `basis_revision` — without choosing a mechanism. Measurement-basis change blocking false comparison is representable (`not_comparable`), and the provider's own in-flight `acos_target` deprecation (visible until 2026-03-30) is a concrete real-world event this design survives. Neither lossy nor a warehouse. PASS, subject to F-2 tightening the wire representation.

### 6. Mercado Livre / Retail Media evidence assessment

Primary sources checked (2026-08-21, current official Mercado Libre developer documentation; the pt-BR mirror blocks non-interactive fetch, the es_AR/en_US mirror served the full content in the browser):

- Product Ads (Mercado Ads guide, last update 30/12/2025, `api-version: 2`): advertiser discovery `GET /advertising/advertisers?product_id=PADS|DISPLAY|BADS`; campaign search + metrics; campaign detail + extended metrics; ads(item) metrics; daily aggregation.
- Visits surface: item visits with `date_from`/`date_to`, unique-per-day counting, availability within 48 hours, query range limited to 150 days.

Conclusions:

- **Visits:** semantics (unique VIP views/day), freshness (≤48 h) and retention (150-day window) match the candidate's evidence-semantics framing; a `visits` count with coverage/freshness/provenance is the honest representation. PASS.
- **Product Ads:** the admitted metric families are all real, current provider vocabulary — `prints/clicks/ctr/cpc/cvr/roas/sov/acos` plus organic and direct/indirect attributed quantities, and the impression-share family including budget/rank loss. The candidate correctly excludes `acos_benchmark` (a provider comparative benchmark — Market-shaped, not first-party performance) and excludes `acos_target`/`strategy`/`budget` (campaign *management* state, correctly out of scope). PASS with the F-6 obligations (per-scope availability asymmetry, 90-day window, daily refresh, `by_ad_rank` naming map).
- **Advertiser binding:** discovery returns a *list* spanning sites (`site_id` MLB/MLM/MLA/MLC in one response), and ads carry `delegated`/`revoked` statuses proving advertiser authority over items is transferable and *not* seller identity. The candidate's fail-closed multi-candidate ceremony under current human + `portfolio.manage` + exact Installation is not defensive fiction — it matches observable provider reality. Keeping it Technical Non-Product with no new Permission is correct; inventing `performance.manage` for a configuration ceremony would have been vocabulary creep. PASS.
- **Attribution:** the docs distinguish direct vs indirect attributed quantities; the exact window is provider-defined and not stated on the metrics page — the contract's optional `attribution_window_days` + opaque `basis_revision` is the right architecture-stage posture (claim only what the source states). PASS.
- **Contractual gates:** `contract_restricted` exists as a first-class unavailable reason and D4 is required to fail honest where derived/exposed use is restricted. At architecture stage that is sufficient; the concrete Developer Program Terms adjudication for each derived statistic belongs to the D7 acquisition lane where the actual calls and exposures are chosen. No additional gate needed now.
- **Cross-provider seam:** `marketplace_kind` enum pinned to `mercado_livre`, per-provider measurement basis mandatory, no generic provider-analytics interface, no cross-Installation aggregate. A later Amazon integration adds an enum value and its own basis — easier, without any pretense of current support or metric equivalence (Amazon "sessions" ≠ ML "visits" stays unrepresentable as equality). PASS.

### 7. D5/OAD proof assessment

What the proof genuinely establishes:

- The 99-operation/30-Permission surface is checked against **authority documents** (W4 matrix + D6-R1 matrix), by exact operation-ID set, per-operation class/Permission/principal-kinds, Installation-scoped paths, required periods, optional paired comparison parameters, canonical `limit`/`cursor`, 28 List/Search, closed coverage states, closed Retail Media scope kinds, integer-string counts, unit-carrying percentages/multiples, `provider_reported` basis, closed `RetailMediaMeasures`, and absence of generic analytics/metrics/strategy paths, metric-selector vocabulary and generic Metric schemas — all on the bundled real OAD, with deterministic double-bundle, deterministic TypeScript and compiled Go projections.
- The 7 performance negative controls mutate the bundled document and prove the validator would catch each drift class. This proves the *checker*, which is the correct falsifier for a static contract gate; my additional source-level mutation (permission swap → attributable red) closes the remaining gap by proving the pipeline end-to-end on the real files.
- Baseline non-regression is real, not ceremonial: the accepted 95/29 verifier is blob-identical to the D5 closeout verifier and is executed against a projection built from the *current* source tree with only the marker-fenced performance block removed and access refs remapped. A performance-ish path added outside the markers would land in the baseline projection and turn the old verifier red (96 ops). The subtraction removes exactly the approved repair; the D2-R1 view change deliberately remains inside the baseline, which is correct because D2-R1 is itself operator-approved.
- The restored review-isolation gate is non-vacuous on the enforcement path (demonstrated red with a wrong base; CI wires `GATE_BASE_REF`/`GATE_CANDIDATE_REF` from the PR so review PRs exercise it for real).

Blind spots, stated honestly: comparison all-or-none/equal-length/no-overlap and reporting-calendar semantics are prose + 422 declarations — OpenAPI cannot express cross-parameter constraints, so their enforcement is a D7 runtime obligation (correctly reported as `NOT_CLAIMED_D7`); the zero-vs-unknown law is prose-level where F-2 applies; and the five tautological repository controls (F-5) contribute count, not falsification. None of these blind spots is concealed by the proof output.

### 8. Frontend coherence assessment

The Portuguese IA is coherent for a real strategy team and free of bounded-context jargon: `ESTRATÉGIA E INTELIGÊNCIA → Performance (Resumo / Publicações / Mídia), Mercado, Economia` reads as three different questions (how are we performing / what is the market doing / what does it mean economically), which is exactly the D1 separation rendered as user language. Exact-Installation context is structural in every Performance screen; no global cross-marketplace KPI exists; individual Listing analysis is discoverable both from Performance/Publicações and as a Performance tab on the Listing detail (S22) without moving Listing ownership; Mídia is analysis-only with no management verbs; route/button visibility is explicitly usability-only; the 99/99 coverage table sums correctly (verified: 5+6+5+12+9+3+4+11+7+3+7+17+3+7 = 99) with no screen-shaped operation.

The wireframes are sufficient low-fi structural proof **after F-1 and F-7**: the signals card must stop implying an un-admitted Product surface, and the blocked-comparison states must appear at least once. With those corrected, proceed to operator adjudication and then topology.

### 9. Method / YAGNI / Global Maximum assessment

Selected architecture: **PROPORTIONATE/YAGNI**. One read/derive boundary, four Qs, one Permission, evidence custody scoped to its own claims, Q/P-only communication.

Rejected alternatives, classified:
- expand Market Intelligence into generic analytics — correctly rejected; would break D1's own external-market definition (AUTHORITY_CONTRADICTION if taken).
- frontend-only dashboard over raw provider data — correctly rejected; violates D6 invariant 1 and provider-DTO containment.
- generic Analytics/Metric API, Metric entity/store, time-series DSL, `signals[]`/recommendation API — correctly rejected as SPECULATIVE FUTURE ABSTRACTION; no current consumer, and each would have created generic authority the negative controls now forbid.
- Strategy domain/API — correctly rejected; composition suffices (proven by the interaction map needing zero new operations).
- Ads-management domain now — correctly rejected; analysis has a consumer, management does not, and management would drag budget/bid mutation authority plus a much heavier D4 write contract.
- Data Lake / Warehouse now — correctly rejected; D2-R2 states the correctness requirement and leaves mechanism to D7 (Global Maximum: decide meaning at the highest stable point, defer mechanism).
- event-per-KPI — correctly rejected; no independent consumer reaction exists, so E-grammar (D3 §E: real independent consumer) is unmet by definition. Q/P-only is the correct YAGNI decision, not an omission.
- AI/MCP now — correctly rejected; the evidence substrate is AI-ready precisely because it is honest and source-qualified, with zero present binding.

Under-engineering: only F-2 (schema expressiveness below the repository's own proof bar). No over-engineering found — the closed measure taxonomies mirror proven provider vocabulary rather than inventing abstraction.

### 10. Reconstruction decision

**NO.** No material falsifier survives against the boundary choice, the custody rule, the Q/P baseline, the four-Q surface, or the proof method. All 32 targeted falsifiers were attempted; none broke the candidate structurally (F-1/F-2/F-7 are exactly the bounded residue of falsifiers 25, 9 and 22/23, each with a smallest in-scope correction). Current primary provider evidence supports, rather than contradicts, the admitted semantics. Reconstructing or broadly reopening accepted architecture would discard proven meaning for preference — prohibited by the review contract and unjustified by evidence.

### 11. Continuation recommendation

Before D6 frontend topology/dependency work resumes:

1. GPT adjudicates F-1..F-8 against repository authority; apply the accepted bounded fixes on `stage/d6-frontend` (expected: F-2 schema split + regenerated projections re-proved by the existing gate; F-1/F-7 wireframe corrections; F-3/F-4/F-5/F-8 hygiene edits; F-6 recorded as D7 obligations in the D6-R1/D4 lane notes).
2. Operator reviews/adjudicates the corrected D6-B1 interaction map + Portuguese wireframe proof (the roadmap's stated next action) — reviewer output here is Evidence, not that adjudication.
3. Only after operator approval, open frontend topology/dependency adjudication under the D6 research subpack; no D7–D9, no Product implementation, no PR #54 merge without explicit operator authorization.

— Fable, 2026-08-21, `review/d6r1-fable`
