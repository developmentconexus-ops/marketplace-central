# Marketplace Central — Documentation Index

<!-- program-status-authority -->

> **Role:** sole current-program status, allowed/blocked-work and exact-next-action authority. This file routes readers to detailed semantic authority; it does not duplicate that authority.

## Default bootstrap

A fresh session reads only:

1. [`AGENTS.md`](../AGENTS.md)
2. this file

Then stop and select the smallest additional read set from the tables below. Do not recursively read every D-stage, ADR, evidence file, code tree or Git history.

Stage-local status and read-order prose accepted before this checkpoint is superseded by this index; it remains historical context, not an active route or second status authority.

## Current checkpoint

| Field | Current value |
|---|---|
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D5 — API — OPEN / ACTIVE** |
| Accepted through | D0–D4, D4-R1, D5-B1, and D5-B2 W1/W2/W3/W4 + Technical Ingress + final Problem/media + OpenAPI authority/tooling |
| Exact next action | **Author and prove the canonical Product OpenAPI Description** |
| Entry document to create | `contracts/api/product/openapi.yaml` |
| Product surface | **95 operations · 29 ordinary Permissions · Principal kinds H/A/S only** |
| Stable Problem origin | `https://conexus.fun` |
| Product Problem namespace | `https://conexus.fun/marketplace-central/problems/product/{slug}` |
| Technical Problem namespace | `https://conexus.fun/marketplace-central/problems/technical/{surface}/{slug}` when a technical contract actually needs it |
| Ngrok | preview/tunnel only; forbidden in canonical `Problem.type` values |
| Implementation | **BLOCKED until D9 is accepted** |

The 95 admitted Product operations define the Product 1.0 contract surface. They do **not** require 95 runtime handlers before the first useful internal vertical slice.

## Program stages

| Stage | Status | Canonical home | Read when the task asks… |
|---|---|---|---|
| D0 — Product / System Definition | ACCEPTED / CLOSED | [D0](engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md) | what the product is, its actors, scope or Product 1.0 boundary |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED | [D1](engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md) | who owns a business meaning/state/capability or which semantic edge is valid |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED | [D2](engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md) | Organization, Principal, IDs, source qualification, persistence ownership or knowledge/provenance semantics |
| D3 — Communication / Events | ACCEPTED / CLOSED | [D3](engineering/rebaseline/D3-COMMUNICATION-EVENTS.md) | Q/C/E/P, propagation, idempotency, recovery, projections or cross-owner coordination |
| D4 — External Integrations | ACCEPTED / CLOSED | [D4](engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md) | Mercado Livre, Sankhya, adapters, provider capability/coverage, external effects or reconciliation |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL | [D4-R1](engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md) | Product→listing authoring, `ListingIntent`, readiness, source following or explicit override |
| D5 — API | OPEN / ACTIVE | [D5](engineering/rebaseline/D5-API.md) | semantic Product API laws and current API-stage ownership |
| D6 — Frontend | BLOCKED BY D5 | not opened | screens, navigation, editor/query composition or frontend package topology |
| D7 — Runtime / Jobs / Transactions | BLOCKED | not opened | processes, router, transactions, queues, workers, schedulers, retries, storage or deployment |
| D8 — Golden Flows | BLOCKED | not opened | end-to-end operational proof and real external-effect lanes |
| D9 — Adversarial Architecture Review | BLOCKED | not opened | final whole-system contradiction/overbuild/under-specification review |
| Implementation | BLOCKED UNTIL D9 | future implementation DAG | code delivery after accepted D9 |

## Selective read routes

### Current D5 Product OAD work

For the exact next action, read this bounded package rather than the entire rebaseline:

| Question | Read |
|---|---|
| Parent semantic API laws | [D5 API](engineering/rebaseline/D5-API.md) |
| Admitted Product/resource surface | [Product Operation Surface](engineering/rebaseline/D5-B2-PRODUCT-OPERATION-SURFACE.md) and [Operation Admission Matrix](engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md) |
| Resource names, paths and HTTP grammar | [W1 Wire Contract](engineering/rebaseline/D5-B2-WIRE-CONTRACT.md) |
| Requests, responses, Problems, idempotency, revision and media grammar | [W2 Schema Grammar](engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md) |
| Collections, filters and cursor grammar | [W3 Collection Grammar](engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md) |
| Exact 95-operation naming, Permissions and Principal admission | [W4 Permission / Client-Class Enforcement](engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md) |
| Provider/OAuth acquisition outside the Product API | [Technical Ingress](engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md) |
| One OAD authority, tool pins, generation and retirement proof | [OpenAPI Wire Authority / Tooling](engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md) |

Read D0–D4 only when a concrete OAD question cannot be answered by the D5 package without reopening earlier semantics.

### Other questions

| Task concerns | Smallest starting read |
|---|---|
| Stable cross-stage platform constraint | [`ARCHITECTURE.md`](../ARCHITECTURE.md) |
| Which decision generation is current after iterative rebaseline work | [Decision Reconciliation Baseline](engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md) — on demand, never default-read |
| ADR file status or unresolved legacy residue | [ADR Registry](architecture/decisions/README.md), then only the named ADR |
| Engineering reasoning for material work | [DevelopmentConexus Engineering Method](engineering/standards/root-cause-global-maximum-method.md) — non-authoritative local availability copy of `developmentconexus-ops/conexus-methodology/METHOD.md` v1.0.0; replace manually only after operator-approved adoption, without sync machinery |
| Historical/current-state facts needed for a specific decision | [Evidence Register](engineering/rebaseline/EVIDENCE-REGISTER.md), then only the cited primary evidence/code |
| Observed legacy defect/proof classes | [Defect Class Catalog](engineering/defect-class-catalog.md) — supporting reference only |
| Current read-only Oracle validation procedure | [Live Oracle Docker](operations/live-oracle-docker.md) — current-state procedure, not target architecture |
| Machine-selectable knowledge routes | [`knowledge-routes.json`](../contracts/governance/knowledge-routes.json) — tooling aid, not status authority |
| Current runtime behavior | relevant code, schema, OpenAPI and tests after the semantic authority above |

## Exact next action: canonical Product OAD authoring/proof

The current architectural sub-batch may:

1. create `contracts/api/product/openapi.yaml` and repository-local `paths/` / `components/` source files;
2. encode exactly the accepted 95 Product operations and 29-Permission projection without adding business meaning;
3. use exact Product Problem constants under `https://conexus.fun/marketplace-central/problems/product/{slug}`;
4. add the accepted rules-only Redocly configuration and deterministic lint/bundle/generation proof;
5. generate and compile the derived TypeScript and Go contract projections;
6. retire or rehome the measured legacy OpenAPI/manual-SDK authority seam atomically;
7. execute negative controls and repository gates;
8. stop and reopen only the smallest authority if executable proof discovers a material contradiction.

It may **not**:

- begin D6–D9 or product implementation;
- implement all 95 runtime handlers merely because the OAD contains them;
- choose a Go router/server framework, runtime validator, deployment topology, queue, storage, CDN or cache;
- add a 96th Product operation, 30th Permission, fourth Principal kind, standalone Product-media CRUD or `/v1`;
- place provider/OAuth acquisition, provider callbacks/webhooks or authored-media byte delivery in the Product OAD/SDK;
- convert MPC Permissions into OAuth scopes;
- use ngrok or another temporary host in canonical Product Problem identifiers;
- preserve the legacy OpenAPI/manual Product SDK as a parallel target authority.

## Authority by scope

Use this order only for the scope actually in question:

1. operator-ratified decision recorded in accepted authority;
2. stable constraint in [`ARCHITECTURE.md`](../ARCHITECTURE.md);
3. accepted D-stage artifact for the relevant semantic scope;
4. accepted ADR not explicitly reopened, with status from the [ADR Registry](architecture/decisions/README.md);
5. current code, schemas, APIs, tests and runtime as current-state evidence.

`docs/README.md` owns status/routing, not detailed product semantics. A material conflict must be surfaced and adjudicated; never choose silently because one file is newer or easier to implement.

## Delivery discipline after D9

Target architecture is not first-release implementation scope. Future implementation planning classifies material work as:

```text
BUILD NOW   required by the current consumer/golden flow or correctness
SEAM NOW    preserve ownership/boundary now; future capability later
PROVE FIRST smallest bounded spike before committing mechanism/dependency
DEFER       no current consumer/failure class; record the reopen trigger
```

Preserve correctness, Organization isolation, ownership, source-qualified identity, security boundaries, migrations, secrets, consequential idempotency/reconciliation, audit/recovery evidence and adapter boundaries when their later retrofit would be unsafe. Defer speculative scale/platform machinery. Build proven end-to-end vertical slices after D9 rather than every repository/service/handler layer horizontally.

## D-stage and pull-request lifecycle

After the current exceptional checkpoint reaches `main`:

```text
main
→ one branch / one PR for the active D-stage
→ lead analysis and operator discussion
→ consolidated candidate
→ one final material Fable review when warranted
→ GPT adjudication; Round 2 only for surviving material contradiction
→ operator ratification
→ absorb durable outcomes into canonical authority
→ delete temporary review/candidate/handoff files
→ gates + cold scope review
→ explicit operator-authorized merge
→ next D-stage starts from clean main
```

Do not stack later D-stages on an unmerged earlier stage.

## Documentation lifecycle

Every active document must have one responsibility and one of these dispositions:

| Disposition | Rule |
|---|---|
| CANONICAL | owns durable meaning in a named scope and is indexed here |
| SUPPORTING | evidence/reference read only for a concrete need; cannot override canonical authority |
| TEMPORARY | exists only on the active branch/PR; must be absorbed or deleted before merge |
| HISTORY | recoverable through Git; not mirrored into `old/`, archive, handoff or roadmap trees |

| Document | Disposition | Responsibility |
|---|---|---|
| `README.md` | SUPPORTING | GitHub landing pointer only; no independent program status or next action |

Binding rules:

- Current status and exact next action live only in this file.
- `AGENTS.md` is an operating bootstrap, not a roadmap or product specification.
- Do not create `docs/superpowers/`, parallel roadmaps, permanent session handoffs, active archive trees or a second status dashboard/cockpit.
- `AI-DIALOG.md` is created only for an active bounded Fable review cycle and is deleted before merge; Git history preserves the cycle.
- External research is absorbed into the responsible D-stage decision or retained as clearly time-bounded supporting evidence.
- Superseded candidates are removed after their surviving meaning is filed.
- `engineering/rebaseline/README.md` is a temporary compatibility pointer for old links only; it owns no status and is not part of the read order.

## Machine/runtime authorities

The following describe current mechanics and are not automatically target architecture:

- `contracts/api/marketplace-central.openapi.yaml` — legacy/current HTTP contract until the canonical Product OAD replacement lands;
- `contracts/governance/` — current machine-enforced repository governance;
- `contracts/gate/` — current verification ratchets;
- `scripts/gate.ps1` — shared local/CI gate implementation;
- application code and migrations — evidence of what currently runs.

## Fresh-session success test

Reading only [`AGENTS.md`](../AGENTS.md) and this file must be enough to conclude:

- D0–D4 and D5-B1 are accepted;
- D5-B2 is active and canonical Product OAD authoring/proof is next;
- the Product surface is 95 operations / 29 Permissions / H-A-S Principals;
- `https://conexus.fun` is the stable origin and ngrok is preview-only;
- implementation remains blocked until D9;
- detailed documents are selected by question rather than read recursively;
- target architecture does not require building every future capability in the first release;
- Fable review and `AI-DIALOG.md` are bounded/temporary;
- a future D-stage lands through one PR before the next stage begins.

If the two-file bootstrap cannot produce those conclusions, the active authority surface is defective.
