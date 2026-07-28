# DISPATCH LEDGER — CHIP-IMPORT-CHAIN (MIS-006 / M-06 / wave 2)

Every dispatch this chip made, with the PATH of the persisted artifact. A row whose artifact is
0 bytes is an unrecoverable verdict — that is the CHIP-ANCHORS-2 failure this ledger exists to
prevent. Byte counts below are the proof the artifacts are real, not streamed away.

Paths are relative to `.mnfs/MIS-006-integracao-fundacao/_chip-import-chain/`.

| # | Role | Model / effort | Sandbox | Slice | Started (UTC) | Ended (UTC) | Exit | Artifacts (bytes) |
|---|---|---|---|---|---|---|---|---|
| 1 | implementer-standard | gpt-5.6-luna / high | workspace-write | S1 — `/importacoes` route + section promoted out of `pages/vinculos/` | 2026-07-28T15:45:56Z | 2026-07-28T15:52:41Z | 0 | `dispatches/agent__implementer-standard-20260728T154550Z-410e32e1.last.md` (1060) · `.log` (262653) · `.result.json` |
| 2 | implementer-standard | gpt-5.6-luna / high | workspace-write | S2 — real `getErpImportChain` consumption + detail screen | 2026-07-28T16:01:09Z | 2026-07-28T16:07:25Z | 0 | `dispatches/agent__implementer-standard-20260728T160103Z-c0849f51.last.md` (955) · `.log` (392302) · `.result.json` |
| 3 | gate-review (adversarial, cold) | Claude Opus 5 / read-only tools | n/a | Full-chip review vs `validation-contract.md` I1..I9 | 2026-07-28 | 2026-07-28 | see artifact | `REVIEW-ADVERSARIAL.md` |

Inputs each dispatch carried (all persisted, all under `dispatches/`):
`pack-impl-v1.0.0.md` (canonical impl-pack, copied verbatim from MIS-004 M-03) ·
`bindings-import-chain.md` (repo/env bindings + the skill-discovery denylist clause verbatim) ·
`card-s1-route.md` / `card-s2-chain.md` (slice cards) · `prompt-s1.md` / `prompt-s2.md` (the
assembled prompts, dispatched from file — never inline) · `dispatch-registry.jsonl`.

## Commits produced

| Commit | Author of the code | Subject |
|---|---|---|
| `4b76a287` | dispatch #1 | `feat(importacoes): add dedicated history route` |
| `1bffdcfd` | dispatch #2 | `feat(web): wire ERP import chain detail` |
| `67e4a3d` | **chip session (glue, disclosed)** | `test(importacoes): repair the suites the chain link turned red` |

`67e4a3d` is the only code this orchestrating session wrote: ~16 lines across three test files,
wrapping two render helpers in `MemoryRouter` and giving one assertion the time the hook's own
retry actually takes. No assertion was changed or weakened, no production behavior added. It is
recorded here rather than folded silently into a worker's commit.

## Field findings (candidates for profile ratification)

- **F-1 — sandbox vitest blindness has teeth.** `codex exec --sandbox workspace-write` on Windows
  cannot run vitest (esbuild `Access is denied`), so both workers committed on typecheck alone.
  Dispatch #2 added a `<Link>`, which made `ImportacaoSection` require a router context and turned
  `ImportacoesPage.test.tsx` and `IntegracoesPage.test.tsx` red — invisible to the worker, which
  reported "typecheck clean, committed". The chip's post-dispatch vitest re-run caught 2 broken
  files plus 1 genuine timing defect. The re-run is not a formality; it is the verification of record.
- **F-2 — `New-DispatchPrompt.ps1` does not validate the role string.** A role typo
  (`ImplementStandard`) assembles cleanly and only fails later at `Invoke-CodexDispatch.ps1` with
  `ROLE-UNKNOWN`. Cheap fix: validate against `roles.psd1` at assembly time.
- **F-3 — Windows PowerShell 5.1 mangles the UTF-8 `HarnessDispatch.psm1`** (em-dash inside a
  string → `Token 're-vendor' inesperado`). `pwsh` (PowerShell 7) is required, which is what the
  profile binds anyway — worth stating as a hard requirement rather than a preference.
- **F-4 — vitest in a junction-only worktree needs an absolute `setupFiles` path plus
  `server.fs.strict: false`.** The junction's realpath resolves outside the vite root and the fs
  guard rejects it, even though the file exists.
- **F-5 — codex binary available here is 0.144.4**, while profile F-ENV-9 recommends
  0.145.0-alpha.18 (not installed). Both dispatches completed exit 0; noted, not a blocker.
