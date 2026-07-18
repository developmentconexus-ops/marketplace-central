# P7 QA — FINAL VERDICT — M-04 vinculos-import-ui (MIS-004-mvp-demo)

## VERDICT: **PASS**

Independent QA validator, 2026-07-18. Combines PART A (UI live-drive, `p7-qa-live-drive.md`) + PART B (API contract sweep, `p7-qa-api-sweep.md`, executed live against `http://localhost:8080`). Scope = hub-ratified scoped-QA. Zero-ML-writes guaranteed (dispatcher OFF); no provider write attempted.

## Scope items — per-item result

| # | Scope criterion | Result | Evidence (status codes) |
|---|---|---|---|
| 1 | Screen renders (shell/nav/queue/tabs/KPI) | PASS | PART A live accessibility tree; /vinculos + Fila/Resolvidos tabs + KPI tiles + queue columns rendered |
| 2 | NO_CANDIDATE honesty live (ADR-17) | PASS | PART A "sem confiança sem candidato" + honest UNAVAILABLE anchors ↔ PART B step1 API confidence=0 int, reasons[] {anchor,direction,detail}, 34 items, zero fabricated ACCEPT (GET 200) |
| 3 | approve mechanics (manual-resolve→2002, leaves queue) | PASS | manual-resolve 200, current_link→resolved (drops from Pendentes); undo-to-none cycle on MLB3758134297 200/200 (state→none) |
| 4 | undo audit chain (reversal + guards) | PASS | reversal 200 (append-only undo audit); SUPERSEDED 409; ALREADY_UNDONE 409; AUDIT_NOT_FOUND 404 |
| 5 | batch-preview dry-run (zero-persist) | PASS | 200, items FAILED+cause, re-GET snapshot UNCHANGED |
| 6 | batch-apply honesty (ADR-17 no-audit-for-all-failed) | PASS | 200, applied=[], failed[] w/ honest cause, batch_id="" |
| 7 | empty batch guard | PASS | 422 PRODUCT_LINKS_BATCH_APPROVALS_REQUIRED |
| 8 | approve-candidate refuses NO_CANDIDATE (no fabricated write) | PASS | 400 PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE (honest refusal; code-nuance note below) |
| 9 | error-code contract table | PASS | 200/409/409/404/422/400/404 all honest + consistent; two prompt-expectation mismatches reconciled to OpenAPI (below) |
| 10 | C04 import history + rejection report | PASS | #001-E 55/0/0, #002-E 3/1/2; rejected_rows[row/code/detail] + itemized warnings (GET 200) |

## Notes (non-blocking, reconciled against authoritative OpenAPI — repo truth > prompt assumption)
- **manual-resolve → 999999 returned 200, not 404.** OpenAPI declares manual-resolve with **only a 200 response**; `ValidateInternalProductID` validates positivity only (no catalog-existence check by design — operator override). NOT a contract violation. OBSERVATION forwarded: no referential validation of internal_product_id (UI product-picker mitigates); owner may reconcile the step-8 expectation with OpenAPI or add existence validation. Not a QA fail.
- **approve-candidate on NO_CANDIDATE returned 400 (code PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE), not 422/INSUFFICIENT_ANCHORS.** 422/INSUFFICIENT_ANCHORS is declared only for empty-batch on batch endpoints. Honest refusal with a consistent domain code; 400-vs-422 is a transport-mapping nuance. Not a fail.
- **manual-resolve idempotent re-resolve returned 200 (no `cause`).** Contract declares only 200; per-item `cause` is batch semantics. Idempotent + honest. Creates a resolved→resolved no-op audit row (append-only; slightly noisy). Not a fail.
- Band-spectrum happy-path + numeric "95%" confidence chip: OUT of live scope per hub ruling (synthetic catalog = zero overlap ⇒ all NO_CANDIDATE is honest). Covered by 10 green IC-01-A2 unit fixtures + P6 gate. Live absence recorded, NOT failed.

## Risks / residual
- Live-state residual: `MLB3758134295` left `current_link.state=resolved→2002` (valid synthetic product; the transient 999999 write was reversed). Full unresolve-to-pending is unreachable via API (superseded original audit + no unresolve endpoint) — a test-sequence artifact, benign. RECOMMEND hub reset this link (regen candidates/DB) if the demo needs it back in the Pendentes queue. All other test writes fully reversed; `MLB3758134297` = clean (none).
- Design observation (not a defect): once an audit is superseded, its undo is permanently blocked and there is no unresolve-to-pending path. Acceptable under "undo = reverse the latest action" semantics; flagged for product awareness.

## Fields
- **Status:** Complete
- **Validation verdict:** PASS
- **Contract checked:** M-04 vinculos-import-ui hub-ratified scoped-QA (screen render; approve/reject/batch-preview/batch/undo contract; NO_CANDIDATE honesty; C04 import history+rejection report; resolve/undo mechanics; error-code contract table), grounded in `contracts/api/marketplace-central.openapi.yaml` + product_links domain/application/transport source.
- **Artifact paths:** `.mnfs/MIS-004-mvp-demo/M-04-vinculos-import-ui/evidence/p7-qa-api-sweep.md`, `.../p7-qa-live-drive.md`, `.../p7-qa-verdict.md`.
- **Evidence/commands:** live curl/python against http://localhost:8080 (GET link-candidates, POST manual-resolve/approve-candidate/batch-preview/batch/undo, GET link-workflows, GET /erp/imports[/{id}]); source grounding via OpenAPI + resolution_service.go / internal_product_id.go / http_handler.go.
- **Blocking failures:** NONE.
- **Recommended correction scope:** None required for PASS. Optional (owner decision, feature-scope): reconcile QA step-8/step-12 expected error codes with OpenAPI, or add internal_product_id catalog-existence validation to manual-resolve; provide an unresolve/queue-reset path; hub reset MLB3758134295 pre-demo.
- **Required next inputs:** none to close the gate.
- **Handoff reason:** Scoped-QA verdict issued; hand to Milestone Orchestrator / hub for M-04 close decision + optional pre-demo link reset.
