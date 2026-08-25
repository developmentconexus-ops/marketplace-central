# Marketplace Central — Documentation Router

> **Role:** navigation aid only. Mutable program status lives only in [`roadmap.md`](roadmap.md). This document provides starting points, not reading limits.

## Start

Read [`../AGENTS.md`](../AGENTS.md) and [`roadmap.md`](roadmap.md), then follow the adopted local method that fits the work:

- material engineering: [`development/engineering-method.md`](development/engineering-method.md)
- frontend Product Experience planning: [`development/frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md)

Use the routes below to find useful owners. There is no fixed file count or owner count. Expand into any Product, architecture, contract, Evidence, research, Git history, code, runtime, or external source that can materially change or falsify the conclusion.

## Task routes

| Task / question | Useful starting references |
| --- | --- |
| Current stage, allowed/blocked work, exact next action | [`roadmap.md`](roadmap.md) |
| Engineering reasoning, Global Maximum, root cause, invariant, alternatives, proof | [`development/engineering-method.md`](development/engineering-method.md) |
| Product mission, scope, actors, Product 1.0 boundary | [D0 Product/System Definition](engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md) |
| Domains, semantic owners, allowed ownership edges | [D1 Domains/Boundaries](engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md); [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md) when relevant |
| Organization, Principal, identity, data ownership, knowledge/provenance | [D2 Identity/Tenant/Data Ownership](engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md) |
| Human-readable Organization/Principal/AccessRole presentation identity | [D2-R1 Presentation Identity](engineering/rebaseline/D2-R1-PRESENTATION-IDENTITY.md) |
| Communication, events, recovery, projections | [D3 Communication/Events](engineering/rebaseline/D3-COMMUNICATION-EVENTS.md) |
| Mercado Livre, Sankhya, provider boundaries/effects/reconciliation | [D4 External Integrations](engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md) |
| Publication input, `ListingIntent`, readiness/source-following/override | [D4-R1 Publication Input](engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md) |
| Product API semantic laws | [D5 API](engineering/rebaseline/D5-API.md) |
| Product OAD authoring and wire proof | [D5 API](engineering/rebaseline/D5-API.md), [OpenAPI Wire Authority / Tooling](engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md), and the relevant D5-B2 operation/schema/access owner |
| Operational collection/read projection repair for D6-R2 triage | [D5-R2 Operational Read Projection Repair](engineering/rebaseline/D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) |
| Frontend interaction/authority model and topology | [D6 Frontend](engineering/rebaseline/D6-FRONTEND.md) |
| Frontend Product Experience planning / D6-R2 | [`development/frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md), [D6 Frontend](engineering/rebaseline/D6-FRONTEND.md), and the current block's relevant authority/evidence |
| Current D6-R2 accumulated authority/evidence | [D6-R2 Authority Route](engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md) |
| Runtime/process topology and whole-D7 integration | [D7 Runtime / Jobs / Transactions](engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md) |
| PostgreSQL Organization isolation, RLS, transactions, idempotency, ETag/revision | [D7-B PostgreSQL Isolation & Transactions](engineering/rebaseline/D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) |
| Durable work, retries, external-effect ambiguity/reconciliation | [D7-C Durable Work & External Effects](engineering/rebaseline/D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) |
| Human session/CSRF/OIDC and machine bearer realization | [D7-D Authentication / Session / CSRF](engineering/rebaseline/D7-D-AUTHENTICATION-SESSION-CSRF.md) |
| HTTP validation, storage, secrets, migrations, telemetry, deployment/backup/proof | [D7-E Operability / Deployment / Proof](engineering/rebaseline/D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) |
| Golden-flow GF-02 revalidation | [D8-R2 GF-02 Operational Read Revalidation](engineering/rebaseline/D8-R2-OPERATIONAL-READ-REVALIDATION.md) |
| Current decision-generation reconciliation | [Decision Reconciliation Baseline](engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md) |
| Stable cross-stage architecture constraints | [`../ARCHITECTURE.md`](../ARCHITECTURE.md) |
| ADR disposition / retirement trigger | [ADR Registry](architecture/decisions/README.md), then the relevant ADR |
| Evidence for a concrete decision | [Evidence Register](engineering/rebaseline/EVIDENCE-REGISTER.md), then the cited evidence |
| Repository-local Git/CI/proof rules | [`development/engineering-rules.md`](development/engineering-rules.md) |
| Production coding, technology/dependency research, integration/proof | [`development/engineering-method.md`](development/engineering-method.md) plus the optional [Evidence-Grounded Production Engineering guide](development/evidence-grounded-production-engineering-for-llm-agents.md) when its detailed technology/proof lenses help |

## Navigation principle

These links are entry points, not fences. Do not omit materially relevant context because it sits outside a suggested starting set. Do not read irrelevant material merely to satisfy ceremony.

Current repository authority owns Product meaning. The adopted local methods govern engineering/frontend reasoning. Evidence and research may challenge accepted authority through those methods; they do not silently replace it.
