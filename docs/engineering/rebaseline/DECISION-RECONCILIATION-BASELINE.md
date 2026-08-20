# Decision Reconciliation Baseline

> **Status:** ACCEPTED / CANONICAL  
> **Role:** on-demand routing map from superseded architectural ideas to their current accepted semantic homes  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Accepted:** 2026-08-18  
> **Operator ratification:** 2026-08-18

## 1. When to read this file

Do **not** read this file during the default fresh-session bootstrap. Start with `AGENTS.md` and `docs/README.md`.

Read this map only when a concrete task needs to answer one of these questions:

- which generation of an iterated decision is current;
- where an older ADR, module shape or design idea was rehomed;
- whether a technical proposal would silently resurrect superseded authority;
- which later stage owns a genuinely deferred mechanism or proof.

This file routes decisions. It is not a second semantic architecture, ADR-status registry, roadmap or current-program status page.

## 2. Authority by scope

1. `docs/README.md` — current program status, allowed/blocked work, exact next action and selective read routes;
2. `ARCHITECTURE.md` — stable cross-stage constraints;
3. this baseline — decision-generation and legacy-idea routing only;
4. `docs/architecture/decisions/README.md` — sole ADR file status/disposition authority;
5. accepted D-stage artifacts — detailed semantic authority in their stage scope;
6. Evidence Register, code, schemas, APIs, tests, runtime and Git history — supporting evidence unless accepted authority explicitly carries their meaning.

If this map disagrees with an accepted semantic home, this map is defective and must shrink or be corrected; it never overrides that home. Absence from this map is never permission to violate accepted authority.

## 3. Current semantic homes

| Scope | Current semantic home | Technical stages must not infer from… |
|---|---|---|
| Product mission, actors and Product 1.0 boundary | D0 + `ARCHITECTURE.md` | dashboard-only or generic-provider-hub legacy shapes |
| Business authorities and ownership edges | D1 | legacy packages/modules or provider resource grouping |
| Organization, Principal, identities, knowledge and durable data ownership | D2 | bare external IDs, a Product mirror or ambient tenant defaults |
| Q/C/E/P, propagation, idempotency, recovery and projections | D3 | transport choice, polling phases or cross-context SQL |
| Mercado Livre, Sankhya, provider capability/coverage and external effects | D4 | provider DTOs, Direct Oracle fallback or a generic integration owner |
| Publication input and listing authoring | D4-R1 | Product/PIM mastery, a rules engine or a separate publication aggregate |
| Semantic Product API laws | D5-B1 | provider/integration ontology or current legacy routes |
| Concrete Product operation/wire surface | accepted D5-B2 W1–W4 + Technical Ingress + OpenAPI tooling | manual SDK/OpenAPI duplication or runtime framework preference |

These homes define semantics, not a required number of services, processes, databases, repositories or first-release features.

## 4. Historical idea → current home

| Superseded or older idea | Current home / implementation instruction |
|---|---|
| Marketplace dashboard as the product | D0 control-plane product loop; do not implement dashboard-only architecture |
| Product mirror / MPC Product master | D2 external source-qualified Product identity + D1 Readiness consumers; do not create an MPC Product master |
| Separate canonical Tenant | D2 Organization root |
| Generic Integration business domain | Portfolio business meaning + D4 protocol + D7 mechanics |
| Provider plugin/self-registration business framework | concrete D4 adapters; shared technical machinery only when later proven |
| Direct Oracle/godror Sankhya target | D4 sanctioned API Gateway only; Direct Oracle is current-state evidence, not fallback |
| `SELLER_SKU == CODPROD` identity law | D2/D4 evidence + Readiness correspondence; never canonical identity |
| CODPROD+EAN unattended-link formula | D2 corroboration safety + Readiness policy |
| Generic Mutation business owner | domain-local intents + Governance + D3/D7 execution safety |
| Generic divergence ledger as truth owner | source-domain correctness + Operational Work lifecycle |
| `sync`/polling phase as Product semantics | domain freshness/coverage + D4/D7 acquisition/runtime mechanics |
| Provider DTO/resource model in core | consumer-owned semantic ports + adapter-local protocol |
| Lowest-common-denominator marketplace model | shared semantics where meanings align + source-qualified provider-rich evidence |
| Generic CollectorPort/market-source framework | explicit source admissibility + later bounded D7 mechanism if proven |
| Generic Customer/Address master for Sankhya | bounded Materialization Party/Destination realization |
| MPC tax engine | sanctioned Sankhya fiscal engine + Commercial Economics interpretation |
| Global Fee/Payment/Settlement MPC entity | source-qualified external movements + Economics attribution/reconciliation |
| `PublicationPreparation` aggregate | Offering-owned `ListingIntent` draft + Readiness query |
| `SourceProductObservation` business service | D4 acquisition feeding consumer-owned ports; no generic source-Product owner |
| Generic listing transformation/rules engine | `FOLLOW_SOURCE | EXPLICIT_OVERRIDE`; reopen only on repeated evidenced need |
| Dedicated AI business/API authority | D2 automation Principal + ordinary Product API/Governance |
| Product API listing quantity owned by Offering | Availability-owned meaning; joint technical realization only |
| Manual OpenAPI + manual SDK authority | one OpenAPI wire authority + derived/conformant client/server projections |
| Compatibility/versioning for the legacy API | hard cutover until a real entitled production consumer appears |
| Generic Workflow/SLA/Policy domain | owner-local semantics + Governance/Work; mechanism never becomes business authority |
| Provider callback/webhook as current truth | D4 acquisition pointer/evidence + authoritative reread when material |
| Provider `2xx` as convergence | accepted/submitted outcome only; owner-specific reread and reconciliation establish convergence |

Git history preserves the original formulations and retired ADR files. Do not restore them into the active tree merely to retain history.

## 5. High-risk implementation reconciliation guard

The list below is intentionally compact and non-exhaustive. Detailed meaning remains in `ARCHITECTURE.md` and the accepted D-stage artifacts.

Technical stages must not silently re-decide by convenience:

1. Product is external/source-qualified, not an MPC master.
2. Organization is the tenant/isolation root and cross-Organization paths fail closed.
3. One business meaning has one semantic authority; mechanism does not acquire authority.
4. D1 business authorities and accepted semantic edges remain binding.
5. Q/C/E/P meaning, recoverable consequential propagation and material occurrence recovery remain distinct.
6. Domain-local intents do not become a generic mutation/workflow authority.
7. Governance authorization, ordinary Permission and owner business disposition remain separate.
8. Work lifecycle never becomes the originating business truth; material actionable conditions cannot become ownerless silent state.
9. External/provider identities remain source-qualified and distinct from MPC canonical identities.
10. Consumers own semantic ports; adapters own provider/business-system protocol.
11. Sankhya target transport is the sanctioned API Gateway; Direct Oracle is not a convenience fallback.
12. Possible external acceptance is never followed by blind retry.
13. Provider richness is preserved for named needs without raw DTO mirroring or fabricated equivalence.
14. Provider PII is minimized.
15. Missing evidence never authorizes fabricated data or an unadjudicated scraping/source path.
16. Externally governed obligations and policy provenance cannot be silently relaxed by MPC.
17. Offering, Readiness and Availability retain their accepted publication ownership split.
18. `ListingIntent` is the one create/edit listing-authoring identity.
19. `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` is the baseline; recurring automation never silently reverses a standing human decision.
20. No hidden PIM, source-observation, rules, connector-platform or AI-specific business framework is introduced.
21. Joint technical realization may serialize multiple owner-issued meanings without transferring ownership or inventing a semantic edge.
22. Each owner evaluates its own convergence after provider effects.
23. Earlier approval is not permanent execution authority after materially governing conditions drift.
24. Multi-target intent, authorization, attempted scope and member-level outcomes remain distinguishable.
25. Cutover/recovery cannot silently discard pending reactions required for accepted lifecycle progression.
26. Product API keeps explicit Organization scope and source-qualified external identity.
27. OpenAPI remains the one machine-readable Product wire authority; no second manual SDK/wire authority survives.
28. Hard cutover remains valid while no entitled production compatibility consumer exists.

A technical design that appears to require violating one of these is evidence for a targeted architecture reopen—not permission for a workaround.

## 6. Deferred owner/stage map

| Later stage | Questions deliberately owned there |
|---|---|
| D5 | current canonical Product OAD authoring/proof and D5-B2 closure only |
| D6 | screens/navigation/editor topology, client composition and frontend package structure |
| D7 | router/server, processes, workers, schedulers, transactions, outbox/queue, cursor/lease/lock, secrets, RLS realization, idempotency persistence, retry/rate control, media storage and deployment topology |
| D8 | selected end-to-end success/failure/retry/reconciliation proofs, including real controlled external effects when authorized |
| D9 | final adversarial whole-architecture contradiction, overbuild and under-specification review |

Unknown or deferred is never permission to invent a plausible default. Reopen only the smallest responsible authority when real evidence invalidates an assumption.

## 7. ADR and active-tree reconciliation

The ADR registry is the only authority for ADR file status/disposition. Fully rehomed pre-rebaseline ADRs remain in Git history, not the active tree. Retained legacy residues remain only because a later stage still owns a concrete unresolved mechanism/proof/transition question.

`ADR-035` remains transition authority through D0–D9. Its embedded 2026-08-14 status tables are historical snapshot evidence; current ADR disposition comes from the registry and later accepted stages. ADR numbers are never reused.

## 8. Usage and reopen triggers

Fresh-session routing is:

```text
AGENTS.md
→ docs/README.md
→ only the task-specific authority/evidence selected there
→ this reconciliation map only when legacy decision-generation conflict exists
```

Reconcile this map only when:

- an accepted D-stage decision is amended or reopened;
- a new target ADR changes decision routing;
- executable proof invalidates an accepted architectural assumption;
- a Product requirement changes D0/D1 ownership;
- a retained legacy ADR is adjudicated and leaves the active tree;
- D0–D9 closes and transition machinery retires.

Do not reopen it for implementation naming, package layout, framework preference or rediscovery of a retired Git-history decision.

## 9. Reconciliation verdict

The accepted D0→D4/D4-R1 + D5 decision set remains the current coherent target authority. This file preserves routing and implementation guards without requiring a fresh session to reread or reconstruct every accepted stage.

Exact program status and the next action live only in `docs/README.md`.
