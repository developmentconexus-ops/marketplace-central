# F-04 preview-confirm-ui — validation record

Chip: CHIP-M03 · branch `chip/m-03-mutation-envelope-writes` · closed 2026-07-17.
Ungated by hub rebase trigger naming M-02 F-03; chip rebased onto main @ 79d6787f (D-63)
before any F-04 work. Plan D-64 (Sol medium): 5 slices, DAG S1→S2→{S3,S4}, S4→S5.

## Slices → commits → reviews

| Slice | Commits | Review (ledger) | Verdict |
| --- | --- | --- | --- |
| S1 modal intent+preview state machine (complex) | fe268655 | D-65/D-66 | ACCEPT (unconditional) |
| S2 approve/poll/terminal/invalidate-once (complex) | a84f0b44 (incl. corrective after worker capacity death, D-68 deviation) | D-67/D-68/D-69 | ACCEPT (unconditional) |
| S3 Anúncios launchers (standard, ∥S4) | 19075cd7 | D-70/D-71 | ACCEPT (unconditional) |
| S4 /protocolos/:id detail (complex, ∥S3) | 42574869 | D-70/D-71 | ACCEPT (unconditional) |
| S5 retry nav + AppRouter route swap (standard) | 1d7a5bad | D-72/D-73 | ACCEPT (unconditional) |

Corrective round M03-COR-1 (post dual-gate round 1) touched F-04 surface:
2ba61fa3 (COR-1d: failure code literal + pt-BR copy in item row; canRetry via protocol
totals `failed>0` instead of current-page scan — closes the D-73 non-blocking carry).

## Contract evidence (M03-C10, C11 served; C05 surface)

- C10 confirm-flow: modal intent entry for all six enabled types (IC-03 intent shapes
  byte-checked vs contract at D-66); preview step mandatory — confirm control genuinely
  absent pre-preview (structural, not disabled-button theater); approve sends literal
  `execute:true`; `preview_stale` structurally swaps approve → "Gerar prévia novamente"
  (no silent re-approve path).
- C11 protocol detail: `/protocolos/:protocolId` real page (route swap 1d7a5bad, element
  swap only per seam); header + type label + actor + retried_from back-link; ADR-17
  unknowns render "—" (never 0); items table keyset pagination via cursor; failed item
  row shows failure code literal (mono) + pt-BR copy (`failureCopy`), `message_pt` NOT
  rendered, provider message only behind "▸ técnico" disclosure (C05 no-leak surface);
  polling resumes on direct access, stops at terminal state, cache invalidation exactly
  once per terminal transition (WeakMap<QueryClient,Set<protocolId>> guard).
- Retry: exactly-once dispatch (synchronous ref guard, double-click proven), disabled
  while pending, navigates to new protocol preserving location.search, retried_from
  back-link on the new page; retry visibility from protocol totals (`failed>0`,
  non-finite/missing totals → hidden, honest null — COR-1d).
- Installation scoping: launchers clear selection + close modal on installation change;
  detail page useInstallation-only rule respected (D-73 verified through full AppRouter
  with InstallationProvider live).

## Lanes

- Web vitest: green every slice, re-run cold by reviewers (S1 11/11 + full 144; S2 18/18
  3x deterministic; S3 29 tests; S4 4/4 twice; S5 23/23 twice; suite at F-04 close
  25 files / 164 tests). `npm run build` green (S5 + ladder re-run).
- Timer discipline (D-68 doctrine): pure `vi.useFakeTimers()`; no waitFor with fake
  timers; no shouldAdvanceTime (proved flaky at 1999/2000 boundary); explicit
  advanceTimersByTimeAsync / runOnlyPendingTimersAsync choreography; cadence proven at
  exact 1999/+1 boundary + stop-at-terminal.
- Backend lanes untouched by F-04 (FE-only write set); full ladder re-run at chip head
  covered go sweep + integration + governance (LADDER.md).

## Live-write posture (security-critical)

Unchanged: UI drives the stub-gated mutation API only; `MPC_PROVIDER_WRITES_ENABLED`
untouched; zero live Mercado Livre writes. Live ML execution remains gated on explicit
operator authorization via hub ESCALATION.

## Deviations / carries out of F-04

- D-68: codex corrective worker died at capacity mid-run ("Selected model is at
  capacity"); fallback per core §1 (sonnet cavecrew-builder + orchestrator empirical
  test stabilization, assertions preserved). D-72: same-pattern orchestrator
  stabilization (`<output>` implicit role=status collision → div; flush rounds).
- Workers blind throughout (codex sandbox vitest denial — esbuild directory read);
  orchestrator ran all suites; reviewers re-ran independently.
- D-73 carry (canRetry current-page-only) CLOSED by COR-1d (2ba61fa3).
- S4 local ItemsTable vs S1 MutationItemsTable duplication = correctly scoped (D-71).
- Modal-close-during-applying: no dedicated automated test (D-69 non-blocking);
  structural guarantee + hub browser QA cover.
