# M-02-marketplace-capability-framework

```yaml
id: M-02
type: milestone
status: passed
owner: QA Validator
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Outcome

MPC exposes marketplace capability ports that business modules can use without provider-specific coupling, and Mercado Livre has the first adapter spine for listing, stock, and order capabilities.

## Why This Milestone Exists

Mercado Livre is first, but the platform must remain a marketplace hub. Capability ports prevent rules from being trapped inside Mercado Livre code.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Capability port contract | Add small Go ports for listing, stock, order, shipment/question placeholders only where needed. |
| F-02 | Provider capability registration | Align provider definitions/capability health with new business capability names. |
| F-03 | Mercado Livre adapter spine | Implement documented shape mapping for listing/stock/order reads and a guarded stock write interface using direct HTTP seams/stubs. |

## Dependencies

- M-01 removes VTEX contradictions.
- IC-001 marketplace capability interface contract.

## Risks

- Over-abstracting provider differences.
- Under-abstracting and locking business modules to Mercado Livre.

## Done Means

- Business modules can depend on capability ports.
- Mercado Livre adapter maps documented provider fields to normalized snapshots.
- Tests prove unsupported capabilities are explicit.

## Handoff

- Current status: passed.
- Next owner: Mission Strategist.
- Next action: Continue mission execution using M-02 capability contracts, provider declarations, and Mercado Livre adapter spine as execution truth.
- Required files/evidence: `validation-result.md` plus accepted F-01/F-02/F-03 artifacts.
- Blockers or open decisions: Live Mercado Livre runtime validation still requires operator-controlled credentials/listings and explicit approval; M-02 pass covers adapter-spine readiness, not live-provider proof.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: Not applicable.
- Attempts used/remaining: 0/2.
- Next artifact: M-02/validation-result.md.
- Revalidation evidence required: Go tests and import-boundary checks.
