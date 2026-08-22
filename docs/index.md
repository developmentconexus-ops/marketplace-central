# Marketplace Central — Documentation Router

> **Role:** task/intention routing only. Mutable program status lives only in [`roadmap.md`](roadmap.md).

## Fresh actor route

Read, in order:

1. [`../AGENTS.md`](../AGENTS.md)
2. this file
3. [`roadmap.md`](roadmap.md)
4. one or two owning documents for the concrete task

Default task pack: **5 files or fewer**. Do not recursively read phase history, ADRs, Evidence, Git history, or removed runtime unless a concrete question requires them.

Pre-standard current-status/read-order prose inside accepted D-stage, architecture or retained ADR artifacts is a **frozen routing snapshot**. This router and `docs/roadmap.md` supersede that prose for navigation/status only; the owning documents' accepted Product/architecture semantics remain unchanged.

## Task routes

| Task / question | Smallest starting authority | Do not read by default |
| --- | --- | --- |
| Current stage, allowed/blocked work, exact next action | [`roadmap.md`](roadmap.md) | phase history |
| Product mission, scope, actors, Product 1.0 boundary | [D0 Product/System Definition](engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md) | D1–D5 |
| Domains, semantic owners, allowed ownership edges | [D1 Domains/Boundaries](engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md); for the current bounded Performance amendment also [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md) | implementation history |
| Organization, Principal, identity, data ownership, knowledge/provenance | [D2 Identity/Tenant/Data Ownership](engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md) | provider/runtime details |
| Human-readable Organization/Principal/AccessRole presentation identity | [D2-R1 Presentation Identity](engineering/rebaseline/D2-R1-PRESENTATION-IDENTITY.md) | profile/directory/runtime mechanics |
| Marketplace Performance Intelligence bounded D0–D5 repair, historical performance evidence, retail-media read contract | [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md) | D7 runtime/storage realization, Ads management, AI implementation |
| Q/C/E/P, communication, events, recovery, projections | [D3 Communication/Events](engineering/rebaseline/D3-COMMUNICATION-EVENTS.md) | D4 provider detail |
| Mercado Livre, Sankhya, provider boundaries/effects/reconciliation | [D4 External Integrations](engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md) | old adapters/runtime |
| Publication input, `ListingIntent`, readiness/source-following/override | [D4-R1 Publication Input](engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md) | generic PIM assumptions |
| Product API semantic laws | [D5 API](engineering/rebaseline/D5-API.md) | D6–D9 |
| Operational collection/read projection repair for D6-R2 triage (`OP-READ-01`) | [D5-R2 Operational Read Projection Repair](engineering/rebaseline/D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) | broad D5 history; implementation/runtime |
| Frontend interaction/authority model, screen→Product capability mapping, frontend topology | [D6 Frontend](engineering/rebaseline/D6-FRONTEND.md) | D7–D9 and removed frontend runtime |
| Frontend product-experience planning process / D6-R2 execution method | [Frontend Product Experience Planning Method v2.1](development/frontend-product-experience-planning-method.md) — reusable methodology, not stage/status authority | stage semantics; use `roadmap.md` + owning authority |
| Runtime/process topology and whole-D7 integration | [D7 Runtime / Jobs / Transactions](engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md) | D8–D9 and removed runtime |
| PostgreSQL Organization isolation, RLS, transaction scope, idempotency, ETag/revision | [D7-B PostgreSQL Isolation & Transactions](engineering/rebaseline/D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) | D7-C–E and full schema census |
| Durable work, River, retries, external-effect ambiguity/reconciliation | [D7-C Durable Work & External Effects](engineering/rebaseline/D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) | auth/deployment and generic workflow assumptions |
| Human session/CSRF/OIDC and A/S machine bearer realization | [D7-D Authentication / Session / CSRF](engineering/rebaseline/D7-D-AUTHENTICATION-SESSION-CSRF.md) | D7-E and generic IAM |
| HTTP runtime validation, byte storage, secrets, migrations, telemetry, deployment/backup/proof | [D7-E Operability / Deployment / Proof](engineering/rebaseline/D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) | D8–D9 and provider-specific infra by preference |
| Current decision-generation reconciliation | [Decision Reconciliation Baseline](engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md) | all phase history |
| Stable cross-stage architecture constraints | [`../ARCHITECTURE.md`](../ARCHITECTURE.md) | all phase documents |
| ADR disposition / retirement trigger | [ADR Registry](architecture/decisions/README.md), then only the named ADR | all ADRs |
| Historical/current-state Evidence for a concrete decision | [Evidence Register](engineering/rebaseline/EVIDENCE-REGISTER.md), then exact cited ref/path | repository archaeology |
| Repository-local engineering/Git/CI/proof rules | [`development/engineering-rules.md`](development/engineering-rules.md) | organizational Method text copied locally |
| Production coding, technology research, framework/dependency evaluation and source hierarchy | [Evidence-Grounded Production Engineering for LLM Agents](development/evidence-grounded-production-engineering-for-llm-agents.md) — derived guide, non-authoritative | broad recursive research |

## Product OAD bounded subpacks

Canonical Product OAD authoring is material and spans several accepted D5 homes. Do **not** open all of them at once. Use one subpack at a time; every subpack remains within the five-file repository limit.

**Start / authority and tooling — 4 files total with the bootstrap:**

- this router + [`roadmap.md`](roadmap.md);
- [D5 API](engineering/rebaseline/D5-API.md);
- [OpenAPI Wire Authority / Tooling](engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md).

**Operation/path authoring — add at most 3 owning files:**

- [Operation Admission Matrix](engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md);
- [Product Operation Surface](engineering/rebaseline/D5-B2-PRODUCT-OPERATION-SURFACE.md);
- [W1 Wire Contract](engineering/rebaseline/D5-B2-WIRE-CONTRACT.md).

**Schema/collection authoring — add at most 2 owning files:**

- [W2 Schema Grammar](engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md);
- [W3 Collection Grammar](engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md).

**Access/technical exclusion proof — add at most 2 owning files:**

- [W4 Permission / Client-Class Enforcement](engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md);
- [Technical Ingress](engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md).

For the D6-discovered Performance repair, start from [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md) and switch to one D5 subpack only when the exact wire/proof question requires it. D6-R1 supersedes only its explicitly bounded 95→99 / 29→30 / 26→28 amendments; it does not generally reopen accepted D5.

For the D6-R2-discovered operational triage repair, start from [D5-R2 Operational Read Projection Repair](engineering/rebaseline/D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) and switch to W2/W3 only for exact schema/collection grammar. D5-R2 supersedes only its bounded operational ListItem/filter amendments.

Switch subpacks as the concrete authoring question changes. Read D0–D4 only when the D5 package cannot answer a concrete contract question without reopening earlier semantics.

## D6 Frontend bounded subpacks

D6 must consume accepted Product/API authority without re-reading the full rebaseline. Use one subpack at a time.

**Start / interaction authority — 5 files total with the bootstrap:**

- [`../AGENTS.md`](../AGENTS.md);
- this router;
- [`roadmap.md`](roadmap.md);
- [D6 Frontend](engineering/rebaseline/D6-FRONTEND.md);
- [`../ARCHITECTURE.md`](../ARCHITECTURE.md).

**Screen / Product capability mapping — replace `ARCHITECTURE.md` with the machine-readable Product authority:**

- [D6 Frontend](engineering/rebaseline/D6-FRONTEND.md);
- [`../contracts/api/product/openapi.yaml`](../contracts/api/product/openapi.yaml).

Use D5 narrative documents only when the canonical OAD plus D6 authority cannot answer a concrete wire-semantic question. For the bounded human-readable Organization/Principal/AccessRole label rule, switch to D2-R1 rather than inventing frontend-local identity. For the strategy/performance workspace and repaired 99/30 Product surface, switch to D6-R1 rather than inventing frontend-local analytics.

**Frontend dependency/topology research — replace the OAD/architecture owner with the derived research guide only when a concrete technology question exists:**

- [D6 Frontend](engineering/rebaseline/D6-FRONTEND.md);
- [Evidence-Grounded Production Engineering for LLM Agents](development/evidence-grounded-production-engineering-for-llm-agents.md).

Do not select D7 runtime/router/database/deployment mechanics from D6 research.

## D7 Runtime bounded subpacks

Start from the D7 task-route row above and add only the one accepted prior owner needed for the concrete invariant. Do not read D7-A→D7-E recursively by default; whole-D7 review is the only task that intentionally composes the entire set.

Do not begin D8 golden-flow choreography or Product implementation from D7 research.

## Organizational standards

Canonical, external to this repository:

- `developmentconexus-ops/conexus-methodology/METHOD.md` — **DevelopmentConexus Engineering Method v1.0.0**
- `developmentconexus-ops/conexus-methodology/REPOSITORY-STANDARD.md` — **DevelopmentConexus Repository Standard v1.0.0**

No local copy is Product or repository authority.
