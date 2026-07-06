# Execution Guide

```yaml
id: M-02
type: execution-guide
status: in_progress
owner: Milestone Orchestrator
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: support
```

## Scope

Mission/Milestone ID: MIS-001 / M-02 Marketplace capability framework.

M-02 creates the provider capability layer that lets future `product_links`, `inventory`, `orders`, and `profitability` modules depend on normalized marketplace operations instead of Mercado Livre endpoint payloads.

## How To Resume

Read these files first:

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-01-vtex-removal-architecture-reset/validation-result.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/milestone.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/validation-contract.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/research/marketplace-capability-interface-contract.md`
- the active feature brief under `M-02-marketplace-capability-framework/F-*/feature.md`

## Execution Order

1. F-01 Capability port contract: define small Go ports and normalized DTOs for listing, stock, and order capability use.
2. F-02 Provider capability registration: align provider definitions and capability health with the new business capability names.
3. F-03 Mercado Livre adapter spine: map documented listing, variation stock, order, and guarded stock-write shapes into normalized capability outputs.
4. Milestone gate: after all three features are accepted, run the independent M-02 validation gate and persist `validation-result.md`.

F-01 must run before F-02/F-03. F-02 and F-03 may overlap only after F-01 returns accepted contracts, because they share names, errors, and DTOs.

## Context Package For Next Session

- Required files: F-01 `feature.md`, M-02 `milestone.md`, M-02 `validation-contract.md`, mission `validation-contract.md`, and IC-001.
- Required contracts: M-02-C01 and M-02-C02; IC-001 operations, required fields, enums/statuses, error matrix, and compatibility rules.
- Current status: M-01 passed; M-02 in progress; F-01 is briefed and ready for Feature Implementer.
- Next action: execute F-01 and return `spec.md`, `plan.md`, changed paths, and `validation.md`.

## F-01 Dispatch Packet

Feature: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-01-capability-port-contract/feature.md`

Build the smallest durable capability contract that satisfies IC-001:

- provider ids stay strings;
- provider endpoint payload structs stay outside business domain/application code;
- unsupported capability, missing installation, credential unavailable, rate limit, validation rejection, unsupported shape, transient failure, and invalid payload map to structured errors from IC-001;
- tests prove a business-facing service can use fake capability implementations.

Do not implement Stock Seguro policy, product links, order profitability, UI, or real Mercado Livre HTTP calls in F-01.

## Handoff

- Current status: in_progress.
- Next owner: Feature Implementer.
- Next action: execute F-01.
- Required files/evidence: F-01 `spec.md`, `plan.md`, changed paths, and `validation.md`; later M-02 `validation-result.md`.
- Blockers or open decisions: none for F-01.

## Advancement Rules

- A feature can be accepted only after its `spec.md`, `plan.md`, changed paths, and `validation.md` are inspectable.
- Acceptance evidence must cite commands that ran; `assumed` or `could-not-run` evidence cannot satisfy load-bearing criteria.
- F-03 provider behavior claims must cite official Mercado Livre documentation or local adapter tests using documented shapes.
- API changes require OpenAPI and `packages/sdk-runtime` updates together.
- After F-01, run import-boundary checks before dispatching adapter work.

## Failure Rules

- Missing F-01 execution artifacts blocks feature acceptance.
- Any provider payload crossing into business domain/application code rejects the feature output.
- Unsupported capability returning nil success, panic, or ambiguous behavior rejects the feature output.
- Live provider stock writes remain blocked unless the human owner provides controlled credentials/listings and explicit approval.

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none
- correction scope owner: Milestone Orchestrator
- blocked-report trigger: unresolved dependency, missing feature evidence, failed milestone gate after retry exhaustion, or required human decision.
- evidence required before revalidation: accepted feature validations, changed paths, command evidence, and M-02 gate output.

## Notes For Future Automation

- Generate a fresh-session F-01 context pack directly from this guide, F-01 `feature.md`, M-02 `validation-contract.md`, and IC-001.
- Keep M-02 validation result derived from the independent milestone gate, not from this orchestration guide.
