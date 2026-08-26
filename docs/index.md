# Marketplace Central — Documentation Router

> **Role:** selective navigation aid. Mutable program status lives only in [`roadmap.md`](roadmap.md). This file points to current owners/evidence; Git history preserves how those decisions were reached.

## Start

Read [`../AGENTS.md`](../AGENTS.md) and [`roadmap.md`](roadmap.md), then use the local method that fits the work:

- material engineering: [`development/engineering-method.md`](development/engineering-method.md)
- frontend Product Experience planning: [`development/frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md)

Start with the smallest likely owner set below. Expand only when another owner, Evidence, Git history, runtime behavior or external source can materially change or falsify the conclusion.

## Current-owner routes

| Task / question | Current starting owner(s) |
| --- | --- |
| Current stage, active increment, allowed/blocked work, exact next action | [`roadmap.md`](roadmap.md) |
| Engineering reasoning, Global Maximum, root cause, alternatives, proof | [`development/engineering-method.md`](development/engineering-method.md) |
| Product mission, Product 1.0 scope and actors | [D0 Product/System Definition](engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md) |
| Domain ownership and semantic edges | [D1 Domains/Boundaries](engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md); [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md) only when that owner is relevant |
| Organization, Principal, source-qualified identity, data ownership, knowledge/provenance | [D2 Identity/Tenant/Data Ownership](engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md) |
| Human-readable Organization/Principal/AccessRole presentation identity | [D2-R1 Presentation Identity](engineering/rebaseline/D2-R1-PRESENTATION-IDENTITY.md) |
| Communication, events, recovery and projections | [D3 Communication/Events](engineering/rebaseline/D3-COMMUNICATION-EVENTS.md) |
| Mercado Livre, Sankhya, provider/source boundary and external effects | [D4 External Integrations](engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md) |
| Publication requirements, Readiness, source-following/override, ListingIntent integration seam | [D4-R1 Publication Input](engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md) |
| Product API semantic laws | [D5 API](engineering/rebaseline/D5-API.md) |
| Product operation/path grammar | [D5 W1 Wire Contract](engineering/rebaseline/D5-B2-WIRE-CONTRACT.md) |
| Product schema/read-write grammar | [D5 W2 Schema Grammar](engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md) |
| Product collections/search/cursor grammar | [D5 W3 Collection Grammar](engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md) |
| Product Permission / Principal-client-class enforcement | [D5 W4 Permission Enforcement](engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md) |
| Product OAD authoring/tooling/proof | [OpenAPI Wire Authority / Tooling](engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md) + [`contracts/api/product/openapi.yaml`](../contracts/api/product/openapi.yaml) |
| Operational collection/read repair used by D6-R2 | [D5-R2 Operational Read Projection Repair](engineering/rebaseline/D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md) |
| Frontend architecture, client-state authority and topology | [D6 Frontend](engineering/rebaseline/D6-FRONTEND.md) + [Frontend Method](development/frontend-product-experience-planning-method.md) |
| Current complete screen/surface inventory | [D6-R2 P5 Screen/Surface Inventory](engineering/rebaseline/D6-R2-P5-SCREEN-SURFACE-INVENTORY.md) |
| Current P8 operator-LOCK registry | [D6-R2 P8 Block Ledger](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md); follow only the block-specific evidence needed for the task |
| B10 Preparação accepted structure / current bounded reopen baseline | [B10 P8 Ratification](engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md), [bounded correspondence revalidation](engineering/rebaseline/D6-R2-P8-B10-CORRESPONDENCE-REVALIDATION.md), [B10 P9 Screen Contract](engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md), [`b10-preparation.html`](../qualification/d6-r2-wireframes/b10-preparation.html) |
| Notifications shell/Inbox/routing accepted structure | [Notifications P8 Ratification](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md), [final Notifications/Approvals P9 Contract](engineering/rebaseline/D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md), then the relevant locked HTML |
| Approvals / AuthorizationRequest accepted frontend structure | [B110 P8 Ratification](engineering/rebaseline/D6-R2-P8-B110-APPROVALS-RATIFICATION.md), [final Notifications/Approvals P9 Contract](engineering/rebaseline/D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md), [`b110-approvals.html`](../qualification/d6-r2-wireframes/b110-approvals.html) |
| Runtime/process/jobs/transactions | [D7 Runtime / Jobs / Transactions](engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md); use D7-B/C/D/E only for the relevant runtime concern |
| PostgreSQL isolation/RLS/transaction/revision mechanics | [D7-B PostgreSQL Isolation & Transactions](engineering/rebaseline/D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) |
| Durable work, retry/external-effect ambiguity and reconciliation | [D7-C Durable Work & External Effects](engineering/rebaseline/D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) |
| Human session/CSRF/OIDC and machine bearer realization | [D7-D Authentication / Session / CSRF](engineering/rebaseline/D7-D-AUTHENTICATION-SESSION-CSRF.md) |
| HTTP validation, storage, telemetry, deployment/backup/proof | [D7-E Operability / Deployment / Proof](engineering/rebaseline/D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) |
| Golden-flow business proof | [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md); [D8-R2 Operational Read Revalidation](engineering/rebaseline/D8-R2-OPERATIONAL-READ-REVALIDATION.md) only for GF-02 read revalidation |
| Stable cross-stage architecture constraints | [`../ARCHITECTURE.md`](../ARCHITECTURE.md) |
| ADR status / retirement condition | [ADR Registry](architecture/decisions/README.md), then only the retained ADR required by the question |
| Concrete evidence for a named decision | [Evidence Register](engineering/rebaseline/EVIDENCE-REGISTER.md), targeted proof, or Git history when required |
| Repository-local Git/CI/proof/document-lifecycle rules | [`development/engineering-rules.md`](development/engineering-rules.md) |
| Production coding / technology/dependency research | [`development/engineering-method.md`](development/engineering-method.md) plus the optional [Evidence-Grounded Production Engineering guide](development/evidence-grounded-production-engineering-for-llm-agents.md) when its detailed technology/proof lenses materially help |

## Navigation principle

The index answers **where current truth is owned**, not how every decision was historically produced. Findings, reviews, ratifications, plans and Git history are consulted when their evidence is materially needed; they are not a default recursive read set.

A route is an entry point, not a correctness fence. If current Evidence contradicts an accepted owner, use the applicable method and reopen the smallest owning authority instead of silently patching around it.
