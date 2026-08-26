# D6-R2 P9 — B10 Preparação Screen Contract

> **Status:** DERIVED / PASS — bounded correspondence rerun completed 2026-08-26 after operator re-LOCK; BACKEND SUFFICIENT; UPSTREAM FINDING NONE
> **Block:** B10 — Preparação / R10
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Locked P8 evidence:** `qualification/d6-r2-wireframes/b10-preparation.html`
> **Canonical Product OAD:** `contracts/api/product/openapi.yaml`
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P9 result

P9 was rerun only after the operator re-LOCKED the simplified B10 candidate.

The locked human job is:

```text
exact Organization + Marketplace Installation
→ search exact source-qualified product
→ inspect marketplace requirements + source values/evidence
→ resolve/clear correspondence when needed
→ authoritative reread after correspondence effect
→ continue to downstream ListingIntent authoring
```

The current Product contract is sufficient for this job. No screen-shaped endpoint, per-requirement satisfaction field, generic mapping engine, provider-field bag or new frontend business authority is required.

**P9 verdict: PASS / BACKEND SUFFICIENT / UPSTREAM FINDING NONE.**

## 2. Route, identity and client-state ownership

Production route family:

`/org/:organizationId/preparacao`

Material carriers:

```text
organization_id
marketplace_installation_id
q
source_instance_id                 optional search narrowing only
selected_source_instance_id
selected_native_product_key
```

Client-state classes:

| State class | B10 ownership |
| --- | --- |
| `GLOBAL_WORKSPACE_CONTEXT` | `organization_id`; changing Organization invalidates incompatible Installation, selected-subject and server state |
| `URL_NAVIGATION_STATE` | `marketplace_installation_id`, `q`, optional `source_instance_id`, selected source/native product identity when deep-linkable |
| `SERVER_STATE` | Installation collection, source search results, ProductChannelReadiness/correspondence, PublicationRequirements and their revisions/evidence |
| `LOCAL_EPHEMERAL` | mobile navigation disclosure, technical-details disclosure and other non-business presentation state |

Prototype-only scenario selectors are Evidence controls, not production Product state.

URL/router state never becomes business authority. TanStack Query remains the accepted production owner of server state; no normalized frontend business-entity mirror is introduced.

## 3. Product operation / access binding

| Screen need | Operation | Semantic owner | Permission | Principal kinds | P9 disposition |
| --- | --- | --- | --- | --- | --- |
| populate exact account selector | `ListMarketplaceInstallations` | MarketplacePortfolio | `portfolio.read` | H/A/S | admitted read; human surface consumes H |
| search source products | `SearchSourceProductsForMarketplace` | ProductChannelReadiness | `readiness.read` | H/A/S | exact source-qualified results; optional SourceInstance is narrowing only |
| read exact subject/correspondence | `GetProductChannelReadiness` | ProductChannelReadiness | `readiness.read` | H/A/S | supplies correspondence/current-read basis; no mandatory per-field readiness label |
| read marketplace fields/source evidence | `GetPublicationRequirements` | ProductChannelReadiness | `readiness.read` | H/A/S | complete provider/context requirement census + source evidence |
| define correspondence | `ResolveProductChannelCorrespondence` | ProductChannelReadiness | `readiness.manage` | H/A | explicit consequential action |
| remove correspondence | `ClearProductChannelCorrespondence` | ProductChannelReadiness | `readiness.manage` | H/A | explicit consequential action |
| downstream draft creation | `CreateListingIntentDraft` | Offering | `listing.manage` | H/A | **no B10 mutation home**; future ListingIntent authoring boundary only |

Human browser realization remains server-side session + CSRF. Permission-conditioned control visibility is usability only; server authorization is authoritative.

`CreateListingIntentDraft` has `Idempotency-Key` semantics in the OAD, but B10 does not call it merely to leave Preparação.

## 4. Material screen contract

### 4.1 Organization + marketplace account context

**GOAL / FLOW:** establish the exact operating context before source search.  
**ROUTE / SURFACE:** global Organization shell + B10 account selector.  
**INFORMATION ROLE:** scope/navigation context.  
**OWNER + READ TRUTH:** Organization access context + MarketplacePortfolio `ListMarketplaceInstallations`.  
**WRITE CONTROL:** none.  
**IDENTITY SOURCE:** canonical `organization_id` + exact `marketplace_installation_id`.  
**CLIENT STATE CLASS:** global workspace + URL navigation + server collection state.  
**WIRE MECHANICS:** account list may populate the selector; no first/default account is selected.  
**MATERIAL FAILURES:** no accessible Organization; exact Installation missing/invalid/inaccessible; Installation list unavailable.  
**FAILURE MESSAGE INTENT:** name missing/unavailable context without inventing fallback.  
**SUCCESS CONSEQUENCE:** exact Installation context established; changing Organization/Installation invalidates selected-subject state.  
**AUTHZ / DISCLOSURE:** `portfolio.read` for population; server auth remains authoritative.  
**FORBIDDEN FRONTEND AUTHORITY:** ambient/default tenant/account or account-derived authorization.  
**BACKEND SUFFICIENCY:** sufficient.

### 4.2 Source-product search and selection

**GOAL / FLOW:** find the exact source product intended for this marketplace context.  
**ROUTE / SURFACE:** B10 search/results.  
**INFORMATION ROLE:** source-qualified discovery evidence.  
**OWNER + READ TRUTH:** ProductChannelReadiness `SearchSourceProductsForMarketplace`.  
**WRITE CONTROL:** none.  
**IDENTITY SOURCE:** returned SourceInstance + native product key; optional `source_instance_id` only narrows search.  
**CLIENT STATE CLASS:** URL navigation (`q`, optional source filter) + server search state.  
**WIRE MECHANICS:** bounded search; omission of source filter searches admitted Organization-scoped sources and never picks a hidden default.  
**MATERIAL FAILURES:** blank query validation; known-empty; unavailable search; inaccessible/not-found source identity.  
**FAILURE MESSAGE INTENT:** preserve `known-empty != unavailable`; never present transport/source failure as no product.  
**SUCCESS CONSEQUENCE:** one exact source-qualified subject is selected for detail.  
**AUTHZ / DISCLOSURE:** `readiness.read`; technical source keys may be secondary support detail.  
**FORBIDDEN FRONTEND AUTHORITY:** MPC Product master/PIM, source identity merge, hidden source default.  
**BACKEND SUFFICIENCY:** sufficient.

### 4.3 Marketplace fields + source values

**GOAL / FLOW:** show what the marketplace asks for and what source value/evidence exists before authoring.  
**ROUTE / SURFACE:** locked four-column `Campos para o marketplace` region.  
**INFORMATION ROLE:** provider/context requirement + source evidence comparison, not final draft validation.  
**OWNER + READ TRUTH:** ProductChannelReadiness `GetPublicationRequirements`.  
**WRITE CONTROL:** none in B10.  
**IDENTITY SOURCE:** exact Organization + Installation + SourceInstance + native product key; optional provider `category_key` / `product_type_key` context; `requirements_revision` remains server truth.  
**CLIENT STATE CLASS:** server state.  
**WIRE MECHANICS:** preserve independent requirement class/applicability/value specification/source evidence; human projection is `Campo do marketplace / Exigência / Valor encontrado / Na configuração do anúncio`.  
**MATERIAL FAILURES:** source evidence `missing`, `conflicting`, `unknown`, `unavailable`, `unsupported`; provider/context requirement read unavailable or invalid.  
**FAILURE MESSAGE INTENT:** distinguish genuine absence from inability to determine/consult and from multiple candidates; do not label any as generic `Atendido`/`Não atendido`.  
**SUCCESS CONSEQUENCE:** operator understands available values and unresolved authoring work; missing source value does not block ListingIntent entry by itself.  
**AUTHZ / DISCLOSURE:** `readiness.read`; technical evidence/candidate keys remain secondary.  
**FORBIDDEN FRONTEND AUTHORITY:** `source_sufficiency`, per-field `satisfied/met/ready`, generic provider rule engine, raw provider field bag.  
**BACKEND SUFFICIENCY:** sufficient; no new wire field required.

### 4.4 Product↔channel correspondence

**GOAL / FLOW:** establish or remove exact product correspondence when needed, then return to authoritative current truth.  
**ROUTE / SURFACE:** `Vínculo com o marketplace`.  
**INFORMATION ROLE:** current correspondence + explicit correction control.  
**OWNER + READ TRUTH:** ProductChannelReadiness `GetProductChannelReadiness`.  
**WRITE CONTROL:** `ResolveProductChannelCorrespondence` / `ClearProductChannelCorrespondence`.  
**IDENTITY SOURCE:** exact Organization + Installation + SourceInstance + native product key + current `correspondence_etag` carried by the accepted write schema.  
**CLIENT STATE CLASS:** server state plus local pending-control state only while submitting.  
**WIRE MECHANICS:** explicit action; current correspondence validator; successful or potentially accepted effect is followed by authoritative reread; requirements reread when current context/revision can change.  
**MATERIAL FAILURES:** 401/403/404/409/422/500; stale/conflicting correspondence; transport outcome where acceptance cannot be ruled out.  
**FAILURE MESSAGE INTENT:** never claim known failure after ambiguous potential acceptance; direct operator to refresh/reconcile current state.  
**SUCCESS CONSEQUENCE:** refreshed current correspondence/requirement context; continuation is re-enabled only from current reread truth.  
**AUTHZ / DISCLOSURE:** `readiness.manage`, H/A; browser POSTs remain CSRF-protected; hidden control is not authorization.  
**FORBIDDEN FRONTEND AUTHORITY:** blind retry, local correspondence truth, write-success assumption without reread.  
**BACKEND SUFFICIENCY:** sufficient.

### 4.5 Continue to ListingIntent authoring

**GOAL / FLOW:** leave preparation and continue the human job in Offering authoring.  
**ROUTE / SURFACE:** `Continuar para configurar o anúncio` → explicit unopened ListingIntent/B23 boundary.  
**INFORMATION ROLE:** navigation handoff only.  
**OWNER + READ TRUTH:** Offering owns ListingIntent desired state; B10 owns none of it.  
**WRITE CONTROL:** none in B10.  
**IDENTITY SOURCE:** exact prepared subject/context may be carried as navigation context; ListingIntent identity is created/owned downstream.  
**CLIENT STATE CLASS:** URL/navigation only.  
**WIRE MECHANICS:** no B10 call to `CreateListingIntentDraft`; the unopened B23 block will decide exact downstream authoring route/mechanics.  
**MATERIAL FAILURES:** correspondence authority not established; downstream access not permitted; stale context before authoring.  
**FAILURE MESSAGE INTENT:** explain why continuation is unsafe/unavailable without pretending an advertisement was created or published.  
**SUCCESS CONSEQUENCE:** B10 ends; no marketplace effect has occurred.  
**AUTHZ / DISCLOSURE:** downstream authoring requires `listing.manage`; server authorization at the downstream surface remains authoritative.  
**FORBIDDEN FRONTEND AUTHORITY:** creating desired ListingIntent state as a navigation side effect, direct provider publication, treating navigation as provider validation.  
**BACKEND SUFFICIENCY:** sufficient for B10; exact B23 UX remains future block design, not a B10 backend gap.

## 5. Bidirectional trace — PASS

### Frontend → Product/backend

```text
Organization + Conta do marketplace
→ access/MarketplacePortfolio context
→ ListMarketplaceInstallations when population is required

Buscar produto + cadastro de origem
→ SearchSourceProductsForMarketplace

Produto selecionado + vínculo
→ GetProductChannelReadiness

Campo do marketplace / Exigência / Valor encontrado
→ GetPublicationRequirements

Definir vínculo / Remover vínculo
→ ResolveProductChannelCorrespondence / ClearProductChannelCorrespondence
→ authoritative reread

Continuar para configurar o anúncio
→ downstream Offering boundary
→ NOT a B10 CreateListingIntentDraft call
```

### Product/backend → Frontend

```text
ListMarketplaceInstallations
→ exact account selector choices

SearchSourceProductsForMarketplace
→ source-qualified result list

GetPublicationRequirements
→ full marketplace requirement census + source values/evidence

GetProductChannelReadiness
→ exact subject/correspondence + current reread basis
→ overall readiness does not force a per-field/status UI

Resolve/Clear correspondence
→ explicit write controls + reread/reconciliation behavior

CreateListingIntentDraft
→ no B10 mutation home
→ downstream Offering authoring only
```

No admitted B10 operation is orphaned from the locked screen job, and no material locked screen control requires a missing Product operation.

## 6. Adversarial checks

P9 explicitly rejects these regressions:

- `known` source evidence becoming `Atendido`;
- missing source value becoming publication-impossible;
- frontend-computed `source_sufficiency`/`satisfied`;
- hidden/default Marketplace Installation or SourceInstance;
- unknown/unavailable collapsing into known-empty;
- correspondence write treated as locally final without reread;
- blind retry after ambiguous possible acceptance;
- B10 creating a ListingIntent merely to navigate;
- direct marketplace publication/validation authority inside B10;
- reopening Product/OAD merely to make the screen easier to code.

All are avoidable with the accepted wire and the operator-LOCKED projection.

## 7. P9 closure

```text
P8 OPERATOR-RATIFIED / LOCKED
→ exact route/state/identity binding
→ exact owner/operation/Permission binding
→ frontend → backend trace PASS
→ backend → frontend trace PASS
→ adversarial shortcuts rejected
→ BACKEND SUFFICIENT
→ UPSTREAM FINDING NONE
```

PR #70 repaired the correspondence candidate read projection and triggered the smallest declared P8/P9 reopen. The operator re-LOCKED the bounded correspondence candidate on 2026-08-26 (`aprovado` → `LOCK`, recorded in [`D6-R2-P8-B10-CORRESPONDENCE-REVALIDATION.md`](D6-R2-P8-B10-CORRESPONDENCE-REVALIDATION.md)), and the bounded correspondence trace was rerun against the integrated canonical OAD:

- `GetProductChannelReadiness` carries `subject_presentation`, `correspondence` (resolved/unresolved/conflicting/unknown/unavailable), `correspondence_candidate_population` (known with `candidates[]`, unknown, unavailable) and `correspondence_etag`;
- each known candidate is exactly `candidate_key + display_label`; no label is ever a write carrier;
- `ResolveProductChannelCorrespondence` accepts only `subject + correspondence_etag + candidate_key`; `ClearProductChannelCorrespondence` only `subject + correspondence_etag`;
- the locked region 4.4 mechanics (explicit selection, disabled resolve until choice, mandatory authoritative reread after consequential effect) bind one-to-one to these schemas with no missing operation and no orphaned admitted operation.

**P9: PASS / CLOSED for B10 against the integrated PR #70 wire.**

P10 may now consolidate only patterns already repeated in LOCKED evidence. P11, Pre-D9/D9 and Product implementation remain outside this P9 closure; B20 resumption follows integration of this increment per `docs/roadmap.md`.
