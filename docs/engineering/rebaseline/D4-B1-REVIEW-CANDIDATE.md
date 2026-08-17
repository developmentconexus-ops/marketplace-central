# D4 Batch B1 — External Contract Grounding — Independent Review Candidate

> **STATUS: REVIEW CANDIDATE — NOT ARCHITECTURE AUTHORITY**  
> **Stage:** D4 — External Integrations  
> **Router state remains:** D4 `NEXT / NOT YET OPENED`; implementation remains blocked until D9  
> **Operator posture:** D4 intake/decomposition direction approved; independent Fable challenge required before any canonical D4 filing  
> **Parent authorities:** accepted D0–D3 + stable `ARCHITECTURE.md` constraints, subject only to explicit material reopen  
> **Disposable:** delete after adjudication; durable meaning belongs only in a future operator-approved canonical D4 artifact / target ADR where justified  
> **Date:** 2026-08-17

## Reviewer bootstrap

Reconstruct the current authority path independently before reviewing this candidate:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md`
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
7. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
8. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
9. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
10. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
11. only the specific provider/ERP documentation, legacy ADRs and current code needed as evidence.

Do **not** treat this file, chat history, prior GPT summaries, reviewer statements or `AI-DIALOG.md` as architecture authority.

After independent review, append findings to `AI-DIALOG.md` only as review evidence. Do not modify canonical architecture directly.

---

## 1. Batch question and boundary

D0–D3 already decide product scope, business authority, identity/tenant/data ownership and Q/C/E/P communication/failure semantics.

B1 answers only the smallest grounding question required before provider-specific operational contracts:

> **What must every concrete Mercado Livre or Sankhya acquisition/effect contract prove about source namespace, authentication, observation authority, coverage, capability and reconciliation so provider mechanics cannot silently become MPC business authority?**

B1 includes concrete external namespace binding, auth/credential semantics needed to establish that binding, notification/pointer versus authoritative reread, source coverage/completeness, the D1 three-level capability fence, minimum later-write reconciliation obligations, legacy-ADR adjudication relevant to B1, and one material current Sankhya transport conflict.

B1 explicitly excludes the full Mercado Livre operation matrix (B2), concrete Sankhya fact/command matrix (B3), market/economics/settlement (B4), D5 API, D6 frontend, D7 runtime/jobs/transactions, D8 golden flows and implementation.

---

## 2. Evidence classification

### 2.1 Already authoritative

B1 imports rather than re-decides these accepted meanings:

1. Business consumers own semantic ports; external adapters own protocol translation.
2. Provider DTOs/protocol/auth/pagination do not become domain models.
3. `Marketplace Installation` is MPC identity for marketplace participation/configuration; provider seller/account identity stays external.
4. `SourceInstance` identifies one logical externally authoritative business-system/source namespace when Installation is not the qualifier.
5. Credential/token rotation is not by itself an identity change.
6. External Listing/Variation, Sale/Order, Shipment and native financial identities remain source-qualified per D2.
7. Provider notification/callback/poll evidence is not automatically an MPC domain event.
8. D3 preserves `known value`, `known empty`, `unknown` and `unavailable`; failure cannot silently become empty/false/zero.
9. `accepted != completed != externally applied != converged`.
10. Ambiguous potentially accepted external writes are not blindly retried.
11. Partial observation absence does not prove closure/deletion.
12. Provider PII is minimized.
13. Mercado Livre proves the first marketplace loop; speculative generic integration frameworks are not Product 1.0 scope.

### 2.2 External evidence reverified for B1 — official sources only

#### Mercado Livre — rechecked 2026-08-17

Current official documentation establishes:

- private access uses OAuth 2.0 server-side authorization and bearer access tokens;
- token response carries provider `user_id`, and `/users/me` returns authenticated user identity;
- refresh tokens rotate and only the latest refresh token remains valid for another refresh;
- current user IDs can exceed Int32;
- notifications identify a source resource/topic/user/application and are followed by a resource query for full/current details;
- notifications may repeat and `missed_feeds` is a bounded recovery aid rather than infinite authoritative history;
- seller item enumeration above 1000 records requires provider scan/scroll behavior rather than ordinary offset paging;
- public item-search quantity can be referential, so endpoint choice is part of fact authority.

Official source family:

- https://developers.mercadolivre.com.br/devcenter/autenticacao-e-autorizacao
- https://developers.mercadolivre.com.br/pt_br/publicacao-de-produtos/gestao-de-identidades-e-acessos-oauth-e-tokens
- https://developers.mercadolivre.com.br/pt_br/usuarios-e-aplicativos/consulta-de-usuarios
- https://developers.mercadolivre.com.br/pt_br/produto-consulta-de-usuarios/produto-receba-notificacoes
- https://developers.mercadolivre.com.br/pt_br/itens-e-buscas

#### Sankhya — rechecked 2026-08-17

Current official documentation establishes:

- current authentication uses OAuth 2.0 `client_credentials` at `/authenticate`, plus `X-Token`;
- production and sandbox use distinct credentials/tokens;
- the OAuth flow replaces the old appkey/login model for new integration work;
- `loadRecords` exposes pagination through `offsetPage` and `hasMoreResult`;
- `modifiedSince` depends on `logAlteracoesTabelas`; without matching logged information the call can return zero records, so zero delta is not automatically proof of zero underlying changes unless logging coverage is established;
- current integration guidance says integrations should use the API Gateway as the exchange standard and characterizes direct database integrations as outside the stated integration guidelines for the applicable client population.

Official source family:

- https://developer.sankhya.com.br/reference/post_authenticate
- https://developer.sankhya.com.br/reference/guia-integracao
- https://developer.sankhya.com.br/reference/get_loadrecords
- https://developer.sankhya.com.br/reference/get_logalteracoestabelas
- https://developer.sankhya.com.br/reference/boas-pr%C3%A1ticas-para-integra%C3%A7%C3%A3o
- https://developer.sankhya.com.br/changelog/novo-fluxo-de-autenticacao-com-oauth-20

### 2.3 Current-state repository evidence only

Current code/legacy ADRs show shapes B1 must not inherit automatically:

- legacy generic connector/provider capability aggregation;
- Mercado Livre `/users/me`, seller item search, scan/scroll, point rereads and order search + point reread;
- legacy write code mapping successful HTTP 2xx directly to `applied`;
- a current `sankhyaoracle` adapter family and carried ADR-006/ADR-007 direct-Oracle/godror target language;
- mission-bounded polling-only, manual-refresh, plugin-registry and catalog-offer-flag decisions.

These prove reachable behavior and defects. They are not D4 target contracts by existence.

---

## 3. Root cause and target invariant

### Root cause

Without D4 grounding, provider mechanics can acquire business authority by convenience: credentials can become identity, a token can silently rebind an Installation, callback payloads can become truth, incomplete pages/deltas can become complete absence, transport failure can become empty data, static adapter support can become business permission, HTTP acceptance can become convergence, unsupported direct-access paths can persist because they already work, and a generic integration framework can become a duplicate semantic owner.

### Target invariant

> **Every external fact or effect entering MPC is qualified by the correct Organization + source namespace, acquired through a contract whose authority and coverage are explicit enough to preserve honest knowledge state, translated through a consumer-owned semantic boundary, and—when it can cause an external effect—correlated to an authoritative reconciliation surface. Provider protocol never acquires MPC business authority.**

---

## 4. Alternatives / Global Maximum

### A — Preserve current integration foundation and patch it

Keep generic connector/capability registry, direct Oracle reads, polling/manual-refresh assumptions and current provider result shapes, then add missing operations.

**Rejected.** This is a local maximum inside legacy structure and now retains a direct-Sankhya-DB target despite current official integration guidance.

### B — Build one universal external-source/provider framework

Normalize Mercado Livre, Sankhya and future providers behind universal Provider/Source/Resource/Capability/Operation objects.

**Rejected.** D0 rejects speculative generic integration frameworks and universal ERP models. Provider differences justify concrete contracts, not a lowest-common-denominator meta-model.

### C — Ground concrete providers behind consumer-owned ports

Use Mercado Livre and Sankhya to prove only the source-identity, coverage, capability, auth and reconciliation rules actually required; keep provider protocol provider-local and leave runtime machinery to D7.

**Recommended.** This preserves essential provider complexity without creating a new integration authority.

---

## 5. Candidate B1.1 — External protocol / semantic-boundary fence

D4 does not create an `Integrations` business domain.

1. A business consumer owns the semantic port it needs.
2. The concrete external adapter owns wire protocol, auth headers, endpoint vocabulary, pagination tokens, transport errors and DTO mapping.
3. Provider DTOs/status vocabularies do not cross into business contexts merely for convenience.
4. Provider-local shared HTTP/auth/pagination machinery is allowed as **mechanism**, but owns no business meaning.
5. D4 may describe Integration Support for a concrete operation; it does not create a universal provider capability registry without a proven current consumer/failure class.
6. One provider adapter may implement multiple consumer-owned ports; that does not imply one global generic `Provider` business port.
7. Assembly/registration mechanics remain D7/composition concerns.

Legacy plugin/self-registration/factory shapes therefore have no target authority by inheritance.

---

## 6. Candidate B1.2 — Mercado Livre Installation ↔ seller binding

A Mercado Livre `Marketplace Installation` binds one Organization participation/configuration to one authoritative external seller/account namespace.

1. External seller/user identity is provider-owned; it is neither Installation ID nor Organization.
2. Initial authorization/re-authorization must establish authenticated seller identity from provider-authoritative OAuth/user evidence (`user_id` and/or `/users/me` as appropriate).
3. Credential/token refresh for the same seller does not create a new Installation.
4. Authorization for a different seller must never silently rebind an existing Installation; mismatch fails closed and requires explicit portfolio/integration reconfiguration under owning business semantics.
5. Seller/user identifiers must not be narrowed to Int32; they remain opaque/source-qualified external references internally.
6. `site_id`, nickname and similar attributes are provider observations/configuration evidence, not tenant/Installation identity roots.
7. Downstream resource identities remain Installation + provider-native resource identity per D2.

Without this binding, later listing/order/shipment/write contracts can target the wrong seller while still being technically authenticated.

---

## 7. Candidate B1.3 — Sankhya SourceInstance binding

1. `SourceInstance` is Organization-qualified and identifies the logical Sankhya business-system namespace/environment whose native keys are referenced.
2. Production and sandbox/test are distinct external environments and must not be silently collapsed merely because entity/table names match.
3. Rotating `client_id`, `client_secret`, `X-Token` or connection credentials does not by itself create a new SourceInstance.
4. Changing to a materially different business-system namespace/environment requires explicit source rebinding/new SourceInstance treatment per D2.
5. Source identity is independent of concrete transport. Transport is protocol; SourceInstance is source identity.
6. B3 owns exact native-key/fact mapping and company/location/as-of qualification.

This lets D4 re-adjudicate Sankhya transport without rewriting D2 identity.

---

## 8. Candidate B1.4 — Material Sankhya transport conflict

### Finding

Official current Sankhya guidance says integrations should use API Gateway as the exchange standard and characterizes direct database integration as outside its integration guidelines for the applicable client population.

That materially conflicts with carried lower-level architecture:

- `ARCHITECTURE.md` stable constraint 5 points ERP reads to Oracle/ADR-006;
- ADR-006 says MPC reads ERP facts directly from Oracle;
- ADR-007 makes godror/OCI the canonical Oracle runtime.

D0 explicitly says exact Sankhya read/write capability/transport is D4 evidence to ratify independently. This is therefore a D4 challenge to a carried transport constraint, not evidence that D0–D3 business semantics are wrong.

### Proposed outcome

**RESTRUCTURE NOW — targeted carried-constraint correction, subject to independent review and operator approval.**

1. For new MPC↔Sankhya integration contracts, provider-sanctioned API Gateway is the default target transport.
2. Direct Oracle must not remain normative target merely because current code/ADR-006/007 already exist.
3. If B3 finds a materially required Product 1.0 fact that cannot be obtained correctly/fully through sanctioned APIs, it must surface that evidence explicitly.
4. A direct-DB exception then requires explicit current provider/customer entitlement/support evidence and a targeted architecture decision; it cannot appear as silent fallback.
5. Until operator-approved adjudication, existing `ARCHITECTURE.md`/ADR-006/ADR-007 remain current authority; this candidate modifies none of them.
6. No D0/D1/D2/D3 reopen is presently required because accepted semantics are transport-independent and D0 delegates exact Sankhya mechanics to D4.

B3 should decide concrete ERP facts against an admissible transport surface, not assume a transport that provider guidance challenges.

Independent review should test whether the public guidance applies fully to this MPC/client context and whether repository evidence contains a contractual/on-premise exception.

---

## 9. Candidate B1.5 — Authentication/credential lifecycle is protocol, not identity

### Mercado Livre

1. Authorization uses the provider's current server-side OAuth contract.
2. Access/refresh credentials are adapter/runtime secrets, never business-domain data or domain-event payload.
3. Token/user evidence verifies seller binding; credentials themselves are not seller identity.
4. Refresh-token rotation is provider credential lifecycle; an old refresh token is not assumed reusable after refresh.
5. Token lifetime is consumed from provider response/protocol rather than frozen as an MPC business constant.
6. Expired/revoked/invalid credentials make acquisition auth-invalid/unavailable, not source data absent.

### Sankhya

1. Current Gateway auth is OAuth 2.0 `client_credentials` plus `X-Token`.
2. Sandbox and production credentials/tokens are environment-bound and cannot be mixed.
3. Legacy appkey/login behavior is historical evidence, not default target auth for new work.
4. Credential changes do not alter SourceInstance when the same source namespace remains authorized.
5. Authentication failure is availability/auth state, not business absence.

### D7 fence

B1 does not choose secret-manager technology, credential schema/encryption, token-refresh locking/scheduling, callback route, retry/backoff or process placement.

---

## 10. Candidate B1.6 — Notification is trigger; authoritative reread owns current provider observation

For Mercado Livre notification-capable resources:

1. Notification receipt is external acquisition evidence, not MPC domain truth.
2. Resource/topic/provider-user/application metadata may locate/qualify the provider observation.
3. The authoritative resource read is performed before an owning business domain commits MPC meaning when current provider state matters.
4. Repeated notifications may cause repeated rereads without creating duplicate business truth.
5. Notification arrival order is not provider business order and cannot regress a newer authoritative observation merely because it arrived later.
6. `missed_feeds` is a bounded recovery aid, not complete durable provider history.
7. Notification outage/gap does not become a completeness claim; another authoritative observation path is required when completeness is material.
8. Non-authoritative callback fields stay acquisition/provenance metadata rather than canonical business status.

Webhook exposure/acknowledgement, queueing, retry receiver, reconciliation polling and cursor persistence remain D7.

---

## 11. Candidate B1.7 — Coverage/completeness is operation-scoped

D4 does not create one global sync-phase vocabulary. Each acquisition contract must expose enough information for its consumer to know what was actually observed.

### Point observation

- concerns one explicitly source-qualified resource/key;
- success establishes only what that endpoint legitimately authorizes for that resource at observation time;
- `not found` becomes known absent only when endpoint/identity semantics genuinely prove absence;
- auth/transport/rate-limit/timeout/parsing failure is unavailable, not absent.

### Enumerated observation

- coverage is only for the provider-defined scope actually traversed;
- all required pages/scan segments must complete before a completed-enumeration claim;
- provider cursors are protocol mechanics, not MPC identities;
- early stop/page failure/depth cap makes the observation partial;
- completed enumeration does not invent stronger snapshot-isolation semantics than provider docs establish;
- absence outside the enumerated scope proves nothing.

Mercado Livre makes this material because seller item enumeration above 1000 records uses `search_type=scan`/scroll rather than ordinary paging.

### Delta/change observation

- proves only the source-defined change set under explicit preconditions/window;
- does not by itself prove every current object was observed;
- source prerequisites for change tracking are part of the coverage claim;
- empty delta with unknown/disabled prerequisites must not become “nothing changed”.

Sankhya makes this concrete because `modifiedSince` depends on `logAlteracoesTabelas` and may return zero when no log information is present.

### Notification trigger

Carries no global completeness claim by itself.

### Contract-shape rule

B1 does not require a universal `Coverage` entity/enum. Each consumer-owned result carries/references the smallest operation-specific provenance/coverage needed for correctness. Freshness remains consumer/use-sensitive under D3.

---

## 12. Candidate B1.8 — Preserve D3 knowledge/failure semantics

1. Adapters may normalize provider errors for consumers, but may not map unavailability to plausible business values.
2. Auth failure, rate limiting, timeout, provider/gateway outage, malformed response and incomplete traversal do not become known empty/zero/false.
3. `known empty/absent` requires affirmative source semantics for the exact queried scope.
4. Unsupported integration operation is explicit rather than masquerading as empty data.
5. If provider behavior establishes only uncertainty, uncertainty remains explicit.
6. Provider-native error DTO/text stays at the external boundary; consumer-facing meaning is consumer-owned.
7. PII/secrets are not retained/propagated merely because payload/error contains them.
8. Provenance preserves enough source identity and acquisition/source time for D2/D3 lineage/freshness without requiring universal raw-payload storage.

Exact error class hierarchy and MPC HTTP encoding are later concerns.

---

## 13. Candidate B1.9 — Integration Support != Provider Effective Capability != Effective Business Capability

### Level 1 — Integration Support / Descriptor — D4 technical meaning

Can this concrete adapter/protocol attempt/read/write this external operation class at all?

### Level 2 — Provider Effective Capability / Requirement — source evidence translated by D4

For this seller/resource/mode/context, what does the provider currently allow/require, and what provider prerequisite/artifact/state applies?

### Level 3 — Effective Business Capability — consuming business authority

Given provider evidence plus MPC readiness/policy/authorization/current business state, may/should MPC perform the action now?

D4 does not own Level 3.

No provider `capability=true` may bypass Readiness, Offering, Availability, Fulfillment, Governance or another D1 owner. A business domain likewise cannot fabricate provider support when D4 cannot establish it.

---

## 14. Candidate B1.10 — Admission gate for later external-effect contracts

B1 does not choose concrete writes. Every later B2/B3 external-effect contract must identify at least:

1. target Installation/SourceInstance-qualified resource;
2. consumer-owned semantic intent/correlation anchor;
3. material provider preconditions/requirements;
4. what the provider response actually proves — rejection, accepted submission, synchronous effect, pending work or ambiguity;
5. when outcome may be ambiguous after possible acceptance;
6. authoritative reread/reconciliation surface;
7. member-level outcome semantics when multi-target work can partially succeed;
8. provider/source occurrence/result discriminator only where same-vs-distinct correctness needs it.

> **Transport success is not promoted to `converged` merely because it returned 2xx.**

Current legacy ML stock-write code mapping 2xx directly to `applied` is evidence to re-adjudicate in B2, not target precedent.

Retry/backoff/idempotency-store/attempt-table mechanisms remain D7; D3 no-blind-retry/ambiguity semantics remain authority.

---

## 15. Legacy ADR disposition proposed by B1

No registry is changed by this candidate.

### ADR-003 — `reopened — D4/D9`

Rehome only the real D4 prerequisite: authenticated provider calls require valid credential lifecycle and provider identity binding. Do not carry the old strict `OAuth -> fee sync -> frontend` implementation sequence. Fee belongs B4; frontend D6; D9 residue stays D9.

### ADR-004 — `reopened — D4`

Proposed: supersede legacy plugin/self-registration/framework target meaning. Preserve only provider-specific protocol near provider adapters and consumer-owned semantic ports. No generic auth factory/fee-sync registry/integration catalog becomes target authority.

### ADR-010 — `reopened — D4/D7`

Preserve honest freshness/failed-refresh meaning; reject mission-time `polling-only/no-webhooks` as target provider rule. B1 adopts notification→reread and explicit coverage. Receiver/scheduler/poll cadence stays D7.

### ADR-014 — `reopened — D4`

Deferred to B4. On-demand/local-Docker/no-history mission shape does not constrain B1.

### ADR-015 — `reopened — D4`

Partially constrained by B1; final disposition B2. Preserve source-qualified identity, honest coverage and reread. Do not carry composite ID format, one read-model table/module, manual-refresh runtime or “absent from completed pull = closed”.

### ADR-020 — `reopened — D4`

Deferred to B4. Generic `CollectorPort` interface is not imported into B1.

### ADR-032 — `reopened — D4`

Deferred to B4. Current env flag/default-off behavior is runtime/current-state evidence, not target D4 authority.

### Newly implicated carried ADR-006 / ADR-007

They are not currently `reopened — D4`, but B1 found current official-provider evidence challenging their direct-Oracle target meaning. Proposed: targeted D4 reopen/re-adjudication, without silent registry change; if Gateway/API is accepted, godror/Oracle remain historical/current-state evidence unless an explicit supported exception is later proven.

---

## 16. Explicit deferrals

### D4-B2 — Mercado Livre Operational Contract

Own listings/items/variations/catalog relations, listing enumeration/point reads, stock/availability reads+writes, price/listing controlled writes, sales/orders, fulfillment modes/requirements/artifacts/deadlines, concrete provider write acceptance/reconciliation, and composite behavior only if the first flow requires it.

### D4-B3 — Sankhya Business-System Contract

Sanctioned target surface after B1 transport adjudication; exact product/native key, inventory company/location/as-of facts, cost/tax evidence, Business Order Intent materialization, Invoicing Intent materialization, native result keys and authoritative reread/reconciliation; explicit unsupported facts/commands where necessary.

### D4-B4 — Market / Economics / Settlement Contract

Official market/competitor evidence, catalog-offers/price-to-win where justified, fee evidence, provider financial movement/settlement/adjustment/refund evidence and provenance/completeness/economics correlation.

D5–D8 remain untouched.

---

## 17. Proof strategy / strongest counterexamples

B1 should not be accepted unless it survives at least these counterexamples:

1. Same ML seller, new tokens -> Installation stable.
2. Existing Installation authorized as different ML seller -> fail closed; no silent rebind.
3. ML user ID exceeds Int32 -> source reference still representable.
4. Sankhya sandbox credentials mixed with production -> auth failure without collapsing source environments or converting facts to absent.
5. Sankhya credentials rotate for same production namespace -> SourceInstance stable.
6. Duplicate ML resource notification -> repeated reread safe; no duplicate business truth.
7. Notification missed outside bounded recovery -> no false complete history claim.
8. Seller has >1000 ML items -> ordinary offset success cannot be called complete; provider scan semantics required.
9. One enumeration page/scan segment fails -> observation partial; unseen resources cannot be closed by absence.
10. Sankhya `modifiedSince` returns zero while log coverage is unproven -> zero is not “no changes”.
11. Public ML search exposes referential quantity -> wrong endpoint cannot become exact inventory authority because field name looks useful.
12. Adapter implements price write but seller/resource is not provider-eligible -> Integration Support does not become provider/business capability.
13. Auth expires -> external fact becomes unavailable/stale as appropriate, never zero/empty/false.
14. Endpoint 404 does not semantically prove deletion -> no manufactured known absence.
15. Provider receives write then connection drops -> ambiguity/reconciliation required; blind retry forbidden.
16. Provider returns 2xx for accepted async work -> no convergence claim unless source contract proves it.
17. Multi-target effect partially succeeds -> member outcomes remain representable.
18. A second marketplace arrives -> consumer-owned ports can be implemented without universal provider entity graph today.
19. Direct Oracle is technically faster -> convenience alone cannot override current provider integration guidance.
20. B3-required Sankhya fact is unavailable through Gateway -> stop/exception decision, not silent direct-DB fallback.

---

## 18. Reopen / stop triggers

1. ML authorization cannot establish a stable seller namespace compatible with D2 -> targeted D2 review.
2. Required external resource identity cannot fit Installation/SourceInstance + native-key model -> targeted D2 review.
3. Concrete interaction needs a new semantic business dependency absent from D1 -> targeted D1 review.
4. Provider effect semantics cannot fit D3 accepted/rejected/pending/ambiguous + reread/reconciliation -> targeted D3 review.
5. Provider evidence makes the Product 1.0 loop materially impossible/different -> targeted D0 review.
6. Current Sankhya contractual/customer evidence explicitly supports direct DB for this context contrary to public guidance -> revisit B1.4 before changing carried Oracle authority.
7. B3 proves a required fact cannot be obtained correctly through sanctioned Gateway/API and no supported equivalent exists -> STOP / targeted transport-exception decision.
8. A concrete second provider creates a repeated technical failure unsolved by consumer-owned ports + provider-local mechanism -> consider the smallest shared mechanism then.

Framework preference/current code convenience/speculative future providers are not reopen evidence.

---

## 19. Proposed batch outcome

If independent review finds no material contradiction:

- **D0–D3: CURRENT STRUCTURE CONFIRMED; no reopen.**
- **D4-B1 grounding: bounded target additions only.**
- **ADR-004 plugin framework:** supersede/rehome D4 meaning.
- **ADR-010 polling-only D4 meaning:** supersede/rehome; D7 residue remains.
- **ADR-003 auth prerequisite:** bounded rehome; later-stage residues remain.
- **ADR-015:** final listing disposition deferred B2.
- **ADR-014/020/032:** deferred B4.
- **ADR-006/007 direct-Oracle target meaning:** `RESTRUCTURE NOW` candidate due current Sankhya guidance, with no canonical change until operator approval.

Target shape:

```text
accepted D1/D2 consumer meaning
        ↓ consumer-owned port
D4 provider/business-system adapter
  - source namespace binding
  - auth/protocol
  - operation-specific authority + coverage
  - provider capability/requirement evidence
  - authoritative reread/reconciliation
        ↓
external system

D7 later supplies runtime mechanics around this boundary.
```

No generic Integration domain, Provider entity graph, connector registry, universal capability object, global sync mode or runtime topology is required by B1.

---

## 20. Independent review focus

Please challenge independently:

1. Does B1 preserve D1/D2/D3 authority without provider DTO/status leakage?
2. Is Installation↔seller binding fail-closed without over-specifying OAuth runtime?
3. Is SourceInstance transport-independent in the right way?
4. Is the Sankhya Gateway-vs-direct-Oracle finding strong enough to reopen carried ADR-006/007, or is material client-specific evidence missing?
5. Are point/enumeration/delta/notification coverage distinctions necessary correctness semantics or accidental abstraction?
6. Does any B1 rule leak into D7 runtime?
7. Does the capability fence preserve provider evidence without creating a generic capability framework?
8. Is the external-effect admission gate the minimum necessary D4 contract, or does it prematurely constrain B2/B3?
9. Does any legacy ADR meaning need preservation that this candidate wrongly discards?
10. Is any B1 decision actually B2/B3/B4/D5–D8 and should be deferred?
11. Are there material D0–D3 contradictions missed here?
12. What is the strongest future-cost/YAGNI objection?

Reviewer findings are evidence, not authority. Canonical architecture changes only after independent adjudication and explicit operator approval.
