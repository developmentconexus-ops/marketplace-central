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

## Corrective round 2 (dual-gate divergence @ 13b49e2c → M02-COR-2)

Gate round 2 diverged: Opus PASS, Sol FAIL on M02-C03 app-wide (ProductLinksPage:184,
OrdersPage:329, IntegrationsHubPage:210, MarketplaceSettingsPage:62 still fetched
installations directly; AppRouter test mocks hid the fetches). Operator ruled REMOVAL,
not migration (pages rebuilt from scratch next mission, design handoff
`docs/design/handoff-2026-07/`). Fixed in e4c8ea90 (ledger rows 32–33, reviewer
ACCEPT-WITH-CONDITIONS — both importants process-level, recorded below):

- **Removal, not just unmount:** nothing outside the 4 packages referenced their code
  (verified: only AppRouter import/mock sites, own package tests, package.json/lock rows,
  and historical docs/.mnfs evidence). All 4 packages DELETED entirely (−7312 lines):
  `feature-product-links`, `feature-orders`, `feature-integrations`, `feature-marketplaces`.
- `/vinculos`, `/pedidos`, `/integracoes`, `/marketplaces` mount the existing
  `WorkspacePlaceholder` stub (same pattern as `/catalogo/produtos/:productId` and
  `/protocolos/:protocolId`). The 6 LegacyRedirect routes byte-unchanged (M02-C01).
  `/estoque` untouched (already compliant per COR-1).
- AppRouter test: 4 legacy page mocks deleted; `listIntegrationInstallations` exactly-once
  assertion now honestly app-wide (remaining mocked pages contain zero installation fetches,
  grep-verified) and asserted on every stub route against the real InstallationProvider.
- **PLAN-ADJUDICATION #3 SUPERSEDED by operator directive:** the adjudicated position
  (legacy pages stay mounted as-is during M-02) is void; removal is the ratified end-state.

Ladder re-run at e4c8ea90: L0 build exit 0 (bundle 438.77 → 361.17 kB) · L0 governance
clean-worktree BaseSha 40-hex exit 0 (baseline exceptions unchanged) · L1 web 133/133
(22 files) exit 0. Zero-grep for the 4 package names over apps/+packages/.

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
7. package-lock regen (M02-COR-2): chip ran `npm install --package-lock-only` to sync the
   lockfile after the operator-ordered package removal (−70 lines, deletions only, verified
   scoped to the 4 removed entries). Consequence of the removal directive, not an
   independent dep change — but it bypassed the formal REQUEST path; hub ratification
   requested (reviewer important #1).
8. index.css @source cleanup (M02-COR-2): 4 dead tailwind `@source` rows pointing at the
   deleted package dirs removed by chip — outside the worker's declared scope but required
   for a truthful build. Hub ratification requested (reviewer important #2).

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
