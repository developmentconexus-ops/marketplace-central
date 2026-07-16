# F-01 shell-routes-context — validation evidence

Feature closed 2026-07-16 by CHIP-M02. Plan source: `../plan-batch.md` §F-01 slices (Sol-medium,
ledger row 1). Spec source: `feature.md` + IC-05 (`../../research/frontend-platform-interface-contract.md`).

## Slices

| Slice | Commit | Worker | Review (ledger row) | Result |
| --- | --- | --- | --- | --- |
| F01-S1 InstallationContext + InstallationGate | 72862997 | Sol low (row 11) | row 13 ACCEPT-WITH-CONDITIONS — important fixed (shell-intact gating split) | 6 scoped tests green |
| F01-S2 route map + 6 query-preserving redirects | 7ab8adff | Luna high (row 14) | row 15 ACCEPT-WITH-CONDITIONS — important fixed (2 direct-mount assertions added) | 18 scoped tests green |
| F01-S3 Layout sidebar + ML pill + gated outlet | de37830c | Luna high (row 16) | row 17 ACCEPT-WITH-CONDITIONS — suggestion fixed (option labels) | 5 scoped tests green |

## Criteria coverage (component level; browser QA is hub-owned P7)

- M02-C01 (route map + redirects): AppRouter.test.tsx — table-driven 6-redirect test asserts
  destination path, full query string preserved, history replacement; direct-mount tests for all
  IC-05 routes incl. `/anuncios` placeholder, `:productId`/`:protocolId` placeholders,
  off-sidebar config routes; unknown path renders no workspace. 18/18 green.
- M02-C02 (installation context persistence): InstallationContext.test.tsx — unknown-id fallback
  + URL rewrite, known-id no-rewrite, selection change preserves pathname + unrelated params,
  exactly one `listIntegrationInstallations` across 3 routes + selection change; Layout.test.tsx —
  sidebar NavLinks carry `?installation=`, ML pill switch preserves route/query. Green.
- M02-C03 (context-owned single fetch): single `useQuery` on `installationsQueryKeys.list()`;
  no localStorage/raw fetch/useEffect fetch (review-verified). Legacy pages' own fetch = recorded
  debt per adjudication #3 (ADR-15 + Migration Briefs), not defect.
- Shell-intact rule (feature.md): InstallationGate gates page content only; regression tests in
  InstallationContext.test.tsx (`shell` testid) and Layout.test.tsx (nav visible during empty).

## Suite state at feature close

`cd apps/web; npm test -- --configLoader native` → 15 files, 86/86 green (includes all F-02
package tests). Plain `npm test` fails at config load — pre-existing Node 26/esbuild issue,
reported as field finding in CLOSED payload.

## Deviations / notes

- StrictMode omitted from new test harnesses for deterministic single-request assertions
  (declared by worker, consistent with existing repo test convention; risk noted by reviewer).
- `InstallationGate` exported alongside provider (same file) — mounted by F01-S3 Layout; not dead code.
- ML pill implemented as native `<select>` inside pill wrapper (no new deps).
