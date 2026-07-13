# F-08 Assisted Sankhya Linkage Discovery Spec

```yaml
id: F-08
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-08
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-1
lifecycle_scope: feature
```

## Feature ID

F-08-assisted-sankhya-linkage-discovery

## Problem

MPC cannot yet prove which exact Sankhya TOP 313 document and lines correspond
to an imported Mercado Livre order. Consequently it cannot follow exact line
lineage to TOP 306 or supply order-specific tax provenance. Customer-specific
custom-field names and live TOP/TGFVAR behavior must not be invented when
runtime evidence is unavailable.

## Requirements

- Separate observed repository facts, safely observed runtime facts,
  inferences, and operational unknowns. Acceptance evidence: `discovery.md`
  names the evidence class and limitations for each conclusion.
- Define a deployable, fail-closed Sankhya administrator contract without
  assuming an existing field name. Acceptance evidence:
  `sankhya-admin-spec.md` fixes the field shape, account-scoped uniqueness,
  configuration and validation gates, permissions, and disable/rollback
  behavior.
- Decide whether the assisted-only workflow needs a Sankhya line custom field,
  including exact handling of partial invoices and one-to-many `TGFVAR`
  descendants. Acceptance evidence: the decision and rationale appear in both
  the discovery and administrator specification.
- Define the smallest next-worker contract for one canonical linkage service
  and append-only tenant/account-scoped ledger. Acceptance evidence:
  `implementation-contract.md` fixes boundaries, candidate-versus-proof
  semantics, idempotency, audit, explicit unknowns, and named contract/code/
  migration seams without modifying them.
- Record exact commands, safe target labels, bounded observables, and
  limitations. Acceptance evidence: `validation.md` contains no secret, PII,
  DDL, provider write, Oracle write, or application-data write.

## Non-Goals

- No application, OpenAPI, SDK, migration, Oracle schema, provider, or
  application-data change.
- No creation or alteration of Sankhya custom fields, indexes, copy rules,
  TOPs, or events.
- No automatic matching by partner, buyer, date, product, quantity, price, or
  candidate confidence.
- No claim that an unobserved custom field, uniqueness capability, copy rule,
  TOP behavior, or live lineage is deployed.
- No milestone acceptance or QA verdict.

## Design

The output is a documentation-only contract. Repository and any safely
available bounded SELECT evidence are recorded separately. The administrator
configures explicit header metadata for the immutable installation/account-
scoped provider order key. For the current assisted-only workflow, MPC owns
the exact line selection in an append-only ledger; a Sankhya item custom field
is optional hardening, not linkage proof. The operator must explicitly select
the TOP 313 header and each exact (`NUNOTA`, `SEQUENCIA`) line. MPC validates
configured metadata and uniqueness fail-closed, records actor/reason/source
time, and follows `TGFVAR` origin-to-destination rows for exact 306 descendants.

Candidate queries only narrow the operator's choices. Proof exists only after
explicit confirmation plus successful validation of immutable external key,
tenant, installation/account, exact 313 identities, and non-conflicting ledger
state. Unknown or ambiguous cases remain unlinked and tax remains missing.

## Edge Cases

- No safe live Oracle access or command: runtime field names and lineage facts
  remain `unknown`; the admin contract remains deployable.
- Duplicate, reused, blank, truncated, or account-mismatched external key:
  reject linkage and make no state transition.
- Identical ML lines or mutable positional order: each linkage uses an
  immutable MPC line identity; no attribute tuple becomes identity.
- Split or combined 313 candidates: require explicit selection of every exact
  origin line; unresolved lines remain unlinked.
- Partial invoicing and one-to-many descendants: retain the exact 313 line as
  origin and append each exact 306 (`NUNOTA`, `SEQUENCIA`) descendant with
  `QTDATENDIDA`; never collapse lineage into product/date totals.
- Zero, missing, reversed, or conflicting `TGFVAR` rows: surface an explicit
  lineage state and block complete tax provenance.
- Retry or concurrent confirmation: idempotency returns the existing matching
  mapping; a conflicting mapping fails closed and is audited, never overwritten.
- Disable/rollback: disable new confirmations without deleting ledger history
  or clearing ERP metadata.

## Acceptance Criteria

### F08-AC01

- Criterion: Discovery separates repository facts, runtime facts,
  inferences, and unknowns and names exact TOP/`TGFVAR` predicates only when
  evidenced.
  - Traces to milestone criterion ID: M-06-C02 — Margin Quality Honesty.
  - Proven by: `git diff --check` plus manual inspection of `discovery.md` and
    the bounded live-discovery evidence classification in `validation.md`.
### F08-AC02

- Criterion: The administrator contract is deployable without an
  invented field name and includes field shape, scope, uniqueness,
  configuration validation, permissions, and disable/rollback behavior.
  - Traces to milestone criterion ID: M-06-C02 — Margin Quality Honesty.
  - Proven by: manual fail-closed contract inspection recorded in
    `validation.md`.
### F08-AC03

- Criterion: The line-field decision covers assisted-only exact
  mapping, partial invoicing, and one-to-many descendants without promoting a
  candidate to proof.
  - Traces to milestone criterion ID: M-06-C02 — Margin Quality Honesty.
  - Proven by: manual cross-artifact lineage inspection recorded in
    `validation.md`.
### F08-AC04

- Criterion: The implementation contract fixes one canonical
  service, tenant/account-scoped idempotency, append-only audit, explicit
  unknowns, and paired API/SDK/migration seams for the next worker.
  - Traces to milestone criterion ID: M-06-C01 — Idempotent Order Ingestion;
    M-06-C02 — Margin Quality Honesty; M-06-C03 — Manual Adjustment Audit.
  - Proven by: manual boundary/traceability inspection and scoped Git checks
    recorded in `validation.md`.
### F08-AC05

- Criterion: Evidence names exact commands, target labels, bounded
  observables, and limitations and contains no DDL, write, secret, or PII
  output.
  - Traces to milestone criterion ID: M-06-C02 — Margin Quality Honesty.
  - Proven by: registered-command resolution, bounded execution where safe,
    artifact content scans, and `git diff --check` recorded in `validation.md`.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Execute `plan.md` and record bounded discovery evidence.
- Required files/evidence: feature brief, validation contract, architecture
  research, approved execution lanes, spec
- Blockers or open decisions: Live customer field names, supported uniqueness
  mechanism, copy behavior, and observed TOP lineage are operational unknowns
  unless the approved bounded discovery safely proves them.
