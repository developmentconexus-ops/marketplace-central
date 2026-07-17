# M-02 Plan Adjudication — chip rulings on planner findings

Planner: gpt-5.6-sol/medium (ledger row 1; full plan `plan-batch.md`). 12 blocking findings
triaged 2026-07-16 by CHIP-M02 against repository truth order (profile §8). Rulings below are
chip-scope (milestone-internal sequencing, reversible bindings, truth-order resolutions);
items needing hub action are marked HUB.

| # | Finding | Ruling |
| --- | --- | --- |
| 1 | IC-05 lacks installation query key | ADD additive `installations: ["installations"]` namespace + `installationsQueryKeys.list()`; allowed by IC-05 Compatibility Rules ("new namespaces/key builders extend"). Submitted for IC-05 ratification at CLOSED. |
| 2 | F-01 needs EmptyState/ErrorState owned by F-02 | Slice-level resequence, single writer unchanged: F02-S6-discovery → F02-S1..S5 → F01-S1..S3 → F02-S6-proxy → F03-S1..S5. Feature EVIDENCE still closes per feature; one-writer rationale of the F-01→F-02→F-03 order is preserved (same sole writer). |
| 3 | M02-C03 vs legacy pages' own installation fetch | Satisfiable: mission ADR-15 accepted trade-off + `## Migration Briefs` explicitly keep legacy direct fetch until rebuild. C03 is proven over new-code routes (context + stubs + /anuncios); legacy pages mounting at new paths keep their fetch, recorded debt not defect. |
| 4 | Stubs vs legacy-at-new-paths contradiction | Legacy pages wrap-mounted UNMODIFIED at their new paths (F-01 brief + milestone RK-05 mitigation); "em construção" stub only where no legacy page exists (`/catalogo/produtos/:id`, `/protocolos/:id`). pt-BR blocking rule binds new/rebuilt surfaces only (IC-05 glossary wording). |
| 5 | Undeclared workspace deps (apps/web→ui; ui→web-query) + package-lock | HUB — REQUEST event sent (dep change gate, profile §9). Blocks F02-S5 and F-01; F02-S1..S4 + discovery slice proceed meanwhile. |
| 6 | 12 failureCopy pt-BR literals unauthored | Chip authors strings from IC-03 taxonomy semantics. VC M02-C07 requires only non-empty pt-BR per code + byte-exact fallback `Falha desconhecida ({code})`. Literals submitted for IC-05 ratification at CLOSED. |
| 7 | Enrichment discriminator literal undefined | Pin `product_enrichment` (snake_case, consistent with IC-03 types). VC row read semantically. Flagged for IC-05 ratification. |
| 8 | `/orders` proxy vs SPA redirect collision | Proxy row WITH HTML bypass (requests with `Accept: text/html` fall through to SPA → redirect; API requests proxied). Both behaviors tested. Satisfies IC-05 "proxy list must include /orders" and M02-C01. |
| 9 | IC-02 research doc stale vs merged M-01 OpenAPI/SDK | Truth order (profile §8): OpenAPI+SDK > .mnfs research. Consume current SDK shapes (`below_margin_worst_case`, nullable summary counters, extended ListingStatus). IC-02 doc amendment flagged to hub/strategist — chip does not edit contracts. |
| 10 | Flat vs grouped 2a table | F-03 feature brief is binding at feature grain: flat `listListings` table; by-product view absent from Expected Output → out of M-02 scope. |
| 11 | `invalid_filter` lacks offending-key detail | Retry clears ALL `filter.*` params (conservative superset of "clears offending filter"); `q` and `tab` preserved. |
| 12 | Refresh-run polling key unspecified | Use IC-05's existing `syncQueryKeys.runs(installationId, filters)` over `listIntegrationOperationRuns`; `refetchInterval` 2s while run non-terminal (mirrors IC-03 protocolo polling), stop at terminal. Flagged for IC-05 ratification. |

Hub directive 2026-07-16 (post-plan): proxy rows `/dashboard` + `/sync` reassigned from
M-05 F-01 to M-02 (CHIP-SAT is zero-frontend in the W1 replan; hub owns seam assignment).
Folded into F02-S6b — its proxy row set is now `/listings /mutations /market /orders
/profitability /dashboard /sync`. Supersedes IC-05 Transport writer-sequence note. ACKed.

Execution slice order (ruling 2): S0=F02-S6a vitest discovery only · S1..S3 = F02-S1..S3 ·
S4 = F02-S4 · [dep REQUEST resolved] S5 = F02-S5 · S6..S8 = F01-S1..S3 · S9 = F02-S6b
(tailwind @source + proxy rows + /orders bypass) · S10..S14 = F03-S1..S5.
