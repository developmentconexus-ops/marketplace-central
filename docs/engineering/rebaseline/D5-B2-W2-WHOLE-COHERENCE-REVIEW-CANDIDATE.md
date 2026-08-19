# D5-B2 — Whole-W2 Global Coherence Review Candidate

> **Status:** NON-AUTHORITATIVE LEAD REVIEW CANDIDATE  
> **Review subject:** accepted-in-stage W2-A/B/C/D/E as one schema/wire system  
> **Authority:** none — review evidence only until operator ratification and canonical filing  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Prepared:** 2026-08-18

## 1. Review scope and outcome

This candidate records the lead Whole-W2 adversarial review. It does not modify accepted W2-A/B/C/D/E by implication.

**Lead outcome:** `RESTRUCTURE NOW — W2-local corrections required; no parent-stage reopen currently proven.`

The local W2 sections are individually coherent, but global review found several places where history, admitted-operation coverage, HTTP precondition semantics and real provider requirement shapes are not yet globally closed.

Do not invoke Fable or advance to later Wire Contract batches until the operator ratifies the lead correction direction and the bounded review package is ready for independent challenge.

## 2. Evidence / known / unknown

### Known from accepted authority

- D4-R1 requires a material historical ListingIntent snapshot before consequential provider dispatch, preserving exact desired Offering values, source/override provenance, provider requirement revision, media selection, decision/authorization context, other-owner correlations and provider attempt/result/convergence evidence.
- W2-B currently exposes current resolved FOLLOW_SOURCE values and says correlation is preserved in historical effect/freeze context, but it does not define the concrete historical dispatch-basis schema.
- The ratified operation matrix admits `List/GetMarketplaceListings`, `List/GetFulfillmentNodes` and `GetEconomicPerformanceSummary`.
- W2-E claims all admitted Product families have a W2 schema home, but W2-A→E do not yet define concrete schemas for those admitted Qs/resources.
- W2-E correctly distinguishes direct-resource `If-Match` from referenced-resource preconditions.
- W2-E media intake currently uses `POST .../listing-intents/{id}/media` while applying the ListingIntent ETag as `If-Match`, although the selected HTTP target is the media collection rather than the ListingIntent resource.
- Current Mercado Livre official documentation exposes publication attribute types including `number_unit`; it also documents explicit N/A attribute semantics and a Brazil-available conditional-required endpoint whose evaluation request contains the concrete listing payload, including fields such as title, price, available quantity, listing type and description.

### Unknown / deferred

- Exact future marketplace/provider requirement types beyond current proven classes remain unknown.
- Exact D7 storage/hash/CDN implementation remains deferred.
- Whether a real selected Product 1.0 category exercises conditional-required attributes at the first D8 write remains to be measured; the provider contract nevertheless makes the failure class reachable for supported marketplace publication.
- Exact production need for multiple parallel Fulfillment execution scopes remains unproven.

## 3. Root cause

W2 was intentionally derived owner-by-owner. That preserved semantic boundaries, but local derivation left three global defect classes reachable:

1. **current-state schemas without sufficient historical effect basis**;
2. **family-level coverage assertions masking admitted operations whose concrete read/resource schema was never actually crystallized**;
3. **cross-cutting mechanics applied locally in a way that violates their own global rule or underfits the selected provider's real dynamic requirement contract**.

The correction must remain W2-local unless implementation/evidence proves an accepted parent semantic dependency is genuinely missing.

## 4. Target invariant

> **Every admitted Product operation has one faithful wire/schema home; every consequential publication remains historically explainable after source/provider change; cross-cutting HTTP safety rules apply to the actual selected request resource; and current provider requirement variability is representable without importing provider DTO/rules ontology or moving owner authority.**

## 5. F-W2-G1 — REVISE: historical publication dispatch basis is under-specified

### Finding

W2-B preserves current desired/resolved values, current candidate/option references and broad effect/convergence axes, but does not define the immutable historical basis required by D4-R1 §11.

A context-bound `source_candidate_key`, `option_key` or current `requirements_revision` is insufficient by itself after source/provider schema changes. A later GET must be able to explain what material Offering values/media/requirement revision and owner correlations were actually used for a consequential provider attempt without re-resolving history from current state.

### Corrected direction

Add an **owner-local immutable publication dispatch/attempt basis** inside ListingIntent history/effect meaning, without creating a new canonical Product resource or generic Operation.

For each consequential publication attempt, preserve proportionately:

- ListingIntent identity + exact submitted/material revision;
- source-qualified Product / Installation / target Listing;
- exact resolved Offering-owned PublicationValues actually used, including resolved FOLLOW_SOURCE value + source provenance and EXPLICIT_OVERRIDE provenance;
- provider requirement/schema revision materially used;
- media selection/order/role + material source/authored provenance;
- current PriceIntent identity/revision used, never an embedded ListingIntent-owned price;
- Availability-issued input/intent correlation used, preserving owner distinction;
- decision/disposition/AuthorizationDecision references materially used;
- intended/authorized/attempted scope where material;
- provider attempt/member/aspect outcomes and authoritative convergence evidence.

This is historical context, not a payload archive/PIM mirror. No standalone `PublicationAttempt` Product CRUD/resource is admitted merely for uniformity.

### Reopen trigger

If the historical basis cannot be represented without ListingIntent owning Price/Availability current meaning or without a genuinely new semantic identity, targeted D1/D2/D4-R1 review is required. No such need is currently proven.

## 6. F-W2-G2 — REVISE: “complete W2 family coverage” is currently false

### Finding

The matrix admits operations whose concrete W2 schema is absent or only named at family level:

1. `List/GetMarketplaceListings` — no concrete `MarketplaceListing` actual-state schema exists in W2-B;
2. `List/GetFulfillmentNodes` plus create/update/deactivate — W2-D references FulfillmentNode but does not define its own read/create/update schema discipline;
3. `GetEconomicPerformanceSummary` — W2-C defines ExpectedEconomics/SaleEconomics but no bounded performance-summary Q schema.

W2-E §15 therefore overclaims completeness.

### Corrected direction

Add the smallest missing schema homes before Whole-W2 can pass:

#### MarketplaceListing

Source-qualified Listing Q, no synthetic MPC Listing ID. Preserve only Offering-owned/provider-observed meaning required by Product 1.0, proportionately:

- qualified Listing ref;
- normalized observed lifecycle;
- publication context;
- bounded observed listing representation values/media relevant to current Offering semantics, using accepted PublicationValue/knowledge grammar rather than raw provider attributes;
- observed current marketplace price as source-qualified Offering evidence where needed for PriceIntent/convergence, never desired price authority;
- source freshness/provenance and representation convergence evidence;
- no Sellable Availability owner value/provider DTO mirror.

#### FulfillmentNode

Closed MPC-owned resource/read/create/update shape with opaque Node ID, organization-facing display/lifecycle and only the bounded Fulfillment-owned eligibility/capability/configuration needed by the claimed path. InventorySource/native warehouse/company/location identity must remain distinct; no generic capability graph/WMS model.

#### EconomicPerformanceSummary

Keyed/period-scoped Q with no synthetic summary ID. Preserve explicit period/scope, coverage/partiality, exact monetary aggregates/rates/counts required by the real consumer, and evaluation/provenance. It never becomes a finance ledger or finality claim stronger than SaleEconomics coverage.

### Reopen trigger

If any missing Q requires a new owner/identity rather than a bounded schema under accepted owners, return only to that implicated D1/D2 decision. No such need is currently proven.

## 7. F-W2-G3 — REVISE: media route contradicts W2-E's own `If-Match` rule

### Finding

W2-E says `If-Match` applies only when the selected HTTP request resource is the stale-state authority; a different referenced resource uses a referenced-resource precondition.

But W2-E media intake uses:

```text
POST /listing-intents/{listing_intent_id}/media
If-Match: <ListingIntent ETag>
```

The selected request target is the contained media collection, while the ETag belongs to the parent ListingIntent. That applies HTTP conditional semantics to a different resource than the request target — the exact misuse W2-E rejects elsewhere.

### Credible alternatives

A. keep collection POST and invent a parent/custom precondition header/body field — works but adds a second precondition spelling for the parent solely because of path choice;

B. **selected:** make media intake an owner-specific capability on the exact mutable ListingIntent resource:

```text
POST /organizations/{org}/listing-intents/{id}:create-media
If-Match: <ListingIntent ETag>
Idempotency-Key: ...
Content-Type: multipart/form-data
```

This lets standard `If-Match` protect the selected resource and still returns a ListingIntent-scoped media descriptor; no standalone media collection authority is implied.

### Hardening

For binary/multipart intake, semantic idempotency equivalence must include the binary content identity plus material semantic metadata. D7 selects digest/hash/storage mechanics; same key with different file bytes is a materially different request.

No parent reopen.

## 8. F-W2-G4 — REVISE: PublicationValue underfits current Mercado Livre requirement types

### External evidence

Current official Mercado Livre attribute documentation lists `number_unit` alongside string/number/boolean/list and exposes allowed/default units. It also documents explicit provider N/A attribute semantics.

### Finding

W2-B PublicationValue currently supports text, exact_decimal, boolean, option, text_list and option_list. It cannot represent a typed exact number + allowed unit without falling back to text or provider-specific formatting, which would lose semantic validation and push truth selection toward adapter/client formatting.

### Corrected direction

Extend only the bounded PublicationValue/RequirementValueSpec family:

- add `number_unit` / quantity-like publication value with `ExactDecimalString` + **requirement-scoped opaque `unit_key`**;
- Readiness requirement value spec exposes the allowed unit keys/default unit when materially required;
- adapter maps unit key to provider-native unit representation;
- no generic unit conversion/UoM engine is introduced;
- add a bounded `not_applicable` PublicationValue variant only when the current Readiness requirement explicitly permits N/A; adapter maps it to provider-native sentinel semantics.

The two D4-R1 resolution modes remain unchanged: FOLLOW_SOURCE or EXPLICIT_OVERRIDE. `not_applicable` is an explicit override value meaning, not a third resolution authority.

### Reopen trigger

A future real requirement value shape not representable without information loss triggers the smallest PublicationValue extension; never an arbitrary JSON/property bag by default.

## 9. F-W2-G5 — REVISE: provider conditional requirements can be draft-dependent

### External evidence

Current Mercado Livre documentation exposes conditional-required attribute evaluation in Brazil and instructs callers to send the concrete listing information; the example includes title, category, price, currency, available quantity, listing type, description and other item information.

### Finding

W2-C currently implies Readiness can always return an effective requirement/applicability from source Product + Marketplace Installation + publication context alone while hiding provider condition rules.

That is too strong for a provider whose conditional applicability may depend on the concrete draft and other owner-issued inputs such as price/quantity.

Exposing raw provider conditions as a rules DSL would be wrong, but pretending applicability is context-only can make a locally valid draft fail at provider validation.

### Corrected direction

Preserve current ownership and add a bounded **draft-dependent requirement evaluation seam**:

- Readiness owns requirement definition/key/value-spec/source candidates and may mark an applicability class as unconditional/currently-known versus `draft_dependent` where provider evidence requires concrete listing input;
- Product API does **not** expose provider condition/expression DSL;
- Offering remains owner of ListingIntent and draft dispatchability;
- for draft-dependent validation, server-side D4 technical machinery may compose current owner-issued ListingIntent + PriceIntent + Availability meanings solely to evaluate the provider's conditional-requirement contract, analogous to accepted joint technical realization;
- the translated validation result feeds Offering dispatchability/blockers while requirement definition/source candidate authority remains Readiness-owned;
- client never submits raw provider validation payload, price/quantity inside ListingIntent, or a provider condition result as truth.

### Parent-edge check

No D1/D3 reopen is required **if** this remains a pure D4 validation/effect-support mechanism using already owner-issued inputs and does not cause Readiness to own/read mutable ListingIntent state or Offering to own provider requirement definition.

Targeted D1/D3/D4-R1 reopen is required only if real implementation proves provider conditional applicability cannot be represented without a new semantic dependency/owner meaning.

## 10. F-W2-G6 — ACCEPT WITH SHARPENING: FulfillmentExecution identity must not become speculative/duplicate identity

### Finding

W2-D's need for an independently addressable Fulfillment-owned lifecycle is credible, but its justification partly leans on hypothetical future split fulfillment (`0..N` executions). That future is not a Product 1.0 requirement and must not be the reason a canonical identity exists.

D2 already prepared one durable Fulfillment routing/dispatch-intent seam where materially required. W2 must not later create both a `FulfillmentIntentId` and a separate `FulfillmentExecutionId` for the same lifecycle.

### Corrected invariant

- `FulfillmentExecutionId` is accepted only as the **one concrete D2 durable Fulfillment lifecycle/intent identity** backing the admitted physical checkpoint/history surface;
- no parallel FulfillmentIntent/Workflow identity may coexist for the same meaning;
- current justification is durable owner history/correlation/checkpoint/artifact/Materialization reference, not speculative future split fulfillment;
- explicit Sale-relative scope remains the future seam;
- a later real need may permit multiple executions by normal use of the existing identity; baseline must not implement split-routing policy now.

Reopen/drop trigger: if later design shows Product 1.0 Fulfillment state is fully and safely keyed by Sale+scope with no independent durable lifecycle/history identity, remove the synthetic ID before OpenAPI rather than preserving it for REST aesthetics.

No parent reopen presently required.

## 11. F-W2-G7 — ACCEPT WITH HARDENING: local opaque selector keys become historical references once persisted

Context-bound keys such as `sale_line_key` are correctly not global entities. However, once one appears in durable Fulfillment/Post-Sale/Work history it must never later be rebound/recycled to denote a different Sale line merely because provider ordering/current interpretation changed.

Similarly, transient requirement/source/media candidate and option keys may remain context-bound, but consequential publication history must preserve the resolved material value/provenance under F-W2-G1 rather than depend on those keys remaining resolvable forever.

This creates no universal identity graph.

## 12. Global checks after findings

### Passes

- owner authority boundaries remain coherent;
- no generic Product/PIM/Integration/Workflow/Finance/Task/Operation authority is required;
- Money/exact decimal and knowledge/freshness/coverage semantics remain coherent;
- PriceIntent and Availability remain separate from ListingIntent through initial publication;
- client/server authority fence remains sound;
- direct vs referenced-resource precondition concept remains sound after media correction;
- idempotency-before-stale-precondition retry rule remains sound;
- Problem Details vs semantic owner outcomes remains coherent;
- physical-fact client-class fence remains sound at D5 altitude; D7 must realize the trusted system-Principal/source binding without turning it into a second business Permission authority;
- Work closure-path audit remains PASS under current evidence;
- provider richness remains bounded rather than lowest-common-denominator or DTO mirror.

### Parent-stage check

No current finding proves D0/D1/D2/D3/D4/D4-R1/D5-B1 must reopen. F-W2-G5 carries an explicit targeted reopen trigger if the real conditional-requirement mechanism cannot preserve current ownership.

## 13. Structural inversion / future-cost result

The corrections remain correct if legacy routes/OpenAPI/controllers are opposite in every relevant respect.

They do not add versioning, generic bulk, public event streams, PIM/media master, rules engine, workflow engine, finance ledger, provider graph, generic UoM system or D6 BFF.

The added seams are tied to current evidence/failure classes:

- historical action explainability;
- admitted operation reachability;
- correct HTTP conditional semantics;
- current provider attribute/value contract;
- current provider conditional requirement contract;
- one non-duplicated Fulfillment lifecycle identity.

## 14. Lead disposition

```text
D0→D5-B1 / ratified B2 operation authority    CURRENT STRUCTURE CONFIRMED
W1 / W2-A→E local structure                    RESTRUCTURE NOW — bounded W2 corrections above
Parent-stage reopen                            NONE CURRENTLY PROVEN
Whole-W2 ready for independent review          NO — operator ratification of lead direction first
```

## 15. Proposed next action

Operator reviews/ratifies F-W2-G1..G7 as the lead correction direction.

If ratified:

1. keep accepted W2-A→E as authority until independent review adjudication; do not silently canonicalize review findings first;
2. use this file as the bounded **NON-AUTHORITATIVE W2 review candidate**;
3. open `AI-DIALOG.md` for one coherent independent Fable Whole-W2 review;
4. require Fable to reconstruct authority and attack every finding plus search for missing contradictions;
5. GPT adjudicates every material reviewer finding;
6. if no contradiction survives, operator ratifies the converged package;
7. only then consolidate corrections into canonical W2 artifacts, delete this candidate/reset AI-DIALOG, and advance the router to the next Wire Contract batch.

Implementation remains blocked until D9.
