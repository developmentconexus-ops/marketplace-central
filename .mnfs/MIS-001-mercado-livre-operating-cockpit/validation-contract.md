# Mission Validation Contract — MVP Replan

```yaml
id: MIS-001
type: mission-validation-contract
status: ready
owner: Mission Strategist
parent: none
created: 2026-07-06
updated: 2026-07-13
validation_level: QA-4
lifecycle_scope: mission
```

## Required Final State

Marketplace Central demonstrates one coherent trusted-local operator journey backed
by PostgreSQL, real read-only Mercado Livre and Sankhya data, canonical CODPROD
identity, and a visibly simulated stock correction. No provider write, production
authentication, or post-MVP commercial breadth is required.

Mission QA evaluates the accepted results of M-09, M-13, and M-14. It does not rerun
historical passed milestones solely to prove preservation and never rewrites M-06's
failed production-gate verdict.

## Required Criteria

### MIS-001-C01 — Architecture and project conventions

- Evidence: accepted milestone reviews, targeted/broad Go tests, SDK/web tests and
  build, and source-boundary inspection.
- Expected: domain/application/ports/adapters/transport boundaries remain intact;
  tenant queries scope `tenant_id`; web uses SDK runtime; OpenAPI/SDK move together.
- Blocks: provider/transport leakage into business packages or direct production
  browser bypass around SDK runtime.

### MIS-001-C02 — Canonical and honest product facts

- Evidence: passed M-09 criteria and its real Oracle read.
- Expected: positive CODPROD is the internal identity; provider/EAN/reference values
  remain separate; missing numeric facts are null/unknown, never zero defaults.
- Blocks: guessed mapping, MSDB still required by the MVP catalog, or unknown=zero.

### MIS-001-C03 — Durable duplicate-safe application state

- Evidence: passed M-14 PostgreSQL integration/idempotency proof.
- Expected: repeated bounded listing/order import preserves one natural identity and
  browser reload retains MPC-owned links, snapshots, and journey context.
- Blocks: duplicate durable identities or fixture-only persistence claims.

### MIS-001-C04 — Safe integrated operator journey

- Evidence: passed M-13 workspace/browser criteria and M-14 vertical drive.
- Expected: Visão geral, Produtos, Anúncios, Vendas, Operações and deep links connect
  attention, listing, Product 360, sale/margin, and stock simulation.
- Blocks: lost identity/context, hidden load-bearing error/conflict, or client-side
  stock/margin calculation.

### MIS-001-C05 — Real-read evidence honesty

- Evidence: M-14 sanitized Oracle, Mercado Livre, PostgreSQL, and browser artifacts.
- Expected: every real claim names source/time/read-only status and every fixture is
  labeled deterministic; mocks never prove integration.
- Blocks: fake-as-real evidence or missing source provenance.

### MIS-001-C06 — Secret, PII, and provider-write safety

- Evidence: M-13 simulation proof and M-14 evidence/network safety fold.
- Expected: no secret or buyer PII in UI/logs/evidence; no provider mutation is
  reachable or reported; simulation says `executed=false`.
- Blocks: credential/PII disclosure, provider write, or misleading success copy.

### MIS-001-C07 — Historical roadmap integrity

- Evidence: mission roadmap and preserved milestone artifacts.
- Expected: M-01–M-05/M-08 remain historical passed foundations; M-06 remains a
  failed production gate whose order/margin implementation may be reused; M-07 waits
  until after M-14; M-10/M-11/M-12 remain post-MVP.
- Blocks: rewriting M-06 as passed or reopening its C03 auth criterion for this MVP.

## Proportional Evidence Rules

- Feature evidence: `<feature-root>/validation.md`.
- Milestone QA: `<milestone-root>/validation-result.md`.
- Mission QA: mission `validation-result.md` after M-14.
- Commands name cwd, reviewed SHA, exit code, and honest limitations.
- API changes include OpenAPI and SDK evidence in the same commit.
- Real evidence contains no credentials, buyer PII, raw provider payloads, or full
  Oracle rows.
- Browser evidence records the primary journey, a representative negative state,
  reload/persistence, and desktop/mobile screenshots.
- Atomic wrappers, hash manifests, OCR automation, and external-library claim ledgers
  are optional tooling, not acceptance criteria.

## Blocking Failures

- Unknown business numeric becomes zero.
- A provider or Oracle write occurs or is required for Pass.
- A secret/PII value enters response, log, screenshot, or evidence.
- A mock is reported as real integration.
- The browser loses canonical product/listing/sale identity.
- M-06's failed historical result is relabeled or erased.

## Retry Policy

Maximum two scoped correction attempts per new milestone. Historical M-06 attempts
are not reset. A missing real source/sample becomes `externally_blocked` and follows
the harness terminal callback.

## Handoff

- Current status: Ready.
- Next owner: M-09 Milestone Orchestrator.
- Next action: execute M-09, then M-13, then M-14 through separate user-started
  Milestone tasks and QA-owned verdicts.
- Open decisions: none for M-09 dispatch.
