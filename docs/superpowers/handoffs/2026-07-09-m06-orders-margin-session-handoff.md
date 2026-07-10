# M-06 Orders And Margin Session Handoff

Date: 2026-07-09
Branch: `main`
Status: Implementation in progress; F-01/F-02 corrections accepted, F-03 Task 1 stopped at RED
Owner plan: `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`
Mission: `.mnfs/MIS-001-mercado-livre-operating-cockpit`

## Goal

Continue the full marketplace operating-cockpit plan using `superpowers:subagent-driven-development`, optimizing for the durable global maximum rather than existing legacy abstractions. Completion requires code review plus real PostgreSQL, Mercado Livre, Oracle, and built-in-browser QA; mocks and tests alone cannot produce a pass.

## Current State

- M-04 product-link flow and M-05 stock-safe flow were implemented before this handoff.
- M-06 F-01 correction is complete: atomic timestamp upsert, guarded child replacement, buyer-PII minimization, real PostgreSQL 16 evidence, and independent SPEC/QUALITY approval.
- M-06 F-02 correction is complete: required actor audit, scope/category invariants, cryptographic event IDs, append-only freight/commission history, forward-safe PostgreSQL checks, SDK 35/35, and independent SPEC/QUALITY approval.
- M-06 F-03 design and implementation plan are committed:
  - `c29315c docs(profitability): define realization states`
  - `fa88d3e docs(profitability): plan realization states`
- F-03 Task 1 is partial and intentionally stopped at RED. No production implementation or GREEN exists yet.

## F-03 Task 1 Partial Work

Changed by the stopped worker:

- `apps/server_core/internal/modules/profitability/adapters/orders/order_reader_test.go` (new)
- `apps/server_core/internal/modules/profitability/application/service_test.go` (modified)

Observed RED command:

```powershell
cd C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/adapters/orders ./internal/modules/profitability/application -count=1
```

Observed failure was expected: undefined profitability-owned `OrderFact`, `OrderItemFact`, realization-state constants, and link-quality constants. The next session must rerun RED and preserve the output before production edits because no task report was written.

## Immediate Objective

Finish only F-03 Task 1 from the existing RED:

1. Create `apps/server_core/internal/modules/profitability/domain/order_fact.go` with the exact contract from `.superpowers/sdd/m06-f03-task-1-brief.md`.
2. Change `profitability/ports.OrderReader` to return profitability-owned facts.
3. Translate orders application/domain values only inside `profitability/adapters/orders`.
4. Update profitability application signatures and tests without implementing snapshot realization yet.
5. Run GREEN, `gofmt`, and the boundary `rg` from the brief.
6. Write `.superpowers/sdd/m06-f03-task-1-report.md`.
7. Dispatch an independent task reviewer; do not begin Task 2 until SPEC and QUALITY are approved.

## Remaining F-03 Order

After Task 1 review is clean:

1. Task 2: realization-aware calculation (`realized`, `not_realized`, `unknown`); cancelled/unknown contribution and margin remain `nil`.
2. Task 3: migration `0031`, PostgreSQL repository, OpenAPI and SDK realization contract.
3. Task 4: orders UI semantics and frontend regression.
4. Full real gate: PostgreSQL 16, live Mercado Livre paid/cancelled orders, one resolved link with real Oracle `CUSSEMICM` and taxes, built-in browser desktop/mobile, then the full independent M-06 gate.

## Guardrails

- Read `C:\Users\leandro.theodoro\.codex\attachments\2a1a5352-6aee-406c-bf04-2d066266687f\goal-objective.md` before continuing.
- Read `AGENTS.md`, `ARCHITECTURE.md`, `wiki/README.md`, `.brain/system-pulse.md`, and `.brain/roadmap.json` at session start.
- User explicitly approved working on `main`; do not create a worktree unless the user changes that decision.
- The worktree is heavily dirty with prior M-03/M-04/M-05/M-06 work. Never revert, clean, reset, or overwrite unrelated changes.
- Existing/legacy implementation is not architectural truth. Follow current spec/design/plan and module boundaries.
- Only adapters may import another module's application/domain types. Ports and application code own and consume module-local contracts.
- Unknown cost, tax, fee, freight, commission, link, or realization state stays explicit and never becomes zero or realized by default.
- API changes require OpenAPI and `packages/sdk-runtime` updates together.
- Fake/mock tests prove deterministic behavior only. Real provider/database/runtime claims require real evidence and must be labelled accurately.
- Use one implementation subagent at a time, then an independent reviewer. Close agents after use.
- Preserve TDD RED before production changes and run fresh GREEN before claiming success.
- Do not use or install WSL. The Superpowers `task-brief` helper is Bash-only; the Task 1 brief already exists and later briefs can be created with `apply_patch` on Windows.
- Do not expose `.env`, OAuth tokens, credentials, or raw buyer payloads.
- Do not mark M-06 passed without persisted `validation-result.md`, full cold gate, live browser evidence, and no failing starred criteria.

## Important Files

- `docs/superpowers/specs/2026-07-09-m06-f03-order-realization-design.md`
- `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`
- `.superpowers/sdd/m06-f03-task-1-brief.md`
- `.superpowers/sdd/m06-f03-correction-report.md`
- `.superpowers/sdd/m06-f01-correction-report.md`
- `.superpowers/sdd/m06-f02-correction-report.md`
- `.superpowers/sdd/progress.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-03-profit-snapshot-calculation/spec.md`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `apps/server_core/internal/modules/profitability/ports/order_reader.go`
- `apps/server_core/internal/modules/profitability/adapters/orders/order_reader.go`
- `apps/server_core/internal/modules/profitability/adapters/orders/order_reader_test.go`

## Environment Notes

- Repository: `C:\Users\leandro.theodoro\Documents\marketplace-central`
- Local date/timezone: 2026-07-09, America/Sao_Paulo.
- A disposable Docker PostgreSQL container named `mpc-m06-test-postgres` was used earlier. Verify whether it still exists before relying on it; recreate it if needed.
- Do not treat the disposable database as live provider proof; it is real PostgreSQL integration evidence only.

## Next Session Prompt

```text
Work in C:\Users\leandro.theodoro\Documents\marketplace-central on main. Read AGENTS.md and C:\Users\leandro.theodoro\.codex\attachments\2a1a5352-6aee-406c-bf04-2d066266687f\goal-objective.md first, then read docs/superpowers/handoffs/2026-07-09-m06-orders-margin-session-handoff.md, the F-03 design, plan, and .superpowers/sdd/m06-f03-task-1-brief.md.

Resume M-06 F-03 using superpowers:subagent-driven-development. Task 1 is stopped at RED: only order_reader_test.go and service_test.go were changed; no production code or GREEN exists. Re-run and preserve RED, implement the exact profitability-owned OrderFact/OrderItemFact contract, keep orders imports only in profitability/adapters/orders, run GREEN + gofmt + boundary rg, write the Task 1 report, and obtain independent SPEC/QUALITY approval before Task 2.

Guardrails: preserve the heavily dirty worktree; never reset/revert unrelated changes; do not use WSL; do not expose secrets; no unknown-to-zero or unknown-to-realized; OpenAPI/SDK together; mocks are not live evidence. After Tasks 1-4, run real PostgreSQL, Mercado Livre + Oracle, built-in-browser desktop/mobile QA, and the full cold M-06 gate. Do not claim M-06 passed before all real evidence and validation artifacts exist.
```
