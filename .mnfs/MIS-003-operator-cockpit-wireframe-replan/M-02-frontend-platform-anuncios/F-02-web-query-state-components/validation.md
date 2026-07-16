# F-02 web-query-state-components — validation evidence

Feature closed 2026-07-16 by CHIP-M02. Plan source: `../plan-batch.md` §F-02 slices (Sol-medium,
ledger row 1). Spec: `feature.md` + IC-05 (`../../research/frontend-platform-interface-contract.md`).

## Slices

| Slice | Commit | Worker | Review (ledger row) | Result |
| --- | --- | --- | --- | --- |
| F02-S1 registry/namespaces/key builders (+ S6a vitest discovery) | 46c93652 | Luna high (row 2) | row 4 ACCEPT-WITH-CONDITIONS (suggestion-level) | 45 tests green |
| F02-S2 invalidateAfterMutation crosswalk | d1eaeab6 | Luna high (row 5) | row 7 ACCEPT (suggestion only) | 12 tests green |
| F02-S3 failureCopy (12 IC-03 codes) | a5c915de | Luna high (row 8) | row 9 ACCEPT-WITH-CONDITIONS — important fixed (failureCodes derived from map keys) | 3 tests green |
| F02-S4 LoadingState/ErrorState/EmptyState | 93bb9d5e | Luna high (row 3) | row 6 ACCEPT-WITH-CONDITIONS — important fixed (Button reuse) | 4 tests green |
| F02-S5 UnknownValue/ConflictTag/FreshnessIndicator | 4300db5d | Luna high (row 10) | row 12 ACCEPT-WITH-CONDITIONS — suggestions applied | 5 tests green |
| F02-S6b tailwind @source + 7 proxy rows + /orders bypass | (this commit) | Luna high (row 18) | row 19 ACCEPT | 1 test file green |

Supporting dep commit: 4e385b1f (hub-granted workspace rows web→ui, ui→web-query; lock additive-only).

## Criteria coverage (component level)

- M02-C04 (invalidation crosswalk): invalidation.test.ts — exact namespace set per IC-05 row
  (presence AND absence), unknown type → typed error with zero invalidations. Exhaustive via
  `satisfies Record<MutationInvalidationType,…>`. `product_enrichment` discriminator per
  adjudication #7.
- M02-C05 (staleTime registry): index.test.ts snapshot — exact QUERY_STALE_TIME incl. listings
  45_000 / mutations 5_000 / orders 120_000 / sync 30_000 / market 300_000; key builders verbatim.
- M02-C06 (six state components): LoadStates.test.tsx + FactStates.test.tsx — byte-exact pt-BR
  copy (Carregando… / Erro ao carregar. + Tentar novamente / Nenhum registro encontrado. /
  — + hint / divergente amber / dados de HH:MM:SS via formatAsOf).
- M02-C07 (failureCopy exhaustive): failureCopy.test.ts — all 12 codes exact literals, unknown →
  byte-exact `Falha desconhecida ({code})`. Literals chip-authored per adjudication #6, submitted
  for IC-05 ratification at CLOSED.
- M02-C10 enabler (proxy contract): viteProxy.test.ts — 7 IC-05 rows (incl. hub-directive
  /dashboard + /sync), /orders HTML-Accept bypass both ways, SPA-route-key guard.

## Rulings applied

- ui FreshnessIndicator is canonical (pt-BR aria-label); legacy web-query FreshnessIndicator kept
  untouched for feature-page consumers (re-export = package cycle); migration via Migration Briefs.
- installations namespace/key = adjudication #1, submitted for IC-05 ratification.

## Suite state at feature close

`cd apps/web; npm test -- --configLoader native` → 15 files, 86/86 green. Plain `npm test`
config-load failure = pre-existing Node 26/esbuild issue (field finding in CLOSED payload).
No typecheck lane exists for packages/* (S3 reviewer finding — field finding in CLOSED payload).
