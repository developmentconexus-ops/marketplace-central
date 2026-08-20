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
| Domains, semantic owners, allowed ownership edges | [D1 Domains/Boundaries](engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md) | implementation history |
| Organization, Principal, identity, data ownership, knowledge/provenance | [D2 Identity/Tenant/Data Ownership](engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md) | provider/runtime details |
| Q/C/E/P, communication, events, recovery, projections | [D3 Communication/Events](engineering/rebaseline/D3-COMMUNICATION-EVENTS.md) | D4 provider detail |
| Mercado Livre, Sankhya, provider boundaries/effects/reconciliation | [D4 External Integrations](engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md) | old adapters/runtime |
| Publication input, `ListingIntent`, readiness/source-following/override | [D4-R1 Publication Input](engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md) | generic PIM assumptions |
| Product API semantic laws | [D5 API](engineering/rebaseline/D5-API.md) | D6–D9 |
| Stable cross-stage architecture constraints | [`../ARCHITECTURE.md`](../ARCHITECTURE.md) | all phase documents |
| ADR disposition / retirement trigger | [ADR Registry](architecture/decisions/README.md), then only the named ADR | all ADRs |
| Historical/current-state Evidence for a concrete decision | [Evidence Register](engineering/rebaseline/EVIDENCE-REGISTER.md), then exact cited ref/path | repository archaeology |
| Repository-local engineering/Git/CI/proof rules | [`development/engineering-rules.md`](development/engineering-rules.md) | organizational Method text copied locally |

## Product OAD task pack

For canonical Product OAD authoring, start with this bounded pack:

1. [D5 API](engineering/rebaseline/D5-API.md)
2. [Product Operation Surface](engineering/rebaseline/D5-B2-PRODUCT-OPERATION-SURFACE.md) + [Operation Admission Matrix](engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md)
3. [W1 Wire Contract](engineering/rebaseline/D5-B2-WIRE-CONTRACT.md)
4. [W2 Schema Grammar](engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md) + [W3 Collection Grammar](engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md)
5. [W4 Permission / Client-Class Enforcement](engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md) + [Technical Ingress](engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md) + [OpenAPI Wire Authority / Tooling](engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md)

Read D0–D4 only when the D5 package cannot answer a concrete contract question without reopening earlier semantics.

## Organizational standards

Canonical, external to this repository:

- `developmentconexus-ops/conexus-methodology/METHOD.md` — **DevelopmentConexus Engineering Method v1.0.0**
- `developmentconexus-ops/conexus-methodology/REPOSITORY-STANDARD.md` — **DevelopmentConexus Repository Standard v1.0.0**

No local copy is Product or repository authority.
