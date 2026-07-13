# F-08 Assisted Sankhya Linkage Discovery

```yaml
id: F-08
type: feature
status: planned
owner: Feature Implementer
parent: M-06
```

## Outcome

Produce the bounded, read-only discovery and deployable administrator contract
needed for MPC to assist an operator linking an imported Mercado Livre order to
a manually created Sankhya TOP 313 order and to trace exact lines through
`TGFVAR` to TOP 306.

## Brief

Produce a bounded read-only discovery record and fail-closed administrator and
implementation contract for operator-assisted Mercado Livre order linkage to
exact Sankhya TOP 313 header/lines and their TOP 306 `TGFVAR` descendants.

## Scope

- Reconcile current repository Oracle metadata/custom-field conventions and
  the approved live-Oracle read lane.
- When safely available, run only bounded SELECT discovery for supported custom
  header metadata and TOP 313 -> 306 `TGFVAR` behavior. Never output secrets,
  buyer PII, customer names, or unbounded row data.
- Name an exact deployed field only when repository or runtime evidence proves
  it. Otherwise mark the name unknown and produce a deployable admin
  specification for an explicit configured field that MPC validates and uses
  fail-closed.
- Decide whether a Sankhya line custom field is necessary for the current
  assisted workflow or whether an explicit audited MPC mapping to exact
  (`NUNOTA`,`SEQUENCIA`) is sufficient.
- Define the smallest implementation contract for one canonical linkage
  service and append-only tenant/account-scoped ledger. Candidate filters are
  never linkage proof.

## Acceptance criteria

1. `discovery.md` distinguishes observed repository/runtime facts, inferences,
   and unknowns, including exact TOP and `TGFVAR` predicates where proven.
2. `sankhya-admin-spec.md` is deployable without inventing a live field name;
   it defines field type/length, uniqueness/account scope, configuration,
   validation, permissions, and rollback/disable behavior.
3. The line-field decision is explicit and justified for the assisted-only
   flow, including partial invoice and one-to-many `TGFVAR` descendants.
4. `implementation-contract.md` fixes domain/ports/adapters boundaries,
   candidate versus proof semantics, idempotency, append-only audit, explicit
   unknowns, and exact API/SDK/migration seams for the next worker.
5. `validation.md` records exact commands, target labels, bounded observables,
   and limitations. No Oracle/provider write, DDL, secret, or PII is produced.

## Expected Output

- `discovery.md` classifies repository facts, bounded runtime facts,
  inferences, and unknowns.
- `sankhya-admin-spec.md` defines a deployable fail-closed configuration
  without inventing a live field name.
- `implementation-contract.md` fixes the next-worker boundaries, seams,
  idempotency, audit, and unknown states.
- `validation.md` records exact bounded evidence and limitations.
