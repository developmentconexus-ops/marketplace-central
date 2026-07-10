# M-06 F-01 Correction Brief

## Context

F-01 order ingestion exists but was rejected during milestone acceptance because the evidence and tests do not prove the complete timestamp freshness policy or buyer-PII minimization required by its specification.

## Required Behavior

1. Persist an absent order snapshot.
2. Replace an existing snapshot when the incoming `provider_updated_at` is newer.
3. Replace it when timestamps are equal, without creating duplicate order, item, or payment rows.
4. Reject an older incoming snapshot.
5. Reject an incoming snapshot with unknown `provider_updated_at` when the stored snapshot has a known timestamp, so known freshness is never erased.
6. Allow replacement when the stored timestamp is unknown, including when both timestamps are unknown.
7. Keep raw provider data limited to safe references; buyer PII must not enter the orders domain, SQL schema, OpenAPI order schemas, SDK order types, or API readback.

## TDD Contract

- Add focused failing tests before production changes and record the RED failure.
- Implement the smallest domain/application/repository policy needed to pass.
- Run targeted orders tests and relevant regression tests.
- Do not use a fake-only result as proof of Postgres behavior; identify the exact real-DB command/readback needed for later gate execution.

## Scope

- `apps/server_core/internal/modules/orders/**`
- F-01 `validation.md` only if fresh evidence is actually run and can be cited.
- `.superpowers/sdd/m06-f01-correction-report.md`

Do not edit other features, OpenAPI, SDK, composition, roadmap, or unrelated dirty files. Do not revert other contributors' work. Do not commit in the shared dirty worktree; report changed paths for controller review.

## Report Contract

Write `.superpowers/sdd/m06-f01-correction-report.md` with:

- status: DONE, DONE_WITH_CONCERNS, NEEDS_CONTEXT, or BLOCKED
- root cause
- RED command and expected failure excerpt
- GREEN command and result
- changed paths
- real-DB verification still required
- self-review concerns
