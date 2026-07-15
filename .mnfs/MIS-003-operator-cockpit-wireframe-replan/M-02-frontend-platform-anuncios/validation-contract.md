# Milestone Validation Contract

```yaml
id: M-02
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: milestone
```

## Milestone ID

M-02-frontend-platform-anuncios

## QA Level

QA-2 — vitest component/unit suites + manual browser walkthrough with screenshots (browser automation lane unconfigured this mission).

## Required Outcome

IC-05 platform seam live (routes, redirects, context, query contract, state components) proven by the rebuilt Anúncios workspace.

## Criteria

## Criterion: Redirect map complete
ID: M02-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest router test + browser visit of each legacy URL with query string
- Expected: `/products→/catalogo`, `/product-links→/vinculos`, `/inventory/stock-seguro→/estoque`, `/orders→/pedidos`, `/integrations→/integracoes`, `/simulator→/precos`; each preserves the full query string (e.g. `?installation=inst_1` present after redirect)
- Actual:
- Artifact: `F-01-shell-routes-context/validation.md` + screenshots
Drive (UI — agent-browser; UI criteria only; omit for non-UI):
- Fixture: dev seed with installation `inst_1`
- Steps:
  - open http://localhost:5174/inventory/stock-seguro?installation=inst_1
  - assert url ~ /estoque\?installation=inst_1
- Expected: browser lands on /estoque with query intact
Blocking failure: any redirect missing or dropping query
Blocking failure observed: No
Owner: QA Validator

## Criterion: Installation context URL persistence
ID: M02-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest: unknown `?installation=inst_ghost`; switch installation via pill; route change
- Expected: unknown id → first installation selected AND URL rewritten to its id; pill switch updates `?installation=` without reload; param carried across route navigation
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: silent mismatch (UI shows A, URL says B) or param lost on navigation
Blocking failure observed: No
Owner: QA Validator

## Criterion: Context single fetch
ID: M02-C03
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: vitest with fetch spy across 3 route changes + 1 installation switch
- Expected: exactly 1 call to `listIntegrationInstallations`; zero direct installation fetches from pages
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: >1 call or page-level installation fetch
Blocking failure observed: No
Owner: QA Validator

## Criterion: Invalidation crosswalk exact
ID: M02-C04
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: vitest unit test with `invalidateQueries` spy, one case per IC-05 crosswalk row
- Expected: `price_update`→{listings,mutations}; `stock_correct`→{listings,inventory,mutations}; `link_apply`→{listings,linkage,catalog,mutations}; `listing_pause|listing_resync|listing_edit`→{listings,mutations}; product-enrichment→{catalog}; unknown type → thrown typed error
- Actual:
- Artifact: `F-02-web-query-state-components/validation.md`
Blocking failure: any row invalidating a different namespace set
Blocking failure observed: No
Owner: QA Validator

## Criterion: staleTime registry mirrors contract
ID: M02-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: vitest snapshot of `QUERY_STALE_TIME`
- Expected: `{catalog:300000, stock:45000, pricecost:120000, listings:45000, mutations:5000, orders:120000, sync:30000, market:300000}` — existing three unchanged
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: any value differing or existing key mutated
Blocking failure observed: No
Owner: QA Validator

## Criterion: State vocabulary components render fixed copy
ID: M02-C06
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest render tests per component
- Expected: LoadingState "Carregando…"; ErrorState "Erro ao carregar." + button "Tentar novamente"; EmptyState "Nenhum registro encontrado."; UnknownValue "—" + hint tooltip; ConflictTag "divergente" amber; FreshnessIndicator "dados de HH:MM:SS" via formatAsOf
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: copy divergence from IC-05 table (byte-exact) or English string
Blocking failure observed: No
Owner: QA Validator

## Criterion: failureCopy exhaustive
ID: M02-C07
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest test iterating all 12 IC-03 failure codes + one unknown code
- Expected: each code returns a non-empty pt-BR string; unknown code → "Falha desconhecida ({code})"
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: missing code or English fallback
Blocking failure observed: No
Owner: QA Validator

## Criterion: Anúncios deep link survives reload (Q5)
ID: M02-C08
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser walkthrough on seeded data
- Expected: `/anuncios?installation=X&tab=pendencia&filter.exception=sync_error` after F5 renders Pendências tab active, filter chip visible, and only sync-error rows; sidebar shows 8 items in IC-05 order
- Actual:
- Artifact: `F-03-anuncios-workspace/validation.md` before/after-F5 screenshots
Drive (UI — agent-browser; UI criteria only; omit for non-UI):
- Fixture: IC-02 seed MLBTEST0001..0006, installation inst_1
- Steps:
  - open http://localhost:5174/anuncios?installation=inst_1&tab=pendencia&filter.exception=sync_error
  - assert text "com erro"
  - assert url ~ tab=pendencia
- Expected: identical filtered view pre/post reload
Blocking failure: tab/filter/installation reset on reload
Blocking failure observed: No
Owner: QA Validator

## Criterion: Anúncios renders unknowns and 409-attach
ID: M02-C09
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest: null-cost row render; refresh-while-running mock returning 409
- Expected: margem cell renders "—" with hint "sem custo no ERP → não simulado" (never "0"); on 409 `refresh_in_progress` UI attaches to the running operation (progress shown, no error toast)
- Actual:
- Artifact: `F-03.../validation.md`
Blocking failure: zero-rendering or 409 surfaced as error
Blocking failure observed: No
Owner: QA Validator

## Criterion: No direct fetch in new code (Q6)
ID: M02-C10
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: grep for `useEffect`-fetch / raw `fetch(`/`client.` calls outside web-query hooks in files added by M-02
- Expected: zero matches; all server state through `listingsQueryKeys`/context; governance lanes green at milestone SHA
- Actual:
- Artifact: milestone `validation-result.md` grep transcript + lane logs
Blocking failure: any direct fetch in new/rebuilt code
Blocking failure observed: No
Owner: QA Validator

## Criterion: Anúncios first data paint (Q1)
ID: M02-C11
Level: Milestone
Type: Performance
Required: Yes
Status: Pending
Evidence:
- Command: browser walkthrough on IC-02 seeded data (2k rows), dev-local build
- Expected: from navigation to `/anuncios?installation=inst_1` until first table rows render < 2s (mission Q1 "UI first data paint < 2s dev-local"); measured via browser performance trace or timestamped screenshots
- Actual:
- Artifact: `F-03-anuncios-workspace/validation.md` timed trace/screenshot pair
Drive (UI — agent-browser; UI criteria only; omit for non-UI):
- Fixture: IC-02 seed at 2k listings, installation inst_1
- Steps:
  - open http://localhost:5174/anuncios?installation=inst_1
  - assert table rows visible; record elapsed time from navigation start
- Expected: first data paint elapsed < 2s
Blocking failure: first data paint ≥ 2s on seeded dev-local run
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Feature proofs `F-0*/validation.md`; rollup `validation-result.md` with fixed SHA, dual-gate records, screenshot set.

## Blocking Failures

Any criterion blocking failure; any English copy string in new surfaces (Q5) blocks regardless.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned.
- Next owner: QA Validator (after F-01..F-03 accepted).
- Next action: execute criteria at fixed SHA.
- Required files/evidence: as above.
- Blockers or open decisions: none.
