# D6-R2 P9 — B10 Preparação Screen Contract

> **Status:** DERIVED / BLOCKED — P8 REOPEN REQUIRED
> **Block:** B10 — Preparação / R10
> **P8 input:** [B10 P8 ratification](D6-R2-P8-B10-PREPARATION-RATIFICATION.md) — OPERATOR-RATIFIED / LOCKED
> **Method profile:** `developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535` → `METHOD + FRONTEND-METHOD`
> **Canonical Product OAD:** `contracts/api/product/openapi.yaml`
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P9 purpose and verdict

P9 binds the LOCKED human-first B10 screen to accepted Product/backend authority in both directions. It does not redesign the screen, add endpoints, create a PIM, reinterpret provider evidence, or begin implementation.

The binding is structurally viable, but P9 falsified one visible P8 statement:

**FRONTEND LOCK FINDING: `F-P9-B10-01` — `known` source evidence does not equal requirement satisfied.**

D4-R1 makes ProductChannelReadiness the owner of source-level readiness/sufficiency and says source candidates **may** satisfy a requirement. PR #68 exposes source evidence and value constraints, but it does not issue a per-requirement `satisfied/met` fact. The locked P8 currently renders `Atendido` and `requisitos atendidos` for `source_evidence.state = known`; that creates a stronger business conclusion in the frontend than the owner issued.

**UPSTREAM FINDING: NONE.** Product authority is sufficient to show that source information exists. The smallest lawful repair is a bounded P8 language correction:

```text
Atendido                 → Informação disponível
2 requisitos atendidos   → 2 com informação disponível
```

No layout, flow, Product contract, requirement census, operation, Permission or owner changes.

Until that wording is operator-adjudicated and the P8 LOCK is re-ratified, **P9 remains BLOCKED and P10 must not start**.

## 2. Canonical route and identity carrier

B10 remains one workspace:

`/preparacao`

P9 does not introduce a second detail route. Search mode and selected-subject mode are durable states inside R10.

### 2.1 Identity/state carriers

| State | Class | Source / law |
| --- | --- | --- |
| current `organization_id` | `GLOBAL_WORKSPACE_CONTEXT` | inherited from the LOCKED B00 Organization workspace; B10 never invents or defaults it |
| `marketplace_installation_id` | `URL_NAVIGATION_STATE` | exact page-local marketplace account context; required before search/detail |
| `q` | `URL_NAVIGATION_STATE` | committed search term; local typing may precede URL commit |
| `source_instance_id` | `URL_NAVIGATION_STATE` | optional search narrowing only; omission means admitted multi-source search, never a default source |
| `selected_source_instance_id` | `URL_NAVIGATION_STATE` | exact selected subject identity in detail mode |
| `selected_native_product_key` | `URL_NAVIGATION_STATE` | exact selected subject identity in detail mode |
| source search collection + `next_cursor` | `SERVER_STATE` | `SearchSourceProductsForMarketplace` |
| readiness + correspondence + `correspondence_etag` + `requirements_revision` | `SERVER_STATE` | `GetProductChannelReadiness` |
| requirement census + source candidates + source media | `SERVER_STATE` | `GetPublicationRequirements` |
| search input before submit, mutation-pending UI, technical-details expanded/collapsed | `LOCAL_EPHEMERAL` | frontend interaction state only; never Product truth |

Prototype-only scenario switches (`Normal`, `Nenhum resultado`, `Busca indisponível`, etc.) are qualification harness controls and **do not become production state**.

### 2.2 URL invariants

A lawful detail URL may conceptually carry:

```text
/preparacao
?marketplace_installation_id=<opaque>
&q=<search-term>
&source_instance_id=<optional-search-filter>
&selected_source_instance_id=<opaque>
&selected_native_product_key=<opaque>
```

Rules:

- `selected_source_instance_id` and `selected_native_product_key` are a pair; partial selected identity is invalid navigation state and returns to a recoverable search state rather than guessing.
- changing Organization invalidates the complete page-local Installation, search and selected-subject context;
- changing Installation invalidates selected-subject server state and rereads under the new exact account;
- removing selected identity returns to search mode without destroying committed `q` / optional source filter;
- URL values are identity/navigation carriers only; server authorization remains authoritative.

## 3. Screen contract by material region/control

| Region / control | Goal / information role | Owner + read truth | Write/control | Identity source | State class | Material failures / consequence | Forbidden frontend authority |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Organization shell | establish global workspace | inherited B00 access authority | Organization switch is navigation/context, not B10 Product write | current Organization | `GLOBAL_WORKSPACE_CONTEXT` | switch clears Installation + selected subject | no ambient cross-org reuse |
| Marketplace account selector | establish exact provider/account context | `ListMarketplaceInstallations` / MarketplacePortfolio where selector population is needed | selection only | Organization + chosen Installation | `URL_NAVIGATION_STATE` | inaccessible/stale account → explicit missing/invalid context | no first/default account |
| search form | find source products admitted for marketplace operation | `SearchSourceProductsForMarketplace` / ProductChannelReadiness | GET query only | Organization + Installation + `q` + optional source filter | URL + server | 400/422 search problem; 403 access denied; 500 unavailable | no unfiltered Product-master browser |
| results list | choose exact source-qualified subject | same search collection | selecting row changes navigation state only | returned SourceProductRef | URL + server | known empty distinct from unavailable | no silent first-row selection |
| selected-product header | preserve exact subject | readiness subject + selected SourceProductRef | none | selected source instance + native key + Installation | URL + server | partial/stale selected identity fails closed | no display name as identity |
| preparation summary | orient operator without inventing business truth | ProductChannelReadiness + PublicationRequirements | none | exact selected subject | derived presentation over `SERVER_STATE` | owner unknown/unavailable stays explicit | **must not synthesize per-requirement satisfied/met** |
| requirements table | show complete applicable marketplace requirement census | `GetPublicationRequirements` / ProductChannelReadiness | none in B10 | subject + provider/category/product-type context + revision | `SERVER_STATE` | unknown/missing/conflicting/unavailable/unsupported remain distinct | no provider field bag; no frontend readiness computation |
| correspondence region | inspect/change Product↔channel correspondence | `GetProductChannelReadiness` | Resolve/Clear only with `readiness.manage` | exact subject + `correspondence_etag` | server + local pending | 409 stale/conflict → reread; ambiguous transport → reread before any retry | no generic Product edit; no blind retry |
| “Continuar para configurar o anúncio” | leave B10 for downstream ListingIntent authoring flow | no B10 Product mutation | navigation boundary only | source subject + Installation context | `URL_NAVIGATION_STATE` | disabled when typed B10 preconditions cannot be established | B10 does not call `CreateListingIntentDraft` |
| “Ver detalhes técnicos” | support/audit disclosure | already-loaded server/trace evidence | expand/collapse only | current selected subject | `LOCAL_EPHEMERAL` | disclosure failure cannot mutate Product truth | technical vocabulary cannot become primary workflow |

## 4. Product/backend operation contracts

### 4.1 Installation context population

`ListMarketplaceInstallations`

- owner: MarketplacePortfolio;
- read Permission: `portfolio.read`;
- admitted Principals: H / A / S;
- Organization-scoped collection;
- B10 uses it only to populate the exact page-local marketplace-account selector when that population is not already supplied by the locked shell context.

No Installation becomes ambient authority merely because it appears first.

### 4.2 Source search

`SearchSourceProductsForMarketplace`

- class: Q;
- owner: `ProductChannelReadiness`;
- Permission: `readiness.read`;
- Principals: H / A / S;
- request identity: `organization_id` + required `marketplace_installation_id` + required search `query`;
- optional narrowing: `source_instance_id`;
- scale: `limit` + `cursor`;
- 200: `SourceProductsSearchCollection`, where every hit remains source-qualified;
- material failures: 400 / 401 / 403 / 404 / 422 / 500.

Omitting `source_instance_id` is multi-source admitted search, not source selection.

### 4.3 Selected-subject readiness

`GetProductChannelReadiness`

- class: Q;
- owner: `ProductChannelReadiness`;
- Permission: `readiness.read`;
- Principals: H / A / S;
- exact request identity: Organization + `marketplace_installation_id` + `source_instance_id` + `native_product_key`;
- 200 server truth includes subject, correspondence, `correspondence_etag`, readiness, blockers, optional `requirements_revision`, evaluated time;
- material failures: 401 / 403 / 404 / 422 / 500.

The frontend may present owner-issued readiness/correspondence truth but must not parse free-form blocker strings into new business rules. If a future pre-ListingIntent blocker needs typed behavior not expressible by accepted fields, P9 must raise an upstream Finding rather than teach the client a private blocker vocabulary.

### 4.4 Publication requirements

`GetPublicationRequirements`

- class: Q;
- owner: `ProductChannelReadiness`;
- Permission: `readiness.read`;
- Principals: H / A / S;
- exact subject: Organization + Installation + SourceInstance + native product key;
- optional context qualifiers: category key / product type key only when established by current authority; no frontend default;
- 200: complete current `PublicationRequirements` projection with `requirements_revision`, requirement census, source evidence/candidates and source-media candidates;
- material failures: 401 / 403 / 404 / 422 / 500.

P9 preserves PR #68 laws:

```text
requirement_class != applicability
known / missing / conflicting / unknown / unavailable / unsupported
7 bounded value_spec families
not_applicable_allowed is override capability, not source fact
candidate evidence remains opaque/context-bound
```

**Important:** `known` means source evidence is known. It is not a Product-issued `satisfied/met` field. This is the basis of `F-P9-B10-01`.

### 4.5 Resolve correspondence

`ResolveProductChannelCorrespondence`

- class: C;
- owner: `ProductChannelReadiness`;
- Permission: `readiness.manage`;
- Principals: H / A;
- request body: exact subject + `correspondence_etag` + selected `candidate_key`;
- 200: current `ProductChannelReadiness`;
- material failures: 401 / 403 / 404 / 409 / 422 / 500.

The write has no client-issued idempotency identity in the accepted wire. Therefore an ambiguous transport outcome gets **no blind retry**. The recovery path is authoritative reread of readiness/correspondence, then a new operator decision only if the current state still requires it.

### 4.6 Clear correspondence

`ClearProductChannelCorrespondence`

- class: C;
- owner: `ProductChannelReadiness`;
- Permission: `readiness.manage`;
- Principals: H / A;
- request body: exact subject + `correspondence_etag`;
- 200: current `ProductChannelReadiness`;
- material failures: 401 / 403 / 404 / 409 / 422 / 500;
- same **no blind retry** / authoritative-reread recovery law as Resolve.

### 4.7 Downstream ListingIntent boundary

`CreateListingIntentDraft` belongs to Offering, not B10:

- owner: `Offering`;
- class: C;
- Permission: `listing.manage`;
- Principals: H / A;
- requires `Idempotency-Key`;
- request requires `source_product`, `target`, `desired`; may carry `requirements_revision`;
- 201 returns the created ListingIntent plus Location/ETag semantics;
- material failures include 400 / 401 / 403 / 404 / 409 / 422 / 500.

**B10 does not call `CreateListingIntentDraft`.** The locked B10 button is a navigation handoff into the downstream authoring flow because B10 does not own `desired` listing state. The downstream block binds actual draft creation when it has the Offering-owned inputs. `data-next-operation="CreateListingIntentDraft"` remains trace evidence, not permission for B10 to manufacture an empty/default desired draft.

## 5. Correspondence mutation lifecycle

A lawful Resolve/Clear cycle is:

```text
current exact subject
+ current correspondence_etag
+ explicit operator action
→ POST Resolve/Clear
→ do not infer convergence from click
→ authoritative GetProductChannelReadiness reread
→ GetPublicationRequirements reread when current correspondence/context/revision can affect the requirement basis
→ render current owner truth
```

On 409, tell the operator that information changed and refresh. On 422, explain that the selected action is no longer valid for the current state. On 403, remove/disable the mutation affordance after the authoritative denial while preserving any read access independently. On ambiguous network failure, say in human terms that the result could not be confirmed and require refresh before another attempt.

Permission-conditioned controls reduce noise only; server authorization remains authoritative.

## 6. B10 progression law

B10 progression is a **frontend workflow transition**, not a new Product readiness status.

The continue control fails closed when:

- exact Organization/Installation/selected subject cannot be established;
- correspondence is in a state that requires explicit resolution before progression;
- required Readiness/PublicationRequirements authority is unknown/unavailable or the reads failed such that the current basis cannot be established;
- server access denies the required read.

B10 does **not** block merely because one requirement has source evidence `missing`, `conflicting` or `unsupported`; D4-R1 allows Offering-owned ListingIntent authoring to resolve desired values later without rewriting source truth.

Likewise, B10 must not turn the free-form `blockers` array into a client rules engine. Any future Product condition that must block B10 but lacks typed authority is an upstream Finding.

## 7. Failure-message intent

| Condition | Operator message intent | Recovery |
| --- | --- | --- |
| no Installation | “Selecione uma conta do marketplace.” | choose exact account |
| search known-empty | “Nenhum produto encontrado.” | change query/filter |
| search/read unavailable | “Não conseguimos consultar as informações agora.” | retry read later; do not show empty |
| 403 read | “Você não tem acesso a esta preparação.” | return to accessible surface |
| selected subject 404 | “Este produto não está disponível neste contexto.” | return to search |
| correspondence 409 | “As informações mudaram. Atualize antes de continuar.” | authoritative reread |
| correspondence 422 | “Essa opção não é válida para o estado atual.” | reread + choose current option |
| ambiguous correspondence transport | “Não foi possível confirmar se a alteração foi aplicada.” | reread; **no blind retry** |
| source missing | “Falta informação no cadastro de origem.” | continue to downstream authoring when other B10 gates permit |
| source conflicting | “Há informações diferentes.” | show evidence; do not auto-choose |
| source unknown/unavailable | “Não foi possível verificar.” | preserve uncertainty; fail closed where basis is required |

Errors never reveal inaccessible resource existence beyond the Product problem response.

## 8. Bidirectional trace

### frontend → backend

| Frontend surface/control | Contract / authority |
| --- | --- |
| Organization context | inherited B00 global access context; becomes Product `organization_id` input |
| marketplace account selector | `ListMarketplaceInstallations` / exact selected Installation |
| search submit | `SearchSourceProductsForMarketplace` |
| optional source filter | optional `source_instance_id` on search; never default source |
| open result | navigation only; stores exact selected source identity |
| preparation summary | `GetProductChannelReadiness` + `GetPublicationRequirements`; presentation only |
| requirements table | `GetPublicationRequirements` |
| Define vínculo | `ResolveProductChannelCorrespondence` with `correspondence_etag` |
| Remover vínculo | `ClearProductChannelCorrespondence` with `correspondence_etag` |
| Atualizar informações | authoritative readiness + requirement reread |
| Continuar para configurar o anúncio | navigation boundary only; **not** `CreateListingIntentDraft` inside B10 |
| Ver detalhes técnicos | local disclosure over already-authorized evidence |

### backend → frontend

| Product/backend contract | B10 home |
| --- | --- |
| `ListMarketplaceInstallations` | page-local account selector |
| `SearchSourceProductsForMarketplace` | search results / known-empty / unavailable state |
| `GetProductChannelReadiness` | exact subject, correspondence, owner readiness evidence, progression inputs |
| `GetPublicationRequirements` | complete requirement table + technical requirement evidence |
| `ResolveProductChannelCorrespondence` | explicit “Definir vínculo” mutation + reread state |
| `ClearProductChannelCorrespondence` | explicit “Remover vínculo” mutation + reread state |
| `CreateListingIntentDraft` | **no B10 mutation home**; trace points only to downstream authoring boundary |

Both directions resolve without orphan Product operations inside the B10 boundary and without a frontend control inventing an unowned mutation.

## 9. Backend sufficiency and forbidden convenience

### Backend sufficiency

For the locked B10 task, backend authority is sufficient for:

- exact multi-source search;
- selected-subject readiness/correspondence;
- complete provider/context-specific publication requirements;
- correspondence Resolve/Clear with concurrency identity;
- a lawful downstream ListingIntent creation seam.

**UPSTREAM FINDING: NONE.**

### Forbidden frontend authority

B10 must not:

- edit Sankhya/source Product master data;
- create a generic Product/PIM attribute bag;
- select a default Organization, Installation, SourceInstance, correspondence candidate or requirement value;
- turn `known` into `satisfied/met`;
- recompute source-level Readiness/sufficiency;
- parse free-form blockers into a client business rules engine;
- synthesize `not_applicable` as source evidence;
- infer unknown/unavailable as empty/zero;
- blindly retry ambiguous Resolve/Clear effects;
- call `CreateListingIntentDraft` with invented/default Offering-owned `desired` state;
- decide final ListingIntent dispatchability.

## 10. P9 disposition

All route/state/operation/write/failure bindings are otherwise coherent.

`F-P9-B10-01` is the only blocking issue found:

```text
locked visible wording: Atendido / requisitos atendidos
owner-issued truth: known source evidence / source candidates that may satisfy
lawful visible wording: Informação disponível / com informação disponível
```

**Disposition:** `DERIVED / BLOCKED — P8 REOPEN REQUIRED`.

**UPSTREAM FINDING: NONE.**

The next action is operator adjudication of the bounded wording correction. If approved, reopen only B10 P8 language, preserve the rest of the LOCK, rerun P8 proof, re-ratify the corrected candidate, then close P9 against that exact evidence. Do not start P10 before P9 is GREEN.
