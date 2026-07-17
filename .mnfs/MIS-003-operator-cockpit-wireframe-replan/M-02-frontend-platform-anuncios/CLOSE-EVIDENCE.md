# M-02 frontend-platform-anuncios — chip close evidence (CLOSED payload)

Chip: CHIP-M02 · branch `chip/m-02-frontend-platform-anuncios` · base a49168e641ffd6f61932ca57c29b1d1bdcde2fb0
Closed 2026-07-17. Milestone close (P6 dual-gate + P7 browser QA) is HUB-owned; chip stops here.

## Corrective round (dual-gate FAIL @ 8744c280 → M02-COR-1)

First dual gate returned FAIL with 2 blockings; both fixed in 13b49e2c (ledger rows 30–31,
reviewer ACCEPT-WITH-CONDITIONS, both suggestions applied):

- **BLOCKING-1 (M02-C03):** StockSeguroPage fetched installations directly. Fix:
  `listIntegrationInstallations` removed from StockSeguroClient interface (page is
  type-incapable of the fetch), fetch effect + page default-installation fallback deleted
  (InstallationProvider owns both), installations = required prop injected by the /estoque
  wrapper from `useInstallation()`. AppRouter test asserts exactly-once app-wide.
  Supersedes the narrowing previously declared in `F-01-shell-routes-context/validation.md`.
- **BLOCKING-2 (M02-C08):** no visible chip for active exception filter. Fix: dismissible
  chips for exception / sync_state / link_state / listing_type_code with exhaustive pt-BR
  maps (`satisfies`); dismiss removes only that filter key, preserves tab/q/installation/
  remaining filters. New chip label copy (Exceção: Erro de sync / Desatualizado / Sem vínculo /
  Abaixo da margem; Sync:/Vínculo:/Modalidade: prefixes; aria `Remover filtro {label}`) —
  added to ratification list below.

Ladder re-run at 13b49e2c: L0 build exit 0 · L0 governance clean-worktree BaseSha 40-hex
exit 0 (baseline exceptions unchanged) · L1 web 133/133 (22 files) exit 0. L0 typecheck
unchanged pre-existing F4.

## Features

| Feature | Evidence | COMMITTED |
| --- | --- | --- |
| F-01 shell-routes-context | `F-01-shell-routes-context/validation.md` (37adc286) | ✅ |
| F-02 web-query-state-components | `F-02-web-query-state-components/validation.md` (b3ddfcfd) | ✅ |
| F-03 anuncios-workspace | `F-03-anuncios-workspace/validation.md` (2b5a7059) | ✅ |

Dispatch ledger: `DISPATCH-LEDGER.md` rows 1–29 (14 implement workers, 1 planner, 14 reviews;
every slice reviewed before dependents; all verdicts ACCEPT or ACCEPT-WITH-CONDITIONS with
conditions fixed in-slice).

## Verification ladder (profile §2)

- **L0 build:** `npm run web:build` → vite build exit 0 (1850 modules, 438.77 kB js).
- **L0 typecheck:** `npx tsc --noEmit` in apps/web fails `TS2688 Cannot find type definition
  file for 'node'` — IDENTICAL failure on main checkout at 449156f4 → PRE-EXISTING, not chip
  drift (field finding F4 below). No working web typecheck lane exists.
- **L0 governance:** `scripts/harness.ps1 -Command governance -BaseSha a49168e641ffd6f61932ca57c29b1d1bdcde2fb0`
  from clean detached worktree at 2b5a7059 → exit 0, baseline exceptions only (14, unchanged),
  `artifact_path=contracts/governance`. Anchor = milestone accepted base SHA per profile §2.
- **L1 web:** `cd apps/web; npm test -- --configLoader native` → 22 files, 130/130 green,
  exit 0 (re-verified after every slice). packages/ui DetailPanel 8/8 scoped.
- **L1 Go:** not applicable — zero Go/backend files touched (branch diff outside
  apps/web + packages + .mnfs = package-lock.json +4 lines only, hub-granted).
- **L2:** hub-owned (chip never boots stack). Deferred to hub QA live-drive.

## Ratification requests (IC-05 gap-fills + copy)

1. `installations` namespace/key (`installationsQueryKeys.list`) — adjudication #1.
2. failureCopy pt-BR literals for 12 IC-03 codes + `Falha desconhecida ({code})` fallback —
   adjudication #6.
3. `product_enrichment` invalidation discriminator — adjudication #7.
4. `syncQueryKeys.runs` polling pattern (refetchInterval 2s, stop at terminal) — adjudication #12.
5. Refresh run status copy chip-pinned pt-BR: na fila / em andamento / concluído / falhou /
   cancelado (F03-S5).
6. Filter chip copy chip-pinned pt-BR (M02-COR-1): exception map Erro de sync / Desatualizado /
   Sem vínculo / Abaixo da margem; chip prefixes Exceção/Sync/Vínculo/Modalidade; dismiss aria
   `Remover filtro {label}`.

## Field findings

- **F1 Node 26/esbuild:** plain `npm test` fails at vitest config load; `--configLoader native`
  required. Affects every worker prompt (documented in each).
- **F2 npm teardown flake:** rare exit-1 with all-green vitest output (seen 2×, F02-S2/F03-S3);
  one re-run always clean. Candidate for L1 allowlist annotation.
- **F3 packages/* typecheck gap:** no typecheck lane covers packages/ui / web-query (F02-S3
  reviewer finding).
- **F4 apps/web typecheck broken pre-existing:** `tsc --noEmit` TS2688 ('node' types) on main
  AND chip — L0 typecheck lane for web is currently impossible repo-wide. Backlog: fix
  tsconfig types entry or add @types/node to workspace (dep change = hub-owned).
- **F5 codex worker teardown kill:** OS-process worker killed at session restart AFTER edits
  complete (F03-S3) — tree-as-delivery recovery worked; `.done` sentinel absence + complete log
  is the signature.

## Notes for hub

- package-lock diff: +4 lines, additive-only, the two hub-granted workspace rows (4e385b1f).
- IC-05 `/dashboard` + `/sync` proxy rows added per hub directive (supersession note in F-02
  evidence).
- SummaryCounter vs shared StatCard divergence is deliberate (wireframe-2a; StatCard.value
  can't host UnknownValue JSX).
- docker/dev/*.sh EOL-normalized by hub in place — never committed by chip.
